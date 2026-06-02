package cmd

import (
	"errors"
	"fmt"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"

	"github.com/rkriad585/neovector/internal/config"
	"github.com/rkriad585/neovector/internal/theme"
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

var configThemeCmd = &cobra.Command{
	Use:   "theme [name]",
	Short: "View or change the color theme",
	Long: `Display the current theme or switch to a different one.

With no argument, prints the current theme.
Provide a theme name to switch directly.
Use --edit to open the interactive TUI picker.

Available themes:
  dark                    Dark Theme
  light                   Light Theme
  sunny_beach_day         Sunny Beach Day (Default)
  olive_garden_feast      Olive Garden Feast
  summer_ocean_breeze     Summer Ocean Breeze
  refreshing_summer_fun   Refreshing Summer Fun
  black_gold_elegance     Black & Gold Elegance
  vibrant_color_fiesta    Vibrant Color Fiesta
  light_steel             Light Steel
  golden_twilight         Golden Twilight
  deep_sea                Deep Sea
  bright_green            Bright Green
  vivid_nightfall         Vivid Nightfall`,
	Args: cobra.MaximumNArgs(1),
	RunE: configThemeRun,
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.Flags().Bool("edit", false, "Edit configuration interactively")
	configCmd.Flags().String("format", "", "Set default output format (txt or json)")
	configCmd.Flags().String("proxy", "", "Set proxy URL for self-update")
	configCmd.AddCommand(configThemeCmd)
	configThemeCmd.Flags().Bool("edit", false, "Choose a theme interactively")
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
	Cyan.Println("  [theme]")
	fmt.Printf("    name = %q\n", cfg.Theme.Name)
	fmt.Printf("    mode = %q\n", cfg.Theme.Mode)
	fmt.Println()
	Yellow.Println("  Run 'neovector config --edit' or 'neovector config theme --edit' to modify.")
	return nil
}

func configEditForm() error {
	cfg := config.Get()

	originalFormat := cfg.General.DefaultFormat
	originalProxy := cfg.Network.Proxy
	originalTheme := cfg.Theme.Name

	themeNames := theme.Names()
	themeLabels := theme.Labels()
	themeOpts := make([]huh.Option[string], len(themeNames))
	for i := range themeNames {
		themeOpts[i] = huh.NewOption(themeLabels[i], themeNames[i])
	}

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

			huh.NewSelect[string]().
				Title("Color Theme").
				Description("Choose a color theme for the terminal UI").
				Options(themeOpts...).
				Value(&cfg.Theme.Name),
		),
	)

	if err := form.Run(); err != nil {
		if errors.Is(err, huh.ErrUserAborted) {
			Yellow.Println("  Canceled.")
			return nil
		}
		return fmt.Errorf("failed to open config editor: %w", err)
	}

	if cfg.General.DefaultFormat == originalFormat &&
		cfg.Network.Proxy == originalProxy &&
		cfg.Theme.Name == originalTheme {
		Yellow.Println("  No changes made.")
		return nil
	}

	if err := config.Set(cfg); err != nil {
		return fmt.Errorf("failed to save configuration: %w", err)
	}

	applyTheme(cfg.Theme.Name)
	Green.Println("  Configuration saved.")
	return nil
}

func configThemeRun(cmd *cobra.Command, args []string) error {
	edit, _ := cmd.Flags().GetBool("edit")

	if edit {
		return configThemeEditForm()
	}

	if len(args) == 1 {
		name := args[0]
		if _, ok := theme.Find(name); !ok {
			return fmt.Errorf("unknown theme %q\nRun 'neovector config theme' to see available themes.", name)
		}
		cfg := config.Get()
		cfg.Theme.Name = name
		if err := config.Set(cfg); err != nil {
			return fmt.Errorf("failed to save theme: %w", err)
		}
		applyTheme(name)
		Green.Printf("  Theme set to %q.\n", name)
		return nil
	}

	cfg := config.Get()
	current := cfg.Theme.Name
	Yellow.Println("━━ Current Theme ━━")
	fmt.Println()
	Cyan.Printf("  Name: %s\n", current)
	if t, ok := theme.Find(current); ok {
		for i, role := range []string{"Primary", "Success", "Warning", "Error", "Accent"} {
			hex := t.Hex(theme.Role(i))
			fmt.Printf("    %-8s %s\n", role, theme.Colorize("████", hex))
		}
	}
	fmt.Println()
	Green.Println("  Available themes:")
	for _, t := range theme.Themes {
		mark := " "
		if t.Name == current {
			mark = "▸"
		}
		palette := ""
		for _, hex := range t.Colors {
			palette += theme.Colorize("██", hex)
		}
		fmt.Printf("  %s %-25s %s\n", mark, t.Label, palette)
	}
	fmt.Println()
	Yellow.Println("  Switch: neovector config theme <name>")
	Yellow.Println("  Browse: neovector config theme --edit")
	return nil
}

func configThemeEditForm() error {
	cfg := config.Get()
	original := cfg.Theme.Name

	themeOpts := make([]huh.Option[string], len(theme.Themes))
	for i, t := range theme.Themes {
		label := t.Label
		if t.Name == original {
			label += " (current)"
		}
		palette := ""
		for _, hex := range t.Colors {
			palette += theme.Colorize("██", hex)
		}
		themeOpts[i] = huh.NewOption(label+"  "+palette, t.Name)
	}

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Choose a Theme").
				Description("The new theme is applied immediately after saving").
				Options(themeOpts...).
				Value(&cfg.Theme.Name),
		),
	)

	if err := form.Run(); err != nil {
		if errors.Is(err, huh.ErrUserAborted) {
			Yellow.Println("  Canceled.")
			return nil
		}
		return fmt.Errorf("failed to open theme picker: %w", err)
	}

	if cfg.Theme.Name == original {
		Yellow.Println("  No changes made.")
		return nil
	}

	if err := config.Set(cfg); err != nil {
		return fmt.Errorf("failed to save theme: %w", err)
	}

	applyTheme(cfg.Theme.Name)
	Green.Printf("  Theme set to %q.\n", cfg.Theme.Name)
	return nil
}
