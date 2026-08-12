package server

import (
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/hclareth7/blogo/internal/config"
	"github.com/hclareth7/blogo/internal/navigation"
	"github.com/hclareth7/blogo/internal/parser"
	"github.com/hclareth7/blogo/internal/renderer"
	"github.com/hclareth7/blogo/web/static"
	"github.com/hclareth7/blogo/web/templates"
)

func testServer(t *testing.T) *Server {
	t.Helper()

	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError}))

	doc := &parser.Document{
		Title:    "System Design",
		RepoSlug: "system-design",
		Author:   "Karan Pratap Singh",
		Sections: []*parser.Section{
			{
				ID: "intro", Title: "Introduction", Level: 1,
				Content: "<p>Welcome to system design.</p>", RawText: "Welcome to system design.",
				WordCount: 5, Order: 0,
			},
			{
				ID: "ip", Title: "IP", Level: 1,
				Content: "<p>Internet Protocol.</p>", RawText: "Internet Protocol.",
				WordCount: 2, Order: 1,
				Children: []*parser.Section{
					{
						ID: "versions", Title: "Versions", Level: 2,
						Content: "<p>IPv4 and IPv6.</p>", RawText: "IPv4 and IPv6.",
						WordCount: 3, Order: 2,
					},
				},
			},
		},
	}

	index := parser.NewIndex(doc)
	navBuilder := navigation.NewBuilderForRepo(logger, "system-design")
	navTree := navBuilder.BuildTree(doc.Sections)

	rend, err := renderer.New(templates.FS, logger)
	if err != nil {
		t.Fatalf("renderer.New() error: %v", err)
	}

	cfg := &config.Config{
		Port: 0,
		Repos: []config.RepoConfig{
			{Name: "System Design", Type: "single-md", Author: "Karan Pratap Singh"},
		},
	}

	repos := map[string]*RepoState{
		"system-design": {
			Doc:     doc,
			Index:   index,
			NavTree: navTree,
			NavBld:  navBuilder,
			Meta: renderer.RepoMeta{
				Name:   "System Design",
				Slug:   "system-design",
				Author: "Karan Pratap Singh",
				Type:   "single-md",
			},
		},
	}
	repoOrder := []string{"system-design"}

	return New(cfg, repos, repoOrder, rend, static.FS, logger)
}

func TestHealthz(t *testing.T) {
	t.Parallel()
	srv := testServer(t)
	r := srv.Routes()

	req := httptest.NewRequest("GET", "/healthz", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("healthz status = %d, want %d", w.Code, http.StatusOK)
	}
}

func TestReadyz(t *testing.T) {
	t.Parallel()
	srv := testServer(t)
	r := srv.Routes()

	req := httptest.NewRequest("GET", "/readyz", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("readyz status = %d, want %d", w.Code, http.StatusOK)
	}
}

func TestHomePage(t *testing.T) {
	t.Parallel()
	srv := testServer(t)
	r := srv.Routes()

	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("home status = %d, want %d", w.Code, http.StatusOK)
	}
	if ct := w.Header().Get("Content-Type"); ct != "text/html; charset=utf-8" {
		t.Errorf("Content-Type = %q, want text/html", ct)
	}
}

func TestRepoHomePage(t *testing.T) {
	t.Parallel()
	srv := testServer(t)
	r := srv.Routes()

	req := httptest.NewRequest("GET", "/system-design", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("repo home status = %d, want %d", w.Code, http.StatusOK)
	}
}

func TestRepoSectionPage(t *testing.T) {
	t.Parallel()
	srv := testServer(t)
	r := srv.Routes()

	req := httptest.NewRequest("GET", "/system-design/ip", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("repo section status = %d, want %d", w.Code, http.StatusOK)
	}
}

func TestRepoSubsectionPage(t *testing.T) {
	t.Parallel()
	srv := testServer(t)
	r := srv.Routes()

	req := httptest.NewRequest("GET", "/system-design/ip/versions", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("repo subsection status = %d, want %d", w.Code, http.StatusOK)
	}
}

func TestNotFound(t *testing.T) {
	t.Parallel()
	srv := testServer(t)
	r := srv.Routes()

	req := httptest.NewRequest("GET", "/nonexistent", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("404 status = %d, want %d", w.Code, http.StatusNotFound)
	}
}

func TestStaticAssets(t *testing.T) {
	t.Parallel()
	srv := testServer(t)
	r := srv.Routes()

	req := httptest.NewRequest("GET", "/static/css/styles.css", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("static CSS status = %d, want %d", w.Code, http.StatusOK)
	}
}

func TestHTMXPartialResponse(t *testing.T) {
	t.Parallel()
	srv := testServer(t)
	r := srv.Routes()

	req := httptest.NewRequest("GET", "/system-design/ip", nil)
	req.Header.Set("HX-Request", "true")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("HTMX section status = %d, want %d", w.Code, http.StatusOK)
	}

	body := w.Body.String()
	if !strings.Contains(body, "sidebar-nav") {
		t.Error("HTMX response missing sidebar OOB swap")
	}
	if strings.Contains(body, "<!DOCTYPE html>") {
		t.Error("HTMX response should not contain full HTML document")
	}
}

func TestSecurityHeaders(t *testing.T) {
	t.Parallel()
	srv := testServer(t)
	r := srv.Routes()

	req := httptest.NewRequest("GET", "/healthz", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	headers := map[string]string{
		"X-Content-Type-Options": "nosniff",
		"X-Frame-Options":       "DENY",
		"Referrer-Policy":       "strict-origin-when-cross-origin",
	}

	for key, want := range headers {
		got := w.Header().Get(key)
		if got != want {
			t.Errorf("header %s = %q, want %q", key, got, want)
		}
	}
}
