package db

import "gorm.io/gorm"

// User represents the object of individual.
type User struct {
	gorm.Model
	Name      string `gorm:"NOT NULL"`
	LoginName string `gorm:"UNIQUE"`
	Email     string `gorm:"NOT NULL"`
	Password  string `gorm:"NOT NULL"`
	Salt      string `gorm:"TYPE:VARCHAR(10)"`

	// Permissions
	IsAdmin bool

	// Avatar
	Avatar      string `gorm:"TYPE:VARCHAR(2048);NOT NULL"`
	AvatarEmail string `gorm:"NOT NULL"`
}
