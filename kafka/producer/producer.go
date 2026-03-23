package producer

import (
	"context"
	"crypto/tls"
	"errors"
	"log"
	"strconv"
	"time"

	"kafka/config"
	"kafka/retry"

	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/plain"
)

// KafkaProducer produces a message with retry logic
func KafkaProducer(ctx context.Context, w *kafka.Writer, Data []byte) error {
	if ctx == nil {
		return errors.New("context cannot be nil")
	}

	log.Printf("Starting Kafka producer operation with retry - Message size: %d bytes", len(Data))

	// Use retry mechanism for better reliability
	retryConfig := retry.DefaultRetryConfig()
	err := retry.RetryWithBackoff(ctx, func() error {
		return w.WriteMessages(ctx, kafka.Message{Value: Data})
	}, retryConfig)

	if err != nil {
		log.Printf("Kafka producer operation failed after retries: %v", err)
		return err
	}

	log.Printf("Kafka producer operation completed successfully")
	return nil
}

// KafkaProducerWithRetry produces a message with custom retry configuration
func KafkaProducerWithRetry(ctx *config.ContextStruct, w *kafka.Writer, data []byte, retryConfig retry.RetryConfig) error {
	// Protect against nil context
	if ctx == nil || ctx.EchoContext == nil {
		log.Printf("Invalid context provided to KafkaProducerWithRetry")
		return errors.New("invalid context: context cannot be nil")
	}

	log.Printf("Producing Kafka message with custom retry configuration - Max retries: %d, Initial backoff: %v", retryConfig.MaxRetries, retryConfig.InitialBackoff)

	operation := func() error {
		return w.WriteMessages(ctx.EchoContext, kafka.Message{Value: data})
	}

	return retry.RetryWithBackoff(ctx.EchoContext, operation, retryConfig)
}

// KafkaWriter creates a basic Kafka writer
func KafkaWriter(topic string, batchSize int, batchTimeout time.Duration) *kafka.Writer {
	w := &kafka.Writer{
		Addr:                   kafka.TCP(config.AppConfig.Kafka.Host + ":" + strconv.Itoa(config.AppConfig.Kafka.Port)),
		Topic:                  topic,
		AllowAutoTopicCreation: true,
		Transport: &kafka.Transport{
			SASL: plain.Mechanism{
				Username: config.AppConfig.Kafka.Username,
				Password: config.AppConfig.Kafka.Password,
			},
			TLS: &tls.Config{
				ClientAuth: 1,
			},
		},
		BatchSize:    batchSize,
		BatchTimeout: batchTimeout,
		// Add retry-friendly configuration
		RequiredAcks: kafka.RequireAll,
		Async:        false, // Ensure synchronous writes for better error handling
	}

	return w
}

// KafkaWriterWithRetry creates a Kafka writer with retry configuration
func KafkaWriterWithRetry(topic string, batchSize int, batchTimeout time.Duration, retryConfig retry.RetryConfig) *kafka.Writer {
	log.Printf("Creating Kafka writer with retry configuration - Topic: %s, BatchSize: %d, BatchTimeout: %v", topic, batchSize, batchTimeout)

	w := &kafka.Writer{
		Addr:                   kafka.TCP(config.AppConfig.Kafka.Host + ":" + strconv.Itoa(config.AppConfig.Kafka.Port)),
		Topic:                  topic,
		AllowAutoTopicCreation: true,
		Transport: &kafka.Transport{
			SASL: plain.Mechanism{
				Username: config.AppConfig.Kafka.Username,
				Password: config.AppConfig.Kafka.Password,
			},
			TLS: &tls.Config{
				ClientAuth: 1,
			},
		},
		BatchSize:    batchSize,
		BatchTimeout: batchTimeout,
		// Add retry configuration to the writer
		RequiredAcks: kafka.RequireAll,
		Async:        false, // Ensure synchronous writes for better error handling
	}

	log.Printf("Kafka writer created successfully - Topic: %s", topic)
	return w
}

// KafkaWriterWithDefaultRetry creates a Kafka writer with default retry configuration
func KafkaWriterWithDefaultRetry(topic string, batchSize int, batchTimeout time.Duration) *kafka.Writer {
	return KafkaWriterWithRetry(topic, batchSize, batchTimeout, retry.DefaultRetryConfig())
}
