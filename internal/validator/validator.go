package validator

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

// Config represents the configuration for prompt validation
type Config struct {
	UseCase         string            `json:"use_case"`
	SafetyLevel     string            `json:"safety_level"`
	AllowedDomains  []string          `json:"allowed_domains"`
	BlockedPatterns []string          `json:"blocked_patterns"`
	MaxLength       int               `json:"max_length"`
	MinLength       int               `json:"min_length"`
	RequireApproval bool              `json:"require_approval"`
	CustomRules     map[string]string `json:"custom_rules"`
	LastUpdated     time.Time         `json:"last_updated"`
}

// ValidationResult represents the result of prompt validation
type ValidationResult struct {
	IsValid         bool                   `json:"is_valid"`
	Score           int                    `json:"score"`
	Issues          []ValidationIssue      `json:"issues"`
	Recommendations []string               `json:"recommendations"`
	Metadata        map[string]interface{} `json:"metadata"`
	Timestamp       time.Time              `json:"timestamp"`
}

// ValidationIssue represents a specific validation issue
type ValidationIssue struct {
	Type       string `json:"type"`
	Severity   string `json:"severity"`
	Message    string `json:"message"`
	Suggestion string `json:"suggestion,omitempty"`
	Line       int    `json:"line,omitempty"`
	Column     int    `json:"column,omitempty"`
}

// ComprehensiveValidationResult extends ValidationResult with additional analysis
type ComprehensiveValidationResult struct {
	ValidationResult
	SecurityAnalysis   SecurityAnalysis   `json:"security_analysis"`
	ComplianceCheck    ComplianceCheck    `json:"compliance_check"`
	PerformanceMetrics PerformanceMetrics `json:"performance_metrics"`
}

// SecurityAnalysis contains security-related validation results
type SecurityAnalysis struct {
	HasInjectionAttempts bool     `json:"has_injection_attempts"`
	HasSensitiveData     bool     `json:"has_sensitive_data"`
	RiskLevel            string   `json:"risk_level"`
	Threats              []string `json:"threats"`
}

// ComplianceCheck contains compliance-related validation results
type ComplianceCheck struct {
	GDPRCompliant    bool     `json:"gdpr_compliant"`
	HIPAACompliant   bool     `json:"hipaa_compliant"`
	SOXCompliant     bool     `json:"sox_compliant"`
	ComplianceIssues []string `json:"compliance_issues"`
}

// PerformanceMetrics contains performance-related analysis
type PerformanceMetrics struct {
	EstimatedTokens   int     `json:"estimated_tokens"`
	ComplexityScore   float64 `json:"complexity_score"`
	ProcessingTime    int64   `json:"processing_time_ms"`
	ResourceIntensive bool    `json:"resource_intensive"`
}

// DefaultConfig returns a default configuration
func DefaultConfig() *Config {
	return &Config{
		UseCase:        "general",
		SafetyLevel:    "medium",
		AllowedDomains: []string{},
		BlockedPatterns: []string{
			`(?i)(password|secret|key|token|credential)`,
			`(?i)(inject|exploit|hack|attack)`,
			`(?i)(illegal|unlawful|criminal)`,
			`(?i)(violence|harm|kill|destroy)`,
		},
		MaxLength:       10000,
		MinLength:       1,
		RequireApproval: false,
		CustomRules:     make(map[string]string),
		LastUpdated:     time.Now(),
	}
}

// ValidatePrompt performs basic prompt validation
func ValidatePrompt(prompt string, config *Config) (*ValidationResult, error) {
	startTime := time.Now()
	result := &ValidationResult{
		IsValid:   true,
		Score:     100,
		Issues:    []ValidationIssue{},
		Metadata:  make(map[string]interface{}),
		Timestamp: time.Now(),
	}

	// Length validation
	if len(prompt) < config.MinLength {
		result.Issues = append(result.Issues, ValidationIssue{
			Type:       "length",
			Severity:   "error",
			Message:    fmt.Sprintf("Prompt too short (minimum %d characters)", config.MinLength),
			Suggestion: "Add more content to your prompt",
		})
		result.IsValid = false
		result.Score -= 20
	}

	if len(prompt) > config.MaxLength {
		result.Issues = append(result.Issues, ValidationIssue{
			Type:       "length",
			Severity:   "error",
			Message:    fmt.Sprintf("Prompt too long (maximum %d characters)", config.MaxLength),
			Suggestion: "Shorten your prompt",
		})
		result.IsValid = false
		result.Score -= 20
	}

	// Pattern validation
	for _, pattern := range config.BlockedPatterns {
		matched, err := regexp.MatchString(pattern, prompt)
		if err != nil {
			continue // Skip invalid patterns
		}
		if matched {
			result.Issues = append(result.Issues, ValidationIssue{
				Type:       "pattern",
				Severity:   "warning",
				Message:    fmt.Sprintf("Prompt contains blocked pattern: %s", pattern),
				Suggestion: "Review and modify the flagged content",
			})
			result.Score -= 10
		}
	}

	// Use case specific validation
	if err := validateUseCase(prompt, config.UseCase, result); err != nil {
		return nil, err
	}

	// Custom rules validation
	if err := validateCustomRules(prompt, config.CustomRules, result); err != nil {
		return nil, err
	}

	// Generate recommendations
	result.Recommendations = generateRecommendations(result)

	// Add metadata
	result.Metadata["processing_time_ms"] = time.Since(startTime).Milliseconds()
	result.Metadata["prompt_length"] = len(prompt)
	result.Metadata["use_case"] = config.UseCase

	return result, nil
}

// ValidatePromptComprehensive performs comprehensive prompt validation
func ValidatePromptComprehensive(prompt string, config *Config) (*ComprehensiveValidationResult, error) {
	// Perform basic validation first
	basicResult, err := ValidatePrompt(prompt, config)
	if err != nil {
		return nil, err
	}

	comprehensiveResult := &ComprehensiveValidationResult{
		ValidationResult:   *basicResult,
		SecurityAnalysis:   performSecurityAnalysis(prompt),
		ComplianceCheck:    performComplianceCheck(prompt, config),
		PerformanceMetrics: calculatePerformanceMetrics(prompt),
	}

	return comprehensiveResult, nil
}

// validateUseCase performs use case specific validation
func validateUseCase(prompt string, useCase string, result *ValidationResult) error {
	switch useCase {
	case "educational":
		// Educational prompts should be informative and safe
		if strings.Contains(strings.ToLower(prompt), "harmful") {
			result.Issues = append(result.Issues, ValidationIssue{
				Type:       "use_case",
				Severity:   "warning",
				Message:    "Educational prompts should avoid potentially harmful content",
				Suggestion: "Reframe the prompt to focus on learning objectives",
			})
			result.Score -= 5
		}
	case "business":
		// Business prompts should be professional
		if strings.Contains(strings.ToLower(prompt), "personal") {
			result.Issues = append(result.Issues, ValidationIssue{
				Type:       "use_case",
				Severity:   "info",
				Message:    "Business prompts should focus on professional objectives",
				Suggestion: "Consider removing personal references",
			})
		}
	case "creative":
		// Creative prompts have more flexibility
		// No specific restrictions for creative use cases
	}
	return nil
}

// validateCustomRules validates against custom rules
func validateCustomRules(prompt string, customRules map[string]string, result *ValidationResult) error {
	for ruleName, pattern := range customRules {
		matched, err := regexp.MatchString(pattern, prompt)
		if err != nil {
			continue // Skip invalid patterns
		}
		if matched {
			result.Issues = append(result.Issues, ValidationIssue{
				Type:       "custom_rule",
				Severity:   "warning",
				Message:    fmt.Sprintf("Custom rule '%s' triggered", ruleName),
				Suggestion: "Review the custom rule configuration",
			})
			result.Score -= 5
		}
	}
	return nil
}

// performSecurityAnalysis performs security analysis on the prompt
func performSecurityAnalysis(prompt string) SecurityAnalysis {
	analysis := SecurityAnalysis{
		HasInjectionAttempts: false,
		HasSensitiveData:     false,
		RiskLevel:            "low",
		Threats:              []string{},
	}

	// Check for injection attempts
	injectionPatterns := []string{
		`(?i)(union|select|insert|update|delete|drop)`,
		`(?i)(<script|javascript:|onload=)`,
		`(?i)(exec|eval|system|shell)`,
	}

	for _, pattern := range injectionPatterns {
		if matched, _ := regexp.MatchString(pattern, prompt); matched {
			analysis.HasInjectionAttempts = true
			analysis.Threats = append(analysis.Threats, "Potential injection attempt")
			analysis.RiskLevel = "high"
		}
	}

	// Check for sensitive data
	sensitivePatterns := []string{
		`(?i)(ssn|social security)`,
		`(?i)(credit card|card number)`,
		`(?i)(password|passwd)`,
		`(?i)(api key|secret key)`,
	}

	for _, pattern := range sensitivePatterns {
		if matched, _ := regexp.MatchString(pattern, prompt); matched {
			analysis.HasSensitiveData = true
			analysis.Threats = append(analysis.Threats, "Potential sensitive data exposure")
			if analysis.RiskLevel == "low" {
				analysis.RiskLevel = "medium"
			}
		}
	}

	return analysis
}

// performComplianceCheck performs compliance validation
func performComplianceCheck(prompt string, config *Config) ComplianceCheck {
	check := ComplianceCheck{
		GDPRCompliant:    true,
		HIPAACompliant:   true,
		SOXCompliant:     true,
		ComplianceIssues: []string{},
	}

	// GDPR compliance check
	if strings.Contains(strings.ToLower(prompt), "personal data") {
		check.GDPRCompliant = false
		check.ComplianceIssues = append(check.ComplianceIssues, "Potential GDPR violation: personal data handling")
	}

	// HIPAA compliance check
	hipaaPatterns := []string{
		`(?i)(patient|medical|health|diagnosis)`,
		`(?i)(phi|protected health information)`,
	}

	for _, pattern := range hipaaPatterns {
		if matched, _ := regexp.MatchString(pattern, prompt); matched {
			check.HIPAACompliant = false
			check.ComplianceIssues = append(check.ComplianceIssues, "Potential HIPAA violation: health information")
			break
		}
	}

	// SOX compliance check
	if strings.Contains(strings.ToLower(prompt), "financial") {
		check.SOXCompliant = false
		check.ComplianceIssues = append(check.ComplianceIssues, "Potential SOX violation: financial information")
	}

	return check
}

// calculatePerformanceMetrics calculates performance-related metrics
func calculatePerformanceMetrics(prompt string) PerformanceMetrics {
	// Simple token estimation (rough approximation)
	words := strings.Fields(prompt)
	estimatedTokens := int(float64(len(words)) * 1.3) // Rough approximation

	// Calculate complexity score based on various factors
	complexityScore := 0.0
	complexityScore += float64(len(strings.Fields(prompt))) * 0.01
	complexityScore += float64(len(strings.Split(prompt, "\n"))) * 0.1
	complexityScore += float64(len(strings.Split(prompt, "."))) * 0.05

	// Check if resource intensive
	resourceIntensive := len(prompt) > 5000 || estimatedTokens > 2000

	return PerformanceMetrics{
		EstimatedTokens:   int(estimatedTokens),
		ComplexityScore:   complexityScore,
		ProcessingTime:    time.Now().UnixNano() / 1000000, // Mock processing time
		ResourceIntensive: resourceIntensive,
	}
}

// generateRecommendations generates recommendations based on validation results
func generateRecommendations(result *ValidationResult) []string {
	recommendations := []string{}

	if result.Score < 80 {
		recommendations = append(recommendations, "Consider reviewing flagged content for safety")
	}

	if len(result.Issues) > 3 {
		recommendations = append(recommendations, "Multiple issues detected - consider breaking down the prompt")
	}

	if result.Score > 90 {
		recommendations = append(recommendations, "Prompt looks good! Consider adding more specific instructions for better results")
	}

	return recommendations
}
