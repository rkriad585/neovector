package main

import (
	_ "embed"
	"strings"

	"github.com/rkriad585/neovector/cmd"
)

//go:embed .version
var versionFile string

func init() {
	if v := strings.TrimSpace(versionFile); v != "" {
		cmd.Version = v
	}
}

func main() {
	cmd.Execute()
}
