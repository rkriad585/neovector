package cmd

import (
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/rkriad585/neovector/internal/banner"
	"github.com/rkriad585/neovector/internal/config"
	"github.com/rkriad585/neovector/internal/theme"
	"github.com/rkriad585/neovector/internal/uninstall"
	"github.com/rkriad585/neovector/internal/version"
)

var (
	Cyan    *color.Color
	Green   *color.Color
	Red     *color.Color
	Yellow  *color.Color
	Magenta *color.Color
)

func init() {
	applyTheme("sunny_beach_day")
}

func applyTheme(name string) {
	c := theme.Resolve(name)
	Cyan = c.Primary
	Green = c.Success
	Yellow = c.Warning
	Red = c.Error
	Magenta = c.Accent
}

var rootCmd = &cobra.Command{
	Use:     "neovector",
	Short:   "Convert between images and their numerical vector representations",
	Long: `neovector is a CLI tool for seamless conversion between raster images
and their numerical vector representations (flat lists of RGB pixel values).

It supports converting images to vectors, reconstructing images from vectors,
and inspecting dimensions of images or vector files.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return initSystem()
	},
}

func initSystem() error {
	cfg := config.Get()
	if cfg.Theme.Name != "" {
		applyTheme(cfg.Theme.Name)
	}
	return nil
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version and build information",
	RunE: func(cmd *cobra.Command, args []string) error {
		Cyan.Printf("  neovector %s (commit: %s)\n", version.Version, version.Commit)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

func Execute() {
	for _, arg := range os.Args[1:] {
		if arg == "--selfuninstall" || arg == "--uninstall" || arg == "--self-uninstall" {
			if isTopLevelUninstall() {
				os.Exit(uninstall.Run())
			}
		}
		if arg == "-u" {
			if isTopLevelUninstall() && !hasHelpFlags() {
				os.Exit(uninstall.Run())
			}
		}
	}

	banner.Print()

	rootCmd.Version = version.Version
	if err := rootCmd.Execute(); err != nil {
		Red.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

// isTopLevelUninstall ensures uninstall flags are not consumed by a subcommand.
func isTopLevelUninstall() bool {
	for _, arg := range os.Args[1:] {
		if !strings.HasPrefix(arg, "-") {
			return false
		}
	}
	return true
}

// hasHelpFlags checks if --help or --version is present alongside -u.
func hasHelpFlags() bool {
	for _, arg := range os.Args[1:] {
		if arg == "--help" || arg == "-h" || arg == "--version" || arg == "-v" {
			return true
		}
	}
	return false
}
