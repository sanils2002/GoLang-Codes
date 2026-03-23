package circuitbreaker

import (
	"errors"
	"testing"
	"time"
)

func TestCircuitBreaker(t *testing.T) {
	config := CircuitBreakerConfig{
		FailureThreshold: 2,
		ResetTimeout:     100 * time.Millisecond,
		HalfOpenLimit:    1,
	}

	cb := NewCircuitBreaker(config)

	// Test successful operation
	err := cb.Execute(func() error {
		return nil
	})

	if err != nil {
		t.Errorf("Expected success, got error: %v", err)
	}

	// Test failure threshold
	for i := 0; i < 2; i++ {
		err := cb.Execute(func() error {
			return errors.New("test error")
		})

		if err == nil {
			t.Error("Expected error, got nil")
		}
	}

	// Circuit should be open now
	err = cb.Execute(func() error {
		return nil
	})

	if err == nil {
		t.Error("Expected circuit breaker to be open, got success")
	}

	// Wait for reset timeout
	time.Sleep(150 * time.Millisecond)

	// Should be half-open now
	err = cb.Execute(func() error {
		return nil
	})

	if err != nil {
		t.Errorf("Expected success in half-open state, got error: %v", err)
	}
}

func TestDefaultCircuitBreakerConfig(t *testing.T) {
	config := DefaultCircuitBreakerConfig()

	if config.FailureThreshold != 5 {
		t.Errorf("Expected FailureThreshold to be 5, got %d", config.FailureThreshold)
	}

	if config.ResetTimeout != 30*time.Second {
		t.Errorf("Expected ResetTimeout to be 30s, got %v", config.ResetTimeout)
	}

	if config.HalfOpenLimit != 2 {
		t.Errorf("Expected HalfOpenLimit to be 2, got %d", config.HalfOpenLimit)
	}
}
