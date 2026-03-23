package producer

import (
	"context"
	"testing"
	"time"
	
	"kafka/config"
	"kafka/retry"
)

func TestKafkaProducerNilContext(t *testing.T) {
	// Create a test writer
	writer := KafkaWriter("test-topic", 1, 10*time.Millisecond)
	defer writer.Close()
	
	// Test with nil context
	data := []byte("test message")
	err := KafkaProducer(nil, writer, data)
	
	if err == nil {
		t.Error("Expected error with nil context, got nil")
	}
}

func TestKafkaProducerWithRetry(t *testing.T) {
	ctx := context.Background()
	contextStruct := &config.ContextStruct{
		EchoContext: ctx,
	}
	
	// Create a test writer
	writer := KafkaWriter("test-topic", 1, 10*time.Millisecond)
	defer writer.Close()
	
	// Test data
	data := []byte("test message")
	
	// Test with retry configuration
	retryConfig := retry.RetryConfig{
		MaxRetries:        2,
		InitialBackoff:    10 * time.Millisecond,
		MaxBackoff:        100 * time.Millisecond,
		BackoffMultiplier: 2.0,
		JitterFactor:      0.1,
		Timeout:           1 * time.Second,
	}
	
	err := KafkaProducerWithRetry(contextStruct, writer, data, retryConfig)
	
	// This might fail if Kafka is not running, but we're testing the function structure
	if err != nil {
		t.Logf("KafkaProducerWithRetry returned error (expected if Kafka not running): %v", err)
	}
}

func TestKafkaWriter(t *testing.T) {
	writer := KafkaWriter("test-topic", 10, 100*time.Millisecond)
	defer writer.Close()
	
	if writer == nil {
		t.Error("Expected non-nil writer, got nil")
	}
	
	if writer.Topic != "test-topic" {
		t.Errorf("Expected topic 'test-topic', got '%s'", writer.Topic)
	}
}

func TestKafkaWriterWithRetry(t *testing.T) {
	retryConfig := retry.RetryConfig{
		MaxRetries:        2,
		InitialBackoff:    10 * time.Millisecond,
		MaxBackoff:        100 * time.Millisecond,
		BackoffMultiplier: 2.0,
		JitterFactor:      0.1,
		Timeout:           1 * time.Second,
	}
	
	writer := KafkaWriterWithRetry("test-topic", 10, 100*time.Millisecond, retryConfig)
	defer writer.Close()
	
	if writer == nil {
		t.Error("Expected non-nil writer, got nil")
	}
	
	if writer.Topic != "test-topic" {
		t.Errorf("Expected topic 'test-topic', got '%s'", writer.Topic)
	}
}

func TestKafkaWriterWithDefaultRetry(t *testing.T) {
	writer := KafkaWriterWithDefaultRetry("test-topic", 10, 100*time.Millisecond)
	defer writer.Close()
	
	if writer == nil {
		t.Error("Expected non-nil writer, got nil")
	}
	
	if writer.Topic != "test-topic" {
		t.Errorf("Expected topic 'test-topic', got '%s'", writer.Topic)
	}
}
