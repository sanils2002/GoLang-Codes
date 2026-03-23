package consumer

import (
	"testing"
)

func TestKafkaRead(t *testing.T) {
	reader := KafkaRead("test-topic", "test-group", 10e3, 10e6)
	defer reader.Close()
	
	if reader == nil {
		t.Error("Expected non-nil reader, got nil")
	}
}

func TestKafkaReadWithRetry(t *testing.T) {
	reader := KafkaReadWithRetry("test-topic", "test-group", 10e3, 10e6)
	defer reader.Close()
	
	if reader == nil {
		t.Error("Expected non-nil reader, got nil")
	}
}
