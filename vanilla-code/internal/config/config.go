package config

import (
	"fmt"
	"strings"
	"time"
)

type Config struct {
	HTTP            HTTPConfig
	Database        DatabaseConfig
	Auth            AuthConfig
	Observability   ObservabilityConfig
}

type HTTPConfig struct {
	Host            string
	Port            int
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	ShutdownTimeout time.Duration
}

type DatabaseConfig struct {
	DSN             string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

type AuthConfig struct {
	JWTSecret       string
	TokenTTL        time.Duration
	RefreshTokenTTL time.Duration
}

type ObservabilityConfig struct {
	LogLevel     string
	OTLPEndpoint string
	ServiceName  string
}

func (c Config) Validate() error {
	var errs []string

	if c.Database.DSN == "" {
		errs = append(errs, "DATABASE_DSN is required")
	}
	if c.Auth.JWTSecret == "" {
		errs = append(errs, "AUTH_JWT_SECRET is required")
	}
	if len(c.Auth.JWTSecret) < 32 {
		errs = append(errs, "AUTH_JWT_SECRET must be at least 32 characters")
	}
	if c.HTTP.Port < 1 || c.HTTP.Port > 65535 {
		errs = append(errs, fmt.Sprintf("HTTP_PORT %d is out of range", c.HTTP.Port))
	}
	if c.Database.MaxOpenConns < c.Database.MaxIdleConns {
		errs = append(errs, "DATABASE_MAX_OPEN_CONNS must be >= DATABASE_MAX_IDLE_CONNS")
	}

	if len(errs) > 0 {
		return fmt.Errorf("configuration errors:\n  - %s", strings.Join(errs, "\n  - "))
	}
	return nil
}
