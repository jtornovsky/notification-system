package handlers

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"delivery-service/internal/models"
)

// DeliveryConfig holds configuration for simulating a delivery type
type DeliveryConfig struct {
	Name         string  // "Email", "SMS", "Push"
	MinDelayMs   int     // Minimum delivery time
	MaxDelayMs   int     // Maximum delivery time
	FailureRate  float32 // 0.10 = 10% failure rate
	ErrorMessage string  // Error message when fails
	SuccessEmoji string  // Emoji for success logs
}

var (
	emailConfig = DeliveryConfig{
		Name:         "Email",
		MinDelayMs:   50,
		MaxDelayMs:   200,
		FailureRate:  0.10,
		ErrorMessage: "SMTP connection timeout",
		SuccessEmoji: "📧",
	}

	smsConfig = DeliveryConfig{
		Name:         "SMS",
		MinDelayMs:   30,
		MaxDelayMs:   100,
		FailureRate:  0.10,
		ErrorMessage: "carrier gateway unreachable",
		SuccessEmoji: "📱",
	}

	pushConfig = DeliveryConfig{
		Name:         "Push",
		MinDelayMs:   20,
		MaxDelayMs:   80,
		FailureRate:  0.10,
		ErrorMessage: "device token invalid",
		SuccessEmoji: "🔔",
	}
)

// simulateDelivery is the generic delivery simulation function
func simulateDelivery(notification models.Notification, config DeliveryConfig) (string, int64, error) {
	startTime := time.Now()

	// Simulate network delay
	delayRange := config.MaxDelayMs - config.MinDelayMs
	delay := time.Duration(config.MinDelayMs+rand.Intn(delayRange)) * time.Millisecond
	time.Sleep(delay)

	deliveryTime := time.Since(startTime).Milliseconds()

	// Random failure
	if rand.Float32() < config.FailureRate {
		return "FAILED", deliveryTime, fmt.Errorf(config.ErrorMessage)
	}

	// Success
	log.Printf("%s %s sent to %s (took %dms)\n", config.SuccessEmoji, config.Name, notification.Recipient, deliveryTime)
	return "SENT", deliveryTime, nil
}

// SimulateEmailDelivery simulates sending an email with random failures
func SimulateEmailDelivery(notification models.Notification) (string, int64, error) {
	return simulateDelivery(notification, emailConfig)
}

// SimulateSmsDelivery simulates sending an SMS with random failures
func SimulateSmsDelivery(notification models.Notification) (string, int64, error) {
	return simulateDelivery(notification, smsConfig)
}

// SimulatePushDelivery simulates sending a push notification with random failures
func SimulatePushDelivery(notification models.Notification) (string, int64, error) {
	return simulateDelivery(notification, pushConfig)
}
