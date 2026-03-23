package circuitbreaker

import (
	"errors"
	"log"
	"sync"
	"time"
)

// CircuitBreakerConfig holds circuit breaker configuration
type CircuitBreakerConfig struct {
	FailureThreshold int           `json:"failure_threshold"`
	ResetTimeout     time.Duration `json:"reset_timeout"`
	HalfOpenLimit    int           `json:"half_open_limit"`
}

// CircuitBreaker implements a simple circuit breaker pattern
type CircuitBreaker struct {
	config       CircuitBreakerConfig
	failureCount int
	lastFailure  time.Time
	state        string // "closed", "open", "half-open"
	mu           sync.RWMutex
}

// NewCircuitBreaker creates a new circuit breaker
func NewCircuitBreaker(config CircuitBreakerConfig) *CircuitBreaker {
	return &CircuitBreaker{
		config: config,
		state:  "closed",
	}
}

// Execute runs the operation with circuit breaker protection
func (cb *CircuitBreaker) Execute(operation func() error) error {
	cb.mu.RLock()
	state := cb.state
	cb.mu.RUnlock()

	switch state {
	case "open":
		if time.Since(cb.lastFailure) > cb.config.ResetTimeout {
			cb.mu.Lock()
			cb.state = "half-open"
			cb.mu.Unlock()
			log.Printf("Circuit breaker transitioning to half-open state")
		} else {
			remainingTime := cb.config.ResetTimeout - time.Since(cb.lastFailure)
			log.Printf("Circuit breaker is open, rejecting operation for %v more", remainingTime)
			return errors.New("circuit breaker is open")
		}
	case "half-open":
		// Allow limited operations in half-open state
		cb.mu.Lock()
		if cb.failureCount >= cb.config.HalfOpenLimit {
			cb.state = "open"
			cb.mu.Unlock()
			log.Printf("Circuit breaker reopened due to failures in half-open state")
			return errors.New("circuit breaker is open")
		}
		cb.mu.Unlock()
	}

	// Execute the operation
	startTime := time.Now()
	err := operation()
	operationDuration := time.Since(startTime)

	cb.mu.Lock()
	defer cb.mu.Unlock()

	if err != nil {
		cb.failureCount++
		cb.lastFailure = time.Now()

		log.Printf("Circuit breaker operation failed in %v: %v (failures: %d/%d)",
			operationDuration, err, cb.failureCount, cb.config.FailureThreshold)

		if cb.failureCount >= cb.config.FailureThreshold {
			cb.state = "open"
			log.Printf("Circuit breaker opened after %d failures. Will reset in %v", cb.failureCount, cb.config.ResetTimeout)
		}
	} else {
		// Success - reset circuit breaker
		previousState := cb.state
		cb.failureCount = 0
		cb.state = "closed"

		if previousState == "half-open" {
			log.Printf("Circuit breaker closed after successful operation in half-open state in %v", operationDuration)
		} else {
			log.Printf("Circuit breaker operation succeeded in %v", operationDuration)
		}
	}

	return err
}

// DefaultCircuitBreakerConfig returns a sensible default circuit breaker configuration
func DefaultCircuitBreakerConfig() CircuitBreakerConfig {
	return CircuitBreakerConfig{
		FailureThreshold: 5,
		ResetTimeout:     30 * time.Second,
		HalfOpenLimit:    2,
	}
}
