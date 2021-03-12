package form

import (
	"github.com/go-macaron/binding"
	"gopkg.in/macaron.v1"
)

// Profile 为个人信息设置
type Profile struct {
	Name             string `binding:"Required;MaxSize(35)" locale:"昵称"`
	Email            string `binding:"Required;Email;MaxSize(254)" locale:"电子邮箱"`
	KeepEmailPrivate bool
	AvatarEmail      string `binding:"Email;MaxSize(254)" locale:"Avatar 电子邮箱"`
}

func (f *Profile) Validate(ctx *macaron.Context, errs binding.Errors) binding.Errors {
	return validate(errs, ctx.Data, f)
}
