package examples

import (
	"testing"
	"time"
)

func TestProducerExample(t *testing.T) {
	// Test that the function runs without panicking
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("ProducerExample panicked: %v", r)
		}
	}()

	// Run the example
	ProducerExample()
}

func TestConsumerExample(t *testing.T) {
	// Test that the function runs without panicking
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("ConsumerExample panicked: %v", r)
		}
	}()

	// Run the example
	ConsumerExample()
}

func TestExampleMessage(t *testing.T) {
	// Test ExampleMessage struct creation
	message := ExampleMessage{
		ID:      "test-123",
		Type:    "test",
		Data:    map[string]interface{}{"key": "value"},
		Retry:   0,
		Created: time.Now(),
	}

	if message.ID != "test-123" {
		t.Errorf("Expected ID 'test-123', got '%s'", message.ID)
	}

	if message.Type != "test" {
		t.Errorf("Expected Type 'test', got '%s'", message.Type)
	}

	if message.Retry != 0 {
		t.Errorf("Expected Retry 0, got %d", message.Retry)
	}
}
