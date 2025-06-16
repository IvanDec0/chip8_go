#!/bin/bash

# CHIP-8 Emulator Cross-Platform Build Script
# Builds binaries for Linux, Windows, and macOS

set -euo pipefail

# Configuration
APP_NAME="chip8-emulator"
VERSION=$(git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_DIR="build"
DIST_DIR="dist"

# Color output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Logging functions
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Function to check if required tools are installed
check_dependencies() {
    log_info "Checking build dependencies..."
    
    if ! command -v go &> /dev/null; then
        log_error "Go is not installed or not in PATH"
        exit 1
    fi
    
    GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
    log_info "Go version: $GO_VERSION"
    
    if ! command -v git &> /dev/null; then
        log_warning "Git is not installed - version will be 'dev'"
    fi
    
    # Check for SDL2 development libraries
    if ! pkg-config --exists sdl2; then
        log_warning "SDL2 development libraries not found - some features may not work"
    fi
    
    log_success "Dependencies check complete"
}

# Function to prepare build environment
prepare_build() {
    log_info "Preparing build environment..."
    
    # Clean previous builds
    rm -rf "$BUILD_DIR" "$DIST_DIR"
    mkdir -p "$BUILD_DIR" "$DIST_DIR"
    
    # Get module dependencies
    log_info "Downloading Go modules..."
    go mod download
    go mod tidy
    
    log_success "Build environment prepared"
}

# Function to run tests before building
run_tests() {
    log_info "Running tests..."
    
    # Run unit tests
    go test ./test/unit/... -v || {
        log_error "Unit tests failed"
        return 1
    }
    
    # Run integration tests (if they exist and are working)
    if [ -d "test/integration" ]; then
        log_info "Running integration tests..."
        go test ./test/integration/... -v || {
            log_warning "Integration tests failed - continuing with build"
        }
    fi
    
    # Run benchmarks (quick mode)
    if [ -d "test/benchmarks" ]; then
        log_info "Running benchmark tests..."
        go test ./test/benchmarks/... -bench=. -benchtime=1s || {
            log_warning "Benchmark tests failed - continuing with build"
        }
    fi
    
    log_success "Tests completed"
}

# Function to build for a specific platform
build_platform() {
    local GOOS=$1
    local GOARCH=$2
    local SUFFIX=$3
    
    log_info "Building for $GOOS/$GOARCH..."
    
    local OUTPUT_NAME="${APP_NAME}_${VERSION}_${GOOS}_${GOARCH}${SUFFIX}"
    local OUTPUT_PATH="$BUILD_DIR/$OUTPUT_NAME"
    
    # Set build environment
    export GOOS GOARCH
    export CGO_ENABLED=1  # Required for SDL2
    
    # Platform-specific CGO configuration
    case "$GOOS" in
        "windows")
            if [ "$GOARCH" = "amd64" ]; then
                export CC=x86_64-w64-mingw32-gcc
                export CXX=x86_64-w64-mingw32-g++
            else
                export CC=i686-w64-mingw32-gcc
                export CXX=i686-w64-mingw32-g++
            fi
            ;;
        "darwin")
            export CC=clang
            export CXX=clang++
            ;;
        "linux")
            # Use default system compiler
            unset CC CXX
            ;;
    esac
    
    # Build flags
    local LDFLAGS="-X main.version=$VERSION -X main.buildTime=$(date -u +%Y-%m-%dT%H:%M:%SZ)"
    
    # Build the binary
    go build \
        -ldflags="$LDFLAGS" \
        -o "$OUTPUT_PATH" \
        . || {
        log_error "Build failed for $GOOS/$GOARCH"
        return 1
    }
    
    # Create distribution package
    local DIST_NAME="${APP_NAME}_${VERSION}_${GOOS}_${GOARCH}"
    local DIST_PATH="$DIST_DIR/$DIST_NAME"
    
    mkdir -p "$DIST_PATH"
    
    # Copy binary
    cp "$OUTPUT_PATH" "$DIST_PATH/"
    
    # Copy additional files
    cp README.md "$DIST_PATH/" 2>/dev/null || true
    cp LICENSE "$DIST_PATH/" 2>/dev/null || true
    
    # Copy documentation
    if [ -d "docs" ]; then
        cp -r docs "$DIST_PATH/"
    fi
    
    # Copy sample ROMs (if they exist)
    if [ -d "roms" ]; then
        mkdir -p "$DIST_PATH/roms"
        cp -r roms/demos "$DIST_PATH/roms/" 2>/dev/null || true
    fi
    
    # Create platform-specific package
    case "$GOOS" in
        "windows")
            # Create ZIP for Windows
            (cd "$DIST_DIR" && zip -r "${DIST_NAME}.zip" "$DIST_NAME/")
            ;;
        *)
            # Create tar.gz for Unix-like systems
            (cd "$DIST_DIR" && tar -czf "${DIST_NAME}.tar.gz" "$DIST_NAME/")
            ;;
    esac
    
    # Clean up directory (keep the archive)
    rm -rf "$DIST_PATH"
    
    log_success "Built $OUTPUT_NAME"
}

# Function to build all platforms
build_all() {
    log_info "Building for all platforms..."
    
    # Build for common platforms
    build_platform "linux" "amd64" ""
    build_platform "linux" "386" ""
    build_platform "darwin" "amd64" ""
    build_platform "darwin" "arm64" ""
    build_platform "windows" "amd64" ".exe"
    build_platform "windows" "386" ".exe"
    
    log_success "All builds completed"
}

# Function to build for current platform only
build_current() {
    log_info "Building for current platform..."
    
    local CURRENT_OS=$(go env GOOS)
    local CURRENT_ARCH=$(go env GOARCH)
    local SUFFIX=""
    
    if [ "$CURRENT_OS" = "windows" ]; then
        SUFFIX=".exe"
    fi
    
    build_platform "$CURRENT_OS" "$CURRENT_ARCH" "$SUFFIX"
    
    log_success "Current platform build completed"
}

# Function to create a development build
build_dev() {
    log_info "Creating development build..."
    
    # Quick build for development
    go build -o "chip8-dev" . || {
        log_error "Development build failed"
        exit 1
    }
    
    log_success "Development build created: ./chip8-dev"
}

# Function to print build information
print_build_info() {
    log_info "Build Information:"
    echo "  Version: $VERSION"
    echo "  Go version: $(go version | awk '{print $3}')"
    echo "  Build time: $(date)"
    echo "  Build directory: $BUILD_DIR"
    echo "  Distribution directory: $DIST_DIR"
}

# Function to clean build artifacts
clean() {
    log_info "Cleaning build artifacts..."
    rm -rf "$BUILD_DIR" "$DIST_DIR"
    rm -f "chip8-dev"
    log_success "Clean completed"
}

# Function to show usage
usage() {
    echo "Usage: $0 [COMMAND]"
    echo ""
    echo "Commands:"
    echo "  all        Build for all supported platforms"
    echo "  current    Build for current platform only"
    echo "  dev        Create development build"
    echo "  test       Run tests only"
    echo "  clean      Clean build artifacts"
    echo "  info       Show build information"
    echo "  help       Show this help message"
    echo ""
    echo "Environment variables:"
    echo "  SKIP_TESTS=1    Skip running tests before build"
    echo "  VERBOSE=1       Enable verbose output"
}

# Main script logic
main() {
    local COMMAND=${1:-"current"}
    
    case "$COMMAND" in
        "all")
            check_dependencies
            prepare_build
            if [ "${SKIP_TESTS:-0}" != "1" ]; then
                run_tests
            fi
            build_all
            print_build_info
            ;;
        "current")
            check_dependencies
            prepare_build
            if [ "${SKIP_TESTS:-0}" != "1" ]; then
                run_tests
            fi
            build_current
            print_build_info
            ;;
        "dev")
            check_dependencies
            build_dev
            ;;
        "test")
            check_dependencies
            run_tests
            ;;
        "clean")
            clean
            ;;
        "info")
            print_build_info
            ;;
        "help"|"-h"|"--help")
            usage
            ;;
        *)
            log_error "Unknown command: $COMMAND"
            usage
            exit 1
            ;;
    esac
}

# Run main function
main "$@"
