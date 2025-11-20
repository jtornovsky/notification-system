package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"delivery-service/internal/handlers"
	"delivery-service/internal/mongo"
)

func main() {
	log.Println("🚀 Starting Delivery Service...")

	// Configuration
	kafkaBrokers := []string{"localhost:9092"}
	mongoURI := "mongodb://admin:password@localhost:27017"
	mongoDatabase := "notifications"
	mongoCollection := "delivery_results"

	// Initialize MongoDB client
	mongoClient, err := mongo.NewClient(mongoURI, mongoDatabase, mongoCollection)
	if err != nil {
		log.Fatalf("❌ Failed to connect to MongoDB: %v", err)
	}
	defer func() {
		if err := mongoClient.Close(); err != nil {
			log.Printf("⚠️ Error closing MongoDB: %v\n", err)
		}
	}()

	// Create handlers
	emailHandler, err := handlers.NewHandler(
		"Email",
		kafkaBrokers,
		"email-notifications",
		"delivery-service-email",
		mongoClient,
		handlers.SimulateEmailDelivery,
	)
	if err != nil {
		log.Fatalf("❌ Failed to create email handler: %v", err)
	}

	smsHandler, err := handlers.NewHandler(
		"SMS",
		kafkaBrokers,
		"sms-notifications",
		"delivery-service-sms",
		mongoClient,
		handlers.SimulateSmsDelivery,
	)
	if err != nil {
		log.Fatalf("❌ Failed to create SMS handler: %v", err)
	}

	pushHandler, err := handlers.NewHandler(
		"Push",
		kafkaBrokers,
		"push-notifications",
		"delivery-service-push",
		mongoClient,
		handlers.SimulatePushDelivery,
	)
	if err != nil {
		log.Fatalf("❌ Failed to create push handler: %v", err)
	}

	// Start all handlers in separate goroutines
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup
	allHandlers := []*handlers.Handler{emailHandler, smsHandler, pushHandler}

	for _, handler := range allHandlers {
		wg.Add(1)
		go func(h *handlers.Handler) {
			defer wg.Done()
			if err := h.Start(ctx); err != nil {
				log.Printf("❌ Handler error: %v", err)
			}
		}(handler)
	}

	// Give goroutines time to start and log
	time.Sleep(100 * time.Millisecond)

	log.Println("✅ All handlers started successfully")

	// Wait for interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	<-sigChan
	log.Println("\n📭 Shutdown signal received...")

	// Graceful shutdown
	cancel() // Signal context cancellation

	// Close all handlers
	for _, handler := range allHandlers {
		if err := handler.Close(); err != nil {
			log.Printf("⚠️ Error closing handler: %v\n", err)
		}
	}

	// Wait for all goroutines to finish
	wg.Wait()

	log.Println("✅ Delivery Service stopped")
}
