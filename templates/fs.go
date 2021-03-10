package templates

import (
	"embed"
)

//go:embed base document share status user home.tmpl
var FS embed.FS
