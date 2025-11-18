# Goverter Release Build Script

# This script creates release binaries for all platforms

set -e

VERSION=${1:-"1.0.0"}
BUILD_DIR="build"
RELEASE_DIR="releases"

echo "Building Goverter v$VERSION for all platforms..."

# Clean previous builds
rm -rf "$BUILD_DIR"
rm -rf "$RELEASE_DIR"
mkdir -p "$BUILD_DIR"
mkdir -p "$RELEASE_DIR"

# Build for different platforms
build_platform() {
    local os=$1
    local arch=$2
    local ext=$3
    
    echo "Building for $os/$arch..."
    
    # CLI
    GOOS=$os GOARCH=$arch go build -ldflags="-s -w" -o "$BUILD_DIR/goverter-cli-$os-$arch$ext" ./cmd/cli
    
    # GUI (only for amd64 due to Fyne dependencies)
    if [[ "$arch" == "amd64" ]]; then
        GOOS=$os GOARCH=$arch go build -ldflags="-s -w" -o "$BUILD_DIR/goverter-gui-$os-$arch$ext" ./cmd/gui
    fi
}

# Build for all platforms
build_platform "linux" "amd64" ""
build_platform "linux" "arm64" ""
build_platform "darwin" "amd64" ""
build_platform "darwin" "arm64" ""
build_platform "windows" "amd64" ".exe"

# Create release archives
create_archive() {
    local os=$1
    local arch=$2
    local ext=$3
    local archive_ext=$4
    
    local archive_name="goverter-${VERSION}-${os}-${arch}${ext}"
    local archive_path="$RELEASE_DIR/${archive_name}${archive_ext}"
    
    echo "Creating archive: $archive_path"
    
    mkdir -p "$BUILD_DIR/$archive_name"
    
    # Copy binaries
    cp "$BUILD_DIR/goverter-cli-$os-$arch$ext" "$BUILD_DIR/$archive_name/goverter-cli$ext"
    
    if [[ -f "$BUILD_DIR/goverter-gui-$os-$arch$ext" ]]; then
        cp "$BUILD_DIR/goverter-gui-$os-$arch$ext" "$BUILD_DIR/$archive_name/goverter-gui$ext"
    fi
    
    # Copy additional files
    cp README.md "$BUILD_DIR/$archive_name/"
    cp install.sh "$BUILD_DIR/$archive_name/"
    
    # Create archive
    if [[ "$os" == "windows" ]]; then
        cd "$BUILD_DIR" && zip -r "../$RELEASE_DIR/${archive_name}.zip" "$archive_name"
    else
        cd "$BUILD_DIR" && tar -czf "../$RELEASE_DIR/${archive_name}.tar.gz" "$archive_name"
    fi
}

# Create archives for each platform
create_archive "linux" "amd64" "" ".tar.gz"
create_archive "linux" "arm64" "" ".tar.gz"
create_archive "darwin" "amd64" "" ".tar.gz"
create_archive "darwin" "arm64" "" ".tar.gz"
create_archive "windows" "amd64" ".exe" ".zip"

# Create checksums
echo "Creating checksums..."
cd "$RELEASE_DIR"
sha256sum * > SHA256SUMS

echo ""
echo "Build completed! Release files:"
ls -la "$RELEASE_DIR"

echo ""
echo "To create a GitHub release:"
echo "1. Go to: https://github.com/ezraclintoc/goverter/releases/new"
echo "2. Tag: v$VERSION"
echo "3. Upload all files from $RELEASE_DIR"
echo "4. Include SHA256SUMS for verification"