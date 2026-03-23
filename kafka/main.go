package kafka

import (
	"kafka/examples"
	"log"
)

// RunAllExamples demonstrates how to use the Kafka retry functionality
func RunAllExamples() {
	log.Printf("Starting Kafka Examples...")

	// Example 1: Basic retry with default configuration
	examples.BasicRetryExample()

	// Example 2: Custom retry configuration for high-priority messages
	examples.HighPriorityRetryExample()

	// Example 3: Circuit breaker with retry
	examples.CircuitBreakerExample()

	// Example 4: Different retry configurations for different use cases
	examples.DifferentUseCasesExample()

	// Example 5: Retry-enabled connection functions
	examples.ConnectionRetryExample()

	// Example 6: Producer and Consumer examples
	examples.ProducerExample()
	examples.ConsumerExample()
}
