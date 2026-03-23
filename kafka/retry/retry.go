package retry

import (
	"context"
	"errors"
	"log"
	"math"
	"math/rand"
	"time"
)

// RetryConfig holds retry configuration parameters
type RetryConfig struct {
	MaxRetries        int           `json:"max_retries"`
	InitialBackoff    time.Duration `json:"initial_backoff"`
	MaxBackoff        time.Duration `json:"max_backoff"`
	BackoffMultiplier float64       `json:"backoff_multiplier"`
	JitterFactor      float64       `json:"jitter_factor"`
	Timeout           time.Duration `json:"timeout"`
}

// RetryableError represents errors that should trigger retries
type RetryableError struct {
	Err error
}

func (e *RetryableError) Error() string {
	return e.Err.Error()
}

func (e *RetryableError) Unwrap() error {
	return e.Err
}

// IsRetryableError checks if an error should trigger a retry
func IsRetryableError(err error) bool {
	if err == nil {
		return false
	}

	// Check for network-related errors
	var netErr interface{ Timeout() bool }
	if errors.As(err, &netErr) {
		log.Printf("Network error detected as retryable: %v", err)
		return true
	}

	// Check for context timeout/cancellation
	if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
		log.Printf("Context error detected as non-retryable: %v", err)
		return false // Don't retry context errors
	}

	// For Kafka errors, we'll retry most errors except for permanent ones
	// This is a conservative approach - retry on most errors
	log.Printf("Error classified as retryable (default behavior): %v", err)
	return true
}

// calculateBackoff calculates the backoff duration with exponential backoff and jitter
func calculateBackoff(attempt int, config RetryConfig) time.Duration {
	// Exponential backoff: base * multiplier^attempt
	backoff := float64(config.InitialBackoff) * math.Pow(config.BackoffMultiplier, float64(attempt-1))

	// Cap at max backoff
	if backoff > float64(config.MaxBackoff) {
		backoff = float64(config.MaxBackoff)
	}

	// Add jitter to prevent thundering herd
	jitter := backoff * config.JitterFactor * (rand.Float64() - 0.5)
	backoff += jitter

	// Ensure minimum backoff
	if backoff < float64(config.InitialBackoff) {
		backoff = float64(config.InitialBackoff)
	}

	return time.Duration(backoff)
}

// RetryWithBackoff executes a function with retry logic and exponential backoff
func RetryWithBackoff(ctx context.Context, operation func() error, config RetryConfig) error {
	// Protect against nil context
	if ctx == nil {
		ctx = context.Background()
		log.Printf("Warning: Nil context provided, using background context")
	}

	var lastErr error
	startTime := time.Now()

	for attempt := 1; attempt <= config.MaxRetries+1; attempt++ {
		// Check context cancellation
		select {
		case <-ctx.Done():
			log.Printf("Context cancelled during retry: %v", ctx.Err())
			return ctx.Err()
		default:
		}

		// Execute the operation
		operationStart := time.Now()
		err := operation()
		operationDuration := time.Since(operationStart)

		if err == nil {
			totalDuration := time.Since(startTime)
			log.Printf("Operation succeeded after %d attempts in %v", attempt, totalDuration)
			return nil // Success
		}

		lastErr = err

		// Check if error is retryable
		if !IsRetryableError(err) {
			log.Printf("Non-retryable error encountered: %v (attempt %d)", err, attempt)
			return err
		}

		// If this was the last attempt, return the error
		if attempt > config.MaxRetries {
			totalDuration := time.Since(startTime)
			log.Printf("Max retries (%d) exceeded after %v. Last error: %v", config.MaxRetries, totalDuration, err)
			return err
		}

		// Calculate backoff duration
		backoff := calculateBackoff(attempt, config)

		log.Printf("Retrying operation - attempt %d/%d with backoff %v after error: %v (operation took %v)", 
			attempt, config.MaxRetries+1, backoff, err, operationDuration)

		// Wait before retry
		select {
		case <-time.After(backoff):
		case <-ctx.Done():
			log.Printf("Context cancelled during backoff: %v", ctx.Err())
			return ctx.Err()
		}
	}

	return lastErr
}

// DefaultRetryConfig returns a sensible default retry configuration
func DefaultRetryConfig() RetryConfig {
	return RetryConfig{
		MaxRetries:        3,
		InitialBackoff:    100 * time.Millisecond,
		MaxBackoff:        5 * time.Second,
		BackoffMultiplier: 2.0,
		JitterFactor:      0.1,
		Timeout:           30 * time.Second,
	}
}
