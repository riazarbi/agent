#!/bin/bash
# Test Architecture Validation Script
# Validates that the testing architecture follows established standards

set -euo pipefail

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

ERRORS=0
WARNINGS=0

echo "🧪 Validating Test Architecture..."
echo "=================================="

# Function to report errors
error() {
    echo -e "${RED}❌ ERROR: $1${NC}"
    ((ERRORS++))
}

# Function to report warnings
warning() {
    echo -e "${YELLOW}⚠️  WARNING: $1${NC}"
    ((WARNINGS++))
}

# Function to report success
success() {
    echo -e "${GREEN}✅ $1${NC}"
}

# Validate directory structure
validate_directory_structure() {
    echo "📁 Validating directory structure..."
    
    # Required directories
    required_dirs=(
        "test/helpers"
        "test/fixtures"
        "test/integration"
    )
    
    for dir in "${required_dirs[@]}"; do
        if [ ! -d "$dir" ]; then
            error "Missing required directory: $dir"
        else
            success "Directory exists: $dir"
        fi
    done
    
    # Check for test files alongside source
    find internal/ -name "*.go" -not -name "*_test.go" | while read -r source_file; do
        test_file="${source_file%%.go}_test.go"
        if [ ! -f "$test_file" ]; then
            warning "Missing test file for $source_file"
        fi
    done
}

# Validate test naming conventions
validate_test_naming() {
    echo "📝 Validating test naming conventions..."
    
    # Find all test functions and validate naming
    find . -name "*_test.go" -exec grep -l "func Test" {} \; | while read -r file; do
        # Check for proper Test function naming
        if ! grep -q "func Test[A-Z]" "$file"; then
            warning "Test functions in $file may not follow naming conventions"
        fi
        
        # Check for table-driven test patterns
        if grep -q "tests := \[\]struct" "$file"; then
            success "Table-driven tests found in $file"
        fi
    done
}

# Validate no API mocking
validate_no_api_mocking() {
    echo "🚫 Validating no API mocking..."
    
    # Look for common mocking patterns
    mock_patterns=(
        "mock.*[Cc]lient"
        "Mock[A-Z]"
        "testify/mock"
        "gomock"
        "httptest.NewServer.*api"
    )
    
    for pattern in "${mock_patterns[@]}"; do
        if grep -r "$pattern" --include="*_test.go" .; then
            error "Potential API mocking found (pattern: $pattern). Use real APIs with skip logic instead."
        fi
    done
    
    # Look for proper skip patterns
    if grep -r "t.Skip.*API" --include="*_test.go" .; then
        success "Found proper API test skip logic"
    fi
}

# Validate build tags
validate_build_tags() {
    echo "🏷️  Validating build tags..."
    
    # Check integration tests have build tags
    if [ -d "test/integration" ]; then
        find test/integration -name "*_test.go" | while read -r file; do
            if ! head -5 "$file" | grep -q "//go:build.*integration"; then
                warning "Integration test $file missing build tag"
            else
                success "Build tag found in $file"
            fi
        done
    fi
    
    # Check API tests have build tags
    if [ -d "test/api" ]; then
        find test/api -name "*_test.go" | while read -r file; do
            if ! head -5 "$file" | grep -q "//go:build.*api"; then
                warning "API test $file missing build tag"
            else
                success "Build tag found in $file"
            fi
        done
    fi
}

# Validate test helpers
validate_test_helpers() {
    echo "🔧 Validating test helpers..."
    
    required_helpers=(
        "test/helpers/config.go"
        "test/helpers/tempdir.go"
    )
    
    for helper in "${required_helpers[@]}"; do
        if [ ! -f "$helper" ]; then
            error "Missing required test helper: $helper"
        else
            success "Test helper exists: $helper"
        fi
    done
    
    # Check for t.Helper() usage in test helpers
    find test/helpers -name "*.go" | while read -r file; do
        if grep -q "func.*\*testing.T" "$file" && ! grep -q "t.Helper()" "$file"; then
            warning "Test helper $file missing t.Helper() calls"
        fi
    done
}

# Validate test isolation
validate_test_isolation() {
    echo "🏝️  Validating test isolation..."
    
    # Check for proper cleanup patterns
    cleanup_patterns=(
        "t.Cleanup"
        "defer.*Remove"
        "defer.*Close"
    )
    
    for pattern in "${cleanup_patterns[@]}"; do
        if ! grep -r "$pattern" --include="*_test.go" . >/dev/null; then
            warning "No cleanup pattern '$pattern' found in tests"
        fi
    done
    
    # Check for potential shared state issues
    shared_state_patterns=(
        "var.*=.*make("
        "init()"
    )
    
    for pattern in "${shared_state_patterns[@]}"; do
        if grep -r "$pattern" --include="*_test.go" .; then
            warning "Potential shared state found (pattern: $pattern)"
        fi
    done
}

# Validate test coverage requirements
validate_coverage() {
    echo "📊 Validating test coverage..."
    
    # Run tests with coverage
    if go test -coverprofile=coverage.out ./... >/dev/null 2>&1; then
        coverage=$(go tool cover -func=coverage.out | grep total | awk '{print $3}' | sed 's/%//')
        
        if (( $(echo "$coverage < 80" | bc -l) )); then
            error "Test coverage is ${coverage}%, minimum required is 80%"
        elif (( $(echo "$coverage < 90" | bc -l) )); then
            warning "Test coverage is ${coverage}%, target is 90%"
        else
            success "Test coverage is ${coverage}%"
        fi
        
        rm -f coverage.out
    else
        error "Failed to generate test coverage report"
    fi
}

# Validate test execution
validate_test_execution() {
    echo "🚀 Validating test execution..."
    
    # Test different categories
    test_commands=(
        "go test ./..."
        "go test -tags=integration ./..."
    )
    
    for cmd in "${test_commands[@]}"; do
        echo "Running: $cmd"
        if $cmd >/dev/null 2>&1; then
            success "Tests pass: $cmd"
        else
            error "Tests fail: $cmd"
        fi
    done
}

# Validate dependencies
validate_dependencies() {
    echo "📦 Validating test dependencies..."
    
    # Check for required testing dependencies
    if ! go list -m github.com/stretchr/testify >/dev/null 2>&1; then
        error "Missing required dependency: github.com/stretchr/testify"
    else
        success "testify dependency found"
    fi
}

# Main validation function
main() {
    validate_directory_structure
    validate_test_naming
    validate_no_api_mocking
    validate_build_tags
    validate_test_helpers
    validate_test_isolation
    validate_coverage
    validate_test_execution
    validate_dependencies
    
    echo ""
    echo "=================================="
    echo "🧪 Test Architecture Validation Complete"
    
    if [ $ERRORS -gt 0 ]; then
        echo -e "${RED}❌ Found $ERRORS errors${NC}"
        exit 1
    elif [ $WARNINGS -gt 0 ]; then
        echo -e "${YELLOW}⚠️  Found $WARNINGS warnings${NC}"
        echo -e "${GREEN}✅ No critical issues found${NC}"
        exit 0
    else
        echo -e "${GREEN}✅ All validations passed!${NC}"
        exit 0
    fi
}

# Run main function
main "$@"