package config

import (
	"os"
	"path/filepath"
)

const (
	ProjectName  = "neovector"
	appDirName   = "neostore"
	configDir    = ".config"
)

func ConfigDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, configDir, appDirName, ProjectName)
}

func EnsureConfigDir() error {
	return os.MkdirAll(ConfigDir(), 0755)
}

func ConfigFile(name string) string {
	return filepath.Join(ConfigDir(), name)
}
