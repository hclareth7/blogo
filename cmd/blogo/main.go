package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/hclareth7/blogo/internal/config"
	"github.com/hclareth7/blogo/internal/content"
	"github.com/hclareth7/blogo/internal/navigation"
	"github.com/hclareth7/blogo/internal/parser"
	"github.com/hclareth7/blogo/internal/renderer"
	"github.com/hclareth7/blogo/internal/server"
	"github.com/hclareth7/blogo/web/static"
	"github.com/hclareth7/blogo/web/templates"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		slog.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	logger := setupLogger(cfg.LogLevel, cfg.LogFormat)

	fetcher := content.NewFetcher(cfg.ContentURL, cfg.ContentDir, logger)
	if cfg.FetchOnStart {
		if _, err := fetcher.Fetch(context.Background()); err != nil {
			logger.Error("failed to fetch content", "error", err)
			os.Exit(1)
		}
	}

	source, err := fetcher.ReadContent()
	if err != nil {
		logger.Error("failed to read content", "error", err)
		os.Exit(1)
	}

	p := parser.NewParser(logger)
	doc, err := p.Parse(source)
	if err != nil {
		logger.Error("failed to parse content", "error", err)
		os.Exit(1)
	}

	index := parser.NewIndex(doc)

	navBuilder := navigation.NewBuilder(logger)
	navTree := navBuilder.BuildTree(doc.Sections)

	rend, err := renderer.New(templates.FS, logger)
	if err != nil {
		logger.Error("failed to initialize renderer", "error", err)
		os.Exit(1)
	}

	srv := server.New(cfg, doc, index, navTree, navBuilder, rend, static.FS, logger)
	if err := srv.Start(context.Background()); err != nil {
		logger.Error("server error", "error", err)
		os.Exit(1)
	}
}

func setupLogger(level, format string) *slog.Logger {
	var lvl slog.Level
	switch level {
	case "debug":
		lvl = slog.LevelDebug
	case "warn":
		lvl = slog.LevelWarn
	case "error":
		lvl = slog.LevelError
	default:
		lvl = slog.LevelInfo
	}

	opts := &slog.HandlerOptions{Level: lvl}

	var handler slog.Handler
	if format == "json" {
		handler = slog.NewJSONHandler(os.Stdout, opts)
	} else {
		handler = slog.NewTextHandler(os.Stdout, opts)
	}

	return slog.New(handler)
}
