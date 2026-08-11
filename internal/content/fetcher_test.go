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

func TestFetchNewContent(t *testing.T) {
	t.Parallel()
	content := "# System Design\n\nHello world"
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(content))
	}))
	defer srv.Close()

	dir := t.TempDir()
	f := NewFetcher(srv.URL, dir, testLogger())

	changed, err := f.Fetch(context.Background())
	if err != nil {
		t.Fatalf("Fetch() error: %v", err)
	}
	if !changed {
		t.Error("Fetch() returned changed=false, want true for new content")
	}

	data, err := f.ReadContent()
	if err != nil {
		t.Fatalf("ReadContent() error: %v", err)
	}
	if string(data) != content {
		t.Errorf("ReadContent() = %q, want %q", string(data), content)
	}
}

func TestFetchUnchangedContent(t *testing.T) {
	t.Parallel()
	content := "# System Design\n\nSame content"
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(content))
	}))
	defer srv.Close()

	dir := t.TempDir()
	f := NewFetcher(srv.URL, dir, testLogger())

	if _, err := f.Fetch(context.Background()); err != nil {
		t.Fatalf("first Fetch() error: %v", err)
	}

	changed, err := f.Fetch(context.Background())
	if err != nil {
		t.Fatalf("second Fetch() error: %v", err)
	}
	if changed {
		t.Error("Fetch() returned changed=true, want false for same content")
	}
}

func TestFetchChangedContent(t *testing.T) {
	t.Parallel()
	callCount := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		if callCount == 1 {
			w.Write([]byte("version 1"))
		} else {
			w.Write([]byte("version 2"))
		}
	}))
	defer srv.Close()

	dir := t.TempDir()
	f := NewFetcher(srv.URL, dir, testLogger())

	if _, err := f.Fetch(context.Background()); err != nil {
		t.Fatalf("first Fetch() error: %v", err)
	}

	changed, err := f.Fetch(context.Background())
	if err != nil {
		t.Fatalf("second Fetch() error: %v", err)
	}
	if !changed {
		t.Error("Fetch() returned changed=false, want true for different content")
	}

	data, err := f.ReadContent()
	if err != nil {
		t.Fatalf("ReadContent() error: %v", err)
	}
	if string(data) != "version 2" {
		t.Errorf("ReadContent() = %q, want %q", string(data), "version 2")
	}
}

func TestFetchServerError(t *testing.T) {
	t.Parallel()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	dir := t.TempDir()
	f := NewFetcher(srv.URL, dir, testLogger())

	_, err := f.Fetch(context.Background())
	if err == nil {
		t.Fatal("Fetch() should return error for 500 response")
	}
}

func TestFetchAtomicWrite(t *testing.T) {
	t.Parallel()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("content"))
	}))
	defer srv.Close()

	dir := t.TempDir()
	f := NewFetcher(srv.URL, dir, testLogger())

	if _, err := f.Fetch(context.Background()); err != nil {
		t.Fatalf("Fetch() error: %v", err)
	}

	tmp := filepath.Join(dir, filename+".tmp")
	if _, err := os.Stat(tmp); !os.IsNotExist(err) {
		t.Error("temp file should not exist after successful fetch")
	}
}

func TestReadContentMissing(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	f := NewFetcher("http://unused", dir, testLogger())

	_, err := f.ReadContent()
	if err == nil {
		t.Fatal("ReadContent() should return error when file doesn't exist")
	}
}
