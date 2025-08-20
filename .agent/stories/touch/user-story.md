# Implement touch Tool

The agent currently relies on creative and often inefficient multi-step solutions to perform basic file system operations. This user story aims to provide a focused, intuitive `touch` tool for creating and updating file timestamps, reducing the cognitive load and execution time for the agent.

## Past Attempts

N/A - This is a new feature set.

## Requirements

*   **Implement `touch` tool:**
    *   Accepts a `path` argument (string).
    *   Creates an empty file at the specified `path` if it does not exist.
    *   If the file already exists, it updates the file's modification timestamp without altering its content.
    *   Fails with an error message (e.g., "No such file or directory") if any parent directory in the `path` does not exist.
    *   Shells out to the system's `touch` command for execution, ensuring identical behavior to the standard GNU/Linux utility.

## Rules

*   The `touch` tool MUST be implemented by shelling out to the corresponding system command.
*   All tools must return clear, concise success messages upon completion.
*   All tools must return clear, concise error messages, closely mimicking standard GNU/Linux command error outputs (e.g., "No such file or directory").

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

**YOU CANNOT TEST THESE NEW TOOLS, A NEW BINARY MUST BE BUILT FIRST. WRITE INSTRUCTIONS FOR TESTING TO A FILE CALLED check.tct, OVERWRITING PREVIOUS CONTENT**


*   **Integration Tests:** For the `touch` tool, integration tests are crucial. These should:
    *   Run against a real file system.
    *   Verify the exact behavior for all described requirements (e.g., `touch` failing on non-existent parent directory).
    *   Create and clean up isolated temporary directories for each test case to prevent test interference.
    *   Specifically verify that `touch` produces the same output/errors as its GNU/Linux counterpart. This might involve capturing output of both the tool and a direct shell call and comparing.
    *   Verify correct error messages are returned for all failure scenarios.

## Implementation Notes

*   All shelled-out commands should use `os/exec` package in Go, ensuring careful handling of arguments to prevent shell injection (i.e., pass arguments as separate strings, not a single command string).
*   Error wrapping should be used to provide context where native Go errors are returned, maintaining clarity for the agent.

## Specification by Example

### `touch`
*   **Example 1: Create new file**
    *   `touch(path="new_file.txt")`
    *   Expected output: `{"message": "Successfully touched 'new_file.txt'."}`
    *   Expected file system state: `new_file.txt` exists.
*   **Example 2: Update existing file**
    *   `touch(path="existing_file.txt")`
    *   Expected output: `{"message": "Successfully touched 'existing_file.txt'."}`
    *   Expected file system state: `existing_file.txt` exists, timestamp updated.
*   **Example 3: Fail on non-existent parent**
    *   `touch(path="nonexistent_dir/new_file.txt")`
    *   Expected output (error): `{"error": "No such file or directory: 'nonexistent_dir/new_file.txt'"}`

## Verification

- [ ] `touch` tool is implemented and available to the agent.
- [ ] `touch` creates a new empty file if it doesn't exist.
- [ ] `touch` updates the timestamp of an existing file without changing its content.
- [ ] `touch` fails gracefully with a "No such file or directory" error if parent directories are missing.
- [ ] `touch`'s behavior, arguments, and error messages are identical to the standard GNU/Linux `touch` command.