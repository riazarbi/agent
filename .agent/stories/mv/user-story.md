# Implement mv Tool

The agent currently lacks a direct and efficient way to move files and directories. This user story aims to provide a focused, intuitive `mv` tool for relocating file system items, reducing the cognitive load and execution time for the agent.

## Past Attempts

N/A - This is a new feature set.

## Requirements

*   **Implement `mv` tool:**
    *   Accepts `source` (string) and `destination` (string) arguments.
    *   Moves the file or directory located at `source` to `destination`.
    *   If `source` is a directory, it moves the directory and all its contents recursively.
    *   If `destination` is an existing directory, `source` is moved *into* that directory.
    *   Fails with an error message (e.g., "No such file or directory") if `source` does not exist.
    *   Fails with an error message (e.g., "File exists") if `destination` is an existing file and `source` is also a file (without an explicit overwrite flag, which is not being added for `mv` at this time).
    *   Shells out to the system's `mv` command for execution, ensuring identical behavior to the standard GNU/Linux utility.

## Rules

*   The `mv` tool MUST be implemented by shelling out to the corresponding system command.
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

**YOU CANNOT TEST THESE NEW TOOLS, A NEW BINARY MUST BE BUILT FIRST. WRITE INSTRUCTIONS FOR TESTING TO A FILE CALLED check.tct, OVERWRITING PREVIOUS CONTENT**


*   **Integration Tests:** For the `mv` tool, integration tests are crucial. These should:
    *   Run against a real file system.
    *   Verify the exact behavior for all described requirements.
    *   Create and clean up isolated temporary directories for each test case to prevent test interference.
    *   Specifically verify that `mv` produces the same output/errors as its GNU/Linux counterpart. This might involve capturing output of both the tool and a direct shell call and comparing.
    *   Test moving files to files, files to directories, directories to directories.
    *   Verify correct error messages are returned for all failure scenarios.

## Implementation Notes

*   All shelled-out commands should use `os/exec` package in Go, ensuring careful handling of arguments to prevent shell injection (i.e., pass arguments as separate strings, not a single command string).
*   Error wrapping should be used to provide context where native Go errors are returned, maintaining clarity for the agent.

## Specification by Example

### `mv`
*   **Example 1: Move file to new name**
    *   Pre-condition: `old_file.txt` exists.
    *   `mv(source="old_file.txt", destination="new_file.txt")`
    *   Expected output: `{"message": "Successfully moved 'old_file.txt' to 'new_file.txt'."}`
    *   Expected file system state: `old_file.txt` does not exist, `new_file.txt` exists with original content.
*   **Example 2: Move file into existing directory**
    *   Pre-condition: `file.txt` exists, `my_dir/` exists.
    *   `mv(source="file.txt", destination="my_dir/")`
    *   Expected output: `{"message": "Successfully moved 'file.txt' to 'my_dir/file.txt'."}`
    *   Expected file system state: `file.txt` does not exist, `my_dir/file.txt` exists.
*   **Example 3: Move directory**
    *   Pre-condition: `my_old_dir/` exists with content.
    *   `mv(source="my_old_dir/", destination="my_new_dir/")`
    *   Expected output: `{"message": "Successfully moved 'my_old_dir/' to 'my_new_dir/'."}`
    *   Expected file system state: `my_old_dir/` does not exist, `my_new_dir/` exists with original content.
*   **Example 4: Fail on non-existent source**
    *   `mv(source="nonexistent.txt", destination="target.txt")`
    *   Expected output (error): `{"error": "No such file or directory: 'nonexistent.txt'"}`
*   **Example 5: Fail on destination file exists (no overwrite)**
    *   Pre-condition: `source.txt` exists, `destination.txt` exists.
    *   `mv(source="source.txt", destination="destination.txt")`
    *   Expected output (error): `{"error": "File exists: 'destination.txt'"}`

## Verification

- [ ] `mv` tool is implemented and available to the agent.
- [ ] `mv` moves files correctly.
- [ ] `mv` moves directories and their contents recursively.
- [ ] `mv` moves source items into an existing destination directory.
- [ ] `mv` fails gracefully with a "No such file or directory" error if the source does not exist.
- [ ] `mv` fails gracefully with a "File exists" error if the destination file exists and no overwrite is specified.
- [ ] `mv`'s behavior, arguments, and error messages are identical to the standard GNU/Linux `mv` command.