package examples

import (
	"context"
	"log"
	"time"

	"kafka/producer"
	"kafka/retry"

	"github.com/segmentio/kafka-go"
)

// BasicRetryExample demonstrates basic retry functionality
func BasicRetryExample() {
	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Use default retry configuration
	retryConfig := retry.DefaultRetryConfig()

	// Example operation that might fail
	operation := func() error {
		// Simulate a Kafka write operation
		// In real usage, this would be your actual Kafka writer
		return nil // or some error
	}

	// Execute with retry
	err := retry.RetryWithBackoff(ctx, operation, retryConfig)
	if err != nil {
		log.Printf("Basic retry failed: %v", err)
		return
	}

	log.Printf("Basic retry example completed successfully")
}

// HighPriorityRetryExample demonstrates high-priority message retry
func HighPriorityRetryExample() {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	data := []byte(`{"priority": "high", "message": "Urgent notification"}`)

	// Custom retry configuration for high-priority messages
	highPriorityConfig := retry.RetryConfig{
		MaxRetries:        5,
		InitialBackoff:    50 * time.Millisecond,
		MaxBackoff:        2 * time.Second,
		BackoffMultiplier: 1.5,
		JitterFactor:      0.1,
		Timeout:           60 * time.Second,
	}

	// Example: Using with a Kafka writer
	writer := producer.KafkaWriter("high-priority-topic", 1, 10*time.Millisecond)
	defer writer.Close()

	err := retry.RetryWithBackoff(ctx, func() error {
		return writer.WriteMessages(ctx, kafka.Message{Value: data})
	}, highPriorityConfig)

	if err != nil {
		log.Printf("High priority retry failed: %v", err)
		return
	}

	log.Printf("High priority retry example completed successfully")
}
