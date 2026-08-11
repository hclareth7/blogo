package parser

import (
	"regexp"
	"strings"
)

var calloutPattern = regexp.MustCompile(`<blockquote>\s*<p>\s*<strong>(Note|Warning|Tip)</strong>\s*`)

func TransformCallouts(html string) string {
	html = calloutPattern.ReplaceAllStringFunc(html, func(match string) string {
		var calloutType string
		lower := strings.ToLower(match)
		switch {
		case strings.Contains(lower, "note"):
			calloutType = "note"
		case strings.Contains(lower, "warning"):
			calloutType = "warning"
		case strings.Contains(lower, "tip"):
			calloutType = "tip"
		default:
			return match
		}

		title := strings.Title(calloutType)
		return `<div class="callout callout-` + calloutType + `"><div class="callout-title">` +
			calloutIcon(calloutType) + " " + title + `</div><p>`
	})

	html = strings.ReplaceAll(html, "</blockquote>", "</div>")

	return html
}

func calloutIcon(calloutType string) string {
	switch calloutType {
	case "note":
		return `<svg class="inline w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"/></svg>`
	case "warning":
		return `<svg class="inline w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-2.5L13.732 4c-.77-.833-1.964-.833-2.732 0L4.072 16.5c-.77.833.192 2.5 1.732 2.5z"/></svg>`
	case "tip":
		return `<svg class="inline w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9.663 17h4.673M12 3v1m6.364 1.636l-.707.707M21 12h-1M4 12H3m3.343-5.657l-.707-.707m2.828 9.9a5 5 0 117.072 0l-.548.547A3.374 3.374 0 0014 18.469V19a2 2 0 11-4 0v-.531c0-.895-.356-1.754-.988-2.386l-.548-.547z"/></svg>`
	}
	return ""
}
