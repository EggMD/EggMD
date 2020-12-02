package context

import (
	"net/http"
	"strings"
	"time"

	"github.com/EggMD/EggMD/internal/conf"
	"github.com/EggMD/EggMD/internal/db"
	"github.com/EggMD/EggMD/internal/template"
	"github.com/go-macaron/csrf"
	"github.com/go-macaron/session"
	"gopkg.in/macaron.v1"
	log "unknwon.dev/clog/v2"
)

// Context represents context of a request.
type Context struct {
	*macaron.Context
	csrf    csrf.CSRF
	Flash   *session.Flash
	Session session.Store

	Link     string // Current request URL
	User     *db.User
	IsLogged bool
}

// Title sets the "Title" field in template data.
func (c *Context) Title(title string) {
	c.Data["Title"] = title
}

// PageIs sets "PageIsxxx" field in template data.
func (c *Context) PageIs(name string) {
	c.Data["PageIs"+name] = true
}

// FormErr sets "Err_xxx" field in template data.
func (c *Context) FormErr(names ...string) {
	for i := range names {
		c.Data["Err_"+names[i]] = true
	}
}

func (c *Context) GetErrMsg() string {
	return c.Data["ErrorMsg"].(string)
}

// HasError returns true if error occurs in form validation.
func (c *Context) HasError() bool {
	hasErr, ok := c.Data["HasError"]
	if !ok {
		return false
	}
	c.Flash.ErrorMsg = c.Data["ErrorMsg"].(string)
	c.Data["Flash"] = c.Flash
	return hasErr.(bool)
}

// Success responses template with status http.StatusOK.
func (c *Context) Success(name string) {
	c.HTML(http.StatusOK, name)
}

// RedirectSubpath responses redirection with given location and status.
// It prepends setting.Server.Subpath to the location string.
func (c *Context) RedirectSubpath(location string, status ...int) {
	c.Redirect(conf.Server.Subpath+location, status...)
}

// Contexter initializes a classic context for a request.
func Contexter() macaron.Handler {
	return func(ctx *macaron.Context, sess session.Store, f *session.Flash, x csrf.CSRF) {
		c := &Context{
			Context: ctx,
			csrf:    x,
			Flash:   f,
			Session: sess,
			Link:    conf.Server.Subpath + strings.TrimSuffix(ctx.Req.URL.Path, "/"),
		}
		c.Data["Link"] = template.EscapePound(c.Link)
		c.Data["PageStartTime"] = time.Now()

		c.Data["CSRFToken"] = x.GetToken()
		c.Data["CSRFTokenHTML"] = template.Safe(`<input type="hidden" name="_csrf" value="` + x.GetToken() + `">`)
		log.Trace("Session ID: %s", sess.ID())
		log.Trace("CSRF Token: %v", c.Data["CSRFToken"])

		ctx.Map(c)
	}
}
