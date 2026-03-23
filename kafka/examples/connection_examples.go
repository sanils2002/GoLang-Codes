package examples

import (
	"context"
	"log"
	"strconv"
	"time"

	"kafka/config"
	"kafka/consumer"
	"kafka/producer"
	"kafka/retry"
	"kafka/utils"
)

// ConnectionRetryExample demonstrates how to use retry-enabled connection functions
func ConnectionRetryExample() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	address := config.AppConfig.Kafka.Host + ":" + strconv.Itoa(int(config.AppConfig.Kafka.Port))

	// Example 1: Basic connection with retry
	conn, err := utils.KafkaDialWithRetry(ctx, address)
	if err != nil {
		log.Printf("Error: Failed to establish basic Kafka connection: %v", err)
		return
	}
	defer conn.Close()

	// Example 2: Leader connection with retry
	leaderConn, err := utils.KafkaDialLeaderWithRetry(ctx, address, "my-topic", 0)
	if err != nil {
		log.Printf("Error: Failed to establish leader connection: %v", err)
		return
	}
	defer leaderConn.Close()

	// Example 3: Using retry-enabled writer
	writer := producer.KafkaWriterWithDefaultRetry("my-topic", 1, 10*time.Millisecond)
	defer writer.Close()

	// Example 4: Using retry-enabled reader
	reader := consumer.KafkaReadWithRetry("my-topic", "my-group", 10e3, 10e6)
	defer reader.Close()

	log.Printf("Info: All retry-enabled connections established successfully")
}

// DifferentUseCasesExample demonstrates different retry configurations for different use cases
func DifferentUseCasesExample() {
	// Example 1: User notification (high priority)
	userNotificationConfig := retry.RetryConfig{
		MaxRetries:        5,
		InitialBackoff:    50 * time.Millisecond,
		MaxBackoff:        2 * time.Second,
		BackoffMultiplier: 1.5,
		JitterFactor:      0.1,
		Timeout:           60 * time.Second,
	}

	// Example 2: Analytics data (low priority)
	analyticsConfig := retry.RetryConfig{
		MaxRetries:        2,
		InitialBackoff:    200 * time.Millisecond,
		MaxBackoff:        10 * time.Second,
		BackoffMultiplier: 2.0,
		JitterFactor:      0.2,
		Timeout:           30 * time.Second,
	}

	// Example 3: System health check (critical)
	healthCheckConfig := retry.RetryConfig{
		MaxRetries:        3,
		InitialBackoff:    100 * time.Millisecond,
		MaxBackoff:        1 * time.Second,
		BackoffMultiplier: 1.2,
		JitterFactor:      0.05,
		Timeout:           10 * time.Second,
	}

	// Use different configurations based on message type
	messageType := "user_notification"
	var config retry.RetryConfig

	switch messageType {
	case "user_notification":
		config = userNotificationConfig
	case "analytics":
		config = analyticsConfig
	case "health_check":
		config = healthCheckConfig
	default:
		config = retry.DefaultRetryConfig()
	}

	// Use the appropriate configuration
	ctx, cancel := context.WithTimeout(context.Background(), config.Timeout)
	defer cancel()

	err := retry.RetryWithBackoff(ctx, func() error {
		// Your Kafka operation here
		return nil
	}, config)

	if err != nil {
		log.Printf("Different use cases example failed: %v", err)
		return
	}

	log.Printf("Different use cases example completed successfully")
}
