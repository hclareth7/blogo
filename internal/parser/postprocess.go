package parser

import (
	"regexp"
	"strings"
)

var (
	externalLinkRe = regexp.MustCompile(`<a href="(https?://[^"]+)"`)
	imgTagRe       = regexp.MustCompile(`<img `)
)

func PostProcess(html string) string {
	html = addExternalLinkAttrs(html)
	html = addLazyImages(html)
	html = TransformCallouts(html)
	return html
}

func addExternalLinkAttrs(html string) string {
	return externalLinkRe.ReplaceAllStringFunc(html, func(match string) string {
		if strings.Contains(match, "target=") {
			return match
		}
		return match + ` target="_blank" rel="noopener noreferrer"`
	})
}

func addLazyImages(html string) string {
	return imgTagRe.ReplaceAllStringFunc(html, func(match string) string {
		if strings.Contains(match, "loading=") {
			return match
		}
		return `<img loading="lazy" `
	})
}
