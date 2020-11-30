package context

import (
	"net/http"
	"net/url"

	"github.com/EggMD/EggMD/internal/conf"
	"gopkg.in/macaron.v1"
)

type ToggleOptions struct {
	SignInRequired  bool
	SignOutRequired bool
	AdminRequired   bool
	DisableCSRF     bool
}

func Toggle(options *ToggleOptions) macaron.Handler {
	return func(c *Context) {
		// Check non-logged users landing page.
		if !c.IsLogged && c.Req.RequestURI == "/" && conf.Server.LandingURL != "/" {
			c.RedirectSubpath(conf.Server.LandingURL)
			return
		}

		// Redirect to dashboard if user tries to visit any non-login page.
		if options.SignOutRequired && c.IsLogged && c.Req.RequestURI != "/" {
			c.RedirectSubpath("/")
			return
		}

		if options.SignInRequired {
			if !c.IsLogged {
				c.SetCookie("redirect_to", url.QueryEscape(conf.Server.Subpath+c.Req.RequestURI), 0, conf.Server.Subpath)
				c.RedirectSubpath("/user/login")
				return
			}
		}

		if options.AdminRequired {
			if !c.User.IsAdmin {
				c.Status(http.StatusForbidden)
				return
			}
			c.PageIs("Admin")
		}
	}
}
