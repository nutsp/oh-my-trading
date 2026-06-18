package config

import (
	"os"
	"time"
)

type Config struct {
	HTTPAddr        string
	Environment     string
	DatabaseURL     string
	ShutdownTimeout time.Duration
	APIMockMode     bool
}

func Load() Config {
	return Config{
		HTTPAddr:        envString("OMT_HTTP_ADDR", ":8080"),
		Environment:     envString("OMT_ENV", "development"),
		DatabaseURL:     envString("OMT_DATABASE_URL", "postgres://omt:omt_local_password@localhost:15432/oh_my_trading?sslmode=disable"),
		ShutdownTimeout: envDuration("OMT_SHUTDOWN_TIMEOUT", 10*time.Second),
		APIMockMode:     envBool("OMT_API_MOCK_MODE", false),
	}
}

func envString(key string, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func envDuration(key string, fallback time.Duration) time.Duration {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	duration, err := time.ParseDuration(value)
	if err != nil {
		return fallback
	}
	return duration
}

func envBool(key string, fallback bool) bool {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	switch value {
	case "1", "true", "TRUE", "True", "yes", "YES", "on", "ON":
		return true
	case "0", "false", "FALSE", "False", "no", "NO", "off", "OFF":
		return false
	default:
		return fallback
	}
}
