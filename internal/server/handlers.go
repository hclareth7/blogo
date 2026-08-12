package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/hclareth7/blogo/internal/parser"
	"github.com/hclareth7/blogo/internal/renderer"
)

func (s *Server) handleHome(w http.ResponseWriter, r *http.Request) {
	rs := s.defaultRepo()
	if rs == nil {
		http.Error(w, "no content available", http.StatusInternalServerError)
		return
	}

	sections := rs.Index.Sections()
	if len(sections) == 0 {
		http.Error(w, "no content available", http.StatusInternalServerError)
		return
	}

	first := sections[0]
	data := s.buildPageData(rs, first, "", r.URL.Path)
	s.render(w, r, "page.html", data)
}

// handleSlug1 resolves /{slug1} — could be a repo slug or a section in the default repo.
func (s *Server) handleSlug1(w http.ResponseWriter, r *http.Request) {
	slug1 := chi.URLParam(r, "slug1")

	// Check if it's a repo slug
	if rs, ok := s.repos[slug1]; ok {
		sections := rs.Index.Sections()
		if len(sections) == 0 {
			s.render404(w, r)
			return
		}
		first := sections[0]
		data := s.buildPageData(rs, first, "", r.URL.Path)
		s.render(w, r, "page.html", data)
		return
	}

	// Fallback: section in default repo
	rs := s.defaultRepo()
	if rs == nil {
		s.render404(w, r)
		return
	}
	section, ok := rs.Index.Lookup(slug1)
	if !ok {
		s.render404(w, r)
		return
	}
	data := s.buildPageData(rs, section, "", r.URL.Path)
	s.render(w, r, "page.html", data)
}

// handleSlug2 resolves /{slug1}/{slug2} — repo/section or section/subsection in default repo.
func (s *Server) handleSlug2(w http.ResponseWriter, r *http.Request) {
	slug1 := chi.URLParam(r, "slug1")
	slug2 := chi.URLParam(r, "slug2")

	// Check if slug1 is a repo
	if rs, ok := s.repos[slug1]; ok {
		section, ok := rs.Index.Lookup(slug2)
		if !ok {
			s.render404(w, r)
			return
		}
		data := s.buildPageData(rs, section, "", r.URL.Path)
		s.render(w, r, "page.html", data)
		return
	}

	// Fallback: subsection in default repo
	rs := s.defaultRepo()
	if rs == nil {
		s.render404(w, r)
		return
	}
	route := slug1 + "/" + slug2
	section, ok := rs.Index.Lookup(route)
	if !ok {
		s.render404(w, r)
		return
	}
	data := s.buildPageData(rs, section, slug1, r.URL.Path)
	s.render(w, r, "page.html", data)
}

// handleSlug3 resolves /{repo}/{section}/{subsection}.
func (s *Server) handleSlug3(w http.ResponseWriter, r *http.Request) {
	slug1 := chi.URLParam(r, "slug1")
	slug2 := chi.URLParam(r, "slug2")
	slug3 := chi.URLParam(r, "slug3")

	rs, ok := s.repos[slug1]
	if !ok {
		s.render404(w, r)
		return
	}

	route := slug2 + "/" + slug3
	section, ok := rs.Index.Lookup(route)
	if !ok {
		s.render404(w, r)
		return
	}

	data := s.buildPageData(rs, section, slug2, r.URL.Path)
	s.render(w, r, "page.html", data)
}

func (s *Server) handleHealthz(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}

func (s *Server) handleReadyz(w http.ResponseWriter, r *http.Request) {
	ready := false
	for _, rs := range s.repos {
		if len(rs.Index.Sections()) > 0 {
			ready = true
			break
		}
	}
	w.Header().Set("Content-Type", "text/plain")
	if !ready {
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Write([]byte("not ready"))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ready"))
}

func (s *Server) buildPageData(rs *RepoState, section *parser.Section, parentID, currentPath string) *renderer.PageData {
	sections := rs.Index.Sections()

	originalURL := buildOriginalURL(rs, section)

	rawText := section.RawText
	if rawText == "" && len(section.Children) > 0 {
		for _, c := range section.Children {
			rawText += c.RawText + " "
		}
	}

	metaDesc := rawText
	if len(metaDesc) > 160 {
		metaDesc = metaDesc[:160]
	}

	return &renderer.PageData{
		Title:       section.Title,
		Section:     section,
		NavTree:     rs.NavTree,
		Breadcrumbs: rs.NavBld.BuildBreadcrumbs(section.ID, parentID, sections),
		PrevNext:    rs.NavBld.BuildPrevNext(sections, section.ID),
		ReadingTime: readingTime(section.WordCount),
		Author:      rs.Meta.Author,
		OriginalURL: originalURL,
		MetaDesc:    metaDesc,
		CurrentPath: currentPath,
		DocTitle:    rs.Doc.Title,
		RepoSlug:    rs.Meta.Slug,
		RepoList:    s.repoList(rs.Meta.Slug),
		IsMultiRepo: s.cfg.IsMultiRepo(),
	}
}

func buildOriginalURL(rs *RepoState, section *parser.Section) string {
	for _, repo := range []struct{ slug, base string }{
		{"system-design", "https://github.com/karanpratapsingh/system-design#"},
	} {
		if rs.Meta.Slug == repo.slug {
			return repo.base + section.ID
		}
	}
	return ""
}

func (s *Server) render(w http.ResponseWriter, r *http.Request, name string, data *renderer.PageData) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	if r.Header.Get("HX-Request") == "true" && name == "page.html" {
		name = "content-fragment"
	}

	if err := s.renderer.RenderPage(w, name, data); err != nil {
		s.logger.Error("template render error", "template", name, "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}

func (s *Server) render404(w http.ResponseWriter, r *http.Request) {
	rs := s.defaultRepo()
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusNotFound)

	data := &renderer.PageData{
		Title:       "Page Not Found",
		CurrentPath: r.URL.Path,
		IsMultiRepo: s.cfg.IsMultiRepo(),
	}

	if rs != nil {
		data.NavTree = rs.NavTree
		data.DocTitle = rs.Doc.Title
		data.RepoSlug = rs.Meta.Slug
		data.RepoList = s.repoList(rs.Meta.Slug)
	}

	if err := s.renderer.RenderPage(w, "404.html", data); err != nil {
		s.logger.Error("404 template render error", "error", err)
		http.Error(w, "page not found", http.StatusNotFound)
	}
}

func readingTime(wordCount int) int {
	m := wordCount / 200
	if m < 1 {
		return 1
	}
	return m
}
