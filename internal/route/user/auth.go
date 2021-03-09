package user

import (
	"net/url"

	"github.com/pkg/errors"

	"github.com/EggMD/EggMD/internal/conf"
	"github.com/EggMD/EggMD/internal/context"
	"github.com/EggMD/EggMD/internal/db"
	"github.com/EggMD/EggMD/internal/form"
	"github.com/EggMD/EggMD/internal/tool"
)

const (
	LOGIN   = "user/auth/login"
	SIGNUP  = "user/auth/signup"
	PROFILE = "user/profile/profile"
)

// Login 为用户登录页面。
func Login(c *context.Context) {
	c.Title("登录")

	redirectTo := c.Query("redirect_to")
	if len(redirectTo) > 0 {
		c.SetCookie("redirect_to", redirectTo, 0, conf.Server.SubPath)
	}

	c.Success(LOGIN)
}

// LoginPost 处理用户提交的登录表单。
func LoginPost(c *context.Context, f form.SignIn) {
	c.Title("登录")

	if c.HasError() {
		c.Success(LOGIN)
		return
	}

	u, err := db.Users.Authenticate(f.Email, f.Password)
	if err != nil {
		switch errors.Cause(err) {
		case db.ErrBadCredentials:
			c.FormErr("UserName", "Password")
			c.RenderWithErr("用户名或密码错误", LOGIN, &f)
		default:
			c.Error(err)
		}
		return
	}

	// 登录成功
	_ = c.Session.Set("uid", u.ID)
	_ = c.Session.Set("uname", u.Name)

	// 清除现在的 CSRF Token，强制生成一个新的。
	c.SetCookie(conf.Session.CSRFCookieName, "", -1, conf.Server.SubPath)

	redirectTo, _ := url.QueryUnescape(c.GetCookie("redirect_to"))
	c.SetCookie("redirect_to", "", -1, conf.Server.SubPath)
	if tool.IsSameSiteURLPath(redirectTo) {
		c.Redirect(redirectTo)
		return
	}

	c.RedirectSubPath("/")
}

// SignUp 为用户注册页面。
func SignUp(c *context.Context) {
	c.Title("注册")
	c.Success(SIGNUP)
}

// SignUpPost 处理用户提交的注册表单。
func SignUpPost(c *context.Context, f form.Register) {
	c.Title("注册")

	if c.HasError() {
		c.Success(SIGNUP)
		return
	}

	if f.Password != f.Retype {
		c.FormErr("Password")
		c.RenderWithErr("两次输入的密码不匹配", SIGNUP, &f)
		return
	}

	if _, err := db.Users.Create(db.CreateUserOpts{
		Name:      f.Name,
		LoginName: f.LoginName,
		Email:     f.Email,
		Password:  f.Password,
		Admin:     false,
	}); err != nil {
		switch err {
		case db.ErrUserAlreadyExists:
			c.FormErr("LoginName")
			c.RenderWithErr("用户名已存在", SIGNUP, &f)
		case db.ErrEmailAlreadyUsed:
			c.FormErr("Email")
			c.RenderWithErr("电子邮箱已注册", SIGNUP, &f)
		default:
			c.Error(err)
		}
		return
	}

	c.Flash.Success("注册成功！")
	c.RedirectSubPath("/user/login")
}

// SignOut 处理用户的登出操作。
func SignOut(c *context.Context) {
	_ = c.Session.Flush()
	_ = c.Session.Destory(c.Context)
	c.SetCookie(conf.Session.CSRFCookieName, "", -1, conf.Server.SubPath)
	c.RedirectSubPath("/")
}
