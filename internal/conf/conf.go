package conf

import (
	"github.com/BurntSushi/toml"
)

// Init 解析配置文件并初始化配置信息。
func Init(configPath string) error {
	var conf struct {
		Security SecurityOpts
		Session  SessionOpts
		Server   ServerOpts
		Database DatabaseOpts
	}

	_, err := toml.DecodeFile(configPath, &conf)
	if err != nil {
		return err
	}

	Security = conf.Security
	Session = conf.Session
	Server = conf.Server
	Database = conf.Database

	return nil
}
