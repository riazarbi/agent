package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/go-task/task/v3"
	"github.com/go-task/task/v3/taskfile/ast"
)

// RunCommandInput structure for running a command
type RunCommandInput struct {
	Command string `json:"command" jsonschema_description:"The name of the command to run (e.g., 'build', 'test', 'lint')."`
	Args    string `json:"args,omitempty" jsonschema_description:"Optional arguments to pass to the command. These arguments are made available within the Taskfile via the '{{.CLI_ARGS}}' variable."`
}

// ListCommandsOutput structure for listing commands
type ListCommandsOutput struct {
	Commands []struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	} `json:"commands"`
	Message string `json:"message"` // A message indicating the status of the listing (e.g., "Successfully listed X commands")
}

const TaskfilePath = ".agent/Taskfile.yml"

// taskfileExists checks if the Taskfile exists
func taskfileExists() bool {
	_, err := os.Stat(TaskfilePath)
	return err == nil
}

// NewCommandTools returns command tools if Taskfile exists
func NewCommandTools() []Tool {
	if !taskfileExists() {
		return []Tool{} // Return empty slice if no Taskfile
	}

	return []Tool{
		{
			Name:        "list_commands",
			Description: "Lists all available development commands from .agent/Taskfile.yml. Commands can be added/removed during session, so list regularly to see current options.",
			InputSchema: map[string]any{
				"type":       "object",
				"properties": map[string]any{},
			},
			Handler: handleListCommands,
		},
		{
			Name:        "run_command",
			Description: "Executes a development command from .agent/Taskfile.yml (e.g., build, test, lint) with optional arguments for code verification workflows",
			InputSchema: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"command": map[string]any{
						"type":        "string",
						"description": "The name of the command to run (e.g., 'build', 'test', 'lint').",
					},
					"args": map[string]any{
						"type":        "string",
						"description": "Optional arguments to pass to the command. These arguments are made available within the Taskfile via the '{{.CLI_ARGS}}' variable.",
					},
				},
				"required": []string{"command"},
			},
			Handler: handleRunCommand,
		},
	}
}

// handleListCommands lists all available commands
func handleListCommands(input json.RawMessage) (string, error) {
	// Create a new executor
	executor := task.NewExecutor(
		task.WithDir("."),
		task.WithEntrypoint(TaskfilePath),
		task.WithSilent(true),
	)

	// Setup the executor
	if err := executor.Setup(); err != nil {
		return "", fmt.Errorf("failed to setup executor: %w", err)
	}

	// Get the taskfile
	taskfile := executor.Taskfile
	if taskfile == nil || taskfile.Tasks == nil {
		return "", fmt.Errorf("no tasks found in Taskfile")
	}

	var output ListCommandsOutput
	output.Commands = []struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}{}

	// Iterate through tasks and collect those with descriptions
	for name, task := range taskfile.Tasks.All(nil) {
		// Only include tasks that are not internal and have descriptions
		if task.Desc != "" {
			output.Commands = append(output.Commands, struct {
				Name        string `json:"name"`
				Description string `json:"description"`
			}{
				Name:        name,
				Description: task.Desc,
			})
		}
	}

	output.Message = fmt.Sprintf("Successfully listed %d development commands from %s.", len(output.Commands), TaskfilePath)

	result, err := json.Marshal(output)
	if err != nil {
		return "", fmt.Errorf("failed to marshal task list: %w", err)
	}

	return string(result), nil
}

// handleRunCommand executes a specific command
func handleRunCommand(input json.RawMessage) (string, error) {
	var commandInput RunCommandInput
	if err := json.Unmarshal(input, &commandInput); err != nil {
		return "", fmt.Errorf("failed to parse input: %w", err)
	}

	if commandInput.Command == "" {
		return "", fmt.Errorf("command cannot be empty")
	}

	// Create a new executor for isolation
	executor := task.NewExecutor(
		task.WithDir("."),
		task.WithEntrypoint(TaskfilePath),
		task.WithSilent(true),
	)

	// Setup the executor
	if err := executor.Setup(); err != nil {
		return "", fmt.Errorf("failed to setup executor: %w", err)
	}

	// Check if task exists
	taskfile := executor.Taskfile
	if taskfile == nil || taskfile.Tasks == nil {
		return "", fmt.Errorf("no tasks found in Taskfile")
	}

	_, exists := taskfile.Tasks.Get(commandInput.Command)
	if !exists {
		return "", fmt.Errorf("command '%s' not found", commandInput.Command)
	}

	// Prepare the task call
	call := &task.Call{
		Task: commandInput.Command,
	}

	// Add CLI_ARGS variable if args provided
	if commandInput.Args != "" {
		vars := ast.NewVars()
		vars.Set("CLI_ARGS", ast.Var{
			Value: commandInput.Args,
		})
		call.Vars = vars
	}

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	// Execute the task
	err := executor.Run(ctx, call)

	// Format the response based on success/failure
	if err != nil {
		return fmt.Sprintf("Command '%s' failed. Output:\n\nError Output:\n\nExit Code: 1\nError: %s", 
			commandInput.Command, err.Error()), nil
	}

	return fmt.Sprintf("Command '%s' completed successfully. Output:\n\nError Output:\n", 
		commandInput.Command), nil
}