package cli

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
)

var (
	// Global flags
	secretKey    string
	outputFormat string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "authkit",
	Short: "AuthKit CLI - Authentication library for Go",
	Long: `AuthKit CLI provides command-line tools for managing users, tokens, 
and authentication operations using the AuthKit library.

Examples:
  authkit user register --email user@example.com --password pass123 --name "John Doe"
  authkit token generate --user-id user123 --secret mySecret
  authkit token validate --token "eyJhbGc..." --secret mySecret`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func init() {
	// Global flags
	rootCmd.PersistentFlags().StringVarP(&secretKey, "secret", "s", "", "JWT secret key (required)")
	rootCmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", "json", "Output format (json, table, yaml)")

	// Mark secret as required for commands that need it
	rootCmd.MarkPersistentFlagRequired("secret")
}

// Common helper functions
func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func printJSON(data interface{}) {
	// Implementation for JSON output
	fmt.Printf("JSON output: %+v\n", data)
}

func printTable(data interface{}) {
	// Implementation for table output
	fmt.Printf("Table output: %+v\n", data)
}

func printOutput(data interface{}) {
	switch outputFormat {
	case "json":
		printJSON(data)
	case "table":
		printTable(data)
	default:
		fmt.Printf("%+v\n", data)
	}
}
