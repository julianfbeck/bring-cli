package cmd

import (
	"github.com/julianfbeck/bring-cli/internal/config"
	"github.com/spf13/cobra"
)

var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Remove stored credentials",
	Long: `Remove stored Bring credentials from your system.

Example:
  bring logout`,
	RunE: runLogout,
}

func init() {
	rootCmd.AddCommand(logoutCmd)
}

func runLogout(cmd *cobra.Command, args []string) error {
	if err := config.ClearCredentials(); err != nil {
		return err
	}

	if isJSON() {
		return printJSON(map[string]bool{"success": true})
	}

	printSuccess("Logged out successfully")
	return nil
}
