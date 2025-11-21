package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	pb "github.com/jtornovsky/notification-system/proto"
	"github.com/redis/go-redis/v9"
	"github.com/segmentio/kafka-go"
	"google.golang.org/protobuf/proto"
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
	router.GET("/analytics", getAnalytics)
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

	log.Printf("📝 Created notification: ID=%s, Type=%s, Recipient=%s",
		notification.ID, notification.Type, notification.Recipient)

	// Cache in Redis (still using JSON for easy retrieval via HTTP)
	notificationJSON, _ := json.Marshal(notification)
	redisClient.Set(ctx, "notification:"+notification.ID, notificationJSON, time.Hour)
	log.Printf("💾 Cached to Redis: notification:%s", notification.ID)

	// Create Protobuf message for Kafka
	pbNotification := &pb.Notification{
		Id:        notification.ID,
		UserId:    notification.UserID,
		Type:      pb.NotificationType(pb.NotificationType_value[notification.Type]),
		Recipient: notification.Recipient,
		Subject:   notification.Subject,
		Message:   notification.Message,
		CreatedAt: notification.CreatedAt.UnixMilli(),
		Status:    pb.NotificationStatus_PENDING,
	}

	// Marshal to Protobuf bytes
	pbData, err := proto.Marshal(pbNotification)
	if err != nil {
		log.Printf("❌ Failed to marshal protobuf: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process notification"})
		return
	}

	log.Printf("📦 Marshaled to Protobuf: %d bytes (vs JSON: %d bytes)",
		len(pbData), len(notificationJSON))

	// Send Protobuf bytes to Kafka
	err = kafkaWriter.WriteMessages(ctx, kafka.Message{
		Key:   []byte(notification.ID),
		Value: pbData, // ← Protobuf bytes instead of JSON
	})

	if err != nil {
		log.Printf("❌ Failed to write to Kafka: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send notification"})
		return
	}

	log.Printf("✅ Sent to Kafka topic 'notifications': ID=%s", notification.ID)

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

func getAnalytics(c *gin.Context) {
	// Proxy request to Elasticsearch
	resp, err := http.Post(
		"http://localhost:9200/notification-analytics/_search",
		"application/json",
		strings.NewReader(`{
			"size": 0,
			"aggs": {
				"by_status": {"terms": {"field": "status"}},
				"by_type": {"terms": {"field": "type"}}
			}
		}`),
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Printf("⚠️ Error to close body: %v", err)
		}
	}(resp.Body)

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		log.Printf("⚠️ Error to parse body: %v", err)
		return
	}

	c.JSON(http.StatusOK, result)
}
