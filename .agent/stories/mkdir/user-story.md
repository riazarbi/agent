# Implement mkdir Tool

The agent currently lacks a direct and efficient way to create directories. This user story aims to provide a focused, intuitive `mkdir` tool for directory creation, including recursive capabilities, reducing the cognitive load and execution time for the agent.

## Past Attempts

N/A - This is a new feature set.

## Requirements

*   **Implement `mkdir` tool:**
    *   Accepts `path` (string) and an optional `parents` (boolean) argument, defaulting to `false`.
    *   Creates the directory at `path`.
    *   If `parents` is `true`, it creates any necessary parent directories along the specified `path` if they do not exist.
    *   If `parents` is `true` and the directory at `path` already exists, the operation succeeds without making changes (idempotent).
    *   If `parents` is `false` (or omitted), it fails with an error message (e.g., "File exists") if the directory already exists.
    *   If `parents` is `false` (or omitted), it fails with an error message (e.g., "No such file or directory") if any parent directory in the `path` does not exist.
    *   Shells out to the system's `mkdir` command for execution, ensuring identical behavior to the standard GNU/Linux utility (including its `-p` flag behavior).

## Rules

*   The `mkdir` tool MUST be implemented by shelling out to the corresponding system command.
*   All tools must return clear, concise success messages upon completion.
*   All tools must return clear, concise error messages, closely mimicking standard GNU/Linux command error outputs (e.g., "No such file or directory", "File exists").

## Domain

```
// Filesystem operations
type FileSystemTool interface {
    Execute(args map[string]interface{}) (string, error)
}
```

## Extra Considerations

*   Error messages for shelling out tools should capture the underlying system command's `stderr` as accurately as possible to ensure fidelity.

## Testing Considerations

**YOU CANNOT TEST THESE NEW TOOLS, A NEW BINARY MUST BE BUILT FIRST. PROVIDE THE USER WITH INSTRUCTIONS FOR TESTING**


*   **Integration Tests:** For the `mkdir` tool, integration tests are crucial. These should:
    *   Run against a real file system.
    *   Verify the exact behavior for all described requirements.
    *   Create and clean up isolated temporary directories for each test case to prevent test interference.
    *   Specifically verify that `mkdir` produces the same output/errors as its GNU/Linux counterpart. This might involve capturing output of both the tool and a direct shell call and comparing.
    *   Test both `parents=true` and `parents=false` scenarios, including creation of existing directories.
    *   Verify correct error messages are returned for all failure scenarios.

## Implementation Notes

*   All shelled-out commands should use `os/exec` package in Go, ensuring careful handling of arguments to prevent shell injection (i.e., pass arguments as separate strings, not a single command string).
*   Error wrapping should be used to provide context where native Go errors are returned, maintaining clarity for the agent.

## Specification by Example

### `mkdir`
*   **Example 1: Create single directory**
    *   `mkdir(path="new_dir/")`
    *   Expected output: `{"message": "Successfully created directory 'new_dir/'."}`
    *   Expected file system state: `new_dir/` exists.
*   **Example 2: Create recursive directories (`parents=true`)**
    *   `mkdir(path="path/to/new_recursive_dir/", parents=true)`
    *   Expected output: `{"message": "Successfully created directory 'path/to/new_recursive_dir/'."}`
    *   Expected file system state: `path/`, `path/to/`, and `path/to/new_recursive_dir/` all exist.
*   **Example 3: Idempotent creation (`parents=true` and directory exists)**
    *   Pre-condition: `existing_dir/` exists.
    *   `mkdir(path="existing_dir/", parents=true)`
    *   Expected output: `{"message": "Successfully ensured directory 'existing_dir/' exists."}`
    *   Expected file system state: `existing_dir/` still exists (no change).
*   **Example 4: Fail on existing directory (default behavior)**
    *   Pre-condition: `existing_dir/` exists.
    *   `mkdir(path="existing_dir/")`
    *   Expected output (error): `{"error": "File exists: 'existing_dir/'"}`
*   **Example 5: Fail on non-existent parent (default behavior)**
    *   `mkdir(path="nonexistent_parent/new_dir/")`
    *   Expected output (error): `{"error": "No such file or directory: 'nonexistent_parent/new_dir/'"}`

## Verification

- [ ] `mkdir` tool is implemented and available to the agent.
- [ ] `mkdir` creates a single directory correctly.
- [ ] `mkdir` creates parent directories recursively when `parents=true`.
- [ ] `mkdir` is idempotent when `parents=true` and the directory already exists.
- [ ] `mkdir` fails gracefully with a "File exists" error if the directory already exists and `parents=false`.
- [ ] `mkdir` fails gracefully with a "No such file or directory" error if parent directories are missing and `parents=false`.
- [ ] `mkdir`'s behavior, arguments, and error messages are identical to the standard GNU/Linux `mkdir` command (including `-p` flag behavior).