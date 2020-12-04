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
	UpdateByShortID(shortID string, content string)
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
		Title:      "",
		ShortID:    shortID,
		Owner:      ownerID,
		Content:    "",
		Permission: 0,
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

func (db *documents) UpdateByShortID(shortID string, content string) {
	db.Model(&Document{}).Where("short_id = ?", shortID).Update("content", content)
}

// GetShortID returns a random user salt token.
func GetShortID() (string, error) {
	return strutil.RandomChars(20)
}
