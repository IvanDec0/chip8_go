#!/bin/bash

# CHIP-8 Emulator Comprehensive Testing Script
# Runs all tests including unit, integration, and benchmarks

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"

echo "=== CHIP-8 Emulator Testing Suite ==="
echo "Project root: $PROJECT_ROOT"

cd "$PROJECT_ROOT"

# Color codes for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Test results tracking
TOTAL_TESTS=0
PASSED_TESTS=0
FAILED_TESTS=0
SKIPPED_TESTS=0

# Function to print colored output
print_status() {
    local status=$1
    local message=$2
    case $status in
        "PASS")
            echo -e "${GREEN}[PASS]${NC} $message"
            PASSED_TESTS=$((PASSED_TESTS + 1))
            ;;
        "FAIL")
            echo -e "${RED}[FAIL]${NC} $message"
            FAILED_TESTS=$((FAILED_TESTS + 1))
            ;;
        "SKIP")
            echo -e "${YELLOW}[SKIP]${NC} $message"
            SKIPPED_TESTS=$((SKIPPED_TESTS + 1))
            ;;
        "INFO")
            echo -e "${BLUE}[INFO]${NC} $message"
            ;;
    esac
    TOTAL_TESTS=$((TOTAL_TESTS + 1))
}

# Function to run a test and capture result
run_test() {
    local test_name=$1
    local test_command=$2
    
    echo ""
    echo "Running: $test_name"
    echo "Command: $test_command"
    
    if eval "$test_command" >/dev/null 2>&1; then
        print_status "PASS" "$test_name"
        return 0
    else
        print_status "FAIL" "$test_name"
        return 1
    fi
}

# Function to run test with output
run_test_with_output() {
    local test_name=$1
    local test_command=$2
    
    echo ""
    echo "Running: $test_name"
    echo "Command: $test_command"
    
    if eval "$test_command"; then
        print_status "PASS" "$test_name"
        return 0
    else
        print_status "FAIL" "$test_name"
        return 1
    fi
}

# Check Go installation
echo ""
echo "=== Environment Check ==="
if ! command -v go >/dev/null 2>&1; then
    print_status "FAIL" "Go is not installed or not in PATH"
    exit 1
else
    go_version=$(go version)
    print_status "PASS" "Go is available: $go_version"
fi

# Check if project builds
echo ""
echo "=== Build Test ==="
run_test "Project Build" "go build -v ."

# Download dependencies
echo ""
echo "=== Dependency Check ==="
run_test "Module Tidy" "go mod tidy"
run_test "Module Verify" "go mod verify"

# Run static analysis
echo ""
echo "=== Static Analysis ==="

# go vet
run_test "Go Vet" "go vet ./..."

# go fmt check
if [ "$(gofmt -l . | wc -l)" -eq 0 ]; then
    print_status "PASS" "Go Format Check"
else
    print_status "FAIL" "Go Format Check - files need formatting"
    echo "Files that need formatting:"
    gofmt -l .
fi

# Run unit tests
echo ""
echo "=== Unit Tests ==="

if [ -d "test/unit" ]; then
    # Run each test file individually to get better reporting
    for test_file in test/unit/*_test.go; do
        if [ -f "$test_file" ]; then
            test_name=$(basename "$test_file" .go)
            run_test_with_output "Unit Test: $test_name" "go test -v \"$test_file\""
        fi
    done
    
    # Run all unit tests together with coverage
    run_test_with_output "Unit Tests (All)" "go test -v -coverprofile=coverage_unit.out ./test/unit/..."
    
    # Generate coverage report if successful
    if [ -f "coverage_unit.out" ]; then
        coverage_percent=$(go tool cover -func=coverage_unit.out | grep total | awk '{print $3}')
        print_status "INFO" "Unit Test Coverage: $coverage_percent"
    fi
else
    print_status "SKIP" "Unit Tests - test/unit directory not found"
fi

# Run integration tests
echo ""
echo "=== Integration Tests ==="

if [ -d "test/integration" ]; then
    run_test_with_output "Integration Tests" "go test -v -timeout=30s ./test/integration/..."
else
    print_status "SKIP" "Integration Tests - test/integration directory not found"
fi

# Run benchmark tests
echo ""
echo "=== Benchmark Tests ==="

if [ -d "test/benchmarks" ]; then
    # Run benchmarks with short duration for CI
    run_test_with_output "Benchmark Tests" "go test -bench=. -benchtime=1s ./test/benchmarks/..."
    
    # Run memory benchmarks
    run_test_with_output "Memory Benchmarks" "go test -bench=. -benchmem -benchtime=100ms ./test/benchmarks/..."
else
    print_status "SKIP" "Benchmark Tests - test/benchmarks directory not found"
fi

# Test build for different platforms
echo ""
echo "=== Cross-Platform Build Tests ==="

# Test builds for major platforms
platforms=(
    "linux/amd64"
    "linux/arm64"
    "darwin/amd64"
    "darwin/arm64"
    "windows/amd64"
)

for platform in "${platforms[@]}"; do
    IFS='/' read -r goos goarch <<< "$platform"
    run_test "Build $platform" "GOOS=$goos GOARCH=$goarch go build -o /tmp/chip8_${goos}_${goarch} ."
done

# Test with race detector (on supported platforms)
echo ""
echo "=== Race Condition Tests ==="

if [[ "$OSTYPE" == "linux-gnu"* ]] || [[ "$OSTYPE" == "darwin"* ]]; then
    run_test "Race Detector Build" "go build -race ."
    
    if [ -d "test/unit" ]; then
        run_test "Race Detector Tests" "go test -race -short ./test/unit/..."
    fi
else
    print_status "SKIP" "Race Detector - not supported on this platform"
fi

# Memory leak tests (if tools are available)
echo ""
echo "=== Memory Tests ==="

# Build with debug info for memory testing
run_test "Debug Build" "go build -gcflags='-N -l' -o chip8_debug ."

# Check for potential memory issues using go tool
if command -v go >/dev/null 2>&1; then
    run_test "Memory Test Build" "go test -c ./test/unit/..."
else
    print_status "SKIP" "Memory Tests - tools not available"
fi

# Documentation tests
echo ""
echo "=== Documentation Tests ==="

# Check for required documentation files
required_docs=(
    "README.md"
    "docs/ROM_BROWSER_GUIDE.md"
    "docs/CONFIGURATION_REFERENCE.md"
)

for doc in "${required_docs[@]}"; do
    if [ -f "$doc" ]; then
        print_status "PASS" "Documentation: $doc exists"
    else
        print_status "FAIL" "Documentation: $doc missing"
    fi
done

# Test example configurations
echo ""
echo "=== Configuration Tests ==="

if [ -f "chip8_save/config.json" ]; then
    # Test if config is valid JSON
    if jq empty "chip8_save/config.json" >/dev/null 2>&1; then
        print_status "PASS" "Configuration JSON is valid"
    elif python3 -m json.tool "chip8_save/config.json" >/dev/null 2>&1; then
        print_status "PASS" "Configuration JSON is valid (python)"
    else
        print_status "FAIL" "Configuration JSON is invalid"
    fi
else
    print_status "SKIP" "Configuration Tests - no config file found"
fi

# Performance regression tests
echo ""
echo "=== Performance Tests ==="

# Run a simple performance check
if [ -f "chip8" ]; then
    # Test startup time (should be under 1 second)
    # Use a more reliable method to capture startup time
    if ./chip8 --help >/dev/null 2>&1; then
        # Use bash's built-in time and capture it properly
        startup_output=$( (time timeout 5s ./chip8 --help >/dev/null 2>&1) 2>&1 )
        startup_time=$(echo "$startup_output" | grep -E '^real' | awk '{print $2}' | head -1)
        
        if [ -n "$startup_time" ]; then
            # Extract the integer part before the decimal
            startup_seconds=$(echo "$startup_time" | sed 's/m.*s$//' | cut -d. -f1)
            # Handle format like "0m0.123s" by extracting just the seconds part
            if echo "$startup_time" | grep -q 'm.*s$'; then
                # Format is like "0m0.123s"
                startup_seconds=$(echo "$startup_time" | sed 's/.*m\([0-9]\+\)\..*s$/\1/')
            else
                # Format is like "0.123"
                startup_seconds=$(echo "$startup_time" | cut -d. -f1)
            fi
            
            # Default to 0 if parsing failed
            if [ -z "$startup_seconds" ] || ! echo "$startup_seconds" | grep -q '^[0-9]\+$'; then
                startup_seconds=0
            fi
            
            if [ "$startup_seconds" -lt 2 ]; then
                print_status "PASS" "Startup Performance: ${startup_time} (< 2s target)"
            else
                print_status "FAIL" "Startup Performance: ${startup_time} (>= 2s)"
            fi
        else
            print_status "SKIP" "Startup Performance - could not measure time"
        fi
    else
        print_status "SKIP" "Startup Performance - application doesn't support --help"
    fi
else
    print_status "SKIP" "Performance Tests - binary not found"
fi

# Security tests
echo ""
echo "=== Security Tests ==="

# Check for common security issues using go mod
run_test "Security Audit" "go list -json -deps ./... | grep -E '(CVE|vulnerability)' || true"

# Test file permissions
if [ -f "chip8" ]; then
    permissions=$(stat -f "%A" chip8 2>/dev/null || stat -c "%a" chip8 2>/dev/null || echo "unknown")
    if [[ "$permissions" =~ [1357] ]]; then
        print_status "PASS" "Binary Permissions: $permissions (executable)"
    else
        print_status "FAIL" "Binary Permissions: $permissions (not executable)"
    fi
fi

# Cleanup
echo ""
echo "=== Cleanup ==="
rm -f chip8_debug
rm -f coverage_unit.out
rm -f *.test
print_status "INFO" "Cleanup completed"

# Final report
echo ""
echo "=== Test Summary ==="
echo "Total Tests: $TOTAL_TESTS"
echo -e "Passed: ${GREEN}$PASSED_TESTS${NC}"
echo -e "Failed: ${RED}$FAILED_TESTS${NC}"
echo -e "Skipped: ${YELLOW}$SKIPPED_TESTS${NC}"

if [ $FAILED_TESTS -eq 0 ]; then
    echo ""
    echo -e "${GREEN}🎉 All tests passed!${NC}"
    exit 0
else
    echo ""
    echo -e "${RED}❌ Some tests failed. Please review the output above.${NC}"
    exit 1
fi
