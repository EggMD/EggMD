package public

import (
	"embed"
)

//go:embed assets css js plugins
var FS embed.FS
