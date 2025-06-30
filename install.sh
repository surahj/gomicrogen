#!/bin/bash

# gomicrogen Installation Script
# This script downloads and installs gomicrogen CLI tool

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
REPO="surahj/gomicrogen"
BINARY_NAME="gomicrogen"
INSTALL_DIR="/usr/local/bin"
TEMP_DIR="/tmp/gomicrogen-install"

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

# Function to detect OS and architecture
detect_platform() {
    OS="$(uname -s | tr '[:upper:]' '[:lower:]')"
    ARCH="$(uname -m)"
    
    case "$ARCH" in
        x86_64) ARCH="amd64" ;;
        aarch64) ARCH="arm64" ;;
        armv7l) ARCH="arm" ;;
        armv8l) ARCH="arm64" ;;
    esac
    
    case "$OS" in
        darwin) OS="darwin" ;;
        linux) OS="linux" ;;
        msys*|cygwin*|mingw*) OS="windows" ;;
        *) print_error "Unsupported OS: $OS"; exit 1 ;;
    esac
    
    print_status "Detected platform: $OS/$ARCH"
}

# Function to get latest version
get_latest_version() {
    print_status "Fetching latest version..."
    
    # Try to get latest release from GitHub API
    if command -v curl >/dev/null 2>&1; then
        LATEST_VERSION=$(curl -s "https://api.github.com/repos/$REPO/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')
    elif command -v wget >/dev/null 2>&1; then
        LATEST_VERSION=$(wget -qO- "https://api.github.com/repos/$REPO/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')
    else
        print_warning "Could not fetch latest version, using 'latest'"
        LATEST_VERSION="latest"
    fi
    
    if [ -z "$LATEST_VERSION" ] || [ "$LATEST_VERSION" = "null" ]; then
        print_warning "Could not determine latest version, using 'latest'"
        LATEST_VERSION="latest"
    fi
    
    print_status "Latest version: $LATEST_VERSION"
}

# Function to download binary
download_binary() {
    local version=$1
    local download_url=""
    
    if [ "$OS" = "windows" ]; then
        download_url="https://github.com/$REPO/releases/download/$version/${BINARY_NAME}-${OS}-${ARCH}.zip"
        local_filename="${BINARY_NAME}-${OS}-${ARCH}.zip"
    else
        download_url="https://github.com/$REPO/releases/download/$version/${BINARY_NAME}-${OS}-${ARCH}.tar.gz"
        local_filename="${BINARY_NAME}-${OS}-${ARCH}.tar.gz"
    fi
    
    print_status "Downloading from: $download_url"
    
    # Create temp directory
    mkdir -p "$TEMP_DIR"
    cd "$TEMP_DIR"
    
    # Download binary
    if command -v curl >/dev/null 2>&1; then
        curl -L -o "$local_filename" "$download_url"
    elif command -v wget >/dev/null 2>&1; then
        wget -O "$local_filename" "$download_url"
    else
        print_error "Neither curl nor wget is available. Please install one of them."
        exit 1
    fi
    
    # Check if download was successful
    if [ ! -f "$local_filename" ]; then
        print_error "Failed to download binary"
        exit 1
    fi
    
    print_success "Download completed"
}

# Function to install binary
install_binary() {
    local local_filename=$1
    
    print_status "Extracting package..."
    
    # Extract the package
    if [ "$OS" = "windows" ]; then
        if command -v unzip >/dev/null 2>&1; then
            unzip -q "$local_filename"
        else
            print_error "unzip is required but not installed. Please install unzip."
            exit 1
        fi
    else
        if command -v tar >/dev/null 2>&1; then
            tar -xzf "$local_filename"
        else
            print_error "tar is required but not installed. Please install tar."
            exit 1
        fi
    fi
    
    # Find the extracted directory
    local package_dir=""
    if [ "$OS" = "windows" ]; then
        package_dir=$(find . -name "${BINARY_NAME}-${OS}-${ARCH}-package" -type d | head -1)
    else
        package_dir=$(find . -name "${BINARY_NAME}-${OS}-${ARCH}-package" -type d | head -1)
    fi
    
    if [ -z "$package_dir" ]; then
        print_error "Could not find extracted package directory"
        exit 1
    fi
    
    print_status "Installing to $INSTALL_DIR..."
    
    # Check if we have write permissions
    if [ ! -w "$INSTALL_DIR" ] && [ "$(id -u)" != "0" ]; then
        print_warning "No write permission to $INSTALL_DIR, using sudo"
        if command -v sudo >/dev/null 2>&1; then
            # Install binary
            if [ "$OS" = "windows" ]; then
                sudo cp "$package_dir/${BINARY_NAME}-${OS}-${ARCH}.exe" "$INSTALL_DIR/$BINARY_NAME"
            else
                sudo cp "$package_dir/${BINARY_NAME}-${OS}-${ARCH}" "$INSTALL_DIR/$BINARY_NAME"
            fi
            sudo chmod +x "$INSTALL_DIR/$BINARY_NAME"
            
            # Install templates directory
            sudo mkdir -p "$INSTALL_DIR/templates"
            sudo cp -r "$package_dir/templates"/* "$INSTALL_DIR/templates/"
        else
            print_error "sudo is not available. Please run this script as root or install sudo."
            exit 1
        fi
    else
        # Install binary
        if [ "$OS" = "windows" ]; then
            cp "$package_dir/${BINARY_NAME}-${OS}-${ARCH}.exe" "$INSTALL_DIR/$BINARY_NAME"
        else
            cp "$package_dir/${BINARY_NAME}-${OS}-${ARCH}" "$INSTALL_DIR/$BINARY_NAME"
        fi
        chmod +x "$INSTALL_DIR/$BINARY_NAME"
        
        # Install templates directory
        mkdir -p "$INSTALL_DIR/templates"
        cp -r "$package_dir/templates"/* "$INSTALL_DIR/templates/"
    fi
    
    print_success "Installation completed"
}

# Function to verify installation
verify_installation() {
    print_status "Verifying installation..."
    
    if command -v "$BINARY_NAME" >/dev/null 2>&1; then
        print_success "$BINARY_NAME is now available in your PATH"
        print_status "Version: $($BINARY_NAME version 2>/dev/null || echo 'Version info not available')"
    else
        print_error "Installation verification failed. $BINARY_NAME is not found in PATH"
        exit 1
    fi
}

# Function to cleanup
cleanup() {
    print_status "Cleaning up temporary files..."
    rm -rf "$TEMP_DIR"
}

# Function to show usage
show_usage() {
    echo "Usage: $0 [OPTIONS]"
    echo ""
    echo "Options:"
    echo "  -v, --version VERSION    Install specific version (default: latest)"
    echo "  -d, --directory DIR      Install to specific directory (default: $INSTALL_DIR)"
    echo "  -h, --help              Show this help message"
    echo ""
    echo "Examples:"
    echo "  $0                      # Install latest version"
    echo "  $0 -v v1.0.0           # Install specific version"
    echo "  $0 -d ~/.local/bin     # Install to custom directory"
}

# Main function
main() {
    local version=""
    local install_dir="$INSTALL_DIR"
    
    # Parse command line arguments
    while [[ $# -gt 0 ]]; do
        case $1 in
            -v|--version)
                version="$2"
                shift 2
                ;;
            -d|--directory)
                install_dir="$2"
                shift 2
                ;;
            -h|--help)
                show_usage
                exit 0
                ;;
            *)
                print_error "Unknown option: $1"
                show_usage
                exit 1
                ;;
        esac
    done
    
    print_status "Starting gomicrogen installation..."
    
    # Detect platform
    detect_platform
    
    # Get version if not specified
    if [ -z "$version" ]; then
        get_latest_version
        version="$LATEST_VERSION"
    fi
    
    # Download binary
    download_binary "$version"
    
    # Install binary
    INSTALL_DIR="$install_dir"
    install_binary "$local_filename"
    
    # Verify installation
    verify_installation
    
    # Cleanup
    cleanup
    
    print_success "Installation completed successfully!"
    echo ""
    echo "You can now use gomicrogen:"
    echo "  gomicrogen --help"
    echo ""
    echo "Example usage:"
    echo "  gomicrogen new my-service --module github.com/your-org/my-service"
}

# Trap to cleanup on exit
trap cleanup EXIT

# Run main function
main "$@" 