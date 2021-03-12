package form

import (
	"github.com/go-macaron/binding"
	"gopkg.in/macaron.v1"
)

// Register 是用户注册表单。
type Register struct {
	Name      string `binding:"Required;MaxSize(35)" locale:"昵称"`
	LoginName string `binding:"Required;AlphaDashDot;MaxSize(35)" locale:"用户名"`
	Email     string `binding:"Required;Email;MaxSize(254)" locale:"电子邮箱"`
	Password  string `binding:"Required;MaxSize(255)" locale:"密码"`
	Retype    string
}

func (f *Register) Validate(ctx *macaron.Context, errs binding.Errors) binding.Errors {
	return validate(errs, ctx.Data, f)
}

// SignIn 是用户登录表单。
type SignIn struct {
	Email    string `binding:"Required;MaxSize(254)" locale:"电子邮箱"`
	Password string `binding:"Required;MaxSize(255)" locale:"密码"`
}

func (f *SignIn) Validate(ctx *macaron.Context, errs binding.Errors) binding.Errors {
	return validate(errs, ctx.Data, f)
}
