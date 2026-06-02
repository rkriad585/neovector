package cmd

import (
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	Cyan   = color.New(color.FgCyan, color.Bold)
	Green  = color.New(color.FgGreen, color.Bold)
	Red    = color.New(color.FgRed, color.Bold)
	Yellow = color.New(color.FgYellow)
	Magenta = color.New(color.FgMagenta)

	Version = "dev"
)

var rootCmd = &cobra.Command{
	Use:     "neovector",
	Short:   "Convert between images and their numerical vector representations",
	Long: `neovector is a CLI tool for seamless conversion between raster images
and their numerical vector representations (flat lists of RGB pixel values).

It supports converting images to vectors, reconstructing images from vectors,
and inspecting dimensions of images or vector files.`,
	Version: "dev",
}

func Execute() {
	if Version != "" {
		rootCmd.Version = Version
	}
	if err := rootCmd.Execute(); err != nil {
		Red.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
