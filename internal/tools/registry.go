package tools

import (
	"encoding/json"
	"fmt"

	"github.com/openai/openai-go/v2"
)

// Tool represents a single tool that can be executed by the agent
type Tool struct {
	Name        string                     `json:"name"`
	Description string                     `json:"description"`
	InputSchema openai.FunctionParameters `json:"input_schema"`
	Handler     func(input json.RawMessage) (string, error)
}

// Registry manages a collection of tools available to the agent
type Registry struct {
	tools map[string]Tool
}

// NewRegistry creates a new tool registry with default tools
func NewRegistry() *Registry {
	r := &Registry{
		tools: make(map[string]Tool),
	}

	// Register default tools by category
	r.registerDefaultTools()
	return r
}

// registerDefaultTools registers all default tools
func (r *Registry) registerDefaultTools() {
	r.Register(NewFileTools()...)
	r.Register(NewWebTools()...)
	r.Register(NewGitTools()...)
	r.Register(NewTodoTools()...)
}

// Register adds one or more tools to the registry
func (r *Registry) Register(tools ...Tool) {
	for _, tool := range tools {
		r.tools[tool.Name] = tool
	}
}

// Execute runs a tool with the given name and input
func (r *Registry) Execute(name string, input json.RawMessage) (string, error) {
	tool, exists := r.tools[name]
	if !exists {
		return "", fmt.Errorf("tool not found: %s", name)
	}

	return tool.Handler(input)
}

// Get retrieves a tool by name
func (r *Registry) Get(name string) (Tool, bool) {
	tool, exists := r.tools[name]
	return tool, exists
}

// List returns all registered tools
func (r *Registry) List() []Tool {
	tools := make([]Tool, 0, len(r.tools))
	for _, tool := range r.tools {
		tools = append(tools, tool)
	}
	return tools
}

// Names returns all registered tool names
func (r *Registry) Names() []string {
	names := make([]string, 0, len(r.tools))
	for name := range r.tools {
		names = append(names, name)
	}
	return names
}

// Count returns the number of registered tools
func (r *Registry) Count() int {
	return len(r.tools)
}