package user

import (
	"net/url"

	"github.com/EggMD/EggMD/internal/conf"
	"github.com/EggMD/EggMD/internal/context"
)

func Login(c *context.Context) {
	c.Title("登录")

	redirectTo := c.Query("redirect_to")
	if len(redirectTo) > 0 {
		c.SetCookie("redirect_to", redirectTo, 0, conf.Server.Subpath)
	} else {
		redirectTo, _ = url.QueryUnescape(c.GetCookie("redirect_to"))
	}

	c.Success("user/auth/login")
}

func SignUp(c *context.Context) {
	c.Title("注册")

	c.Success("user/auth/signup")
}
