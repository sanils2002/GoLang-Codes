package examples

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"kafka/config"
	"kafka/consumer"
	"kafka/producer"
	"kafka/retry"
)

// ExampleMessage represents a sample message structure
type ExampleMessage struct {
	ID      string                 `json:"id"`
	Type    string                 `json:"type"`
	Data    map[string]interface{} `json:"data"`
	Retry   int                    `json:"retry_count"`
	Created time.Time              `json:"created_at"`
}

// ProducerExample shows how to create a producer with retry
func ProducerExample() {
	// Create a message
	message := ExampleMessage{
		ID:      "msg-123",
		Type:    "notification",
		Data:    map[string]interface{}{"text": "Hello World"},
		Retry:   0,
		Created: time.Now(),
	}

	// Marshal to JSON
	data, err := json.Marshal(message)
	if err != nil {
		log.Printf("Failed to marshal message: %v", err)
		return
	}

	// Create context
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	contextStruct := &config.ContextStruct{
		EchoContext: ctx,
	}

	// Create Kafka writer
	writer := producer.KafkaWriter("example-topic", 1, 10*time.Millisecond)
	defer writer.Close()

	// Produce message with retry
	err = producer.KafkaProducerWithRetry(contextStruct, writer, data, retry.DefaultRetryConfig())
	if err != nil {
		log.Printf("Producer failed: %v", err)
		return
	}

	log.Printf("Producer example completed successfully")
}

// ConsumerExample shows how to handle consumer errors with retry
func ConsumerExample() {
	// Create reader
	reader := consumer.KafkaRead("example-topic", "example-group", 10e3, 10e6)
	defer reader.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Retry configuration for consumer operations
	retryConfig := retry.RetryConfig{
		MaxRetries:        3,
		InitialBackoff:    100 * time.Millisecond,
		MaxBackoff:        5 * time.Second,
		BackoffMultiplier: 2.0,
		JitterFactor:      0.1,
		Timeout:           30 * time.Second,
	}

	// Read message with retry
	err := retry.RetryWithBackoff(ctx, func() error {
		_, err := reader.ReadMessage(context.Background())
		return err
	}, retryConfig)

	if err != nil {
		log.Printf("Consumer failed: %v", err)
		return
	}

	log.Printf("Consumer example completed successfully")
}
