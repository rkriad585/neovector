package config

import (
	"os"
	"path/filepath"
	"sync"

	"github.com/BurntSushi/toml"
)

const (
	ProjectName = "neovector"
	appDirName  = "neostore"
	configDir   = ".config"
)

type Config struct {
	General GeneralConfig `toml:"general"`
	Network NetworkConfig `toml:"network"`
	Theme   ThemeConfig   `toml:"theme"`
}

type GeneralConfig struct {
	DefaultFormat string `toml:"default_format"`
}

type NetworkConfig struct {
	Proxy string `toml:"proxy"`
}

type ThemeConfig struct {
	Name string `toml:"name"`
	Mode string `toml:"mode"`
}

var (
	global *Config
	mu     sync.RWMutex
)

func DefaultConfig() *Config {
	return &Config{
		General: GeneralConfig{
			DefaultFormat: "txt",
		},
		Network: NetworkConfig{
			Proxy: "",
		},
		Theme: ThemeConfig{
			Name: "sunny_beach_day",
			Mode: "dark",
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

func Get() *Config {
	mu.RLock()
	if global != nil {
		mu.RUnlock()
		return global
	}
	mu.RUnlock()

	mu.Lock()
	defer mu.Unlock()
	if global != nil {
		return global
	}

	cfg, err := Load()
	if err != nil {
		cfg = DefaultConfig()
	}
	global = cfg
	return cfg
}

func Load() (*Config, error) {
	path := ConfigPath()
	if _, err := os.Stat(path); os.IsNotExist(err) {
		cfg := DefaultConfig()
		if err := Save(cfg); err != nil {
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

func Save(cfg *Config) error {
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

// Set saves cfg to disk and replaces the in-memory global.
func Set(cfg *Config) error {
	if err := Save(cfg); err != nil {
		return err
	}
	mu.Lock()
	global = cfg
	mu.Unlock()
	return nil
}

// LoadConfig is the legacy loader kept for backward compatibility.
// New code should use Get() or Load() instead.
func LoadConfig() (*Config, error) {
	return Load()
}

// SaveConfig is the legacy saver kept for backward compatibility.
// New code should use Save() or Set() instead.
func SaveConfig(cfg *Config) error {
	return Save(cfg)
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
