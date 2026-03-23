package config

import "time"

func Load() (Config, error) {
	if err := loadOptionalEnvFile(); err != nil {
		return Config{}, err
	}

	cfg := Config{
		HTTP: HTTPConfig{
			Host:            env("HTTP_HOST", "0.0.0.0"),
			Port:            envInt("HTTP_PORT", 8080),
			ReadTimeout:     envDuration("HTTP_READ_TIMEOUT", 5*time.Second),
			WriteTimeout:    envDuration("HTTP_WRITE_TIMEOUT", 10*time.Second),
			ShutdownTimeout: envDuration("HTTP_SHUTDOWN_TIMEOUT", 30*time.Second),
		},
		Database: DatabaseConfig{
			DSN:             envRequired("DATABASE_DSN"),
			MaxOpenConns:    envInt("DATABASE_MAX_OPEN_CONNS", 25),
			MaxIdleConns:    envInt("DATABASE_MAX_IDLE_CONNS", 5),
			ConnMaxLifetime: envDuration("DATABASE_CONN_MAX_LIFETIME", 5*time.Minute),
		},
		Auth: AuthConfig{
			JWTSecret:       envRequired("AUTH_JWT_SECRET"),
			TokenTTL:        envDuration("AUTH_TOKEN_TTL", 15*time.Minute),
			RefreshTokenTTL: envDuration("AUTH_REFRESH_TOKEN_TTL", 7*24*time.Hour),
		},
		Observability: ObservabilityConfig{
			LogLevel:     env("LOG_LEVEL", "info"),
			OTLPEndpoint: env("OTEL_EXPORTER_OTLP_ENDPOINT", ""),
			ServiceName:  env("OTEL_SERVICE_NAME", "api"),
		},
	}

	return cfg, cfg.Validate()
}
