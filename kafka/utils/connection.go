package utils

import (
	"context"
	"log"
	"strconv"

	"kafka/config"
	"kafka/retry"

	"github.com/segmentio/kafka-go"
)

// KafkaConnectionWithRetry creates a Kafka connection with retry logic
func KafkaConnectionWithRetry(ctx context.Context, connectionType string, address string, topic string, partition int) (*kafka.Conn, error) {
	log.Printf("Creating Kafka connection with retry - Type: %s, Address: %s, Topic: %s, Partition: %d", connectionType, address, topic, partition)

	// Use retry for connection
	retryConfig := retry.DefaultRetryConfig()
	var conn *kafka.Conn
	err := retry.RetryWithBackoff(ctx, func() error {
		var dialErr error
		switch connectionType {
		case "leader":
			conn, dialErr = kafka.DialLeader(ctx, "tcp", address, topic, partition)
		default:
			conn, dialErr = kafka.Dial("tcp", address)
		}
		return dialErr
	}, retryConfig)

	if err != nil {
		log.Printf("Failed to establish Kafka connection after retries - Type: %s, Address: %s, Error: %v", connectionType, address, err)
		return nil, err
	}

	log.Printf("Kafka connection established successfully - Type: %s, Address: %s, Topic: %s", connectionType, address, topic)
	return conn, nil
}

// KafkaDialLeaderWithRetry creates a Kafka leader connection with retry logic
func KafkaDialLeaderWithRetry(ctx context.Context, address string, topic string, partition int) (*kafka.Conn, error) {
	return KafkaConnectionWithRetry(ctx, "leader", address, topic, partition)
}

// KafkaDialWithRetry creates a basic Kafka connection with retry logic
func KafkaDialWithRetry(ctx context.Context, address string) (*kafka.Conn, error) {
	return KafkaConnectionWithRetry(ctx, "basic", address, "", 0)
}

// GetKafkaAddress returns the formatted Kafka address
func GetKafkaAddress() string {
	return config.AppConfig.Kafka.Host + ":" + strconv.Itoa(config.AppConfig.Kafka.Port)
}
