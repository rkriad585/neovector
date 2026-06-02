package cmd

import (
	"fmt"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"

	"github.com/rkriad585/neovector/internal/config"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "View or edit neovector configuration",
	Long: `Display the current configuration or modify it interactively.

Without flags, prints the current configuration values.
Use --edit to open the interactive TUI form.
Use --format or --proxy to change a single value directly.`,
	RunE: configRun,
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.Flags().Bool("edit", false, "Edit configuration interactively")
	configCmd.Flags().String("format", "", "Set default output format (txt or json)")
	configCmd.Flags().String("proxy", "", "Set proxy URL for self-update")
}

func configRun(cmd *cobra.Command, args []string) error {
	edit, _ := cmd.Flags().GetBool("edit")
	format, _ := cmd.Flags().GetString("format")
	proxy, _ := cmd.Flags().GetString("proxy")

	formatChanged := cmd.Flags().Changed("format")
	proxyChanged := cmd.Flags().Changed("proxy")

	if edit {
		return configEditForm()
	}

	if formatChanged || proxyChanged {
		cfg := config.Get()

		if formatChanged {
			if format != "txt" && format != "json" {
				return fmt.Errorf("invalid format %q: must be 'txt' or 'json'", format)
			}
			cfg.General.DefaultFormat = format
		}
		if proxyChanged {
			cfg.Network.Proxy = proxy
		}

		if err := config.Set(cfg); err != nil {
			return fmt.Errorf("failed to save configuration: %w", err)
		}
		Green.Println("  Configuration updated.")
		return nil
	}

	cfg := config.Get()
	Yellow.Println("━━ neovector Configuration ━━")
	fmt.Println()
	Cyan.Println("  Config file:")
	fmt.Printf("    %s\n", config.ConfigPath())
	fmt.Println()
	Cyan.Println("  [general]")
	fmt.Printf("    default_format = %q\n", cfg.General.DefaultFormat)
	fmt.Println()
	Cyan.Println("  [network]")
	fmt.Printf("    proxy = %q\n", cfg.Network.Proxy)
	fmt.Println()
	Yellow.Println("  Run 'neovector config --edit' to modify interactively.")
	return nil
}

func configEditForm() error {
	cfg := config.Get()

	originalFormat := cfg.General.DefaultFormat
	originalProxy := cfg.Network.Proxy

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Default Output Format").
				Description("Output format for image-to-vector conversion").
				Options(
					huh.NewOption("Text (.txt)", "txt"),
					huh.NewOption("JSON (.json)", "json"),
				).
				Value(&cfg.General.DefaultFormat),

			huh.NewInput().
				Title("Proxy URL").
				Description("Optional HTTP proxy for self-update (leave empty for none)").
				Placeholder("http://proxy:8080").
				Value(&cfg.Network.Proxy),
		),
	)

	if err := form.Run(); err != nil {
		return fmt.Errorf("failed to open config editor: %w", err)
	}

	if cfg.General.DefaultFormat == originalFormat && cfg.Network.Proxy == originalProxy {
		Yellow.Println("  No changes made.")
		return nil
	}

	if err := config.Set(cfg); err != nil {
		return fmt.Errorf("failed to save configuration: %w", err)
	}

	Green.Println("  Configuration saved.")
	return nil
}
