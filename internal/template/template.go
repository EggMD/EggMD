package template

import (
	"html/template"
	"sync"

	"github.com/EggMD/EggMD/internal/conf"
)

var (
	funcMap     []template.FuncMap
	funcMapOnce sync.Once
)

// FuncMap returns a list of user-defined template functions.
func FuncMap() []template.FuncMap {
	funcMapOnce.Do(func() {
		funcMap = []template.FuncMap{map[string]interface{}{
			"AppSubURL": func() string {
				return conf.Server.Subpath
			},
			"Safe": Safe,
		}}
	})
	return funcMap
}

func Safe(raw string) template.HTML {
	return template.HTML(raw)
}
