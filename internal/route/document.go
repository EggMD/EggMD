// Copyright 2022 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package route

import (
	"net/http"

	"github.com/pkg/errors"
	log "unknwon.dev/clog/v2"

	"github.com/EggMD/EggMD/internal/context"
	"github.com/EggMD/EggMD/internal/db"
	"github.com/EggMD/EggMD/internal/form"
)

var Document documentRouter

type documentRouter struct{}

func (documentRouter) Documenter(ctx context.Context) error {
	uid := ctx.Param("uid")
	document, err := db.Documents.GetByUID(ctx.Request().Context(), uid)
	if err != nil {
		if errors.Is(err, db.ErrDocumentNotFound) {
			return ctx.Error(40400, "文档不存在")
		}
		log.Error("Failed to get document by UID: %v", err)
		return ctx.ServerError()
	}

	if !ctx.IsLogged {
		u := &db.User{}
		ctx.User = u
		ctx.Map(u)
	}

	//if document.OwnerID != ctx.User.ID {
	//	return ctx.Error(40400, "文档不存在")
	//}

	if ctx.Request().Method != http.MethodGet && ctx.User.ID != document.OwnerID {
		return ctx.Error(40300, "无权对文档进行操作")
	}

	ctx.Map(document)
	return nil
}

func (documentRouter) New(ctx context.Context) error {
	document, err := db.Documents.Create(ctx.Request().Context(), ctx.User.ID)
	if err != nil {
		log.Error("Failed to create new document: %v", err)
		return ctx.ServerError()
	}
	return ctx.Success(document.UID)
}

func (documentRouter) List(ctx context.Context) error {
	documents, err := db.Documents.GetUserDocuments(ctx.Request().Context(), ctx.User.ID)
	if err != nil {
		log.Error("Failed to get user documents: %v", err)
		return ctx.ServerError()
	}
	return ctx.Success(documents)
}

func (documentRouter) Meta(ctx context.Context, document *db.Document) error {
	return ctx.Success(document)
}

func (documentRouter) Content(ctx context.Context, document *db.Document) error {
	return ctx.Success(document.Content)
}

func (documentRouter) Save(ctx context.Context, document *db.Document) error {
	requestBody, err := ctx.Request().Body().Bytes()
	if err != nil {
		log.Error("Failed to get request body: %v", err)
		return ctx.ServerError()
	}
	if err := db.Documents.UpdateContentByUID(ctx.Request().Context(), document.UID, requestBody); err != nil {
		log.Error("Failed to update document content: %v", err)
		return ctx.ServerError()
	}
	return ctx.Success()
}

func (documentRouter) Delete(ctx context.Context, document *db.Document) error {
	if err := db.Documents.DeleteByUID(ctx.Request().Context(), document.UID); err != nil {
		log.Error("Failed to delete document by UID: %v", err)
		return ctx.ServerError()
	}
	return ctx.Success("文档已删除")
}

func (documentRouter) GetSetting(ctx context.Context, document *db.Document) error {
	return ctx.Success(map[string]interface{}{
		"Title": document.Title,
	})
}

func (documentRouter) UpdateSetting(ctx context.Context, document *db.Document, f form.UpdateDocumentSetting) error {
	if err := db.Documents.UpdateMetaByUID(ctx.Request().Context(), document.UID, db.UpdateMetaOptions{
		Title: f.Title,
	}); err != nil {
		log.Error("Failed to update document meta: %v", err)
		return ctx.ServerError()
	}
	return ctx.Success("文档设置已更新")
}
