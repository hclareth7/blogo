package server

import (
	"context"
	"fmt"
	"io/fs"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/hclareth7/blogo/internal/config"
	"github.com/hclareth7/blogo/internal/navigation"
	"github.com/hclareth7/blogo/internal/parser"
	"github.com/hclareth7/blogo/internal/renderer"
)

type RepoState struct {
	Doc     *parser.Document
	Index   *parser.Index
	NavTree *navigation.NavTree
	NavBld  *navigation.Builder
	Meta    renderer.RepoMeta
}

type Server struct {
	cfg       *config.Config
	repos     map[string]*RepoState
	repoOrder []string
	renderer  *renderer.Renderer
	staticFS  fs.FS
	contentFS http.Handler
	logger    *slog.Logger
}

func New(
	cfg *config.Config,
	repos map[string]*RepoState,
	repoOrder []string,
	rend *renderer.Renderer,
	staticFS fs.FS,
	logger *slog.Logger,
) *Server {
	var contentFS http.Handler
	contentDir := cfg.ContentDir
	if contentDir == "" {
		contentDir = "./content"
	}
	if _, err := os.Stat(contentDir); err == nil {
		contentFS = http.StripPrefix("/static/content/", http.FileServer(http.Dir(contentDir)))
	}

	return &Server{
		cfg:       cfg,
		repos:     repos,
		repoOrder: repoOrder,
		renderer:  rend,
		staticFS:  staticFS,
		contentFS: contentFS,
		logger:    logger,
	}
}

func (s *Server) Start(ctx context.Context) error {
	httpServer := &http.Server{
		Addr:         fmt.Sprintf(":%d", s.cfg.Port),
		Handler:      s.Routes(),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	errCh := make(chan error, 1)
	go func() {
		s.logger.Info("starting server", "port", s.cfg.Port, "repos", len(s.repos))
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- err
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-errCh:
		return fmt.Errorf("server error: %w", err)
	case sig := <-quit:
		s.logger.Info("shutting down", "signal", sig.String())
	case <-ctx.Done():
		s.logger.Info("context cancelled")
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), s.cfg.ShutdownTimeout)
	defer cancel()

	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("shutdown error: %w", err)
	}

	s.logger.Info("server stopped")
	return nil
}

func (s *Server) defaultRepo() *RepoState {
	if len(s.repoOrder) == 0 {
		return nil
	}
	return s.repos[s.repoOrder[0]]
}

func (s *Server) repoList(activeSlug string) []renderer.RepoMeta {
	list := make([]renderer.RepoMeta, 0, len(s.repoOrder))
	for _, slug := range s.repoOrder {
		rs := s.repos[slug]
		list = append(list, renderer.RepoMeta{
			Name:      rs.Meta.Name,
			Slug:      rs.Meta.Slug,
			Author:    rs.Meta.Author,
			Type:      rs.Meta.Type,
			Active:    slug == activeSlug,
			AvatarURL: rs.Meta.AvatarURL,
		})
	}
	return list
}
