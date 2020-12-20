package context

import (
	"github.com/EggMD/EggMD/internal/db"
	"gopkg.in/macaron.v1"
)

func DocumentAssignment() macaron.Handler {
	return func(c *Context) {
		uid := c.Params(":uid")
		doc, err := db.Documents.GetDocByUID(uid)
		if err != nil {
			c.Success("404")
			return
		}

		c.Doc = doc
		c.Data["Doc"] = c.Doc
	}
}
