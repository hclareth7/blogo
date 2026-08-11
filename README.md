# 🌍 BLOGO

> A modern and interactive System Design documentation platform built with Go.

BLOGO transforms the excellent System Design knowledge base created by Karan Pratap Singh into a faster, more searchable, and more navigable learning experience.

Instead of reading a massive README from top to bottom, BLOGO allows engineers to explore concepts through structured navigation, deep links, full-text search, and an optimized reading experience.

BLOGO is the first planet inside the **hclareth7.space** universe.

---

## Vision

BLOGO exists for two reasons:

1. Create a better learning experience for System Design.
2. Serve as a playground for learning:
   - Go
   - Concurrency
   - Search indexing
   - Specification Driven Development
   - Open Source architecture

The long-term vision of **hclareth7.space** is to build a universe of engineering projects where every project is represented as a celestial body.

```
hclareth7.space

├── 🌍 BLOGO
├── ⭐ Future Project
├── 🌌 Future Platform
└── 🚀 More to come
```

---

## Features

### Documentation Experience

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

## Content Source

BLOGO uses the public System Design repository maintained by Karan Pratap Singh as its content source.

Original project:

https://github.com/karanpratapsingh/system-design

---

## Attribution

This project does not claim ownership of the original educational content.

All educational material belongs to its original author:

**Karan Pratap Singh**

BLOGO focuses on improving:

- Navigation
- Searchability
- Accessibility
- Reading experience

while preserving attribution and respecting the original license.

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
GitHub README
       │
       ▼
Markdown Parser
       │
       ▼
Document Index
       │
       ├── Navigation Tree
       ├── Search Index
       ├── SEO Routes
       └── Content Renderer
       │
       ▼
BLOGO UI
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
├── specs/
│   ├── vision.md
│   ├── architecture.md
│   ├── navigation.md
│   ├── search.md
│   └── ui.md
│
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

```bash
git clone https://github.com/G-HNL7/blogo.git

cd blogo

go run ./cmd/blogo
```

Application will be available at:

```text
http://localhost:8080
```

---

## Contributing

Contributions, ideas, and feedback are welcome.

Please read:

- CONTRIBUTING.md
- CODE_OF_CONDUCT.md

before opening issues or pull requests.

---

## Acknowledgments

- Karan Pratap Singh for the original System Design content
- The Go community
- Open Source maintainers worldwide

---

## Author

Hernando Clareth

https://hclareth7.space

---

## Status

🚧 Early Development