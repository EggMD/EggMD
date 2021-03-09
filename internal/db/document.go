package db

import (
	"gorm.io/gorm"
)

const (
	// 	用户权限
	// 	[游客]
	//	可查看: 0 1 3
	//	可编辑: 0
	//
	// 	[注册用户]
	//	可查看: 0 1 2 3 4
	//	可编辑: 0 1 2

	FREELY    = iota // 0 任何人都可以查看、编辑
	EDITABLE         // 1 任何人都可以查看，注册用户可编辑
	LIMITED          // 2 注册用户可以查看、编辑
	LOCKED           // 3 任何人都可以查看，文档作者可编辑
	PROTECTED        // 4 注册用户可以查看，文档作者可编辑
	PRIVATE          // 5 仅文档作者可以查看、编辑
)

// Document 为独立的一篇文档对象。
type Document struct {
	gorm.Model
	Title      string `gorm:"NOT NULL"`
	UID        string `gorm:"UNIQUE"` // 文档 UID，显示在编辑页面的 URL 中
	ShortID    string `gorm:"UNIQUE"` // 文档短链接，用于文档分享链接
	Content    string
	Permission uint `gorm:"NOT NULL"`

	OwnerID uint  `gorm:"NOT NULL"`
	Owner   *User `gorm:"ForeignKey:OwnerID"`

	Users []User `gorm:"many2many:document_users;"`

	LastModifiedUserID uint
	LastModifiedUser   *User `gorm:"-"`
}

// Permission 为文档权限对象。
type Permission struct {
	canRead  bool
	canWrite bool
}

// MakePermission 根据传入的参数构造并返回一个 Permission 对象。
func MakePermission(canRead, canWrite bool) *Permission {
	return &Permission{
		canRead:  canRead,
		canWrite: canWrite,
	}
}

// CanRead 返回该权限下是否含可读权限。
func (p *Permission) CanRead() bool {
	return p.canRead
}

// CanWrite 返回该权限下是否含可写权限。
func (p *Permission) CanWrite() bool {
	return p.canWrite
}

// HasPermission 检查指定 userID 所对应的用户是否有相应的操作权限。
func (d *Document) HasPermission(userID uint) *Permission {
	if userID == d.OwnerID {
		return MakePermission(true, true)
	}

	// 用户 ID 为 0，为未登录游客。
	if userID == 0 {
		switch d.Permission {
		case 0:
			return MakePermission(true, true)
		case 1, 3:
			return MakePermission(true, false)
		}
		return MakePermission(false, false)
	}

	// 注册用户
	switch d.Permission {
	case 0, 1, 2:
		return MakePermission(true, true)
	case 3, 4:
		return MakePermission(true, false)
	}
	return MakePermission(false, false)
}
