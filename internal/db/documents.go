package db

import (
	"github.com/pkg/errors"
	"github.com/satori/go.uuid"
	"gorm.io/gorm"

	"github.com/EggMD/EggMD/internal/strutil"
)

var (
	// ErrDocumentNotFound 为文档不存在错误。
	ErrDocumentNotFound = errors.New("document not found")
)

// DocumentsStore 是 Documents 文档操作的实现接口。
type DocumentsStore interface {
	// GetDocByUID 根据指定的 uid 查找并返回一篇文档。
	GetDocByUID(uid string) (*Document, error)
	// GetDocByShortID 根据指定的 shortID 查找并返回一篇文档。
	GetDocByShortID(shortID string) (*Document, error)
	// GetUserDocuments 返回属于指定用户的文档列表。
	GetUserDocuments(opts *UserDocOptions) (DocumentList, error)

	// Create 创建一个属于给定 ownerID 用户的新文档。
	Create(ownerID uint) (*Document, error)
	// Remove 根据指定的 uid 删除文档。
	Remove(uid string) error
	// UpdateByUID 根据指定的 uid 查找并更新一篇文档。
	UpdateByUID(uid string, opts UpdateDocOptions) error

	// SetPermission 设置一篇文档的权限。
	SetPermission(uid string, permission uint) error

	// AppendContributor 为一篇文档增加一名协作者。
	AppendContributor(userID, documentID uint) error
	// RemoveContributor 从文档中删除用户 userID 的协作者关系。
	RemoveContributor(userID, documentID uint) error
}

var Documents DocumentsStore

var _ DocumentsStore = (*documents)(nil)

// DocumentList 为文档列表
type DocumentList []*Document

type documents struct {
	*gorm.DB
}

func (db *documents) GetDocByUID(uid string) (*Document, error) {
	var d Document
	if err := db.Preload("Owner").Model(&Document{}).Where("uid = ?", uid).First(&d).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrDocumentNotFound
		}
		return nil, err
	}
	return &d, nil
}

func (db *documents) GetDocByShortID(shortID string) (*Document, error) {
	var d Document
	if err := db.Preload("Owner").Model(&Document{}).Where("short_id = ?", shortID).First(&d).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrDocumentNotFound
		}
		return nil, err
	}
	return &d, nil
}

type UserDocOptions struct {
	UserID      uint
	ShowPrivate bool
	Page        int
	PageSize    int
}

func (db *documents) GetUserDocuments(opts *UserDocOptions) (DocumentList, error) {
	if opts.Page <= 0 {
		opts.Page = 1
	}

	var permission = PROTECTED
	if opts.ShowPrivate {
		permission = PRIVATE
	}

	docs := make(DocumentList, 0, opts.PageSize)
	err := db.Preload("Owner").Model(&User{
		Model: gorm.Model{
			ID: opts.UserID,
		},
	}).Offset((opts.Page-1)*opts.PageSize).Limit(opts.PageSize).
		Order("`updated_at` DESC").Where("`permission` <= ?", permission).Association("Documents").Find(&docs)
	if err != nil {
		return nil, err
	}

	if err = docs.loadAttributes(db.DB); err != nil {
		return nil, err
	}
	return docs, err
}

func (db *documents) Create(ownerID uint) (*Document, error) {
	shortID, err := GetShortID()
	if err != nil {
		return nil, err
	}

	newDocument := &Document{
		Title:              "新的文档",
		UID:                uuid.NewV4().String(),
		ShortID:            shortID,
		OwnerID:            ownerID,
		Content:            "",
		LastModifiedUserID: ownerID,
		Permission:         0, // TODO: 支持设置新建文档时的默认权限
		Users: []User{{
			Model: gorm.Model{
				ID: ownerID,
			},
		}},
	}

	if err := db.Model(&Document{}).Create(newDocument).Error; err != nil {
		return nil, err
	}
	return newDocument, nil
}

func (db *documents) Remove(uid string) error {
	doc, err := db.GetDocByUID(uid)
	if err != nil {
		return err
	}

	tx := db.Begin()
	// 删除所有该文档与注册用户的对应关系。
	if err = tx.Model(doc).Association("Users").Clear(); err != nil {
		tx.Rollback()
		return errors.Wrap(err, "remove document and users association")
	}

	if err := tx.Model(&Document{}).Delete(&Document{}, "uid = ?", uid).Error; err != nil {
		tx.Rollback()
		return errors.Wrap(err, "remove document")
	}
	return tx.Commit().Error
}

type UpdateDocOptions struct {
	Title              string
	Content            string
	LastModifiedUserID uint
}

func (db *documents) UpdateByUID(uid string, opts UpdateDocOptions) error {
	return db.Model(&Document{}).Where("uid = ?", uid).Updates(map[string]interface{}{
		"title":                 opts.Title,
		"content":               opts.Content,
		"last_modified_user_id": opts.LastModifiedUserID,
	}).Error
}

func (docs DocumentList) loadAttributes(db *gorm.DB) error {
	if len(docs) == 0 {
		return nil
	}

	// 补全文档用户信息
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
	}

	return nil
}

func (db *documents) SetPermission(uid string, permission uint) error {
	if permission < FREELY || permission > PRIVATE {
		return errors.Errorf("unexpected permission type: %d", permission)
	}
	return db.Model(&Document{}).Where("uid = ?", uid).Update("permission", permission).Error
}

func (db *documents) AppendContributor(userID, documentID uint) error {
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

func (db *documents) RemoveContributor(userID, documentID uint) error {
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

// GetShortID 生成并返回长度为 9 的随机 shortID。
func GetShortID() (string, error) {
	return strutil.RandomChars(9)
}
