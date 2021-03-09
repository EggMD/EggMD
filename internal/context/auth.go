package context

import (
	"net/http"
	"net/url"

	"github.com/go-macaron/csrf"
	"github.com/go-macaron/session"
	"gopkg.in/macaron.v1"
	log "unknwon.dev/clog/v2"

	"github.com/EggMD/EggMD/internal/conf"
	"github.com/EggMD/EggMD/internal/db"
)

type ToggleOptions struct {
	SignInRequired  bool
	SignOutRequired bool
	AdminRequired   bool
	DisableCSRF     bool
}

func Toggle(options *ToggleOptions) macaron.Handler {
	return func(c *Context) {
		// 检查 CSRF Token。
		if c.Req.Method == http.MethodPost {
			csrf.Validate(c.Context, c.csrf)
			if c.Written() {
				return
			}
		}

		// 已登录用户尝试访问未登录页面，跳转至用户仪表盘。
		if options.SignOutRequired && c.IsLogged && c.Req.RequestURI != "/" {
			c.RedirectSubpath("/")
			return
		}

		if options.SignInRequired {
			// 未登录用户尝试访问需要登录的页面，跳转到用户登录页面。
			if !c.IsLogged {
				c.SetCookie("redirect_to", url.QueryEscape(conf.Server.SubPath+c.Req.RequestURI), 0, conf.Server.SubPath)
				c.RedirectSubpath("/user/login")
				return
			}
		}

		// 非管理员权限用户访问管理员页面，提示禁止访问。
		if options.AdminRequired {
			if !c.User.IsAdmin {
				c.Status(http.StatusForbidden)
				return
			}
			c.PageIs("Admin")
		}
	}
}

func DocumentToggle() macaron.Handler {
	return func(c *Context) {
		userID := uint(0)
		if c.IsLogged {
			userID = c.User.ID
		}

		permission := c.Doc.HasPermission(userID)
		// 文档无可读权限，跳转至 404 页面。
		if !permission.CanRead() {
			c.NotFound()
			return
		}
	}
}

// authenticatedUser 从当前 Session 中尝试获取已登录用户信息。
// 若用户未登录，则返回 nil。
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
