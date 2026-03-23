# Kafka Package - Organized Structure

This package provides a clean, organized Kafka client implementation with retry logic, circuit breakers, simplified logging, and comprehensive examples.

## 🔄 Why Retry Mechanisms Are Essential

### **Network Reliability Challenges**
In distributed systems, network operations are inherently unreliable. Kafka clients frequently encounter temporary failures due to:
- Network timeouts and connection drops
- Temporary service unavailability
- Load balancer issues
- DNS resolution problems
- Infrastructure maintenance windows

### **Kafka-Specific Failure Scenarios**
Kafka operations are particularly susceptible to transient failures:
- Broker leader elections causing temporary unavailability
- Partition reassignments during cluster scaling
- Consumer group rebalancing events
- High load causing broker overload
- Disk I/O bottlenecks on brokers

### **Business Impact of Failures**
Without proper retry mechanisms, temporary failures can cause:
- **Data Loss**: Failed message production without retry
- **Service Disruption**: Consumer applications stopping on temporary errors
- **Poor User Experience**: Application timeouts and errors
- **Operational Overhead**: Manual intervention for transient issues
- **Inconsistent State**: Partial data processing due to failures

### **Retry Strategy Benefits**
Implementing intelligent retry mechanisms provides:
- **Improved Reliability**: Automatic recovery from transient failures
- **Better User Experience**: Seamless operation despite temporary issues
- **Reduced Operational Overhead**: Self-healing systems
- **Data Consistency**: Ensured message delivery and processing
- **System Resilience**: Graceful handling of infrastructure changes

## 📁 Directory Structure

```
kafka/
├── config/
│   └── config.go                    # Configuration structures and settings
├── retry/
│   ├── retry.go                     # Retry logic with exponential backoff
│   └── retry_test.go                # Retry package tests
├── circuitbreaker/
│   ├── circuitbreaker.go            # Circuit breaker pattern implementation
│   └── circuitbreaker_test.go       # Circuit breaker package tests
├── producer/
│   ├── producer.go                  # Kafka producer functionality
│   └── producer_test.go             # Producer package tests
├── consumer/
│   ├── consumer.go                  # Kafka consumer functionality
│   └── consumer_test.go             # Consumer package tests
├── utils/
│   ├── connection.go                # Connection utilities
│   └── connection_test.go           # Utils package tests
├── examples/
│   ├── basic_examples.go            # Basic retry examples
│   ├── basic_examples_test.go       # Basic examples tests
│   ├── circuit_breaker_examples.go  # Circuit breaker examples
│   ├── circuit_breaker_examples_test.go # Circuit breaker examples tests
│   ├── producer_consumer_examples.go # Producer and consumer examples
│   ├── producer_consumer_examples_test.go # Producer/consumer examples tests
│   ├── connection_examples.go       # Connection examples
│   └── connection_examples_test.go  # Connection examples tests
├── main.go                          # Main orchestrator for examples
├── go.mod                           # Go module definition
└── README.md                        # This file
```

## 🏗️ Architecture

### **Config Package** (`config/`)
- **Config struct**: Main configuration structure
- **KafkaConfig struct**: Kafka-specific settings
- **ContextStruct**: Application context structure
- **AppConfig**: Global configuration instance

### **Retry Package** (`retry/`)
- **RetryConfig**: Retry configuration parameters
- **RetryWithBackoff**: Exponential backoff with jitter
- **IsRetryableError**: Error classification logic
- **DefaultRetryConfig**: Sensible defaults

### **Circuit Breaker Package** (`circuitbreaker/`)
- **CircuitBreaker**: Circuit breaker implementation
- **CircuitBreakerConfig**: Circuit breaker settings
- **Execute**: Protected operation execution
- **DefaultCircuitBreakerConfig**: Default settings

### **Producer Package** (`producer/`)
- **KafkaProducer**: Basic message producer
- **KafkaProducerWithRetry**: Producer with retry logic
- **KafkaWriter**: Basic Kafka writer
- **KafkaWriterWithRetry**: Writer with retry configuration

### **Consumer Package** (`consumer/`)
- **KafkaRead**: Basic Kafka reader
- **KafkaReadWithRetry**: Reader with retry logic

### **Utils Package** (`utils/`)
- **KafkaConnectionWithRetry**: Connection with retry
- **KafkaDialLeaderWithRetry**: Leader connection
- **KafkaDialWithRetry**: Basic connection
- **GetKafkaAddress**: Address formatting

### **Examples Package** (`examples/`)
- **Basic Examples**: Basic retry functionality demonstrations
- **Circuit Breaker Examples**: Circuit breaker pattern usage
- **Producer/Consumer Examples**: Message production and consumption
- **Connection Examples**: Connection management and utilities
- **Comprehensive Tests**: Full test coverage for all examples

### **Main Orchestrator** (`main.go`)
- **RunAllExamples**: Orchestrates all example demonstrations
- **Example Integration**: Shows how to use all packages together

## 🚀 Usage

### Running Examples
The package includes comprehensive examples demonstrating all functionality:

- **Basic Retry Examples**: Simple retry operations with default configuration
- **High Priority Retry Examples**: Custom retry configurations for critical operations
- **Circuit Breaker Examples**: Fault tolerance pattern implementation
- **Producer/Consumer Examples**: Message production and consumption with retry logic
- **Connection Examples**: Robust connection management with retry capabilities

### Package Usage
Each package can be used independently or together:

- **Config Package**: Centralized configuration management
- **Retry Package**: Exponential backoff with jitter for resilient operations
- **Circuit Breaker Package**: Fault tolerance pattern for system protection
- **Producer Package**: Kafka message production with retry capabilities
- **Consumer Package**: Kafka message consumption with retry logic
- **Utils Package**: Connection utilities with built-in retry mechanisms

### Testing
All packages include comprehensive test coverage:

- **Unit Tests**: Individual package functionality testing
- **Integration Tests**: Cross-package integration testing
- **Example Tests**: Real-world usage scenario testing
- **Error Handling Tests**: Robust error handling validation

## ✨ Features

- ✅ **Modular Design**: Clean separation of concerns with organized package structure
- ✅ **Retry Logic**: Exponential backoff with jitter for resilient operations
- ✅ **Circuit Breaker**: Fault tolerance pattern for system protection
- ✅ **Simple Logging**: Standard Go logging without complex dependencies
- ✅ **Configuration**: Centralized configuration management
- ✅ **Connection Management**: Robust connection handling with retry mechanisms
- ✅ **Comprehensive Examples**: Real-world usage demonstrations
- ✅ **Full Test Coverage**: Unit tests, integration tests, and example tests
- ✅ **Error Handling**: Robust error classification and handling
- ✅ **Context Support**: Proper context management and timeout handling

## 🔧 Dependencies

```go
require (
    github.com/segmentio/kafka-go v0.4.47
)
```

## 📝 Notes

- All complex enterprise features (tracing, custom logging) have been removed for simplicity
- Uses standard Go `log` package for consistent logging
- Configuration is centralized and easily modifiable
- Retry logic is configurable and robust with exponential backoff
- Circuit breaker provides fault tolerance and system protection
- Clean, maintainable code structure with comprehensive documentation
- All examples include proper error handling and logging
- Test coverage ensures reliability and correctness
- Modular design allows for easy extension and customization
