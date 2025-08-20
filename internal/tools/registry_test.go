package tools

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"

	"agent/internal/session"
	"agent/test/helpers"
)

func TestNewRegistry(t *testing.T) {
	cfg := &RegistryConfig{}
	registry := NewRegistry(cfg)

	assert.NotNil(t, registry)
	assert.NotNil(t, registry.tools)

	// Should have registered default tools
	tools := registry.List()
	assert.NotEmpty(t, tools)

	// Should have at least file, web, and git tools
	names := registry.Names()
	assert.Contains(t, names, "read_file")
	assert.Contains(t, names, "edit_file")
	assert.Contains(t, names, "web_fetch")
	assert.Contains(t, names, "git_diff")
}

func TestRegistryWithSessionManager(t *testing.T) {
	// Create a session manager for testing
	sessionDir := helpers.TempDir(t)
	sessionManager := session.NewManager(session.Config{SessionsDir: sessionDir})

	cfg := &RegistryConfig{
		SessionManager:   sessionManager,
		CurrentSessionID: "test-session-123",
	}
	registry := NewRegistry(cfg)

	names := registry.Names()

	// Should have todo tools when session manager is provided
	assert.Contains(t, names, "todowrite")
	assert.Contains(t, names, "todoread")

	// Verify todo tools have valid handlers
	todoWriteTool, exists := registry.Get("todowrite")
	assert.True(t, exists)
	assert.NotNil(t, todoWriteTool.Handler)
}

func TestRegistry_Register(t *testing.T) {
	cfg := &RegistryConfig{}
	registry := NewRegistry(cfg)

	initialCount := registry.Count()

	// Create a test tool
	testTool := Tool{
		Name:        "test_tool",
		Description: "A test tool",
		InputSchema: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"message": map[string]interface{}{
					"type": "string",
				},
			},
		},
		Handler: func(input json.RawMessage) (string, error) {
			return "test result", nil
		},
	}

	registry.Register(testTool)

	// Verify tool was registered
	assert.Equal(t, initialCount+1, registry.Count())

	tool, exists := registry.Get("test_tool")
	assert.True(t, exists)
	assert.Equal(t, "test_tool", tool.Name)
	assert.Equal(t, "A test tool", tool.Description)

	names := registry.Names()
	assert.Contains(t, names, "test_tool")
}

func TestRegistry_RegisterMultiple(t *testing.T) {
	cfg := &RegistryConfig{}
	registry := NewRegistry(cfg)

	initialCount := registry.Count()

	// Create multiple test tools
	tools := []Tool{
		{
			Name:        "tool1",
			Description: "First tool",
			InputSchema: map[string]interface{}{"type": "object"},
			Handler: func(input json.RawMessage) (string, error) {
				return "result1", nil
			},
		},
		{
			Name:        "tool2",
			Description: "Second tool",
			InputSchema: map[string]interface{}{"type": "object"},
			Handler: func(input json.RawMessage) (string, error) {
				return "result2", nil
			},
		},
	}

	registry.Register(tools...)

	// Verify both tools were registered
	assert.Equal(t, initialCount+2, registry.Count())

	tool1, exists1 := registry.Get("tool1")
	assert.True(t, exists1)
	assert.Equal(t, "tool1", tool1.Name)

	tool2, exists2 := registry.Get("tool2")
	assert.True(t, exists2)
	assert.Equal(t, "tool2", tool2.Name)
}

func TestRegistry_Execute(t *testing.T) {
	cfg := &RegistryConfig{}
	registry := NewRegistry(cfg)

	// Register a test tool
	testTool := Tool{
		Name:        "echo_tool",
		Description: "Echoes the input message",
		InputSchema: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"message": map[string]interface{}{
					"type": "string",
				},
			},
		},
		Handler: func(input json.RawMessage) (string, error) {
			var data struct {
				Message string `json:"message"`
			}
			if err := json.Unmarshal(input, &data); err != nil {
				return "", err
			}
			return "Echo: " + data.Message, nil
		},
	}

	registry.Register(testTool)

	// Test successful execution
	input := json.RawMessage(`{"message": "hello world"}`)
	result, err := registry.Execute("echo_tool", input)

	assert.NoError(t, err)
	assert.Equal(t, "Echo: hello world", result)
}

func TestRegistry_Execute_ToolNotFound(t *testing.T) {
	cfg := &RegistryConfig{}
	registry := NewRegistry(cfg)

	input := json.RawMessage(`{"test": "data"}`)
	result, err := registry.Execute("nonexistent_tool", input)

	assert.Error(t, err)
	assert.Empty(t, result)
	assert.Contains(t, err.Error(), "tool not found: nonexistent_tool")
}

func TestRegistry_Execute_HandlerError(t *testing.T) {
	cfg := &RegistryConfig{}
	registry := NewRegistry(cfg)

	// Register a tool that always returns an error
	errorTool := Tool{
		Name:        "error_tool",
		Description: "A tool that errors",
		InputSchema: map[string]interface{}{"type": "object"},
		Handler: func(input json.RawMessage) (string, error) {
			return "", assert.AnError
		},
	}

	registry.Register(errorTool)

	input := json.RawMessage(`{}`)
	result, err := registry.Execute("error_tool", input)

	assert.Error(t, err)
	assert.Empty(t, result)
	assert.Equal(t, assert.AnError, err)
}

func TestRegistry_Get(t *testing.T) {
	cfg := &RegistryConfig{}
	registry := NewRegistry(cfg)

	// Test getting an existing tool (file tools should be registered by default)
	tool, exists := registry.Get("read_file")
	assert.True(t, exists)
	assert.Equal(t, "read_file", tool.Name)
	assert.NotEmpty(t, tool.Description)
	assert.NotNil(t, tool.Handler)

	// Test getting a non-existent tool
	_, exists = registry.Get("nonexistent_tool")
	assert.False(t, exists)
}

func TestRegistry_List(t *testing.T) {
	cfg := &RegistryConfig{}
	registry := NewRegistry(cfg)

	tools := registry.List()
	assert.NotEmpty(t, tools)

	// All tools should have required fields
	for _, tool := range tools {
		assert.NotEmpty(t, tool.Name)
		assert.NotEmpty(t, tool.Description)
		assert.NotNil(t, tool.InputSchema)
		assert.NotNil(t, tool.Handler)
	}

	// Should include default file tools
	var toolNames []string
	for _, tool := range tools {
		toolNames = append(toolNames, tool.Name)
	}

	assert.Contains(t, toolNames, "read_file")
	assert.Contains(t, toolNames, "edit_file")
	assert.Contains(t, toolNames, "delete_file")
}

func TestRegistry_Names(t *testing.T) {
	cfg := &RegistryConfig{}
	registry := NewRegistry(cfg)

	names := registry.Names()
	assert.NotEmpty(t, names)

	// Should include default tool names
	assert.Contains(t, names, "read_file")
	assert.Contains(t, names, "edit_file")
	assert.Contains(t, names, "web_fetch")
	assert.Contains(t, names, "git_diff")
}

func TestRegistry_Count(t *testing.T) {
	cfg := &RegistryConfig{}
	registry := NewRegistry(cfg)

	initialCount := registry.Count()
	assert.Greater(t, initialCount, 0)

	// Add a tool and verify count increases
	testTool := Tool{
		Name:        "count_test_tool",
		Description: "Tool for testing count",
		InputSchema: map[string]interface{}{"type": "object"},
		Handler: func(input json.RawMessage) (string, error) {
			return "test", nil
		},
	}

	registry.Register(testTool)
	assert.Equal(t, initialCount+1, registry.Count())
}

func TestRegistry_RegisterDuplicateTool(t *testing.T) {
	cfg := &RegistryConfig{}
	registry := NewRegistry(cfg)

	// Register the same tool twice
	testTool := Tool{
		Name:        "duplicate_tool",
		Description: "Original description",
		InputSchema: map[string]interface{}{"type": "object"},
		Handler: func(input json.RawMessage) (string, error) {
			return "original", nil
		},
	}

	registry.Register(testTool)
	initialCount := registry.Count()

	// Register tool with same name but different description
	testTool.Description = "Updated description"
	testTool.Handler = func(input json.RawMessage) (string, error) {
		return "updated", nil
	}

	registry.Register(testTool)

	// Count should remain the same (overwrite, not add)
	assert.Equal(t, initialCount, registry.Count())

	// Tool should have updated description
	tool, exists := registry.Get("duplicate_tool")
	assert.True(t, exists)
	assert.Equal(t, "Updated description", tool.Description)

	// Test execution uses the updated handler
	result, err := registry.Execute("duplicate_tool", json.RawMessage(`{}`))
	assert.NoError(t, err)
	assert.Equal(t, "updated", result)
}

func TestRegistry_registerDefaultTools(t *testing.T) {
	// Test that default tools registration works correctly
	cfg := &RegistryConfig{}
	registry := &Registry{
		tools: make(map[string]Tool),
	}

	// Should start empty
	assert.Equal(t, 0, registry.Count())

	// Register default tools
	registry.registerDefaultTools(cfg)

	// Should now have tools
	assert.Greater(t, registry.Count(), 0)

	// Should have specific tool categories
	names := registry.Names()

	// File tools
	assert.Contains(t, names, "read_file")
	assert.Contains(t, names, "edit_file")
	assert.Contains(t, names, "delete_file")

	// Web tools
	assert.Contains(t, names, "web_fetch")
	assert.Contains(t, names, "html_to_markdown")

	// Git tools
	assert.Contains(t, names, "git_diff")
	assert.Contains(t, names, "rg")
	assert.Contains(t, names, "glob")
}

func TestRegistryConfig(t *testing.T) {
	// Test that RegistryConfig can be created with different configurations
	sessionDir := helpers.TempDir(t)
	sessionManager := session.NewManager(session.Config{SessionsDir: sessionDir})

	cfg := RegistryConfig{
		SessionManager: sessionManager,
	}

	assert.Equal(t, sessionManager, cfg.SessionManager)

	// Test that empty config works too
	emptyCfg := RegistryConfig{}
	assert.Nil(t, emptyCfg.SessionManager)
}
