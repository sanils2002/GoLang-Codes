package examples

import (
	"testing"
)

func TestCircuitBreakerExample(t *testing.T) {
	// Test that the function runs without panicking
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("CircuitBreakerExample panicked: %v", r)
		}
	}()

	// Run the example
	CircuitBreakerExample()
}

func TestCircuitBreakerWithRetryExample(t *testing.T) {
	// Test that the function runs without panicking
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("CircuitBreakerWithRetryExample panicked: %v", r)
		}
	}()

	// Run the example
	CircuitBreakerWithRetryExample()
}
