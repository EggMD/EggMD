package db

import (
	"github.com/EggMD/EggMD/internal/strutil"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

var (
	ErrDocumentNotFound = errors.New("document not found")
)

// DocumentsStore is the persistent interface for users.
type DocumentsStore interface {
	Create(ownerID uint) (*Document, error)

	GetDocByShortID(shortID string) (*Document, error)
	UpdateByShortID(shortID string, opts UpdateDocOptions) error
	GetUserDocuments(opts *UserDocOptions) (DocumentList, error)
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
		ShortID:            shortID,
		OwnerID:            ownerID,
		Content:            "",
		LastModifiedUserID: ownerID,
		Permission:         0,
	}
	err = db.Model(&Document{}).Create(d).Error
	return d, err
}

func (db *documents) GetDocByShortID(shortID string) (*Document, error) {
	d := new(Document)
	err := db.Model(&Document{}).Where("short_id = ?", shortID).First(d).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrDocumentNotFound
		}
		return nil, err
	}
	return d, nil
}

type UpdateDocOptions struct {
	Title              string
	Content            string
	LastModifiedUserID uint
}

func (db *documents) UpdateByShortID(shortID string, opts UpdateDocOptions) error {
	tx := db.Begin()
	sess := tx.Model(&Document{}).Where("short_id = ?", shortID).Updates(map[string]interface{}{
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

	// Load modified users
	userSet := make(map[uint]*User)
	for i := range docs {
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
	}

	return nil
}

func (db *documents) GetUserDocuments(opts *UserDocOptions) (DocumentList, error) {
	if opts.Page <= 0 {
		opts.Page = 1
	}

	docs := make(DocumentList, 0, opts.PageSize)
	err := db.Debug().Model(&Document{}).Where("owner_id = ?", opts.UserID).
		Offset((opts.Page - 1) * opts.PageSize).Limit(opts.PageSize).
		Order("`updated_at` DESC").Find(&docs).Error
	if err != nil {
		return nil, err
	}

	if err = docs.loadAttributes(db.DB); err != nil {
		return nil, err
	}

	return docs, err
}

// GetShortID returns a random user salt token.
func GetShortID() (string, error) {
	return strutil.RandomChars(20)
}
