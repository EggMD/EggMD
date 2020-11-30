package cmd

import (
	"fmt"
	"net/http"

	"github.com/EggMD/EggMD/internal/conf"
	"github.com/EggMD/EggMD/internal/context"
	"github.com/EggMD/EggMD/internal/route"
	"github.com/EggMD/EggMD/internal/template"
	"github.com/go-macaron/csrf"
	"github.com/go-macaron/session"
	"github.com/urfave/cli"
	"gopkg.in/macaron.v1"
	log "unknwon.dev/clog/v2"
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

// newMacaron initializes Macaron instance.
func newMacaron() *macaron.Macaron {
	m := macaron.Classic()

	return m
}

func runWeb(c *cli.Context) error {
	conf.Init()

	m := newMacaron()

	renderOpt := macaron.RenderOptions{
		Directory:         "templates",
		AppendDirectories: []string{"templates"},
		Funcs:             template.FuncMap(),
		IndentJSON:        macaron.Env != macaron.PROD,
	}
	m.Use(macaron.Renderer(renderOpt))

	//reqSignIn := context.Toggle(&context.ToggleOptions{SignInRequired: true})
	//reqSignOut := context.Toggle(&context.ToggleOptions{SignOutRequired: true})
	//
	//bindIgnErr := binding.BindIgnErr

	m.Group("", func() {
		m.Get("/", route.Home)
	},
		session.Sessioner(session.Options{
			CookieName:  conf.Session.CookieName,
			CookiePath:  conf.Server.Subpath,
			Gclifetime:  conf.Session.GCInterval,
			Maxlifetime: conf.Session.MaxLifeTime,
		}),
		csrf.Csrfer(csrf.Options{
			Secret:         conf.Security.SecretKey,
			Header:         "X-CSRF-Token",
			Cookie:         conf.Session.CSRFCookieName,
			CookiePath:     conf.Server.Subpath,
			CookieHttpOnly: true,
			SetCookie:      true,
		}),
		context.Contexter(),
	)

	if c.IsSet("port") {
		conf.Server.HTTPPort = c.String("port")
	} else {
		conf.Server.HTTPPort = "1999"
	}

	listenAddr := fmt.Sprintf("%s:%s", conf.Server.HTTPAddr, conf.Server.HTTPPort)
	log.Info("Listen on http://%s", listenAddr)

	return http.ListenAndServe(listenAddr, m)
}
