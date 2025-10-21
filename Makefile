# PromptSentinel CLI Makefile

# Variables
BINARY_NAME=promptsentinel
BUILD_DIR=build
VERSION=1.0.0
GO_VERSION=1.22

# Default target
.PHONY: all
all: clean build

# Clean build artifacts
.PHONY: clean
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf $(BUILD_DIR)
	@go clean

# Install dependencies
.PHONY: deps
deps:
	@echo "Installing dependencies..."
	@go mod tidy
	@go mod download

# Build the binary
.PHONY: build
build: deps
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	@go build -ldflags "-X main.version=$(VERSION)" -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/promptsentinel

# Build for multiple platforms
.PHONY: build-all
build-all: deps
	@echo "Building for multiple platforms..."
	@mkdir -p $(BUILD_DIR)
	@GOOS=linux GOARCH=amd64 go build -ldflags "-X main.version=$(VERSION)" -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 ./cmd/promptsentinel
	@GOOS=darwin GOARCH=amd64 go build -ldflags "-X main.version=$(VERSION)" -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 ./cmd/promptsentinel
	@GOOS=darwin GOARCH=arm64 go build -ldflags "-X main.version=$(VERSION)" -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 ./cmd/promptsentinel
	@GOOS=windows GOARCH=amd64 go build -ldflags "-X main.version=$(VERSION)" -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe ./cmd/promptsentinel

# Install the binary to GOPATH/bin
.PHONY: install
install: build
	@echo "Installing $(BINARY_NAME)..."
	@go install ./cmd/promptsentinel

# Run tests
.PHONY: test
test:
	@echo "Running tests..."
	@go test -v ./...

# Run tests with coverage
.PHONY: test-coverage
test-coverage:
	@echo "Running tests with coverage..."
	@go test -v -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Run linting
.PHONY: lint
lint:
	@echo "Running linter..."
	@go vet ./...
	@go fmt ./...

# Create a release package
.PHONY: release
release: build-all
	@echo "Creating release package..."
	@mkdir -p $(BUILD_DIR)/release
	@cp $(BUILD_DIR)/$(BINARY_NAME)-* $(BUILD_DIR)/release/
	@cp README.md $(BUILD_DIR)/release/
	@cp LICENSE $(BUILD_DIR)/release/ 2>/dev/null || true
	@cd $(BUILD_DIR) && tar -czf release.tar.gz release/
	@echo "Release package created: $(BUILD_DIR)/release.tar.gz"

# Show help
.PHONY: help
help:
	@echo "PromptSentinel CLI Build System"
	@echo "==============================="
	@echo ""
	@echo "Available targets:"
	@echo "  all          - Clean and build the binary"
	@echo "  clean        - Remove build artifacts"
	@echo "  deps         - Install dependencies"
	@echo "  build        - Build the binary for current platform"
	@echo "  build-all    - Build for multiple platforms"
	@echo "  install      - Install binary to GOPATH/bin"
	@echo "  test         - Run tests"
	@echo "  test-coverage- Run tests with coverage report"
	@echo "  lint         - Run linter and formatter"
	@echo "  release      - Create release package"
	@echo "  help         - Show this help message"
	@echo ""
	@echo "Examples:"
	@echo "  make build        # Build for current platform"
	@echo "  make install      # Install to system"
	@echo "  make test         # Run all tests"
