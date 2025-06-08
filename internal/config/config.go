package config

import (
	"errors"
	"io"
	"log"
	"os"

	"github.com/darwinfroese/scribe/internal/theme"
	"github.com/pelletier/go-toml"
)

type Config struct {
	Theme *theme.Theme
}

func Load() *Config {
	path := getConfigPath()
	config := &Config{}

	_, err := os.Stat(path)
	if !errors.Is(err, os.ErrNotExist) && err != nil {
		config.defaults()
		return config
	}

	if errors.Is(err, os.ErrNotExist) {
		config.defaults()
		return config
	}

	file, err := os.Open(path)
	if err != nil {
		config.defaults()
		return config
	}

	defer func() {
		err := file.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	contents, err := io.ReadAll(file)
	if err != nil {
		config.defaults()
		return config
	}

	config.parse(contents)
	config.Theme = theme.Load(config.Theme)

	return config
}

func (config *Config) defaults() {
	config.Theme = theme.Load(&theme.Theme{Base: "default"})
}

func (config *Config) parse(contents []byte) {
	err := toml.Unmarshal(contents, config)
	if err != nil {
		config.defaults()
	}
}
