package context

import (
	"gopkg.in/macaron.v1"

	"github.com/EggMD/EggMD/internal/db"
)

// DocumentUIDAssignment 提取 URL 中的 uid 并用于查找文档是否存在。
func DocumentUIDAssignment() macaron.Handler {
	return func(c *Context) {
		uid := c.Params(":uid")
		doc, err := db.Documents.GetDocByUID(uid)
		if err != nil {
			c.NotFound()
			return
		}

		c.Doc = doc
		c.Data["Doc"] = c.Doc
	}
}

// DocumentUIDAssignment 提取 URL 中的 shortID 并用于查找文档是否存在。
func DocumentShortIDAssignment() macaron.Handler {
	return func(c *Context) {
		shortID := c.Params(":shortID")
		doc, err := db.Documents.GetDocByShortID(shortID)
		if err != nil {
			c.NotFound()
			return
		}

		c.Doc = doc
		c.Data["Doc"] = c.Doc
	}
}
