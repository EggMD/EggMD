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

func Login(c *context.Context) {
	c.Title("登录")

	redirectTo := c.Query("redirect_to")
	if len(redirectTo) > 0 {
		c.SetCookie("redirect_to", redirectTo, 0, conf.Server.Subpath)
	}

	c.Success("user/auth/login")
}

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

	// Login successfully.
	_ = c.Session.Set("uid", u.ID)
	_ = c.Session.Set("uname", u.Name)

	// Clear whatever CSRF has right now, force to generate a new one
	c.SetCookie(conf.Session.CSRFCookieName, "", -1, conf.Server.Subpath)

	redirectTo, _ := url.QueryUnescape(c.GetCookie("redirect_to"))
	c.SetCookie("redirect_to", "", -1, conf.Server.Subpath)
	if tool.IsSameSiteURLPath(redirectTo) {
		c.Redirect(redirectTo)
		return
	}

	c.RedirectSubpath("/")
}

func SignUp(c *context.Context) {
	c.Title("注册")

	c.Success(SIGNUP)
}

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
	c.RedirectSubpath("/user/login")
}

func SignOut(c *context.Context) {
	_ = c.Session.Flush()
	_ = c.Session.Destory(c.Context)
	c.SetCookie(conf.Session.CSRFCookieName, "", -1, conf.Server.Subpath)
	c.RedirectSubpath("/")
}

func Profile(c *context.Context) {
	loginName := c.Params(":name")
	if loginName == "" {
		c.Redirect("/user/login")
		return
	}

	profileUser, err := db.Users.GetByLoginName(loginName)
	if err != nil {
		c.NotFound()
		return
	}
	c.Data["Owner"] = profileUser
	c.Title(profileUser.Name)

	showPrivate := c.IsLogged && (profileUser.ID == c.User.ID || c.User.IsAdmin)
	documents, err := db.Documents.GetUserDocuments(&db.UserDocOptions{
		UserID:      profileUser.ID,
		ShowPrivate: showPrivate,
		Page:        0,
		PageSize:    0,
	})
	if err != nil {
		c.Error(err)
		return
	}
	c.Data["Docs"] = documents

	c.Success(PROFILE)
}
