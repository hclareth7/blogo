package navigation

import (
	"log/slog"

	"github.com/hclareth7/blogo/internal/parser"
)

type groupDef struct {
	Name     string
	Icon     string
	Sections []string
}

var defaultGroups = []groupDef{
	{
		Name: "Getting Started",
		Icon: "book-open",
		Sections: []string{
			"what-is-system-design",
		},
	},
	{
		Name: "Fundamentals",
		Icon: "book",
		Sections: []string{
			"ip", "osi-model", "tcp-and-udp", "domain-name-system-dns",
			"load-balancing", "clustering", "caching",
			"content-delivery-network-cdn", "proxy",
			"availability", "scalability", "storage",
		},
	},
	{
		Name: "Databases",
		Icon: "database",
		Sections: []string{
			"databases-and-dbms", "sql-databases", "nosql-databases",
			"sql-vs-nosql-databases", "database-replication", "indexes",
			"normalization-and-denormalization", "acid-and-base-consistency-models",
			"cap-theorem", "pacelc-theorem", "transactions",
			"distributed-transactions", "sharding", "consistent-hashing",
			"database-federation",
		},
	},
	{
		Name: "Architecture",
		Icon: "layers",
		Sections: []string{
			"n-tier-architecture", "message-brokers", "message-queues",
			"publish-subscribe", "enterprise-service-bus-esb",
			"monoliths-and-microservices", "event-driven-architecture-eda",
			"event-sourcing",
			"command-and-query-responsibility-segregation-cqrs",
			"api-gateway", "rest-graphql-grpc",
			"long-polling-websockets-server-sent-events-sse",
		},
	},
	{
		Name: "Infrastructure",
		Icon: "server",
		Sections: []string{
			"geohashing-and-quadtrees", "circuit-breaker", "rate-limiting",
			"service-discovery", "sla-slo-sli", "disaster-recovery",
			"virtual-machines-vms-and-containers",
			"oauth-20-and-openid-connect-oidc",
			"single-sign-on-sso", "ssl-tls-mtls",
		},
	},
	{
		Name: "Case Studies",
		Icon: "briefcase",
		Sections: []string{
			"system-design-interviews",
			"url-shortener", "whatsapp", "twitter", "netflix", "uber",
		},
	},
	{
		Name: "Appendix",
		Icon: "bookmark",
		Sections: []string{
			"next-steps", "references",
		},
	},
}

type Builder struct {
	groups    []groupDef
	logger    *slog.Logger
	urlPrefix string
}

func NewBuilder(logger *slog.Logger) *Builder {
	return &Builder{
		groups: defaultGroups,
		logger: logger,
	}
}

func NewBuilderWithPrefix(logger *slog.Logger, urlPrefix string) *Builder {
	return &Builder{
		groups:    nil,
		logger:    logger,
		urlPrefix: urlPrefix,
	}
}

func NewBuilderForRepo(logger *slog.Logger, repoSlug string, useGroups bool) *Builder {
	b := &Builder{
		logger:    logger,
		urlPrefix: "/" + repoSlug,
	}
	if useGroups {
		b.groups = defaultGroups
	}
	return b
}

func (b *Builder) BuildTree(sections []*parser.Section) *NavTree {
	if len(b.groups) == 0 {
		return b.buildFlatTree(sections)
	}

	sectionMap := make(map[string]*parser.Section)
	for _, s := range sections {
		sectionMap[s.ID] = s
	}

	tree := &NavTree{}
	assigned := make(map[string]bool)

	for _, gd := range b.groups {
		group := &NavGroup{
			Name: gd.Name,
			Icon: gd.Icon,
		}

		for _, slug := range gd.Sections {
			s, ok := sectionMap[slug]
			if !ok {
				continue
			}

			item := b.sectionToNavItem(s)
			group.Items = append(group.Items, item)
			assigned[slug] = true
		}

		if len(group.Items) > 0 {
			tree.Groups = append(tree.Groups, group)
		}
	}

	for _, s := range sections {
		if !assigned[s.ID] {
			tree.Ungrouped = append(tree.Ungrouped, b.sectionToNavItem(s))
		}
	}

	return tree
}

func (b *Builder) buildFlatTree(sections []*parser.Section) *NavTree {
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
