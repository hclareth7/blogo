package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/hclareth7/blogo/internal/parser"
	"github.com/hclareth7/blogo/internal/renderer"
)

func (s *Server) handleHome(w http.ResponseWriter, r *http.Request) {
	sections := s.index.Sections()
	if len(sections) == 0 {
		http.Error(w, "no content available", http.StatusInternalServerError)
		return
	}

	first := sections[0]
	data := s.buildPageData(first, "", r.URL.Path)
	s.render(w, r, "page.html", data)
}

func (s *Server) handleSection(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "section")
	section, ok := s.index.Lookup(slug)
	if !ok {
		s.render404(w, r)
		return
	}

	data := s.buildPageData(section, "", r.URL.Path)
	s.render(w, r, "page.html", data)
}

func (s *Server) handleSubsection(w http.ResponseWriter, r *http.Request) {
	parentSlug := chi.URLParam(r, "section")
	childSlug := chi.URLParam(r, "subsection")
	route := parentSlug + "/" + childSlug

	section, ok := s.index.Lookup(route)
	if !ok {
		s.render404(w, r)
		return
	}

	data := s.buildPageData(section, parentSlug, r.URL.Path)
	s.render(w, r, "page.html", data)
}

func (s *Server) handleHealthz(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}

func (s *Server) handleReadyz(w http.ResponseWriter, r *http.Request) {
	if len(s.index.Sections()) == 0 {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Write([]byte("not ready"))
		return
	}
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ready"))
}

func (s *Server) buildPageData(section *parser.Section, parentID, currentPath string) *renderer.PageData {
	sections := s.index.Sections()

	originalBase := "https://github.com/karanpratapsingh/system-design#"
	originalURL := originalBase + section.ID

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
		NavTree:     s.navTree,
		Breadcrumbs: s.navBld.BuildBreadcrumbs(section.ID, parentID, sections),
		PrevNext:    s.navBld.BuildPrevNext(sections, section.ID),
		ReadingTime: readingTime(section.WordCount),
		Author:      "Karan Pratap Singh",
		OriginalURL: originalURL,
		MetaDesc:    metaDesc,
		CurrentPath: currentPath,
		DocTitle:    s.doc.Title,
	}
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
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusNotFound)
	data := &renderer.PageData{
		Title:       "Page Not Found",
		NavTree:     s.navTree,
		CurrentPath: r.URL.Path,
		DocTitle:    s.doc.Title,
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
