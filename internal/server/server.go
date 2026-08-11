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

type Server struct {
	cfg      *config.Config
	index    *parser.Index
	doc      *parser.Document
	navTree  *navigation.NavTree
	navBld   *navigation.Builder
	renderer *renderer.Renderer
	staticFS fs.FS
	logger   *slog.Logger
}

func New(
	cfg *config.Config,
	doc *parser.Document,
	index *parser.Index,
	navTree *navigation.NavTree,
	navBld *navigation.Builder,
	rend *renderer.Renderer,
	staticFS fs.FS,
	logger *slog.Logger,
) *Server {
	return &Server{
		cfg:      cfg,
		doc:      doc,
		index:    index,
		navTree:  navTree,
		navBld:   navBld,
		renderer: rend,
		staticFS: staticFS,
		logger:   logger,
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
		s.logger.Info("starting server", "port", s.cfg.Port)
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
