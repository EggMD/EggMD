package main

import (
	"os"

	"github.com/urfave/cli"
	log "unknwon.dev/clog/v2"

	"github.com/EggMD/EggMD/internal/cmd"
)

var (
	version = ""
)

func main() {
	defer log.Stop()
	err := log.NewConsole()
	if err != nil {
		panic(err)
	}

	app := cli.NewApp()
	app.Name = "EggMD"
	app.Usage = "Self-hosted collaborative documents service"
	app.Version = version
	app.Commands = []cli.Command{
		cmd.Web,
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal("Failed to start application: %v", err)
	}
}
