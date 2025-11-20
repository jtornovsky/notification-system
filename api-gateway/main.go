package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/segmentio/kafka-go"
)

// Request structures
type CreateNotificationRequest struct {
	UserID    string            `json:"user_id" binding:"required"`
	Type      string            `json:"type" binding:"required,oneof=EMAIL SMS PUSH"`
	Recipient string            `json:"recipient" binding:"required"`
	Subject   string            `json:"subject"`
	Message   string            `json:"message" binding:"required"`
	Metadata  map[string]string `json:"metadata"`
}

type NotificationResponse struct {
	ID        string            `json:"id"`
	UserID    string            `json:"user_id"`
	Type      string            `json:"type"`
	Recipient string            `json:"recipient"`
	Subject   string            `json:"subject"`
	Message   string            `json:"message"`
	Status    string            `json:"status"`
	CreatedAt int64             `json:"created_at"`
	Metadata  map[string]string `json:"metadata"`
}

var (
	redisClient *redis.Client
	kafkaWriter *kafka.Writer
	ctx         = context.Background()
)

func main() {
	// Initialize Redis
	redisClient = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   0,
	})

	// Test Redis connection
	if err := redisClient.Ping(ctx).Err(); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	log.Println("✓ Connected to Redis")

	// Initialize Kafka Writer
	kafkaWriter = &kafka.Writer{
		Addr:     kafka.TCP("localhost:9092"),
		Topic:    "notifications",
		Balancer: &kafka.LeastBytes{},
	}
	defer func(kafkaWriter *kafka.Writer) {
		if err := kafkaWriter.Close(); err != nil {
			log.Fatalf("Error closing Kafka: %v", err)
		}
	}(kafkaWriter)
	log.Println("✓ Connected to Kafka")

	// Initialize Gin router
	router := gin.Default()

	// Routes
	router.POST("/notifications", createNotification)
	router.GET("/notifications/:id", getNotification)
	router.GET("/health", healthCheck)

	// Start server
	log.Println("🚀 API Gateway starting on :8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func createNotification(c *gin.Context) {
	var req CreateNotificationRequest

	// Validate request
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("❌ Validation error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create notification
	notification := NotificationResponse{
		ID:        uuid.New().String(),
		UserID:    req.UserID,
		Type:      req.Type,
		Recipient: req.Recipient,
		Subject:   req.Subject,
		Message:   req.Message,
		Status:    "PENDING",
		CreatedAt: time.Now().UnixMilli(),
		Metadata:  req.Metadata,
	}

	log.Printf("📝 Creating notification: ID=%s, Type=%s, Recipient=%s",
		notification.ID, notification.Type, notification.Recipient)

	// Cache in Redis
	notificationJSON, _ := json.Marshal(notification)
	if err := redisClient.Set(ctx, "notification:"+notification.ID, notificationJSON, 1*time.Hour).Err(); err != nil {
		log.Printf("⚠️  Redis cache failed: %v", err)
	} else {
		log.Printf("✓ Cached in Redis: notification:%s", notification.ID)
	}

	// Publish to Kafka
	kafkaMessage := kafka.Message{
		Key:   []byte(notification.ID),
		Value: notificationJSON,
	}

	if err := kafkaWriter.WriteMessages(ctx, kafkaMessage); err != nil {
		log.Printf("❌ Failed to publish to Kafka: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to queue notification"})
		return
	}

	log.Printf("✓ Published to Kafka topic 'notifications': ID=%s", notification.ID)

	c.JSON(http.StatusCreated, notification)
}

func getNotification(c *gin.Context) {
	id := c.Param("id")

	log.Printf("🔍 Looking up notification: %s", id)

	// Try Redis first
	val, err := redisClient.Get(ctx, "notification:"+id).Result()
	if err == redis.Nil {
		log.Printf("❌ Notification not found: %s", id)
		c.JSON(http.StatusNotFound, gin.H{"error": "Notification not found"})
		return
	} else if err != nil {
		log.Printf("⚠️  Redis error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal error"})
		return
	}

	var notification NotificationResponse
	err = json.Unmarshal([]byte(val), &notification)
	if err != nil {
		return
	}

	log.Printf("✓ Found in Redis: ID=%s, Status=%s", notification.ID, notification.Status)
	c.JSON(http.StatusOK, notification)
}

func healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "healthy",
		"service":   "api-gateway",
		"timestamp": time.Now().Unix(),
	})
}
