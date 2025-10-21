package cli

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"promptsentinel/internal/validator"
)

// readFromStdin reads input from stdin
func readFromStdin() (string, error) {
	scanner := bufio.NewScanner(os.Stdin)
	var lines []string

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	return strings.Join(lines, "\n"), nil
}

// getDefaultConfigPath returns the default configuration file path
func getDefaultConfigPath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		// Fallback to current directory
		return "./promptsentinel.json"
	}
	return filepath.Join(homeDir, ".config", "promptsentinel", "config.json")
}

// getConfigPath returns the configuration file path
func getConfigPath(configFile string) string {
	if configFile != "" {
		return configFile
	}
	return getDefaultConfigPath()
}

// loadConfig loads configuration from file
func loadConfig(configFile string) (*validator.Config, error) {
	configPath := getConfigPath(configFile)

	// Check if config file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// Return default config if file doesn't exist
		return validator.DefaultConfig(), nil
	}

	// Read and parse config file
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config validator.Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return &config, nil
}

// saveConfig saves configuration to file
func saveConfig(config *validator.Config, configPath string) error {
	// Ensure directory exists
	if err := os.MkdirAll(filepath.Dir(configPath), 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Marshal config to JSON
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	// Write to file
	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// setConfigValue sets a configuration value
func setConfigValue(config *validator.Config, key, value string) error {
	switch key {
	case "use_case":
		config.UseCase = value
	case "safety_level":
		config.SafetyLevel = value
	case "max_length":
		// Parse integer value
		var maxLength int
		if _, err := fmt.Sscanf(value, "%d", &maxLength); err != nil {
			return fmt.Errorf("invalid max_length value: %w", err)
		}
		config.MaxLength = maxLength
	case "min_length":
		// Parse integer value
		var minLength int
		if _, err := fmt.Sscanf(value, "%d", &minLength); err != nil {
			return fmt.Errorf("invalid min_length value: %w", err)
		}
		config.MinLength = minLength
	case "require_approval":
		// Parse boolean value
		config.RequireApproval = strings.ToLower(value) == "true"
	default:
		return fmt.Errorf("unknown configuration key: %s", key)
	}

	return nil
}

// displayResults displays validation results
func displayResults(result *validator.ValidationResult) {
	fmt.Printf("\n🔍 Prompt Validation Results\n")
	fmt.Printf("============================\n\n")

	// Status
	status := "❌ FAILED"
	if result.IsValid {
		status = "✅ PASSED"
	}
	fmt.Printf("Status: %s\n", status)
	fmt.Printf("Score: %d/100\n\n", result.Score)

	// Issues
	if len(result.Issues) > 0 {
		fmt.Printf("Issues Found:\n")
		for i, issue := range result.Issues {
			severityIcon := "⚠️"
			if issue.Severity == "error" {
				severityIcon = "❌"
			} else if issue.Severity == "info" {
				severityIcon = "ℹ️"
			}

			fmt.Printf("  %d. %s [%s] %s\n", i+1, severityIcon, strings.ToUpper(issue.Severity), issue.Message)
			if issue.Suggestion != "" {
				fmt.Printf("     💡 Suggestion: %s\n", issue.Suggestion)
			}
		}
		fmt.Println()
	}

	// Recommendations
	if len(result.Recommendations) > 0 {
		fmt.Printf("Recommendations:\n")
		for i, rec := range result.Recommendations {
			fmt.Printf("  %d. %s\n", i+1, rec)
		}
		fmt.Println()
	}

	// Metadata
	if len(result.Metadata) > 0 {
		fmt.Printf("Metadata:\n")
		for key, value := range result.Metadata {
			fmt.Printf("  %s: %v\n", key, value)
		}
	}
}

// displayDetailedResults displays comprehensive validation results
func displayDetailedResults(result *validator.ComprehensiveValidationResult) {
	// Display basic results first
	displayResults(&result.ValidationResult)

	// Security Analysis
	fmt.Printf("\n🔒 Security Analysis\n")
	fmt.Printf("===================\n")
	fmt.Printf("Risk Level: %s\n", strings.ToUpper(result.SecurityAnalysis.RiskLevel))
	fmt.Printf("Injection Attempts: %t\n", result.SecurityAnalysis.HasInjectionAttempts)
	fmt.Printf("Sensitive Data: %t\n", result.SecurityAnalysis.HasSensitiveData)

	if len(result.SecurityAnalysis.Threats) > 0 {
		fmt.Printf("Threats Detected:\n")
		for i, threat := range result.SecurityAnalysis.Threats {
			fmt.Printf("  %d. %s\n", i+1, threat)
		}
	}

	// Compliance Check
	fmt.Printf("\n📋 Compliance Check\n")
	fmt.Printf("==================\n")
	fmt.Printf("GDPR Compliant: %t\n", result.ComplianceCheck.GDPRCompliant)
	fmt.Printf("HIPAA Compliant: %t\n", result.ComplianceCheck.HIPAACompliant)
	fmt.Printf("SOX Compliant: %t\n", result.ComplianceCheck.SOXCompliant)

	if len(result.ComplianceCheck.ComplianceIssues) > 0 {
		fmt.Printf("Compliance Issues:\n")
		for i, issue := range result.ComplianceCheck.ComplianceIssues {
			fmt.Printf("  %d. %s\n", i+1, issue)
		}
	}

	// Performance Metrics
	fmt.Printf("\n⚡ Performance Metrics\n")
	fmt.Printf("=====================\n")
	fmt.Printf("Estimated Tokens: %d\n", result.PerformanceMetrics.EstimatedTokens)
	fmt.Printf("Complexity Score: %.2f\n", result.PerformanceMetrics.ComplexityScore)
	fmt.Printf("Resource Intensive: %t\n", result.PerformanceMetrics.ResourceIntensive)
}

// displayJSONResults displays results in JSON format
func displayJSONResults(result *validator.ComprehensiveValidationResult) {
	data, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		fmt.Printf("Error marshaling results: %v\n", err)
		return
	}
	fmt.Println(string(data))
}

// displayConfig displays configuration
func displayConfig(config *validator.Config) {
	fmt.Printf("\n⚙️  PromptSentinel Configuration\n")
	fmt.Printf("================================\n\n")

	fmt.Printf("Use Case: %s\n", config.UseCase)
	fmt.Printf("Safety Level: %s\n", config.SafetyLevel)
	fmt.Printf("Max Length: %d\n", config.MaxLength)
	fmt.Printf("Min Length: %d\n", config.MinLength)
	fmt.Printf("Require Approval: %t\n", config.RequireApproval)

	if len(config.AllowedDomains) > 0 {
		fmt.Printf("Allowed Domains: %s\n", strings.Join(config.AllowedDomains, ", "))
	}

	if len(config.BlockedPatterns) > 0 {
		fmt.Printf("Blocked Patterns: %d patterns configured\n", len(config.BlockedPatterns))
	}

	if len(config.CustomRules) > 0 {
		fmt.Printf("Custom Rules: %d rules configured\n", len(config.CustomRules))
	}

	fmt.Printf("Last Updated: %s\n", config.LastUpdated.Format("2006-01-02 15:04:05"))
}
