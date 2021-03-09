package db

import (
	"crypto/sha256"
	"crypto/subtle"
	"fmt"

	"golang.org/x/crypto/pbkdf2"
	"gorm.io/gorm"
)

// User 为独立的一个注册用户对象。
type User struct {
	gorm.Model
	Name             string `gorm:"NOT NULL"`
	LoginName        string `gorm:"UNIQUE"`
	Email            string `gorm:"NOT NULL"`
	KeepEmailPrivate bool
	Password         string `gorm:"NOT NULL"`
	Salt             string `gorm:"TYPE:VARCHAR(10)"`

	// 权限
	IsAdmin bool

	// 文档
	Documents []Document `gorm:"many2many:document_users;"`

	// 头像
	Avatar      string `gorm:"TYPE:VARCHAR(2048);NOT NULL"`
	AvatarEmail string `gorm:"NOT NULL"`
}

// GetDocuments 返回属于用户自己的所有文档。
func (u *User) GetDocuments(page, pageSize int) (DocumentList, error) {
	return Documents.GetUserDocuments(&UserDocOptions{
		UserID:    u.ID,
		LoggedUID: u.ID,
		Page:      page,
		PageSize:  pageSize,
	})
}

// GetVisibleDocuments 返回当前用户有权查看的某用户文档列表。
func (u *User) GetVisibleDocuments(loggedUID uint, page, pageSize int) (DocumentList, error) {
	return Documents.GetUserDocuments(&UserDocOptions{
		UserID:    u.ID,
		LoggedUID: loggedUID,
		Page:      page,
		PageSize:  pageSize,
	})
}

// EncodePassword 将密码转换成安全的格式。
func (u *User) EncodePassword() {
	newPasswd := pbkdf2.Key([]byte(u.Password), []byte(u.Salt), 10000, 50, sha256.New)
	u.Password = fmt.Sprintf("%x", newPasswd)
}

// ValidatePassword 检查真实密码与用户输入的密码是否匹配。
func (u *User) ValidatePassword(password string) bool {
	newUser := &User{Password: password, Salt: u.Salt}
	newUser.EncodePassword()
	return subtle.ConstantTimeCompare([]byte(u.Password), []byte(newUser.Password)) == 1
}
