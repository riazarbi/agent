# Implement append_file Tool

The agent currently lacks a direct and efficient way to append content to files. This user story aims to provide a focused, intuitive `append_file` tool for adding content to the end of existing files or creating new files with initial content, reducing the cognitive load and execution time for the agent.

## Past Attempts

N/A - This is a new feature set.

## Requirements

*   **Implement `append_file` tool:**
    *   Accepts `path` (string) and `content` (string) arguments.
    *   Appends `content` to the end of the file at `path`.
    *   If the file at `path` does not exist, it creates the file and writes the `content` as its initial content.
    *   Fails with an error message (e.g., "Is a directory") if `path` refers to an existing directory.
    *   Is implemented natively in Go, providing clear success/failure messages.

## Rules

*   The `append_file` tool MUST be implemented natively in Go, and its name (`append_file`) ensures it does not directly correspond to a standard GNU/Linux command name.
*   All tools must return clear, concise success messages upon completion.
*   All tools must return clear, concise error messages, closely mimicking standard GNU/Linux command error outputs (e.g., "Is a directory").

## Domain

```
// Filesystem operations
type FileSystemTool interface {
    Execute(args map[string]interface{}) (string, error)
}
```

## Extra Considerations

*   Error messages for native Go tools should provide clear context and ideally map to common Unix error types where applicable.

## Testing Considerations

**YOU CANNOT TEST THESE NEW TOOLS, A NEW BINARY MUST BE BUILT FIRST. PROVIDE THE USER WITH INSTRUCTIONS FOR TESTING**


*   **Unit Tests:** For `append_file` (native Go implementation), unit tests should cover:
    *   Successful creation of a new file with content.
    *   Successful appending to an existing file.
    *   Correct error handling when `path` is a directory.
    *   Correct error handling for underlying OS issues (e.g., permissions).
*   **Integration Tests:** Integration tests are crucial. These should:
    *   Run against a real file system.
    *   Verify the exact behavior for all described requirements.
    *   Create and clean up isolated temporary directories for each test case to prevent test interference.
    *   Verify correct error messages are returned for all failure scenarios.

## Implementation Notes

*   Standard Go `os` and `io/ioutil` packages should be used for `append_file` implementation.
*   Error wrapping should be used to provide context where native Go errors are returned, maintaining clarity for the agent.

## Specification by Example

### `append_file`
*   **Example 1: Create new file with content**
    *   `append_file(path="log.txt", content="First log entry.\n")`
    *   Expected output: `{"message": "Successfully appended to 'log.txt'."}`
    *   Expected file system state: `log.txt` exists with "First log entry.\n" as content.
*   **Example 2: Append to existing file**
    *   Pre-condition: `log.txt` contains "First log entry.\n"
    *   `append_file(path="log.txt", content="Second log entry.\n")`
    *   Expected output: `{"message": "Successfully appended to 'log.txt'."}`
    *   Expected file system state: `log.txt` contains "First log entry.\nSecond log entry.\n".
*   **Example 3: Fail on path being a directory**
    *   `append_file(path="my_directory/", content="Some content")`
    *   Expected output (error): `{"error": "Is a directory: 'my_directory/'"}`

## Verification

- [ ] `append_file` tool is implemented and available to the agent.
- [ ] `append_file` creates a new file and writes content if the file does not exist.
- [ ] `append_file` appends content to an existing file.
- [ ] `append_file` fails gracefully with an "Is a directory" error if the path points to a directory.