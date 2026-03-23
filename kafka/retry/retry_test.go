package retry

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestRetryWithBackoff(t *testing.T) {
	ctx := context.Background()
	retryConfig := RetryConfig{
		MaxRetries:        2,
		InitialBackoff:    10 * time.Millisecond,
		MaxBackoff:        100 * time.Millisecond,
		BackoffMultiplier: 2.0,
		JitterFactor:      0.1,
		Timeout:           1 * time.Second,
	}

	// Test successful operation
	attempts := 0
	err := RetryWithBackoff(ctx, func() error {
		attempts++
		if attempts == 1 {
			return errors.New("temporary error")
		}
		return nil
	}, retryConfig)

	if err != nil {
		t.Errorf("Expected success, got error: %v", err)
	}

	if attempts != 2 {
		t.Errorf("Expected 2 attempts, got %d", attempts)
	}
}

func TestRetryWithBackoffMaxRetries(t *testing.T) {
	ctx := context.Background()
	retryConfig := RetryConfig{
		MaxRetries:        2,
		InitialBackoff:    10 * time.Millisecond,
		MaxBackoff:        100 * time.Millisecond,
		BackoffMultiplier: 2.0,
		JitterFactor:      0.1,
		Timeout:           1 * time.Second,
	}

	// Test max retries exceeded
	attempts := 0
	err := RetryWithBackoff(ctx, func() error {
		attempts++
		return errors.New("persistent error")
	}, retryConfig)

	if err == nil {
		t.Error("Expected error, got nil")
	}

	if attempts != 3 { // MaxRetries + 1
		t.Errorf("Expected 3 attempts, got %d", attempts)
	}
}

func TestRetryWithBackoffContextCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	retryConfig := RetryConfig{
		MaxRetries:        5,
		InitialBackoff:    10 * time.Millisecond,
		MaxBackoff:        100 * time.Millisecond,
		BackoffMultiplier: 2.0,
		JitterFactor:      0.1,
		Timeout:           1 * time.Second,
	}

	// Cancel context after first attempt
	attempts := 0
	err := RetryWithBackoff(ctx, func() error {
		attempts++
		if attempts == 1 {
			cancel()
		}
		return errors.New("temporary error")
	}, retryConfig)

	if err == nil {
		t.Error("Expected context cancellation error, got nil")
	}

	if attempts != 1 {
		t.Errorf("Expected 1 attempt, got %d", attempts)
	}
}

func TestIsRetryableError(t *testing.T) {
	// Test retryable error
	retryableErr := errors.New("network timeout")
	if !IsRetryableError(retryableErr) {
		t.Error("Expected network timeout to be retryable")
	}

	// Test non-retryable error (currently all errors are treated as retryable by default)
	nonRetryableErr := errors.New("invalid input")
	if !IsRetryableError(nonRetryableErr) {
		t.Error("Expected invalid input to be retryable (current default behavior)")
	}
}

func TestDefaultRetryConfig(t *testing.T) {
	config := DefaultRetryConfig()

	if config.MaxRetries != 3 {
		t.Errorf("Expected MaxRetries to be 3, got %d", config.MaxRetries)
	}

	if config.InitialBackoff != 100*time.Millisecond {
		t.Errorf("Expected InitialBackoff to be 100ms, got %v", config.InitialBackoff)
	}

	if config.MaxBackoff != 5*time.Second {
		t.Errorf("Expected MaxBackoff to be 5s, got %v", config.MaxBackoff)
	}
}

func TestRetryWithBackoffNilContext(t *testing.T) {
	retryConfig := RetryConfig{
		MaxRetries:        2,
		InitialBackoff:    10 * time.Millisecond,
		MaxBackoff:        100 * time.Millisecond,
		BackoffMultiplier: 2.0,
		JitterFactor:      0.1,
		Timeout:           1 * time.Second,
	}

	// Test with nil context
	err := RetryWithBackoff(nil, func() error {
		return errors.New("test error")
	}, retryConfig)

	if err == nil {
		t.Error("Expected error with nil context, got nil")
	}
}
