package parser

import "testing"

func TestSlugify(t *testing.T) {
	t.Parallel()
	tests := []struct {
		input string
		want  string
	}{
		{"Load Balancing", "load-balancing"},
		{"What is a CDN?", "what-is-a-cdn"},
		{"N-Tier Architecture", "n-tier-architecture"},
		{"CAP Theorem", "cap-theorem"},
		{"HTTP / HTTPS", "http-https"},
		{"DNS (Domain Name System)", "dns-domain-name-system"},
		{"", ""},
		{"---special---chars---", "special-chars"},
		{"ALL CAPS TITLE", "all-caps-title"},
		{"simple", "simple"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			t.Parallel()
			got := Slugify(tt.input)
			if got != tt.want {
				t.Errorf("Slugify(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestSlugifyTruncation(t *testing.T) {
	t.Parallel()
	long := "this is a very long heading that exceeds the maximum slug length of eighty characters and should be truncated"
	slug := Slugify(long)
	if len(slug) > maxSlugLen {
		t.Errorf("slug length = %d, want <= %d", len(slug), maxSlugLen)
	}
}

func TestSlugRegistryUnique(t *testing.T) {
	t.Parallel()
	r := newSlugRegistry()

	s1 := r.Unique("test")
	if s1 != "test" {
		t.Errorf("first = %q, want test", s1)
	}

	s2 := r.Unique("test")
	if s2 != "test-2" {
		t.Errorf("second = %q, want test-2", s2)
	}

	s3 := r.Unique("test")
	if s3 != "test-3" {
		t.Errorf("third = %q, want test-3", s3)
	}
}

func TestSlugRegistryEmptySlug(t *testing.T) {
	t.Parallel()
	r := newSlugRegistry()
	s := r.Unique("")
	if s != "section" {
		t.Errorf("empty slug = %q, want section", s)
	}
}
