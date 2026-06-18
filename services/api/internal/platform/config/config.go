package config

import (
	"os"
	"time"
)

type Config struct {
	HTTPAddr        string
	Environment     string
	ShutdownTimeout time.Duration
}

func Load() Config {
	return Config{
		HTTPAddr:        envString("OMT_HTTP_ADDR", ":8080"),
		Environment:     envString("OMT_ENV", "development"),
		ShutdownTimeout: envDuration("OMT_SHUTDOWN_TIMEOUT", 10*time.Second),
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
