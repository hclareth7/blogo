package content

import (
	"context"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func testLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError}))
}

func TestFetchSingleMD(t *testing.T) {
	t.Parallel()
	body := "# System Design\n\nHello world"
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(body))
	}))
	defer srv.Close()

	dir := t.TempDir()
	f := NewFetcher(testLogger())
	f.client = srv.Client()

	err := f.FetchSingleMD(context.Background(), srv.URL, "main", dir, "test-repo")
	if err != nil {
		t.Fatalf("FetchSingleMD() error: %v", err)
	}

	data, err := f.ReadSingleMD(dir, "test-repo")
	if err != nil {
		t.Fatalf("ReadSingleMD() error: %v", err)
	}
	if string(data) != body {
		t.Errorf("ReadSingleMD() = %q, want %q", string(data), body)
	}
}

func TestRepoSlug(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		want string
	}{
		{"System Design", "system-design"},
		{"System Design Notes", "system-design-notes"},
		{"My Repo!", "my-repo"},
	}
	for _, tt := range tests {
		got := RepoSlug(tt.name)
		if got != tt.want {
			t.Errorf("RepoSlug(%q) = %q, want %q", tt.name, got, tt.want)
		}
	}
}

func TestSlugifyFolder(t *testing.T) {
	t.Parallel()
	tests := []struct {
		input string
		want  string
	}{
		{"01. Scaling", "scaling"},
		{"02. Back Of the Envelope Estimation", "back-of-the-envelope-estimation"},
		{"24. S3-like Object Storage", "s3-like-object-storage"},
		{"Simple Name", "simple-name"},
	}
	for _, tt := range tests {
		got := slugifyFolder(tt.input)
		if got != tt.want {
			t.Errorf("slugifyFolder(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestParseGitHubURL(t *testing.T) {
	t.Parallel()
	tests := []struct {
		url        string
		wantOwner  string
		wantRepo   string
	}{
		{"https://github.com/karanpratapsingh/system-design", "karanpratapsingh", "system-design"},
		{"https://github.com/liquidslr/system-design-notes", "liquidslr", "system-design-notes"},
		{"https://github.com/owner/repo.git", "owner", "repo"},
	}
	for _, tt := range tests {
		owner, repo := parseGitHubURL(tt.url)
		if owner != tt.wantOwner || repo != tt.wantRepo {
			t.Errorf("parseGitHubURL(%q) = (%q, %q), want (%q, %q)", tt.url, owner, repo, tt.wantOwner, tt.wantRepo)
		}
	}
}

func TestAtomicWrite(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	dest := filepath.Join(dir, "test.txt")

	if err := atomicWrite(dest, []byte("hello")); err != nil {
		t.Fatalf("atomicWrite() error: %v", err)
	}

	data, err := os.ReadFile(dest)
	if err != nil {
		t.Fatalf("ReadFile() error: %v", err)
	}
	if string(data) != "hello" {
		t.Errorf("file content = %q, want %q", string(data), "hello")
	}

	tmp := dest + ".tmp"
	if _, err := os.Stat(tmp); !os.IsNotExist(err) {
		t.Error("temp file should not exist after successful write")
	}
}
