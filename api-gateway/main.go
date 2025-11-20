package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/segmentio/kafka-go"
)

var (
	redisClient *redis.Client
	kafkaWriter *kafka.Writer
	ctx         = context.Background()
)

type NotificationRequest struct {
	UserID    string `json:"user_id"`
	Type      string `json:"type"`
	Recipient string `json:"recipient"`
	Subject   string `json:"subject"`
	Message   string `json:"message"`
}

type NotificationResponse struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Type      string    `json:"type"`
	Recipient string    `json:"recipient"`
	Subject   string    `json:"subject"`
	Message   string    `json:"message"`
	CreatedAt time.Time `json:"created_at"`
}

func main() {
	redisClient = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	kafkaWriter = &kafka.Writer{
		Addr:     kafka.TCP("localhost:9092"),
		Topic:    "notifications",
		Balancer: &kafka.LeastBytes{},
	}
	defer func(kafkaWriter *kafka.Writer) {
		err := kafkaWriter.Close()
		if err != nil {

		}
	}(kafkaWriter)

	router := gin.Default()

	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusOK)
			return
		}

		c.Next()
	})

	router.POST("/notifications", createNotification)
	router.GET("/notifications/:id", getNotification)
	router.GET("/health", healthCheck)

	log.Println("API Gateway listening on :8080")
	log.Fatal(router.Run(":8080"))
}

func createNotification(c *gin.Context) {
	var req NotificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	notification := NotificationResponse{
		ID:        generateID(),
		UserID:    req.UserID,
		Type:      req.Type,
		Recipient: req.Recipient,
		Subject:   req.Subject,
		Message:   req.Message,
		CreatedAt: time.Now(),
	}

	notificationJSON, _ := json.Marshal(notification)
	redisClient.Set(ctx, "notification:"+notification.ID, notificationJSON, time.Hour)

	err := kafkaWriter.WriteMessages(ctx, kafka.Message{
		Key:   []byte(notification.ID),
		Value: notificationJSON,
	})

	if err != nil {
		log.Printf("Failed to write to Kafka: %v", err)
	}

	c.JSON(http.StatusCreated, notification)
}

func getNotification(c *gin.Context) {
	id := c.Param("id")

	val, err := redisClient.Get(ctx, "notification:"+id).Result()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Notification not found"})
		return
	}

	var notification NotificationResponse
	err = json.Unmarshal([]byte(val), &notification)
	if err != nil {
		return
	}

	c.JSON(http.StatusOK, notification)
}

func healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func generateID() string {
	return time.Now().Format("20060102150405")
}
