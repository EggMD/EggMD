package context

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/go-macaron/csrf"
	"github.com/go-macaron/session"
	"gopkg.in/macaron.v1"
	log "unknwon.dev/clog/v2"

	"github.com/EggMD/EggMD/internal/conf"
	"github.com/EggMD/EggMD/internal/db"
	"github.com/EggMD/EggMD/internal/form"
	"github.com/EggMD/EggMD/internal/template"
)

// Context 是对 Macaron 中 context 上下文的扩展。
type Context struct {
	*macaron.Context
	csrf    csrf.CSRF
	Flash   *session.Flash
	Session session.Store

	Link     string // 当前请求 URL
	User     *db.User
	IsLogged bool

	// 当前所在文档相关页面，文档信息
	Doc *db.Document
}

// Title 设置模板数据中的 Title 字段。
func (c *Context) Title(title string) {
	c.Data["Title"] = title
}

// PageIs 设置模板数据中的 `PageIsxxx` 字段。
func (c *Context) PageIs(name string) {
	c.Data["PageIs"+name] = true
}

// FormErr 设置模板数据中的 `Err_xxx` 字段。
func (c *Context) FormErr(names ...string) {
	for i := range names {
		c.Data["Err_"+names[i]] = true
	}
}

// HasError 在表单验证发生错误时将返回 true。
func (c *Context) HasError() bool {
	hasErr, ok := c.Data["HasError"]
	if !ok {
		return false
	}
	c.Flash.ErrorMsg = c.Data["ErrorMsg"].(string)
	c.Data["Flash"] = c.Flash
	return hasErr.(bool)
}

// Success 返回 http.StatusOK 状态码模板响应。
func (c *Context) Success(name string) {
	c.HTML(http.StatusOK, name)
}

// Error 返回 http.StatusInternalServerError 状态码的错误响应。
func (c *Context) Error(err error) {
	c.Title("服务内部错误")
	c.HTML(http.StatusInternalServerError, fmt.Sprintf("status/%d", http.StatusInternalServerError))
}

// NotFound 返回页面不存在响应。
func (c *Context) NotFound() {
	c.Title("页面不存在")
	c.HTML(http.StatusNotFound, fmt.Sprintf("status/%d", http.StatusNotFound))
}

// RenderWithErr 用于页面含有表单验证且需要将验证错误返回给用户的场景。
func (c *Context) RenderWithErr(msg, tpl string, f interface{}) {
	if f != nil {
		form.Assign(f, c.Data)
	}
	c.Flash.ErrorMsg = msg
	c.Data["Flash"] = c.Flash

	c.HTML(http.StatusOK, tpl)
}

// RedirectSubPath 根据指定路径与状态码返回页面跳转响应。
// 它会在跳转地址前拼接加上 conf.Server.SubPath 前缀路径。
func (c *Context) RedirectSubPath(location string, status ...int) {
	c.Redirect(conf.Server.SubPath+location, status...)
}

// Contexter 初始化请求的上下文。
func Contexter() macaron.Handler {
	return func(ctx *macaron.Context, sess session.Store, f *session.Flash, x csrf.CSRF) {
		c := &Context{
			Context: ctx,
			csrf:    x,
			Flash:   f,
			Session: sess,
			Link:    conf.Server.SubPath + strings.TrimSuffix(ctx.Req.URL.Path, "/"),
		}

		c.Data["Link"] = template.EscapePound(c.Link)
		c.Data["PageStartTime"] = time.Now()

		// 尝试从 Session 中获取当前登录用户信息。
		c.User = authenticatedUser(c.Session)
		if c.User != nil {
			c.IsLogged = true
			c.Data["IsLogged"] = c.IsLogged
			c.Data["LoggedUser"] = c.User
			c.Data["LoggedUserID"] = c.User.ID
			c.Data["LoggedName"] = c.User.Name
			c.Data["IsAdmin"] = c.User.IsAdmin
		} else {
			c.Data["LoggedUserID"] = 0
			c.Data["LoggedUserName"] = ""
		}

		c.Data["CSRFToken"] = x.GetToken()
		c.Data["CSRFTokenHTML"] = template.Safe(`<input type="hidden" name="_csrf" value="` + x.GetToken() + `">`)

		log.Trace("Session ID: %s", sess.ID())
		log.Trace("CSRF Token: %v", c.Data["CSRFToken"])

		ctx.Map(c)
	}
}
