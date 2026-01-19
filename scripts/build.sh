#!/bin/bash

# Delta Tool Build Script
# This script builds the delta-tool application for the current platform

set -e

echo "========================================"
echo "  Building Delta Tool"
echo "========================================"
echo ""

# Check if wails is installed
if ! command -v wails &> /dev/null; then
    echo "Error: wails is not installed"
    echo "Please install wails first:"
    echo "  go install github.com/wailsapp/wails/v2/cmd/wails@latest"
    exit 1
fi

# Build the application
echo "Building application..."
wails build

echo ""
echo "========================================"
echo "  Build Complete!"
echo "========================================"
echo "Binary location: build/bin/"
echo ""
