//go:build windows
// +build windows

package config

import (
	"os"
	"path/filepath"
)

const (
	configPath = ".config/scribe/scribe.toml"
)

func getConfigPath() string {
	home := os.Getenv("APPDATA")
	path := filepath.Join(home, configPath)

	return path
}
