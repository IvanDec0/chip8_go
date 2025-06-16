#!/bin/bash

# CHIP-8 Emulator Release Packaging Script
# This script packages the built binaries for distribution

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
BUILD_DIR="$PROJECT_ROOT/build"
DIST_DIR="$PROJECT_ROOT/dist"
VERSION=$(cat "$PROJECT_ROOT/VERSION" 2>/dev/null || echo "dev")

echo "=== CHIP-8 Emulator Packaging Script ==="
echo "Version: $VERSION"
echo "Project root: $PROJECT_ROOT"

# Clean previous builds
echo "Cleaning previous builds..."
rm -rf "$DIST_DIR"
mkdir -p "$DIST_DIR"

# Ensure build directory exists
if [ ! -d "$BUILD_DIR" ]; then
    echo "Error: Build directory not found. Run build.sh first."
    exit 1
fi

# Package function
package_release() {
    local os=$1
    local arch=$2
    local ext=$3
    
    echo "Packaging $os-$arch..."
    
    local binary_name="chip8"
    if [ "$os" = "windows" ]; then
        binary_name="chip8.exe"
    fi
    
    local package_name="chip8-v${VERSION}-${os}-${arch}"
    local package_dir="$DIST_DIR/$package_name"
    
    # Create package directory
    mkdir -p "$package_dir"
    
    # Copy binary
    if [ -f "$BUILD_DIR/${os}-${arch}/$binary_name" ]; then
        cp "$BUILD_DIR/${os}-${arch}/$binary_name" "$package_dir/"
    else
        echo "Warning: Binary not found for $os-$arch"
        return
    fi
    
    # Copy documentation
    cp "$PROJECT_ROOT/README.md" "$package_dir/"
    cp "$PROJECT_ROOT/README_ES.md" "$package_dir/"
    
    # Copy documentation files
    if [ -d "$PROJECT_ROOT/docs" ]; then
        cp -r "$PROJECT_ROOT/docs" "$package_dir/"
    fi
    
    # Copy sample ROMs (if they exist)
    if [ -d "$PROJECT_ROOT/roms/demos" ]; then
        mkdir -p "$package_dir/roms"
        cp -r "$PROJECT_ROOT/roms/demos" "$package_dir/roms/"
    fi
    
    # Create installation script for Unix systems
    if [ "$os" != "windows" ]; then
        cat > "$package_dir/install.sh" << 'EOF'
#!/bin/bash
# CHIP-8 Emulator Installation Script

INSTALL_DIR="/usr/local/bin"
BINARY_NAME="chip8"

echo "Installing CHIP-8 Emulator..."

# Check for sudo if installing to system directory
if [ ! -w "$INSTALL_DIR" ]; then
    echo "Installing to $INSTALL_DIR (requires sudo)"
    sudo cp "$BINARY_NAME" "$INSTALL_DIR/"
    sudo chmod +x "$INSTALL_DIR/$BINARY_NAME"
else
    cp "$BINARY_NAME" "$INSTALL_DIR/"
    chmod +x "$INSTALL_DIR/$BINARY_NAME"
fi

echo "Installation complete! Run 'chip8' to start the emulator."
EOF
        chmod +x "$package_dir/install.sh"
    fi
    
    # Create Windows installer batch file
    if [ "$os" = "windows" ]; then
        cat > "$package_dir/install.bat" << 'EOF'
@echo off
echo Installing CHIP-8 Emulator...

REM Create program directory
set INSTALL_DIR=%PROGRAMFILES%\CHIP8
if not exist "%INSTALL_DIR%" mkdir "%INSTALL_DIR%"

REM Copy files
copy chip8.exe "%INSTALL_DIR%\"
copy *.md "%INSTALL_DIR%\"
if exist docs xcopy docs "%INSTALL_DIR%\docs\" /E /I
if exist roms xcopy roms "%INSTALL_DIR%\roms\" /E /I

REM Add to PATH (requires admin)
echo Adding to PATH...
setx PATH "%PATH%;%INSTALL_DIR%" /M

echo Installation complete! You may need to restart your command prompt.
echo Run 'chip8' to start the emulator.
pause
EOF
    fi
    
    # Create LICENSE file
    cat > "$package_dir/LICENSE" << 'EOF'
MIT License

Copyright (c) 2024 CHIP-8 Emulator

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
EOF
    
    # Create package archive
    cd "$DIST_DIR"
    if [ "$os" = "windows" ]; then
        # Create ZIP for Windows
        if command -v zip >/dev/null 2>&1; then
            zip -r "${package_name}.zip" "$package_name"
            echo "Created ${package_name}.zip"
        else
            echo "Warning: zip not found, skipping archive creation for Windows package"
        fi
    else
        # Create tar.gz for Unix systems
        tar -czf "${package_name}.tar.gz" "$package_name"
        echo "Created ${package_name}.tar.gz"
    fi
    cd "$PROJECT_ROOT"
}

# Package all available builds
echo "Packaging releases..."

# Check what builds are available
for build_dir in "$BUILD_DIR"/*; do
    if [ -d "$build_dir" ]; then
        build_name=$(basename "$build_dir")
        
        case "$build_name" in
            "linux-amd64")
                package_release "linux" "amd64" ""
                ;;
            "linux-arm64")
                package_release "linux" "arm64" ""
                ;;
            "darwin-amd64")
                package_release "darwin" "amd64" ""
                ;;
            "darwin-arm64")
                package_release "darwin" "arm64" ""
                ;;
            "windows-amd64")
                package_release "windows" "amd64" ".exe"
                ;;
            "windows-arm64")
                package_release "windows" "arm64" ".exe"
                ;;
            *)
                echo "Unknown build: $build_name"
                ;;
        esac
    fi
done

# Create combined source package
echo "Creating source package..."
source_package="chip8-v${VERSION}-source"
mkdir -p "$DIST_DIR/$source_package"

# Copy source files (excluding build artifacts)
rsync -av \
    --exclude="build/" \
    --exclude="dist/" \
    --exclude=".git/" \
    --exclude="*.log" \
    --exclude="chip8_save/" \
    "$PROJECT_ROOT/" "$DIST_DIR/$source_package/"

# Create source archive
cd "$DIST_DIR"
tar -czf "${source_package}.tar.gz" "$source_package"
echo "Created ${source_package}.tar.gz"
cd "$PROJECT_ROOT"

# Create checksums
echo "Creating checksums..."
cd "$DIST_DIR"
find . -name "*.tar.gz" -o -name "*.zip" | while read file; do
    if command -v sha256sum >/dev/null 2>&1; then
        sha256sum "$file" >> "checksums.sha256"
    elif command -v shasum >/dev/null 2>&1; then
        shasum -a 256 "$file" >> "checksums.sha256"
    fi
done

if [ -f "checksums.sha256" ]; then
    echo "Created checksums.sha256"
fi

cd "$PROJECT_ROOT"

# Generate release notes
echo "Generating release notes..."
cat > "$DIST_DIR/RELEASE_NOTES.md" << EOF
# CHIP-8 Emulator v${VERSION} - Release Notes

## Features

- Full CHIP-8 instruction set support
- Advanced ROM browser with metadata support
- Favorites and recent ROMs management
- Configurable key bindings
- Multiple color themes including high contrast accessibility options
- Statistics and usage tracking
- Export functionality for ROM lists and statistics
- Interactive help system
- Cross-platform support (Linux, macOS, Windows)

## Installation

### Linux/macOS
1. Extract the archive: \`tar -xzf chip8-v${VERSION}-<platform>-<arch>.tar.gz\`
2. Run the installer: \`./install.sh\` (optional)
3. Or run directly: \`./chip8\`

### Windows
1. Extract the ZIP file
2. Run \`install.bat\` as Administrator (optional)
3. Or run \`chip8.exe\` directly

## Configuration

The emulator creates a configuration directory:
- Linux: \`~/.config/chip8/\`
- macOS: \`~/Library/Application Support/chip8/\`
- Windows: \`%APPDATA%\chip8\`

## ROM Compatibility

This emulator supports standard CHIP-8 ROMs (.ch8 and .c8 files). Sample demo ROMs are included in the \`roms/demos/\` directory.

## Controls

- Arrow keys: Navigate menus
- Enter: Select ROM/menu item
- Escape: Go back/exit
- F1: Help
- F: Toggle favorites (in ROM browser)
- /: Search
- Tab: Change sorting

For game controls, see the documentation or press F1 for help.

## Documentation

Complete documentation is available in the \`docs/\` directory:
- \`ROM_BROWSER_GUIDE.md\` - User guide for the ROM browser
- \`CONFIGURATION_REFERENCE.md\` - Configuration options reference

## Support

For issues and feature requests, please visit the project repository.

Generated on $(date)
EOF

echo ""
echo "=== Packaging Complete ==="
echo "Packages created in: $DIST_DIR"
echo ""
echo "Available packages:"
ls -la "$DIST_DIR"/*.tar.gz "$DIST_DIR"/*.zip 2>/dev/null || true

echo ""
echo "Release is ready for distribution!"
