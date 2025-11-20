package models

import "time"

// Notification represents the incoming notification message
type Notification struct {
	ID        string `json:"id"`
	Type      string `json:"type"`
	Recipient string `json:"recipient"`
	Message   string `json:"message"`
	Subject   string `json:"subject,omitempty"`
}

// DeliveryResult represents a notification delivery attempt
// Used for both MongoDB storage (bson tags) and Kafka events (json tags)
type DeliveryResult struct {
	NotificationID string    `json:"notification_id" bson:"notification_id"`
	Type           string    `json:"type" bson:"type"`
	Recipient      string    `json:"recipient" bson:"recipient"`
	Status         string    `json:"status" bson:"status"`
	Timestamp      time.Time `json:"timestamp" bson:"timestamp"`
	DeliveryTimeMs int64     `json:"delivery_time_ms" bson:"delivery_time_ms"`
	ErrorMessage   string    `json:"error_message,omitempty" bson:"error_message,omitempty"`
}
