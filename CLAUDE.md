# BLOGO

A modern, interactive System Design documentation platform built with Go.
Transforms the System Design knowledge base by Karan Pratap Singh into a searchable, navigable learning experience.

## Development Philosophy

This project follows Specification Driven Development (SDD):
Specification -> Review -> Implementation -> Testing -> Release

Every feature starts with a spec in `specs/` before implementation.

## Tech Stack

| Layer           | Technology        |
|-----------------|-------------------|
| Language        | Go 1.22+          |
| Module          | github.com/hclareth7/blogo |
| Router          | Chi               |
| Templates       | Go html/template  |
| Interactivity   | HTMX + Alpine.js  |
| CSS             | Tailwind CSS (standalone CLI) |
| Markdown        | Goldmark + Chroma |
| Search          | Bleve             |
| Diagrams        | Mermaid           |
| Logging         | slog              |
| Metrics         | Prometheus        |

## Project Structure

```
blogo/
├── cmd/blogo/           # Application entrypoint
├── internal/
│   ├── config/          # Config file loading (blogo.yaml)
│   ├── parser/          # Markdown parsing (Goldmark)
│   ├── search/          # Full-text search (Bleve)
│   ├── navigation/      # Sidebar, TOC, prev/next
│   ├── renderer/        # HTML rendering pipeline
│   ├── server/          # HTTP server, middleware, routes
│   └── content/         # Content fetch job
├── web/
│   ├── templates/       # Go HTML templates
│   ├── static/          # CSS, JS, images
│   └── assets/          # Source assets (pre-build)
├── content/             # Fetched Markdown source
├── specs/               # Feature specifications (SDD)
├── deploy/
│   └── k8s/             # OpenShift/K8s manifests (Kustomize)
├── Dockerfile           # Multi-stage scratch-based image
└── .dockerignore        # Build context exclusions
```

## Build Commands

```bash
make dev              # Run with hot reload (air)
make build            # Build static binary
make test             # Run all tests
make lint             # Run golangci-lint
make fmt              # Format code (gofumpt)
make fetch-content    # Fetch content from source repo
make css              # Compile Tailwind CSS
make docker           # Build container image (podman)
make clean            # Remove build artifacts
```

## Container & Deployment

```bash
# Build image
podman build -t quay.io/hclareth/blogo:latest .

# Push to registry
podman push quay.io/hclareth/blogo:latest

# Deploy to OpenShift (Kustomize)
oc apply -k deploy/k8s/
```

## Configuration

Primary configuration via `blogo.yaml` (multi-repo). Fallback to environment variables for single-repo mode.

### Config file (`blogo.yaml`)

```yaml
port: 8080
log_level: info
log_format: json
fetch_on_start: true

repos:
  - name: "System Design"
    url: "https://github.com/karanpratapsingh/system-design"
    type: single-md
    branch: main
    author: "Karan Pratap Singh"

  - name: "System Design Notes"
    url: "https://github.com/liquidslr/system-design-notes"
    type: multi-folder
    branch: main
    author: "liquidslr"
```

### Repo types

| Type           | Description                                          |
|----------------|------------------------------------------------------|
| `single-md`    | Single README.md with all content (H1 per section)   |
| `multi-folder` | Numbered folders, each with its own README.md        |

### Environment variables (legacy fallback)

| Variable              | Default         | Description                     |
|-----------------------|-----------------|---------------------------------|
| BLOGO_CONFIG          | ./blogo.yaml    | Path to config file             |
| BLOGO_PORT            | 8080            | HTTP server port                |
| BLOGO_CONTENT_URL     | (GitHub raw)    | Markdown source URL             |
| BLOGO_CONTENT_DIR     | ./content       | Local content directory         |
| BLOGO_FETCH_ON_START  | true            | Fetch content on startup        |
| BLOGO_LOG_LEVEL       | info            | Log level                       |
| BLOGO_LOG_FORMAT      | json            | Log format (json/text)          |

## Coding Conventions

- Follow standard Go project layout (`cmd/`, `internal/`)
- Use `slog` for structured logging (no fmt.Println in production code)
- All SQL/queries use parameterized statements
- HTML output escaped by default via Go templates
- Tests use stdlib `testing` + testify for assertions
- Tests run with `t.Parallel()` by default
- No global mutable state; pass dependencies via constructor injection

## Content Source

- Repository: https://github.com/karanpratapsingh/system-design
- License: CC BY-NC-ND 4.0 (content only)
- BLOGO code license: Apache 2.0
- Always preserve attribution to Karan Pratap Singh

## Deployment

| Component        | Value                                |
|------------------|--------------------------------------|
| Container image  | quay.io/hclareth/blogo               |
| Container runtime| Podman                               |
| Base image       | scratch (multi-stage, CGO_ENABLED=0) |
| Platform         | OpenShift (restricted-v2 SCC)        |
| TLS              | cert-manager (ClusterIssuer)         |
| URL              | https://blogo.hclareth.space         |
| Manifests        | deploy/k8s/ (Kustomize)             |

### Pull secret

The image registry requires a pull secret (not versioned in git):

```bash
oc create secret docker-registry blogo-pull-secret \
  --namespace=blogo \
  --docker-server=quay.io \
  --docker-username="<ROBOT_USERNAME>" \
  --docker-password="<ROBOT_TOKEN>"
```

## Roadmap

- Phase 1: Markdown parser, document rendering, sidebar navigation
- Phase 1.1: Dockerfile, OpenShift manifests (Kustomize), cert-manager TLS
- Phase 1.2: Multi-repo support (single-md + multi-folder), config file, repo selector UI
- Phase 2: Search engine, deep linking, SEO routes
- Phase 3: Reading progress, keyboard navigation, enhanced UX
- Phase 4: Advanced indexing, performance optimization, content sync
