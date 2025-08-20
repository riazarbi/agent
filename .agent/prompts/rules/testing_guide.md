# Testing Guide for Code Editing Agent

## Overview

This document defines the testing architecture, standards, and best practices for the code-editing-agent project. All contributors must follow these guidelines, and automated validation ensures compliance.

## Testing Philosophy

### Core Principles
1. **No API Mocking** - Test against real APIs with proper skip logic
2. **Real Environment Testing** - Use actual files, networks, and dependencies
3. **Comprehensive Coverage** - Unit, integration, and end-to-end testing
4. **Isolation & Cleanup** - Tests must not affect each other
5. **Fast Feedback** - Tests should run quickly and provide clear results

### Testing Pyramid
```
    🔺 E2E Tests (few, slow, high confidence)
   🔺🔺 Integration Tests (some, medium speed)
  🔺🔺🔺 Unit Tests (many, fast, focused)
```

## Repository Structure

### Required Directory Layout
```
project/
├── internal/
│   └── package/
│       ├── source.go
│       └── source_test.go           # Unit tests alongside source
├── test/
│   ├── helpers/                     # Shared test utilities
│   │   ├── agent.go                # Agent test helpers
│   │   ├── api.go                  # API test helpers
│   │   ├── assertions.go           # Custom assertions
│   │   ├── config.go               # Test configuration
│   │   ├── environment.go          # Test environment setup
│   │   └── tempdir.go              # File system helpers
│   ├── fixtures/                    # Test data and constants
│   │   ├── test_data.go            # Static test data
│   │   └── samples/                # Sample files for testing
│   ├── integration/                 # Integration tests
│   │   ├── tools_test.go           # Tool integration tests
│   │   └── agent_test.go           # Agent integration tests
│   ├── e2e/                        # End-to-end tests
│   │   └── workflows_test.go       # Complete user workflows
│   └── api/                        # API-specific tests
│       └── openai_test.go          # Real API tests with skip logic
└── scripts/
    └── test-runner.sh              # Test execution scripts
```

## Testing Categories & Build Tags

### Build Tags Usage
- `unit` - Fast unit tests (default)
- `integration` - Integration tests requiring setup
- `api` - Tests requiring API keys
- `e2e` - End-to-end workflow tests
- `slow` - Long-running tests

### Example Usage
```go
//go:build integration
// +build integration

func TestToolsIntegration(t *testing.T) {
    // Integration test code
}
```

### Running Tests by Category
```bash
go test ./...                           # Unit tests only
go test -tags=integration ./...         # Unit + Integration
go test -tags="integration,api" ./...   # Unit + Integration + API
go test -tags=e2e ./test/e2e/          # End-to-end only
```

## API Testing Standards

### No Mocking Rule
**NEVER mock API calls.** Instead, test against real APIs with proper skip logic.

### API Test Template
```go
func TestRealAPICall(t *testing.T) {
    // Skip if no API key
    apiKey := helpers.GetTestAPIKey(t)
    if apiKey == "" {
        t.Skip("Skipping API test: no API key provided (set AGENT_API_KEY)")
    }
    
    // Test with real API
    client := openai.NewClient(option.WithAPIKey(apiKey))
    // ... actual test
}
```

### Environment Variables for Testing
- `AGENT_API_KEY` or `ANTHROPIC_API_KEY` - API keys for testing
- `SKIP_API_TESTS=true` - Skip all API tests
- `SKIP_SLOW_TESTS=true` - Skip long-running tests
- `TEST_TIMEOUT=30s` - Override test timeouts

## Required Test Utilities

### Must-Have Helper Functions
```go
// test/helpers/api.go
func GetTestAPIKey(t *testing.T) string
func SkipIfNoAPIKey(t *testing.T) string
func SetupRealAgent(t *testing.T, apiKey string) *agent.Agent

// test/helpers/environment.go
func SetupTestEnvironment(t *testing.T) TestEnv
func CreateTestSession(t *testing.T) *session.Session
func InitTestGitRepo(t *testing.T) string

// test/helpers/assertions.go
func AssertToolSuccess(t *testing.T, result string)
func AssertValidJSON(t *testing.T, data string)
func AssertFileExists(t *testing.T, path string)
func AssertFileContent(t *testing.T, path, expected string)
```

## Test Writing Standards

### Unit Test Requirements
1. **Location**: Alongside source files (`package_test.go`)
2. **Naming**: `TestFunctionName_Scenario`
3. **Structure**: Table-driven when testing multiple scenarios
4. **Isolation**: Use `t.Helper()`, `t.Cleanup()`
5. **Coverage**: Test both success and error paths

### Unit Test Template
```go
func TestFunction_Scenario(t *testing.T) {
    tests := []struct {
        name        string
        input       InputType
        setup       func(t *testing.T) // Optional setup
        want        ExpectedType
        wantErr     bool
        errContains string  // Optional error message check
    }{
        {
            name:    "successful case",
            input:   validInput,
            want:    expectedOutput,
            wantErr: false,
        },
        {
            name:        "error case",
            input:       invalidInput,
            wantErr:     true,
            errContains: "expected error message",
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            if tt.setup != nil {
                tt.setup(t)
            }
            
            got, err := FunctionUnderTest(tt.input)
            
            if tt.wantErr {
                assert.Error(t, err)
                if tt.errContains != "" {
                    assert.Contains(t, err.Error(), tt.errContains)
                }
                return
            }
            
            assert.NoError(t, err)
            assert.Equal(t, tt.want, got)
        })
    }
}
```

### Integration Test Requirements
1. **Location**: `test/integration/` directory
2. **Build Tags**: `//go:build integration`
3. **Real Dependencies**: Use actual external services
4. **Environment Setup**: Create isolated test environments
5. **Comprehensive Coverage**: Test complete workflows

### End-to-End Test Requirements
1. **Location**: `test/e2e/` directory  
2. **Build Tags**: `//go:build e2e`
3. **User Workflows**: Test complete user scenarios
4. **Real Data**: Use realistic test data and scenarios

## Test Data Management

### Fixtures Structure
```go
// test/fixtures/test_data.go
package fixtures

const (
    SampleGoCode = `...`
    SampleMarkdown = `...`
    // ... other constants
)

var (
    TodoItems = []session.TodoItem{...}
    TestConfigs = map[string]config.Config{...}
)
```

### File-Based Test Data
```
test/fixtures/samples/
├── code/
│   ├── valid.go
│   └── invalid.go
├── documents/
│   ├── sample.md
│   └── sample.html
└── configs/
    ├── valid-config.json
    └── minimal-config.json
```

## Quality Gates & Validation

### Required Test Coverage
- **Minimum**: 80% line coverage
- **Target**: 90% line coverage
- **Critical Paths**: 100% coverage for core functionality

### Coverage Commands
```bash
# Generate coverage report
go test -coverprofile=coverage.out ./...

# View coverage HTML report
go tool cover -html=coverage.out -o coverage.html

# Coverage by package
go test -coverprofile=coverage.out ./... && go tool cover -func=coverage.out
```

### Performance Benchmarks
```go
func BenchmarkCriticalFunction(b *testing.B) {
    for i := 0; i < b.N; i++ {
        CriticalFunction(testInput)
    }
}
```

## Continuous Integration Requirements

### Test Execution Matrix
```yaml
# Example CI configuration
strategy:
  matrix:
    test-type: [unit, integration, api, e2e]
    go-version: [1.22.x]
    os: [ubuntu-latest, macos-latest]
```

### Required CI Checks
1. **Unit Tests**: Must pass on all commits
2. **Integration Tests**: Must pass on PRs
3. **API Tests**: Run with provided API keys
4. **Coverage Check**: Minimum 80% coverage
5. **Benchmark Regression**: Performance must not degrade >10%

## Debugging & Troubleshooting

### Test Debugging
```bash
# Run specific test with verbose output
go test -v -run TestSpecificTest ./package

# Run tests with race detection
go test -race ./...

# Run tests with detailed output
go test -v -json ./... | jq

# Test with timeout
go test -timeout 30s ./...
```

### Common Issues & Solutions
1. **API Rate Limits**: Implement exponential backoff
2. **Test Flakiness**: Add retry logic for network tests  
3. **Resource Cleanup**: Always use `t.Cleanup()`
4. **Test Isolation**: Avoid global state, use fresh instances

## Validation & Enforcement

### Automated Checks
- **Test Structure**: Validate directory layout
- **Naming Conventions**: Enforce test naming patterns
- **Coverage Requirements**: Fail CI if coverage drops
- **API Mock Detection**: Reject PRs with API mocks
- **Build Tags**: Ensure proper categorization

### Manual Review Checklist
- [ ] Tests cover both success and error cases
- [ ] No API mocking is used
- [ ] Proper test isolation and cleanup
- [ ] Appropriate build tags applied
- [ ] Test data is realistic and comprehensive
- [ ] Performance implications considered

## Examples & Templates

### Complete Test File Example
See: `internal/tools/file_test.go` - Demonstrates all best practices

### Integration Test Example  
See: `test/integration/tools_test.go` (to be implemented)

### API Test Example
See: `test/api/openai_test.go` (to be implemented)

## Migration Guide

### Converting Existing Tests
1. **Add Build Tags**: Categorize existing tests
2. **Remove Mocks**: Replace with real API calls + skip logic
3. **Add Helpers**: Use standardized test utilities
4. **Improve Isolation**: Ensure proper cleanup
5. **Add Coverage**: Fill gaps in test coverage

### Testing New Features
1. **Start with Unit Tests**: Test individual functions
2. **Add Integration Tests**: Test component interactions
3. **Include API Tests**: Test real external dependencies
4. **Add E2E Tests**: Test complete user workflows
5. **Benchmark Critical Paths**: Ensure performance standards

This guide ensures consistent, reliable, and maintainable testing practices across the entire codebase.