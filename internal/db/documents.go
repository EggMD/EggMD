package db

import (
	"github.com/pkg/errors"
	"github.com/satori/go.uuid"
	"gorm.io/gorm"

	"github.com/EggMD/EggMD/internal/strutil"
)

var (
	ErrDocumentNotFound = errors.New("document not found")
)

// DocumentsStore is the persistent interface for users.
type DocumentsStore interface {
	// Create creates a new document belongs to one user with the given ownerID.
	Create(ownerID uint) (*Document, error)
	// Remove removes the document by given uid.
	Remove(uid string) error
	// GetDocByUID returns a document with the given uid.
	GetDocByUID(uid string) (*Document, error)
	// UpdateByUID updates a document with the given uid.
	UpdateByUID(uid string, opts UpdateDocOptions) error
	// GetDocByShortID returns a document with the given shortID.
	GetDocByShortID(shortID string) (*Document, error)
	// GetUserDocuments returns a user's document list.
	GetUserDocuments(opts *UserDocOptions) (DocumentList, error)
	// SetPermission sets a document's permission.
	SetPermission(uid string, permission uint) error
	// AppendEditor creates a editor user to document relation.
	AppendEditor(userID, documentID uint) error
	// RemoveEditor remove the editor user relation with the given document.
	RemoveEditor(userID, documentID uint) error
}

var Documents DocumentsStore

var _ DocumentsStore = (*documents)(nil)

type documents struct {
	*gorm.DB
}

func (db *documents) Create(ownerID uint) (*Document, error) {
	shortID, err := GetShortID()
	if err != nil {
		return nil, err
	}

	d := &Document{
		Title:              "",
		UID:                uuid.NewV4().String(),
		ShortID:            shortID,
		OwnerID:            ownerID,
		Content:            "",
		LastModifiedUserID: ownerID,
		Permission:         0,
		Users: []User{{
			Model: gorm.Model{
				ID: ownerID,
			},
		}},
	}
	err = db.Model(&Document{}).Create(d).Error
	return d, err
}

func (db *documents) Remove(uid string) error {
	doc, err := db.GetDocByUID(uid)
	if err != nil {
		return err
	}
	tx := db.Begin()
	tx.Model(doc).Association("Users").Clear()
	tx.Model(&Document{}).Delete(&Document{}, "uid = ?", uid)
	return tx.Commit().Error
}

func (db *documents) GetDocByUID(uid string) (*Document, error) {
	var d Document
	err := db.Preload("Owner").Model(&Document{}).Where("uid = ?", uid).First(&d).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrDocumentNotFound
		}
		return nil, err
	}
	return &d, nil
}

type UpdateDocOptions struct {
	Title              string
	Content            string
	LastModifiedUserID uint
}

func (db *documents) UpdateByUID(uid string, opts UpdateDocOptions) error {
	tx := db.Begin()
	sess := tx.Model(&Document{}).Where("uid = ?", uid).Updates(map[string]interface{}{
		"title":                 opts.Title,
		"content":               opts.Content,
		"last_modified_user_id": opts.LastModifiedUserID,
	})
	if err := sess.Error; err != nil {
		sess.Rollback()
		return err
	}
	if sess.RowsAffected != 1 {
		sess.Rollback()
		return nil
	}
	sess.Commit()
	return nil
}

type UserDocOptions struct {
	UserID   uint
	Page     int
	PageSize int
}

type DocumentList []*Document

func (docs DocumentList) loadAttributes(db *gorm.DB) error {
	if len(docs) == 0 {
		return nil
	}

	// Load modified and owner users
	userSet := make(map[uint]*User)
	for i := range docs {
		userSet[docs[i].OwnerID] = nil
		userSet[docs[i].LastModifiedUserID] = nil
	}
	userIDs := make([]uint, 0, len(userSet))
	for userID := range userSet {
		userIDs = append(userIDs, userID)
	}
	users := make([]*User, 0, len(userIDs))
	if err := db.Model(&User{}).Where("`id` IN ?", userIDs).Find(&users).Error; err != nil {
		return errors.Errorf("find users: %v", err)
	}
	for _, u := range users {
		userSet[u.ID] = u
	}
	for i, d := range docs {
		docs[i].LastModifiedUser = userSet[d.LastModifiedUserID]
		//docs[i].Owner = userSet[d.OwnerID]
	}

	return nil
}

func (db *documents) GetDocByShortID(shortID string) (*Document, error) {
	var d Document
	err := db.Preload("Owner").Model(&Document{}).Where("short_id = ?", shortID).First(&d).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrDocumentNotFound
		}
		return nil, err
	}
	return &d, nil
}

func (db *documents) GetUserDocuments(opts *UserDocOptions) (DocumentList, error) {
	if opts.Page <= 0 {
		opts.Page = 1
	}

	docs := make(DocumentList, 0, opts.PageSize)
	err := db.Preload("Owner").Model(&User{
		Model: gorm.Model{
			ID: opts.UserID,
		},
	}).Offset((opts.Page - 1) * opts.PageSize).Limit(opts.PageSize).
		Order("`updated_at` DESC").Association("Documents").Find(&docs)
	if err != nil {
		return nil, err
	}

	if err = docs.loadAttributes(db.DB); err != nil {
		return nil, err
	}
	return docs, err
}

func (db *documents) SetPermission(uid string, permission uint) error {
	if permission < FREELY || permission > PRIVATE {
		return errors.Errorf("unexpected permission type: %d", permission)
	}

	tx := db.Begin()
	tx.Model(&Document{}).Where("uid = ?", uid).Update("permission", permission)
	return tx.Commit().Error
}

// GetShortID returns a random user salt token.
func GetShortID() (string, error) {
	return strutil.RandomChars(9)
}

func (db *documents) AppendEditor(userID, documentID uint) error {
	return db.Model(&Document{
		Model: gorm.Model{
			ID: documentID,
		},
	}).Association("Users").Append(&User{
		Model: gorm.Model{
			ID: userID,
		},
	})
}

func (db *documents) RemoveEditor(userID, documentID uint) error {
	return db.Model(&Document{
		Model: gorm.Model{
			ID: documentID,
		},
	}).Association("Users").Delete(&User{
		Model: gorm.Model{
			ID: userID,
		},
	})
}
