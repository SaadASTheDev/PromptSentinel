#!/bin/bash

# PromptSentinel CLI Installation Script
# This script installs PromptSentinel CLI tool

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
BINARY_NAME="promptsentinel"
INSTALL_DIR="/usr/local/bin"
CONFIG_DIR="$HOME/.config/promptsentinel"
VERSION="1.0.0"

# Function to print colored output
print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Function to check if command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Function to detect OS and architecture
detect_platform() {
    local os
    local arch
    
    case "$(uname -s)" in
        Linux*)     os="linux" ;;
        Darwin*)    os="darwin" ;;
        CYGWIN*)    os="windows" ;;
        MINGW*)     os="windows" ;;
        *)          os="unknown" ;;
    esac
    
    case "$(uname -m)" in
        x86_64)     arch="amd64" ;;
        arm64)      arch="arm64" ;;
        aarch64)    arch="arm64" ;;
        *)          arch="unknown" ;;
    esac
    
    echo "${os}-${arch}"
}

# Function to download binary
download_binary() {
    local platform=$1
    local download_url="https://github.com/your-org/PromptSentinel/releases/download/v${VERSION}/promptsentinel-${platform}"
    
    print_status "Downloading PromptSentinel for ${platform}..."
    
    # Create temporary directory
    local temp_dir=$(mktemp -d)
    cd "$temp_dir"
    
    # Download binary (this is a placeholder - replace with actual download logic)
    print_warning "Download functionality not implemented yet. Please build from source."
    print_status "To build from source:"
    echo "  git clone https://github.com/your-org/PromptSentinel.git"
    echo "  cd PromptSentinel"
    echo "  make install"
    
    # Cleanup
    cd - > /dev/null
    rm -rf "$temp_dir"
}

# Function to install from source
install_from_source() {
    print_status "Installing from source..."
    
    # Check if Go is installed
    if ! command_exists go; then
        print_error "Go is not installed. Please install Go 1.22 or newer."
        print_status "Visit https://golang.org/dl/ to download Go."
        exit 1
    fi
    
    # Check Go version
    local go_version=$(go version | cut -d' ' -f3 | sed 's/go//')
    local required_version="1.22"
    
    if [ "$(printf '%s\n' "$required_version" "$go_version" | sort -V | head -n1)" != "$required_version" ]; then
        print_error "Go version $go_version is too old. Please install Go 1.22 or newer."
        exit 1
    fi
    
    print_success "Go version $go_version is compatible"
    
    # Clone repository
    local repo_dir=$(mktemp -d)
    print_status "Cloning repository to $repo_dir..."
    
    if ! git clone https://github.com/your-org/PromptSentinel.git "$repo_dir"; then
        print_error "Failed to clone repository. Please check your internet connection."
        exit 1
    fi
    
    cd "$repo_dir"
    
    # Build and install
    print_status "Building PromptSentinel..."
    if ! make install; then
        print_error "Failed to build PromptSentinel. Please check the build logs."
        exit 1
    fi
    
    # Cleanup
    cd - > /dev/null
    rm -rf "$repo_dir"
    
    print_success "PromptSentinel built and installed successfully"
}

# Function to create configuration directory
create_config_dir() {
    print_status "Creating configuration directory..."
    
    if [ ! -d "$CONFIG_DIR" ]; then
        mkdir -p "$CONFIG_DIR"
        print_success "Configuration directory created: $CONFIG_DIR"
    else
        print_status "Configuration directory already exists: $CONFIG_DIR"
    fi
}

# Function to initialize configuration
init_config() {
    print_status "Initializing configuration..."
    
    if command_exists "$BINARY_NAME"; then
        if "$BINARY_NAME" config init; then
            print_success "Configuration initialized successfully"
        else
            print_warning "Failed to initialize configuration. You can run 'promptsentinel config init' later."
        fi
    else
        print_warning "PromptSentinel not found in PATH. Configuration will be initialized on first run."
    fi
}

# Function to verify installation
verify_installation() {
    print_status "Verifying installation..."
    
    if command_exists "$BINARY_NAME"; then
        local version=$("$BINARY_NAME" --version 2>/dev/null || echo "unknown")
        print_success "PromptSentinel installed successfully"
        print_status "Version: $version"
        print_status "Location: $(which $BINARY_NAME)"
    else
        print_error "PromptSentinel not found in PATH"
        print_status "Please ensure $INSTALL_DIR is in your PATH"
        exit 1
    fi
}

# Function to show usage instructions
show_usage() {
    print_success "Installation completed!"
    echo
    echo "Usage examples:"
    echo "  $BINARY_NAME check \"Write a story about a cat\""
    echo "  $BINARY_NAME validate \"Your prompt here\" --format json"
    echo "  $BINARY_NAME config show"
    echo "  $BINARY_NAME config init"
    echo
    echo "For more information, run: $BINARY_NAME --help"
}

# Main installation function
main() {
    print_status "PromptSentinel CLI Installation Script"
    print_status "======================================"
    echo
    
    # Check if already installed
    if command_exists "$BINARY_NAME"; then
        print_warning "PromptSentinel is already installed"
        local current_version=$("$BINARY_NAME" --version 2>/dev/null || echo "unknown")
        print_status "Current version: $current_version"
        print_status "Target version: $VERSION"
        
        read -p "Do you want to reinstall? (y/N): " -n 1 -r
        echo
        if [[ ! $REPLY =~ ^[Yy]$ ]]; then
            print_status "Installation cancelled"
            exit 0
        fi
    fi
    
    # Detect platform
    local platform=$(detect_platform)
    print_status "Detected platform: $platform"
    
    # Install from source (since we don't have pre-built binaries yet)
    install_from_source
    
    # Create configuration directory
    create_config_dir
    
    # Initialize configuration
    init_config
    
    # Verify installation
    verify_installation
    
    # Show usage instructions
    show_usage
}

# Run main function
main "$@"
