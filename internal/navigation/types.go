package navigation

type NavItem struct {
	ID       string
	Title    string
	URL      string
	Level    int
	Active   bool
	Expanded bool
	Children []*NavItem
	Order    int
}

type NavTree struct {
	Ungrouped []*NavItem
}

type PrevNext struct {
	Prev *NavItem
	Next *NavItem
}

type Breadcrumb struct {
	Title string
	URL   string
	Last  bool
}
