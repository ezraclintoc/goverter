#!/bin/bash

# Build script for Goverter

echo "Building Goverter..."

# Build CLI
echo "Building CLI..."
go build -o goverter-cli ./cmd/cli
if [ $? -eq 0 ]; then
    echo "✓ CLI built successfully"
else
    echo "✗ CLI build failed"
    exit 1
fi

# Build GUI
echo "Building GUI..."
go build -o goverter-gui ./cmd/gui
if [ $? -eq 0 ]; then
    echo "✓ GUI built successfully"
else
    echo "✗ GUI build failed"
    exit 1
fi

echo "Build complete!"
echo "Usage:"
echo "  ./goverter-cli --help    # CLI interface"
echo "  ./goverter-gui           # GUI interface"