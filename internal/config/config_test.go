package config

import (
	"testing"
	"time"
)

func TestLoadDefaults(t *testing.T) {
	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() returned error: %v", err)
	}

	if cfg.Port != 8080 {
		t.Errorf("Port = %d, want 8080", cfg.Port)
	}
	if cfg.ContentURL != defaultContentURL {
		t.Errorf("ContentURL = %q, want %q", cfg.ContentURL, defaultContentURL)
	}
	if cfg.ContentDir != "./content" {
		t.Errorf("ContentDir = %q, want ./content", cfg.ContentDir)
	}
	if cfg.FetchOnStart != true {
		t.Errorf("FetchOnStart = %v, want true", cfg.FetchOnStart)
	}
	if cfg.LogLevel != "info" {
		t.Errorf("LogLevel = %q, want info", cfg.LogLevel)
	}
	if cfg.LogFormat != "text" {
		t.Errorf("LogFormat = %q, want text", cfg.LogFormat)
	}
	if cfg.ShutdownTimeout != 10*time.Second {
		t.Errorf("ShutdownTimeout = %v, want 10s", cfg.ShutdownTimeout)
	}
}

func TestLoadEnvOverrides(t *testing.T) {
	t.Setenv("BLOGO_PORT", "9090")
	t.Setenv("BLOGO_CONTENT_URL", "https://example.com/test.md")
	t.Setenv("BLOGO_CONTENT_DIR", "/tmp/content")
	t.Setenv("BLOGO_FETCH_ON_START", "false")
	t.Setenv("BLOGO_LOG_LEVEL", "debug")
	t.Setenv("BLOGO_LOG_FORMAT", "json")
	t.Setenv("BLOGO_SHUTDOWN_TIMEOUT", "30s")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() returned error: %v", err)
	}

	if cfg.Port != 9090 {
		t.Errorf("Port = %d, want 9090", cfg.Port)
	}
	if cfg.ContentURL != "https://example.com/test.md" {
		t.Errorf("ContentURL = %q, want https://example.com/test.md", cfg.ContentURL)
	}
	if cfg.ContentDir != "/tmp/content" {
		t.Errorf("ContentDir = %q, want /tmp/content", cfg.ContentDir)
	}
	if cfg.FetchOnStart != false {
		t.Errorf("FetchOnStart = %v, want false", cfg.FetchOnStart)
	}
	if cfg.LogLevel != "debug" {
		t.Errorf("LogLevel = %q, want debug", cfg.LogLevel)
	}
	if cfg.LogFormat != "json" {
		t.Errorf("LogFormat = %q, want json", cfg.LogFormat)
	}
	if cfg.ShutdownTimeout != 30*time.Second {
		t.Errorf("ShutdownTimeout = %v, want 30s", cfg.ShutdownTimeout)
	}
}

func TestLoadInvalidPort(t *testing.T) {
	t.Setenv("BLOGO_PORT", "not-a-number")
	_, err := Load()
	if err == nil {
		t.Fatal("Load() should return error for invalid port")
	}
}

func TestLoadPortOutOfRange(t *testing.T) {
	t.Setenv("BLOGO_PORT", "99999")
	_, err := Load()
	if err == nil {
		t.Fatal("Load() should return error for port out of range")
	}
}

func TestLoadInvalidShutdownTimeout(t *testing.T) {
	t.Setenv("BLOGO_SHUTDOWN_TIMEOUT", "invalid")
	_, err := Load()
	if err == nil {
		t.Fatal("Load() should return error for invalid shutdown timeout")
	}
}
