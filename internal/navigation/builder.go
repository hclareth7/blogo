package navigation

import (
	"log/slog"

	"github.com/hclareth7/blogo/internal/parser"
)

type Builder struct {
	logger    *slog.Logger
	urlPrefix string
}

func NewBuilder(logger *slog.Logger) *Builder {
	return &Builder{
		logger: logger,
	}
}

func NewBuilderForRepo(logger *slog.Logger, repoSlug string) *Builder {
	return &Builder{
		logger:    logger,
		urlPrefix: "/" + repoSlug,
	}
}

func (b *Builder) BuildTree(sections []*parser.Section) *NavTree {
	tree := &NavTree{}
	for _, s := range sections {
		tree.Ungrouped = append(tree.Ungrouped, b.sectionToNavItem(s))
	}
	return tree
}

func (b *Builder) BuildPrevNext(sections []*parser.Section, currentID string) *PrevNext {
	flat := flattenSections(sections)
	pn := &PrevNext{}

	for i, item := range flat {
		if item.ID != currentID {
			continue
		}
		if i > 0 {
			pn.Prev = &NavItem{
				ID:    flat[i-1].ID,
				Title: flat[i-1].Title,
				URL:   b.sectionURL(flat[i-1]),
			}
		}
		if i < len(flat)-1 {
			pn.Next = &NavItem{
				ID:    flat[i+1].ID,
				Title: flat[i+1].Title,
				URL:   b.sectionURL(flat[i+1]),
			}
		}
		break
	}

	return pn
}

func (b *Builder) BuildBreadcrumbs(sectionID, parentID string, sections []*parser.Section) []Breadcrumb {
	crumbs := []Breadcrumb{
		{Title: "Home", URL: "/"},
	}

	if parentID != "" {
		for _, s := range sections {
			if s.ID == parentID {
				crumbs = append(crumbs, Breadcrumb{
					Title: s.Title,
					URL:   b.urlPrefix + "/" + s.ID,
				})
				break
			}
		}
	}

	for _, s := range sections {
		if s.ID == sectionID {
			crumbs = append(crumbs, Breadcrumb{
				Title: s.Title,
				URL:   b.sectionURL(s),
				Last:  true,
			})
			break
		}
		for _, c := range s.Children {
			if c.ID == sectionID {
				if parentID == "" {
					crumbs = append(crumbs, Breadcrumb{
						Title: s.Title,
						URL:   b.urlPrefix + "/" + s.ID,
					})
				}
				crumbs = append(crumbs, Breadcrumb{
					Title: c.Title,
					Last:  true,
				})
				return crumbs
			}
		}
	}

	if len(crumbs) > 0 {
		crumbs[len(crumbs)-1].Last = true
	}

	return crumbs
}

func (b *Builder) FlatSections(sections []*parser.Section) []*parser.Section {
	return flattenSections(sections)
}

func (b *Builder) sectionToNavItem(s *parser.Section) *NavItem {
	item := &NavItem{
		ID:    s.ID,
		Title: s.Title,
		URL:   b.urlPrefix + "/" + s.ID,
		Level: s.Level,
		Order: s.Order,
	}

	for _, c := range s.Children {
		item.Children = append(item.Children, &NavItem{
			ID:    c.ID,
			Title: c.Title,
			URL:   b.urlPrefix + "/" + s.ID + "/" + c.ID,
			Level: c.Level,
			Order: c.Order,
		})
	}

	return item
}

func (b *Builder) sectionURL(s *parser.Section) string {
	return b.urlPrefix + "/" + s.ID
}

func flattenSections(sections []*parser.Section) []*parser.Section {
	var flat []*parser.Section
	for _, s := range sections {
		flat = append(flat, s)
	}
	return flat
}
