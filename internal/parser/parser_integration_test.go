package parser

import (
	"os"
	"testing"
)

func TestParseRealContent(t *testing.T) {
	data, err := os.ReadFile("../../content/README.md")
	if err != nil {
		t.Skip("content/README.md not found, run 'make fetch-content' first")
	}

	p := NewParser(testLogger())
	doc, err := p.Parse(data)
	if err != nil {
		t.Fatalf("Parse() error: %v", err)
	}

	if doc.Title != "System Design" {
		t.Errorf("Title = %q, want System Design", doc.Title)
	}

	if len(doc.Sections) < 40 {
		t.Errorf("len(Sections) = %d, want >= 40", len(doc.Sections))
	}

	slugs := make(map[string]bool)
	for _, s := range doc.Sections {
		if s.Title == "" {
			t.Errorf("Section with empty title at order %d", s.Order)
		}
		if s.ID == "" {
			t.Errorf("Section %q has empty ID", s.Title)
		}
		if slugs[s.ID] {
			t.Errorf("duplicate slug: %q (section %q)", s.ID, s.Title)
		}
		slugs[s.ID] = true

		if s.Content == "" && len(s.Children) == 0 {
			t.Errorf("Section %q has empty content and no children", s.Title)
		}

		for _, c := range s.Children {
			childKey := s.ID + "/" + c.ID
			if slugs[childKey] {
				t.Errorf("duplicate child slug: %q in section %q", c.ID, s.Title)
			}
			slugs[childKey] = true
		}
	}

	t.Logf("Parsed %d sections from %q", len(doc.Sections), doc.Title)
	for _, s := range doc.Sections[:5] {
		t.Logf("  [%s] %s (%d children, %d words)", s.ID, s.Title, len(s.Children), s.WordCount)
	}
}
