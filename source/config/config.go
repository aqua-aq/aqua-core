package config

import (
	"io"

	"github.com/pelletier/go-toml/v2"
	"github.com/vandi37/aqua/env"
)

type Config struct {
	Name         string            `toml:"name"`
	Version      string            `toml:"version"`
	Main         string            `toml:"main"`
	Lib          *struct{}         `toml:"lib"`
	Dependencies map[string]string `toml:"dependencies"`
}

var defaultConfig = Config{
	Name:    "<unknown project>",
	Version: "0.1.0",
	Main:    env.MAIN,
}

func DefaultConfig() Config {
	return defaultConfig
}

func NewConfig(r io.Reader) (Config, error) {
	config := DefaultConfig()
	err := toml.NewDecoder(r).Decode(&config)
	if err != nil {
		return Config{}, err
	}
	return config, nil
}
