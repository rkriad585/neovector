package main

import (
	_ "embed"
	"strings"

	"github.com/rkriad585/neovector/cmd"
	"github.com/rkriad585/neovector/internal/version"
)

//go:embed .version
var versionFile string

func init() {
	if version.Version == "dev" {
		if v := strings.TrimSpace(versionFile); v != "" {
			version.Version = v
		}
	}
}

func main() {
	cmd.Execute()
}
