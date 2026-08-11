# Navigation Specification

## Status

| Field       | Value              |
|-------------|--------------------|
| Version     | 1.0                |
| Status      | Draft              |
| Author      | Hernando Clareth   |
| Created     | 2026-08-11         |

---

## Overview

BLOGO generates all navigation structures from the Markdown heading hierarchy. There is no manual configuration -- the navigation tree, table of contents, deep links, and sequential navigation are all derived from the parsed document.

## Heading Hierarchy Mapping

The source document uses Markdown headings to structure content:

| Markdown | Level | Navigation Role        | URL Pattern               |
|----------|-------|------------------------|---------------------------|
| `# H1`   | 1     | Top-level section      | `/{slug}`                 |
| `## H2`  | 2     | Subsection             | `/{parent-slug}/{slug}`   |
| `### H3` | 3     | Sub-subsection         | `/{parent-slug}/{slug}#anchor` |

Headings below H3 are rendered inline but do not generate navigation entries.

## Slug Generation

Slugs are generated from heading text:

1. Convert to lowercase
2. Replace spaces with hyphens
3. Remove non-alphanumeric characters (except hyphens)
4. Collapse consecutive hyphens
5. Trim leading/trailing hyphens

Examples:

| Heading                     | Slug                        |
|-----------------------------|-----------------------------|
| "Load Balancing"            | `load-balancing`            |
| "What is a CDN?"            | `what-is-a-cdn`             |
| "N-Tier Architecture"       | `n-tier-architecture`       |
| "CAP Theorem"               | `cap-theorem`               |

Slug uniqueness is enforced within the document. If a duplicate is detected, append `-2`, `-3`, etc.

## Sidebar Navigation

The sidebar displays the full navigation tree, visible on all pages.

### Structure

```
Sidebar
├── Section A (H1)
│   ├── Subsection A.1 (H2)
│   ├── Subsection A.2 (H2)
│   └── Subsection A.3 (H2)
├── Section B (H1)
│   ├── Subsection B.1 (H2)
│   └── Subsection B.2 (H2)
└── Section C (H1)
```

### Behavior

- The current section is highlighted in the sidebar
- Sections are collapsible (expand/collapse children)
- On mobile: sidebar is hidden by default, toggled via hamburger menu (Alpine.js)
- On desktop: sidebar is always visible in a fixed left column
- Sidebar scroll position is preserved during navigation

### Data Model

```
NavTree
├── Items []NavItem
│   ├── ID       string     (slug)
│   ├── Title    string     (heading text)
│   ├── URL      string     (full path)
│   ├── Level    int        (1 or 2)
│   ├── Active   bool       (current page)
│   ├── Children []NavItem  (subsections)
│   └── Order    int        (position)
```

## Table of Contents (TOC)

The TOC is not a separate column. Instead, it is integrated into the content reading experience via the sidebar navigation and in-page heading anchors. The sidebar serves as the primary navigation structure, and H3 anchors within pages provide in-page deep linking.

## Previous / Next Navigation

Sequential reading is supported via previous/next links at the bottom of each page.

### Behavior

- Previous and next are determined by document order (flattened heading tree)
- First page has no "Previous" link
- Last page has no "Next" link
- Links show the title of the target section
- Navigation follows H1 sections in order, then H2 within each H1

### Example

```
Section order: [Introduction, IP, TCP, UDP, DNS, Load Balancing, ...]

On "TCP" page:
  ← Previous: IP
  → Next: UDP
```

## Deep Linking

Every section and subsection has a permanent URL.

| Content Level | URL Example                        |
|---------------|------------------------------------|
| H1 Section    | `/load-balancing`                  |
| H2 Subsection | `/load-balancing/algorithms`       |
| H3 Anchor     | `/load-balancing/algorithms#round-robin` |

### Requirements

- All URLs are stable (derived from heading text, which is stable in the source)
- URLs are human-readable and SEO-friendly
- Deep links are shareable
- Server returns 404 for unknown slugs with a helpful "section not found" page

## Reading Progress

### Behavior

- A thin progress bar at the top of the page indicates scroll position within the current page
- Progress is calculated as: `scrollTop / (scrollHeight - clientHeight) * 100`
- Implemented client-side with minimal JS (Alpine.js `x-data`)
- No persistence (progress resets on page reload)

## Keyboard Navigation

| Key        | Action                            |
|------------|-----------------------------------|
| `j` / `↓`  | Next section (in sidebar)         |
| `k` / `↑`  | Previous section (in sidebar)     |
| `Enter`    | Navigate to highlighted section   |
| `/`        | Focus search input                |
| `Escape`   | Close search / unfocus            |

Keyboard shortcuts are active only when no input element is focused.

## URL Structure

```
/                                    → Home (overview / first section)
/{section-slug}                      → Top-level section (H1)
/{section-slug}/{subsection-slug}    → Subsection (H2)
/search?q={query}                    → Search results page
/static/*                            → Static assets
/healthz                             → Liveness probe
/readyz                              → Readiness probe
```

## SEO

Each page generates:

- `<title>` — `{Section Title} | BLOGO`
- `<meta name="description">` — first 160 characters of the section content (plain text)
- `<link rel="canonical">` — full URL of the page
- Open Graph tags: `og:title`, `og:description`, `og:url`, `og:type=article`
- `<link rel="prev">` and `<link rel="next">` for sequential discovery

## Edge Cases

- Empty heading text: skip, do not generate navigation entry
- Heading with only special characters: generate slug as `section-{order}`
- Very long heading: truncate slug at 80 characters
- Duplicate slugs: append numeric suffix (`-2`, `-3`)
- Missing H1 (document starts with H2): treat first heading as root regardless of level
- Single section document: sidebar shows one item, no prev/next links
