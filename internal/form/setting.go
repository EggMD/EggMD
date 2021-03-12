package form

import (
	"github.com/go-macaron/binding"
	"gopkg.in/macaron.v1"
)

// ProfileSettings 为个人信息设置
type ProfileSettings struct {
	Name             string `binding:"Required;MaxSize(35)" locale:"昵称"`
	Email            string `binding:"Required;Email;MaxSize(254)" locale:"电子邮箱"`
	KeepEmailPrivate bool
	AvatarEmail      string `binding:"Email;MaxSize(254)" locale:"Avatar 电子邮箱"`
}

func (f *ProfileSettings) Validate(ctx *macaron.Context, errs binding.Errors) binding.Errors {
	return validate(errs, ctx.Data, f)
}

// AccountSettings 为账号信息设置
type AccountSettings struct {
	Password    string `binding:"Required;MaxSize(255)" locale:"当前密码"`
	NewPassword string `binding:"Required;MaxSize(255)" locale:"新的密码"`
	Retype      string
}

func (f *AccountSettings) Validate(ctx *macaron.Context, errs binding.Errors) binding.Errors {
	return validate(errs, ctx.Data, f)
}
