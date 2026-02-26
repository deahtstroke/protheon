package config

import (
	"os"
	"path/filepath"
)

const (
	appName   = "protheon"
	configDir = "config"
)

func GlobalDbDir() string {
	dataHomeDir := os.Getenv("XDG_DATA_HOME")
	if dataHomeDir == "" {
		homeDir, _ := os.UserHomeDir()
		return filepath.Join(homeDir, ".local/share", appName)
	}
	return filepath.Join(dataHomeDir, appName)
}

func GlobalConfigDataDir() string {
	dataHomeDir := os.Getenv("XDG_DATA_HOME")
	if dataHomeDir == "" {
		homeDir, _ := os.UserHomeDir()
		return filepath.Join(homeDir, ".local/share", appName, configDir)
	}

	return filepath.Join(dataHomeDir, appName, configDir)
}
