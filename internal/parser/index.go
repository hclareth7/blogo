package parser

import "sync"

type Index struct {
	mu      sync.RWMutex
	bySlug  map[string]*Section
	byRoute map[string]*Section
	flat    []*Section
}

func NewIndex(doc *Document) *Index {
	idx := &Index{}
	idx.Rebuild(doc)
	return idx
}

func (idx *Index) Rebuild(doc *Document) {
	idx.mu.Lock()
	defer idx.mu.Unlock()

	bySlug := make(map[string]*Section)
	byRoute := make(map[string]*Section)
	var flat []*Section

	for _, s := range doc.Sections {
		bySlug[s.ID] = s
		byRoute[s.ID] = s
		flat = append(flat, s)

		for _, c := range s.Children {
			route := s.ID + "/" + c.ID
			bySlug[c.ID] = c
			byRoute[route] = c
		}
	}

	idx.bySlug = bySlug
	idx.byRoute = byRoute
	idx.flat = flat
}

func (idx *Index) Lookup(path string) (*Section, bool) {
	idx.mu.RLock()
	defer idx.mu.RUnlock()

	s, ok := idx.byRoute[path]
	return s, ok
}

func (idx *Index) LookupBySlug(slug string) (*Section, bool) {
	idx.mu.RLock()
	defer idx.mu.RUnlock()

	s, ok := idx.bySlug[slug]
	return s, ok
}

func (idx *Index) Sections() []*Section {
	idx.mu.RLock()
	defer idx.mu.RUnlock()

	result := make([]*Section, len(idx.flat))
	copy(result, idx.flat)
	return result
}

func (idx *Index) FindParent(childID string) *Section {
	idx.mu.RLock()
	defer idx.mu.RUnlock()

	for _, s := range idx.flat {
		for _, c := range s.Children {
			if c.ID == childID {
				return s
			}
		}
	}
	return nil
}
