package tools

import (
	"encoding/json"
	"fmt"
	"strings"

	"agent/internal/session"
)

// TodoOperations handles todo-related tool operations
type TodoOperations struct {
	sessionManager   session.SessionManager
	currentSessionID string
}

// TodoWriteInput matches the expected JSON structure from main.go
type TodoWriteInput struct {
	TodosJSON string `json:"todos_json" jsonschema_description:"The updated todo list as JSON string containing array of TodoItem objects"`
}

// TodoReadInput for todoread (no parameters needed)
type TodoReadInput struct {
	// No parameters needed - reads from current session
}

// NewTodoTools returns todo operation tools configured with session manager
func NewTodoTools(sessionManager session.SessionManager, currentSessionID string) []Tool {
	ops := &TodoOperations{
		sessionManager:   sessionManager,
		currentSessionID: currentSessionID,
	}

	return []Tool{
		{
			Name:        "todowrite",
			Description: "Create and manage structured task lists for complex multi-step operations within the current session. Each todo requires: 'id', 'content', 'status' (pending/in_progress/completed/cancelled), 'priority' (high/medium/low). Replaces entire todo list. Data persists within the session.",
			InputSchema: GenerateSchema[TodoWriteInput](),
			Handler:     ops.TodoWrite,
		},
		{
			Name:        "todoread",
			Description: "Read the current todo list from session state. Returns structured todos with IDs, content, status, and priority. Data persists within the session.",
			InputSchema: GenerateSchema[TodoReadInput](),
			Handler:     ops.TodoRead,
		},
	}
}

// TodoWrite manages the todo list for the current session using session manager
func (t *TodoOperations) TodoWrite(input json.RawMessage) (string, error) {
	var todoWriteInput TodoWriteInput
	if err := json.Unmarshal(input, &todoWriteInput); err != nil {
		return "", fmt.Errorf("invalid input: %v", err)
	}

	// Parse the JSON string containing the todos array
	var todos []session.TodoItem
	if todoWriteInput.TodosJSON != "" {
		if err := json.Unmarshal([]byte(todoWriteInput.TodosJSON), &todos); err != nil {
			return "", fmt.Errorf("invalid todos JSON: %v", err)
		}
	}

	// Validate and process todos
	var processedTodos []session.TodoItem
	inProgressCount := 0

	for i, todo := range todos {
		// Generate ID if not provided
		if todo.ID == "" {
			todo.ID = t.sessionManager.GenerateTodoID()
		}

		// Validate status
		if !session.ValidateTodoStatus(todo.Status) {
			return "", fmt.Errorf("invalid status '%s' for todo %d. Must be one of: pending, in_progress, completed, cancelled", todo.Status, i+1)
		}

		// Validate priority
		if !session.ValidateTodoPriority(todo.Priority) {
			return "", fmt.Errorf("invalid priority '%s' for todo %d. Must be one of: high, medium, low", todo.Priority, i+1)
		}

		// Count in_progress todos
		if todo.Status == "in_progress" {
			inProgressCount++
		}

		// Validate content is not empty
		if strings.TrimSpace(todo.Content) == "" {
			return "", fmt.Errorf("todo content cannot be empty for todo %d", i+1)
		}

		processedTodos = append(processedTodos, todo)
	}

	// Enforce single in_progress todo rule
	if inProgressCount > 1 {
		return "", fmt.Errorf("only one todo can be in_progress at a time, found %d", inProgressCount)
	}

	// Save todos using session manager
	if err := t.sessionManager.SaveTodos(t.currentSessionID, processedTodos); err != nil {
		return "", fmt.Errorf("failed to save todos: %v", err)
	}

	// Count non-completed todos for title
	nonCompletedCount := 0
	for _, todo := range processedTodos {
		if todo.Status != "completed" && todo.Status != "cancelled" {
			nonCompletedCount++
		}
	}

	// Return result
	result := map[string]interface{}{
		"title":  fmt.Sprintf("Updated todo list with %d active todos", nonCompletedCount),
		"output": fmt.Sprintf("Successfully updated %d todos", len(processedTodos)),
	}

	resultJSON, err := json.Marshal(result)
	if err != nil {
		return "", fmt.Errorf("failed to marshal result: %v", err)
	}

	return string(resultJSON), nil
}

// TodoRead retrieves the current todo list from session state using session manager
func (t *TodoOperations) TodoRead(input json.RawMessage) (string, error) {
	todos, err := t.sessionManager.GetTodos(t.currentSessionID)
	if err != nil {
		return "", fmt.Errorf("failed to get todos: %v", err)
	}

	if len(todos) == 0 {
		result := map[string]interface{}{
			"title":  "0 todos",
			"output": "[]",
		}
		resultJSON, err := json.Marshal(result)
		if err != nil {
			return "", fmt.Errorf("failed to marshal result: %v", err)
		}
		return string(resultJSON), nil
	}

	// Count non-completed todos for title
	nonCompletedCount := 0
	for _, todo := range todos {
		if todo.Status != "completed" && todo.Status != "cancelled" {
			nonCompletedCount++
		}
	}

	// Convert todos to JSON
	todosJSON, err := json.Marshal(todos)
	if err != nil {
		return "", fmt.Errorf("failed to marshal todos: %v", err)
	}

	result := map[string]interface{}{
		"title":  fmt.Sprintf("%d todos", nonCompletedCount),
		"output": string(todosJSON),
	}

	resultJSON, err := json.Marshal(result)
	if err != nil {
		return "", fmt.Errorf("failed to marshal result: %v", err)
	}

	return string(resultJSON), nil
}
