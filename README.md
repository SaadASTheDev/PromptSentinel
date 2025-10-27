# PromptSentinel - Basic Prompt Checker

**Course**: CISC 4900 - Software Engineering  
**Student**: Saad  
**Instructor**: Katherine Chuang
**Semester**: Fall 2025 
**Project Type**: Final Project  

## Project Overview

This is a simple command-line tool I built for my CISC 4900 final project. It checks if text prompts are "safe" by looking for certain words and patterns. It's pretty basic but shows I understand Go programming, testing, and documentation.

## What It Does

- **Basic Text Checking** – Looks for bad words and suspicious patterns
- **Simple Configuration** – You can change what words to block
- **Command Line Tool** – Type commands to check text
- **Basic Testing** – Has some tests to make sure it works

## How to Run It

1. Make sure you have Go installed
2. Clone this repo: `git clone [your-repo-url]`
3. Go to the folder: `cd PromptSentinel`
4. Build it: `go build ./cmd/promptsentinel`
5. Run it: `./promptsentinel check "your text here"`

## Basic Usage

```bash
# Check some text
./promptsentinel check "Write a story about a cat"

# Check text from a file
echo "Some text" | ./promptsentinel check

# See all commands
./promptsentinel --help
```

## Usage

### Basic Commands

#### Check Command
Validate a prompt for safety and security:
```bash
# Check a prompt directly
promptsentinel check "Your prompt here"

# Check from stdin
echo "Your prompt" | promptsentinel check

# Use custom configuration
promptsentinel check "Your prompt" --config ./my-config.json

# Override use case
promptsentinel check "Your prompt" --use-case educational
```

#### Validate Command
Perform comprehensive validation with detailed analysis:
```bash
# Basic comprehensive validation
promptsentinel validate "Your prompt here"

# Output in JSON format
promptsentinel validate "Your prompt" --format json

# Use custom configuration
promptsentinel validate "Your prompt" --config ./config.json
```

#### Configuration Management
```bash
# Initialize default configuration
promptsentinel config init

# Show current configuration
promptsentinel config show

# Set configuration values
promptsentinel config set use_case educational
promptsentinel config set safety_level high
promptsentinel config set max_length 5000
```

### Configuration

The configuration file is stored at `~/.config/promptsentinel/config.json` by default. You can specify a custom path using the `--config` flag.

#### Configuration Options

- **use_case**: The intended use case for prompts (`general`, `educational`, `business`, `creative`)
- **safety_level**: Safety validation level (`low`, `medium`, `high`)
- **max_length**: Maximum allowed prompt length
- **min_length**: Minimum required prompt length
- **blocked_patterns**: Regex patterns to block
- **custom_rules**: Custom validation rules
- **require_approval**: Whether to require manual approval for certain prompts

#### Example Configuration

```json
{
  "use_case": "educational",
  "safety_level": "high",
  "max_length": 10000,
  "min_length": 10,
  "blocked_patterns": [
    "(?i)(password|secret|key)",
    "(?i)(inject|exploit|hack)"
  ],
  "custom_rules": {
    "no_personal_info": "(?i)(ssn|social security|credit card)"
  },
  "require_approval": false
}
```

## Output Formats

### Text Output (Default)
```
🔍 Prompt Validation Results
============================

Status: ✅ PASSED
Score: 85/100

Issues Found:
  1. ⚠️ [WARNING] Prompt contains blocked pattern: (?i)(password|secret|key)
     💡 Suggestion: Review and modify the flagged content

Recommendations:
  1. Consider reviewing flagged content
  2. Prompt looks good! Consider adding more specific instructions

Metadata:
  processing_time_ms: 15
  prompt_length: 245
  use_case: educational
```

### JSON Output
```json
{
  "is_valid": true,
  "score": 85,
  "issues": [
    {
      "type": "pattern",
      "severity": "warning",
      "message": "Prompt contains blocked pattern: (?i)(password|secret|key)",
      "suggestion": "Review and modify the flagged content"
    }
  ],
  "recommendations": [
    "Consider reviewing flagged content",
    "Prompt looks good! Consider adding more specific instructions"
  ],
  "metadata": {
    "processing_time_ms": 15,
    "prompt_length": 245,
    "use_case": "educational"
  }
}
```

## Development

### Building

```bash
# Build for current platform
make build

# Build for all platforms
make build-all

# Install to system
make install
```

### Testing

```bash
# Run tests
make test

# Run tests with coverage
make test-coverage

# Run linter
make lint
```

### Project Structure

```
.
├── cmd/promptsentinel/     # Main CLI application
├── internal/
│   ├── cli/               # CLI command implementations
│   ├── validator/         # Core validation logic
│   ├── auth/             # API key helpers
│   └── promptdb/         # Database utilities
├── docs/                 # Documentation
├── Makefile             # Build system
└── go.mod               # Go module definition
```

## Use Cases

### Educational Institutions
- Validate student prompts for safety
- Ensure compliance with educational standards
- Block inappropriate content

### Business Applications
- Validate customer-facing prompts
- Ensure brand safety
- Comply with business regulations

### Content Creation
- Validate creative prompts
- Ensure appropriate content
- Maintain creative freedom while ensuring safety

### Research and Development
- Validate research prompts
- Ensure ethical AI usage
- Maintain research integrity

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Run the test suite
6. Submit a pull request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Support

- **Documentation**: [docs/](docs/)
- **Issues**: [GitHub Issues](https://github.com/your-org/PromptSentinel/issues)
- **Discussions**: [GitHub Discussions](https://github.com/your-org/PromptSentinel/discussions)

## Changelog

See [CHANGELOG.md](CHANGELOG.md) for a detailed list of changes and new features.
