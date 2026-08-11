# Vision Specification

## Status

| Field       | Value              |
|-------------|--------------------|
| Version     | 1.0                |
| Status      | Draft              |
| Author      | Hernando Clareth   |
| Created     | 2026-08-11         |

---

## Purpose

BLOGO is a modern, interactive documentation platform that transforms the System Design knowledge base by Karan Pratap Singh into a structured, searchable, and navigable learning experience.

Instead of reading a monolithic README, engineers explore concepts through structured navigation, deep links, full-text search, and an optimized reading experience.

## Problem Statement

The original System Design repository is a single large README file. While the content is excellent, the format has limitations:

- No table of contents or sidebar navigation
- No search capability
- No deep linking to specific sections
- No reading progress tracking
- No optimized reading experience (themes, responsive layout)
- No SEO-friendly URLs for individual topics

## Solution

BLOGO parses the Markdown source, builds a document index, and serves the content through a web application that provides:

1. Structured navigation (sidebar, TOC, previous/next)
2. Full-text search with instant results
3. Deep links to every section and subsection
4. Clean reading experience with light/dark themes
5. SEO-friendly URLs per topic

## Scope

### In Scope

- Rendering System Design content from Karan Pratap Singh's repository
- Markdown parsing with support for code blocks, tables, images, Mermaid diagrams
- Navigation tree generation from document headings
- Full-text search indexing and querying
- Responsive UI with theme support
- SEO routes and metadata
- Content synchronization via fetch job
- Future: webhook-based content updates on push to source main branch

### Out of Scope

- User-generated content or comments
- User authentication or accounts
- Content editing through the UI
- Multi-source content aggregation (single source only)
- Content modification or transformation beyond rendering

## Content Source

| Field          | Value                                              |
|----------------|----------------------------------------------------|
| Repository     | https://github.com/karanpratapsingh/system-design  |
| Format         | Markdown (single README.md)                        |
| License        | CC BY-NC-ND 4.0                                    |
| Strategy       | Fetch job (pull-based)                             |
| Future         | Webhook on push to main (TBD)                      |

## Content License Compliance

The source content is licensed under CC BY-NC-ND 4.0. BLOGO must:

- Preserve full attribution to Karan Pratap Singh
- Avoid commercial usage
- Not create derivative works of the content itself
- Display the original license prominently
- Link back to the original repository

BLOGO's code (the platform) is licensed separately under Apache 2.0.

## Target Audience

- Software engineers studying system design
- Engineers preparing for system design interviews
- Technical leads reviewing architecture concepts
- Anyone who prefers structured navigation over scrolling a long README

## Universe Context

BLOGO is the first project ("planet") in the **hclareth7.space** universe -- a collection of engineering projects where each project is represented as a celestial body.

```
hclareth7.space
├── BLOGO        (this project)
├── Future projects...
```

## Success Criteria

- Content renders correctly with all Markdown features (code, tables, images, Mermaid)
- Navigation tree accurately reflects document structure
- Search returns relevant results in under 100ms
- Lighthouse performance score above 90
- Pages load without requiring JavaScript (progressive enhancement)
- Content can be updated without redeploying the application

## Development Philosophy

BLOGO follows Specification Driven Development (SDD):

```
Specification → Review → Implementation → Testing → Release
```

Every feature starts with a specification before any code is written.
