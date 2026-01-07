package cmd

import (
	"fmt"
	"strings"

	"github.com/julianfbeck/bring-cli/internal/api"
	"github.com/spf13/cobra"
)

var (
	addSpec string
	addList string
)

var addCmd = &cobra.Command{
	Use:   "add <item>...",
	Short: "Add item(s) to a shopping list",
	Long: `Add one or more items to a shopping list.

If no list is specified, uses the default list.

Examples:
  bring add Milk
  bring add Bread --spec "2 loaves, whole wheat"
  bring add Eggs Butter Cheese
  bring add "Orange Juice" --list abc123`,
	Args: cobra.MinimumNArgs(1),
	RunE: runAdd,
}

func init() {
	addCmd.Flags().StringVarP(&addSpec, "spec", "s", "", "item specification (quantity, notes)")
	addCmd.Flags().StringVarP(&addList, "list", "l", "", "target list UUID")
	rootCmd.AddCommand(addCmd)
}

func runAdd(cmd *cobra.Command, args []string) error {
	client, err := getAuthenticatedClient()
	if err != nil {
		return err
	}

	listUUID, err := getDefaultListUUID(addList)
	if err != nil {
		return err
	}

	// Build changes for all items
	var changes []api.ItemChange
	for _, item := range args {
		changes = append(changes, api.ItemChange{
			ItemID:    item,
			Spec:      addSpec,
			Operation: api.OperationAdd,
		})
	}

	if err := client.UpdateItems(listUUID, changes); err != nil {
		return fmt.Errorf("adding items: %w", err)
	}

	if isJSON() {
		return printJSON(map[string]interface{}{
			"success": true,
			"items":   args,
			"list":    listUUID,
		})
	}

	if len(args) == 1 {
		printSuccess("Added %s to list", args[0])
	} else {
		printSuccess("Added %d items to list: %s", len(args), strings.Join(args, ", "))
	}

	return nil
}
