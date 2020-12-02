package db

import (
	"crypto/sha256"
	"crypto/subtle"
	"fmt"

	"golang.org/x/crypto/pbkdf2"
	"gorm.io/gorm"
)

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

// EncodePassword encodes password to safe format.
func (u *User) EncodePassword() {
	newPasswd := pbkdf2.Key([]byte(u.Password), []byte(u.Salt), 10000, 50, sha256.New)
	u.Password = fmt.Sprintf("%x", newPasswd)
}

// ValidatePassword checks if given password matches the one belongs to the user.
func (u *User) ValidatePassword(password string) bool {
	newUser := &User{Password: password, Salt: u.Salt}
	newUser.EncodePassword()
	return subtle.ConstantTimeCompare([]byte(u.Password), []byte(newUser.Password)) == 1
}
