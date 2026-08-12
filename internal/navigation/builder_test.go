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
	b := NewBuilderForRepo(testLogger(), "system-design")
	tree := b.BuildTree(testSections())

	if tree == nil {
		t.Fatal("BuildTree returned nil")
	}

	if len(tree.Ungrouped) != 4 {
		t.Errorf("len(Ungrouped) = %d, want 4", len(tree.Ungrouped))
	}
}

func TestBuildTreeFlatOrder(t *testing.T) {
	t.Parallel()
	b := NewBuilderForRepo(testLogger(), "test-repo")
	sections := []*parser.Section{
		{ID: "intro", Title: "Introduction", Level: 1, Order: 0},
		{ID: "chapter-1", Title: "Chapter 1", Level: 1, Order: 1},
		{ID: "chapter-2", Title: "Chapter 2", Level: 1, Order: 2},
	}

	tree := b.BuildTree(sections)

	if len(tree.Ungrouped) != 3 {
		t.Fatalf("len(Ungrouped) = %d, want 3", len(tree.Ungrouped))
	}
	if tree.Ungrouped[0].Title != "Introduction" {
		t.Errorf("first item = %q, want Introduction", tree.Ungrouped[0].Title)
	}
	if tree.Ungrouped[2].Title != "Chapter 2" {
		t.Errorf("last item = %q, want Chapter 2", tree.Ungrouped[2].Title)
	}
}

func TestBuildTreeChildURLs(t *testing.T) {
	t.Parallel()
	b := NewBuilderForRepo(testLogger(), "system-design")
	sections := []*parser.Section{
		{ID: "ip", Title: "IP", Level: 1, Children: []*parser.Section{
			{ID: "versions", Title: "Versions", Level: 2},
		}},
	}
	tree := b.BuildTree(sections)

	if len(tree.Ungrouped) != 1 {
		t.Fatalf("len(Ungrouped) = %d, want 1", len(tree.Ungrouped))
	}
	item := tree.Ungrouped[0]
	if len(item.Children) == 0 {
		t.Fatal("item should have children")
	}
	want := "/system-design/ip/versions"
	if item.Children[0].URL != want {
		t.Errorf("child URL = %q, want %q", item.Children[0].URL, want)
	}
}

func TestBuildPrevNextFirst(t *testing.T) {
	t.Parallel()
	b := NewBuilderForRepo(testLogger(), "system-design")
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
	b := NewBuilderForRepo(testLogger(), "system-design")
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
	b := NewBuilderForRepo(testLogger(), "system-design")
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
	b := NewBuilderForRepo(testLogger(), "system-design")
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
	b := NewBuilderForRepo(testLogger(), "system-design")
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
