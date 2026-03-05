package config

import (
	"os"
	"path/filepath"
)

const (
	appName   = "protheon"
	dbName    = "protheon.db"
	configDir = "config"
)

func GlobalDatabaseUrl() string {
	dataHomeDir := os.Getenv("XDG_DATA_HOME")
	if dataHomeDir == "" {
		homeDir, _ := os.UserHomeDir()
		path := filepath.Join(homeDir, ".local/share", appName, dbName)
		return "file:" + path
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
