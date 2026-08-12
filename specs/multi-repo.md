# Multi-Repo Support Specification

## Status

| Field       | Value              |
|-------------|--------------------|
| Version     | 1.0                |
| Status      | Draft              |
| Phase       | 1.2                |
| Author      | Hernando Clareth   |
| Created     | 2026-08-12         |

---

## Overview

BLOGO should support reading markdown content from multiple Git repositories with different organizational structures. Each repository is treated as an independent content source with its own navigation tree, and users can switch between repos via a searchable selector in the sidebar.

## Repository Types

### `single-md`

A repository with all content in a single Markdown file (typically `README.md`).

**Example:** [karanpratapsingh/system-design](https://github.com/karanpratapsingh/system-design)

```
repo/
└── README.md          # 58 H1 sections, all content here
```

**Parsing strategy:**
- Fetch single file via raw GitHub URL
- Split by H1 headings into top-level sections
- H2 headings become subsections (children)
- H3+ rendered inline within their parent section
- Skip "Table of Contents" sections

### `multi-folder`

A repository where each topic is a numbered folder containing its own `README.md` (or `Readme.md`) and optional assets like images.

**Example:** [liquidslr/system-design-notes](https://github.com/liquidslr/system-design-notes)

```
repo/
├── 01. Scaling/
│   ├── Readme.md
│   └── images/
├── 02. Back Of the Envelope Estimation/
│   ├── Readme.md
│   └── images/
├── ...
└── 24. S3-like Object Storage/
    ├── README.md
    └── images/
```

**Parsing strategy:**
- Clone or fetch the full repository tree
- Discover folders sorted by numeric prefix (natural sort)
- Each folder's markdown file becomes a top-level section
- The folder name (stripped of numeric prefix) becomes the section title
- H2 headings within each file become subsections
- Resolve relative image paths (`./images/filename.png`) to served static paths

**Detection heuristics (for future auto-detect):**
- No root-level README.md with H1 sections, OR
- Multiple numbered directories containing `.md` files

### Future types (out of scope for Phase 1.2)

| Type           | Description                                          | Example                     |
|----------------|------------------------------------------------------|-----------------------------|
| `flat-files`   | Root directory with multiple .md files, no folders   | `intro.md`, `chapter-1.md`  |
| `nested-tree`  | Hierarchical folder tree with index.md per level     | mdBook, Docusaurus sources  |

These are documented for planning purposes only. Phase 1.2 implements `single-md` (already working) and `multi-folder`.

## Configuration

### Config file: `blogo.yaml`

The repo list is complex enough that a config file is more appropriate than environment variables. BLOGO reads `blogo.yaml` from the working directory (overridable via `BLOGO_CONFIG` env var).

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

### Repo fields

| Field    | Required | Description                                     |
|----------|----------|-------------------------------------------------|
| `name`   | yes      | Display name shown in the sidebar and selector  |
| `url`    | yes      | Git repository HTTPS URL                        |
| `type`   | yes      | Repository structure type (`single-md`, `multi-folder`) |
| `branch` | no       | Branch to fetch (default: `main`)               |
| `author` | no       | Content author (displayed in metadata line)     |

### Backward compatibility

Environment variables (`BLOGO_PORT`, `BLOGO_LOG_LEVEL`, etc.) still work for server-level config. If `blogo.yaml` is absent, BLOGO falls back to legacy behavior (single repo via `BLOGO_CONTENT_URL` env var, type `single-md`).

### Config resolution order

1. `BLOGO_CONFIG` env var (path to config file)
2. `./blogo.yaml` in working directory
3. Env-var fallback (legacy single-repo mode)

## Content Fetching

### `single-md` repos

Fetch the single raw file via HTTP (existing behavior):

```
https://raw.githubusercontent.com/{owner}/{repo}/{branch}/README.md
```

### `multi-folder` repos

Use the GitHub API to discover the folder tree, then fetch each markdown file:

1. **List tree:** `GET /repos/{owner}/{repo}/git/trees/{branch}?recursive=1`
2. **Filter:** Keep only entries matching `*/README.md` or `*/Readme.md` (case-insensitive)
3. **Sort:** Natural sort by folder name (numeric prefix)
4. **Fetch each file:** Download raw content via GitHub raw URL
5. **Fetch images:** Download referenced images and serve them locally

### Storage layout

Each repo gets its own directory under `content/`:

```
content/
├── system-design/                    # slug from repo name
│   └── README.md
└── system-design-notes/              # slug from repo name
    ├── 01-scaling/
    │   ├── Readme.md
    │   └── images/
    ├── 02-back-of-the-envelope-estimation/
    │   ├── Readme.md
    │   └── images/
    └── ...
```

## Parsing

### `single-md` parser (existing)

No changes. The current `Parser.Parse(source []byte)` handles this.

### `multi-folder` parser (new)

New method: `Parser.ParseMultiFolder(repoDir string) (*Document, error)`

1. Walk the repo directory, find all `.md` files
2. Sort folders by natural order (numeric prefix)
3. For each folder:
   - Read the markdown file
   - Extract the folder title from the folder name (strip numeric prefix: `"01. Scaling"` → `"Scaling"`)
   - Parse headings: H1 becomes the section title (or use folder title if no H1), H2 becomes children
   - Resolve relative image paths to `/static/content/{repo-slug}/{folder-slug}/images/`
4. Assemble all sections into a single `Document`

### Document model changes

The `Document` struct adds a `RepoSlug` field:

```go
type Document struct {
    Title    string
    RepoSlug string      // URL-safe identifier for the repo
    Sections []*Section
}
```

## URL Routing

Each repo gets a URL prefix based on its slug:

| Repo                 | Slug                   | URLs                                    |
|----------------------|------------------------|-----------------------------------------|
| System Design        | system-design          | `/system-design/{section}`              |
| System Design Notes  | system-design-notes    | `/system-design-notes/{section}`        |

New routes:

```
/                                          → Home (first repo, first section)
/{repo-slug}                               → Repo home (first section)
/{repo-slug}/{section}                     → Section page
/{repo-slug}/{section}/{subsection}        → Subsection page
/static/content/{repo-slug}/...            → Repo-specific static assets (images)
```

The current routes (`/{section}`) continue to work as a default repo shortcut when only one repo is configured.

## UI Changes

### Sidebar header — Repo selector

Replace the static "System Design / by Author" header with a searchable autocomplete dropdown:

```
┌─────────────────────┐
│ 🔽 System Design    │  ← click to open selector
│ by Karan Pratap S.  │
├─────────────────────┤
│ 🔍 Search repos...  │  ← autocomplete input
│                     │
│ ● System Design     │  ← active repo (highlighted)
│   System Design N.  │
│   [future repos]    │
└─────────────────────┘
```

**Behavior:**
- Clicking the repo name opens a dropdown with search input
- Typing filters the repo list (client-side, Alpine.js)
- Selecting a repo navigates to `/{repo-slug}` (first section)
- The sidebar navigation tree updates to show the selected repo's sections
- Current repo name + author displayed in the sidebar header
- Dropdown closes on selection, click outside, or Escape key

### Metadata line

The author field in the metadata line uses the repo-level `author` value:

```
📖 5 min read  •  👤 {repo.author}  •  Original ↗
```

### HTMX integration

Switching repos triggers an HTMX request to `/{repo-slug}` that swaps `#main-content` and the sidebar navigation via OOB swap (same pattern as section navigation).

## Server Changes

### Multi-repo state

The server holds a map of loaded repos:

```go
type Server struct {
    repos    map[string]*RepoState  // keyed by repo slug
    repoList []*RepoMeta            // ordered list for selector
    // ... existing fields
}

type RepoState struct {
    Doc      *parser.Document
    Index    *parser.Index
    NavTree  *navigation.NavTree
    NavBld   *navigation.Builder
}

type RepoMeta struct {
    Name   string
    Slug   string
    Author string
    Type   string
}
```

### Startup flow

1. Load `blogo.yaml`
2. For each repo in config:
   a. Fetch content (based on type)
   b. Parse into `Document`
   c. Build `Index` and `NavTree`
   d. Store in `repos` map
3. Start HTTP server

## Running the Project

### With config file

```bash
# Create blogo.yaml in the project root
cat > blogo.yaml <<EOF
port: 8080
log_level: info
fetch_on_start: true

repos:
  - name: "System Design"
    url: "https://github.com/karanpratapsingh/system-design"
    type: single-md
    branch: main
    author: "Karan Pratap Singh"
EOF

go run ./cmd/blogo
```

### With environment variables (legacy, single repo)

```bash
BLOGO_CONTENT_URL="https://raw.githubusercontent.com/karanpratapsingh/system-design/main/README.md" \
go run ./cmd/blogo
```

## Migration Path

Phase 1.2 is backward compatible:
- If `blogo.yaml` exists → use multi-repo config
- If no config file → fall back to env vars (single-repo `single-md` mode, legacy behavior)
- All existing URLs (`/{section}`) work when only one repo is configured

## Out of Scope

- Repo auto-detection (type must be explicit in config)
- Webhook-based content updates (future phase)
- Cross-repo search (Phase 2)
- Git clone-based fetching (uses GitHub API/raw URLs)
- Image optimization or caching headers for repo images
