package navigation

import (
	"log/slog"
	"os"
	"testing"

	"github.com/hclareth7/blogo/internal/parser"
)

func testLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError}))
}

func testSections() []*parser.Section {
	return []*parser.Section{
		{ID: "what-is-system-design", Title: "What is system design?", Level: 1, Order: 0},
		{ID: "ip", Title: "IP", Level: 1, Order: 1, Children: []*parser.Section{
			{ID: "versions", Title: "Versions", Level: 2, Order: 2},
			{ID: "types", Title: "Types", Level: 2, Order: 3},
		}},
		{ID: "dns", Title: "DNS", Level: 1, Order: 4},
		{ID: "load-balancing", Title: "Load Balancing", Level: 1, Order: 5},
	}
}

func TestBuildTree(t *testing.T) {
	t.Parallel()
	b := NewBuilderForRepo(testLogger(), "system-design", true)
	tree := b.BuildTree(testSections())

	if tree == nil {
		t.Fatal("BuildTree returned nil")
	}

	totalItems := 0
	for _, g := range tree.Groups {
		totalItems += len(g.Items)
	}
	totalItems += len(tree.Ungrouped)

	if totalItems == 0 {
		t.Error("tree has no items")
	}
}

func TestBuildTreeGrouping(t *testing.T) {
	t.Parallel()
	b := NewBuilderForRepo(testLogger(), "system-design", true)
	sections := []*parser.Section{
		{ID: "what-is-system-design", Title: "What is system design?", Level: 1},
		{ID: "ip", Title: "IP", Level: 1},
		{ID: "load-balancing", Title: "Load Balancing", Level: 1},
		{ID: "sql-databases", Title: "SQL Databases", Level: 1},
	}

	tree := b.BuildTree(sections)

	groupNames := make(map[string]bool)
	for _, g := range tree.Groups {
		groupNames[g.Name] = true
	}

	if !groupNames["Getting Started"] {
		t.Error("missing 'Getting Started' group")
	}
	if !groupNames["Fundamentals"] {
		t.Error("missing 'Fundamentals' group")
	}
}

func TestBuildTreeChildURLs(t *testing.T) {
	t.Parallel()
	b := NewBuilderForRepo(testLogger(), "system-design", true)
	sections := []*parser.Section{
		{ID: "ip", Title: "IP", Level: 1, Children: []*parser.Section{
			{ID: "versions", Title: "Versions", Level: 2},
		}},
	}
	tree := b.BuildTree(sections)

	for _, g := range tree.Groups {
		for _, item := range g.Items {
			if item.ID == "ip" && len(item.Children) > 0 {
				want := "/system-design/ip/versions"
				if item.Children[0].URL != want {
					t.Errorf("child URL = %q, want %q", item.Children[0].URL, want)
				}
			}
		}
	}
}

func TestBuildPrevNextFirst(t *testing.T) {
	t.Parallel()
	b := NewBuilderForRepo(testLogger(), "system-design", true)
	sections := testSections()

	pn := b.BuildPrevNext(sections, "what-is-system-design")

	if pn.Prev != nil {
		t.Errorf("first section should have no Prev, got %q", pn.Prev.Title)
	}
	if pn.Next == nil {
		t.Fatal("first section should have Next")
	}
	if pn.Next.ID != "ip" {
		t.Errorf("Next.ID = %q, want ip", pn.Next.ID)
	}
}

func TestBuildPrevNextLast(t *testing.T) {
	t.Parallel()
	b := NewBuilderForRepo(testLogger(), "system-design", true)
	sections := testSections()

	pn := b.BuildPrevNext(sections, "load-balancing")

	if pn.Next != nil {
		t.Errorf("last section should have no Next, got %q", pn.Next.Title)
	}
	if pn.Prev == nil {
		t.Fatal("last section should have Prev")
	}
}

func TestBuildPrevNextMiddle(t *testing.T) {
	t.Parallel()
	b := NewBuilderForRepo(testLogger(), "system-design", true)
	sections := testSections()

	pn := b.BuildPrevNext(sections, "ip")

	if pn.Prev == nil {
		t.Fatal("middle section should have Prev")
	}
	if pn.Next == nil {
		t.Fatal("middle section should have Next")
	}
}

func TestBuildBreadcrumbsRoot(t *testing.T) {
	t.Parallel()
	b := NewBuilderForRepo(testLogger(), "system-design", true)
	sections := testSections()

	crumbs := b.BuildBreadcrumbs("ip", "", sections)

	if len(crumbs) < 2 {
		t.Fatalf("len(crumbs) = %d, want >= 2", len(crumbs))
	}
	if crumbs[0].Title != "Home" {
		t.Errorf("crumbs[0] = %q, want Home", crumbs[0].Title)
	}
	if !crumbs[len(crumbs)-1].Last {
		t.Error("last crumb should have Last=true")
	}
}

func TestBuildBreadcrumbsChild(t *testing.T) {
	t.Parallel()
	b := NewBuilderForRepo(testLogger(), "system-design", true)
	sections := testSections()

	crumbs := b.BuildBreadcrumbs("versions", "ip", sections)

	if len(crumbs) < 3 {
		t.Fatalf("len(crumbs) = %d, want >= 3", len(crumbs))
	}
	if crumbs[0].Title != "Home" {
		t.Errorf("crumbs[0] = %q, want Home", crumbs[0].Title)
	}
	if crumbs[1].Title != "IP" {
		t.Errorf("crumbs[1] = %q, want IP", crumbs[1].Title)
	}
}
