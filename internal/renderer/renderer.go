package renderer

import (
	"html/template"
	"io"
	"io/fs"
	"log/slog"
	"strings"

	"github.com/hclareth7/blogo/internal/navigation"
	"github.com/hclareth7/blogo/internal/parser"
)

type RepoMeta struct {
	Name      string
	Slug      string
	Author    string
	Type      string
	Active    bool
	AvatarURL string
}

type PageData struct {
	Title       string
	Section     *parser.Section
	NavTree     *navigation.NavTree
	Breadcrumbs []navigation.Breadcrumb
	PrevNext    *navigation.PrevNext
	ReadingTime int
	Author      string
	AvatarURL   string
	OriginalURL string
	MetaDesc    string
	CurrentPath string
	DocTitle    string
	RepoSlug    string
	RepoList    []RepoMeta
	IsMultiRepo bool
}

type Renderer struct {
	templates *template.Template
	logger    *slog.Logger
}

func New(templateFS fs.FS, logger *slog.Logger) (*Renderer, error) {
	funcMap := template.FuncMap{
		"safeHTML":    func(s string) template.HTML { return template.HTML(s) },
		"readingTime": readingTime,
		"truncate":    truncate,
		"isActive":    isActive,
		"hasPrefix":   strings.HasPrefix,
		"add":         func(a, b int) int { return a + b },
		"initial":     initial,
	}

	tmpl, err := template.New("").Funcs(funcMap).ParseFS(templateFS, "*.html", "partials/*.html")
	if err != nil {
		return nil, err
	}

	return &Renderer{
		templates: tmpl,
		logger:    logger,
	}, nil
}

func (r *Renderer) RenderPage(w io.Writer, name string, data *PageData) error {
	return r.templates.ExecuteTemplate(w, name, data)
}

func readingTime(wordCount int) int {
	minutes := wordCount / 200
	if minutes < 1 {
		return 1
	}
	return minutes
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "..."
}

func initial(s string) string {
	for _, r := range s {
		return strings.ToUpper(string(r))
	}
	return ""
}

func isActive(currentPath, itemURL string) bool {
	if currentPath == itemURL {
		return true
	}
	if itemURL != "/" && strings.HasPrefix(currentPath, itemURL+"/") {
		return true
	}
	return false
}
