package content

import (
	"context"
	"crypto/sha256"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

const filename = "README.md"

type Fetcher struct {
	url    string
	dir    string
	client *http.Client
	logger *slog.Logger
}

func NewFetcher(url, dir string, logger *slog.Logger) *Fetcher {
	return &Fetcher{
		url: url,
		dir: dir,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
		logger: logger,
	}
}

func (f *Fetcher) Fetch(ctx context.Context) (bool, error) {
	f.logger.Info("fetching content", "url", f.url)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, f.url, nil)
	if err != nil {
		return false, fmt.Errorf("creating request: %w", err)
	}

	resp, err := f.client.Do(req)
	if err != nil {
		return false, fmt.Errorf("fetching content: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, fmt.Errorf("reading response body: %w", err)
	}

	newHash := checksum(data)

	existing, err := f.ReadContent()
	if err == nil && checksum(existing) == newHash {
		f.logger.Info("content unchanged, skipping write")
		return false, nil
	}

	if err := os.MkdirAll(f.dir, 0o755); err != nil {
		return false, fmt.Errorf("creating content directory: %w", err)
	}

	dest := filepath.Join(f.dir, filename)
	tmp := dest + ".tmp"

	if err := os.WriteFile(tmp, data, 0o644); err != nil {
		return false, fmt.Errorf("writing temp file: %w", err)
	}

	if err := os.Rename(tmp, dest); err != nil {
		os.Remove(tmp)
		return false, fmt.Errorf("renaming temp file: %w", err)
	}

	f.logger.Info("content updated", "bytes", len(data))
	return true, nil
}

func (f *Fetcher) ReadContent() ([]byte, error) {
	return os.ReadFile(filepath.Join(f.dir, filename))
}

func (f *Fetcher) ContentPath() string {
	return filepath.Join(f.dir, filename)
}

func checksum(data []byte) string {
	h := sha256.Sum256(data)
	return fmt.Sprintf("%x", h)
}
