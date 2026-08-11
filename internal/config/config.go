package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

const (
	defaultPort            = 8080
	defaultContentURL      = "https://raw.githubusercontent.com/karanpratapsingh/system-design/main/README.md"
	defaultContentDir      = "./content"
	defaultFetchOnStart    = true
	defaultLogLevel        = "info"
	defaultLogFormat       = "text"
	defaultShutdownTimeout = 10 * time.Second
)

type Config struct {
	Port            int
	ContentURL      string
	ContentDir      string
	FetchOnStart    bool
	LogLevel        string
	LogFormat       string
	ShutdownTimeout time.Duration
}

func Load() (*Config, error) {
	port, err := envInt("BLOGO_PORT", defaultPort)
	if err != nil {
		return nil, fmt.Errorf("invalid BLOGO_PORT: %w", err)
	}
	if port < 1 || port > 65535 {
		return nil, fmt.Errorf("BLOGO_PORT must be between 1 and 65535, got %d", port)
	}

	shutdownTimeout, err := envDuration("BLOGO_SHUTDOWN_TIMEOUT", defaultShutdownTimeout)
	if err != nil {
		return nil, fmt.Errorf("invalid BLOGO_SHUTDOWN_TIMEOUT: %w", err)
	}

	return &Config{
		Port:            port,
		ContentURL:      envString("BLOGO_CONTENT_URL", defaultContentURL),
		ContentDir:      envString("BLOGO_CONTENT_DIR", defaultContentDir),
		FetchOnStart:    envBool("BLOGO_FETCH_ON_START", defaultFetchOnStart),
		LogLevel:        envString("BLOGO_LOG_LEVEL", defaultLogLevel),
		LogFormat:       envString("BLOGO_LOG_FORMAT", defaultLogFormat),
		ShutdownTimeout: shutdownTimeout,
	}, nil
}

func envString(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func envInt(key string, fallback int) (int, error) {
	v := os.Getenv(key)
	if v == "" {
		return fallback, nil
	}
	return strconv.Atoi(v)
}

func envBool(key string, fallback bool) bool {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	b, err := strconv.ParseBool(v)
	if err != nil {
		return fallback
	}
	return b
}

func envDuration(key string, fallback time.Duration) (time.Duration, error) {
	v := os.Getenv(key)
	if v == "" {
		return fallback, nil
	}
	return time.ParseDuration(v)
}
