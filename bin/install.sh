#!/bin/sh
set -e

# WebKit Installation Script
# Usage: curl -sSL https://raw.githubusercontent.com/ainsleydev/webkit/main/bin/install.sh | sh

REPO="ainsleydev/webkit"
VERSION="${VERSION:-latest}"

# Default to user's local bin, fallback to /usr/local/bin if it doesn't exist
if [ -z "$INSTALL_DIR" ]; then
    if [ -d "$HOME/.local/bin" ] || mkdir -p "$HOME/.local/bin" 2>/dev/null; then
        INSTALL_DIR="$HOME/.local/bin"
    else
        INSTALL_DIR="/usr/local/bin"
    fi
fi

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

info() {
    printf "${GREEN}==>${NC} %s\n" "$1" >&2
}

warn() {
    printf "${YELLOW}Warning:${NC} %s\n" "$1" >&2
}

error() {
    printf "${RED}Error:${NC} %s\n" "$1" >&2
    exit 1
}

# Detect OS
detect_os() {
    OS="$(uname -s)"
    case "$OS" in
        Linux*)     OS='linux';;
        Darwin*)    OS='darwin';;
        MINGW*|MSYS*|CYGWIN*) OS='windows';;
        *)          error "Unsupported operating system: $OS";;
    esac
    echo "$OS"
}

# Detect architecture
detect_arch() {
    ARCH="$(uname -m)"
    case "$ARCH" in
        x86_64|amd64)   ARCH='x86_64';;
        aarch64|arm64)  ARCH='arm64';;
        armv7l|armv6l)  ARCH='arm';;
        *)              error "Unsupported architecture: $ARCH";;
    esac
    echo "$ARCH"
}

# Get latest release version from GitHub
get_latest_version() {
    info "Fetching latest release information from GitHub..."

    # Try to get the latest release using GitHub API
    API_RESPONSE=$(curl -sSL "https://api.github.com/repos/$REPO/releases/latest" 2>&1)

    if [ $? -ne 0 ]; then
        error "Failed to connect to GitHub API. Check your internet connection."
    fi

    LATEST_VERSION=$(echo "$API_RESPONSE" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/' | head -n 1)

    if [ -z "$LATEST_VERSION" ]; then
        error "Failed to fetch latest version. API response may be rate-limited or no releases exist.\nTry setting VERSION environment variable explicitly: VERSION=v0.0.3 sh install.sh"
    fi

    echo "$LATEST_VERSION"
}

# Download and install webkit
install_webkit() {
    OS=$(detect_os)
    ARCH=$(detect_arch)

    if [ "$VERSION" = "latest" ]; then
        info "Fetching latest version..."
        VERSION=$(get_latest_version)
    fi

    info "Installing webkit $VERSION for $OS/$ARCH..."

    # Construct download URL and file extension
    # Note: This naming must match GoReleaser's archive naming template
    BINARY_NAME="webkit_${OS}_${ARCH}"

    if [ "$OS" = "windows" ]; then
        ARCHIVE_EXT=".zip"
    else
        ARCHIVE_EXT=".tar.gz"
    fi

    ARCHIVE_NAME="${BINARY_NAME}${ARCHIVE_EXT}"
    DOWNLOAD_URL="https://github.com/$REPO/releases/download/$VERSION/$ARCHIVE_NAME"

    info "Archive name: $ARCHIVE_NAME"
    info "Downloading from: $DOWNLOAD_URL"

    # Create temporary directory
    TMP_DIR=$(mktemp -d)
    trap "rm -rf $TMP_DIR" EXIT

    # Download archive
    if ! curl -sSL -o "$TMP_DIR/$ARCHIVE_NAME" "$DOWNLOAD_URL"; then
        error "Failed to download webkit. Please check the release exists at:\n$DOWNLOAD_URL"
    fi

    # Extract archive
    info "Extracting archive..."
    cd "$TMP_DIR"

    if [ "$OS" = "windows" ]; then
        if ! unzip -q "$ARCHIVE_NAME"; then
            error "Failed to extract archive"
        fi
    else
        if ! tar -xzf "$ARCHIVE_NAME"; then
            error "Failed to extract archive"
        fi
    fi

    # Find the webkit binary
    WEBKIT_BINARY="$TMP_DIR/webkit"
    if [ "$OS" = "windows" ]; then
        WEBKIT_BINARY="$TMP_DIR/webkit.exe"
    fi

    if [ ! -f "$WEBKIT_BINARY" ]; then
        error "webkit binary not found in archive"
    fi

    # Make binary executable
    chmod +x "$WEBKIT_BINARY"

    # Check if we need sudo for installation
    if [ -w "$INSTALL_DIR" ]; then
        SUDO=""
        info "Installing to $INSTALL_DIR (no sudo required)"
    else
        SUDO="sudo"
        warn "Installing to $INSTALL_DIR requires sudo privileges"
    fi

    # Install binary
    TARGET="$INSTALL_DIR/webkit"
    if [ "$OS" = "windows" ]; then
        TARGET="$INSTALL_DIR/webkit.exe"
    fi

    if ! $SUDO mv "$WEBKIT_BINARY" "$TARGET"; then
        error "Failed to install webkit to $INSTALL_DIR"
    fi

    info "✓ webkit installed successfully to $TARGET"

    # Verify installation
    if command -v webkit >/dev/null 2>&1; then
        info "✓ webkit is now available in your PATH"
        webkit version
    else
        warn "webkit was installed but is not in your PATH"

        if [ "$INSTALL_DIR" = "$HOME/.local/bin" ]; then
            echo
            info "Add this to your shell profile (~/.bashrc, ~/.zshrc, etc.):"
            echo "  export PATH=\"\$HOME/.local/bin:\$PATH\""
            echo
            info "Then reload your shell or run:"
            echo "  source ~/.bashrc  # or ~/.zshrc"
        else
            warn "Add $INSTALL_DIR to your PATH or use the full path: $TARGET"
        fi
    fi
}

# Check dependencies
check_dependencies() {
    if ! command -v curl >/dev/null 2>&1; then
        error "curl is required but not installed"
    fi

    OS=$(detect_os)
    if [ "$OS" = "windows" ]; then
        if ! command -v unzip >/dev/null 2>&1; then
            error "unzip is required but not installed"
        fi
    else
        if ! command -v tar >/dev/null 2>&1; then
            error "tar is required but not installed"
        fi
    fi
}

# Main execution
main() {
    info "WebKit Installer"
    info "Repository: https://github.com/$REPO"
    echo

    check_dependencies
    install_webkit

    echo
    info "Installation complete! 🎉"
    info "Run 'webkit --help' to get started"
}

main
