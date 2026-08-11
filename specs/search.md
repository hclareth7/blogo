# Search Specification

## Status

| Field       | Value              |
|-------------|--------------------|
| Version     | 1.0                |
| Status      | Draft              |
| Phase       | Phase 2            |
| Author      | Hernando Clareth   |
| Created     | 2026-08-11         |

---

## Overview

BLOGO provides full-text search across all System Design content using Bleve, a pure Go search library. Search is section-aware -- results link directly to the relevant section, not just the document.

## Search Engine

| Aspect        | Choice    | Rationale                                         |
|---------------|-----------|---------------------------------------------------|
| Library       | Bleve     | Pure Go, no external service, embeddable           |
| Index storage | In-memory | Content is small enough; rebuilt on startup/update |
| Analyzer      | Standard  | English language analyzer with stemming            |

## Indexing

### What Gets Indexed

Each document section (H1, H2, H3) is indexed as a separate search document.

**Index document schema:**

| Field      | Type   | Indexed | Stored | Boost | Description                        |
|------------|--------|---------|--------|-------|------------------------------------|
| id         | string | no      | yes    | -     | Section slug (unique identifier)   |
| title      | text   | yes     | yes    | 2.0   | Heading text (boosted relevance)   |
| body       | text   | yes     | yes    | 1.0   | Section content (plain text, no HTML) |
| url        | string | no      | yes    | -     | Full URL path to the section       |
| parent     | string | yes     | yes    | 0.5   | Parent section title               |
| level      | number | no      | yes    | -     | Heading level (1, 2, 3)            |
| order      | number | no      | yes    | -     | Position in document               |

### Index Build Process

```
Parsed Sections
      │
      ▼
Strip HTML from rendered content → plain text
      │
      ▼
Create Bleve index (in-memory)
      │
      ▼
Batch-index all sections
      │
      ▼
Index ready for queries
```

### Re-indexing

When content is updated (fetch job detects changes):

1. Build new index from updated sections
2. Swap old index reference with new index (atomic pointer swap)
3. Old index is garbage collected

No downtime -- queries against the old index continue until swap completes.

## Query Interface

### API Endpoint

```
GET /api/v1/search?q={query}&limit={limit}&offset={offset}
```

| Parameter | Type   | Default | Description                     |
|-----------|--------|---------|---------------------------------|
| q         | string | (required) | Search query                 |
| limit     | int    | 20      | Maximum results to return       |
| offset    | int    | 0       | Pagination offset               |

### Response Format

```json
{
  "query": "load balancing",
  "total": 5,
  "took_ms": 12,
  "results": [
    {
      "id": "load-balancing",
      "title": "Load Balancing",
      "snippet": "...distributes incoming network traffic across multiple <mark>servers</mark>...",
      "url": "/load-balancing",
      "parent": "System Design",
      "score": 0.95
    }
  ]
}
```

### HTMX Endpoint

```
GET /search?q={query}
```

Returns an HTML fragment (not full page) for HTMX partial update of the search results container.

## Search Features

### Query Processing

- Tokenize query into terms
- Apply same analyzer as indexing (lowercase, stemming)
- Search across `title` and `body` fields
- Title matches rank higher (2x boost)
- Support quoted phrases: `"consistent hashing"` matches exact phrase

### Result Ranking

Bleve's default TF-IDF (Term Frequency-Inverse Document Frequency) scoring with:

- Title field boost: 2.0x
- Phrase matches rank higher than individual term matches
- Results sorted by relevance score (descending)

### Snippet Generation

- Extract fragment around the first match in the body
- Maximum snippet length: 200 characters
- Highlight matched terms with `<mark>` tags
- If match is in title only, show first 200 characters of body as snippet

### Performance Targets

| Metric                         | Target     |
|--------------------------------|------------|
| Query latency (p95)            | < 50ms     |
| Index build time               | < 2s       |
| Index memory footprint         | < 50 MB    |

## Keyboard Shortcuts

| Key       | Action                               |
|-----------|--------------------------------------|
| `/`       | Focus search input                   |
| `Escape`  | Clear search input and close results |
| `↓`       | Navigate to next result              |
| `↑`       | Navigate to previous result          |
| `Enter`   | Go to selected result                |

## UI Behavior

See `specs/ui.md` for detailed search UI specification.

Summary:
- Search input in header with placeholder "Search... (press /)"
- Results via HTMX partial update, debounced at 300ms
- Dropdown on desktop, full-page on mobile
- Maximum 20 results displayed
- Loading spinner during query execution

## Edge Cases

- Empty query: return no results (do not show all sections)
- Single character query: require minimum 2 characters
- Very long query (>200 chars): truncate to 200 characters
- No results: display "No results found for '{query}'" message
- Special characters in query: escape for Bleve query syntax
- Concurrent re-index during query: query completes on old index, next query uses new index
