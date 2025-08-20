# Implement write_file Tool

The agent currently lacks a direct and efficient way to overwrite files with specific content or create new files with content. This user story aims to provide a focused, intuitive `write_file` tool for managing file contents, reducing the cognitive load and execution time for the agent.

## Past Attempts

N/A - This is a new feature set.

## Requirements

*   **Implement `write_file` tool:**
    *   Accepts `path` (string) and `content` (string) arguments.
    *   Overwrites the entire content of the file at `path` with `content`.
    *   If the file at `path` does not exist, it creates the file and writes the `content`.
    *   Fails with an error message (e.g., "Is a directory") if `path` refers to an existing directory.
    *   Is implemented natively in Go, providing clear success/failure messages.

## Rules

*   The `write_file` tool MUST be implemented natively in Go, and its name (`write_file`) ensures it does not directly correspond to a standard GNU/Linux command name.
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


*   **Unit Tests:** For `write_file` (native Go implementation), unit tests should cover:
    *   Successful creation of a new file with content.
    *   Successful overwriting of an existing file with new content.
    *   Correct error handling when `path` is a directory.
    *   Correct error handling for underlying OS issues (e.g., permissions).
*   **Integration Tests:** Integration tests are crucial. These should:
    *   Run against a real file system.
    *   Verify the exact behavior for all described requirements.
    *   Create and clean up isolated temporary directories for each test case to prevent test interference.
    *   Verify correct error messages are returned for all failure scenarios.

## Implementation Notes

*   Standard Go `os` and `io/ioutil` packages should be used for `write_file` implementation.
*   Error wrapping should be used to provide context where native Go errors are returned, maintaining clarity for the agent.

## Specification by Example

### `write_file`
*   **Example 1: Create new file with content**
    *   `write_file(path="my_doc.txt", content="This is my new document.")`
    *   Expected output: `{"message": "Successfully wrote to 'my_doc.txt'."}`
    *   Expected file system state: `my_doc.txt` exists with "This is my new document." as content.
*   **Example 2: Overwrite existing file**
    *   `write_file(path="existing_doc.txt", content="This is the updated content.")` (Assume `existing_doc.txt` previously contained "Old content.")
    *   Expected output: `{"message": "Successfully wrote to 'existing_doc.txt'."}`
    *   Expected file system state: `existing_doc.txt` contains "This is the updated content.".
*   **Example 3: Fail on path being a directory**
    *   `write_file(path="my_directory/", content="Some content")`
    *   Expected output (error): `{"error": "Is a directory: 'my_directory/'"}`

## Verification

- [ ] `write_file` tool is implemented and available to the agent.
- [ ] `write_file` creates a new file and writes content if the file does not exist.
- [ ] `write_file` overwrites an existing file's content.
- [ ] `write_file` fails gracefully with an "Is a directory" error if the path points to a directory.