package handlers

import (
	"context"
	"log"
	"time"

	"delivery-service/internal/models"
	"delivery-service/internal/mongo"

	pb "github.com/jtornovsky/notification-system/proto"
	"github.com/segmentio/kafka-go"
	"google.golang.org/protobuf/proto"
)

// DeliverySimulator is a function that simulates delivery for a specific type
type DeliverySimulator func(models.Notification) (string, int64, error)

// Handler is a generic notification delivery handler
type Handler struct {
	name         string
	consumer     *kafka.Reader
	producer     *kafka.Writer
	mongoClient  *mongo.Client
	simulator    DeliverySimulator
	shutdownChan chan struct{}
}

// NewHandler creates a new generic delivery handler
func NewHandler(
	name string,
	brokers []string,
	topic string,
	groupID string,
	mongoClient *mongo.Client,
	simulator DeliverySimulator,
) (*Handler, error) {
	consumer := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     brokers,
		Topic:       topic,
		GroupID:     groupID,
		MinBytes:    10e3,
		MaxBytes:    10e6,
		StartOffset: kafka.LastOffset,
	})

	producer := &kafka.Writer{
		Addr:     kafka.TCP(brokers...),
		Topic:    "delivery-events",
		Balancer: &kafka.LeastBytes{},
	}

	log.Printf("✓ %s handler initialized for topic: %s\n", name, topic)

	return &Handler{
		name:         name,
		consumer:     consumer,
		producer:     producer,
		mongoClient:  mongoClient,
		simulator:    simulator,
		shutdownChan: make(chan struct{}),
	}, nil
}

// processMessage handles a single notification message
func (h *Handler) processMessage(ctx context.Context, message kafka.Message) error {
	// Deserialize Protobuf message
	pbNotification := &pb.Notification{}
	if err := proto.Unmarshal(message.Value, pbNotification); err != nil {
		log.Printf("❌ [%s] Failed to unmarshal protobuf: %v\n", h.name, err)
		return err
	}

	log.Printf("📦 [%s] Received Protobuf message: %d bytes\n", h.name, len(message.Value))

	// Convert to internal model
	notification := models.Notification{
		ID:        pbNotification.Id,
		Type:      pb.NotificationType_name[int32(pbNotification.Type)],
		Recipient: pbNotification.Recipient,
		Subject:   pbNotification.Subject,
		Message:   pbNotification.Message,
	}

	status, deliveryTimeMs, err := h.simulator(notification)

	errorMessage := ""
	if err != nil {
		errorMessage = err.Error()
		log.Printf("❌ [%s] Delivery failed: %v\n", h.name, err)
	}

	deliveryResult := models.DeliveryResult{
		NotificationID: notification.ID,
		Type:           notification.Type,
		Recipient:      notification.Recipient,
		Status:         status,
		Timestamp:      time.Now(),
		DeliveryTimeMs: deliveryTimeMs,
		ErrorMessage:   errorMessage,
	}

	if err := h.mongoClient.SaveDeliveryResult(ctx, deliveryResult); err != nil {
		log.Printf("❌ [%s] Failed to save to MongoDB: %v\n", h.name, err)
		return err
	}

	log.Printf("✓ [%s] Saved delivery result to MongoDB\n", h.name)

	// Create Protobuf delivery event
	pbDeliveryEvent := &pb.DeliveryEvent{
		NotificationId: deliveryResult.NotificationID,
		Type:           pb.NotificationType(pb.NotificationType_value[deliveryResult.Type]),
		Recipient:      deliveryResult.Recipient,
		Status:         pb.NotificationStatus(pb.NotificationStatus_value[deliveryResult.Status]),
		ProcessedAt:    deliveryResult.Timestamp.UnixMilli(),
		ErrorMessage:   deliveryResult.ErrorMessage,
		DeliveryTimeMs: int32(deliveryResult.DeliveryTimeMs),
	}

	// Marshal to Protobuf bytes
	eventBytes, err := proto.Marshal(pbDeliveryEvent)
	if err != nil {
		log.Printf("❌ [%s] Failed to marshal protobuf delivery event: %v\n", h.name, err)
		return err
	}

	log.Printf("📦 [%s] Publishing Protobuf delivery event: %d bytes\n", h.name, len(eventBytes))

	if err := h.producer.WriteMessages(ctx, kafka.Message{
		Key:   []byte(notification.ID),
		Value: eventBytes,
	}); err != nil {
		log.Printf("❌ [%s] Failed to publish delivery event: %v\n", h.name, err)
		return err
	}

	log.Printf("✓ [%s] Published delivery event to Kafka\n", h.name)
	return nil
}

// Start begins consuming messages
func (h *Handler) Start(ctx context.Context) error {
	log.Printf("🚀 [%s] Handler started, waiting for messages...\n", h.name)

	for {
		select {
		case <-h.shutdownChan:
			log.Printf("📭 [%s] Handler shutting down...\n", h.name)
			return nil
		default:
			message, err := h.consumer.FetchMessage(ctx)
			if err != nil {
				// Don't log if context was canceled (shutdown in progress)
				if ctx.Err() == nil {
					log.Printf("⚠️ [%s] Error fetching message: %v\n", h.name, err)
				}
				continue
			}

			if err := h.processMessage(ctx, message); err != nil {
				log.Printf("❌ [%s] Error processing message: %v\n", h.name, err)
			}

			if err := h.consumer.CommitMessages(ctx, message); err != nil {
				log.Printf("⚠️ [%s] Failed to commit message: %v\n", h.name, err)
			}
		}
	}
}

// Close gracefully shuts down the handler
func (h *Handler) Close() error {
	close(h.shutdownChan)

	if err := h.consumer.Close(); err != nil {
		log.Printf("⚠️ [%s] Error closing consumer: %v\n", h.name, err)
	}

	if err := h.producer.Close(); err != nil {
		log.Printf("⚠️ [%s] Error closing producer: %v\n", h.name, err)
	}

	log.Printf("✓ [%s] Handler closed\n", h.name)
	return nil
}
