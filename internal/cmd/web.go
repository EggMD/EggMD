package cmd

import (
	"fmt"
	"net/http"

	"github.com/flamego/flamego"
	"github.com/flamego/session"
	"github.com/flamego/session/postgres"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
	log "unknwon.dev/clog/v2"

	"github.com/EggMD/EggMD/internal/conf"
	"github.com/EggMD/EggMD/internal/context"
	"github.com/EggMD/EggMD/internal/db"
	"github.com/EggMD/EggMD/internal/form"
	"github.com/EggMD/EggMD/internal/route"
)

var Web = &cli.Command{
	Name:        "web",
	Usage:       "Start web server",
	Description: "",
	Action:      runWeb,
	Flags: []cli.Flag{
		stringFlag("port, p", "1999", "Temporary port number to prevent conflict"),
	},
}

func runWeb(c *cli.Context) error {
	db, err := db.Init()
	if err != nil {
		return errors.Wrap(err, "init database")
	}

	f := flamego.Classic()

	// We prefer to save session into database,
	// if no database configuration, the session will be saved into memory instead.
	var sessionStorage interface{}
	initer := session.MemoryIniter()
	if conf.Database.DSN != "" {
		initer = postgres.Initer()
		sessionStorage = postgres.Config{
			DSN: conf.Database.DSN,
		}
	}

	sessioner := session.Sessioner(session.Options{
		Initer: initer,
		Config: sessionStorage,
	})
	f.Use(sessioner)

	reqSignIn := context.Toggle(&context.ToggleOptions{SignInRequired: true})
	reqSignOut := context.Toggle(&context.ToggleOptions{SignOutRequired: true})

	f.Group("", func() {
		f.Group("/api", func() {
			f.Group("/user", func() {
				f.Group("", func() {
					f.Post("/sign-up", form.Bind(form.SignUp{}), route.User.SignUp)
					f.Post("/sign-in", form.Bind(form.SignIn{}), route.User.SignIn)
				}, reqSignOut)

				f.Group("", func() {
					f.Combo("/profile").Get(route.User.GetProfile).Put()
					f.Post("/sign-out", route.User.SignOut)
				}, reqSignIn)
			})

			f.Group("/doc", func() {
				f.Group("", func() {
					f.Combo("").Get(route.Document.List).Post(route.Document.New)
				}, reqSignIn)
				f.Group("/{uid}", func() {
					f.Get("/meta", route.Document.Meta)
					f.Get("/content", route.Document.Content)
					f.Post("/save", route.Document.Save)
					f.Delete("", route.Document.Delete)
					f.Combo("/setting").Get(route.Document.GetSetting).Post(form.Bind(form.UpdateDocumentSetting{}), route.Document.UpdateSetting)
				}, route.Document.Documenter)
			})

			f.Group("/s", func() {

			})

			f.Group("/config", func() {
				f.Get("/global", func(ctx context.Context) error {
					return ctx.Success("Hi EggMD!!")
				})
			})
		})
	},
		context.Contexter(db),
	)

	if c.IsSet("port") {
		conf.Server.HTTPPort = c.String("port")
	} else {
		conf.Server.HTTPPort = "1999"
	}

	listenAddr := fmt.Sprintf("%s:%s", conf.Server.HTTPAddr, conf.Server.HTTPPort)
	log.Info("Listen on http://%s", listenAddr)

	return http.ListenAndServe(listenAddr, f)
}
