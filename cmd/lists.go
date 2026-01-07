package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/julianfbeck/bring-cli/internal/api"
	"github.com/julianfbeck/bring-cli/internal/config"
	"github.com/spf13/cobra"
)

var listsCmd = &cobra.Command{
	Use:   "lists",
	Short: "List all shopping lists",
	Long: `Display all shopping lists associated with your Bring account.

Example:
  bring lists
  bring lists --json`,
	RunE: runLists,
}

func init() {
	rootCmd.AddCommand(listsCmd)
}

func runLists(cmd *cobra.Command, args []string) error {
	client, err := getAuthenticatedClient()
	if err != nil {
		return err
	}

	lists, err := client.GetLists()
	if err != nil {
		return fmt.Errorf("fetching lists: %w", err)
	}

	if isJSON() {
		return printJSON(lists)
	}

	if len(lists.Lists) == 0 {
		fmt.Println("No shopping lists found")
		return nil
	}

	creds, _ := config.GetCredentials()
	defaultList := ""
	if creds != nil {
		defaultList = creds.DefaultList
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "NAME\tUUID\tDEFAULT")
	for _, list := range lists.Lists {
		isDefault := ""
		if list.ListUUID == defaultList {
			isDefault = "*"
		}
		fmt.Fprintf(w, "%s\t%s\t%s\n", list.Name, list.ListUUID, isDefault)
	}
	w.Flush()

	return nil
}

// getAuthenticatedClient returns an authenticated API client.
func getAuthenticatedClient() (*api.Client, error) {
	creds, err := config.GetCredentials()
	if err != nil {
		return nil, fmt.Errorf("loading credentials: %w", err)
	}
	if creds == nil || creds.AccessToken == "" {
		return nil, fmt.Errorf("not logged in. Run 'bring login' first")
	}

	client := api.NewClient(creds, config.SaveCredentials)
	return client, nil
}
