//go:build windows
// +build windows

package config

import (
	"os"
	"path/filepath"
)

const (
	configPath = ".config/scribe.yml"
)

func getConfigPath() string {
	home := os.Getenv("APPDATA")
	path := filepath.Join(home, configPath)

	return path
}
