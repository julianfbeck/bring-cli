package cmd

import (
	"fmt"

	"github.com/julianfbeck/bring-cli/internal/api"
	"github.com/julianfbeck/bring-cli/internal/config"
	"github.com/spf13/cobra"
)

var (
	notifyType string
	notifyList string
)

var notifyCmd = &cobra.Command{
	Use:   "notify",
	Short: "Send notification to list users",
	Long: `Send a notification to all users sharing the list.

Notification types:
  going-shopping  - Let others know you're heading to the store
  changed-list    - Notify that the list was updated
  shopping-done   - Let others know you finished shopping

If no list is specified, uses the default list.

Examples:
  bring notify
  bring notify --type going-shopping
  bring notify --type shopping-done --list abc123`,
	RunE: runNotify,
}

func init() {
	notifyCmd.Flags().StringVarP(&notifyType, "type", "t", "going-shopping", "notification type: going-shopping, changed-list, shopping-done")
	notifyCmd.Flags().StringVarP(&notifyList, "list", "l", "", "target list UUID")
	rootCmd.AddCommand(notifyCmd)
}

func runNotify(cmd *cobra.Command, args []string) error {
	client, err := getAuthenticatedClient()
	if err != nil {
		return err
	}

	listUUID := notifyList
	if listUUID == "" {
		creds, _ := config.GetCredentials()
		if creds != nil && creds.DefaultList != "" {
			listUUID = creds.DefaultList
		} else {
			return fmt.Errorf("no list specified and no default list configured")
		}
	}

	// Map friendly names to API constants
	var apiNotifyType string
	switch notifyType {
	case "going-shopping":
		apiNotifyType = api.NotifyGoingShopping
	case "changed-list":
		apiNotifyType = api.NotifyChangedList
	case "shopping-done":
		apiNotifyType = api.NotifyShoppingDone
	default:
		return fmt.Errorf("invalid notification type: %s (use: going-shopping, changed-list, shopping-done)", notifyType)
	}

	if err := client.Notify(listUUID, apiNotifyType, nil); err != nil {
		return fmt.Errorf("sending notification: %w", err)
	}

	if isJSON() {
		return printJSON(map[string]interface{}{
			"success": true,
			"type":    notifyType,
			"list":    listUUID,
		})
	}

	var message string
	switch notifyType {
	case "going-shopping":
		message = "Notified list users: Going shopping!"
	case "changed-list":
		message = "Notified list users: List updated"
	case "shopping-done":
		message = "Notified list users: Shopping done!"
	}
	printSuccess(message)

	return nil
}
