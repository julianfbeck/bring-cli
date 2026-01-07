package cmd

import (
	"fmt"
	"strings"

	"github.com/julianfbeck/bring-cli/internal/api"
	"github.com/spf13/cobra"
)

var completeList string

var completeCmd = &cobra.Command{
	Use:   "complete <item>...",
	Short: "Mark item(s) as completed",
	Long: `Mark one or more items as completed (moves to recently bought).

If no list is specified, uses the default list.

Examples:
  bring complete Milk
  bring complete Eggs Butter Cheese
  bring complete "Orange Juice" --list abc123`,
	Args: cobra.MinimumNArgs(1),
	RunE: runComplete,
}

func init() {
	completeCmd.Flags().StringVarP(&completeList, "list", "l", "", "target list UUID")
	rootCmd.AddCommand(completeCmd)
}

func runComplete(cmd *cobra.Command, args []string) error {
	client, err := getAuthenticatedClient()
	if err != nil {
		return err
	}

	listUUID, err := getDefaultListUUID(completeList)
	if err != nil {
		return err
	}

	// Build changes for all items
	var changes []api.ItemChange
	for _, item := range args {
		changes = append(changes, api.ItemChange{
			ItemID:    item,
			Operation: api.OperationComplete,
		})
	}

	if err := client.UpdateItems(listUUID, changes); err != nil {
		return fmt.Errorf("completing items: %w", err)
	}

	if isJSON() {
		return printJSON(map[string]interface{}{
			"success": true,
			"items":   args,
			"list":    listUUID,
		})
	}

	if len(args) == 1 {
		printSuccess("Completed %s", args[0])
	} else {
		printSuccess("Completed %d items: %s", len(args), strings.Join(args, ", "))
	}

	return nil
}
