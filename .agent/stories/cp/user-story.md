# Implement cp Tool

The agent currently lacks a direct and efficient way to copy files and directories. This user story aims to provide a focused, intuitive `cp` tool for duplicating file system items, including recursive copy capabilities, reducing the cognitive load and execution time for the agent.

## Past Attempts

N/A - This is a new feature set.

## Requirements

*   **Implement `cp` tool:**
    *   Accepts `source` (string) and `destination` (string) arguments.
    *   Accepts an optional `recursive` (boolean) argument, defaulting to `false`, required when copying directories.
    *   Copies the file or directory located at `source` to `destination`.
    *   If `source` is a file and `destination` is an existing file, `source` will overwrite `destination`.
    *   If `source` is a file and `destination` is an existing directory, `source` is copied *into* that directory.
    *   If `source` is a directory and `recursive` is `false`, it fails with an error (e.g., "-r not specified; omitting directory 'source'").
    *   If `source` is a directory and `recursive` is `true`, it copies the directory and its contents recursively. If `destination` exists and is a directory, `source` is copied *into* it (e.g., `cp -r source_dir/ destination_dir/` results in `destination_dir/source_dir/`). If `destination` does not exist, a new directory named `destination` is created with contents of `source`.
    *   Fails with an error message (e.g., "No such file or directory") if `source` does not exist.
    *   Shells out to the system's `cp` command for execution, ensuring identical behavior to the standard GNU/Linux utility.

## Rules

*   The `cp` tool MUST be implemented by shelling out to the corresponding system command.
*   All tools must return clear, concise success messages upon completion.
*   All tools must return clear, concise error messages, closely mimicking standard GNU/Linux command error outputs (e.g., "No such file or directory", "-r not specified; omitting directory 'source'").

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

*   **Integration Tests:** For the `cp` tool, integration tests are crucial. These should:
    *   Run against a real file system.
    *   Verify the exact behavior for all described requirements.
    *   Create and clean up isolated temporary directories for each test case to prevent test interference.
    *   Specifically verify that `cp` produces the same output/errors as its GNU/Linux counterpart. This might involve capturing output of both the tool and a direct shell call and comparing.
    *   Test copying files to files (overwrite), files to directories, and directories to directories (with and without `recursive=true`).
    *   Verify correct error messages are returned for all failure scenarios.

## Implementation Notes

*   All shelled-out commands should use `os/exec` package in Go, ensuring careful handling of arguments to prevent shell injection (i.e., pass arguments as separate strings, not a single command string).
*   Error wrapping should be used to provide context where native Go errors are returned, maintaining clarity for the agent.

## Specification by Example

### `cp`
*   **Example 1: Copy file to new file (overwrite if exists)**
    *   Pre-condition: `source.txt` exists. `destination.txt` may or may not exist.
    *   `cp(source="source.txt", destination="destination.txt")`
    *   Expected output: `{"message": "Successfully copied 'source.txt' to 'destination.txt'."}`
    *   Expected file system state: `source.txt` exists, `destination.txt` exists with content of `source.txt`.
*   **Example 2: Copy file into existing directory**
    *   Pre-condition: `file.txt` exists, `my_dir/` exists.
    *   `cp(source="file.txt", destination="my_dir/")`
    *   Expected output: `{"message": "Successfully copied 'file.txt' to 'my_dir/file.txt'."}`
    *   Expected file system state: `file.txt` exists, `my_dir/file.txt` exists with content of `file.txt`.
*   **Example 3: Fail on copying directory without `recursive=true`**
    *   Pre-condition: `my_dir_source/` exists.
    *   `cp(source="my_dir_source/", destination="my_dir_destination/")`
    *   Expected output (error): `{"error": "-r not specified; omitting directory 'my_dir_source/'"}` (or similar system error message).
*   **Example 4: Copy directory recursively to non-existent destination**
    *   Pre-condition: `source_dir/` exists with content. `new_destination_dir/` does NOT exist.
    *   `cp(source="source_dir/", destination="new_destination_dir/", recursive=true)`
    *   Expected output: `{"message": "Successfully copied 'source_dir/' to 'new_destination_dir/' recursively."}`
    *   Expected file system state: `source_dir/` exists, `new_destination_dir/` exists with the content of `source_dir/`.
*   **Example 5: Copy directory recursively into existing directory**
    *   Pre-condition: `source_dir/` exists with content. `existing_destination_dir/` exists.
    *   `cp(source="source_dir/", destination="existing_destination_dir/", recursive=true)`
    *   Expected output: `{"message": "Successfully copied 'source_dir/' into 'existing_destination_dir/' recursively."}`
    *   Expected file system state: `source_dir/` exists, `existing_destination_dir/source_dir/` exists with the content of `source_dir/`.
*   **Example 6: Fail on non-existent source**
    *   `cp(source="nonexistent.txt", destination="target.txt")`
    *   Expected output (error): `{"error": "No such file or directory: 'nonexistent.txt'"}`

## Verification

- [ ] `cp` tool is implemented and available to the agent.
- [ ] `cp` copies files correctly, overwriting destination file if it exists.
- [ ] `cp` copies files into existing directories correctly.
- [ ] `cp` fails gracefully if a directory is copied without `recursive=true`.
- [ ] `cp` copies directories and their contents recursively when `recursive=true`.
- [ ] `cp` copies directories into an existing destination directory when `recursive=true`.
- [ ] `cp` creates a new destination directory with contents when copying a directory recursively to a non-existent path.
- [ ] `cp` fails gracefully with a "No such file or directory" error if the source does not exist.
- [ ] `cp`'s behavior, arguments, and error messages are identical to the standard GNU/Linux `cp` command.