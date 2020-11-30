package conf

import (
	"github.com/BurntSushi/toml"
	log "unknwon.dev/clog/v2"
)

func init() {
	_ = log.NewConsole()
}

func Init() {
	var conf struct {
		Security SecurityOpts
		Session  SessionOpts
		Server   ServerOpts
		Database DatabaseOpts
	}

	_, err := toml.DecodeFile("./conf/app.toml", &conf)
	if err != nil {
		log.Fatal("Failed to load config: %v", err)
	}

	Security = conf.Security
	Session = conf.Session
	Server = conf.Server
	Database = conf.Database
}
