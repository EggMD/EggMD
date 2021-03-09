package db

import (
	"strings"
	"unicode/utf8"

	"github.com/pkg/errors"
	"gorm.io/gorm"

	"github.com/EggMD/EggMD/internal/cryptoutil"
	"github.com/EggMD/EggMD/internal/strutil"
)

var (
	// ErrUserNotFound 为用户不存在错误。
	ErrUserNotFound = errors.New("user not found")
	// ErrUserAlreadyExists 为用户已存在（已重复）错误。
	ErrUserAlreadyExists = errors.New("user already exists")
	// ErrEmptyUserName 为用户名为空错误。
	ErrEmptyUserName = errors.New("user name is empty")
	// ErrEmailAlreadyUsed 为电子邮箱地址已存在错误。
	ErrEmailAlreadyUsed = errors.New("email is already used")
	// ErrBadCredentials 为用户登录凭证错误。
	ErrBadCredentials = errors.New("bad credentials")
)

// UsersStore 是 Users 用户操作的实现接口。
type UsersStore interface {
	// Authenticate 检查根据输入的邮箱与密码验证用户。
	Authenticate(email, password string) (*User, error)
	// Create 在数据库中创建一个新的用户。
	// 若用户名已存在，则会返回 ErrUserAlreadyExists 错误，若电子邮箱已被其他用户使用，则返回 ErrEmailAlreadyUsed 错误。
	Create(opts CreateUserOpts) (*User, error)
	// GetByEmail 根据输入的电子邮箱地址查找用户并返回。
	GetByEmail(email string) (*User, error)
	// GetByID 根据输入的用户 ID 查找用户并返回，若用户不存在则返回 ErrUserNotFound 错误。
	GetByID(id uint) (*User, error)
	// GetByLoginName 根据输入的用户登录名查找用户并返回，若用户不存在则返回 ErrUserNotFound 错误。
	GetByLoginName(loginName string) (*User, error)
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

	if err == nil && user.ValidatePassword(password) {
		return user, nil
	}
	// 用户凭证错误，或用户不存在
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
		Name:             name,
		Email:            email,
		KeepEmailPrivate: true, // 默认隐藏用户电子邮箱地址
		Password:         opts.Password,
		LoginName:        opts.LoginName,
		IsAdmin:          opts.Admin,
		Avatar:           cryptoutil.MD5(email),
		AvatarEmail:      email,
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

// GetUserSalt 返回一个随机生成的盐。
func GetUserSalt() (string, error) {
	return strutil.RandomChars(10)
}

// isUsernameAllowed 检查用户名是否合法。
func isUsernameAllowed(name string) error {
	name = strings.TrimSpace(strings.ToLower(name))
	if utf8.RuneCountInString(name) == 0 {
		return ErrEmptyUserName
	}
	return nil
}
