package cmd

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/rkriad585/neovector/internal/uninstall"
)

var selfUninstallCmd = &cobra.Command{
	Use:   "self-uninstall",
	Short: "Uninstall neovector from the system",
	Long: `Removes neovector binaries and configuration files from the system.

On Windows, a deferred batch script handles deleting the running binary
since Windows cannot delete in-use executables.

After uninstalling, remove the neovector PATH entry from your shell rc
or run the installer with --selfuninstall to clean PATH automatically.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		os.Exit(uninstall.Run())
		return nil
	},
}

func init() {
	rootCmd.AddCommand(selfUninstallCmd)
}
