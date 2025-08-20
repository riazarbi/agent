# Go Best Practices for Agent-Friendly Codebases

*Rules for writing Go code optimized for LLM agent interaction patterns*

## Core Principles

1. **Agent Discoverability**: Structure code for grep/glob exploration before expensive file reads
2. **Token Efficiency**: Keep files focused and small enough for complete agent reading
3. **Search-Friendly Patterns**: Use consistent naming enabling precise grep patterns
4. **Focused Responsibility**: Each file has single, clear purpose inferrable from name/location

## Critical Rules (Must Have)

- **File size limits** - 250 lines (warning), 400 lines (error) - forces focused responsibility
- **File names indicate purpose** - `tool_definitions.go`, `agent_runner.go`, `cli_flags.go`
- **Exports have searchable keywords** - comments include domain terms for grep
- **Package boundaries clear** - group related functionality (tools, agent, cli)

## Project Structure

```
project/
├── cmd/                    # Main applications  
├── internal/               # Private code (tools, agent, config)
├── pkg/                    # Public libraries
├── test/                   # Test utilities
└── docs/                   # Documentation
```

## Naming Conventions

- **Packages**: lowercase, descriptive: `user`, `config` (not `util`, `common`)
- **Types**: PascalCase with domain keywords: `ToolDefinition`, `AgentRunner`
- **Functions**: Include searchable terms: `ExecuteTool`, `ValidateInput`
- **Interfaces**: End with `-er`: `Reader`, `ToolExecutor`

## Code Organization

### Separation of Concerns
```go
// Good: Clear boundaries
package user
type Service struct { repo UserRepository }

package http  
type Handler struct { userService *user.Service }

// Bad: Mixed concerns in single file
package main
func CreateUser() { ... }     // Business logic
func ServeHTTP() { ... }      // HTTP handling
```

### Agent-Friendly Structure
```go
// File: tools/definitions.go
// Keywords: tool, definition, schema, validation

// ToolDefinition defines agent tool structure
// Keywords: tool, interface, execution
type ToolDefinition struct {
    Name     string
    Function func(json.RawMessage) (string, error)
}
```

## Error Handling

```go
// Custom error types with context
type ValidationError struct {
    Field   string
    Message string
}

// Always wrap errors with context
func (s *Service) CreateUser(user User) error {
    if err := s.validate(user); err != nil {
        return fmt.Errorf("validating user: %w", err)
    }
    return nil
}

// Sentinel errors for expected conditions
var (
    ErrUserNotFound = errors.New("user not found")
    ErrInvalidInput = errors.New("invalid input")
)
```

## Testing

### Table-Driven Tests
```go
func TestUserService_Create(t *testing.T) {
    tests := []struct {
        name    string
        user    User
        wantErr bool
    }{
        {"valid user", User{Name: "John"}, false},
        {"missing name", User{}, true},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := service.Create(tt.user)
            if tt.wantErr != (err != nil) {
                t.Errorf("unexpected error state")
            }
        })
    }
}
```

### Test Helpers
```go
func setupTestService(t *testing.T) *UserService {
    t.Helper()
    return NewUserService(&mockRepo{}, &mockLogger{})
}
```

## Dependency Injection

```go
// Define interfaces for dependencies
type UserRepository interface {
    Create(User) error
    GetByID(string) (User, error)
}

// Constructor functions
func NewUserService(repo UserRepository) *UserService {
    return &UserService{repo: repo}
}
```

## Performance & Concurrency

```go
// Use context for cancellation
func (s *Service) ProcessWithTimeout(ctx context.Context) error {
    ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
    defer cancel()
    return s.process(ctx)
}

// Protect shared state
type Cache struct {
    mu    sync.RWMutex
    items map[string]Item
}

// Reuse slices/builders
var sb strings.Builder
sb.Grow(estimatedSize)
```

## Security

```go
// Validate all inputs
func validateEmail(email string) error {
    if !emailRegex.MatchString(email) {
        return ValidationError{Field: "email", Message: "invalid"}
    }
    return nil
}

// Secure configuration
type Config struct {
    APIKey string `json:"-"` // Don't marshal secrets
}

// Prevent path traversal
func validatePath(path string) error {
    if strings.Contains(path, "..") {
        return errors.New("path traversal not allowed")
    }
    return nil
}
```

## Documentation

```go
// Package documentation with usage
// Package user provides user management.
// 
// Basic usage:
//   service := user.NewService(repo, logger)
//   err := service.Create(user)
package user

// Function docs with examples
// Create creates a new user with validation.
// Returns ValidationError if invalid, ErrUserExists if duplicate.
func (s *Service) Create(user User) error { ... }
```

## Quality Standards

- Use `gofmt`/`goimports` for formatting
- Run `golangci-lint` for comprehensive linting  
- Maintain 80%+ test coverage
- All exports must have comments with searchable keywords
- Package docs should list main types/functions

## Decision Framework

When adding functionality:
1. Can agent find this via grep without reading unrelated code? → Create focused file
2. Does file exceed 250 lines? → Split by responsibility (warning at 250+, error at 400+)
3. Are domain keywords in names/comments? → Add searchable terms

**Key Question**: Can an agent understand this functionality without reading unrelated code?