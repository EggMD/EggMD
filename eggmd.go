package main

import (
	"os"

	"github.com/urfave/cli"
	log "unknwon.dev/clog/v2"

	"github.com/EggMD/EggMD/internal/cmd"
	"github.com/EggMD/EggMD/internal/conf"
	"github.com/EggMD/EggMD/internal/db"
)

var (
	Version = "development"
)

func main() {
	defer log.Stop()
	err := log.NewConsole()
	if err != nil {
		panic(err)
	}

	err = conf.Init("./conf/app.toml")
	if err != nil {
		log.Fatal("Config error: %v", err)
	}
	conf.Server.AppVersion = Version

	err = db.Init()
	if err != nil {
		log.Fatal("Database error: %v", err)
	}

	app := cli.NewApp()
	app.Name = "EggMD"
	app.Usage = "Self-hosted collaborative documents service"
	app.Version = Version
	app.Commands = []cli.Command{
		cmd.Web,
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal("Failed to start application: %v", err)
	}
}
