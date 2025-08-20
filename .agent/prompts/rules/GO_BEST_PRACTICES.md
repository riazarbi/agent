# Go Best Practices Guide

## Overview
This document establishes coding standards, conventions, and best practices for Go development. These guidelines are designed to be timeless and should remain consistent regardless of project evolution.

## Table of Contents
1. [Project Structure](#project-structure)
2. [Code Organization](#code-organization)
3. [Naming Conventions](#naming-conventions)
4. [Error Handling](#error-handling)
5. [Testing](#testing)
6. [Package Design](#package-design)
7. [Performance](#performance)
8. [Concurrency](#concurrency)
9. [Security](#security)
10. [Documentation](#documentation)

## Project Structure

### Standard Layout
Follow the [Standard Go Project Layout](https://github.com/golang-standards/project-layout):

```
project/
├── cmd/                    # Main applications
│   └── appname/
│       └── main.go
├── internal/               # Private application code
│   ├── app/               # Application logic
│   ├── pkg/               # Internal libraries
│   └── config/            # Configuration
├── pkg/                    # Public library code
├── api/                    # API definitions
├── web/                    # Web application assets
├── configs/               # Configuration files
├── deployments/           # Deployment configurations
├── test/                  # Test utilities and data
├── docs/                  # Documentation
├── scripts/               # Build scripts
├── .github/               # GitHub workflows
├── go.mod
├── go.sum
├── README.md
└── Makefile
```

### Key Principles
- Use `cmd/` for main applications
- Use `internal/` for private code that shouldn't be imported
- Use `pkg/` for public libraries
- Keep `main.go` minimal - only entry point logic

## Code Organization

### 1. Separation of Concerns
Each package should have a single, well-defined responsibility:

```go
// Good: Clear separation
package user
type Service struct { ... }
func (s *Service) Create(user User) error { ... }

package http
type Handler struct { userService *user.Service }
func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) { ... }

// Bad: Mixed concerns
package server
func CreateUser() { ... }           // Business logic
func ServeHTTP() { ... }            // HTTP handling
func ValidateUser() { ... }         // Validation
```

### 2. Dependency Injection
Use constructor functions and interfaces:

```go
// Define interfaces for dependencies
type UserRepository interface {
    Create(user User) error
    GetByID(id string) (User, error)
}

type UserService struct {
    repo UserRepository
    logger Logger
}

// Constructor function
func NewUserService(repo UserRepository, logger Logger) *UserService {
    return &UserService{
        repo:   repo,
        logger: logger,
    }
}
```

### 3. Configuration Management
Centralize configuration handling:

```go
type Config struct {
    Server   ServerConfig   `json:"server"`
    Database DatabaseConfig `json:"database"`
    Logger   LoggerConfig   `json:"logger"`
}

func LoadConfig() (*Config, error) {
    // Load from environment variables, files, etc.
}
```

## Naming Conventions

### Packages
- Use lowercase, single words when possible
- Use short, descriptive names: `user`, `http`, `config`
- Avoid generic names: `util`, `common`, `helper`

### Types and Functions
- Use PascalCase for exported identifiers
- Use camelCase for unexported identifiers
- Use descriptive names: `UserService` not `US`
- Interface names should end with `-er`: `Reader`, `Writer`, `UserCreator`

### Variables
- Use short names for short-lived variables: `i`, `err`, `ctx`
- Use descriptive names for longer-lived variables: `userService`, `httpClient`
- Avoid Hungarian notation

### Constants
- Use PascalCase for exported constants
- Use camelCase for unexported constants
- Group related constants in blocks

```go
const (
    StatusPending = "pending"
    StatusActive  = "active"
    StatusInactive = "inactive"
)
```

## Error Handling

### 1. Error Types
Define custom error types for different categories:

```go
type ValidationError struct {
    Field   string
    Message string
}

func (e ValidationError) Error() string {
    return fmt.Sprintf("validation failed for %s: %s", e.Field, e.Message)
}

// Implement Unwrap for error chaining
func (e ValidationError) Unwrap() error {
    return e.Wrapped
}
```

### 2. Error Wrapping
Always provide context when wrapping errors:

```go
func (s *Service) CreateUser(user User) error {
    if err := s.validate(user); err != nil {
        return fmt.Errorf("validating user: %w", err)
    }
    
    if err := s.repo.Create(user); err != nil {
        return fmt.Errorf("creating user in repository: %w", err)
    }
    
    return nil
}
```

### 3. Sentinel Errors
Use sentinel errors for expected conditions:

```go
var (
    ErrUserNotFound    = errors.New("user not found")
    ErrInvalidInput    = errors.New("invalid input")
    ErrUnauthorized    = errors.New("unauthorized")
)

// Check with errors.Is()
if errors.Is(err, ErrUserNotFound) {
    // Handle specific error
}
```

## Testing

### 1. Test Organization
- Place test files alongside source files: `user.go` → `user_test.go`
- Use `testdata/` directories for test fixtures
- Separate unit, integration, and end-to-end tests

### 2. Test Structure
Use table-driven tests for multiple scenarios:

```go
func TestUserService_Create(t *testing.T) {
    tests := []struct {
        name    string
        user    User
        wantErr bool
        errType error
    }{
        {
            name: "valid user",
            user: User{Name: "John", Email: "john@example.com"},
        },
        {
            name:    "missing name",
            user:    User{Email: "john@example.com"},
            wantErr: true,
            errType: ValidationError{},
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            service := setupTestService(t)
            err := service.Create(tt.user)
            
            if tt.wantErr {
                assert.Error(t, err)
                if tt.errType != nil {
                    assert.IsType(t, tt.errType, err)
                }
                return
            }
            
            assert.NoError(t, err)
        })
    }
}
```

### 3. Mocking
Use interfaces for testable code:

```go
type mockUserRepository struct {
    users map[string]User
    err   error
}

func (m *mockUserRepository) Create(user User) error {
    if m.err != nil {
        return m.err
    }
    m.users[user.ID] = user
    return nil
}
```

### 4. Test Helpers
Create helper functions to reduce boilerplate:

```go
func setupTestService(t *testing.T) *UserService {
    t.Helper()
    
    mockRepo := &mockUserRepository{users: make(map[string]User)}
    mockLogger := &mockLogger{}
    
    return NewUserService(mockRepo, mockLogger)
}
```

## Package Design

### 1. Interface Segregation
Keep interfaces small and focused:

```go
// Good: Small, focused interfaces
type Reader interface {
    Read() ([]byte, error)
}

type Writer interface {
    Write(data []byte) error
}

// Bad: Large interface
type FileHandler interface {
    Read() ([]byte, error)
    Write(data []byte) error
    Delete() error
    Compress() error
    Encrypt() error
}
```

### 2. Package Dependencies
- Avoid circular dependencies
- Depend on interfaces, not concrete types
- Keep dependency direction clear: `app → service → repository`

### 3. Exported API
Be conservative with exports:

```go
// Export only what's necessary
type Service struct {
    // Unexported fields
    repo   repository
    logger logger
}

// Exported methods only
func (s *Service) Create(user User) error { ... }
func (s *Service) GetByID(id string) (User, error) { ... }
```

## Performance

### 1. Memory Management
- Reuse slices and maps when possible
- Use object pools for frequently allocated objects
- Be mindful of slice capacity vs. length

```go
// Reuse slice
func (s *Service) ProcessItems(items []Item) {
    if cap(s.buffer) < len(items) {
        s.buffer = make([]ProcessedItem, len(items))
    }
    s.buffer = s.buffer[:len(items)]
    
    // Process into buffer
}
```

### 2. String Operations
Use `strings.Builder` for string concatenation:

```go
func buildMessage(parts []string) string {
    var sb strings.Builder
    sb.Grow(estimateSize(parts)) // Pre-allocate if size is known
    
    for _, part := range parts {
        sb.WriteString(part)
    }
    return sb.String()
}
```

### 3. Profiling
- Use built-in profiling tools: `go test -bench` and `pprof`
- Profile before optimizing
- Focus on algorithmic improvements first

## Concurrency

### 1. Context Usage
Always use context for cancellation and timeouts:

```go
func (s *Service) ProcessWithTimeout(ctx context.Context, data Data) error {
    ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
    defer cancel()
    
    return s.process(ctx, data)
}
```

### 2. Synchronization
Protect shared state appropriately:

```go
// Use sync.RWMutex for read-heavy workloads
type Cache struct {
    mu    sync.RWMutex
    items map[string]Item
}

func (c *Cache) Get(key string) (Item, bool) {
    c.mu.RLock()
    defer c.mu.RUnlock()
    
    item, exists := c.items[key]
    return item, exists
}
```

### 3. Channel Usage
- Use channels for communication, mutexes for protecting state
- Close channels from sender side
- Use buffered channels appropriately

```go
// Worker pool pattern
func (s *Service) ProcessConcurrently(tasks []Task) []Result {
    const numWorkers = 5
    taskChan := make(chan Task, len(tasks))
    resultChan := make(chan Result, len(tasks))
    
    // Start workers
    for i := 0; i < numWorkers; i++ {
        go func() {
            for task := range taskChan {
                result := s.processTask(task)
                resultChan <- result
            }
        }()
    }
    
    // Send tasks
    for _, task := range tasks {
        taskChan <- task
    }
    close(taskChan)
    
    // Collect results
    results := make([]Result, len(tasks))
    for i := range results {
        results[i] = <-resultChan
    }
    
    return results
}
```

## Security

### 1. Input Validation
Validate all inputs at boundaries:

```go
func validateEmail(email string) error {
    if len(email) == 0 {
        return ValidationError{Field: "email", Message: "required"}
    }
    
    if !emailRegex.MatchString(email) {
        return ValidationError{Field: "email", Message: "invalid format"}
    }
    
    return nil
}
```

### 2. Secure Configuration
- Never hardcode secrets
- Use environment variables or secure storage
- Don't log sensitive information

```go
type Config struct {
    APIKey  string `json:"-"` // Don't marshal secrets
    BaseURL string `json:"base_url"`
}

func LoadConfig() Config {
    return Config{
        APIKey:  os.Getenv("API_KEY"),
        BaseURL: getEnvOrDefault("BASE_URL", "https://api.example.com"),
    }
}
```

### 3. Path Validation
Prevent directory traversal attacks:

```go
func validatePath(path string) error {
    if strings.Contains(path, "..") {
        return errors.New("path traversal not allowed")
    }
    
    absPath, err := filepath.Abs(path)
    if err != nil {
        return err
    }
    
    workDir, _ := os.Getwd()
    if !strings.HasPrefix(absPath, workDir) {
        return errors.New("access outside working directory not allowed")
    }
    
    return nil
}
```

## Documentation

### 1. Package Documentation
Document packages with their purpose and usage:

```go
// Package user provides user management functionality.
//
// This package handles user creation, validation, and retrieval.
// It provides a service layer that can be used by HTTP handlers
// and other application components.
//
// Basic usage:
//
//	repo := repository.NewUserRepository(db)
//	logger := log.New(os.Stdout, "", log.LstdFlags)
//	service := user.NewService(repo, logger)
//	
//	err := service.Create(user.User{Name: "John", Email: "john@example.com"})
//
package user
```

### 2. Function Documentation
Document exported functions with their behavior:

```go
// Create creates a new user in the system.
//
// The user must have a valid name and email address. If validation
// fails, a ValidationError is returned. If the user already exists,
// an ErrUserExists error is returned.
//
// Example:
//
//	user := User{Name: "John", Email: "john@example.com"}
//	err := service.Create(user)
//	if err != nil {
//		log.Fatal(err)
//	}
//
func (s *Service) Create(user User) error {
    // Implementation...
}
```

### 3. Code Examples
Provide working examples in documentation:

```go
// Example demonstrates basic service usage
func Example() {
    repo := &mockRepository{}
    logger := log.New(os.Stdout, "", log.LstdFlags)
    service := NewService(repo, logger)
    
    user := User{Name: "John", Email: "john@example.com"}
    err := service.Create(user)
    if err != nil {
        log.Printf("Error creating user: %v", err)
        return
    }
    
    fmt.Println("User created successfully")
    // Output: User created successfully
}
```

## Code Quality Standards

### 1. Formatting
- Use `gofmt` or `goimports` for consistent formatting
- Configure your editor to format on save
- Use `golangci-lint` for comprehensive linting

### 2. Code Review
- Review for correctness, readability, and maintainability
- Ensure error handling is appropriate
- Verify tests cover new functionality
- Check for security vulnerabilities

### 3. Continuous Integration
- Run tests, linting, and formatting checks in CI
- Maintain test coverage above 80%
- Use static analysis tools
- Enforce these standards in pull requests

This guide establishes the foundation for writing high-quality Go code that is maintainable, testable, and secure.