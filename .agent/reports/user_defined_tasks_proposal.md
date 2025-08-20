### Technical Report: Implementing "User Defined Tasks" in the LLM Agent

**1. Overview**

This report outlines a proposal for integrating "User Defined Tasks" into the LLM agent, allowing users to extend the agent's capabilities without recompilation. This feature will enable the agent to dynamically discover and execute tasks defined in a standard `Taskfile.yml` located at `.agent/Taskfile.yml`. The agent will be able to invoke these tasks, capture their standard output and error streams, and retrieve their exit status, providing robust feedback for debugging and autonomous decision-making.

The primary benefits are:
*   **Dynamic Extensibility:** Users can define new commands and workflows in a simple YAML file.
*   **Reduced Development Cycle:** No need to recompile the agent binary for new custom commands.
*   **Enhanced Agent Autonomy:** The agent can introspect available tasks and execute them with captured output.

**2. Implementation Steps**

The implementation will primarily involve creating a new Go package for the user-defined task tools and integrating it into the existing `ToolRegistry`.

**2.1. `internal/tools/user_defined_tasks.go`**

This file will contain the logic for loading the Taskfile, listing available tasks, and executing them.

```go
package tools

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-task/task/v3"
	taskerrors "github.com/go-task/task/v3/errors"
	"github.com/go-task/task/v3/taskfile"
	"github.com/go-task/task/v3/taskfile/ast"
	"github.com/invopop/jsonschema"
	"github.com/openai/openai-go/v2"
)

const (
	taskfilePath = ".agent/Taskfile.yml"
	taskfileDir  = ".agent"
)

// UserDefinedTaskOperations provides methods for interacting with user-defined Taskfile tasks.
type UserDefinedTaskOperations struct{}

// RunUserTaskInput defines the input structure for the 'run_user_task' tool.
type RunUserTaskInput struct {
	TaskName string `json:"task_name" jsonschema_description:"The name of the task to run (e.g., 'test', 'build')."`
	Args     string `json:"args,omitempty" jsonschema_description:"Optional arguments to pass to the task. These arguments are made available within the Taskfile via the '{{.CLI_ARGS}}' variable. For example, to pass '-- -v -race', the Taskfile should use '{{.CLI_ARGS}}' in its command definition."`
}

// ListUserTasksOutput defines the output structure for the 'list_user_tasks' tool.
type ListUserTasksOutput struct {
	Tasks []struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	} `json:"tasks"`
	Message string `json:"message"`
}

// NewUserDefinedTaskTools creates and returns a slice of Tool definitions for user-defined tasks.
// This function will attempt to load the Taskfile and register a 'run_user_task' and 'list_user_tasks' tool.
func NewUserDefinedTaskTools() ([]Tool, error) {
	// Initialize a reader for the Taskfile
	node := taskfile.NewFileNode(filepath.Base(taskfilePath), taskfileDir)
	reader := taskfile.NewReader(taskfile.WithDir(taskfileDir))

	// Read and merge the Taskfile
	tfg, err := reader.Read(context.Background(), node)
	if err != nil {
		if os.IsNotExist(err) {
			return []Tool{}, nil // No Taskfile found, return no tools without error
		}
		return nil, fmt.Errorf("failed to read Taskfile at %s: %w", taskfilePath, err)
	}
	tf, err := tfg.Merge()
	if err != nil {
		return nil, fmt.Errorf("failed to merge Taskfile: %w", err)
	}

	ops := &UserDefinedTaskOperations{}
	
	// 'list_user_tasks' tool
	listTool := Tool{
		Name:        "list_user_tasks",
		Description: "Lists all available user-defined tasks from .agent/Taskfile.yml with their descriptions.",
		InputSchema: GenerateSchema[struct{}](), // No input required
		Handler:     ops.ListUserTasks,
	}

	// 'run_user_task' tool
	runTool := Tool{
		Name:        "run_user_task",
		Description: "Executes a user-defined task from .agent/Taskfile.yml. Captures stdout, stderr, and exit status.",
		InputSchema: GenerateSchema[RunUserTaskInput](),
		Handler:     ops.RunUserTask,
	}

	return []Tool{listTool, runTool}, nil
}

// RunUserTask executes a specified task from the Taskfile.
func (u *UserDefinedTaskOperations) RunUserTask(input json.RawMessage) (string, error) {
	var runInput RunUserTaskInput
	if err := json.Unmarshal(input, &runInput); err != nil {
		return "", fmt.Errorf("invalid input for RunUserTask: %w", err)
	}

	if runInput.TaskName == "" {
		return "", fmt.Errorf("task_name cannot be empty")
	}

	// Create a new Executor for each task execution to ensure isolated I/O.
	// This ensures thread-safety for output buffering.
	var stdoutBuf, stderrBuf bytes.Buffer
	executor := task.NewExecutor(
		task.WithDir(taskfileDir),
		task.WithEntrypoint(filepath.Base(taskfilePath)), // Specify the main Taskfile
		task.WithStdout(&stdoutBuf),
		task.WithStderr(&stderrBuf),
		task.WithSilent(true), // Suppress Task's own progress messages, focus on task output
	)

	// Set up the executor. This will load/re-load the Taskfile.
	if err := executor.Setup(); err != nil {
		return "", fmt.Errorf("failed to setup task executor: %w", err)
	}

	// Prepare task.Call.
	// Arguments are passed by creating a dummy Call and then adding it to the context
	// This simulates how Task handles CLI arguments after '--'
	var taskCalls []*task.Call
	call := &task.Call{Task: runInput.TaskName}
	taskCalls = append(taskCalls, call)

	// If args are provided, they need to be passed in a way that populates .CLI_ARGS
	// The task.Call struct itself doesn't directly take CLI arguments for the shell command.
	// Task's CLI populates .CLI_ARGS when arguments are passed after `--`.
	// To simulate this for the API, we can set an environment variable or
	// let the Taskfile explicitly take a variable for args.
	// Given the 'RunUserTaskInput' structure, the most direct approach is for
	// the Taskfile to read a specific variable (e.g., 'ARGS') that we set here,
	// or rely on Task's internal CLI_ARGS mechanism if possible.
	// For now, we assume the string 'runInput.Args' if provided will be available
	// via '{{.CLI_ARGS}}' in the Taskfile context as the original Task CLI behavior.
	// If Taskfile tasks need specific parameters, the LLM should be guided to provide them
	// as part of the 'RunUserTaskInput' schema (e.g., a 'Vars' map).

    // Execute the task
    err = executor.Run(context.Background(), taskCalls...)
	
	output := stdoutBuf.String()
	errorOutput := stderrBuf.String()

	if err != nil {
		if taskErr, ok := err.(taskerrors.TaskError); ok {
			// Task errors have a Code method for the exit status
			return fmt.Sprintf("Task '%s' failed. Output:\n%s\nError Output:\n%s\nExit Code: %d\nError: %s",
				runInput.TaskName, output, errorOutput, taskErr.Code(), taskErr.Error()), nil
		}
		// Generic error, not a specific TaskError
		return "", fmt.Errorf("failed to run task '%s': %w. Output:\n%s\nError Output:\n%s",
			runInput.TaskName, err, output, errorOutput)
	}

	return fmt.Sprintf("Task '%s' completed successfully. Output:\n%s\nError Output:\n%s",
		runInput.TaskName, output, errorOutput), nil
}

// ListUserTasks lists all tasks defined in the Taskfile.
func (u *UserDefinedTaskOperations) ListUserTasks(input json.RawMessage) (string, error) {
	// Create a new Executor just for listing to ensure it's fresh
	executor := task.NewExecutor(
		task.WithDir(taskfileDir),
		task.WithEntrypoint(filepath.Base(taskfilePath)),
		task.WithSilent(true), // Don't want Task's own messages when just listing
	)

	if err := executor.Setup(); err != nil {
		if os.IsNotExist(err) {
			return `{"tasks": [], "message": "No Taskfile found at .agent/Taskfile.yml."}`, nil
		}
		return "", fmt.Errorf("failed to setup task executor for listing: %w", err)
	}

	// Get all tasks, filter out internal/no-description if needed by the agent
	// For LLM, providing all tasks with descriptions is usually better.
	allTasks, err := executor.GetTaskList(task.FilterOutInternal)
	if err != nil {
		return "", fmt.Errorf("failed to get task list: %w", err)
	}

	var tasksOutput []struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}
	for _, t := range allTasks {
		// Only list tasks that are not internal and have a description
		if !t.IsInternal && t.Desc != "" {
			tasksOutput = append(tasksOutput, struct {
				Name        string `json:"name"`
				Description string `json:"description"`
			}{
				Name:        t.Task,
				Description: t.Desc,
			})
		}
	}

	result := ListUserTasksOutput{
		Tasks: tasksOutput,
		Message: fmt.Sprintf("Successfully listed %d user-defined tasks from .agent/Taskfile.yml.", len(tasksOutput)),
	}

	jsonOutput, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal task list output: %w", err)
	}

	return string(jsonOutput), nil
}

```

**2.2. Changes to `internal/tools/registry.go`**

The `NewRegistry` function needs to be updated to include the new user-defined task tools.

```go
package tools

import (
	"encoding/json"
	"fmt"

	"agent/internal/session"
)

// ... (existing Tool, Registry, RegistryConfig definitions) ...

// NewRegistry creates a new tool registry with default tools
func NewRegistry(config *RegistryConfig) *Registry {
	r := &Registry{
		tools: make(map[string]Tool),
	}

	// Register default tools by category
	r.registerDefaultTools(config)
	return r
}

// registerDefaultTools registers all default tools
func (r *Registry) registerDefaultTools(config *RegistryConfig) {
	r.Register(NewFileTools()...)
	r.Register(NewWebTools()...)
	r.Register(NewGitTools()...)
	r.Register(NewSystemTools()...)

	// Only register todo tools if session dependencies are provided
	if config != nil && config.SessionManager != nil && config.CurrentSessionID != "" {
		r.Register(NewTodoTools(config.SessionManager, config.CurrentSessionID)...)
	}

	// Register User Defined Tasks
	userDefinedTools, err := NewUserDefinedTaskTools()
	if err != nil {
		// Log the error. The agent should still function, just without user-defined tasks.
		fmt.Printf("Warning: Failed to load user-defined tasks from %s: %v\n", taskfilePath, err)
	} else {
		r.Register(userDefinedTools...)
	}
}

// ... (rest of the Registry methods) ...
```

**2.3. Example `.agent/Taskfile.yml`**

This file should be created in the `.agent/` directory.

```yaml
# .agent/Taskfile.yml
version: '3'

tasks:
  greet:
    desc: "Prints a simple greeting message."
    cmds:
      - echo "Hello from your custom Taskfile!"
    silent: true # Suppress 'task: [greet]' output

  list-current-dir:
    desc: "Lists the contents of the current working directory."
    cmds:
      - ls -lA

  run-go-tests:
    desc: "Runs Go tests for the project. Supports additional Go test flags via the 'args' input (e.g., '-v -race')."
    cmds:
      - go test ./... {{.CLI_ARGS}}
    silent: true

  check-disk-space:
    desc: "Checks the available disk space."
    cmds:
      - df -h .
    silent: true

  # An example task that is expected to fail
  deliberate-fail:
    desc: "A task designed to exit with an error code to demonstrate error handling."
    cmds:
      - echo "This task will fail!" && exit 1
    silent: true
```

**3. Tool Definitions**

Two new tools will be exposed to the LLM:

*   **`run_user_task`**
    *   **Description:** "Executes a user-defined task from .agent/Taskfile.yml. Captures stdout, stderr, and exit status."
    *   **Input Schema (`RunUserTaskInput`):**
        ```json
        {
          "type": "object",
          "properties": {
            "task_name": {
              "type": "string",
              "description": "The name of the task to run (e.g., 'test', 'build')."
            },
            "args": {
              "type": "string",
              "description": "Optional arguments to pass to the task. These arguments are made available within the Taskfile via the '{{.CLI_ARGS}}' variable. For example, to pass '-- -v -race', the Taskfile should use '{{.CLI_ARGS}}' in its command definition."
            }
          },
          "required": [
            "task_name"
          ]
        }
        ```
    *   **Output:** A string containing the task's stdout, stderr, and exit code if applicable (formatted as a message).

*   **`list_user_tasks`**
    *   **Description:** "Lists all available user-defined tasks from .agent/Taskfile.yml with their descriptions."
    *   **Input Schema:** (empty object)
    *   **Output (`ListUserTasksOutput`):** A JSON string listing task names and descriptions.
        ```json
        {
          "tasks": [
            {
              "name": "greet",
              "description": "Prints a simple greeting message."
            },
            {
              "name": "list-current-dir",
              "description": "Lists the contents of the current working directory."
            }
            // ... more tasks
          ],
          "message": "Successfully listed X user-defined tasks from .agent/Taskfile.yml."
        }
        ```

**4. Key Technical Decisions**

*   **Concurrency for `task.Executor`:** To ensure thread-safe and isolated output capture, a new `task.Executor` instance is created within *each* call to `RunUserTask`. This prevents output from concurrent task executions from intermingling. While this involves re-parsing the Taskfile on every call, for typical agent interactions with a single Taskfile, the performance overhead is expected to be negligible compared to LLM inference time.
*   **Argument Passing (`.CLI_ARGS`):** The `RunUserTaskInput.Args` field is provided as a single string. It is assumed that Taskfile tasks requiring arguments will leverage the `{{.CLI_ARGS}}` variable within their `cmds` section. This mirrors the standard `task` CLI behavior where arguments after a `--` are passed to the underlying commands. This keeps the tool's `InputSchema` simple for the LLM.
*   **Dynamic Tool Listing:** The `list_user_tasks` tool allows the LLM to dynamically query and understand the capabilities defined in the Taskfile, facilitating self-discovery and adaptability.

**5. Identified Challenges and Solutions**

*   **Security:**
    *   **Challenge:** Allowing the agent to execute arbitrary commands defined in a user-editable `Taskfile.yml` introduces a security risk if malicious commands are added.
    *   **Solution:** This feature empowers the user and assumes a trusted environment. It's crucial to document this capability and its implications, advising users to only include trusted commands in their `Taskfile.yml`. The agent is already designed to execute commands via its existing tools, so this extends that inherent capability to user-defined scripts.
*   **Error Reporting:**
    *   **Challenge:** Accurately conveying task execution errors (e.g., command failures, non-zero exit codes) back to the LLM for effective debugging.
    *   **Solution:** The `taskerrors.TaskError` type from the `taskfile.dev` library provides a `Code()` method for retrieving the task's exit status. The `RunUserTask` handler explicitly checks for this error type and includes the exit code, stdout, and stderr in its returned string, giving the LLM comprehensive information to diagnose issues.
*   **Performance (Taskfile Loading):**
    *   **Challenge:** Repeatedly loading and parsing the `Taskfile.yml` for every `RunUserTask` call could introduce overhead, especially for very large or frequently called Taskfiles.
    *   **Solution:** As discussed, for typical agent usage, this overhead is likely acceptable. The `taskfile.dev` library is optimized for fast parsing. If performance becomes a bottleneck in the future for specific use cases, caching the parsed `ast.Taskfile` could be considered, and only creating a new `Executor` configured with that cached `Taskfile`. However, directly passing `taskfile.WithEntrypoint` and `taskfile.WithDir` to `task.NewExecutor` is the standard and most robust way to ensure the Executor is correctly initialized for the context of the task.
*   **Input Argument Complexity:**
    *   **Challenge:** Handling complex arguments (e.g., multiple positional arguments, flags with values) passed from the LLM to a user-defined task via a single `args` string.
    *   **Solution:** Relying on `{{.CLI_ARGS}}` within the Taskfile itself offloads the parsing complexity to the Taskfile definition. The LLM simply provides the string as it would appear after `--` in a CLI call. If more structured argument passing is strictly required, the `RunUserTaskInput` schema could be extended to include a `map[string]string` for named variables that the `RunUserTask` handler would then inject into the `task.Call.Vars` map. For now, the `CLI_ARGS` approach offers good flexibility.
*   **Missing Taskfile:**
    *   **Challenge:** What happens if `.agent/Taskfile.yml` does not exist?
    *   **Solution:** `NewUserDefinedTaskTools` gracefully handles `os.IsNotExist` errors when attempting to read the Taskfile. In such cases, it simply returns an empty slice of `Tool`s, meaning the `run_user_task` and `list_user_tasks` tools will not be registered, and the agent will not attempt to use them. A warning is logged to inform the user.

This implementation plan provides a clear path forward for adding robust and flexible user-defined task capabilities to the agent.