package cmd

import (
	"fmt"
	"strings"

	"github.com/julianfbeck/bring-cli/internal/api"
	"github.com/julianfbeck/bring-cli/internal/config"
	"github.com/spf13/cobra"
)

var removeList string

var removeCmd = &cobra.Command{
	Use:   "remove <item>...",
	Short: "Remove item(s) from a shopping list",
	Long: `Remove one or more items from a shopping list.

If no list is specified, uses the default list.

Examples:
  bring remove Milk
  bring remove Eggs Butter Cheese
  bring remove "Orange Juice" --list abc123`,
	Args: cobra.MinimumNArgs(1),
	RunE: runRemove,
}

func init() {
	removeCmd.Flags().StringVarP(&removeList, "list", "l", "", "target list UUID")
	rootCmd.AddCommand(removeCmd)
}

func runRemove(cmd *cobra.Command, args []string) error {
	client, err := getAuthenticatedClient()
	if err != nil {
		return err
	}

	listUUID := removeList
	if listUUID == "" {
		creds, _ := config.GetCredentials()
		if creds != nil && creds.DefaultList != "" {
			listUUID = creds.DefaultList
		} else {
			return fmt.Errorf("no list specified and no default list configured")
		}
	}

	// Build changes for all items
	var changes []api.ItemChange
	for _, item := range args {
		changes = append(changes, api.ItemChange{
			ItemID:    item,
			Operation: api.OperationRemove,
		})
	}

	if err := client.UpdateItems(listUUID, changes); err != nil {
		return fmt.Errorf("removing items: %w", err)
	}

	if isJSON() {
		return printJSON(map[string]interface{}{
			"success": true,
			"items":   args,
			"list":    listUUID,
		})
	}

	if len(args) == 1 {
		printSuccess("Removed %s from list", args[0])
	} else {
		printSuccess("Removed %d items from list: %s", len(args), strings.Join(args, ", "))
	}

	return nil
}
