package context

import (
	"github.com/EggMD/EggMD/internal/db"
	"gopkg.in/macaron.v1"
)

func DocumentUIDAssignment() macaron.Handler {
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

func DocumentShortIDAssignment() macaron.Handler {
	return func(c *Context) {
		shortID := c.Params(":shortID")
		doc, err := db.Documents.GetDocByShortID(shortID)
		if err != nil {
			c.Success("404")
			return
		}

		c.Doc = doc
		c.Data["Doc"] = c.Doc
	}
}