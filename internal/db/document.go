package db

import (
	"gorm.io/gorm"
)

const (
	// Permission
	// Guest View: 0 1 3 Edit: 0
	// User View: 0 1 2 3 4 Edit: 0 1 2
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

// HasPermission checks if the user has permission to do the operation.
func (d *Document) HasPermission(userID uint) (view, edit bool) {
	// Guest
	if userID == 0 {
		switch d.Permission {
		case 0:
			view = true
			edit = true
		case 1, 3:
			view = true
		}
		return
	}

	// SignedIn User
	switch d.Permission {
	case 0, 1, 2:
		view = true
		edit = true
	case 3, 4:
		view = true
	}
	return
}
