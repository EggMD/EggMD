package document

import (
	"github.com/EggMD/EggMD/internal/context"
	"github.com/EggMD/EggMD/internal/db"
)

const (
	DOCUMENT_EDITOR = "document/editor"
)

func New(c *context.Context) {
	newDoc, err := db.Documents.Create(c.User.ID)
	if err != nil {
		c.Error(500, err.Error())
		return
	}
	c.RedirectSubpath("/" + newDoc.ShortID)
}

func Editor(c *context.Context) {
	c.Success(DOCUMENT_EDITOR)
}
