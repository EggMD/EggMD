package context

import (
	"github.com/EggMD/EggMD/internal/db"
	"gopkg.in/macaron.v1"
)

func DocumentAssignment() macaron.Handler {
	return func(c *Context) {
		shortID := c.Params(":shortid")
		doc, err := db.Documents.GetDocByShortID(shortID)
		if err != nil {
			c.Success("404")
			return
		}

		c.Doc = doc
		c.Data["Doc"] = c.Doc
	}
}
