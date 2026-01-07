package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list [list-uuid]",
	Short: "Show items in a shopping list",
	Long: `Display all items in a shopping list.

If no list UUID is provided, uses the default list.

Examples:
  bring list
  bring list abc123-def456
  bring list --json`,
	RunE: runList,
}

func init() {
	rootCmd.AddCommand(listCmd)
}

func runList(cmd *cobra.Command, args []string) error {
	client, err := getAuthenticatedClient()
	if err != nil {
		return err
	}

	argListUUID := ""
	if len(args) > 0 {
		argListUUID = args[0]
	}
	listUUID, err := getDefaultListUUID(argListUUID)
	if err != nil {
		return err
	}

	items, err := client.GetListItems(listUUID)
	if err != nil {
		return fmt.Errorf("fetching list items: %w", err)
	}

	if isJSON() {
		return printJSON(items)
	}

	if len(items.Items.Purchase) == 0 && len(items.Items.Recently) == 0 {
		fmt.Println("List is empty")
		return nil
	}

	if len(items.Items.Purchase) > 0 {
		fmt.Println("To Buy:")
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "  ITEM\tSPECIFICATION")
		for _, item := range items.Items.Purchase {
			fmt.Fprintf(w, "  %s\t%s\n", item.ItemID, item.Specification)
		}
		w.Flush()
	}

	if len(items.Items.Recently) > 0 {
		if len(items.Items.Purchase) > 0 {
			fmt.Println()
		}
		fmt.Println("Recently Completed:")
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "  ITEM\tSPECIFICATION")
		for _, item := range items.Items.Recently {
			fmt.Fprintf(w, "  %s\t%s\n", item.ItemID, item.Specification)
		}
		w.Flush()
	}

	return nil
}
