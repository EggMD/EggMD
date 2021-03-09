package template

import (
	"html/template"
	"net/url"
	"sync"
	"time"

	"github.com/EggMD/EggMD/internal/conf"
	"github.com/EggMD/EggMD/internal/tool"
)

var (
	funcMap     []template.FuncMap
	funcMapOnce sync.Once
)

// FuncMap 返回用户自定义的模板函数。
func FuncMap() []template.FuncMap {
	funcMapOnce.Do(func() {
		funcMap = []template.FuncMap{map[string]interface{}{
			"AppSubURL": func() string {
				return conf.Server.SubPath
			},
			"Safe": Safe,
			"DateFmtLong": func(t time.Time) string {
				return t.Format(time.RFC1123Z)
			},
			"DateFmtShort": func(t time.Time) string {
				return t.Format("Jan 02, 2006")
			},
			"EscapePound": EscapePound,
			"AvatarLink":  tool.AvatarLink,
		}}
	})
	return funcMap
}

func Safe(raw string) template.HTML {
	return template.HTML(raw)
}

func EscapePound(str string) string {
	return url.PathEscape(str)
}
