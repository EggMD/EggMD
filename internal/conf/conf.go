package conf

import (
	"github.com/pelletier/go-toml"
	"github.com/pkg/errors"
)

func Init(configPath string) error {
	config, err := toml.LoadFile(configPath)
	if err != nil {
		return errors.Wrap(err, "load toml config file")
	}
	return parse(config)
}

func parse(config *toml.Tree) error {
	if err := config.Get("Server").(*toml.Tree).Unmarshal(&Server); err != nil {
		return errors.Wrap(err, "mapping [App] section")
	}

	if err := config.Get("Database").(*toml.Tree).Unmarshal(&Database); err != nil {
		return errors.Wrap(err, "mapping [Database] section")
	}

	return nil
}
