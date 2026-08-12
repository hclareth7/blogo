# Architecture Specification

## Status

| Field       | Value              |
|-------------|--------------------|
| Version     | 2.0                |
| Status      | Draft              |
| Author      | Hernando Clareth   |
| Created     | 2026-08-11         |

---

## Overview

BLOGO is a server-side rendered Go application that parses Markdown sources from multiple Git repositories, builds per-repo in-memory document indexes, and serves them as navigable HTML pages with full-text search.

```
blogo.yaml (repo list)
       │
       ▼
  Config Loader
       │
       ▼
  ┌────────────────────────────────────┐
  │  For each repo in config:          │
  │                                    │
  │  GitHub Repository (Markdown)      │
  │         │                          │
  │         ▼                          │
  │    Fetch Job (per type)            │
  │    ├── single-md: HTTP GET raw     │
  │    └── multi-folder: GitHub API    │
  │         │                          │
  │         ▼                          │
  │    content/{repo-slug}/            │
  │         │                          │
  │         ▼                          │
  │    Markdown Parser (Goldmark)      │
  │         │                          │
  │         ▼                          │
  │    RepoState (in-memory)           │
  │    ├── Document Index              │
  │    ├── Navigation Tree             │
  │    ├── Search Index (Bleve)        │
  │    └── Rendered HTML Fragments     │
  └────────────────────────────────────┘
       │
       ▼
  HTTP Server (Chi + Go Templates + HTMX)
       │
       ▼
  Browser (repo selector + per-repo nav)
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

- Load configuration from `blogo.yaml` (or fall back to env vars)
- For each configured repo: execute content fetch job, parse content, build index
- Store per-repo state in `Server.repos` map
- Initialize search index (per repo)
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
├── Title    string
├── RepoSlug string             (URL-safe identifier for the repo)
├── Sections []Section
│   ├── ID        string    (URL slug)
│   ├── Title     string    (heading text)
│   ├── Level     int       (1, 2, 3)
│   ├── Content   string    (rendered HTML)
│   ├── Children  []Section (subsections)
│   └── Order     int       (position in document)
```

**Parser methods:**
- `Parse(source []byte) (*Document, error)` — single-md repos (existing)
- `ParseMultiFolder(repoDir string) (*Document, error)` — multi-folder repos (Phase 1.2)

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
- Multi-repo state management (`repos` map, `repoList` for selector)
- Repo-scoped route resolution (`/{repo-slug}/{section}`)
- Graceful shutdown with configurable drain timeout

**Multi-repo state:**
```
Server
├── repos    map[string]*RepoState  (keyed by repo slug)
├── repoList []*RepoMeta            (ordered list for selector UI)
├── ...existing fields

RepoState
├── Doc      *parser.Document
├── Index    *parser.Index
├── NavTree  *navigation.NavTree
├── NavBld   *navigation.Builder

RepoMeta
├── Name   string
├── Slug   string
├── Author string
├── Type   string
```

### internal/config/ (Phase 1.2)

Configuration loading from `blogo.yaml`.

**Responsibilities:**
- Load and validate `blogo.yaml` config file
- Fall back to environment variables when no config file is present
- Parse repo list with type, branch, and author fields
- Config resolution order: `BLOGO_CONFIG` env var → `./blogo.yaml` → env-var fallback

### internal/content/

Content acquisition and management.

**Responsibilities:**
- Fetch job per repo type:
  - `single-md`: download raw README.md via HTTP (existing)
  - `multi-folder`: discover folder tree via GitHub API, fetch each markdown file and images
- Store content in `content/{repo-slug}/` directory
- Detect content changes (checksum comparison)
- Resolve relative image paths for `multi-folder` repos
- Trigger re-parse and re-index on content update
- Future: webhook endpoint to receive push notifications

## Route Design

| Method | Path                                   | Handler          | Description                         |
|--------|----------------------------------------|------------------|-------------------------------------|
| GET    | /                                      | Public           | Home (first repo, first section)    |
| GET    | /{repo-slug}                           | Public           | Repo home (first section)           |
| GET    | /{repo-slug}/{section}                 | Public           | Top-level section page              |
| GET    | /{repo-slug}/{section}/{subsection}    | Public           | Subsection page                     |
| GET    | /search                                | Search           | Search results page (HTMX)         |
| GET    | /api/v1/search                         | Search API       | JSON search endpoint                |
| GET    | /static/*                              | Static           | Global static assets (CSS, JS)      |
| GET    | /static/content/{repo-slug}/*          | Static           | Repo-specific images                |
| GET    | /healthz                               | Health           | Liveness probe                      |
| GET    | /readyz                                | Health           | Readiness probe                     |

When only one repo is configured (or in env-var fallback mode), routes without `{repo-slug}` prefix continue to work for backward compatibility.

## Configuration

### Config file: `blogo.yaml` (Phase 1.2)

Primary configuration via `blogo.yaml` for multi-repo support. See `specs/multi-repo.md` for full format.

**Resolution order:**
1. `BLOGO_CONFIG` env var (path to config file)
2. `./blogo.yaml` in working directory
3. Env-var fallback (legacy single-repo mode)

### Environment variables (legacy / server-level)

| Variable              | Default                  | Description                          |
|-----------------------|--------------------------|--------------------------------------|
| BLOGO_CONFIG          | (none)                   | Path to config file                  |
| BLOGO_PORT            | 8080                     | HTTP server port                     |
| BLOGO_CONTENT_URL     | (GitHub raw URL)         | URL to fetch Markdown source (legacy)|
| BLOGO_CONTENT_DIR     | ./content                | Local content directory              |
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

- Multi-stage Dockerfile: build stage (`golang:1.25-alpine`) → runtime stage (`scratch`)
- Build flags: `CGO_ENABLED=0 -ldflags="-s -w" -trimpath`
- Embed templates and static assets via `go:embed`
- Final image size: ~15 MB (uncompressed)
- Expose port 8080
- Default `USER 65534:65534` (overridden by OpenShift SCC)
- `.dockerignore` excludes specs, deploy, .git, content from build context

```bash
podman build -t quay.io/hclareth/blogo:latest .
podman push quay.io/hclareth/blogo:latest
```

### Health Checks

- `/healthz` — liveness: returns 200 if the process is running
- `/readyz` — readiness: returns 200 if content is loaded and search index is ready

### OpenShift Deployment (Kustomize)

Manifests in `deploy/k8s/`, applied via `oc apply -k deploy/k8s/`.

| Manifest           | Kind        | Description                                            |
|--------------------|-------------|--------------------------------------------------------|
| namespace.yaml     | Namespace   | `blogo` namespace                                      |
| deployment.yaml    | Deployment  | 2 replicas, RollingUpdate, probes, security hardening  |
| service.yaml       | Service     | ClusterIP on port 8080                                 |
| certificate.yaml   | Certificate | cert-manager TLS for `blogo.hclareth.space`            |
| route.yaml         | Route       | OpenShift Route, TLS edge, externalCertificate         |
| kustomization.yaml | Kustomize   | Namespace, commonLabels, image tag management           |

### Security Hardening (restricted-v2 SCC)

Pod level:
- `runAsNonRoot: true` (UID assigned by OpenShift from namespace range)
- `seccompProfile: RuntimeDefault`
- `automountServiceAccountToken: false`

Container level:
- `allowPrivilegeEscalation: false`
- `readOnlyRootFilesystem: true`
- `capabilities.drop: [ALL]`
- `seccompProfile: RuntimeDefault`

Writable volumes: `emptyDir` for `/content` and `/tmp`.

### Image Pull Secret

The container registry (quay.io) requires authentication. The pull secret is created manually (never versioned):

```bash
oc create secret docker-registry blogo-pull-secret \
  --namespace=blogo \
  --docker-server=quay.io \
  --docker-username="<ROBOT_USERNAME>" \
  --docker-password="<ROBOT_TOKEN>"
```

The Deployment references it via `spec.template.spec.imagePullSecrets`.

### TLS

- cert-manager `Certificate` resource requests a TLS certificate from ClusterIssuer `ca-cluster-issue-letsencrypt`
- Secret `blogo-tls` is referenced by the OpenShift Route via `tls.externalCertificate`
- Route enforces HTTPS with `insecureEdgeTerminationPolicy: Redirect`

### Kubernetes Compatibility

- Liveness and readiness probes mapped to health endpoints
- Graceful shutdown on SIGTERM with configurable drain timeout (`terminationGracePeriodSeconds: 30`)
- Configuration via environment variables (ConfigMap/Secret friendly)
- Resource requests/limits defined in Deployment spec

## Security Considerations

- No user input persisted (read-only application)
- HTML output escaped by default via Go templates
- Static content served with appropriate cache headers
- Security headers: CSP, X-Content-Type-Options, X-Frame-Options
- Content fetch uses HTTPS only
- Webhook endpoint (future) validates HMAC-SHA256 signatures
