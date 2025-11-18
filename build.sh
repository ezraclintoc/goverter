#!/bin/bash

# Enhanced build script for Goverter
# Downloads dependencies if missing and warns user about missing tools

set -e

VERSION="1.0.0"
BUILD_DIR="build"

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

# Check if command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Check and install Go dependencies
check_go_deps() {
    print_status "Checking Go dependencies..."
    
    # Check if go.mod exists
    if [[ ! -f "go.mod" ]]; then
        print_error "go.mod not found. Please run from project root."
        exit 1
    fi
    
    # Download dependencies
    print_status "Downloading Go dependencies..."
    if ! go mod download; then
        print_error "Failed to download Go dependencies"
        exit 1
    fi
    
    # Verify dependencies
    if ! go mod verify; then
        print_warning "Some dependencies could not be verified"
    fi
    
    print_success "Go dependencies ready"
}

# Check system dependencies
check_system_deps() {
    print_status "Checking system dependencies..."
    
    local missing_deps=()
    local optional_missing=()
    
    # Check required tools
    if ! command_exists go; then
        print_error "Go is required but not installed"
        print_status "Please install Go from: https://golang.org/dl/"
        exit 1
    fi
    
    # Check optional tools with warnings
    if ! command_exists ffmpeg; then
        optional_missing+=("ffmpeg (video/audio conversion)")
    fi
    
    if ! command_exists magick; then
        optional_missing+=("imagemagick (image conversion)")
    fi
    
    if ! command_exists pandoc; then
        optional_missing+=("pandoc (document conversion)")
    fi
    
    # Report missing dependencies
    if [[ ${#optional_missing[@]} -gt 0 ]]; then
        print_warning "Missing optional dependencies:"
        for dep in "${optional_missing[@]}"; do
            echo "  • $dep"
        done
        
        echo ""
        show_manual_install_commands
        echo ""
        
        read -p "Continue building anyway? (y/N): " -n 1 -r
        echo
        if [[ ! $REPLY =~ ^[Yy]$ ]]; then
            print_status "Build cancelled. Install dependencies and try again."
            exit 0
        fi
    else
        print_success "All dependencies found!"
    fi
}

# Try to install missing tools automatically (Linux only)
try_auto_install() {
    if [[ "$OSTYPE" != "linux-gnu"* ]] && [[ "$OSTYPE" != "linux"* ]]; then
        return
    fi
    
    print_status "Attempting automatic dependency installation..."
    
    local packages=()
    
    # Check missing packages
    if ! command_exists ffmpeg; then
        packages+=("ffmpeg")
    fi
    if ! command_exists magick; then
        packages+=("imagemagick")
    fi
    if ! command_exists pandoc; then
        packages+=("pandoc")
    fi
    
    if [[ ${#packages[@]} -eq 0 ]]; then
        print_success "All dependencies already installed!"
        return
    fi
    
    # Detect distribution and package manager
    if command_exists apt-get; then
        # Ubuntu/Debian/Pop!_OS etc.
        install_with_apt "${packages[@]}"
    elif command_exists dnf; then
        # Fedora/RHEL/CentOS Stream
        install_with_dnf "${packages[@]}"
    elif command_exists yum; then
        # CentOS 7/older RHEL
        install_with_yum "${packages[@]}"
    elif command_exists pacman; then
        # Arch Linux
        install_with_pacman "${packages[@]}"
    elif command_exists zypper; then
        # openSUSE
        install_with_zypper "${packages[@]}"
    else
        print_warning "Unsupported package manager. Please install manually:"
        show_manual_install_commands
        return 1
    fi
}

# Install with apt-get (Ubuntu/Debian)
install_with_apt() {
    local packages=("$@")
    print_status "Using apt-get to install missing dependencies..."
    print_status "Installing: ${packages[*]}"
    
    if sudo apt-get update && sudo apt-get install -y "${packages[@]}"; then
        print_success "Dependencies installed successfully!"
    else
        print_warning "APT installation failed. Please install manually."
        show_manual_install_commands
        return 1
    fi
}

# Install with dnf (Fedora/RHEL)
install_with_dnf() {
    local packages=("$@")
    print_status "Using DNF to install missing dependencies..."
    
    # Enable RPM Fusion for multimedia packages
    print_status "Enabling RPM Fusion repository..."
    sudo dnf install -y epel-release || true
    sudo dnf install -y https://download1.rpmfusion.org/free/fedora/rpmfusion-free-release-$(rpm -E %fedora).noarch.rpm
    sudo dnf install -y https://download1.rpmfusion.org/nonfree/fedora/rpmfusion-nonfree-release-$(rpm -E %fedora).noarch.rpm
    
    print_status "Installing: ${packages[*]}"
    if sudo dnf install -y "${packages[@]}"; then
        print_success "Dependencies installed successfully!"
    else
        print_warning "DNF installation failed. Please install manually."
        show_manual_install_commands
        return 1
    fi
}

# Install with yum (CentOS 7)
install_with_yum() {
    local packages=("$@")
    print_status "Using YUM to install missing dependencies..."
    
    # Enable EPEL repository
    print_status "Enabling EPEL repository..."
    sudo yum install -y epel-release
    
    print_status "Installing: ${packages[*]}"
    if sudo yum install -y "${packages[@]}"; then
        print_success "Dependencies installed successfully!"
    else
        print_warning "YUM installation failed. Please install manually."
        show_manual_install_commands
        return 1
    fi
}

# Install with pacman (Arch Linux)
install_with_pacman() {
    local packages=("$@")
    print_status "Using Pacman to install missing dependencies..."
    
    # Convert package names for Arch
    local arch_packages=()
    for pkg in "${packages[@]}"; do
        case $pkg in
            "imagemagick")
                arch_packages+=("imagemagick")
                ;;
            *)
                arch_packages+=("$pkg")
                ;;
        esac
    done
    
    print_status "Installing: ${arch_packages[*]}"
    if sudo pacman -S --noconfirm "${arch_packages[@]}"; then
        print_success "Dependencies installed successfully!"
    else
        print_warning "Pacman installation failed. Please install manually."
        show_manual_install_commands
        return 1
    fi
}

# Install with zypper (openSUSE)
install_with_zypper() {
    local packages=("$@")
    print_status "Using Zypper to install missing dependencies..."
    
    # Add repositories if needed
    sudo zypper --non-interactive addrepo -f https://download.opensuse.org/repositories/multimedia/openSUSE_Leap_$(rpm -E %suse_version).repo || true
    
    print_status "Installing: ${packages[*]}"
    if sudo zypper --non-interactive install -y "${packages[@]}"; then
        print_success "Dependencies installed successfully!"
    else
        print_warning "Zypper installation failed. Please install manually."
        show_manual_install_commands
        return 1
    fi
}

# Show manual installation commands
show_manual_install_commands() {
    echo ""
    print_status "Manual installation commands:"
    echo ""
    echo "🐧 Ubuntu/Debian/Pop!_OS:"
    echo "  sudo apt update"
    echo "  sudo apt install ffmpeg imagemagick pandoc"
    echo ""
    echo "🐧 Fedora:"
    echo "  sudo dnf install https://download1.rpmfusion.org/free/fedora/rpmfusion-free-release-\$(rpm -E %fedora).noarch.rpm"
    echo "  sudo dnf install https://download1.rpmfusion.org/nonfree/fedora/rpmfusion-nonfree-release-\$(rpm -E %fedora).noarch.rpm"
    echo "  sudo dnf install ffmpeg ImageMagick pandoc"
    echo ""
    echo "🐧 CentOS/RHEL:"
    echo "  sudo yum install epel-release"
    echo "  sudo yum localinstall --nogpgcheck https://download1.rpmfusion.org/free/el/rpmfusion-free-release-7.noarch.rpm"
    echo "  sudo yum localinstall --nogpgcheck https://download1.rpmfusion.org/nonfree/el/rpmfusion-nonfree-release-7.noarch.rpm"
    echo "  sudo yum install ffmpeg ImageMagick pandoc"
    echo ""
    echo "🐧 Arch Linux:"
    echo "  sudo pacman -S ffmpeg imagemagick pandoc"
    echo ""
    echo "🐧 openSUSE:"
    echo "  sudo zypper addrepo https://download.opensuse.org/repositories/multimedia/openSUSE_Leap_\$(rpm -E %suse_version).repo"
    echo "  sudo zypper install ffmpeg ImageMagick pandoc"
    echo ""
    echo "🍎 macOS:"
    echo "  brew install ffmpeg imagemagick pandoc"
    echo ""
    echo "🪟 Windows:"
    echo "  Download and install from official websites:"
    echo "  - FFmpeg: https://ffmpeg.org/download.html"
    echo "  - ImageMagick: https://imagemagick.org/script/download.php"
    echo "  - Pandoc: https://pandoc.org/installing.html"
    echo ""
}

# Build for current platform
build_current() {
    print_status "Building for current platform..."
    
    # Clean previous builds
    rm -rf "$BUILD_DIR"
    mkdir -p "$BUILD_DIR"
    
    # Build CLI
    print_status "Building CLI..."
    if go build -ldflags="-s -w" -o "$BUILD_DIR/goverter-cli" ./cmd/cli; then
        print_success "CLI built successfully"
    else
        print_error "CLI build failed"
        exit 1
    fi
    
    # Build GUI (may fail on some systems without GUI libraries)
    print_status "Building GUI..."
    if go build -ldflags="-s -w" -o "$BUILD_DIR/goverter-gui" ./cmd/gui 2>/dev/null; then
        print_success "GUI built successfully"
    else
        print_warning "GUI build failed (this is normal on some systems)"
        print_status "GUI requires graphical libraries. CLI will still work."
        
        # Create a simple GUI fallback
        create_simple_gui
    fi
    
    # Make executables
    chmod +x "$BUILD_DIR/goverter-cli" 2>/dev/null || true
    chmod +x "$BUILD_DIR/goverter-gui" 2>/dev/null || true
}

# Create simple GUI fallback
create_simple_gui() {
    print_status "Creating simple GUI fallback..."
    
    cat > "$BUILD_DIR/goverter-gui" << 'EOF'
#!/bin/bash
echo "Goverter GUI"
echo "=============="
echo ""
echo "The full GUI requires graphical libraries."
echo "Please use the CLI for full functionality:"
echo ""
echo "  ./goverter-cli --help"
echo ""
echo "Or install dependencies:"
echo "  Ubuntu/Debian: sudo apt install ffmpeg imagemagick pandoc"
echo "  macOS: brew install ffmpeg imagemagick pandoc"
EOF
    
    chmod +x "$BUILD_DIR/goverter-gui"
    print_success "Simple GUI fallback created"
}

# Run tests if available
run_tests() {
    if [[ -d "test" ]] || [[ -f "*_test.go" ]]; then
        print_status "Running tests..."
        if go test ./...; then
            print_success "All tests passed!"
        else
            print_warning "Some tests failed"
        fi
    else
        print_status "No tests found, skipping..."
    fi
}

# Show build summary
show_summary() {
    echo ""
    print_success "Build completed!"
    echo ""
    echo "Built binaries:"
    ls -la "$BUILD_DIR"/
    echo ""
    echo "Usage:"
    echo "  CLI: $BUILD_DIR/goverter-cli --help"
    echo "  GUI: $BUILD_DIR/goverter-gui"
    echo ""
    
    if command_exists "$BUILD_DIR/goverter-cli"; then
        echo "Quick test:"
        "$BUILD_DIR/goverter-cli" --version 2>/dev/null || echo "  (version info not available)"
    fi
}

# Main build function
main() {
    echo "Goverter Build Script v$VERSION"
    echo "==========================="
    echo ""
    
    # Check if we're in the right directory
    if [[ ! -f "go.mod" ]] || [[ ! -d "cmd" ]]; then
        print_error "Please run this script from the project root directory"
        exit 1
    fi
    
    # Check dependencies
    check_go_deps
    check_system_deps
    
    # Try auto-install if requested
    if [[ "$1" == "--install-deps" ]]; then
        try_auto_install
        if [[ $? -eq 0 ]]; then
            print_success "Dependencies installed successfully!"
        fi
        check_system_deps  # Re-check after installation
    fi
    
    # Build
    build_current
    
    # Run tests
    if [[ "$1" != "--skip-tests" ]]; then
        run_tests
    fi
    
    # Show summary
    show_summary
}

# Help function
show_help() {
    echo "Goverter Build Script"
    echo "Usage: $0 [options]"
    echo ""
    echo "Options:"
    echo "  --help, -h         Show this help message"
    echo "  --install-deps      Try to automatically install missing dependencies"
    echo "  --skip-tests        Skip running tests"
    echo "  --clean             Clean build artifacts before building"
    echo ""
    echo "Examples:"
    echo "  $0                  # Build with dependency checks"
    echo "  $0 --install-deps   # Try to install missing dependencies"
    echo "  $0 --skip-tests     # Build without running tests"
    echo ""
    echo "🐧 Supported Linux distributions:"
    echo "  • Ubuntu/Debian/Pop!_OS (APT)"
    echo "  • Fedora (DNF)"
    echo "  • CentOS/RHEL (YUM)"
    echo "  • Arch Linux (Pacman)"
    echo "  • openSUSE (Zypper)"
    echo ""
}

# Clean build artifacts
clean_build() {
    print_status "Cleaning build artifacts..."
    rm -rf "$BUILD_DIR"
    rm -f goverter-cli goverter-gui
    go clean -cache
    print_success "Clean completed"
}

# Parse command line arguments
case "$1" in
    --help|-h)
        show_help
        exit 0
        ;;
    --clean)
        clean_build
        exit 0
        ;;
esac

# Run main build
main "$@"