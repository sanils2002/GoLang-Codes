package examples

import (
	"context"
	"log"
	"time"

	"kafka/circuitbreaker"
	"kafka/retry"
)

// CircuitBreakerExample demonstrates circuit breaker functionality
func CircuitBreakerExample() {
	// Create circuit breaker
	circuitBreakerConfig := circuitbreaker.CircuitBreakerConfig{
		FailureThreshold: 3,
		ResetTimeout:     30 * time.Second,
		HalfOpenLimit:    1,
	}
	circuitBreaker := circuitbreaker.NewCircuitBreaker(circuitBreakerConfig)

	// Operation that might fail
	operation := func() error {
		// Simulate Kafka operation
		return nil
	}

	// Execute with circuit breaker
	err := circuitBreaker.Execute(operation)
	if err != nil {
		log.Printf("Circuit breaker operation failed: %v", err)
		return
	}

	log.Printf("Circuit breaker example completed successfully")
}

// CircuitBreakerWithRetryExample demonstrates combining circuit breaker with retry
func CircuitBreakerWithRetryExample() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Create circuit breaker
	circuitBreakerConfig := circuitbreaker.CircuitBreakerConfig{
		FailureThreshold: 3,
		ResetTimeout:     30 * time.Second,
		HalfOpenLimit:    1,
	}
	circuitBreaker := circuitbreaker.NewCircuitBreaker(circuitBreakerConfig)

	// Retry configuration
	retryConfig := retry.RetryConfig{
		MaxRetries:        3,
		InitialBackoff:    100 * time.Millisecond,
		MaxBackoff:        5 * time.Second,
		BackoffMultiplier: 2.0,
		JitterFactor:      0.1,
		Timeout:           30 * time.Second,
	}

	// Operation that might fail
	operation := func() error {
		// Simulate Kafka operation
		return nil
	}

	// Execute with retry first, then circuit breaker
	err := retry.RetryWithBackoff(ctx, func() error {
		return circuitBreaker.Execute(operation)
	}, retryConfig)

	if err != nil {
		log.Printf("Circuit breaker with retry failed: %v", err)
		return
	}

	log.Printf("Circuit breaker with retry example completed successfully")
}
