package parser

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRewriteImagePaths_Markdown(t *testing.T) {
	t.Parallel()
	source := []byte(`# Test
![diagram](./images/arch.png)
![other](images/flow.png)
`)
	result := rewriteImagePaths(source, "/static/content/repo/folder")
	got := string(result)

	if !strings.Contains(got, "/static/content/repo/folder/images/arch.png") {
		t.Errorf("should rewrite ./images/ markdown path, got %s", got)
	}
	if !strings.Contains(got, "/static/content/repo/folder/images/flow.png") {
		t.Errorf("should rewrite images/ markdown path, got %s", got)
	}
}

func TestRewriteImagePaths_HTML(t *testing.T) {
	t.Parallel()
	source := []byte(`# Test
<img src="./images/arch.png" width="400" />
<img src="images/flow.png" />
`)
	result := rewriteImagePaths(source, "/static/content/repo/folder")
	got := string(result)

	if !strings.Contains(got, `src="/static/content/repo/folder/images/arch.png"`) {
		t.Errorf("should rewrite ./images/ HTML path, got %s", got)
	}
	if !strings.Contains(got, `src="/static/content/repo/folder/images/flow.png"`) {
		t.Errorf("should rewrite images/ HTML path, got %s", got)
	}
}

func TestParseMultiFolder_Ordering(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()

	folders := []struct {
		name    string
		content string
	}{
		{"03-third", "# Third Chapter\nContent three."},
		{"01-first", "# First Chapter\nContent one."},
		{"02-second", "# Second Chapter\nContent two."},
	}

	for _, f := range folders {
		fDir := filepath.Join(dir, f.name)
		os.MkdirAll(fDir, 0o755)
		os.WriteFile(filepath.Join(fDir, "README.md"), []byte(f.content), 0o644)
	}

	p := NewParser(testLogger())
	doc, err := p.ParseMultiFolder(dir, "test-repo")
	if err != nil {
		t.Fatalf("ParseMultiFolder() error: %v", err)
	}

	if len(doc.Sections) != 3 {
		t.Fatalf("len(Sections) = %d, want 3", len(doc.Sections))
	}

	expected := []string{"First Chapter", "Second Chapter", "Third Chapter"}
	for i, want := range expected {
		if doc.Sections[i].Title != want {
			t.Errorf("Section[%d].Title = %q, want %q", i, doc.Sections[i].Title, want)
		}
	}
}

func TestParseMultiFolder_RootReadme(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()

	os.WriteFile(filepath.Join(dir, "README.md"), []byte("# My Repo\nIntroduction content."), 0o644)

	fDir := filepath.Join(dir, "01-chapter")
	os.MkdirAll(fDir, 0o755)
	os.WriteFile(filepath.Join(fDir, "README.md"), []byte("# Chapter One\nChapter content."), 0o644)

	p := NewParser(testLogger())
	doc, err := p.ParseMultiFolder(dir, "test-repo")
	if err != nil {
		t.Fatalf("ParseMultiFolder() error: %v", err)
	}

	if len(doc.Sections) != 2 {
		t.Fatalf("len(Sections) = %d, want 2", len(doc.Sections))
	}

	if doc.Sections[0].Title != "My Repo" {
		t.Errorf("first section should be root README, got %q", doc.Sections[0].Title)
	}
	if doc.Sections[1].Title != "Chapter One" {
		t.Errorf("second section should be chapter, got %q", doc.Sections[1].Title)
	}
}

func TestExtractFolderNum_WithPrefix(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		want int
	}{
		{"01-scaling", 1},
		{"02-caching", 2},
		{"24-object-storage", 24},
		{"no-number", 0},
	}
	for _, tt := range tests {
		got := extractFolderNum(tt.name)
		if got != tt.want {
			t.Errorf("extractFolderNum(%q) = %d, want %d", tt.name, got, tt.want)
		}
	}
}

func TestStripNumericPrefix(t *testing.T) {
	t.Parallel()
	tests := []struct {
		input string
		want  string
	}{
		{"01-scaling", "scaling"},
		{"02. Back Of the Envelope", "Back Of the Envelope"},
		{"24-s3-like-object-storage", "s3-like-object-storage"},
		{"simple-name", "simple-name"},
	}
	for _, tt := range tests {
		got := stripNumericPrefix(tt.input)
		if got != tt.want {
			t.Errorf("stripNumericPrefix(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}
