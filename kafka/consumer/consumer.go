package consumer

import (
	"crypto/tls"
	"log"
	"strconv"
	"time"

	"kafka/config"

	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/plain"
)

// KafkaRead creates a basic Kafka reader
func KafkaRead(topic string, consumerGroupId string, minBytes int, maxBytes int) *kafka.Reader {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{config.AppConfig.Kafka.Host + ":" + strconv.Itoa(config.AppConfig.Kafka.Port)},
		GroupID:  consumerGroupId,
		Topic:    topic,
		MaxBytes: maxBytes, // 10MB
		Dialer: &kafka.Dialer{
			Timeout:   10 * time.Second,
			DualStack: true,
			SASLMechanism: plain.Mechanism{
				Username: config.AppConfig.Kafka.Username,
				Password: config.AppConfig.Kafka.Password,
			},
			TLS: &tls.Config{
				ClientAuth: 1,
			},
		},
	})

	return r
}

// KafkaReadWithRetry creates a Kafka reader with retry logic for connection failures
func KafkaReadWithRetry(topic string, consumerGroupId string, minBytes int, maxBytes int) *kafka.Reader {
	log.Printf("Creating Kafka reader with retry logic - Topic: %s, Group: %s", topic, consumerGroupId)

	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{config.AppConfig.Kafka.Host + ":" + strconv.Itoa(config.AppConfig.Kafka.Port)},
		GroupID:  consumerGroupId,
		Topic:    topic,
		MaxBytes: maxBytes, // 10MB
		Dialer: &kafka.Dialer{
			Timeout:   10 * time.Second,
			DualStack: true,
			SASLMechanism: plain.Mechanism{
				Username: config.AppConfig.Kafka.Username,
				Password: config.AppConfig.Kafka.Password,
			},
			TLS: &tls.Config{
				ClientAuth: 1,
			},
		},
	})

	log.Printf("Kafka reader created successfully - Topic: %s, Group: %s", topic, consumerGroupId)
	return r
}
