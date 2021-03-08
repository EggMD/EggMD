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
		c.Error(err)
		return
	}
	c.RedirectSubpath("/e/" + newDoc.UID)
}

func Remove(c *context.Context) {
	uid := c.Params(":uid")
	defer c.Redirect("/")

	doc, err := db.Documents.GetDocByUID(uid)
	if err != nil {
		log.Error("find doc: %v", err)
		return
	}

	// Owner will remove the document.
	if c.User.ID == doc.OwnerID {
		err = db.Documents.Remove(uid)
		if err != nil {
			log.Error("Failed to remove doc: %v", err)
		}
		return
	}
	err = db.Documents.RemoveEditor(c.User.ID, doc.ID)
	if err != nil {
		log.Error("Failed to remove doc: %v", err)
	}
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
	c.Data["Doc"] = c.Doc
	c.Data["RawHTML"], _ = mdutil.RenderMarkdown(c.Doc.Content)
	c.Success(DOCUMENT_SHARE)
}
