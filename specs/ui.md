# UI Specification

## Status

| Field       | Value              |
|-------------|--------------------|
| Version     | 3.0                |
| Status      | Implemented        |
| Author      | Hernando Clareth   |
| Created     | 2026-08-11         |
| Updated     | 2026-08-11         |
| Mockup      | `specs/mockups/ui-layout.png` |

---

## Overview

BLOGO's UI is server-side rendered HTML enhanced with HTMX for partial page updates and Alpine.js for client-side behavior (theme toggle, mobile menu, scroll tracking). The UI must be fully functional without JavaScript -- HTMX and Alpine.js are progressive enhancements.

The design follows a two-column layout (sidebar + content) with no separate TOC column. Navigation between sections uses HTMX partial rendering to avoid full page reloads -- only the content area and sidebar active state are swapped via out-of-band (OOB) updates.

## Layout

### Two-Column Layout (Desktop >= 768px)

```
┌──────────────────────────────────────────────────────────────┐
│  Header: [icon] BLOGO    [🔍 Search anything...  ⌘K]  [⚙/🌙]│
├─────────────────┬────────────────────────────────────────────┤
│                 │                                            │
│  System Design  │  Breadcrumb: Home > Introduction           │
│  by K.P. Singh  │                                            │
│                 │  # Introduction                            │
│  ● Introduction │                                            │
│    Getting...   │  📖 5 min read • 👤 Author • Original ↗    │
│                 │                                            │
│  📖 Fundamentals│  [Content body -- max-w-5xl, ml/mr 14rem] │
│    What is...   │                                            │
│    Why System...│  ┌─ Note ─────────────────────────────┐    │
│    How to...    │  │ ℹ The content is sourced from...   │    │
│    Design...    │  └─────────────────────────────────────┘    │
│    Common...    │                                            │
│                 │  ## What is System Design?                  │
│  🌐 Network     │  [body text...]                            │
│    OSI Model    │                                            │
│    DNS...       │  ## Why System Design?                      │
│    HTTP/HTTPS   │  • bullet list...                           │
│    Load Bal...  │                                            │
│    CDN...       │  ┌──────────────┬──────────────────┐       │
│                 │  │ ← Previous   │        Next →     │       │
│  💻 Compute     │  │ No previous  │ Getting Started   │       │
│    Servers      │  └──────────────┴──────────────────┘       │
│    Containers   │                                            │
│    Virtualiz... │                                            │
│                 │                                            │
│  🗄 Storage      │                                            │
│    Databases... │                                            │
│    SQL...       │                                            │
│    NoSQL...     │                                            │
│    CAP Theorem  │                                            │
├─────────────────┴────────────────────────────────────────────┤
│ 🪐 BLOGO  Planet ID: BLG-001  │ Open-source...  │ © K.P.S   │
│    hclareth.space ↗           │ Made with ❤ Go  │ CC BY-NC  │
└──────────────────────────────────────────────────────────────┘
```

### Content Area Spacing

- Content max-width: `max-w-5xl` (1024px)
- Main element margins: `margin-left: 14rem`, `margin-right: 14rem` (md breakpoint)
- Internal padding: `px-8 py-8`
- Content is not centered (`mx-auto` removed); it flows naturally from the sidebar edge with balanced right margin

### Single-Column Layout (Mobile < 768px)

```
┌────────────────────────┐
│  Header (hamburger,    │
│  search, theme)        │
├────────────────────────┤
│  Breadcrumb            │
│                        │
│  Content               │
│  (full width)          │
│                        │
├────────────────────────┤
│  Prev / Next Nav       │
├────────────────────────┤
│  Footer                │
└────────────────────────┘
```

Sidebar is hidden, toggled via hamburger menu (slide-in overlay).

## Responsive Breakpoints

| Breakpoint | Name    | Layout        | Sidebar    |
|------------|---------|---------------|------------|
| < 768px    | Mobile  | Single column | Overlay    |
| >= 768px   | Desktop | Two columns   | Visible    |

## Header

```
┌──────────────────────────────────────────────────────┐
│  [🚀] BLOGO       [🔍 Search anything...  ⌘K]  [⚙/🌙] │
└──────────────────────────────────────────────────────┘
```

Components:
- **Logo**: rocket/planet icon + "BLOGO" text -- links to home
- **Search input**: centered, with magnifying glass icon, placeholder "Search anything...", `⌘K` keyboard shortcut badge
- **Theme toggle**: gear icon (light mode) / moon icon (dark mode)
- **Hamburger menu**: mobile only, replaces logo area, toggles sidebar overlay

## Sidebar

### Structure

```
┌─────────────────────┐
│ System Design        │
│ by Karan Pratap Singh│
│                      │
│ ● Introduction    ·  │  ← active: blue bg + dot indicator
│   Getting Started    │
│                      │
│ 📖 Fundamentals    ▾ │  ← collapsible group with icon
│   What is System D?  │
│   Why System Design? │
│   How to Approach?   │
│   Design Principles  │
│   Common Pitfalls    │
│                      │
│ 🌐 Network         ▾ │
│   OSI Model          │
│   DNS (Domain Name)  │
│   HTTP / HTTPS       │
│   Load Balancer      │
│   CDN (Content Del.) │
│                      │
│ 💻 Compute          ▾ │
│   Servers            │
│   Containers         │
│   Virtualization     │
│                      │
│ 🗄 Storage           ▾ │
│   Databases Overview │
│   SQL Databases      │
│   NoSQL Databases    │
│   CAP Theorem        │
└─────────────────────┘
```

### Section Group Icons

Each top-level category has an icon:

| Group          | Icon                  | Description           |
|----------------|-----------------------|-----------------------|
| Fundamentals   | Book (📖)             | Foundational concepts |
| Network        | Globe (🌐)            | Network topics        |
| Compute        | Server/Monitor (💻)   | Compute resources     |
| Storage        | Database (🗄)          | Storage topics        |

Icons are rendered as SVG or icon font (Lucide/Heroicons), not emoji. The mockup uses emoji for illustration only.

### Active Section Indicator

- Active section has a **blue/accent background** with rounded corners
- A small **dot indicator** (●) appears to the right of the active section name
- Active state uses `bg-blue-500/10` (light) / `bg-blue-500/20` (dark) with `text-blue-600` / `text-blue-400`

### Sidebar Header

- Title: "System Design" (bold)
- Subtitle: "by Karan Pratap Singh" (muted text, smaller)

### Collapsible Groups

- Click the group header or chevron to expand/collapse children
- Chevron rotates: ▸ (collapsed) → ▾ (expanded)
- Expanded state persisted in `sessionStorage`
- Groups default to expanded on first visit

### Sidebar Width

- Desktop: ~260px fixed
- Scrollable independently of content

### Mobile Sidebar

- Hidden by default
- Hamburger icon in header toggles slide-in from left
- Semi-transparent backdrop overlay
- Close on: backdrop click, Escape key, navigation to a page
- Trap focus within sidebar when open (accessibility)

## Breadcrumbs

Displayed at the top of the content area, below the header.

```
Home > Introduction
Home > Fundamentals > What is System Design?
```

- "Home" always links to the root page
- Current page (last item) is not a link, displayed as plain text
- Separator: `>` character
- Muted text color, smaller font size

## Content Metadata Line

Below the page title, a metadata line provides context:

```
📖 5 min read  •  👤 Karan Pratap Singh  •  Original README ↗
```

| Element            | Description                                    |
|--------------------|------------------------------------------------|
| Reading time       | Clock icon + estimated minutes (word count / 200 WPM) |
| Author             | Person icon + "Karan Pratap Singh"             |
| Original README    | Link to the source section in the GitHub repo, external link icon |

- Items separated by bullet (•)
- Muted text color, smaller font size

## Callout / Note Blocks

Styled admonition blocks for important notes:

```
┌──────────────────────────────────────────┐
│ ℹ  Note                                  │
│                                          │
│ The content is sourced from the original │
│ README of the system-design repository   │
│ by Karan Pratap Singh.                   │
└──────────────────────────────────────────┘
```

### Styles

| Type    | Icon | Border/BG Color              |
|---------|------|------------------------------|
| Note    | ℹ    | Blue border-left, light blue bg |
| Warning | ⚠    | Yellow border-left, light yellow bg |
| Tip     | 💡   | Green border-left, light green bg |

- Left border: 4px solid accent color
- Background: subtle tint of the accent color
- Title: bold, same color as border
- Body: normal text
- Rounded corners on right side
- Dark theme: darker tint backgrounds

### Markdown Syntax

Callouts are triggered by blockquotes with a specific format:

```markdown
> **Note**
> The content is sourced from...
```

The parser detects blockquotes starting with bold "Note", "Warning", or "Tip" and renders them as styled callout blocks.

## Theme Support

### Light Theme

- Page background: white (`#ffffff`)
- Content text: dark gray (`#1a1a2e`)
- Sidebar background: light gray (`#f8f9fa`)
- Sidebar text: dark gray
- Header background: white with subtle bottom border
- Active nav item: blue tinted background (`#eff6ff`)
- Links: blue (`#2563eb`)
- Code blocks: light gray background
- Callout backgrounds: subtle blue/yellow/green tint
- Footer background: light gray

### Dark Theme

- Page background: deep navy (`#0a0f1a` / `rgb(10, 15, 26)`)
- Content text: light gray (`#e2e8f0`)
- Sidebar background: slightly lighter navy (`#0d1220` / `rgb(13, 18, 32)`)
- Header background: matches sidebar (`#0d1220`)
- Footer background: matches sidebar (`#0d1220`)
- Border color: subtle navy (`#141c2e` / `rgb(20, 28, 46)`)
- Active nav item: blue tinted background (`bg-blue-500/10`) with `text-blue-400`
- Links: light blue (`#60a5fa`)
- Code blocks: dark background (`slate-950`)
- Callout backgrounds: darker tinted variants

### Implementation

- Theme preference stored in `localStorage`
- On first visit: respect `prefers-color-scheme` media query
- Toggle switches between `light` and `dark` class on `<html>` element
- Tailwind CSS `dark:` variant for all theme-dependent styles
- No flash of wrong theme (inline script in `<head>` reads preference before render)

## Content Rendering

### Typography

- Font: system font stack (`-apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, ...`)
- Code font: monospace stack (`"JetBrains Mono", "Fira Code", Consolas, monospace`)
- Body text: 16px / 1.7 line height
- H1 (page title): 2.25rem, bold, extra bottom margin
- H2: 1.5rem, bold
- H3: 1.25rem, semibold
- Max content width: 1024px (`max-w-5xl`), left-aligned with balanced margins

### Code Blocks

- Syntax highlighting via Chroma (server-side rendered)
- Language label displayed in top-right corner of code block
- Copy button on hover (Alpine.js + Clipboard API)
- Horizontal scroll for long lines (no wrapping)
- Rounded corners, subtle border
- Inline code: `code` styled with background tint and rounded corners

### Tables

- Full width within content area
- Horizontal scroll on overflow (mobile)
- Alternating row colors
- Sticky header row

### Images

- Max width: 100% of content area
- Centered with optional caption
- Lazy loading (`loading="lazy"`)
- Responsive sizing

### Mermaid Diagrams

- Rendered client-side by Mermaid.js
- Wrapped in `<pre class="mermaid">` by the Markdown parser
- Respects current theme (light/dark)
- Fallback: display raw Mermaid code if JS is disabled

### Links

- External links: open in new tab (`target="_blank" rel="noopener noreferrer"`)
- External links display a small arrow icon (↗) after the text
- Internal links: navigate within BLOGO (same tab)
- Link color: blue (`#2563eb` light / `#60a5fa` dark)

## Search UI

### Input

- Located in the header, centered
- Magnifying glass icon on the left
- Placeholder: "Search anything..."
- `⌘K` badge displayed on the right side of the input (macOS) / `Ctrl+K` (other OS)
- Keyboard shortcut `⌘K` / `Ctrl+K` opens/focuses search
- `Escape` clears and blurs

### Results

- Displayed via HTMX partial update (no full page reload)
- Results appear below the search input as a dropdown overlay (desktop) or full page (mobile)
- Each result shows: section title, snippet with highlighted match, link
- Maximum 20 results displayed
- "No results found" message when empty
- Debounced input: 300ms delay before sending request

### HTMX Integration

```html
<input type="search"
       name="q"
       placeholder="Search anything..."
       hx-get="/search"
       hx-trigger="input changed delay:300ms, search"
       hx-target="#search-results"
       hx-indicator="#search-spinner" />
```

## Reading Progress Bar

- Thin bar (3px) at the very top of the viewport, above the header
- Color: accent color (matches link color per theme)
- Width: percentage of page scrolled
- Implemented with Alpine.js `@scroll.window`
- Only visible on content pages (not home/search)

## Previous / Next Navigation

Displayed at the bottom of the content area, before the footer.

```
┌───────────────────────┬───────────────────────┐
│  ← Previous           │              Next →    │
│  No previous          │    Getting Started     │
└───────────────────────┴───────────────────────┘
```

- Two cards side by side with subtle border
- **Previous** (left): left arrow + "Previous" label + section title below
- **Next** (right): "Next" label + right arrow + section title below (blue/accent color)
- First page: Previous shows "No previous" (muted, non-clickable)
- Last page: Next shows "No next" (muted, non-clickable)
- Cards have hover effect (subtle background change)
- On mobile: stack vertically

## Footer

Three-column footer with the hclareth.space universe theme:

```
┌──────────────────────────────────────────────────────────────┐
│ 🪐                        │                    │             │
│ You are exploring          │ Open-source        │ Source      │
│ BLOGO                      │ documentation      │ content ©   │
│ Planet ID: BLG-001         │ platform.          │ Karan Pratap│
│ hclareth.space ↗          │ Made with ❤ and Go.│ Singh       │
│                            │                    │ Licensed    │
│                            │                    │ under CC    │
│                            │                    │ BY-NC-ND 4.0│
│                            │                    │         [GH]│
└──────────────────────────────────────────────────────────────┘
```

### Left Column — Identity

- Planet/space icon (rendered, not emoji)
- "You are exploring" (muted, small)
- **BLOGO** (bold, larger)
- "Planet ID: BLG-001" (muted, monospace)
- "hclareth.space" with external link icon -- links to hclareth.space

### Center Column — Description

- "Open-source documentation platform."
- "Made with ❤️ and Go." (Go is a link to golang.org)

### Right Column — Attribution

- "Source content © Karan Pratap Singh"
- "Licensed under CC BY-NC-ND 4.0"
- GitHub icon (links to BLOGO repository)

### Footer Styling

- **Position**: fixed to the bottom of the viewport (`position: fixed; bottom: 0`)
- **Height**: 60px
- Background: matches sidebar (`#0d1220` in dark mode, light gray in light mode)
- Top border: subtle separator line (`#141c2e` in dark mode)
- Main content has `pb-[60px]` to prevent overlap
- On mobile: columns stack vertically, centered

## HTMX Client-Side Navigation

Section navigation uses HTMX to swap only the content area instead of reloading the full page.

### How It Works

1. Sidebar and prev/next links carry `hx-get`, `hx-target="#main-content"`, `hx-push-url="true"`, and `hx-swap="innerHTML scroll:window:top"`
2. The server detects the `HX-Request: true` header and returns the `content-fragment` template instead of the full `page.html`
3. The fragment includes:
   - A `<title>` tag (HTMX updates the document title)
   - The article content (breadcrumbs, heading, metadata, prose, prev/next)
   - A sidebar OOB (Out-of-Band) swap (`hx-swap-oob="innerHTML"` on `#sidebar-nav`) to update the active section indicator
4. Regular navigation (direct URL, browser refresh) still serves the full page with `base.html`

### Fallback

All links retain standard `href` attributes, so navigation works without JavaScript -- HTMX is a progressive enhancement.

## Favicon

- SVG favicon (`/static/favicon.svg`)
- Blue planet design matching the BLOGO space/planet theme
- Referenced via `<link rel="icon" type="image/svg+xml" href="/static/favicon.svg">`

## Custom Scrollbars

Both the sidebar and main content area use styled scrollbars that match the overall UI theme.

### Behavior

- **Visibility**: scrollbars are invisible by default, appearing only on hover over the scrollable container (auto-reveal pattern)
- **Width**: 6px (thin profile)
- **Shape**: fully rounded (`border-radius: 9999px`)
- **Track**: transparent (no visible track)

### Colors

| State                      | Light Mode                       | Dark Mode                        |
|----------------------------|----------------------------------|----------------------------------|
| Thumb (on container hover) | `rgb(148 163 184 / 0.3)` (slate-400) | `rgb(100 116 139 / 0.3)` (slate-500) |
| Thumb (on thumb hover)     | `rgb(148 163 184 / 0.5)`        | `rgb(100 116 139 / 0.5)`        |
| Track                      | transparent                      | transparent                      |

### Implementation

- Firefox: `scrollbar-width: thin` + `scrollbar-color` property
- WebKit/Blink: `::-webkit-scrollbar`, `::-webkit-scrollbar-track`, `::-webkit-scrollbar-thumb` pseudo-elements
- Applied globally (`*`) so all scrollable containers (sidebar, content, code blocks) share the same style

## Accessibility

- Semantic HTML5 elements (`<nav>`, `<main>`, `<article>`, `<aside>`, `<header>`, `<footer>`)
- ARIA labels on interactive elements (search, sidebar toggle, theme toggle)
- Skip-to-content link (visually hidden, visible on focus)
- Focus management: visible focus rings on all interactive elements
- Color contrast: WCAG 2.1 AA minimum (4.5:1 for text, 3:1 for large text)
- Reduced motion: respect `prefers-reduced-motion` (disable smooth scroll, progress bar animations)
- Screen reader: all images have alt text, icons have aria-labels
- Sidebar icons have `aria-hidden="true"` (decorative, text labels provide meaning)

## Performance Targets

| Metric                     | Target     |
|----------------------------|------------|
| Lighthouse Performance     | > 90       |
| First Contentful Paint     | < 1s       |
| Largest Contentful Paint   | < 2s       |
| Total JS payload (gzipped) | < 30 KB    |
| Total CSS payload (gzipped)| < 15 KB    |
| Time to Interactive        | < 2s       |

## Asset Strategy

- Tailwind CSS: compiled at build time via standalone CLI, single output file
- HTMX v2.0.4: served from static assets (~14 KB gzipped)
- Alpine.js v3: served from static assets (~8 KB gzipped)
- Mermaid.js: loaded only on pages with diagrams (lazy load)
- Icon library: Lucide or Heroicons (SVG, tree-shakeable)
- All assets embedded in binary via `go:embed`
- Cache-busting via content hash in filename (e.g., `styles.a1b2c3.css`)
