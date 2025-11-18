#!/bin/bash

# Goverter Installation Script
# Supports Linux, macOS, and Windows (via WSL or Git Bash)

set -e

VERSION="1.0.0"
INSTALL_DIR="$HOME/.local/bin"
CONFIG_DIR="$HOME/.config/goverter"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

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

# Detect OS
detect_os() {
    if [[ "$OSTYPE" == "linux-gnu"* ]]; then
        echo "linux"
    elif [[ "$OSTYPE" == "darwin"* ]]; then
        echo "macos"
    elif [[ "$OSTYPE" == "msys" ]] || [[ "$OSTYPE" == "cygwin" ]]; then
        echo "windows"
    else
        echo "unknown"
    fi
}

# Install dependencies for Linux
install_linux_deps() {
    print_status "Installing dependencies for Linux..."
    
    if command -v apt-get &> /dev/null; then
        # Ubuntu/Debian
        sudo apt-get update
        sudo apt-get install -y ffmpeg imagemagick pandoc
    elif command -v yum &> /dev/null; then
        # CentOS/RHEL/Fedora
        sudo yum install -y epel-release
        sudo yum install -y ffmpeg ImageMagick pandoc
    elif command -v pacman &> /dev/null; then
        # Arch Linux
        sudo pacman -S --noconfirm ffmpeg imagemagick pandoc
    elif command -v zypper &> /dev/null; then
        # openSUSE
        sudo zypper install -y ffmpeg ImageMagick pandoc
    else
        print_warning "Unable to detect package manager. Please install manually:"
        echo "  - FFmpeg: https://ffmpeg.org/download.html"
        echo "  - ImageMagick: https://imagemagick.org/script/download.php"
        echo "  - Pandoc: https://pandoc.org/installing.html"
    fi
}

# Install dependencies for macOS
install_macos_deps() {
    print_status "Installing dependencies for macOS..."
    
    if command -v brew &> /dev/null; then
        brew install ffmpeg imagemagick pandoc
    else
        print_warning "Homebrew not found. Please install Homebrew first:"
        echo "  /bin/bash -c \"\$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)\""
        echo "Then run this installer again."
        exit 1
    fi
}

# Install dependencies for Windows
install_windows_deps() {
    print_status "Windows detected. Please install dependencies manually:"
    echo ""
    echo "1. FFmpeg:"
    echo "   Download from: https://ffmpeg.org/download.html"
    echo "   Extract and add to PATH"
    echo ""
    echo "2. ImageMagick:"
    echo "   Download from: https://imagemagick.org/script/download.php"
    echo "   Install and add to PATH"
    echo ""
    echo "3. Pandoc:"
    echo "   Download from: https://pandoc.org/installing.html"
    echo "   Install and add to PATH"
    echo ""
    read -p "Press Enter after installing dependencies..."
}

# Create directories
create_dirs() {
    print_status "Creating directories..."
    mkdir -p "$INSTALL_DIR"
    mkdir -p "$CONFIG_DIR"
}

# Download and install Goverter
install_goverter() {
    print_status "Downloading Goverter..."
    
    # Get the latest release URL
    RELEASE_URL="https://github.com/ezraclintoc/goverter/releases/latest/download/goverter-${VERSION}-$(detect_os)-amd64.tar.gz"
    
    if command -v wget &> /dev/null; then
        wget -O goverter.tar.gz "$RELEASE_URL"
    elif command -v curl &> /dev/null; then
        curl -L -o goverter.tar.gz "$RELEASE_URL"
    else
        print_error "Neither wget nor curl found. Please install one of them."
        exit 1
    fi
    
    print_status "Extracting Goverter..."
    tar -xzf goverter.tar.gz
    
    print_status "Installing binaries..."
    cp goverter-cli "$INSTALL_DIR/"
    cp goverter-gui "$INSTALL_DIR/"
    
    # Make executables
    chmod +x "$INSTALL_DIR/goverter-cli"
    chmod +x "$INSTALL_DIR/goverter-gui"
    
    # Clean up
    rm goverter.tar.gz
    
    print_success "Goverter installed successfully!"
}

# Update PATH
update_path() {
    OS=$(detect_os)
    case $OS in
        "linux"|"macos")
            if [[ ":$PATH:" != *":$INSTALL_DIR:"* ]]; then
                print_status "Adding $INSTALL_DIR to PATH..."
                
                # Add to shell profile
                SHELL_RC=""
                if [[ -f "$HOME/.bashrc" ]]; then
                    SHELL_RC="$HOME/.bashrc"
                elif [[ -f "$HOME/.zshrc" ]]; then
                    SHELL_RC="$HOME/.zshrc"
                fi
                
                if [[ -n "$SHELL_RC" ]]; then
                    echo "" >> "$SHELL_RC"
                    echo "# Goverter" >> "$SHELL_RC"
                    echo "export PATH=\"\$PATH:$INSTALL_DIR\"" >> "$SHELL_RC"
                    print_status "Added to $SHELL_RC"
                fi
            fi
            ;;
        "windows")
            print_warning "Please add $INSTALL_DIR to your Windows PATH manually."
            ;;
    esac
}

# Create desktop entry for Linux
create_desktop_entry() {
    if [[ "$OSTYPE" == "linux-gnu"* ]]; then
        print_status "Creating desktop entry..."
        
        cat > "$HOME/.local/share/applications/goverter.desktop" << EOF
[Desktop Entry]
Version=1.0
Type=Application
Name=Goverter
Comment=A comprehensive file conversion tool
Exec=$INSTALL_DIR/goverter-gui
Icon=$CONFIG_DIR/goverter.png
Terminal=false
Categories=AudioVideo;AudioVideoEditing;Graphics;Utility;
EOF
        
        print_success "Desktop entry created!"
    fi
}

# Main installation flow
main() {
    echo "Goverter Installer v$VERSION"
    echo "=========================="
    echo ""
    
    OS=$(detect_os)
    print_status "Detected OS: $OS"
    
    # Check if already installed
    if [[ -f "$INSTALL_DIR/goverter-cli" ]]; then
        print_warning "Goverter is already installed!"
        read -p "Do you want to reinstall? (y/N): " -n 1 -r
        echo
        if [[ ! $REPLY =~ ^[Yy]$ ]]; then
            print_status "Installation cancelled."
            exit 0
        fi
    fi
    
    # Install dependencies
    case $OS in
        "linux")
            install_linux_deps
            ;;
        "macos")
            install_macos_deps
            ;;
        "windows")
            install_windows_deps
            ;;
        *)
            print_error "Unsupported OS: $OS"
            exit 1
            ;;
    esac
    
    # Install Goverter
    create_dirs
    install_goverter
    update_path
    create_desktop_entry
    
    echo ""
    print_success "Installation completed!"
    echo ""
    echo "Usage:"
    echo "  CLI: goverter-cli --help"
    echo "  GUI: goverter-gui"
    echo ""
    echo "Please restart your terminal or run:"
    echo "  export PATH=\"\$PATH:$INSTALL_DIR\""
    echo ""
    
    if [[ "$OS" != "windows" ]]; then
        print_status "You can now run 'goverter-gui' to start the GUI application."
    fi
}

# Check for help flag
if [[ "$1" == "--help" ]] || [[ "$1" == "-h" ]]; then
    echo "Goverter Installer"
    echo "Usage: $0 [options]"
    echo ""
    echo "Options:"
    echo "  --help, -h     Show this help message"
    echo "  --uninstall     Uninstall Goverter"
    echo ""
    exit 0
fi

# Check for uninstall flag
if [[ "$1" == "--uninstall" ]]; then
    print_status "Uninstalling Goverter..."
    
    # Remove binaries
    rm -f "$INSTALL_DIR/goverter-cli"
    rm -f "$INSTALL_DIR/goverter-gui"
    
    # Remove desktop entry
    rm -f "$HOME/.local/share/applications/goverter.desktop"
    
    # Remove config directory
    rm -rf "$CONFIG_DIR"
    
    print_success "Goverter uninstalled successfully!"
    exit 0
fi

# Run main installation
main "$@"