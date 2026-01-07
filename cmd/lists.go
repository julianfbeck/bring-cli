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
// It first checks for BRING_EMAIL and BRING_PASSWORD environment variables.
// If not set, it falls back to stored credentials from config file.
func getAuthenticatedClient() (*api.Client, error) {
	// Check for environment variables first
	email := os.Getenv("BRING_EMAIL")
	password := os.Getenv("BRING_PASSWORD")

	if email != "" && password != "" {
		// Login with environment variables
		client := api.NewClient(nil, nil)
		_, err := client.Login(email, password)
		if err != nil {
			return nil, fmt.Errorf("login with environment variables failed: %w", err)
		}
		return client, nil
	}

	// Fall back to stored credentials
	creds, err := config.GetCredentials()
	if err != nil {
		return nil, fmt.Errorf("loading credentials: %w", err)
	}
	if creds == nil || creds.AccessToken == "" {
		return nil, fmt.Errorf("not logged in. Set BRING_EMAIL and BRING_PASSWORD environment variables, or run 'bring login'")
	}

	client := api.NewClient(creds, config.SaveCredentials)
	return client, nil
}

// getDefaultListUUID returns the list UUID to use.
// Priority: flag > BRING_LIST env var > stored default list
func getDefaultListUUID(flagValue string) (string, error) {
	// 1. Command line flag has highest priority
	if flagValue != "" {
		return flagValue, nil
	}

	// 2. Check BRING_LIST environment variable
	if envList := os.Getenv("BRING_LIST"); envList != "" {
		return envList, nil
	}

	// 3. Fall back to stored default list
	creds, _ := config.GetCredentials()
	if creds != nil && creds.DefaultList != "" {
		return creds.DefaultList, nil
	}

	return "", fmt.Errorf("no list specified. Use --list flag, set BRING_LIST environment variable, or run 'bring login'")
}
