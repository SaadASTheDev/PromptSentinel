package main

import (
	"fmt"
	"os"

	"promptsentinel/internal/cli"

	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "promptsentinel",
		Short: "PromptSentinel CLI - Validate and secure your AI prompts",
		Long: `PromptSentinel is a command-line tool for validating AI prompts for safety and security.
It helps you check if your prompts are safe for your specific use case and manage
configuration settings for prompt validation.`,
		Version: "1.0.0",
	}

	// Add subcommands
	rootCmd.AddCommand(cli.NewCheckCommand())
	rootCmd.AddCommand(cli.NewConfigCommand())
	rootCmd.AddCommand(cli.NewValidateCommand())

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
