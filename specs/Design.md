# Design System

## Status

| Field   | Value            |
|---------|------------------|
| Version | 1.0              |
| Status  | Implemented      |
| Author  | Hernando Clareth |
| Created | 2026-08-12       |
| Source  | https://blogo.hclareth.space/ |

---

## Overview

This document captures the exact visual design of the BLOGO production site. It is the single source of truth for all UI implementation. Every component is documented with its precise HTML structure, Tailwind CSS classes (both light and dark variants), spacing, and interactive behavior.

**Theme strategy**: Tailwind `darkMode: "class"` on `<html>`. Theme stored in `localStorage`, respects `prefers-color-scheme` on first visit. Inline `<script>` in `<head>` prevents flash of wrong theme.

**Tooling**: Tailwind CSS (standalone CLI), `@tailwindcss/typography` plugin, Alpine.js (interactivity), HTMX (partial navigation).

---

## 1. Design Tokens

### 1.1 Color Palette

#### Backgrounds

| Surface         | Light                | Dark                          |
|-----------------|----------------------|-------------------------------|
| Page            | `bg-white`           | `bg-[#0a0f1a]`               |
| Header          | `bg-white`           | `bg-[#0d1220]`               |
| Sidebar         | `bg-slate-50`        | `bg-[#0d1220]`               |
| Footer          | `bg-slate-50`        | `bg-[#0d1220]`               |
| Dropdown/panel  | `bg-white`           | `bg-[#111827]`               |
| Input field     | `bg-slate-50`        | `bg-white/5`                 |
| Hover (generic) | `hover:bg-slate-100` | `dark:hover:bg-white/5`      |
| Hover (header)  | `hover:bg-slate-100` | `dark:hover:bg-slate-800`    |
| Active nav      | `bg-blue-500/10`     | `bg-blue-500/10`             |
| Active repo     | `bg-blue-50`         | `bg-gradient-to-r from-blue-900/50 to-blue-800/20` |

#### Borders

| Element        | Light                  | Dark                     |
|----------------|------------------------|--------------------------|
| Layout divider | `border-slate-200`     | `border-[#141c2e]`      |
| Input          | `border-slate-200`     | `border-white/10`        |
| Search input   | `border-slate-300`     | `border-slate-600`       |
| Selector glow  | `border-blue-400/50`   | `border-blue-500/60`     |
| Prev/next card | `border-slate-200`     | `border-[#141c2e]`      |
| Prev/next hover| `hover:border-blue-300`| `dark:hover:border-blue-600` |
| Disabled card  | `border-slate-100`     | `border-slate-800`       |

#### Text

| Role            | Light                | Dark                    |
|-----------------|----------------------|-------------------------|
| Body            | `text-slate-800`     | `text-slate-200`        |
| Heading         | (inherits body)      | (inherits body)         |
| Muted           | `text-slate-500`     | `text-slate-400`        |
| Extra muted     | `text-slate-400`     | `text-slate-500`        |
| Link            | `text-blue-600`      | `text-blue-400`         |
| Active nav text | `text-blue-600`      | `text-blue-400`         |
| Active repo text| `text-blue-700`      | `text-white`            |
| Separator       | `text-slate-300`     | `text-slate-600`        |
| Disabled        | `text-slate-400`     | `text-slate-600`        |
| Monospace label | `text-slate-400`     | `text-slate-500`        |

#### Accent

| Use              | Value                |
|------------------|----------------------|
| Primary accent   | `blue-500`           |
| Avatar bg        | `bg-blue-600/20`     |
| Avatar text      | `text-blue-500`      |
| Checkmark bg     | `bg-blue-500`        |
| Progress bar     | `bg-blue-500`        |
| Heart            | `text-red-500`       |

### 1.2 Shadows

| Element           | Light                                            | Dark                                              |
|-------------------|--------------------------------------------------|---------------------------------------------------|
| Repo selector btn | `shadow-[0_0_15px_rgba(59,130,246,0.12)]`        | `shadow-[0_0_20px_rgba(37,99,235,0.18)]`         |
| Dropdown panel    | `shadow-xl`                                      | `shadow-xl`                                       |

### 1.3 Border Radius

| Element             | Value          |
|---------------------|----------------|
| Avatar (image/init) | `rounded-full` |
| Selector button     | `rounded-2xl`  |
| Dropdown panel      | `rounded-2xl`  |
| Dropdown items      | `rounded-2xl`  |
| Search input        | `rounded-lg`   |
| Filter input        | `rounded-xl`   |
| Chevron container   | `rounded-xl` (w-7 h-7) |
| Nav link            | `rounded-lg`   |
| Prev/next card      | `rounded-lg`   |
| Code inline         | `rounded`      |
| Code block          | `rounded-lg`   |
| Buttons             | `rounded-lg`   |
| Checkmark circle    | `rounded-full` |
| Nav dot             | `rounded-full` |
| Callout block       | `rounded-r-lg` |
| Progress bar        | (none)         |
| Kbd badge           | `rounded`      |

### 1.4 Spacing Scale

| Context                | Value           |
|------------------------|-----------------|
| Header height          | `h-14` (56px)   |
| Footer height          | `h-[60px]`      |
| Sidebar width          | `w-72` (288px)  |
| Content max-width      | `max-w-5xl`     |
| Content padding        | `px-8 py-8`     |
| Content margin (md)    | `md:ml-[14rem] md:mr-[14rem]` |
| Content bottom pad     | `pb-[60px]`     |
| Sidebar padding        | `p-3`           |
| Sidebar header padding | `p-4`           |
| Selector btn padding   | `px-3 py-2.5`   |
| Dropdown padding       | `p-3`           |
| Dropdown item padding  | `px-3 py-2.5`   |
| Nav link padding       | `px-3 py-1.5`   |
| Footer inner padding   | `px-6`          |
| Footer grid gap        | `gap-4`         |
| Prev/next padding      | `p-4`           |
| Prev/next top margin   | `mt-12 pt-8`    |

### 1.5 Typography

| Element          | Size / Weight                                  |
|------------------|-------------------------------------------------|
| Page title (H1)  | `text-3xl font-bold`                            |
| H2               | `text-xl font-bold` (via prose)                 |
| H3               | `text-lg font-semibold` (via prose)             |
| Body text        | System font, 16px base via `prose`              |
| Sidebar label    | `text-[10px] font-semibold uppercase tracking-[0.2em]` |
| Repo name (btn)  | `text-sm font-semibold`                         |
| Repo author (btn)| `text-[11px]`                                   |
| Repo name (list) | `font-semibold` (inherits text-sm)              |
| Repo author (list)| `text-[11px]`                                  |
| Nav link         | `text-sm` (space-y-1 between items)             |
| Active nav link  | `text-sm font-medium`                           |
| Breadcrumb       | `text-sm`                                       |
| Metadata line    | `text-sm`                                       |
| Footer text      | `text-xs`                                       |
| Footer brand     | `text-xs font-bold`                             |
| Footer mono      | `text-xs font-mono`                             |
| Kbd badge        | `text-xs`                                       |
| Prev/next label  | `text-xs`                                       |
| Prev/next title  | `text-sm font-medium`                           |

---

## 2. Components

### 2.1 Progress Bar

Position: fixed, top-0, z-50. Full viewport width, 3px height. Blue fill tracks scroll percentage.

```html
<div class="fixed top-0 left-0 w-full h-[3px] z-50"
     x-data="{ progress: 0 }"
     @scroll.window="progress = Math.min(100, ...)">
    <div class="h-full bg-blue-500 transition-[width] duration-150"
         :style="'width:' + progress + '%'"></div>
</div>
```

### 2.2 Header

Sticky, z-40. Height: 56px. Bottom border separating from content.

```html
<header class="sticky top-0 z-40 bg-white dark:bg-[#0d1220]
               border-b border-slate-200 dark:border-[#141c2e]">
    <div class="flex items-center justify-between px-4 h-14">
        <!-- Left: hamburger (mobile) + logo -->
        <!-- Center: search input (hidden sm:flex) -->
        <!-- Right: theme toggle -->
    </div>
</header>
```

#### 2.2.1 Logo

```html
<a href="/" class="flex items-center gap-2 font-bold text-lg">
    <svg class="w-6 h-6 text-blue-500" ...><!-- globe/planet icon --></svg>
    <span>BLOGO</span>
</a>
```

#### 2.2.2 Mobile Hamburger

```html
<button @click="sidebarOpen = !sidebarOpen"
        class="md:hidden p-2 rounded-lg hover:bg-slate-100 dark:hover:bg-slate-800"
        aria-label="Toggle sidebar">
    <svg class="w-5 h-5" ...><!-- 3-line menu icon --></svg>
</button>
```

#### 2.2.3 Search Input

Centered, hidden on mobile. Disabled (Phase 2 feature).

```html
<div class="hidden sm:flex flex-1 max-w-md mx-4">
    <div class="relative w-full">
        <svg class="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-slate-400" .../>
        <input type="search" placeholder="Search anything..." disabled
               class="w-full pl-10 pr-16 py-1.5 rounded-lg
                      border border-slate-300 dark:border-slate-600
                      bg-slate-50 dark:bg-slate-800 text-sm
                      focus:outline-none focus:ring-2 focus:ring-blue-500
                      disabled:opacity-60 disabled:cursor-not-allowed">
        <kbd class="absolute right-3 top-1/2 -translate-y-1/2 text-xs text-slate-400
                    border border-slate-300 dark:border-slate-600 rounded px-1.5 py-0.5">
            ⌘K
        </kbd>
    </div>
</div>
```

#### 2.2.4 Theme Toggle

```html
<button x-data @click="toggle dark/light classes + localStorage"
        class="p-2 rounded-lg hover:bg-slate-100 dark:hover:bg-slate-800"
        aria-label="Toggle theme">
    <!-- Sun icon: hidden dark:block -->
    <svg class="w-5 h-5 hidden dark:block" .../>
    <!-- Moon icon: block dark:hidden -->
    <svg class="w-5 h-5 block dark:hidden" .../>
</button>
```

### 2.3 Sidebar

Fixed on mobile, sticky on desktop. Positioned below header, above footer.

```html
<!-- Mobile backdrop -->
<div x-show="sidebarOpen" @click="sidebarOpen = false"
     class="fixed inset-0 bg-black/50 z-30 md:hidden" x-transition.opacity></div>

<!-- Sidebar -->
<aside :class="sidebarOpen ? 'translate-x-0' : '-translate-x-full md:translate-x-0'"
       class="fixed md:sticky top-14 left-0 z-30 w-72
              h-[calc(100vh-3.5rem-60px)] overflow-y-auto
              bg-slate-50 dark:bg-[#0d1220]
              border-r border-slate-200 dark:border-[#141c2e]
              transition-transform duration-200 flex flex-col">
    <!-- sidebar-header -->
    <!-- sidebar-nav -->
</aside>
```

#### 2.3.1 Repo Selector -- Closed State

Container: `p-4`, bottom border. Alpine state: `{ open: false, search: '' }`.

Label:
```html
<p class="text-[10px] font-semibold uppercase tracking-[0.2em]
          text-slate-400 dark:text-slate-500 mb-3 px-1">
    Current source
</p>
```

Button:
```html
<button @click="open = !open"
        class="w-full flex items-center justify-between rounded-2xl
               border border-blue-400/50 dark:border-blue-500/60
               bg-white dark:bg-[#0B1325]/90
               px-3 py-2.5
               shadow-[0_0_15px_rgba(59,130,246,0.12)]
               dark:shadow-[0_0_20px_rgba(37,99,235,0.18)]
               backdrop-blur transition-colors">

    <div class="flex items-center gap-2.5 min-w-0 flex-1">
        <!-- Avatar: circular, 30x30 -->
        <img src="{{.AvatarURL}}" alt=""
             class="w-[30px] h-[30px] rounded-full flex-shrink-0 object-cover"
             onerror="this.style.display='none';this.nextElementSibling.style.display='flex'">
        <!-- Fallback initial (hidden by default) -->
        <span class="w-[30px] h-[30px] rounded-full bg-blue-600/20
                     items-center justify-center text-blue-500
                     text-sm font-bold flex-shrink-0 hidden">
            {{initial .DocTitle}}
        </span>

        <div class="text-left min-w-0 flex-1">
            <span class="block text-sm font-semibold text-slate-800 dark:text-white truncate">
                {{.DocTitle}}
            </span>
            <span class="block text-[11px] text-slate-500 dark:text-slate-400 truncate">
                by {{.Author}}
            </span>
        </div>
    </div>

    <!-- Chevron button -->
    <div class="flex items-center justify-center w-7 h-7 rounded-xl
                border border-slate-200 dark:border-white/10
                bg-slate-50 dark:bg-white/5 ml-1.5 flex-shrink-0">
        <svg class="w-3.5 h-3.5 text-slate-500 dark:text-white transition-transform"
             :class="open ? 'rotate-180' : ''" ...>
            <path d="M6 9l6 6 6-6"/>
        </svg>
    </div>
</button>
```

#### 2.3.2 Repo Selector -- Dropdown Panel

```html
<div x-show="open" x-transition
     class="absolute left-4 right-4 mt-3
            rounded-2xl border border-slate-200 dark:border-white/10
            bg-white dark:bg-[#111827]
            p-3 backdrop-blur shadow-xl z-50">

    <!-- Header -->
    <p class="text-[10px] font-semibold uppercase tracking-[0.2em]
              text-slate-400 dark:text-slate-500 mb-4 px-1">
        Switch source repository
    </p>

    <!-- Search filter -->
    <div class="mb-3 px-1">
        <input type="text" x-model="search" placeholder="Search repos..."
               class="w-full px-3 py-1.5 text-sm rounded-xl
                      border border-slate-200 dark:border-white/10
                      bg-slate-50 dark:bg-white/5
                      text-slate-800 dark:text-slate-200
                      placeholder-slate-400
                      focus:outline-none focus:ring-1 focus:ring-blue-500"
               @click.stop>
    </div>

    <!-- Repo list -->
    <div class="max-h-64 overflow-y-auto space-y-1.5">
        <!-- items (see 2.3.3 and 2.3.4) -->
    </div>
</div>
```

#### 2.3.3 Repo Item -- Active

```html
<a href="/{{.Slug}}"
   hx-get="/{{.Slug}}" hx-target="#main-content" hx-push-url="true"
   hx-swap="innerHTML scroll:window:top"
   class="flex items-center justify-between rounded-2xl px-3 py-2.5
          text-sm transition-colors
          bg-blue-50 dark:bg-gradient-to-r dark:from-blue-900/50 dark:to-blue-800/20
          hover:bg-blue-100 dark:hover:bg-blue-900/60">

    <div class="flex items-center gap-2.5 min-w-0 flex-1">
        <!-- Avatar: circular, blue bg -->
        <img src="{{.AvatarURL}}" alt=""
             class="w-[30px] h-[30px] rounded-full flex-shrink-0 object-cover" ...>
        <span class="w-[30px] h-[30px] rounded-full bg-blue-600/20
                     items-center justify-center text-sm font-bold
                     flex-shrink-0 hidden text-blue-500">
            {{initial .Name}}
        </span>

        <div class="min-w-0 flex-1">
            <span class="block truncate font-semibold text-blue-700 dark:text-white">
                {{.Name}}
            </span>
            <span class="block text-[11px] text-slate-500 dark:text-slate-400 truncate">
                by {{.Author}}
            </span>
        </div>
    </div>

    <!-- Blue checkmark -->
    <span class="flex items-center justify-center w-6 h-6
                 rounded-full bg-blue-500 flex-shrink-0 ml-2">
        <svg class="w-3.5 h-3.5 text-white" fill="none"
             stroke="currentColor" stroke-width="3" viewBox="0 0 24 24">
            <path d="M5 13l4 4L19 7"/>
        </svg>
    </span>
</a>
```

#### 2.3.4 Repo Item -- Inactive

```html
<a href="/{{.Slug}}" ...
   class="flex items-center justify-between rounded-2xl px-3 py-2.5
          text-sm transition-colors
          hover:bg-slate-100 dark:hover:bg-white/5">

    <div class="flex items-center gap-2.5 min-w-0 flex-1">
        <!-- Avatar: circular, grey bg -->
        <img src="{{.AvatarURL}}" alt=""
             class="w-[30px] h-[30px] rounded-full flex-shrink-0 object-cover" ...>
        <span class="w-[30px] h-[30px] rounded-full bg-slate-200 dark:bg-slate-700
                     items-center justify-center text-sm font-bold
                     flex-shrink-0 hidden text-slate-500 dark:text-slate-300">
            {{initial .Name}}
        </span>

        <div class="min-w-0 flex-1">
            <span class="block truncate font-semibold text-slate-800 dark:text-white">
                {{.Name}}
            </span>
            <span class="block text-[11px] text-slate-500 dark:text-slate-400 truncate">
                by {{.Author}}
            </span>
        </div>
    </div>
    <!-- No checkmark -->
</a>
```

#### 2.3.5 Single Repo Fallback

When only one repo is configured, no dropdown. Static header only.

```html
<div class="p-4 border-b border-slate-200 dark:border-[#141c2e]">
    <h2 class="font-semibold text-sm">{{.DocTitle}}</h2>
    <p class="text-xs text-slate-500 dark:text-slate-400">by {{.Author}}</p>
</div>
```

#### 2.3.6 Navigation Links

Flat list, dot bullet indicators. No grouping, no section icons.

```html
<nav class="flex-1 overflow-y-auto p-3 space-y-1"
     aria-label="Documentation navigation">

    <!-- Active link -->
    <a href="{{.URL}}" hx-get="..." title="{{.Title}}"
       class="flex items-center gap-2 px-3 py-1.5 text-sm rounded-lg
              transition-colors
              bg-blue-500/10 text-blue-600 dark:text-blue-400 font-medium">
        <span class="w-1.5 h-1.5 rounded-full bg-blue-500 flex-shrink-0"></span>
        <span class="truncate">{{.Title}}</span>
    </a>

    <!-- Inactive link -->
    <a href="{{.URL}}" hx-get="..." title="{{.Title}}"
       class="flex items-center gap-2 px-3 py-1.5 text-sm rounded-lg
              transition-colors
              text-slate-600 dark:text-slate-300
              hover:bg-slate-200/50 dark:hover:bg-slate-700/50">
        <span class="w-1.5 h-1.5 rounded-full bg-slate-400 dark:bg-slate-500
                     flex-shrink-0"></span>
        <span class="truncate">{{.Title}}</span>
    </a>
</nav>
```

### 2.4 Content Area

```html
<main class="flex-1 min-w-0 flex flex-col md:ml-[14rem] md:mr-[14rem] pb-[60px]">
    <div id="main-content" class="flex-1">
        <article class="max-w-5xl px-8 py-8 mr-8">
            <!-- breadcrumbs -->
            <!-- h1 -->
            <!-- metadata -->
            <!-- prose content -->
            <!-- prev/next -->
        </article>
    </div>
</main>
```

#### 2.4.1 Breadcrumbs

```html
<nav class="flex items-center gap-1.5 text-sm
            text-slate-500 dark:text-slate-400 mb-4"
     aria-label="Breadcrumb">
    <a href="..." class="hover:text-slate-700 dark:hover:text-slate-200
                         transition-colors">Home</a>
    <span class="text-slate-300 dark:text-slate-600">&gt;</span>
    <span class="text-slate-700 dark:text-slate-300">Current Page</span>
</nav>
```

#### 2.4.2 Page Title

```html
<h1 class="text-3xl font-bold mb-4">{{.Section.Title}}</h1>
```

#### 2.4.3 Metadata Line

```html
<div class="flex flex-wrap items-center gap-3 text-sm
            text-slate-500 dark:text-slate-400 mb-6">
    <span class="flex items-center gap-1">
        <svg class="w-4 h-4" .../><!-- clock icon -->
        {{.ReadingTime}} min read
    </span>
    <span class="text-slate-300 dark:text-slate-600">&bull;</span>
    <span class="flex items-center gap-1">
        <svg class="w-4 h-4" .../><!-- person icon -->
        {{.Author}}
    </span>
    <span class="text-slate-300 dark:text-slate-600">&bull;</span>
    <a href="{{.OriginalURL}}" target="_blank" rel="noopener noreferrer"
       class="flex items-center gap-1
              text-blue-600 dark:text-blue-400 hover:underline">
        Original README
        <svg class="w-3 h-3" .../><!-- external link icon -->
    </a>
</div>
```

#### 2.4.4 Prose Content

```html
<div class="prose prose-slate dark:prose-invert max-w-none
            prose-headings:scroll-mt-20
            prose-h2:text-xl prose-h2:font-bold prose-h2:mt-10 prose-h2:mb-4
            prose-h3:text-lg prose-h3:font-semibold prose-h3:mt-8 prose-h3:mb-3
            prose-a:text-blue-600 dark:prose-a:text-blue-400
            prose-code:text-sm prose-code:bg-slate-100
            dark:prose-code:bg-slate-800
            prose-code:px-1.5 prose-code:py-0.5 prose-code:rounded
            prose-pre:bg-slate-900 dark:prose-pre:bg-slate-950
            prose-pre:rounded-lg
            prose-img:rounded-lg
            prose-table:text-sm
            prose-th:bg-slate-100 dark:prose-th:bg-slate-800">
    <!-- rendered markdown -->
</div>
```

### 2.5 Previous / Next Navigation

```html
<div class="flex flex-col sm:flex-row gap-4
            mt-12 pt-8 border-t border-slate-200 dark:border-[#141c2e]">

    <!-- Previous (enabled) -->
    <a href="..." hx-get="..."
       class="flex-1 group flex items-center gap-3 p-4 rounded-lg
              border border-slate-200 dark:border-[#141c2e]
              hover:border-blue-300 dark:hover:border-blue-600
              transition-colors">
        <svg class="w-5 h-5 text-slate-400
                    group-hover:text-blue-500 transition-colors" .../>
        <div>
            <div class="text-xs text-slate-500 dark:text-slate-400">Previous</div>
            <div class="text-sm font-medium
                        group-hover:text-blue-600
                        dark:group-hover:text-blue-400
                        transition-colors">{{title}}</div>
        </div>
    </a>

    <!-- Previous (disabled) -->
    <div class="flex-1 flex items-center gap-3 p-4 rounded-lg
                border border-slate-100 dark:border-slate-800
                text-slate-400 dark:text-slate-600">
        <svg class="w-5 h-5" .../>
        <div>
            <div class="text-xs">Previous</div>
            <div class="text-sm">No previous</div>
        </div>
    </div>

    <!-- Next (enabled) -->
    <a href="..." hx-get="..."
       class="flex-1 group flex items-center justify-end gap-3 p-4
              rounded-lg border border-slate-200 dark:border-[#141c2e]
              hover:border-blue-300 dark:hover:border-blue-600
              transition-colors text-right">
        <div>
            <div class="text-xs text-slate-500 dark:text-slate-400">Next</div>
            <div class="text-sm font-medium
                        text-blue-600 dark:text-blue-400
                        group-hover:text-blue-700
                        dark:group-hover:text-blue-300
                        transition-colors">{{title}}</div>
        </div>
        <svg class="w-5 h-5 text-blue-500
                    group-hover:text-blue-600 transition-colors" .../>
    </a>

    <!-- Next (disabled) -->
    <div class="flex-1 flex items-center justify-end gap-3 p-4
                rounded-lg border border-slate-100 dark:border-slate-800
                text-slate-400 dark:text-slate-600 text-right">
        <div>
            <div class="text-xs">Next</div>
            <div class="text-sm">No next</div>
        </div>
        <svg class="w-5 h-5" .../>
    </div>
</div>
```

### 2.6 Footer

Fixed at bottom, z-20. Three-column grid on desktop, single column on mobile.

```html
<footer class="fixed bottom-0 left-0 right-0 z-20
               h-[60px] border-t border-slate-200 dark:border-[#141c2e]
               bg-slate-50 dark:bg-[#0d1220]">
    <div class="h-full max-w-7xl mx-auto px-6 flex items-center">
        <div class="grid grid-cols-1 md:grid-cols-3 gap-4 w-full items-center">

            <!-- Left: identity (hidden on mobile) -->
            <div class="hidden md:flex items-center gap-2">
                <svg class="w-6 h-6 text-blue-500 flex-shrink-0" .../>
                <div class="flex items-center gap-2 text-xs">
                    <span class="font-bold">BLOGO</span>
                    <span class="text-slate-400 dark:text-slate-500 font-mono">
                        BLG-001
                    </span>
                    <a href="https://hclareth.space" target="_blank"
                       class="text-blue-600 dark:text-blue-400 hover:underline
                              flex items-center gap-0.5">
                        hclareth.space
                        <svg class="w-2.5 h-2.5" .../>
                    </a>
                </div>
            </div>

            <!-- Center: tagline -->
            <div class="text-xs text-slate-500 dark:text-slate-400 text-center">
                Made with <span class="text-red-500">&hearts;</span> and
                <a href="https://go.dev" target="_blank"
                   class="text-blue-600 dark:text-blue-400 hover:underline">Go</a>
            </div>

            <!-- Right: license + GitHub (hidden on mobile) -->
            <div class="hidden md:flex items-center justify-end gap-3
                        text-xs text-slate-500 dark:text-slate-400">
                <span>Apache 2.0</span>
                <a href="https://github.com/hclareth7/blogo" target="_blank"
                   class="text-slate-400 hover:text-slate-600
                          dark:hover:text-slate-300 transition-colors"
                   aria-label="GitHub repository">
                    <svg class="w-4 h-4" fill="currentColor" .../>
                </a>
            </div>

        </div>
    </div>
</footer>
```

---

## 3. Callout Components

Defined in `input.css` as `@layer components`.

```css
.callout        { @apply border-l-4 rounded-r-lg p-4 my-4; }
.callout-note   { @apply border-blue-500 bg-blue-50 dark:bg-blue-900/20; }
.callout-warning{ @apply border-yellow-500 bg-yellow-50 dark:bg-yellow-900/20; }
.callout-tip    { @apply border-green-500 bg-green-50 dark:bg-green-900/20; }
.callout-title  { @apply font-semibold mb-1; }
```

Title colors per type:

| Type    | Light               | Dark                |
|---------|----------------------|---------------------|
| Note    | `text-blue-700`      | `text-blue-300`     |
| Warning | `text-yellow-700`    | `text-yellow-300`   |
| Tip     | `text-green-700`     | `text-green-300`    |

---

## 4. Custom Scrollbars

Invisible by default. Appear on container hover. Applied globally.

```css
* { scrollbar-width: thin; scrollbar-color: transparent transparent; }
*:hover { scrollbar-color: rgb(148 163 184 / 0.3) transparent; }
.dark *:hover { scrollbar-color: rgb(100 116 139 / 0.3) transparent; }
```

Thumb width: 6px. Shape: `border-radius: 9999px`. Track: transparent.

---

## 5. Syntax Highlighting (Chroma -- Dracula)

Code blocks use Dracula theme, consistent across both light and dark modes.

| Token       | Color     |
|-------------|-----------|
| Background  | `#282a36` |
| Foreground  | `#f8f8f2` |
| Keywords    | `#ff79c6` |
| Types       | `#8be9fd` |
| Functions   | `#50fa7b` |
| Strings     | `#f1fa8c` |
| Numbers     | `#bd93f9` |
| Comments    | `#6272a4` |
| Operators   | `#ff79c6` |

---

## 6. Z-Index Stack

| Layer          | z-index |
|----------------|---------|
| Progress bar   | `z-50`  |
| Header         | `z-40`  |
| Sidebar        | `z-30`  |
| Mobile backdrop| `z-30`  |
| Dropdown panel | `z-50`  |
| Footer         | `z-20`  |

---

## 7. Tailwind Configuration

```js
// tailwind.config.js
module.exports = {
  content: ["./web/templates/**/*.html"],
  darkMode: "class",
  theme: { extend: {} },
  plugins: [require("@tailwindcss/typography")],
}
```

---

## 8. Interactive Behaviors

| Behavior                 | Mechanism                                         |
|--------------------------|---------------------------------------------------|
| Theme toggle             | Alpine.js inline, `classList.toggle('dark'/'light')`, `localStorage` |
| Mobile sidebar           | Alpine.js `sidebarOpen`, slide transform           |
| Repo selector open/close | Alpine.js `open`, `@click.outside`, `@keydown.escape` |
| Repo search filter       | Alpine.js `x-model` + `x-show` string match       |
| Page navigation          | HTMX `hx-get`, `hx-target="#main-content"`, OOB sidebar swap |
| Reading progress         | Alpine.js `@scroll.window`, width percentage       |
| Chevron rotation         | Alpine.js `:class="open ? 'rotate-180' : ''"`     |

---

## 9. HTMX Navigation Pattern

All sidebar links, prev/next, and repo selector items carry:

```html
hx-get="{{.URL}}"
hx-target="#main-content"
hx-push-url="true"
hx-swap="innerHTML scroll:window:top"
```

Server returns `content-fragment` template (not full page) when `HX-Request: true` header is present. Fragment includes `<title>` tag and `sidebar-oob` template with `hx-swap-oob="innerHTML"` for `#sidebar-header` and `#sidebar-nav`.

All links retain `href` for no-JS fallback.
