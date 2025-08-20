# User Defined Tasks for LLM Agent

This story addresses the need for users to extend the LLM agent's capabilities without requiring recompilation of the agent itself. The main objective is to enable the agent to dynamically discover and execute custom tasks defined in a standard `.agent/Taskfile.yml` file, providing robust feedback on their execution.

## Past Attempts

If this user story has been attempted before, the changes made will appear in the git diff. Our policy is to only make a single commit per user story, so you can always review the git diff to review progress across attempts.

## Requirements

*   The agent must be able to discover and load tasks defined in a `.agent/Taskfile.yml` file.
*   The agent must expose a tool, `list_user_tasks`, that returns a list of available user-defined tasks along with their descriptions.
*   The agent must expose a tool, `run_user_task`, that can execute a specified user-defined task, accepting optional arguments.
*   The `run_user_task` tool must capture and return the standard output (stdout) of the executed task.
*   The `run_user_task` tool must capture and return the standard error (stderr) of the executed task.
*   The `run_user_task` tool must capture and return the exit status/code of the executed task, clearly indicating success or failure.
*   If `.agent/Taskfile.yml` does not exist, the user-defined task tools (`list_user_tasks` and `run_user_task`) should not be registered, and the agent should gracefully handle this absence without error.

## Rules

*   The `Taskfile.yml` will be located at `.agent/Taskfile.yml`.
*   This feature assumes a trusted environment; users should be advised to only include trusted commands in their `Taskfile.yml` due to security implications of arbitrary command execution.
*   Arguments passed to `run_user_task` will be made available within the Taskfile via the `{{.CLI_ARGS}}` variable.
*   The `list_user_tasks` tool should only list tasks that are not internal and have a description.

## Domain

```
// RunUserTaskInput structure for running a task
type RunUserTaskInput struct {
    TaskName string `json:"task_name" jsonschema_description:"The name of the task to run (e.g., 'test', 'build')."
    Args     string `json:"args,omitempty" jsonschema_description:"Optional arguments to pass to the task. These arguments are made available within the Taskfile via the '{{.CLI_ARGS}}' variable."`
}

// ListUserTasksOutput structure for listing tasks
type ListUserTasksOutput struct {
    Tasks []struct {
        Name        string `json:"name"`
        Description string `json:"description"`
    } `json:"tasks"`
    Message string `json:"message"` // A message indicating the status of the listing (e.g., "Successfully listed X tasks")
}

// Task execution output string format
// "Task 'task_name' completed successfully. Output:\n[stdout]\nError Output:\n[stderr]"
// OR
// "Task 'task_name' failed. Output:\n[stdout]\nError Output:\n[stderr]\nExit Code: [code]\nError: [error_message]"
```

## Extra Considerations

*   **Security**: Ensure clear documentation on the security implications of executing user-defined commands.
*   **Error Reporting**: The agent needs to effectively communicate task execution errors, including non-zero exit codes, stdout, and stderr, to the LLM for proper diagnosis.
*   **Performance (Taskfile Loading)**: While a new `task.Executor` is created for each `RunUserTask` call to ensure isolation, the overhead of re-parsing the Taskfile should be monitored. For typical agent interactions, this is expected to be negligible.
*   **Argument Passing Complexity**: The reliance on `{{.CLI_ARGS}}` offloads parsing complexity to the Taskfile, but the LLM must understand this convention for passing arguments.
*   **Missing Taskfile**: The system should gracefully handle the absence of `.agent/Taskfile.yml` without errors, simply not registering the user-defined task tools.

## Testing Considerations

*   **Positive Scenarios**:
    *   Verify `list_user_tasks` correctly lists tasks from an existing `Taskfile.yml`.
    *   Verify `run_user_task` successfully executes a simple task (e.g., `greet`, `list-current-dir`) and returns expected stdout/stderr.
    *   Verify `run_user_task` passes arguments correctly to a task using `{{.CLI_ARGS}}` (e.g., `run-go-tests`).
*   **Negative Scenarios**:
    *   Verify `run_user_task` reports failure and the correct exit code for a task designed to fail (e.g., `deliberate-fail`).
    *   Verify `run_user_task` returns an error when `task_name` is empty or invalid.
    *   Verify the system behaves correctly when `.agent/Taskfile.yml` does not exist (i.e., tools are not registered).
*   **Concurrency**: Although `task.Executor` instances are isolated, consider testing rapid, sequential calls to `run_user_task` to ensure stability.

## Implementation Notes

*   Implement `internal/tools/user_defined_tasks.go` to handle Taskfile loading, task listing, and task execution using `github.com/go-task/task/v3`.
*   Update `internal/tools/registry.go` to include the new `NewUserDefinedTaskTools()` in the `NewRegistry` function, gracefully handling cases where the Taskfile is not found.
*   Ensure proper error handling and formatting of output messages for `run_user_task` to provide comprehensive feedback to the LLM.
*   The `GenerateSchema` function will be used to create the OpenAPI schemas for the new tools.
*   A new `task.Executor` instance should be created for *each* `RunUserTask` call to ensure isolated I/O and thread safety.

## Specification by Example

**Example `.agent/Taskfile.yml`:**

```yaml
# .agent/Taskfile.yml
version: '3'

tasks:
  greet:
    desc: "Prints a simple greeting message."
    cmds:
      - echo "Hello from your custom Taskfile!"
    silent: true

  list-current-dir:
    desc: "Lists the contents of the current working directory."
    cmds:
      - ls -lA

  run-go-tests:
    desc: "Runs Go tests for the project. Supports additional Go test flags via the 'args' input (e.g., '-v -race')."
    cmds:
      - go test ./... {{.CLI_ARGS}}
    silent: true

  deliberate-fail:
    desc: "A task designed to exit with an error code to demonstrate error handling."
    cmds:
      - echo "This task will fail!" && exit 1
    silent: true
```

**Conceptual Usage by LLM:**

1.  **Listing tasks:**
    *   LLM would call: `list_user_tasks()`
    *   Expected output: `{"tasks": [{"name": "greet", "description": "Prints a simple greeting message."}, ...], "message": "Successfully listed X user-defined tasks from .agent/Taskfile.yml."}`

2.  **Running a successful task:**
    *   LLM would call: `run_user_task(task_name="greet")`
    *   Expected output: `"Task 'greet' completed successfully. Output:\nHello from your custom Taskfile!\nError Output:\n"`

3.  **Running a task with arguments:**
    *   LLM would call: `run_user_task(task_name="run-go-tests", args="-v -race")`
    *   Expected output (example for success): `"Task 'run-go-tests' completed successfully. Output:\n[go test output]\nError Output:\n"`

4.  **Running a failing task:**
    *   LLM would call: `run_user_task(task_name="deliberate-fail")`
    *   Expected output: `"Task 'deliberate-fail' failed. Output:\nThis task will fail!\nError Output:\n\nExit Code: 1\nError: task: Failed to run task 'deliberate-fail'"`

## Verification

- [ ] Verify that `list_user_tasks` tool is available when `.agent/Taskfile.yml` exists.
- [ ] Verify that `run_user_task` tool is available when `.agent/Taskfile.yml` exists.
- [ ] Verify that `list_user_tasks` returns an accurate list of tasks and descriptions from the `Taskfile.yml`.
- [ ] Verify that `run_user_task` successfully executes a simple command and returns its stdout.
- [ ] Verify that `run_user_task` correctly handles tasks that output to stderr.
- [ ] Verify that `run_user_task` accurately reports the exit code for failing tasks.
- [ ] Verify that arguments passed to `run_user_task` are correctly received by the Taskfile via `{{.CLI_ARGS}}`.
- [ ] Verify that if `.agent/Taskfile.yml` does not exist, neither `list_user_tasks` nor `run_user_task` tools are registered, and no errors are thrown.
- [ ] Verify that multiple sequential calls to `run_user_task` function correctly without output intermingling.