package content

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"
)

type Fetcher struct {
	client *http.Client
	logger *slog.Logger
}

func NewFetcher(logger *slog.Logger) *Fetcher {
	return &Fetcher{
		client: &http.Client{Timeout: 60 * time.Second},
		logger: logger,
	}
}

func (f *Fetcher) FetchSingleMD(ctx context.Context, repoURL, branch, contentDir, slug string) error {
	var rawURL string
	if strings.HasPrefix(repoURL, "https://raw.") || strings.HasPrefix(repoURL, "http://127.0.0.1") || strings.HasPrefix(repoURL, "http://localhost") {
		rawURL = repoURL
	} else {
		rawURL = githubRawURL(repoURL, branch, "README.md")
	}
	f.logger.Info("fetching single-md", "url", rawURL, "slug", slug)

	data, err := f.httpGet(ctx, rawURL)
	if err != nil {
		return fmt.Errorf("fetching %s: %w", rawURL, err)
	}

	dir := filepath.Join(contentDir, slug)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("creating dir %s: %w", dir, err)
	}

	dest := filepath.Join(dir, "README.md")

	existing, readErr := os.ReadFile(dest)
	if readErr == nil && checksum(existing) == checksum(data) {
		f.logger.Info("content unchanged", "slug", slug)
		return nil
	}

	if err := atomicWrite(dest, data); err != nil {
		return err
	}

	f.logger.Info("content updated", "slug", slug, "bytes", len(data))
	return nil
}

func (f *Fetcher) FetchMultiFolder(ctx context.Context, repoURL, branch, contentDir, slug string) error {
	owner, repo := parseGitHubURL(repoURL)
	if owner == "" || repo == "" {
		return fmt.Errorf("cannot parse GitHub owner/repo from URL: %s", repoURL)
	}

	f.logger.Info("fetching multi-folder tree", "owner", owner, "repo", repo, "branch", branch)

	treeURL := fmt.Sprintf("https://api.github.com/repos/%s/%s/git/trees/%s?recursive=1", owner, repo, branch)
	treeData, err := f.httpGet(ctx, treeURL)
	if err != nil {
		return fmt.Errorf("fetching tree: %w", err)
	}

	var tree githubTree
	if err := json.Unmarshal(treeData, &tree); err != nil {
		return fmt.Errorf("parsing tree: %w", err)
	}

	rootMDPattern := regexp.MustCompile(`(?i)^readme\.md$`)
	mdPattern := regexp.MustCompile(`(?i)^[^/]+/readme\.md$`)
	imgPattern := regexp.MustCompile(`(?i)^[^/]+/images/.+$`)

	var rootMD string
	var mdFiles []string
	var imgFiles []string

	for _, entry := range tree.Tree {
		if entry.Type == "blob" {
			if rootMDPattern.MatchString(entry.Path) {
				rootMD = entry.Path
			} else if mdPattern.MatchString(entry.Path) {
				mdFiles = append(mdFiles, entry.Path)
			} else if imgPattern.MatchString(entry.Path) {
				imgFiles = append(imgFiles, entry.Path)
			}
		}
	}

	sort.Slice(mdFiles, func(i, j int) bool {
		return naturalLess(mdFiles[i], mdFiles[j])
	})

	repoDir := filepath.Join(contentDir, slug)
	if err := os.RemoveAll(repoDir); err != nil {
		return fmt.Errorf("cleaning repo dir: %w", err)
	}
	if err := os.MkdirAll(repoDir, 0o755); err != nil {
		return fmt.Errorf("creating repo dir: %w", err)
	}

	if rootMD != "" {
		rawURL := githubRawURL(repoURL, branch, rootMD)
		data, err := f.httpGet(ctx, rawURL)
		if err != nil {
			f.logger.Warn("skipping root README", "error", err)
		} else {
			dest := filepath.Join(repoDir, "README.md")
			if err := atomicWrite(dest, data); err != nil {
				return err
			}
			f.logger.Debug("fetched root README")
		}
	}

	for _, mdPath := range mdFiles {
		rawURL := githubRawURL(repoURL, branch, mdPath)
		data, err := f.httpGet(ctx, rawURL)
		if err != nil {
			f.logger.Warn("skipping file", "path", mdPath, "error", err)
			continue
		}

		folderName := filepath.Dir(mdPath)
		folderSlug := slugifyFolder(folderName)
		localDir := filepath.Join(repoDir, folderSlug)
		if err := os.MkdirAll(localDir, 0o755); err != nil {
			return fmt.Errorf("creating folder dir: %w", err)
		}

		dest := filepath.Join(localDir, "README.md")
		if err := atomicWrite(dest, data); err != nil {
			return err
		}

		f.logger.Debug("fetched", "path", mdPath, "folder", folderSlug)
	}

	for _, imgPath := range imgFiles {
		rawURL := githubRawURL(repoURL, branch, imgPath)
		data, err := f.httpGet(ctx, rawURL)
		if err != nil {
			f.logger.Warn("skipping image", "path", imgPath, "error", err)
			continue
		}

		parts := strings.SplitN(imgPath, "/", 2)
		if len(parts) < 2 {
			continue
		}
		folderSlug := slugifyFolder(parts[0])
		localPath := filepath.Join(repoDir, folderSlug, parts[1])
		if err := os.MkdirAll(filepath.Dir(localPath), 0o755); err != nil {
			return fmt.Errorf("creating image dir: %w", err)
		}
		if err := atomicWrite(localPath, data); err != nil {
			return err
		}
	}

	f.logger.Info("multi-folder fetch complete", "slug", slug, "folders", len(mdFiles), "images", len(imgFiles))
	return nil
}

func (f *Fetcher) ReadSingleMD(contentDir, slug string) ([]byte, error) {
	return os.ReadFile(filepath.Join(contentDir, slug, "README.md"))
}

func (f *Fetcher) httpGet(ctx context.Context, url string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := f.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %d from %s", resp.StatusCode, url)
	}

	return io.ReadAll(resp.Body)
}

type githubTree struct {
	Tree []githubTreeEntry `json:"tree"`
}

type githubTreeEntry struct {
	Path string `json:"path"`
	Type string `json:"type"`
}

func parseGitHubURL(url string) (owner, repo string) {
	url = strings.TrimSuffix(url, ".git")
	url = strings.TrimSuffix(url, "/")

	parts := strings.Split(url, "/")
	if len(parts) < 2 {
		return "", ""
	}
	return parts[len(parts)-2], parts[len(parts)-1]
}

func githubRawURL(repoURL, branch, path string) string {
	owner, repo := parseGitHubURL(repoURL)
	return fmt.Sprintf("https://raw.githubusercontent.com/%s/%s/%s/%s", owner, repo, branch, path)
}

func slugifyFolder(name string) string {
	name = strings.ToLower(name)
	name = strings.ReplaceAll(name, " ", "-")
	name = regexp.MustCompile(`[^a-z0-9-]`).ReplaceAllString(name, "")
	name = regexp.MustCompile(`-+`).ReplaceAllString(name, "-")
	name = strings.Trim(name, "-")
	return name
}

func naturalLess(a, b string) bool {
	return extractNum(a) < extractNum(b)
}

func extractNum(path string) int {
	parts := strings.SplitN(filepath.Base(filepath.Dir(path)), ".", 2)
	if len(parts) == 0 {
		return 0
	}
	n := 0
	for _, c := range parts[0] {
		if c >= '0' && c <= '9' {
			n = n*10 + int(c-'0')
		} else {
			break
		}
	}
	return n
}

func atomicWrite(dest string, data []byte) error {
	tmp := dest + ".tmp"
	if err := os.WriteFile(tmp, data, 0o644); err != nil {
		return fmt.Errorf("writing temp file %s: %w", tmp, err)
	}
	if err := os.Rename(tmp, dest); err != nil {
		os.Remove(tmp)
		return fmt.Errorf("renaming %s: %w", tmp, err)
	}
	return nil
}

func checksum(data []byte) string {
	h := sha256.Sum256(data)
	return fmt.Sprintf("%x", h)
}

func RepoSlug(name string) string {
	s := strings.ToLower(name)
	s = strings.ReplaceAll(s, " ", "-")
	s = regexp.MustCompile(`[^a-z0-9-]`).ReplaceAllString(s, "")
	s = regexp.MustCompile(`-+`).ReplaceAllString(s, "-")
	s = strings.Trim(s, "-")
	return s
}
