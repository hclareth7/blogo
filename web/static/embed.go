package static

import "embed"

//go:embed all:css all:js favicon.svg
var FS embed.FS
