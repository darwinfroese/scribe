package config

import (
	"errors"
	"os"

	"github.com/darwinfroese/scribe/internal/theme"
)

type Config struct {
	Theme *theme.Theme
}

func Load() *Config {
	path := getConfigPath()

	_, err := os.Stat(path)
	if !errors.Is(err, os.ErrNotExist) && err != nil {
		return defaults()
	}

	if errors.Is(err, os.ErrNotExist) {
		return defaults()
	}

	return &Config{}
}

func defaults() *Config {
	return &Config{
		Theme: theme.Load("default"),
	}
}
