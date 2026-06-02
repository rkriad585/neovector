package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestConfigDir(t *testing.T) {
	dir := ConfigDir()
	if dir == "" {
		t.Fatal("ConfigDir() returned empty string")
	}

	home, err := os.UserHomeDir()
	if err != nil {
		t.Fatal(err)
	}

	expected := filepath.Join(home, ".config", "neostore", "neovector")
	if dir != expected {
		t.Errorf("ConfigDir() = %q, want %q", dir, expected)
	}

	t.Logf("Config dir: %s", dir)
}

func TestEnsureConfigDir(t *testing.T) {
	if err := EnsureConfigDir(); err != nil {
		t.Fatal(err)
	}

	dir := ConfigDir()
	info, err := os.Stat(dir)
	if err != nil {
		t.Fatal(err)
	}

	if !info.IsDir() {
		t.Fatal("ConfigDir() is not a directory")
	}

	os.Remove(dir)
}

func TestConfigFile(t *testing.T) {
	name := "settings.json"
	path := ConfigFile(name)
	expected := filepath.Join(ConfigDir(), name)
	if path != expected {
		t.Errorf("ConfigFile(%q) = %q, want %q", name, path, expected)
	}
}

func TestProjectName(t *testing.T) {
	if ProjectName != "neovector" {
		t.Errorf("ProjectName = %q, want %q", ProjectName, "neovector")
	}
}
