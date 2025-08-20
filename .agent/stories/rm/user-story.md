# Implement rm Tool

The agent currently lacks a direct and efficient way to delete files and directories. This user story aims to provide a focused, intuitive `rm` tool for removing file system items, including recursive deletion capabilities with safety precautions, reducing the cognitive load and execution time for the agent.

## Past Attempts

N/A - This is a new feature set.

## Requirements

*   **Implement `rm` tool:**
    *   Accepts `path` (string) and an optional `recursive` (boolean) argument, defaulting to `false`.
    *   Deletes the file or directory at `path`.
    *   If `path` refers to an empty directory, it deletes the directory.
    *   If `path` refers to a non-empty directory and `recursive` is `false` (or omitted), it fails with an error message (e.g., "Directory not empty").
    *   If `path` refers to a non-empty directory and `recursive` is `true`, it deletes the directory and all its contents recursively.
    *   Fails with an error message (e.g., "No such file or directory") if `path` does not exist.
    *   Shells out to the system's `rm` command for execution, ensuring identical behavior to the standard GNU/Linux utility (including its `-r` flag behavior).

## Rules

*   The `rm` tool MUST be implemented by shelling out to the corresponding system command.
*   All tools must return clear, concise success messages upon completion.
*   All tools must return clear, concise error messages, closely mimicking standard GNU/Linux command error outputs (e.g., "No such file or directory", "Directory not empty").

## Domain

```
// Filesystem operations
type FileSystemTool interface {
    Execute(args map[string]interface{}) (string, error)
}
```

## Extra Considerations

*   Error messages for shelling out tools should capture the underlying system command's `stderr` as accurately as possible to ensure fidelity.
*   Due to the destructive nature of `rm -r`, extra care should be taken in testing to ensure proper path validation and that tests are run in isolated, temporary environments.

## Testing Considerations

**YOU CANNOT TEST THESE NEW TOOLS, A NEW BINARY MUST BE BUILT FIRST. WRITE INSTRUCTIONS FOR TESTING TO A FILE CALLED check.tct, OVERWRITING PREVIOUS CONTENT**


*   **Integration Tests:** For the `rm` tool, integration tests are crucial. These should:
    *   Run against a real file system.
    *   Verify the exact behavior for all described requirements.
    *   Create and clean up isolated temporary directories for each test case to prevent test interference.
    *   Specifically verify that `rm` produces the same output/errors as its GNU/Linux counterpart. This might involve capturing output of both the tool and a direct shell call and comparing.
    *   Test deleting files, empty directories, and non-empty directories (with and without `recursive=true`).
    *   Verify correct error messages are returned for all failure scenarios.

## Implementation Notes

*   All shelled-out commands should use `os/exec` package in Go, ensuring careful handling of arguments to prevent shell injection (i.e., pass arguments as separate strings, not a single command string).
*   Error wrapping should be used to provide context where native Go errors are returned, maintaining clarity for the agent.

## Specification by Example

### `rm`
*   **Example 1: Delete file**
    *   Pre-condition: `file_to_delete.txt` exists.
    *   `rm(path="file_to_delete.txt")`
    *   Expected output: `{"message": "Successfully removed 'file_to_delete.txt'."}`
    *   Expected file system state: `file_to_delete.txt` does not exist.
*   **Example 2: Delete empty directory**
    *   Pre-condition: `empty_dir/` exists and is empty.
    *   `rm(path="empty_dir/")`
    *   Expected output: `{"message": "Successfully removed directory 'empty_dir/'."}`
    *   Expected file system state: `empty_dir/` does not exist.
*   **Example 3: Fail on non-empty directory (default behavior)**
    *   Pre-condition: `non_empty_dir/` exists and contains files.
    *   `rm(path="non_empty_dir/")`
    *   Expected output (error): `{"error": "Directory not empty: 'non_empty_dir/'"}`
*   **Example 4: Delete non-empty directory recursively (`recursive=true`)**
    *   Pre-condition: `non_empty_dir/` exists and contains files.
    *   `rm(path="non_empty_dir/", recursive=true)`
    *   Expected output: `{"message": "Successfully removed directory 'non_empty_dir/' recursively."}`
    *   Expected file system state: `non_empty_dir/` does not exist.
*   **Example 5: Fail on non-existent path**
    *   `rm(path="nonexistent_item.txt")`
    *   Expected output (error): `{"error": "No such file or directory: 'nonexistent_item.txt'"}`

## Verification

- [ ] `rm` tool is implemented and available to the agent.
- [ ] `rm` deletes files correctly.
- [ ] `rm` deletes empty directories correctly.
- [ ] `rm` fails gracefully with a "Directory not empty" error if the directory is not empty and `recursive=false`.
- [ ] `rm` deletes non-empty directories and their contents recursively when `recursive=true`.
- [ ] `rm` fails gracefully with a "No such file or directory" error if the path does not exist.
- [ ] `rm`'s behavior, arguments, and error messages are identical to the standard GNU/Linux `rm` command (including `-r` flag behavior).