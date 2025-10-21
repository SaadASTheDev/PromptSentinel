package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"promptsentinel/internal/validator"

	"github.com/spf13/cobra"
)

// NewCheckCommand creates the check command for validating prompts
func NewCheckCommand() *cobra.Command {
	var configFile string
	var useCase string

	cmd := &cobra.Command{
		Use:   "check [prompt]",
		Short: "Check if a prompt is safe for your use case",
		Long: `Check validates a prompt against safety and security criteria.
You can provide the prompt as an argument or via stdin.

Examples:
  promptsentinel check "Write a story about a cat"
  echo "Your prompt here" | promptsentinel check
  promptsentinel check "Your prompt" --config ./config.json`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var prompt string
			var err error

			if len(args) > 0 {
				prompt = args[0]
			} else {
				// Read from stdin
				prompt, err = readFromStdin()
				if err != nil {
					return fmt.Errorf("failed to read from stdin: %w", err)
				}
			}

			if strings.TrimSpace(prompt) == "" {
				return fmt.Errorf("prompt cannot be empty")
			}

			// Load configuration
			config, err := loadConfig(configFile)
			if err != nil {
				return fmt.Errorf("failed to load config: %w", err)
			}

			// Override use case if provided
			if useCase != "" {
				config.UseCase = useCase
			}

			// Validate the prompt
			result, err := validator.ValidatePrompt(prompt, config)
			if err != nil {
				return fmt.Errorf("validation failed: %w", err)
			}

			// Display results
			displayResults(result)
			return nil
		},
	}

	cmd.Flags().StringVarP(&configFile, "config", "c", "", "Path to configuration file")
	cmd.Flags().StringVarP(&useCase, "use-case", "u", "", "Override the use case for validation")

	return cmd
}

// NewConfigCommand creates the config command for managing settings
func NewConfigCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Manage PromptSentinel configuration",
		Long:  "Manage configuration settings for prompt validation.",
	}

	cmd.AddCommand(newConfigInitCommand())
	cmd.AddCommand(newConfigShowCommand())
	cmd.AddCommand(newConfigSetCommand())

	return cmd
}

// NewValidateCommand creates the validate command for comprehensive validation
func NewValidateCommand() *cobra.Command {
	var configFile string
	var outputFormat string

	cmd := &cobra.Command{
		Use:   "validate [prompt]",
		Short: "Perform comprehensive prompt validation",
		Long: `Validate performs a comprehensive analysis of a prompt including
safety checks, security analysis, and compliance validation.

Examples:
  promptsentinel validate "Your prompt here"
  promptsentinel validate "Your prompt" --format json`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var prompt string
			var err error

			if len(args) > 0 {
				prompt = args[0]
			} else {
				prompt, err = readFromStdin()
				if err != nil {
					return fmt.Errorf("failed to read from stdin: %w", err)
				}
			}

			if strings.TrimSpace(prompt) == "" {
				return fmt.Errorf("prompt cannot be empty")
			}

			// Load configuration
			config, err := loadConfig(configFile)
			if err != nil {
				return fmt.Errorf("failed to load config: %w", err)
			}

			// Perform comprehensive validation
			result, err := validator.ValidatePromptComprehensive(prompt, config)
			if err != nil {
				return fmt.Errorf("validation failed: %w", err)
			}

			// Display results in requested format
			if outputFormat == "json" {
				displayJSONResults(result)
			} else {
				displayDetailedResults(result)
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&configFile, "config", "c", "", "Path to configuration file")
	cmd.Flags().StringVarP(&outputFormat, "format", "f", "text", "Output format (text, json)")

	return cmd
}

func newConfigInitCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "init",
		Short: "Initialize a new configuration file",
		Long:  "Create a new configuration file with default settings.",
		RunE: func(cmd *cobra.Command, args []string) error {
			configPath := getDefaultConfigPath()

			// Check if config already exists
			if _, err := os.Stat(configPath); err == nil {
				return fmt.Errorf("configuration file already exists at %s", configPath)
			}

			// Create default configuration
			config := validator.DefaultConfig()

			// Ensure config directory exists
			if err := os.MkdirAll(filepath.Dir(configPath), 0755); err != nil {
				return fmt.Errorf("failed to create config directory: %w", err)
			}

			// Write configuration file
			if err := saveConfig(config, configPath); err != nil {
				return fmt.Errorf("failed to save config: %w", err)
			}

			fmt.Printf("Configuration initialized at: %s\n", configPath)
			return nil
		},
	}
}

func newConfigShowCommand() *cobra.Command {
	var configFile string

	cmd := &cobra.Command{
		Use:   "show",
		Short: "Show current configuration",
		Long:  "Display the current configuration settings.",
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := loadConfig(configFile)
			if err != nil {
				return fmt.Errorf("failed to load config: %w", err)
			}

			displayConfig(config)
			return nil
		},
	}

	cmd.Flags().StringVarP(&configFile, "config", "c", "", "Path to configuration file")

	return cmd
}

func newConfigSetCommand() *cobra.Command {
	var configFile string
	var key string
	var value string

	cmd := &cobra.Command{
		Use:   "set <key> <value>",
		Short: "Set a configuration value",
		Long:  "Set a specific configuration value.",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			key = args[0]
			value = args[1]

			config, err := loadConfig(configFile)
			if err != nil {
				return fmt.Errorf("failed to load config: %w", err)
			}

			// Update configuration
			if err := setConfigValue(config, key, value); err != nil {
				return fmt.Errorf("failed to set config value: %w", err)
			}

			// Save updated configuration
			configPath := getConfigPath(configFile)
			if err := saveConfig(config, configPath); err != nil {
				return fmt.Errorf("failed to save config: %w", err)
			}

			fmt.Printf("Set %s = %s\n", key, value)
			return nil
		},
	}

	cmd.Flags().StringVarP(&configFile, "config", "c", "", "Path to configuration file")

	return cmd
}
