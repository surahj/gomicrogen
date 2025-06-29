#!/bin/bash
# One-liner installation script for gomicrogen
# Usage: curl -fsSL https://raw.githubusercontent.com/surahj/gomicrogen/main/install-oneline.sh | bash

set -e

# Configuration
REPO="surahj/gomicrogen"
BINARY_NAME="gomicrogen"
INSTALL_DIR="/usr/local/bin"
TEMP_DIR="/tmp/gomicrogen-install"

# Colors
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${BLUE}🚀 Installing gomicrogen...${NC}"

# Detect platform
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

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
    *) echo "❌ Unsupported OS: $OS"; exit 1 ;;
esac

# Get latest version
LATEST_VERSION=$(curl -s "https://api.github.com/repos/$REPO/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')

if [ -z "$LATEST_VERSION" ] || [ "$LATEST_VERSION" = "null" ]; then
    LATEST_VERSION="latest"
fi

# Download URL
if [ "$OS" = "windows" ]; then
    DOWNLOAD_URL="https://github.com/$REPO/releases/download/$LATEST_VERSION/gomicrogen-windows-$ARCH.zip"
    local_filename="gomicrogen-windows-$ARCH.zip"
else
    DOWNLOAD_URL="https://github.com/$REPO/releases/download/$LATEST_VERSION/gomicrogen-$OS-$ARCH.tar.gz"
    local_filename="gomicrogen-$OS-$ARCH.tar.gz"
fi

echo -e "${BLUE}📦 Downloading gomicrogen $LATEST_VERSION for $OS/$ARCH...${NC}"

# Download and install
TEMP_DIR=$(mktemp -d)
cd "$TEMP_DIR"

if command -v curl >/dev/null 2>&1; then
    curl -L -o "$local_filename" "$DOWNLOAD_URL"
elif command -v wget >/dev/null 2>&1; then
    wget -O "$local_filename" "$DOWNLOAD_URL"
else
    echo "❌ Neither curl nor wget is available. Please install one of them."
    exit 1
fi

# Extract the binary
if [ "$OS" = "windows" ]; then
    unzip -q "$local_filename"
    BINARY_NAME="gomicrogen.exe"
else
    tar -xzf "$local_filename"
    BINARY_NAME="gomicrogen"
fi

# Find the extracted binary (exclude archive files)
EXTRACTED_BINARY=$(find . -name "gomicrogen*" -type f ! -name "*.tar.gz" ! -name "*.zip" | head -1)
if [ -z "$EXTRACTED_BINARY" ]; then
    echo "❌ Failed to extract binary from archive"
    exit 1
fi

# Get just the filename without path
BINARY_FILENAME=$(basename "$EXTRACTED_BINARY")

if [ "$OS" != "windows" ]; then
    chmod +x "$EXTRACTED_BINARY"
fi

# Install to /usr/local/bin (or ~/.local/bin if no write permission)
if [ -w "/usr/local/bin" ] || command -v sudo >/dev/null 2>&1; then
    if [ -w "/usr/local/bin" ]; then
        cp "$EXTRACTED_BINARY" "/usr/local/bin/gomicrogen"
        chmod +x "/usr/local/bin/gomicrogen"
    else
        sudo cp "$EXTRACTED_BINARY" "/usr/local/bin/gomicrogen"
        sudo chmod +x "/usr/local/bin/gomicrogen"
    fi
    echo -e "${GREEN}✅ Installed to /usr/local/bin/${NC}"
else
    mkdir -p "$HOME/.local/bin"
    cp "$EXTRACTED_BINARY" "$HOME/.local/bin/gomicrogen"
    chmod +x "$HOME/.local/bin/gomicrogen"
    echo -e "${GREEN}✅ Installed to $HOME/.local/bin/${NC}"
    echo -e "${BLUE}💡 Please add $HOME/.local/bin to your PATH${NC}"
fi

# Cleanup
rm -rf "$TEMP_DIR"

echo -e "${GREEN}🎉 Installation completed!${NC}"
echo ""
echo "You can now use gomicrogen:"
echo "  gomicrogen --help"
echo ""
echo "Example usage:"
echo "  gomicrogen new my-service --module github.com/your-org/my-service" 