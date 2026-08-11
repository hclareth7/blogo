package parser

import "testing"

func TestIndexLookup(t *testing.T) {
	t.Parallel()
	doc := &Document{
		Title: "Test",
		Sections: []*Section{
			{ID: "ip", Title: "IP", Level: 1, Children: []*Section{
				{ID: "versions", Title: "Versions", Level: 2},
			}},
			{ID: "dns", Title: "DNS", Level: 1},
		},
	}

	idx := NewIndex(doc)

	s, ok := idx.Lookup("ip")
	if !ok || s.Title != "IP" {
		t.Errorf("Lookup(ip) = (%v, %v), want (IP, true)", s, ok)
	}

	s, ok = idx.Lookup("ip/versions")
	if !ok || s.Title != "Versions" {
		t.Errorf("Lookup(ip/versions) = (%v, %v), want (Versions, true)", s, ok)
	}

	_, ok = idx.Lookup("nonexistent")
	if ok {
		t.Error("Lookup(nonexistent) should return false")
	}
}

func TestIndexSections(t *testing.T) {
	t.Parallel()
	doc := &Document{
		Title: "Test",
		Sections: []*Section{
			{ID: "a", Title: "A", Level: 1},
			{ID: "b", Title: "B", Level: 1},
			{ID: "c", Title: "C", Level: 1},
		},
	}

	idx := NewIndex(doc)
	sections := idx.Sections()

	if len(sections) != 3 {
		t.Fatalf("len(Sections()) = %d, want 3", len(sections))
	}
}

func TestIndexFindParent(t *testing.T) {
	t.Parallel()
	doc := &Document{
		Title: "Test",
		Sections: []*Section{
			{ID: "ip", Title: "IP", Level: 1, Children: []*Section{
				{ID: "versions", Title: "Versions", Level: 2},
			}},
		},
	}

	idx := NewIndex(doc)

	parent := idx.FindParent("versions")
	if parent == nil || parent.ID != "ip" {
		t.Errorf("FindParent(versions) = %v, want IP", parent)
	}

	parent = idx.FindParent("ip")
	if parent != nil {
		t.Errorf("FindParent(ip) = %v, want nil", parent)
	}
}

func TestIndexRebuild(t *testing.T) {
	t.Parallel()
	doc1 := &Document{
		Title:    "V1",
		Sections: []*Section{{ID: "a", Title: "A", Level: 1}},
	}
	doc2 := &Document{
		Title:    "V2",
		Sections: []*Section{{ID: "b", Title: "B", Level: 1}},
	}

	idx := NewIndex(doc1)

	_, ok := idx.Lookup("a")
	if !ok {
		t.Error("should find 'a' before rebuild")
	}

	idx.Rebuild(doc2)

	_, ok = idx.Lookup("a")
	if ok {
		t.Error("should not find 'a' after rebuild")
	}

	_, ok = idx.Lookup("b")
	if !ok {
		t.Error("should find 'b' after rebuild")
	}
}
