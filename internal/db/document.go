package db

import (
	"gorm.io/gorm"
)

const (
	// Permission
	FREELY    = iota // Anyone can view & edit
	EDITABLE         // Anyone can view, Signed-in people can edit
	LIMITED          // Signed-in people can view & edit
	LOCKED           // Anyone can view, Only owner can edit
	PROTECTED        // Signed-in people can view, Only owner can edit
	PRIVATE          // Only owner can view & edit
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
