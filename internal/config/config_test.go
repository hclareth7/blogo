package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestLoadDefaults(t *testing.T) {
	t.Setenv("BLOGO_CONFIG", "/nonexistent/blogo.yaml")

	origDir, _ := os.Getwd()
	tmpDir := t.TempDir()
	os.Chdir(tmpDir)
	defer os.Chdir(origDir)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() returned error: %v", err)
	}

	if cfg.Port != 8080 {
		t.Errorf("Port = %d, want 8080", cfg.Port)
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
	if len(cfg.Repos) != 1 {
		t.Fatalf("Repos count = %d, want 1", len(cfg.Repos))
	}
	if cfg.Repos[0].Type != "single-md" {
		t.Errorf("Repos[0].Type = %q, want single-md", cfg.Repos[0].Type)
	}
}

func TestLoadFromYAML(t *testing.T) {
	yaml := `
port: 9090
log_level: debug
fetch_on_start: false
repos:
  - name: "Test Repo"
    url: "https://github.com/test/repo"
    type: multi-folder
    branch: develop
    author: "Test Author"
  - name: "Other Repo"
    url: "https://github.com/test/other"
    type: single-md
    author: "Other Author"
`
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "blogo.yaml")
	os.WriteFile(configPath, []byte(yaml), 0o644)
	t.Setenv("BLOGO_CONFIG", configPath)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() returned error: %v", err)
	}

	if cfg.Port != 9090 {
		t.Errorf("Port = %d, want 9090", cfg.Port)
	}
	if cfg.LogLevel != "debug" {
		t.Errorf("LogLevel = %q, want debug", cfg.LogLevel)
	}
	if cfg.FetchOnStart != false {
		t.Errorf("FetchOnStart = %v, want false", cfg.FetchOnStart)
	}
	if len(cfg.Repos) != 2 {
		t.Fatalf("Repos count = %d, want 2", len(cfg.Repos))
	}
	if cfg.Repos[0].Name != "Test Repo" {
		t.Errorf("Repos[0].Name = %q, want Test Repo", cfg.Repos[0].Name)
	}
	if cfg.Repos[0].Type != "multi-folder" {
		t.Errorf("Repos[0].Type = %q, want multi-folder", cfg.Repos[0].Type)
	}
	if cfg.Repos[0].Branch != "develop" {
		t.Errorf("Repos[0].Branch = %q, want develop", cfg.Repos[0].Branch)
	}
	if cfg.Repos[1].Branch != "main" {
		t.Errorf("Repos[1].Branch = %q, want main (default)", cfg.Repos[1].Branch)
	}
}

func TestLoadEnvOverrides(t *testing.T) {
	t.Setenv("BLOGO_CONFIG", "/nonexistent/blogo.yaml")
	t.Setenv("BLOGO_PORT", "9090")
	t.Setenv("BLOGO_CONTENT_URL", "https://example.com/test.md")
	t.Setenv("BLOGO_CONTENT_DIR", "/tmp/content")
	t.Setenv("BLOGO_FETCH_ON_START", "false")
	t.Setenv("BLOGO_LOG_LEVEL", "debug")
	t.Setenv("BLOGO_LOG_FORMAT", "json")
	t.Setenv("BLOGO_SHUTDOWN_TIMEOUT", "30s")

	origDir, _ := os.Getwd()
	tmpDir := t.TempDir()
	os.Chdir(tmpDir)
	defer os.Chdir(origDir)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() returned error: %v", err)
	}

	if cfg.Port != 9090 {
		t.Errorf("Port = %d, want 9090", cfg.Port)
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
	t.Setenv("BLOGO_CONFIG", "/nonexistent/blogo.yaml")
	t.Setenv("BLOGO_PORT", "not-a-number")

	origDir, _ := os.Getwd()
	tmpDir := t.TempDir()
	os.Chdir(tmpDir)
	defer os.Chdir(origDir)

	_, err := Load()
	if err == nil {
		t.Fatal("Load() should return error for invalid port")
	}
}

func TestLoadPortOutOfRange(t *testing.T) {
	t.Setenv("BLOGO_CONFIG", "/nonexistent/blogo.yaml")
	t.Setenv("BLOGO_PORT", "99999")

	origDir, _ := os.Getwd()
	tmpDir := t.TempDir()
	os.Chdir(tmpDir)
	defer os.Chdir(origDir)

	_, err := Load()
	if err == nil {
		t.Fatal("Load() should return error for port out of range")
	}
}

func TestIsMultiRepo(t *testing.T) {
	cfg := &Config{
		Repos: []RepoConfig{
			{Name: "A"},
			{Name: "B"},
		},
	}
	if !cfg.IsMultiRepo() {
		t.Error("IsMultiRepo() = false, want true for 2 repos")
	}

	cfg.Repos = cfg.Repos[:1]
	if cfg.IsMultiRepo() {
		t.Error("IsMultiRepo() = true, want false for 1 repo")
	}
}
