package parser

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
)

func (p *Parser) ParseMultiFolder(repoDir, repoSlug string) (*Document, error) {
	entries, err := os.ReadDir(repoDir)
	if err != nil {
		return nil, fmt.Errorf("reading repo dir %s: %w", repoDir, err)
	}

	type folderEntry struct {
		path  string
		name  string
		order int
	}

	var folders []folderEntry
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		mdPath := findReadme(filepath.Join(repoDir, e.Name()))
		if mdPath == "" {
			continue
		}
		folders = append(folders, folderEntry{
			path:  mdPath,
			name:  e.Name(),
			order: extractFolderNum(e.Name()),
		})
	}

	sort.Slice(folders, func(i, j int) bool {
		return folders[i].order < folders[j].order
	})

	doc := &Document{RepoSlug: repoSlug}
	slugs := newSlugRegistry()
	imgBase := fmt.Sprintf("/static/content/%s", repoSlug)

	rootReadme := findReadme(repoDir)
	if rootReadme != "" {
		source, err := os.ReadFile(rootReadme)
		if err == nil {
			reader := text.NewReader(source)
			tree := p.md.Parser().Parse(reader)

			title := "Introduction"
			node := tree.FirstChild()
			if h, ok := node.(*ast.Heading); ok && h.Level == 1 {
				title = extractHeadingText(h, source)
				node = node.NextSibling()
			}

			rootSlug := slugs.Unique(Slugify(title))
			section := &Section{
				ID:    rootSlug,
				Title: title,
				Level: 1,
				Order: -1,
			}

			var contentNodes []ast.Node
			for node != nil {
				contentNodes = append(contentNodes, node)
				node = node.NextSibling()
			}

			section.Content = PostProcess(p.renderNodes(contentNodes, source))
			section.RawText = p.extractText(contentNodes, source)
			section.WordCount = countWords(section.RawText)

			doc.Title = title
			doc.Sections = append(doc.Sections, section)
		}
	}

	for i, folder := range folders {
		source, err := os.ReadFile(folder.path)
		if err != nil {
			p.logger.Warn("skipping folder", "folder", folder.name, "error", err)
			continue
		}

		title := stripNumericPrefix(folder.name)
		folderSlug := slugs.Unique(Slugify(title))
		folderDir := filepath.Dir(folder.path)
		folderBase := filepath.Base(folderDir)

		source = rewriteImagePaths(source, imgBase+"/"+folderBase)

		reader := text.NewReader(source)
		tree := p.md.Parser().Parse(reader)

		section := &Section{
			ID:    folderSlug,
			Title: title,
			Level: 1,
			Order: i,
		}

		var contentNodes []ast.Node
		childSlugs := newSlugRegistry()

		node := tree.FirstChild()
		for node != nil {
			heading, ok := node.(*ast.Heading)
			if !ok {
				contentNodes = append(contentNodes, node)
				node = node.NextSibling()
				continue
			}

			headingText := extractHeadingText(heading, source)

			if heading.Level == 1 {
				if headingText != "" {
					section.Title = headingText
					title = headingText
				}
				node = node.NextSibling()
				continue
			}

			if heading.Level >= 3 {
				contentNodes = append(contentNodes, node)
				node = node.NextSibling()
				continue
			}

			// H2 — child section
			childSlug := childSlugs.Unique(Slugify(headingText))
			child := &Section{
				ID:    childSlug,
				Title: headingText,
				Level: 2,
				Order: section.Order*100 + len(section.Children),
			}

			var childNodes []ast.Node
			node = node.NextSibling()
			for node != nil {
				if h, ok := node.(*ast.Heading); ok && h.Level <= 2 {
					break
				}
				childNodes = append(childNodes, node)
				node = node.NextSibling()
			}

			child.Content = PostProcess(p.renderNodes(childNodes, source))
			child.RawText = p.extractText(childNodes, source)
			child.WordCount = countWords(child.RawText)
			section.Children = append(section.Children, child)
			section.WordCount += child.WordCount
		}

		section.Content = PostProcess(p.renderNodes(contentNodes, source))
		section.RawText = p.extractText(contentNodes, source)
		section.WordCount += countWords(section.RawText)

		if doc.Title == "" {
			doc.Title = title
		}

		doc.Sections = append(doc.Sections, section)
	}

	p.logger.Info("parsed multi-folder", "slug", repoSlug, "sections", len(doc.Sections))
	return doc, nil
}

func findReadme(dir string) string {
	candidates := []string{"README.md", "Readme.md", "readme.md"}
	for _, name := range candidates {
		path := filepath.Join(dir, name)
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}
	return ""
}

func stripNumericPrefix(name string) string {
	re := regexp.MustCompile(`^\d+[\.\-\s]+\s*`)
	return strings.TrimSpace(re.ReplaceAllString(name, ""))
}

func extractFolderNum(name string) int {
	n := 0
	for _, c := range name {
		if c >= '0' && c <= '9' {
			n = n*10 + int(c-'0')
		} else {
			break
		}
	}
	return n
}

func rewriteImagePaths(source []byte, basePath string) []byte {
	mdRe := regexp.MustCompile(`(\!\[[^\]]*\]\()(\./)?images/`)
	source = mdRe.ReplaceAll(source, []byte("${1}"+basePath+"/images/"))

	htmlRe := regexp.MustCompile(`(src=["'])(\./)?images/`)
	source = htmlRe.ReplaceAll(source, []byte("${1}"+basePath+"/images/"))

	return source
}
