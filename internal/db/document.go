package db

import (
	"gorm.io/gorm"
)

// Document represents the object of individual.
type Document struct {
	gorm.Model
	Title      string `gorm:"NOT NULL"`
	UID        string `gorm:"UNIQUE"`
	ShortID    string `gorm:"UNIQUE"`
	OwnerID    uint   `gorm:"NOT NULL"`
	Content    string
	Permission uint `gorm:"NOT NULL"`

	LastModifiedUserID uint
	LastModifiedUser   *User `gorm:"-"`
}
