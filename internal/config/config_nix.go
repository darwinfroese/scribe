//go:build linux || darwin
// +build linux darwin

package config

import (
	"os/user"
	"path/filepath"
)

const (
	configPath = ".config/scribe/scribe.yml"
)

func getConfigPath() string {
	usr, _ := user.Current()
	homeDir := usr.HomeDir

	path := filepath.Join(homeDir, configPath)

	return path
}
