package config

import (
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

const (
	ProjectName = "neovector"
	appDirName  = "neostore"
	configDir   = ".config"
)

type Config struct {
	General GeneralConfig `toml:"general"`
}

type GeneralConfig struct {
	DefaultFormat string `toml:"default_format"`
}

func DefaultConfig() *Config {
	return &Config{
		General: GeneralConfig{
			DefaultFormat: "txt",
		},
	}
}

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

func ConfigPath() string {
	return ConfigFile("config.toml")
}

func LoadConfig() (*Config, error) {
	path := ConfigPath()
	if _, err := os.Stat(path); os.IsNotExist(err) {
		cfg := DefaultConfig()
		if err := SaveConfig(cfg); err != nil {
			return nil, err
		}
		return cfg, nil
	}

	var cfg Config
	if _, err := toml.DecodeFile(path, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func SaveConfig(cfg *Config) error {
	path := ConfigPath()
	if err := EnsureConfigDir(); err != nil {
		return err
	}
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	return toml.NewEncoder(file).Encode(cfg)
}

func LogPath() string {
	return ConfigFile("history.log")
}

func OutputDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, "Downloads", appDirName, ProjectName)
}

func EnsureOutputDir() error {
	return os.MkdirAll(OutputDir(), 0755)
}

func OutputFile(name string) string {
	return filepath.Join(OutputDir(), name)
}

func ResolveOutput(path string) string {
	if filepath.IsAbs(path) {
		return path
	}
	if len(path) > 1 && (path[0] == '.' || path[0] == '~') {
		return path
	}
	if filepath.Base(path) != path {
		return path
	}
	return OutputFile(path)
}

func ResolveInput(path string) string {
	if _, err := os.Stat(path); err == nil {
		return path
	}
	if fallback := OutputFile(path); true {
		if _, err := os.Stat(fallback); err == nil {
			return fallback
		}
	}
	return path
}
