package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
)

func (s *Server) Routes() chi.Router {
	r := chi.NewRouter()

	r.Use(chimw.RequestID)
	r.Use(chimw.RealIP)
	r.Use(requestLogger(s.logger))
	r.Use(chimw.Recoverer)
	r.Use(chimw.Compress(5))
	r.Use(securityHeaders)

	r.Handle("/static/*", http.StripPrefix("/static/", http.FileServerFS(s.staticFS)))

	r.Get("/healthz", s.handleHealthz)
	r.Get("/readyz", s.handleReadyz)

	r.Get("/", s.handleHome)
	r.Get("/{section}", s.handleSection)
	r.Get("/{section}/{subsection}", s.handleSubsection)

	return r
}
