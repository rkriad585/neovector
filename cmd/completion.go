package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var completionCmd = &cobra.Command{
	Use:   "completion [bash|zsh|fish|powershell]",
	Short: "Generate shell completion scripts",
	Long: `Generate shell completion scripts for bash, zsh, fish, or PowerShell.

To load completions:

  Bash:
    source <(neovector completion bash)
    # or add to ~/.bashrc:
    echo 'source <(neovector completion bash)' >> ~/.bashrc

  Zsh:
    source <(neovector completion zsh)
    # or add to ~/.zshrc:
    echo 'source <(neovector completion zsh)' >> ~/.zshrc

  Fish:
    neovector completion fish | source
    # or persist:
    neovector completion fish > ~/.config/fish/completions/neovector.fish

  PowerShell:
    neovector completion powershell | Out-String | Invoke-Expression
    # or add to profile:
    neovector completion powershell >> $PROFILE`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

var completionBashCmd = &cobra.Command{
	Use:   "bash",
	Short: "Generate bash completion script",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Root().GenBashCompletion(os.Stdout)
	},
}

var completionZshCmd = &cobra.Command{
	Use:   "zsh",
	Short: "Generate zsh completion script",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Root().GenZshCompletion(os.Stdout)
	},
}

var completionFishCmd = &cobra.Command{
	Use:   "fish",
	Short: "Generate fish completion script",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Root().GenFishCompletion(os.Stdout, true)
	},
}

var completionPowerShellCmd = &cobra.Command{
	Use:   "powershell",
	Short: "Generate PowerShell completion script",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Root().GenPowerShellCompletionWithDesc(os.Stdout)
	},
}

func init() {
	rootCmd.AddCommand(completionCmd)
	completionCmd.AddCommand(completionBashCmd)
	completionCmd.AddCommand(completionZshCmd)
	completionCmd.AddCommand(completionFishCmd)
	completionCmd.AddCommand(completionPowerShellCmd)

	completionCmd.Aliases = []string{"completions"}

	completionCmd.RunE = func(cmd *cobra.Command, args []string) error {
		if len(args) == 1 {
			switch args[0] {
			case "bash":
				return completionBashCmd.RunE(cmd, args)
			case "zsh":
				return completionZshCmd.RunE(cmd, args)
			case "fish":
				return completionFishCmd.RunE(cmd, args)
			case "powershell":
				return completionPowerShellCmd.RunE(cmd, args)
			default:
				return fmt.Errorf("unknown shell %q: supported shells are bash, zsh, fish, powershell", args[0])
			}
		}
		return cmd.Help()
	}
}
