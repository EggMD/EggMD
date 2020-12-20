package document

import (
	"github.com/EggMD/EggMD/internal/context"
	"github.com/EggMD/EggMD/internal/db"
	"github.com/EggMD/EggMD/internal/mdutil"
)

const (
	DOCUMENT_EDITOR = "document/editor"
	DOCUMENT_SHARE  = "share/share"
)

func New(c *context.Context) {
	newDoc, err := db.Documents.Create(c.User.ID)
	if err != nil {
		c.Error(500, err.Error())
		return
	}
	c.RedirectSubpath("/" + newDoc.UID)
}

func Editor(c *context.Context) {
	c.Success(DOCUMENT_EDITOR)
}

func Share(c *context.Context) {
	c.Data["RawHTML"], _ = mdutil.RenderMarkdown(c.Doc.Content)
	c.Success(DOCUMENT_SHARE)
}
