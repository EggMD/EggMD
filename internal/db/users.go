package db

import (
	"strings"
	"unicode/utf8"

	"github.com/EggMD/EggMD/internal/cryptoutil"
	"github.com/EggMD/EggMD/internal/strutil"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrEmptyUserName     = errors.New("user name is empty")
	ErrEmailAlreadyUsed  = errors.New("email is already used")
	ErrBadCredentials    = errors.New("bad credentials")
)

// UsersStore is the persistent interface for users.
type UsersStore interface {
	// Authenticate validates username and password.
	Authenticate(email, password string) (*User, error)
	// Create creates a new user and persists to database.
	// It returns ErrUserAlreadyExists when a user with same name already exists,
	// or ErrEmailAlreadyUsed if the email has been used by another user.
	Create(opts CreateUserOpts) (*User, error)
	// GetByEmail returns the user with given email.
	GetByEmail(email string) (*User, error)
	// GetByID returns the user with given ID. It returns ErrUserNotFound when not found.
	GetByID(id uint) (*User, error)
	// GetByUsername returns the user with given username. It returns ErrUserNotFound when not found.
	GetByLoginName(username string) (*User, error)
}

var Users UsersStore

var _ UsersStore = (*users)(nil)

type users struct {
	*gorm.DB
}

type CreateUserOpts struct {
	Name      string
	LoginName string
	Email     string
	Password  string
	Admin     bool
}

func (db *users) Authenticate(email, password string) (*User, error) {
	user := new(User)
	err := db.Where("email = ?", email).First(user).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, errors.Wrap(err, "get user")
	}

	// User found in the database
	if err == nil && user.ValidatePassword(password) {
		return user, nil
	}
	return nil, ErrBadCredentials
}

func (db *users) Create(opts CreateUserOpts) (*User, error) {
	name := opts.Name
	loginName := opts.LoginName
	email := opts.Email

	err := isUsernameAllowed(name)
	if err != nil {
		return nil, err
	}

	_, err = db.GetByLoginName(loginName)
	if err == nil {
		return nil, ErrUserAlreadyExists
	} else if err != ErrUserNotFound {
		return nil, err
	}

	_, err = db.GetByEmail(email)
	if err == nil {
		return nil, ErrEmailAlreadyUsed
	} else if err != ErrUserNotFound {
		return nil, err
	}

	user := &User{
		Name:        name,
		Email:       email,
		Password:    opts.Password,
		LoginName:   opts.LoginName,
		IsAdmin:     opts.Admin,
		Avatar:      cryptoutil.MD5(email),
		AvatarEmail: email,
	}

	user.Salt, err = GetUserSalt()
	if err != nil {
		return nil, err
	}
	user.EncodePassword()

	return user, db.DB.Create(user).Error
}

func (db *users) GetByID(id uint) (*User, error) {
	user := new(User)
	if err := db.Model(&User{}).Where(&User{
		Model: gorm.Model{ID: id},
	}).First(user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return user, nil
}

func (db *users) GetByLoginName(loginName string) (*User, error) {
	user := new(User)
	if err := db.Model(&User{}).Where(&User{
		LoginName: loginName,
	}).First(user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return user, nil
}

func (db *users) GetByEmail(email string) (*User, error) {
	user := new(User)
	if err := db.Model(&User{}).Where(&User{
		Email: email,
	}).First(user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return user, nil
}

// GetUserSalt returns a random user salt token.
func GetUserSalt() (string, error) {
	return strutil.RandomChars(10)
}

// isUsernameAllowed return an error if given name is a reserved name or pattern for users.
func isUsernameAllowed(name string) error {
	name = strings.TrimSpace(strings.ToLower(name))
	if utf8.RuneCountInString(name) == 0 {
		return ErrEmptyUserName
	}
	return nil
}
