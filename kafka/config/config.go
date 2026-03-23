package config

import (
	"context"
)

// Config represents the main configuration structure
type Config struct {
	Kafka KafkaConfig
}

// KafkaConfig represents Kafka-specific configuration
type KafkaConfig struct {
	Host     string
	Port     int
	Username string
	Password string
}

// ContextStruct represents the context structure used throughout the application
type ContextStruct struct {
	EchoContext      context.Context
	SpanLabelPrefix  string
	SpanLabelPostfix string
	DDCurrentSpan    interface{} // Placeholder for DataDog span
}

// Global configuration instance
var AppConfig = &Config{
	Kafka: KafkaConfig{
		Host:     "localhost",
		Port:     9092,
		Username: "user",
		Password: "password",
	},
}
