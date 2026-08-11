# Architecture Specification

## Status

| Field       | Value              |
|-------------|--------------------|
| Version     | 1.0                |
| Status      | Draft              |
| Author      | Hernando Clareth   |
| Created     | 2026-08-11         |

---

## Overview

BLOGO is a server-side rendered Go application that parses a Markdown source, builds an in-memory document index, and serves it as navigable HTML pages with full-text search.

```
GitHub Repository (Markdown Source)
       │
       ▼
  Fetch Job (HTTP pull)
       │
       ▼
  content/README.md (local copy)
       │
       ▼
  Markdown Parser (Goldmark)
       │
       ▼
  Document Index (in-memory)
       │
       ├── Navigation Tree
       ├── Search Index (Bleve)
       ├── SEO Route Map
       └── Rendered HTML Fragments
       │
       ▼
  HTTP Server (Chi + Go Templates + HTMX)
       │
       ▼
  Browser
```

## Technology Stack

| Layer           | Technology      | Rationale                                                    |
|-----------------|-----------------|--------------------------------------------------------------|
| Language        | Go 1.22+        | Single binary, fast compilation, strong stdlib, concurrency  |
| HTTP Router     | Chi             | Minimal, idiomatic, composable middleware, stdlib-compatible  |
| Templates       | Go html/template| No build step, server-side rendering, secure by default      |
| Interactivity   | HTMX            | Hypermedia-driven, no JS build pipeline, progressive enhance |
| Client behavior | Alpine.js       | Lightweight reactive behavior (theme toggle, mobile menu)    |
| CSS             | Tailwind CSS    | Utility-first, standalone CLI (no Node.js required)          |
| Markdown        | Goldmark        | CommonMark compliant, extensible, pure Go                    |
| Diagrams        | Mermaid          | Client-side rendering of architecture diagrams               |
| Search          | Bleve           | Pure Go full-text search, no external service                |
| Syntax highlight| Chroma          | Go-native syntax highlighter, integrates with Goldmark       |

## Component Architecture

### cmd/blogo/main.go

Application entrypoint. Responsibilities:

- Load configuration from environment variables
- Execute content fetch job (if enabled)
- Initialize Markdown parser and build document index
- Initialize search index
- Configure HTTP server with middleware and routes
- Start server with graceful shutdown on SIGTERM

### internal/parser/

Markdown parsing pipeline using Goldmark.

**Responsibilities:**
- Parse Markdown source into AST (Abstract Syntax Tree)
- Extract heading hierarchy (h1, h2, h3) for navigation tree
- Render Markdown to HTML fragments per section
- Handle extensions: tables, fenced code blocks (Chroma), Mermaid blocks, footnotes
- Generate URL slugs from heading text

**Key types:**
```
Document
├── Sections []Section
│   ├── ID        string    (URL slug)
│   ├── Title     string    (heading text)
│   ├── Level     int       (1, 2, 3)
│   ├── Content   string    (rendered HTML)
│   ├── Children  []Section (subsections)
│   └── Order     int       (position in document)
```

**Goldmark extensions to configure:**
- `extension.Table`
- `extension.Strikethrough`
- `extension.Footnote`
- `highlighting.NewHighlighting` (with Chroma)
- Custom Mermaid block renderer (wraps in `<pre class="mermaid">`)

### internal/navigation/

Builds the navigation tree from parsed sections.

**Responsibilities:**
- Build sidebar tree from heading hierarchy
- Generate previous/next links for sequential reading
- Build table of contents per page
- Resolve deep links (section ID to URL mapping)

See `specs/navigation.md` for full specification.

### internal/search/

Full-text search using Bleve.

**Responsibilities:**
- Build search index from document sections
- Execute search queries with ranking
- Return section-aware results (title, snippet, link)
- Rebuild index on content update

See `specs/search.md` for full specification.

### internal/renderer/

HTML rendering pipeline.

**Responsibilities:**
- Combine rendered Markdown HTML with Go templates
- Inject navigation context (sidebar, TOC, prev/next)
- Generate SEO metadata per page (title, description, Open Graph)
- Handle template composition (base layout, partials)

### internal/server/

HTTP server setup.

**Responsibilities:**
- Chi router configuration
- Middleware stack (logging, recovery, request ID, compression, cache headers)
- Route registration (pages, search API, static assets, health)
- Graceful shutdown with configurable drain timeout

### internal/content/

Content acquisition and management.

**Responsibilities:**
- Fetch job: download README.md from GitHub via HTTP
- Store content in `content/` directory
- Detect content changes (checksum comparison)
- Trigger re-parse and re-index on content update
- Future: webhook endpoint to receive push notifications

## Route Design

| Method | Path                    | Handler          | Description                      |
|--------|-------------------------|------------------|----------------------------------|
| GET    | /                       | Public           | Home / document overview         |
| GET    | /{section}              | Public           | Top-level section page           |
| GET    | /{section}/{subsection} | Public           | Subsection page                  |
| GET    | /search                 | Search           | Search results page (HTMX)      |
| GET    | /api/v1/search          | Search API       | JSON search endpoint             |
| GET    | /static/*               | Static           | CSS, JS, images                  |
| GET    | /healthz                | Health           | Liveness probe                   |
| GET    | /readyz                 | Health           | Readiness probe                  |

## Configuration

All configuration via environment variables following 12-factor principles.

| Variable              | Default                  | Description                          |
|-----------------------|--------------------------|--------------------------------------|
| BLOGO_PORT            | 8080                     | HTTP server port                     |
| BLOGO_CONTENT_URL     | (GitHub raw URL)         | URL to fetch Markdown source         |
| BLOGO_CONTENT_DIR     | ./content                | Local directory for content storage  |
| BLOGO_FETCH_ON_START  | true                     | Fetch content on application start   |
| BLOGO_FETCH_INTERVAL  | 0 (disabled)             | Periodic fetch interval (e.g., 1h)  |
| BLOGO_LOG_LEVEL       | info                     | Log level (debug, info, warn, error) |
| BLOGO_LOG_FORMAT      | json                     | Log format (json, text)              |
| BLOGO_SHUTDOWN_TIMEOUT| 10s                      | Graceful shutdown drain timeout      |

## Content Synchronization

### Strategy 1: Fetch Job (Phase 1)

Pull-based content acquisition:

```
Application Start
       │
       ├── BLOGO_FETCH_ON_START=true → Download README.md from BLOGO_CONTENT_URL
       │
       ├── Compare checksum with existing content/README.md
       │
       ├── If changed → Write to disk → Re-parse → Re-index
       │
       └── If BLOGO_FETCH_INTERVAL > 0 → Schedule periodic fetch
```

- HTTP GET to GitHub raw content URL
- SHA-256 checksum comparison to detect changes
- Atomic file write (write to temp, rename)
- Re-parse and re-index without server restart

### Strategy 2: Webhook (TBD - Future)

Push-based content update:

```
GitHub Push to main
       │
       ▼
  GitHub Webhook (POST /api/v1/webhook)
       │
       ▼
  Validate signature (HMAC-SHA256)
       │
       ▼
  Fetch updated content
       │
       ▼
  Re-parse → Re-index
```

- Endpoint: `POST /api/v1/webhook`
- GitHub webhook secret validation
- Filter: only trigger on pushes that modify the README.md
- Implementation deferred to a later phase

## Concurrency Model

- Single Go binary, single process
- HTTP requests handled concurrently by Go's net/http (goroutine per request)
- Document index is read-heavy: use `sync.RWMutex` for concurrent reads, exclusive lock on content update
- Search index: Bleve handles its own concurrency
- Content fetch job runs in a separate goroutine with a ticker

## Error Handling

- Application fails fast on startup if content cannot be fetched (when `BLOGO_FETCH_ON_START=true`)
- Runtime fetch failures are logged but do not crash the server (serve stale content)
- Invalid Markdown sections are logged and skipped, not fatal
- All errors logged with structured context (slog)

## Observability

| Signal   | Tool      | Endpoint/Output |
|----------|-----------|-----------------|
| Logging  | slog      | stdout (JSON)   |
| Metrics  | Prometheus| /metrics        |
| Health   | Custom    | /healthz, /readyz |

Metrics to expose:
- `blogo_http_requests_total` (method, path, status)
- `blogo_http_request_duration_seconds` (histogram)
- `blogo_content_fetch_total` (status: success, failure)
- `blogo_content_last_fetch_timestamp`
- `blogo_search_queries_total`
- `blogo_search_query_duration_seconds` (histogram)
- `blogo_document_sections_total` (gauge)

## Deployment

### Container Image

- Multi-stage Dockerfile: build stage (Go compiler) → runtime stage (distroless/scratch)
- Embed templates and static assets via `go:embed`
- Target image size: under 30 MB compressed
- Expose port 8080
- Run as non-root user

### Health Checks

- `/healthz` — liveness: returns 200 if the process is running
- `/readyz` — readiness: returns 200 if content is loaded and search index is ready

### Kubernetes Compatibility

- Liveness and readiness probes mapped to health endpoints
- Graceful shutdown on SIGTERM with configurable drain timeout
- Configuration via environment variables (ConfigMap/Secret friendly)
- Resource requests/limits documented in Helm chart values

## Security Considerations

- No user input persisted (read-only application)
- HTML output escaped by default via Go templates
- Static content served with appropriate cache headers
- Security headers: CSP, X-Content-Type-Options, X-Frame-Options
- Content fetch uses HTTPS only
- Webhook endpoint (future) validates HMAC-SHA256 signatures
