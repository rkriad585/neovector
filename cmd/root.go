package cmd

import (
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/rkriad585/neovector/internal/banner"
	"github.com/rkriad585/neovector/internal/config"
	"github.com/rkriad585/neovector/internal/uninstall"
	"github.com/rkriad585/neovector/internal/version"
)

var (
	Cyan   = color.New(color.FgCyan, color.Bold)
	Green  = color.New(color.FgGreen, color.Bold)
	Red    = color.New(color.FgRed, color.Bold)
	Yellow = color.New(color.FgYellow)
	Magenta = color.New(color.FgMagenta)

	appCfg *config.Config
)

var rootCmd = &cobra.Command{
	Use:     "neovector",
	Short:   "Convert between images and their numerical vector representations",
	Long: `neovector is a CLI tool for seamless conversion between raster images
and their numerical vector representations (flat lists of RGB pixel values).

It supports converting images to vectors, reconstructing images from vectors,
and inspecting dimensions of images or vector files.`,
	Version: version.Version,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return initSystem()
	},
}

func initSystem() error {
	if err := config.EnsureConfigDir(); err != nil {
		return err
	}
	var err error
	appCfg, err = config.LoadConfig()
	return err
}

func Execute() {
	banner.Print()

	for _, arg := range os.Args[1:] {
		if arg == "--selfuninstall" || arg == "--uninstall" || arg == "-u" {
			if arg != "-u" || isTopLevelUninstall() {
				os.Exit(uninstall.Run())
			}
		}
	}

	rootCmd.Version = version.Version
	if err := rootCmd.Execute(); err != nil {
		Red.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

// isTopLevelUninstall ensures -u is not consumed by a subcommand.
func isTopLevelUninstall() bool {
	for _, arg := range os.Args[1:] {
		if !strings.HasPrefix(arg, "-") {
			return false
		}
	}
	return true
}
