package db

import (
	"context"
	"encoding/json"

	"github.com/pkg/errors"
	"github.com/thanhpk/randstr"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

var (
	// ErrDocumentNotFound is returned when a document is not found.
	ErrDocumentNotFound = errors.New("document not found")
)

// DocumentsStore is the persistent interface for documents.
type DocumentsStore interface {
	// Create creates a new document belongs to the given user.
	Create(ctx context.Context, ownerID uint) (*Document, error)
	// GetByID returns the document with the given ID.
	GetByID(ctx context.Context, id uint) (*Document, error)
	// GetByUID returns the document with the given UID.
	GetByUID(ctx context.Context, uid string) (*Document, error)
	// GetUserDocuments returns the documents belongs to the given user.
	GetUserDocuments(ctx context.Context, userID uint) (DocumentList, error)
	// UpdateContentByUID updates the content of the document with the given ID.
	UpdateContentByUID(ctx context.Context, uid string, content json.RawMessage) error
	// UpdateMetaByUID updates the meta of the document with the given ID.
	UpdateMetaByUID(ctx context.Context, uid string, opts UpdateMetaOptions) error
	// DeleteByUID deletes the document with the given UID.
	DeleteByUID(ctx context.Context, uid string) error
}

func NewDocumentsStore(db *gorm.DB) DocumentsStore {
	return &documents{db}
}

var Documents DocumentsStore

var _ DocumentsStore = (*documents)(nil)

// DocumentList 为文档列表
type DocumentList []*Document

type documents struct {
	*gorm.DB
}

// Document represents the document.
type Document struct {
	gorm.Model
	UID     string         `gorm:"UNIQUE"`
	Title   string         `gorm:"NOT NULL"`
	Content datatypes.JSON `gorm:"type:jsonb" json:"-"`

	OwnerID uint  `gorm:"NOT NULL"`
	Owner   *User `gorm:"-"`

	LastModifiedUserID uint
	LastModifiedUser   *User `gorm:"-"`
}

func (db *documents) Create(ctx context.Context, ownerID uint) (*Document, error) {
	uid := "doc" + randstr.String(24)
	newDocument := &Document{
		UID:                uid,
		Title:              "新的文档",
		Content:            datatypes.JSON("[]"),
		OwnerID:            ownerID,
		LastModifiedUserID: ownerID,
	}

	if err := db.WithContext(ctx).Model(&Document{}).Create(newDocument).Error; err != nil {
		return nil, err
	}

	documents, err := db.loadAttributes(ctx, newDocument)
	if err != nil {
		return nil, errors.Wrap(err, "load attributes")
	}
	if len(documents) == 0 {
		return nil, errors.New("documents is empty after load attributes")
	}
	return documents[0], nil
}

func (db *documents) getBy(ctx context.Context, whereQuery interface{}, whereArgs ...interface{}) (*Document, error) {
	var document Document
	if err := db.WithContext(ctx).Model(&Document{}).Where(whereQuery, whereArgs...).First(&document).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrDocumentNotFound
		}
		return nil, err
	}
	documents, err := db.loadAttributes(ctx, &document)
	if err != nil {
		return nil, errors.Wrap(err, "load attributes")
	}
	if len(documents) == 0 {
		return nil, errors.New("documents is empty after load attributes")
	}
	return documents[0], nil
}

func (db *documents) GetByID(ctx context.Context, id uint) (*Document, error) {
	return db.getBy(ctx, "id = ?", id)
}

func (db *documents) GetByUID(ctx context.Context, uid string) (*Document, error) {
	return db.getBy(ctx, "uid = ?", uid)
}

func (db *documents) GetUserDocuments(ctx context.Context, userID uint) (DocumentList, error) {
	var documents []*Document
	if err := db.WithContext(ctx).Model(&Document{}).Where("owner_id = ?", userID).Find(&documents).Error; err != nil {
		return nil, errors.Wrap(err, "get user documents")
	}
	return db.loadAttributes(ctx, documents...)
}

func (db *documents) UpdateContentByUID(ctx context.Context, uid string, content json.RawMessage) error {
	doc, err := db.GetByUID(ctx, uid)
	if err != nil {
		return errors.Wrap(err, "get document by uid")
	}

	return db.WithContext(ctx).Model(&Document{}).Where("id = ?", doc.ID).Update("content", content).Error
}

type UpdateMetaOptions struct {
	Title string
}

func (db *documents) UpdateMetaByUID(ctx context.Context, uid string, opts UpdateMetaOptions) error {
	doc, err := db.GetByUID(ctx, uid)
	if err != nil {
		return errors.Wrap(err, "get document by uid")
	}

	return db.WithContext(ctx).Model(&Document{}).Where("id = ?", doc.ID).Updates(&Document{
		Title: opts.Title,
	}).Error
}

func (db *documents) DeleteByUID(ctx context.Context, uid string) error {
	doc, err := db.GetByUID(ctx, uid)
	if err != nil {
		return errors.Wrap(err, "get document by uid")
	}

	if err := db.WithContext(ctx).Model(&Document{}).Delete(&Document{}, "id = ?", doc.ID).Error; err != nil {
		return errors.Wrap(err, "delete document")
	}
	return nil
}

func (db *documents) loadAttributes(ctx context.Context, documents ...*Document) ([]*Document, error) {
	if len(documents) == 0 {
		return documents, nil
	}

	userIDSet := make(map[uint]struct{}, len(documents))
	for _, document := range documents {
		userIDSet[document.OwnerID] = struct{}{}
		userIDSet[document.LastModifiedUserID] = struct{}{}
	}
	userIDs := make([]uint, 0, len(userIDSet))
	for userID := range userIDSet {
		userIDs = append(userIDs, userID)
	}

	usersStore := NewUsersStore(db.DB)
	users, err := usersStore.GetByIDs(ctx, userIDs...)
	if err != nil {
		return nil, errors.Wrap(err, "get users by IDs")
	}

	userSet := make(map[uint]*User)
	for _, user := range users {
		user := user
		userSet[user.ID] = user
	}

	for _, document := range documents {
		document.Owner = userSet[document.OwnerID]
		document.LastModifiedUser = userSet[document.LastModifiedUserID]
	}
	return documents, nil
}
