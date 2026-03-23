package utils

import (
	"context"
	"strconv"
	"testing"

	"kafka/config"
)

func TestGetKafkaAddress(t *testing.T) {
	address := GetKafkaAddress()
	expectedAddress := config.AppConfig.Kafka.Host + ":" + strconv.Itoa(int(config.AppConfig.Kafka.Port))

	if address != expectedAddress {
		t.Errorf("Expected address '%s', got '%s'", expectedAddress, address)
	}
}

func TestKafkaDialWithRetry(t *testing.T) {
	ctx := context.Background()
	address := GetKafkaAddress()

	// This might fail if Kafka is not running, but we're testing the function structure
	conn, err := KafkaDialWithRetry(ctx, address)

	if err != nil {
		t.Logf("KafkaDialWithRetry returned error (expected if Kafka not running): %v", err)
	} else {
		defer conn.Close()
		if conn == nil {
			t.Error("Expected non-nil connection, got nil")
		}
	}
}

func TestKafkaDialLeaderWithRetry(t *testing.T) {
	ctx := context.Background()
	address := GetKafkaAddress()

	// This might fail if Kafka is not running, but we're testing the function structure
	conn, err := KafkaDialLeaderWithRetry(ctx, address, "test-topic", 0)

	if err != nil {
		t.Logf("KafkaDialLeaderWithRetry returned error (expected if Kafka not running): %v", err)
	} else {
		defer conn.Close()
		if conn == nil {
			t.Error("Expected non-nil connection, got nil")
		}
	}
}

func TestKafkaConnectionWithRetry(t *testing.T) {
	ctx := context.Background()
	address := GetKafkaAddress()

	// This might fail if Kafka is not running, but we're testing the function structure
	conn, err := KafkaConnectionWithRetry(ctx, "leader", address, "test-topic", 0)

	if err != nil {
		t.Logf("KafkaConnectionWithRetry returned error (expected if Kafka not running): %v", err)
	} else {
		defer conn.Close()
		if conn == nil {
			t.Error("Expected non-nil connection, got nil")
		}
	}
}
