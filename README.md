# <img src="web/static/favicon.svg" alt="BLOGO" width="32" height="32"> BLOGO

> A multi-repo documentation platform that transforms Markdown repositories into searchable, navigable learning experiences.

BLOGO reads multiple Git repositories with different organizational structures — from a single giant README to folder-per-topic layouts — and serves them as a unified, interactive documentation site with structured navigation, deep links, and full-text search.

BLOGO is the first planet inside the **hclareth.space** universe.

---

## Vision

BLOGO exists for two reasons:

1. Create a better learning experience from community-maintained Markdown repositories.
2. Serve as a playground for learning:
   - Go
   - Concurrency
   - Search indexing
   - Specification Driven Development
   - Open Source architecture

The long-term vision of **hclareth.space** is to build a universe of engineering projects where every project is represented as a celestial body.

```
hclareth.space

├── 🌍 BLOGO
├── ⭐ Future Project
├── 🌌 Future Platform
└── 🚀 More to come
```

---

## Features

### Documentation Experience

- Multi-repo support with searchable repo selector
- Repository types: `single-md` (one big README) and `multi-folder` (folder-per-topic)
- Clean reading experience
- Responsive layout
- Light and Dark themes
- SEO-friendly URLs
- Deep linking
- Previous / Next navigation
- Reading progress tracking

### Search

- Full-text search
- Instant results
- Keyboard shortcuts
- Section-aware navigation

### Content Navigation

- Auto-generated sidebar
- Auto-generated table of contents
- URL per section
- URL per subsection

### Markdown Support

- Syntax highlighted code blocks
- Tables
- Images
- Mermaid diagrams
- External links

---

## Content Sources

BLOGO supports multiple Git repositories as content sources. The default configuration includes:

| Repository | Author | Type |
|---|---|---|
| [karanpratapsingh/system-design](https://github.com/karanpratapsingh/system-design) | Karan Pratap Singh | `single-md` |
| [liquidslr/system-design-notes](https://github.com/liquidslr/system-design-notes) | liquidslr | `multi-folder` |

---

## Attribution

This project does not claim ownership of any educational content.

All educational material belongs to its original authors.

BLOGO focuses on improving navigation, searchability, accessibility, and reading experience while preserving attribution and respecting original licenses.

---

## License Considerations

The source content is licensed under:

**Creative Commons Attribution-NonCommercial-NoDerivatives 4.0 International (CC BY-NC-ND 4.0)**

BLOGO aims to:

- Preserve attribution
- Avoid commercial usage
- Respect original licensing terms

Users should review the original repository and license before reusing content.

---

## Architecture

```text
blogo.yaml (repo list)
       │
       ▼
  ┌─────────────────────────┐
  │  For each repo:         │
  │  GitHub Repo → Fetcher  │
  │  → Markdown Parser      │
  │  → Document Index       │
  │  → Navigation Tree      │
  └─────────────────────────┘
       │
       ▼
  HTTP Server (Chi + HTMX)
       │
       ▼
  Browser (repo selector UI)
```

---

## Technology Stack

### Backend

- Go
- Chi Router
- Go Templates

### Frontend

- HTMX
- Tailwind CSS
- Alpine.js

### Search

- Bleve

### Documentation

- Markdown
- Goldmark
- Mermaid

---

## Project Structure

```text
blogo/

├── cmd/
│   └── blogo/
│
├── internal/
│   ├── config/
│   ├── parser/
│   ├── search/
│   ├── navigation/
│   ├── renderer/
│   └── server/
│
├── web/
│   ├── templates/
│   ├── static/
│   └── assets/
│
├── content/
│   └── README.md
│
├── deploy/
│   └── k8s/              # OpenShift manifests (Kustomize)
│
├── specs/
│   ├── vision.md
│   ├── architecture.md
│   ├── multi-repo.md
│   ├── navigation.md
│   ├── search.md
│   └── ui.md
│
├── blogo.yaml            # Multi-repo configuration
├── Dockerfile
└── README.md
```

---

## Development Philosophy

BLOGO follows a Specification Driven Development approach.

Every feature starts with a specification before implementation.

```text
Specification
      ↓
Review
      ↓
Implementation
      ↓
Testing
      ↓
Release
```

---

## Roadmap

### Phase 1

- Markdown parser
- Documentation rendering
- Sidebar navigation

### Phase 1.1

- Dockerfile (multi-stage, scratch)
- OpenShift manifests (Kustomize)
- cert-manager TLS (blogo.hclareth.space)

### Phase 1.2

- Multi-repo support (`single-md` + `multi-folder` types)
- YAML config file (`blogo.yaml`) for repo list
- Searchable repo selector in sidebar
- Per-repo URL routing (`/{repo-slug}/{section}`)

### Phase 1.2.1

- Chapter ordering by numeric prefix for `multi-folder` repos
- Root README.md support (shown as first sidebar item)
- HTML `<img>` tag image path rewriting
- Dynamic footer copyright per repo author
- Sidebar nav truncation with tooltips for long titles
- Repo selector arrow icon sizing
- Standard repo templates documentation (`specs/repo-standards.md`)

### Phase 1.2.2

- Generic flat navigation (removed repo-specific section groups and icons)
- Dot bullet indicators for all sidebar sections
- Repo selector redesign with initial letter avatars, "Current source" / "Switch source repository" labels

### Phase 2

- Search engine
- Deep linking
- SEO routes

### Phase 3

- Reading progress
- Keyboard navigation
- Enhanced UX

### Phase 4

- Advanced indexing
- Performance optimization
- Content synchronization

---

## Running Locally

### With config file (multi-repo)

```bash
git clone https://github.com/hclareth7/blogo.git
cd blogo

# Create blogo.yaml
cat > blogo.yaml <<EOF
port: 8080
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
EOF

go run ./cmd/blogo
```

### With environment variables (legacy, single repo)

```bash
git clone https://github.com/hclareth7/blogo.git
cd blogo
go run ./cmd/blogo
```

Application will be available at `http://localhost:8080`.

## Deployment

BLOGO deploys to OpenShift using a multi-stage scratch-based container image.

```bash
# Build image
podman build -t quay.io/hclareth/blogo:latest .

# Push to registry
podman push quay.io/hclareth/blogo:latest

# Deploy to OpenShift
oc apply -k deploy/k8s/
```

| Component     | Value                            |
|---------------|----------------------------------|
| Image         | quay.io/hclareth/blogo           |
| Platform      | OpenShift (restricted-v2 SCC)    |
| TLS           | cert-manager + ClusterIssuer     |
| URL           | https://blogo.hclareth.space     |
| Manifests     | Kustomize (`deploy/k8s/`)        |

---

## Contributing

Contributions, ideas, and feedback are welcome.

Please read:

- CONTRIBUTING.md
- CODE_OF_CONDUCT.md

before opening issues or pull requests.

---

## Acknowledgments

- Karan Pratap Singh for the System Design content
- liquidslr for the System Design Notes
- The Go community
- Open Source maintainers worldwide

---

## Author

Hernando Clareth

https://hclareth.space

---

## Status

🚧 Early Development