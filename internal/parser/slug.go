package parser

import (
	"regexp"
	"strings"
)

var (
	nonAlphaNum    = regexp.MustCompile(`[^a-z0-9-]+`)
	multipleHyphen = regexp.MustCompile(`-{2,}`)
)

const maxSlugLen = 80

func Slugify(text string) string {
	s := strings.ToLower(text)
	s = strings.ReplaceAll(s, " ", "-")
	s = nonAlphaNum.ReplaceAllString(s, "")
	s = multipleHyphen.ReplaceAllString(s, "-")
	s = strings.Trim(s, "-")
	if len(s) > maxSlugLen {
		s = s[:maxSlugLen]
		s = strings.TrimRight(s, "-")
	}
	return s
}

type slugRegistry struct {
	seen map[string]int
}

func newSlugRegistry() *slugRegistry {
	return &slugRegistry{seen: make(map[string]int)}
}

func (r *slugRegistry) Unique(slug string) string {
	if slug == "" {
		slug = "section"
	}
	r.seen[slug]++
	count := r.seen[slug]
	if count == 1 {
		return slug
	}
	return Slugify(slug + "-" + strings.Repeat("", 0) + itoa(count))
}

func itoa(n int) string {
	if n < 10 {
		return string(rune('0' + n))
	}
	return itoa(n/10) + string(rune('0'+n%10))
}
