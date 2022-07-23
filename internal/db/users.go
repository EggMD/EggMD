package db

import (
	"context"
	"crypto/sha256"
	"crypto/subtle"
	"fmt"

	"github.com/pkg/errors"
	"github.com/satori/go.uuid"
	"github.com/thanhpk/randstr"
	"golang.org/x/crypto/pbkdf2"
	"gorm.io/gorm"
)

var (
	// ErrUserNotFound is returned when a user is not found.
	ErrUserNotFound = errors.New("user not found")
	// ErrUserAlreadyExists is returned when a user already exists.
	ErrUserAlreadyExists = errors.New("user already exists")
	// ErrLoginNameAlreadyExists is returned when a login name already exists.
	ErrLoginNameAlreadyExists = errors.New("login name already exists")
	// ErrEmailAlreadyExists is returned when an email already exists.
	ErrEmailAlreadyExists = errors.New("email already exists")
	// ErrBadCredentials is returned when a user's credentials are incorrect.
	ErrBadCredentials = errors.New("bad credentials")
)

// UsersStore is the persistent interface for users.
type UsersStore interface {
	// Authenticate validates phone and password.
	// It returns ErrBadCredentials when validate failed.
	Authenticate(ctx context.Context, email, password string) (*User, error)
	// Create creates a new user and persists to database.
	// It returns the user when it created.
	Create(ctx context.Context, opts CreateUserOptions) (*User, error)
	// GetByID returns the user with the given ID.
	GetByID(ctx context.Context, id uint) (*User, error)
	// GetByIDs returns the users with the given IDs.
	GetByIDs(ctx context.Context, ids ...uint) ([]*User, error)
	// GetByUID returns the user with the given UID.
	GetByUID(ctx context.Context, uid string) (*User, error)
	// GetByEmail returns the user with given email.
	GetByEmail(ctx context.Context, email string) (*User, error)
	// Update updates the user with the given ID and options.
	Update(ctx context.Context, id uint, opts UpdateUserOptions) error
	// ChangePassword changes the password of the user.
	ChangePassword(ctx context.Context, id uint, oldPassword, newPassword string) error
	// SetPassword sets the password of the user with the given ID.
	SetPassword(ctx context.Context, id uint, newPassword string) error
}

func NewUsersStore(db *gorm.DB) UsersStore {
	return &users{db}
}

var Users UsersStore

var _ UsersStore = (*users)(nil)

type users struct {
	*gorm.DB
}

// User represents the user.
type User struct {
	gorm.Model `json:"-"`

	UID       string `gorm:"UNIQUE"`
	NickName  string
	LoginName string `gorm:"UNIQUE"`
	Email     string `gorm:"UNIQUE"`
	Password  string `json:"-"`
	Salt      string `json:"-"`
}

// EncodePassword encodes the password with the user's salt.
func (u *User) EncodePassword() {
	newPasswd := pbkdf2.Key([]byte(u.Password), []byte(u.Salt), 10000, 50, sha256.New)
	u.Password = fmt.Sprintf("%x", newPasswd)
}

// ValidatePassword validates the password.
func (u *User) ValidatePassword(password string) bool {
	newUser := &User{Password: password, Salt: u.Salt}
	newUser.EncodePassword()
	return subtle.ConstantTimeCompare([]byte(u.Password), []byte(newUser.Password)) == 1
}

func (db *users) Authenticate(ctx context.Context, email, password string) (*User, error) {
	var user User
	err := db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrBadCredentials
		}
		return nil, errors.Wrap(err, "get user")
	}

	if user.ValidatePassword(password) {
		return &user, nil
	}
	return nil, ErrBadCredentials
}

type CreateUserOptions struct {
	NickName  string
	LoginName string
	Email     string
	Password  string
}

func (db *users) Create(ctx context.Context, opts CreateUserOptions) (*User, error) {
	_, err := db.GetByLoginName(ctx, opts.LoginName)
	if err != nil {
		if !errors.Is(err, ErrUserNotFound) {
			return nil, errors.Wrap(err, "get user by login name")
		}
	} else {
		return nil, ErrLoginNameAlreadyExists
	}

	_, err = db.GetByEmail(ctx, opts.Email)
	if err != nil {
		if !errors.Is(err, ErrUserNotFound) {
			return nil, errors.Wrap(err, "get user by email")
		}
	} else {
		return nil, ErrEmailAlreadyExists
	}

	user := User{
		UID:       uuid.NewV4().String(),
		NickName:  opts.NickName,
		LoginName: opts.LoginName,
		Email:     opts.Email,
		Password:  opts.Password,
	}

	user.Salt = generateUserSalt()
	user.EncodePassword()

	return &user, db.WithContext(ctx).Create(&user).Error
}

func (db *users) getBy(ctx context.Context, whereQuery interface{}, whereArgs ...interface{}) (*User, error) {
	var user User
	if err := db.WithContext(ctx).Model(&User{}).Where(whereQuery, whereArgs...).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (db *users) GetByID(ctx context.Context, id uint) (*User, error) {
	return db.getBy(ctx, "id = ?", id)
}

func (db *users) GetByUID(ctx context.Context, uid string) (*User, error) {
	return db.getBy(ctx, "uid = ?", uid)
}

func (db *users) GetByLoginName(ctx context.Context, loginName string) (*User, error) {
	return db.getBy(ctx, "login_name = ?", loginName)
}

func (db *users) GetByEmail(ctx context.Context, email string) (*User, error) {
	return db.getBy(ctx, "email = ?", email)
}

func (db *users) GetByIDs(ctx context.Context, ids ...uint) ([]*User, error) {
	var users []*User
	if err := db.WithContext(ctx).Model(&User{}).Where("id IN (?)", ids).Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

type UpdateUserOptions struct {
	NickName  string
	LoginName string
}

func (db *users) Update(ctx context.Context, id uint, opts UpdateUserOptions) error {
	_, err := db.GetByID(ctx, id)
	if err != nil {
		return errors.Wrap(err, "get user by ID")
	}

	u, err := db.GetByLoginName(ctx, opts.LoginName)
	if err == nil {
		if u.ID != id {
			return ErrLoginNameAlreadyExists
		}
	} else if err != ErrUserNotFound {
		return errors.Wrap(err, "get user by login name")
	}

	if err := db.WithContext(ctx).Model(&User{}).
		Where("id = ?", id).
		Updates(&User{
			NickName:  opts.NickName,
			LoginName: opts.LoginName,
		}).Error; err != nil {
		return errors.Wrap(err, "update user")
	}

	return nil
}

func (db *users) ChangePassword(ctx context.Context, id uint, oldPassword, newPassword string) error {
	user, err := db.GetByID(ctx, id)
	if err != nil {
		return errors.Wrap(err, "get user by ID")
	}

	if !user.ValidatePassword(oldPassword) {
		return ErrBadCredentials
	}

	user.Password = newPassword
	user.EncodePassword()

	if err := db.WithContext(ctx).Model(&User{}).Where("id = ?", id).Update("password", user.Password).Error; err != nil {
		return errors.Wrap(err, "update user password")
	}
	return nil
}

func (db *users) SetPassword(ctx context.Context, id uint, newPassword string) error {
	user, err := db.GetByID(ctx, id)
	if err != nil {
		return errors.Wrap(err, "get user by ID")
	}

	user.Password = newPassword
	user.EncodePassword()

	if err := db.WithContext(ctx).Model(&User{}).Where("id = ?", id).Update("password", user.Password).Error; err != nil {
		return errors.Wrap(err, "update user password")
	}
	return nil
}

// generateUserSalt generates a random salt.
func generateUserSalt() string {
	return randstr.String(10)
}
