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

	fetcher := content.NewFetcher(logger)
	p := parser.NewParser(logger)

	repos := make(map[string]*server.RepoState)
	var repoOrder []string

	for _, repoCfg := range cfg.Repos {
		slug := content.RepoSlug(repoCfg.Name)
		logger.Info("loading repo", "name", repoCfg.Name, "slug", slug, "type", repoCfg.Type)

		if cfg.FetchOnStart {
			if err := fetchRepo(context.Background(), fetcher, repoCfg, cfg.ContentDir, slug); err != nil {
				logger.Error("failed to fetch content", "repo", repoCfg.Name, "error", err)
				os.Exit(1)
			}
		}

		doc, err := parseRepo(p, repoCfg, cfg.ContentDir, slug, fetcher)
		if err != nil {
			logger.Error("failed to parse content", "repo", repoCfg.Name, "error", err)
			os.Exit(1)
		}

		doc.RepoSlug = slug
		doc.Author = repoCfg.Author

		index := parser.NewIndex(doc)

		navBuilder := navigation.NewBuilderForRepo(logger, slug)
		navTree := navBuilder.BuildTree(doc.Sections)

		repos[slug] = &server.RepoState{
			Doc:     doc,
			Index:   index,
			NavTree: navTree,
			NavBld:  navBuilder,
			Meta: renderer.RepoMeta{
				Name:      repoCfg.Name,
				Slug:      slug,
				Author:    repoCfg.Author,
				Type:      repoCfg.Type,
				AvatarURL: content.OwnerAvatarURL(repoCfg.URL),
			},
		}
		repoOrder = append(repoOrder, slug)

		logger.Info("repo loaded", "name", repoCfg.Name, "sections", len(doc.Sections))
	}

	rend, err := renderer.New(templates.FS, logger)
	if err != nil {
		logger.Error("failed to initialize renderer", "error", err)
		os.Exit(1)
	}

	srv := server.New(cfg, repos, repoOrder, rend, static.FS, logger)
	if err := srv.Start(context.Background()); err != nil {
		logger.Error("server error", "error", err)
		os.Exit(1)
	}
}

func fetchRepo(ctx context.Context, fetcher *content.Fetcher, repoCfg config.RepoConfig, contentDir, slug string) error {
	switch repoCfg.Type {
	case "multi-folder":
		return fetcher.FetchMultiFolder(ctx, repoCfg.URL, repoCfg.Branch, contentDir, slug)
	default:
		return fetcher.FetchSingleMD(ctx, repoCfg.URL, repoCfg.Branch, contentDir, slug)
	}
}

func parseRepo(p *parser.Parser, repoCfg config.RepoConfig, contentDir, slug string, fetcher *content.Fetcher) (*parser.Document, error) {
	switch repoCfg.Type {
	case "multi-folder":
		repoDir := contentDir + "/" + slug
		return p.ParseMultiFolder(repoDir, slug)
	default:
		source, err := fetcher.ReadSingleMD(contentDir, slug)
		if err != nil {
			return nil, err
		}
		return p.Parse(source)
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
