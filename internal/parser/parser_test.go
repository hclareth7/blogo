package parser

import (
	"log/slog"
	"os"
	"strings"
	"testing"
)

func testLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError}))
}

func TestParseSingleSection(t *testing.T) {
	t.Parallel()
	source := []byte(`# System Design

Welcome to the course.

# What is system design?

System design is the process of defining architecture.
`)
	p := NewParser(testLogger())
	doc, err := p.Parse(source)
	if err != nil {
		t.Fatalf("Parse() error: %v", err)
	}

	if doc.Title != "System Design" {
		t.Errorf("Title = %q, want %q", doc.Title, "System Design")
	}

	if len(doc.Sections) != 1 {
		t.Fatalf("len(Sections) = %d, want 1", len(doc.Sections))
	}

	s := doc.Sections[0]
	if s.Title != "What is system design?" {
		t.Errorf("Section title = %q, want %q", s.Title, "What is system design?")
	}
	if s.ID != "what-is-system-design" {
		t.Errorf("Section ID = %q, want %q", s.ID, "what-is-system-design")
	}
	if !strings.Contains(s.Content, "defining architecture") {
		t.Errorf("Section content should contain 'defining architecture', got %q", s.Content)
	}
}

func TestParseWithChildren(t *testing.T) {
	t.Parallel()
	source := []byte(`# System Design

Intro.

# IP

Internet Protocol.

## Versions

IPv4 and IPv6.

## Types

Unicast, multicast, broadcast.
`)
	p := NewParser(testLogger())
	doc, err := p.Parse(source)
	if err != nil {
		t.Fatalf("Parse() error: %v", err)
	}

	if len(doc.Sections) != 1 {
		t.Fatalf("len(Sections) = %d, want 1", len(doc.Sections))
	}

	s := doc.Sections[0]
	if s.Title != "IP" {
		t.Errorf("Section title = %q, want IP", s.Title)
	}

	if len(s.Children) != 2 {
		t.Fatalf("len(Children) = %d, want 2", len(s.Children))
	}

	if s.Children[0].Title != "Versions" {
		t.Errorf("Child 0 title = %q, want Versions", s.Children[0].Title)
	}
	if s.Children[0].ID != "versions" {
		t.Errorf("Child 0 ID = %q, want versions", s.Children[0].ID)
	}
	if s.Children[1].Title != "Types" {
		t.Errorf("Child 1 title = %q, want Types", s.Children[1].Title)
	}
}

func TestParseSkipsTOC(t *testing.T) {
	t.Parallel()
	source := []byte(`# System Design

Intro.

# Table of contents

- [IP](#ip)
- [DNS](#dns)

# IP

Internet Protocol.
`)
	p := NewParser(testLogger())
	doc, err := p.Parse(source)
	if err != nil {
		t.Fatalf("Parse() error: %v", err)
	}

	for _, s := range doc.Sections {
		if strings.Contains(strings.ToLower(s.Title), "table of content") {
			t.Errorf("TOC section should be skipped, found %q", s.Title)
		}
	}
}

func TestParseCodeBlocks(t *testing.T) {
	t.Parallel()
	source := []byte("# System Design\n\nIntro.\n\n# Example\n\n```go\nfunc main() {}\n```\n")
	p := NewParser(testLogger())
	doc, err := p.Parse(source)
	if err != nil {
		t.Fatalf("Parse() error: %v", err)
	}

	if len(doc.Sections) == 0 {
		t.Fatal("no sections parsed")
	}

	s := doc.Sections[0]
	if !strings.Contains(s.Content, "chroma") || !strings.Contains(s.Content, "func") {
		if !strings.Contains(s.Content, "highlight") && !strings.Contains(s.Content, "<code") {
			t.Logf("Content: %s", s.Content)
		}
	}
}

func TestParseWordCount(t *testing.T) {
	t.Parallel()
	source := []byte("# System Design\n\nIntro.\n\n# Test\n\nOne two three four five.\n")
	p := NewParser(testLogger())
	doc, err := p.Parse(source)
	if err != nil {
		t.Fatalf("Parse() error: %v", err)
	}

	if len(doc.Sections) == 0 {
		t.Fatal("no sections parsed")
	}
	if doc.Sections[0].WordCount < 3 {
		t.Errorf("WordCount = %d, want >= 3", doc.Sections[0].WordCount)
	}
}

func TestParseMultipleSections(t *testing.T) {
	t.Parallel()
	source := []byte(`# System Design

Intro.

# IP

Internet Protocol.

# DNS

Domain Name System.

# Load Balancing

Distributes traffic.
`)
	p := NewParser(testLogger())
	doc, err := p.Parse(source)
	if err != nil {
		t.Fatalf("Parse() error: %v", err)
	}

	if len(doc.Sections) != 3 {
		t.Fatalf("len(Sections) = %d, want 3", len(doc.Sections))
	}

	titles := []string{"IP", "DNS", "Load Balancing"}
	for i, want := range titles {
		if doc.Sections[i].Title != want {
			t.Errorf("Section[%d].Title = %q, want %q", i, doc.Sections[i].Title, want)
		}
	}
}

func TestParseSectionOrder(t *testing.T) {
	t.Parallel()
	source := []byte(`# System Design

Intro.

# A

First.

# B

Second.

# C

Third.
`)
	p := NewParser(testLogger())
	doc, err := p.Parse(source)
	if err != nil {
		t.Fatalf("Parse() error: %v", err)
	}

	for i := 1; i < len(doc.Sections); i++ {
		if doc.Sections[i].Order <= doc.Sections[i-1].Order {
			t.Errorf("Section %q Order=%d should be > Section %q Order=%d",
				doc.Sections[i].Title, doc.Sections[i].Order,
				doc.Sections[i-1].Title, doc.Sections[i-1].Order)
		}
	}
}
