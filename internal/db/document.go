package db

import "gorm.io/gorm"

// Document represents the object of individual.
type Document struct {
	gorm.Model
	Title      string `gorm:"NOT NULL"`
	ShortID    string `gorm:"UNIQUE"`
	Owner      uint   `gorm:"NOT NULL"`
	Content    string
	Permission uint `gorm:"NOT NULL"`
}
