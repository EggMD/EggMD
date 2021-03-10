package cmd

import (
	"fmt"
	"net/http"

	"github.com/go-macaron/binding"
	"github.com/go-macaron/csrf"
	"github.com/go-macaron/session"
	"github.com/urfave/cli"
	"gopkg.in/macaron.v1"
	log "unknwon.dev/clog/v2"

	"github.com/EggMD/sockets"

	"github.com/EggMD/EggMD/internal/conf"
	"github.com/EggMD/EggMD/internal/context"
	"github.com/EggMD/EggMD/internal/filesystem"
	"github.com/EggMD/EggMD/internal/form"
	"github.com/EggMD/EggMD/internal/route"
	"github.com/EggMD/EggMD/internal/route/document"
	"github.com/EggMD/EggMD/internal/route/user"
	"github.com/EggMD/EggMD/internal/socket"
	"github.com/EggMD/EggMD/internal/template"
	"github.com/EggMD/EggMD/public"
	"github.com/EggMD/EggMD/templates"
)

var Web = cli.Command{
	Name:        "web",
	Usage:       "Start web server",
	Description: "",
	Action:      runWeb,
	Flags: []cli.Flag{
		stringFlag("port, p", "1999", "Temporary port number to prevent conflict"),
	},
}

// newMacaron 初始化一个新的 Macaron 实例。
func newMacaron() *macaron.Macaron {
	m := macaron.New()
	m.Use(macaron.Logger())
	m.Use(macaron.Recovery())
	m.Use(macaron.Statics(macaron.StaticOptions{
		FileSystem: http.FS(public.FS),
	}, "."))
	
	return m
}

func runWeb(c *cli.Context) error {
	m := newMacaron()

	renderOpt := macaron.RenderOptions{
		Funcs:              template.FuncMap(),
		IndentJSON:         macaron.Env != macaron.PROD,
		TemplateFileSystem: filesystem.NewFS(templates.FS),
	}
	m.Use(macaron.Renderer(renderOpt))

	reqSignIn := context.Toggle(&context.ToggleOptions{SignInRequired: true})
	reqSignOut := context.Toggle(&context.ToggleOptions{SignOutRequired: true})

	m.Group("", func() {
		m.Get("/", route.Home)

		m.Group("/user", func() {
			m.Group("/login", func() {
				m.Combo("").Get(user.Login).
					Post(binding.Bind(form.SignIn{}), user.LoginPost)
			})
			m.Group("/sign_up", func() {
				m.Combo("").Get(user.SignUp).
					Post(binding.Bind(form.Register{}), user.SignUpPost)
			})
		}, reqSignOut)

		m.Group("/user", func() {
			m.Post("/logout", user.SignOut)
		})
		m.Get("/user/:name", user.Profile)

		// 文档
		m.Group("/doc", func() {
			m.Post("/new", document.New)
			m.Post("/remove/:uid", document.Remove)
		}, reqSignIn)

		// 文档在线编辑器
		m.Group("/e", func() {
			// 网页前端
			m.Get("/:uid", document.Editor)
			// Websocket 连接
			m.Get("/socket/:uid", sockets.JSON(socket.EventMessage{}), socket.Handler)

		}, context.DocumentUIDAssignment(), context.DocumentToggle())

		// 分享文档
		m.Group("/s/:shortID", func() {
			m.Get("/", document.Share)
		}, context.DocumentShortIDAssignment(), context.DocumentToggle())
	},

		session.Sessioner(session.Options{
			CookieName:  conf.Session.CookieName,
			CookiePath:  conf.Server.SubPath,
			Gclifetime:  conf.Session.GCInterval,
			Maxlifetime: conf.Session.MaxLifeTime,
		}),
		csrf.Csrfer(csrf.Options{
			Secret:         conf.Security.SecretKey,
			Header:         "X-CSRF-Token",
			Cookie:         conf.Session.CSRFCookieName,
			CookiePath:     conf.Server.SubPath,
			CookieHttpOnly: true,
			SetCookie:      true,
		}),
		context.Contexter(),
	)

	m.NotFound(route.NotFound)

	if c.IsSet("port") {
		conf.Server.HTTPPort = c.String("port")
	} else {
		conf.Server.HTTPPort = "1999"
	}

	listenAddr := fmt.Sprintf("%s:%s", conf.Server.HTTPAddr, conf.Server.HTTPPort)
	log.Info("Listen on http://%s", listenAddr)

	return http.ListenAndServe(listenAddr, m)
}
