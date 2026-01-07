package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	version   = "1.0.0"
	quiet     bool
	jsonOut   bool
	noColor   bool
)

var rootCmd = &cobra.Command{
	Use:   "bring",
	Short: "Manage Bring shopping lists from the command line",
	Long: `bring is a CLI tool for interacting with the Bring shopping list app.

Add items to your shopping list, check off completed items, and send
notifications to other list users - all from your terminal.

Get started:
  bring login           # Authenticate with your Bring account
  bring lists           # View all your shopping lists
  bring add Milk        # Add an item to your default list`,
	SilenceUsage:  true,
	SilenceErrors: true,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		printError(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&quiet, "quiet", "q", false, "suppress non-essential output")
	rootCmd.PersistentFlags().BoolVar(&jsonOut, "json", false, "output as JSON")
	rootCmd.PersistentFlags().BoolVar(&noColor, "no-color", false, "disable color output")

	rootCmd.Version = version
	rootCmd.SetVersionTemplate("bring version {{.Version}}\n")
}

// printError prints an error message to stderr.
func printError(err error) {
	if jsonOut {
		_ = json.NewEncoder(os.Stderr).Encode(map[string]string{"error": err.Error()})
	} else {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	}
}

// printSuccess prints a success message if not in quiet mode.
func printSuccess(format string, args ...interface{}) {
	if !quiet && !jsonOut {
		fmt.Printf(format+"\n", args...)
	}
}

// printJSON prints data as JSON.
func printJSON(v interface{}) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(v)
}

// isQuiet returns true if output should be suppressed.
func isQuiet() bool {
	return quiet
}

// isJSON returns true if JSON output is requested.
func isJSON() bool {
	return jsonOut
}
