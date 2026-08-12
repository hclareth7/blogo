package parser

type Document struct {
	Title    string
	RepoSlug string
	Author   string
	Sections []*Section
}

type Section struct {
	ID        string
	Title     string
	Level     int
	Content   string
	RawText   string
	Children  []*Section
	Order     int
	WordCount int
}
