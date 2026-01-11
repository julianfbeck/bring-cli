package cmd

import (
	"fmt"
	"strings"

	"github.com/julianfbeck/bring-cli/internal/config"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage bring-cli configuration",
	Long:  `Manage bring-cli configuration settings.`,
}

var setListCmd = &cobra.Command{
	Use:   "set-list <list-uuid-or-name>",
	Short: "Set the default shopping list",
	Long: `Set the default shopping list for all commands.

You can specify either the list UUID or the list name.
The default list will be used when no --list flag is provided.

Example:
  bring config set-list b63caa6a-7307-4786-9a9a-7cdc772a1763
  bring config set-list "Zuhause"`,
	Args: cobra.ExactArgs(1),
	RunE: runSetList,
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(setListCmd)
}

func runSetList(cmd *cobra.Command, args []string) error {
	listArg := args[0]

	client, err := getAuthenticatedClient()
	if err != nil {
		return err
	}

	// Fetch all lists to validate and resolve name to UUID
	lists, err := client.GetLists()
	if err != nil {
		return fmt.Errorf("fetching lists: %w", err)
	}

	var listUUID string
	var listName string

	// Check if argument is a UUID or a name
	for _, list := range lists.Lists {
		if list.ListUUID == listArg {
			listUUID = list.ListUUID
			listName = list.Name
			break
		}
		if strings.EqualFold(list.Name, listArg) {
			listUUID = list.ListUUID
			listName = list.Name
			break
		}
	}

	if listUUID == "" {
		return fmt.Errorf("list not found: %s\nRun 'bring lists' to see available lists", listArg)
	}

	if err := config.SetDefaultList(listUUID); err != nil {
		return fmt.Errorf("saving default list: %w", err)
	}

	if !isQuiet() {
		fmt.Printf("Default list set to: %s (%s)\n", listName, listUUID)
	}

	return nil
}
