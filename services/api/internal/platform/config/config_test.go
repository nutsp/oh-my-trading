package config

import (
	"testing"
	"time"
)

func TestLoadUsesDefaults(t *testing.T) {
	t.Setenv("OMT_HTTP_ADDR", "")
	t.Setenv("OMT_ENV", "")
	t.Setenv("OMT_DATABASE_URL", "")
	t.Setenv("OMT_SHUTDOWN_TIMEOUT", "")
	t.Setenv("OMT_API_MOCK_MODE", "")

	cfg := Load()

	if cfg.HTTPAddr != ":8080" {
		t.Fatalf("HTTPAddr = %q, want %q", cfg.HTTPAddr, ":8080")
	}
	if cfg.Environment != "development" {
		t.Fatalf("Environment = %q, want %q", cfg.Environment, "development")
	}
	if cfg.DatabaseURL != "postgres://omt:omt_local_password@localhost:15432/oh_my_trading?sslmode=disable" {
		t.Fatalf("DatabaseURL = %q", cfg.DatabaseURL)
	}
	if cfg.ShutdownTimeout != 10*time.Second {
		t.Fatalf("ShutdownTimeout = %s, want %s", cfg.ShutdownTimeout, 10*time.Second)
	}
	if cfg.APIMockMode {
		t.Fatalf("APIMockMode = %t, want false", cfg.APIMockMode)
	}
}

func TestLoadUsesEnvironmentOverrides(t *testing.T) {
	t.Setenv("OMT_HTTP_ADDR", ":9090")
	t.Setenv("OMT_ENV", "test")
	t.Setenv("OMT_DATABASE_URL", "postgres://example")
	t.Setenv("OMT_SHUTDOWN_TIMEOUT", "3s")
	t.Setenv("OMT_API_MOCK_MODE", "true")

	cfg := Load()

	if cfg.HTTPAddr != ":9090" {
		t.Fatalf("HTTPAddr = %q, want %q", cfg.HTTPAddr, ":9090")
	}
	if cfg.Environment != "test" {
		t.Fatalf("Environment = %q, want %q", cfg.Environment, "test")
	}
	if cfg.DatabaseURL != "postgres://example" {
		t.Fatalf("DatabaseURL = %q, want postgres://example", cfg.DatabaseURL)
	}
	if cfg.ShutdownTimeout != 3*time.Second {
		t.Fatalf("ShutdownTimeout = %s, want %s", cfg.ShutdownTimeout, 3*time.Second)
	}
	if !cfg.APIMockMode {
		t.Fatalf("APIMockMode = %t, want true", cfg.APIMockMode)
	}
}

func TestLoadUsesBoolFallbackOnInvalidValue(t *testing.T) {
	t.Setenv("OMT_API_MOCK_MODE", "not-a-bool")

	cfg := Load()

	if cfg.APIMockMode {
		t.Fatalf("APIMockMode = %t, want false on invalid value", cfg.APIMockMode)
	}
}
