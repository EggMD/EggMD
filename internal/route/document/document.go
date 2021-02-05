package document

import (
	log "unknwon.dev/clog/v2"

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
	if c.IsLogged {
		// Append editor document relation.
		err := db.Documents.AppendEditor(c.User.ID, c.Doc.ID)
		if err != nil {
			log.Error("append editor error: %v", err)
		}
	}

	c.Success(DOCUMENT_EDITOR)
}

func Share(c *context.Context) {
	c.Data["RawHTML"], _ = mdutil.RenderMarkdown(c.Doc.Content)
	c.Success(DOCUMENT_SHARE)
}
