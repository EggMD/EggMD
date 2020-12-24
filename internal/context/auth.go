package context

import (
	"net/http"
	"net/url"

	"github.com/EggMD/EggMD/internal/conf"
	"github.com/EggMD/EggMD/internal/db"
	"github.com/go-macaron/csrf"
	"github.com/go-macaron/session"
	"gopkg.in/macaron.v1"
	log "unknwon.dev/clog/v2"
)

type ToggleOptions struct {
	SignInRequired  bool
	SignOutRequired bool
	AdminRequired   bool
	DisableCSRF     bool
}

func Toggle(options *ToggleOptions) macaron.Handler {
	return func(c *Context) {
		// Check CSRF token.
		if c.Req.Method == "POST" {
			csrf.Validate(c.Context, c.csrf)
			if c.Written() {
				return
			}
		}

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

func DocToggle() macaron.Handler {
	return func(c *Context) {
		switch c.Doc.Permission {
		case db.FREELY: // Anyone can view & edit
			c.Permission.View = true
			c.Permission.Edit = true
		case db.EDITABLE: // Anyone can view, Signed-in people can edit
			c.Permission.View = true
			c.Permission.View = c.IsLogged
		case db.LIMITED: // Signed-in people can view & edit
			if !c.IsLogged {
				c.RedirectSubpath("/")
				return
			}
			c.Permission.View = c.IsLogged
			c.Permission.View = c.IsLogged
		case db.LOCKED: // Anyone can view, Only owner can edit
			c.Permission.View = true
			c.Permission.View = c.IsLogged && c.User.ID == c.Doc.OwnerID
		case db.PROTECTED: // Signed-in people can view, Only owner can edit
			if !c.IsLogged {
				c.RedirectSubpath("/")
				return
			}
			c.Permission.View = c.IsLogged
			c.Permission.View = c.IsLogged && c.User.ID == c.Doc.OwnerID
		case db.PRIVATE: // Only owner can view & edit
			if !c.IsLogged || c.User.ID != c.Doc.OwnerID {
				c.RedirectSubpath("/")
				return
			}
			c.Permission.View = c.IsLogged && c.User.ID == c.Doc.OwnerID
			c.Permission.View = c.IsLogged && c.User.ID == c.Doc.OwnerID
		}
	}
}

// authenticatedUser returns the user object of the authenticated user.
func authenticatedUser(sess session.Store) *db.User {
	uid := sess.Get("uid")
	if uid == nil {
		return nil
	}
	if id, ok := uid.(uint); ok {
		u, err := db.Users.GetByID(id)
		if err != nil {
			if err != db.ErrUserNotFound {
				log.Error("Failed to get user by ID: %v", err)
			}
			return nil
		}
		return u
	}
	return nil
}
