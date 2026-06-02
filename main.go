package main

import (
	_ "embed"
	"strings"

	"github.com/rkriad585/neovector/cmd"
	"github.com/rkriad585/neovector/internal/version"
)

//go:embed .version
var versionFile string

// ldflags injection targets (set by release workflow via -X main.Version=...)
var (
	Version        string
	Commit         string
	PublisherName  string
	PublisherEmail string
)

func init() {
	if Version != "" {
		version.Version = Version
	}
	if Commit != "" {
		version.Commit = Commit
	}
	if PublisherName != "" {
		version.PublisherName = PublisherName
	}
	if PublisherEmail != "" {
		version.PublisherEmail = PublisherEmail
	}
	if version.Version == "dev" {
		if v := strings.TrimSpace(versionFile); v != "" {
			version.Version = v
		}
	}
}

func main() {
	cmd.Execute()
}
