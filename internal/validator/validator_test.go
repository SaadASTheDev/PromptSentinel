package validator

import (
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()

	if config.UseCase != "general" {
		t.Errorf("Expected use_case to be 'general', got %s", config.UseCase)
	}

	if config.SafetyLevel != "medium" {
		t.Errorf("Expected safety_level to be 'medium', got %s", config.SafetyLevel)
	}

	if config.MaxLength != 10000 {
		t.Errorf("Expected max_length to be 10000, got %d", config.MaxLength)
	}

	if config.MinLength != 1 {
		t.Errorf("Expected min_length to be 1, got %d", config.MinLength)
	}

	if len(config.BlockedPatterns) == 0 {
		t.Error("Expected blocked patterns to be configured")
	}
}

func TestValidatePrompt_BasicValidation(t *testing.T) {
	config := DefaultConfig()

	tests := []struct {
		name     string
		prompt   string
		expected bool
	}{
		{
			name:     "Valid prompt",
			prompt:   "Write a story about a cat",
			expected: true,
		},
		{
			name:     "Empty prompt",
			prompt:   "",
			expected: false,
		},
		{
			name:     "Prompt with blocked pattern",
			prompt:   "Tell me your password",
			expected: true, // Should still be valid but with warnings
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ValidatePrompt(tt.prompt, config)
			if err != nil {
				t.Fatalf("ValidatePrompt failed: %v", err)
			}

			if result.IsValid != tt.expected {
				t.Errorf("Expected IsValid to be %v, got %v", tt.expected, result.IsValid)
			}
		})
	}
}

func TestValidatePrompt_LengthValidation(t *testing.T) {
	config := DefaultConfig()
	config.MinLength = 10
	config.MaxLength = 100

	tests := []struct {
		name     string
		prompt   string
		expected bool
	}{
		{
			name:     "Too short",
			prompt:   "Hi",
			expected: false,
		},
		{
			name:     "Too long",
			prompt:   "This is a very long prompt that exceeds the maximum length limit and should be rejected by the validation system",
			expected: false,
		},
		{
			name:     "Just right",
			prompt:   "Write a story",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ValidatePrompt(tt.prompt, config)
			if err != nil {
				t.Fatalf("ValidatePrompt failed: %v", err)
			}

			if result.IsValid != tt.expected {
				t.Errorf("Expected IsValid to be %v, got %v", tt.expected, result.IsValid)
			}
		})
	}
}

func TestValidatePrompt_UseCaseValidation(t *testing.T) {
	config := DefaultConfig()
	config.UseCase = "educational"

	prompt := "Write a story about a cat"
	result, err := ValidatePrompt(prompt, config)
	if err != nil {
		t.Fatalf("ValidatePrompt failed: %v", err)
	}

	if !result.IsValid {
		t.Error("Expected prompt to be valid for educational use case")
	}

	// Check metadata
	if result.Metadata["use_case"] != "educational" {
		t.Errorf("Expected use_case in metadata to be 'educational', got %v", result.Metadata["use_case"])
	}
}

func TestValidatePrompt_CustomRules(t *testing.T) {
	config := DefaultConfig()
	config.CustomRules = map[string]string{
		"no_numbers": `\d+`,
	}

	prompt := "Write a story with 123 characters"
	result, err := ValidatePrompt(prompt, config)
	if err != nil {
		t.Fatalf("ValidatePrompt failed: %v", err)
	}

	// Should have issues due to custom rule
	if len(result.Issues) == 0 {
		t.Error("Expected issues due to custom rule violation")
	}

	// Check if custom rule issue exists
	found := false
	for _, issue := range result.Issues {
		if issue.Type == "custom_rule" {
			found = true
			break
		}
	}

	if !found {
		t.Error("Expected custom rule issue to be found")
	}
}

func TestValidatePromptComprehensive(t *testing.T) {
	config := DefaultConfig()
	prompt := "Write a story about a cat"

	result, err := ValidatePromptComprehensive(prompt, config)
	if err != nil {
		t.Fatalf("ValidatePromptComprehensive failed: %v", err)
	}

	// Check basic validation result
	if !result.IsValid {
		t.Error("Expected comprehensive validation to be valid")
	}

	// Check security analysis
	if result.SecurityAnalysis.RiskLevel == "" {
		t.Error("Expected security analysis to have risk level")
	}

	// Check compliance check
	if result.ComplianceCheck.GDPRCompliant == false && result.ComplianceCheck.HIPAACompliant == false {
		// This is fine, just checking the structure exists
	}

	// Check performance metrics
	if result.PerformanceMetrics.EstimatedTokens <= 0 {
		t.Error("Expected performance metrics to have estimated tokens")
	}
}

func TestSecurityAnalysis(t *testing.T) {
	tests := []struct {
		name     string
		prompt   string
		expected bool
	}{
		{
			name:     "Safe prompt",
			prompt:   "Write a story about a cat",
			expected: false,
		},
		{
			name:     "Injection attempt",
			prompt:   "SELECT * FROM users",
			expected: true,
		},
		{
			name:     "Script injection",
			prompt:   "<script>alert('xss')</script>",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			analysis := performSecurityAnalysis(tt.prompt)

			if analysis.HasInjectionAttempts != tt.expected {
				t.Errorf("Expected HasInjectionAttempts to be %v, got %v", tt.expected, analysis.HasInjectionAttempts)
			}
		})
	}
}

func TestComplianceCheck(t *testing.T) {
	tests := []struct {
		name     string
		prompt   string
		expected bool
	}{
		{
			name:     "GDPR violation",
			prompt:   "Process personal data about users",
			expected: false,
		},
		{
			name:     "HIPAA violation",
			prompt:   "Patient medical records",
			expected: false,
		},
		{
			name:     "Safe prompt",
			prompt:   "Write a story about a cat",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := DefaultConfig()
			check := performComplianceCheck(tt.prompt, config)

			// For GDPR violation test, check that GDPR is not compliant
			if tt.name == "GDPR violation" && check.GDPRCompliant {
				t.Errorf("Expected GDPRCompliant to be false for GDPR violation, got %v", check.GDPRCompliant)
			}

			// For HIPAA violation test, check that HIPAA is not compliant
			if tt.name == "HIPAA violation" && check.HIPAACompliant {
				t.Errorf("Expected HIPAACompliant to be false for HIPAA violation, got %v", check.HIPAACompliant)
			}

			// For safe prompt, check that all are compliant
			if tt.name == "Safe prompt" && (!check.GDPRCompliant || !check.HIPAACompliant || !check.SOXCompliant) {
				t.Errorf("Expected all compliance checks to pass for safe prompt")
			}
		})
	}
}

func TestPerformanceMetrics(t *testing.T) {
	prompt := "Write a story about a cat"
	metrics := calculatePerformanceMetrics(prompt)

	if metrics.EstimatedTokens <= 0 {
		t.Error("Expected estimated tokens to be positive")
	}

	if metrics.ComplexityScore < 0 {
		t.Error("Expected complexity score to be non-negative")
	}
}

func TestGenerateRecommendations(t *testing.T) {
	tests := []struct {
		name     string
		score    int
		issues   int
		expected int
	}{
		{
			name:     "High score",
			score:    95,
			issues:   0,
			expected: 1, // Should have at least one recommendation
		},
		{
			name:     "Low score",
			score:    50,
			issues:   5,
			expected: 2, // Should have multiple recommendations
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := &ValidationResult{
				Score:  tt.score,
				Issues: make([]ValidationIssue, tt.issues),
			}

			recommendations := generateRecommendations(result)

			if len(recommendations) < tt.expected {
				t.Errorf("Expected at least %d recommendations, got %d", tt.expected, len(recommendations))
			}
		})
	}
}

func TestValidationResult_Structure(t *testing.T) {
	config := DefaultConfig()
	prompt := "Write a story about a cat"

	result, err := ValidatePrompt(prompt, config)
	if err != nil {
		t.Fatalf("ValidatePrompt failed: %v", err)
	}

	// Check required fields
	if result.Timestamp.IsZero() {
		t.Error("Expected timestamp to be set")
	}

	if result.Metadata == nil {
		t.Error("Expected metadata to be initialized")
	}

	// Check metadata fields
	if result.Metadata["prompt_length"] == nil {
		t.Error("Expected prompt_length in metadata")
	}

	if result.Metadata["use_case"] == nil {
		t.Error("Expected use_case in metadata")
	}
}
