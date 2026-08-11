package parser

import (
	"bytes"
	"log/slog"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/text"

	highlighting "github.com/yuin/goldmark-highlighting/v2"
	chromahtml "github.com/alecthomas/chroma/v2/formatters/html"
)

var skipSections = map[string]bool{
	"table-of-contents": true,
	"table-of-content":  true,
}

type Parser struct {
	md     goldmark.Markdown
	logger *slog.Logger
}

func NewParser(logger *slog.Logger) *Parser {
	md := goldmark.New(
		goldmark.WithExtensions(
			extension.Table,
			extension.Strikethrough,
			extension.Linkify,
			highlighting.NewHighlighting(
				highlighting.WithStyle("dracula"),
				highlighting.WithFormatOptions(
					chromahtml.WithClasses(true),
					chromahtml.WithLineNumbers(false),
				),
			),
		),
		goldmark.WithRendererOptions(
			html.WithUnsafe(),
		),
	)

	return &Parser{md: md, logger: logger}
}

func (p *Parser) Parse(source []byte) (*Document, error) {
	reader := text.NewReader(source)
	tree := p.md.Parser().Parse(reader)

	doc := &Document{}
	slugs := newSlugRegistry()
	order := 0

	var sections []*rawSection
	var current *rawSection

	for node := tree.FirstChild(); node != nil; node = node.NextSibling() {
		heading, ok := node.(*ast.Heading)
		if !ok {
			if current != nil {
				current.nodes = append(current.nodes, node)
			}
			continue
		}

		title := extractHeadingText(heading, source)

		if heading.Level == 1 {
			if current != nil {
				sections = append(sections, current)
			}
			current = &rawSection{
				title: title,
				level: 1,
			}
		} else if heading.Level == 2 && current != nil {
			current.children = append(current.children, &rawSection{
				title: title,
				level: 2,
				parent: current,
			})
		} else if heading.Level >= 3 && current != nil {
			current.nodes = append(current.nodes, node)
			continue
		}

		if heading.Level == 2 && current != nil && len(current.children) > 0 {
			continue
		}
	}
	if current != nil {
		sections = append(sections, current)
	}

	sections = p.distributeNodes(tree, source, sections)

	for _, raw := range sections {
		slug := slugs.Unique(Slugify(raw.title))

		if raw.title == "" || skipSections[slug] {
			continue
		}

		if doc.Title == "" && raw.level == 1 {
			doc.Title = raw.title
		}

		section := &Section{
			ID:    slug,
			Title: raw.title,
			Level: raw.level,
			Order: order,
		}
		order++

		section.Content = PostProcess(p.renderNodes(raw.contentNodes, source))
		section.RawText = p.extractText(raw.contentNodes, source)
		section.WordCount = countWords(section.RawText)

		childSlugs := newSlugRegistry()
		for _, child := range raw.children {
			childSlug := childSlugs.Unique(Slugify(child.title))
			childSection := &Section{
				ID:        childSlug,
				Title:     child.title,
				Level:     2,
				Order:     order,
				Content:   PostProcess(p.renderNodes(child.contentNodes, source)),
				RawText:   p.extractText(child.contentNodes, source),
				WordCount: countWords(p.extractText(child.contentNodes, source)),
			}
			order++
			section.Children = append(section.Children, childSection)
			section.WordCount += childSection.WordCount
		}

		if doc.Title == raw.title && len(doc.Sections) == 0 {
			continue
		}

		doc.Sections = append(doc.Sections, section)
	}

	p.logger.Info("parsed document", "title", doc.Title, "sections", len(doc.Sections))
	return doc, nil
}

type rawSection struct {
	title        string
	level        int
	nodes        []ast.Node
	contentNodes []ast.Node
	children     []*rawSection
	parent       *rawSection
}

func (p *Parser) distributeNodes(tree ast.Node, source []byte, sections []*rawSection) []*rawSection {
	if len(sections) == 0 {
		return sections
	}

	type headingPos struct {
		sectionIdx int
		childIdx   int
		lineStart  int
	}

	var positions []headingPos

	for i, s := range sections {
		positions = append(positions, headingPos{sectionIdx: i, childIdx: -1})
		for j := range s.children {
			positions = append(positions, headingPos{sectionIdx: i, childIdx: j})
		}
	}

	sIdx := 0
	currentSection := sections[0]
	var currentChild *rawSection

	for node := tree.FirstChild(); node != nil; node = node.NextSibling() {
		heading, isHeading := node.(*ast.Heading)
		if isHeading {
			title := extractHeadingText(heading, source)
			if heading.Level == 1 {
				for i, s := range sections {
					if s.title == title {
						sIdx = i
						currentSection = s
						currentChild = nil
						break
					}
				}
				continue
			}
			if heading.Level == 2 && currentSection != nil {
				found := false
				for _, c := range currentSection.children {
					if c.title == title {
						currentChild = c
						found = true
						break
					}
				}
				if found {
					continue
				}
			}
		}

		if currentChild != nil {
			currentChild.contentNodes = append(currentChild.contentNodes, node)
		} else if currentSection != nil {
			_ = sIdx
			currentSection.contentNodes = append(currentSection.contentNodes, node)
		}
	}

	return sections
}

func (p *Parser) renderNodes(nodes []ast.Node, source []byte) string {
	if len(nodes) == 0 {
		return ""
	}

	var buf bytes.Buffer
	renderer := p.md.Renderer()

	for _, node := range nodes {
		if err := renderer.Render(&buf, source, wrapNode(node)); err != nil {
			p.logger.Error("render error", "error", err)
		}
	}
	return buf.String()
}

func (p *Parser) extractText(nodes []ast.Node, source []byte) string {
	var sb strings.Builder
	for _, node := range nodes {
		extractNodeText(node, source, &sb)
	}
	return sb.String()
}

func extractNodeText(node ast.Node, source []byte, sb *strings.Builder) {
	if node.Type() == ast.TypeBlock || node.Type() == ast.TypeInline {
		if tnode, ok := node.(*ast.Text); ok {
			sb.Write(tnode.Segment.Value(source))
			if tnode.SoftLineBreak() || tnode.HardLineBreak() {
				sb.WriteByte(' ')
			}
		}
		if cblock, ok := node.(*ast.CodeSpan); ok {
			for i := 0; i < cblock.ChildCount(); i++ {
				child := nthChild(cblock, i)
				if child != nil {
					extractNodeText(child, source, sb)
				}
			}
			return
		}
	}
	for child := node.FirstChild(); child != nil; child = child.NextSibling() {
		extractNodeText(child, source, sb)
	}
}

func nthChild(node ast.Node, n int) ast.Node {
	child := node.FirstChild()
	for i := 0; i < n && child != nil; i++ {
		child = child.NextSibling()
	}
	return child
}

func extractHeadingText(heading *ast.Heading, source []byte) string {
	var sb strings.Builder
	for child := heading.FirstChild(); child != nil; child = child.NextSibling() {
		extractNodeText(child, source, &sb)
	}
	return strings.TrimSpace(sb.String())
}

func countWords(text string) int {
	return len(strings.Fields(text))
}

type wrappedDoc struct {
	ast.BaseBlock
	child ast.Node
}

func wrapNode(node ast.Node) *wrappedDoc {
	doc := &wrappedDoc{}
	doc.SetLines(node.Lines())

	clone := shallowRef(node)
	doc.AppendChild(doc, clone)
	return doc
}

func shallowRef(node ast.Node) ast.Node {
	return node
}

func (w *wrappedDoc) Kind() ast.NodeKind {
	return ast.KindDocument
}

func (w *wrappedDoc) Dump(source []byte, level int) {}

func (w *wrappedDoc) Type() ast.NodeType {
	return ast.TypeDocument
}
