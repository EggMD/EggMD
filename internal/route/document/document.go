package document

const (
	DOCUMENT_EDITOR = "document/editor"
	DOCUMENT_SHARE  = "share/share"
)

//
//// New 新建一个新的文档，并跳转到文件编辑页面。
//func New(c *context.Context) {
//	newDoc, err := db.Documents.Create(c.User.ID)
//	if err != nil {
//		c.Error(err)
//		return
//	}
//	c.RedirectSubPath("/e/" + newDoc.UID)
//}
//
//// Remove 从用户列表中移除一个文档，并跳转到用户仪表盘。
//func Remove(c *context.Context) {
//	uid := c.Params(":uid")
//	defer c.Redirect("/")
//
//	doc, err := db.Documents.GetDocByUID(uid)
//	if err != nil {
//		log.Error("Failed to find doc: %v", err)
//		c.Error(err)
//		return
//	}
//
//	// 文档作者移除文档，将会直接将文档删除。
//	if c.User.ID == doc.OwnerID {
//		err = db.Documents.Remove(uid)
//		if err != nil {
//			log.Error("Failed to remove doc: %v", err)
//			return
//		}
//	}
//	// ⚠️ 文档作者也属于文档协作者，因此在删除文档时也需要删除协作者关系。
//	err = db.Documents.RemoveContributor(c.User.ID, doc.ID)
//	if err != nil {
//		log.Error("Failed to remove doc: %v", err)
//		c.Error(err)
//	}
//}
//
//// Editor 为文档在线编辑器页面。
//func Editor(c *context.Context) {
//	if c.IsLogged {
//		// 添加当前登录用户为文档协作者。
//		err := db.Documents.AppendContributor(c.User.ID, c.Doc.ID)
//		if err != nil {
//			log.Error("append editor error: %v", err)
//			c.Error(err)
//			return
//		}
//	}
//	c.Success(DOCUMENT_EDITOR)
//}
//
//// Share 为文档分享页面。
//func Share(c *context.Context) {
//	c.Data["Doc"] = c.Doc
//	c.Data["RawHTML"], _ = mdutil.RenderMarkdown(c.Doc.Content) // TODO: 缓存 Markdown 渲染数据。
//	c.Success(DOCUMENT_SHARE)
//}
