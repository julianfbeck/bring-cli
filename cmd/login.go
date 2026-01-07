package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/julianfbeck/bring-cli/internal/api"
	"github.com/julianfbeck/bring-cli/internal/config"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Authenticate with your Bring account",
	Long: `Authenticate with your Bring account using email and password.

Credentials are stored securely in ~/.config/bring-cli/config.yaml.

Example:
  bring login`,
	RunE: runLogin,
}

func init() {
	rootCmd.AddCommand(loginCmd)
}

func runLogin(cmd *cobra.Command, args []string) error {
	// Check if already logged in
	creds, _ := config.GetCredentials()
	if creds != nil && creds.AccessToken != "" {
		fmt.Println("You are already logged in as", creds.Email)
		fmt.Print("Do you want to log in with a different account? [y/N]: ")
		reader := bufio.NewReader(os.Stdin)
		response, _ := reader.ReadString('\n')
		response = strings.TrimSpace(strings.ToLower(response))
		if response != "y" && response != "yes" {
			return nil
		}
	}

	// Prompt for email
	fmt.Print("Email: ")
	reader := bufio.NewReader(os.Stdin)
	email, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("reading email: %w", err)
	}
	email = strings.TrimSpace(email)

	if email == "" {
		return fmt.Errorf("email is required")
	}

	// Prompt for password (hidden input)
	fmt.Print("Password: ")
	passwordBytes, err := term.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		return fmt.Errorf("reading password: %w", err)
	}
	fmt.Println() // New line after hidden input
	password := string(passwordBytes)

	if password == "" {
		return fmt.Errorf("password is required")
	}

	// Authenticate
	client := api.NewClient(nil, nil)
	authResp, err := client.Login(email, password)
	if err != nil {
		return fmt.Errorf("login failed: %w", err)
	}

	// Save credentials
	if err := config.SaveCredentials(client.GetCredentials()); err != nil {
		return fmt.Errorf("saving credentials: %w", err)
	}

	if isJSON() {
		return printJSON(map[string]interface{}{
			"success":      true,
			"email":        authResp.Email,
			"default_list": authResp.BringListUUID,
		})
	}

	printSuccess("Logged in as %s", authResp.Email)
	printSuccess("Default list: %s", authResp.BringListUUID)

	return nil
}
