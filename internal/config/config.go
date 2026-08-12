package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"gopkg.in/yaml.v3"
)

const (
	defaultPort            = 8080
	defaultContentURL      = "https://raw.githubusercontent.com/karanpratapsingh/system-design/main/README.md"
	defaultContentDir      = "./content"
	defaultFetchOnStart    = true
	defaultLogLevel        = "info"
	defaultLogFormat       = "text"
	defaultShutdownTimeout = 10 * time.Second
	defaultConfigFile      = "blogo.yaml"
)

type RepoConfig struct {
	Name   string `yaml:"name"`
	URL    string `yaml:"url"`
	Type   string `yaml:"type"`
	Branch string `yaml:"branch"`
	Author string `yaml:"author"`
}

type Config struct {
	Port            int           `yaml:"port"`
	ContentDir      string        `yaml:"content_dir"`
	FetchOnStart    bool          `yaml:"fetch_on_start"`
	LogLevel        string        `yaml:"log_level"`
	LogFormat       string        `yaml:"log_format"`
	ShutdownTimeout time.Duration `yaml:"-"`
	Repos           []RepoConfig  `yaml:"repos"`

	// Legacy single-repo fields (env-var fallback)
	ContentURL string `yaml:"-"`
}

func Load() (*Config, error) {
	configPath := envString("BLOGO_CONFIG", defaultConfigFile)

	if data, err := os.ReadFile(configPath); err == nil {
		return loadFromFile(data)
	}

	return loadFromEnv()
}

func loadFromFile(data []byte) (*Config, error) {
	cfg := &Config{
		Port:         defaultPort,
		ContentDir:   defaultContentDir,
		FetchOnStart: defaultFetchOnStart,
		LogLevel:     defaultLogLevel,
		LogFormat:    defaultLogFormat,
	}

	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("parsing config file: %w", err)
	}

	if cfg.Port < 1 || cfg.Port > 65535 {
		return nil, fmt.Errorf("port must be between 1 and 65535, got %d", cfg.Port)
	}

	for i := range cfg.Repos {
		if cfg.Repos[i].Branch == "" {
			cfg.Repos[i].Branch = "main"
		}
		if cfg.Repos[i].Type == "" {
			cfg.Repos[i].Type = "single-md"
		}
	}

	if len(cfg.Repos) == 0 {
		return nil, fmt.Errorf("config file must define at least one repo")
	}

	// Allow env overrides for server-level settings
	if v := os.Getenv("BLOGO_PORT"); v != "" {
		port, err := strconv.Atoi(v)
		if err != nil {
			return nil, fmt.Errorf("invalid BLOGO_PORT: %w", err)
		}
		cfg.Port = port
	}
	if v := os.Getenv("BLOGO_LOG_LEVEL"); v != "" {
		cfg.LogLevel = v
	}
	if v := os.Getenv("BLOGO_LOG_FORMAT"); v != "" {
		cfg.LogFormat = v
	}

	shutdownTimeout, err := envDuration("BLOGO_SHUTDOWN_TIMEOUT", defaultShutdownTimeout)
	if err != nil {
		return nil, fmt.Errorf("invalid BLOGO_SHUTDOWN_TIMEOUT: %w", err)
	}
	cfg.ShutdownTimeout = shutdownTimeout

	return cfg, nil
}

func loadFromEnv() (*Config, error) {
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

	contentURL := envString("BLOGO_CONTENT_URL", defaultContentURL)

	return &Config{
		Port:            port,
		ContentURL:      contentURL,
		ContentDir:      envString("BLOGO_CONTENT_DIR", defaultContentDir),
		FetchOnStart:    envBool("BLOGO_FETCH_ON_START", defaultFetchOnStart),
		LogLevel:        envString("BLOGO_LOG_LEVEL", defaultLogLevel),
		LogFormat:       envString("BLOGO_LOG_FORMAT", defaultLogFormat),
		ShutdownTimeout: shutdownTimeout,
		Repos: []RepoConfig{
			{
				Name:   "System Design",
				URL:    contentURL,
				Type:   "single-md",
				Branch: "main",
				Author: "Karan Pratap Singh",
			},
		},
	}, nil
}

func (c *Config) IsMultiRepo() bool {
	return len(c.Repos) > 1
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
