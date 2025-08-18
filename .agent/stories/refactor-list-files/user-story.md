# Refactor main.go: Extract list_files Tool

*Reduce the size of `main.go` by extracting the `list_files` tool definition into a separate file within a `tools/` directory. This will simplify `main.go` and improve its readability, serving as a proof of concept for future tool extractions.*

## Requirements

- Create a new directory named `tools/` if it doesn't already exist.
- Read the content of `main.go`.
- Identify and extract the code related to the `list_files` tool's definition (structs, functions, etc.) from `main.go`.
- Create a new file `tools/list_files.go`.
- Place the extracted `list_files` tool definition code into `tools/list_files.go`.
- Modify `main.go` to import the new `list_files` package from the `tools/` directory.
- Ensure all existing functionality related to `list_files` remains intact after the refactoring.
- The `go mod tidy` and `go build .` commands should execute successfully after the changes.

## Rules

- The refactoring must adhere to Go coding standards and best practices.
- No external dependencies should be introduced.
- The `list_files` tool's code must be placed solely in `tools/list_files.go`.
- Changes to `main.go` should be minimal and focused only on importing and using the refactored tool.

## Domain

```go
// Represents the `list_files` tool's structure and methods.
// This will be moved from main.go to tools/list_files.go
type ListFilesTool struct {
    // ... relevant fields for list_files
}

// Function signature for list_files
func (a *Agent) listFiles(path string) (string, error) {
    // ... implementation details
}
```

## Extra Considerations

- Determine the correct Go module path from `go.mod` for the import statement.
- Implement robust error handling in the refactoring process, especially when using `edit_file`.
- Create a backup of `main.go` before making any modifications.

## Testing Considerations

*For this initial refactoring, testing will involve manual verification that the `list_files` tool functions correctly after the changes, ensuring its existing functionality remains intact.*

## Implementation Notes

- Use `read_file` to get `main.go` content.
- Use `edit_file` for precise, small changes, verifying after each edit.
- Avoid large block replacements with `edit_file`.
- Use `grep` to confirm the presence or absence of strings before and after edits.

## Specification by Example

*N/A - the existing `list_files` functionality should remain the same.*

## Verification

- [ ] `tools/` directory exists.
- [ ] `tools/list_files.go` exists and contains the correct `list_files` tool definition code.
- [ ] `main.go` no longer contains the `list_files` tool definition and correctly imports the `tools` package.
- [ ] `go mod tidy` and `go build .` execute successfully.
- [ ] The `list_files` tool can be invoked and functions correctly through the agent (manual verification).
- [ ] `main.go` has a reduced line count related to the `list_files` tool's removal.
- [ ] A `git diff` shows only the expected changes.