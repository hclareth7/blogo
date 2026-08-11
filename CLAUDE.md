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
└── deploy/              # Dockerfile, Helm, Terraform
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
make docker           # Build container image
make clean            # Remove build artifacts
```

## Configuration

All configuration via environment variables (12-factor):

| Variable              | Default         | Description                     |
|-----------------------|-----------------|---------------------------------|
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

## Roadmap

- Phase 1: Markdown parser, document rendering, sidebar navigation
- Phase 2: Search engine, deep linking, SEO routes
- Phase 3: Reading progress, keyboard navigation, enhanced UX
- Phase 4: Advanced indexing, performance optimization, content sync
