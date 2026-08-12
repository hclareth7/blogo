# BLOGO Repo Standards

Version: 1.0
Phase: 1.2.1

## Purpose

Defines the standard repository structure for each repo type supported by BLOGO. Any repository owner who wants their content served through BLOGO must follow one of these templates.

---

## Supported Types

| Type           | Config value    | Description                                    |
|----------------|-----------------|------------------------------------------------|
| Single Markdown| `single-md`     | All content in a single `README.md` at root    |
| Multi Folder   | `multi-folder`  | Numbered folders, each with its own `README.md`|

---

## Type: `single-md`

A single `README.md` at the repository root contains all documentation content.

### Required structure

```text
repo/
└── README.md
```

### README.md format

```markdown
# Document Title

Introduction text (skipped as content, used as document title).

# Section One

Content for section one.

## Subsection A

Content for subsection A.

## Subsection B

Content for subsection B.

# Section Two

Content for section two.
```

### Rules

| Rule                        | Detail                                                    |
|-----------------------------|-----------------------------------------------------------|
| Title                       | First `# H1` becomes the document title                  |
| Sections                    | Each subsequent `# H1` becomes a top-level section        |
| Subsections                 | Each `## H2` inside a section becomes a child subsection  |
| Deep headings               | `### H3` and deeper are rendered inline (not navigable)   |
| Table of Contents           | `# Table of Contents` heading is auto-skipped             |
| Images                      | External URLs only (no local images)                      |
| Code blocks                 | Fenced code blocks with language tags for syntax highlight|

### Example

Repository: [karanpratapsingh/system-design](https://github.com/karanpratapsingh/system-design)

---

## Type: `multi-folder`

Each topic lives in its own numbered folder with a `README.md` and optional `images/` directory.

### Required structure

```text
repo/
├── README.md              (optional) Root introduction, shown as first item
├── 01. Topic Name/
│   ├── README.md           Chapter content
│   └── images/             (optional) Chapter images
│       ├── diagram.png
│       └── flow.png
├── 02. Another Topic/
│   ├── README.md
│   └── images/
│       └── overview.png
└── NN. Last Topic/
    └── README.md
```

### Folder naming

| Pattern              | Example                           | Result          |
|----------------------|-----------------------------------|-----------------|
| `NN. Name`           | `01. Scaling`                     | Order: 1        |
| `NN-Name`            | `02-Estimation`                   | Order: 2        |
| `NN Name`            | `03 Framework`                    | Order: 3        |
| No prefix            | `Appendix`                        | Order: 0 (first)|

- Numeric prefix determines display order in the sidebar
- The prefix is stripped from the display title
- If the folder's `README.md` contains an `# H1` heading, that heading overrides the display title

### README.md format (per folder)

```markdown
# Chapter Title

Introduction paragraph.

## Section One

Content with inline images:

<img src="./images/diagram.png" width="400" />

Or markdown images:

![diagram](./images/diagram.png)

## Section Two

More content.
```

### Image references

| Format                              | Supported |
|--------------------------------------|-----------|
| `![alt](./images/file.png)`          | Yes       |
| `![alt](images/file.png)`            | Yes       |
| `<img src="./images/file.png" />`    | Yes       |
| `<img src="images/file.png" />`      | Yes       |
| `![alt](https://external.com/img)`   | Yes       |
| `![alt](../other-folder/images/x)`   | No        |

### Root README.md

- Optional
- If present, displayed as the first item in the sidebar navigation
- Title derived from its `# H1` heading, or "Introduction" if none
- Does not support local images (no `images/` at root level)

### Rules

| Rule                        | Detail                                                    |
|-----------------------------|-----------------------------------------------------------|
| Title                       | First `# H1` in a folder's README becomes the section title |
| Subsections                 | `## H2` headings become navigable child subsections       |
| Deep headings               | `### H3` and deeper are rendered inline                   |
| Ordering                    | Numeric prefix in folder name determines sidebar order    |
| Images                      | Must be in `images/` subdirectory within the same folder  |
| Image paths                 | Must use `./images/` or `images/` relative paths          |

### Example

Repository: [liquidslr/system-design-notes](https://github.com/liquidslr/system-design-notes)

---

## Configuration

Add the repository to `blogo.yaml`:

```yaml
repos:
  - name: "Display Name"
    url: "https://github.com/owner/repo"
    type: single-md        # or multi-folder
    branch: main
    author: "Author Name"
```

### Field reference

| Field    | Required | Description                                |
|----------|----------|--------------------------------------------|
| `name`   | Yes      | Display name in the repo selector          |
| `url`    | Yes      | GitHub repository URL                      |
| `type`   | Yes      | `single-md` or `multi-folder`              |
| `branch` | Yes      | Git branch to fetch content from           |
| `author` | Yes      | Attribution shown in footer and sidebar    |
