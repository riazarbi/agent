# Refactor main.go: Extract All Tools to tools/ Package

*Reduce the size and complexity of `main.go` by extracting all tool definitions (excluding `todowrite` and `todoread`) into separate files within a unified `tools/` package. This refactoring aims to improve code organization, readability, and maintainability.*

## Requirements

- Create a new directory named `tools/` if it doesn't already exist.
- For each tool currently defined in `main.go` (excluding `todowrite` and `todoread`):
    - Read the content of `main.go`.
    - Identify and extract the code related to the tool's definition (structs, functions, etc.).
    - Create a new file `tools/<tool_name>.go` (e.g., `tools/read_file.go`, `tools/edit_file.go`).
    - Place the extracted tool definition code into its respective `tools/<tool_name>.go` file, ensuring the package declared is `package tools`.
    - Modify `main.go` to import the `tools` package.
    - Adapt all calls to the refactored tools in `main.go` to use the `tools.<ToolName>` format (e.g., `agent.tools.ListFiles` becomes `tools.ListFiles`).
- Ensure all existing functionality related to the refactored tools remains intact after the refactoring.
- The `go mod tidy` and `go build .` commands should execute successfully after the changes.
- All refactoring changes should be part of a single Git commit.
- Create a `manual_verification.md` file with a clear prompt for an agent to manually verify the functionality of all refactored tools.

## Rules

- The refactoring must adhere to Go coding standards and best practices.
- No external dependencies should be introduced that are not already present.
- Each refactored tool's code must be placed solely in its respective `tools/<tool_name>.go` file.
- `todowrite` and `todoread` tools must remain in `main.go` and not be refactored as part of this story.
- Changes to `main.go` should be minimal and focused only on importing and using the refactored tools; no new features or unrelated chores should be introduced.
- The `list_files.go` file should also be updated to declare `package tools`.

## Domain

```go
// Represents the various tool structures and methods, now unified under a 'tools' package.
// For example:
// In tools/read_file.go:
// package tools
// type ReadFileTool struct { ... }
// func (t *ReadFileTool) ReadFile(path string) (string, error) { ... }

// In main.go, calls will change from:
// agent.read_file(path)
// To:
// tools.ReadFile(path) // assuming tools package is imported and methods are direct functions or accessed via a tools struct
```

## Extra Considerations


- Implement robust error handling in the refactoring process, especially when using `edit_file`.
- Maintain a todo list to track the progress of refactoring each individual tool.

## Testing Considerations

*Manual verification is required for this refactoring. The executor will create a `manual_verification.md` file that provides a clear prompt for an agent to perform comprehensive manual tests on all refactored tools to ensure their functionality remains intact.*

## Implementation Notes

- Use `read_file` to get `main.go` content.
- Use `edit_file` for precise, small changes, verifying after each edit.
- Avoid large block replacements with `edit_file`.
- Use `grep` to confirm the presence or absence of strings before and after edits.
- Use the `todowrite` and `todoread` tools to manage the list of tools to be refactored and track their progress.

## Specification by Example

*N/A - The existing functionality of all tools should remain the same after refactoring.*

## Verification

- [ ] `tools/` directory exists.
- [ ] For each refactored tool (`read_file`, `edit_file`, `delete_file`, `grep`, `glob`, `git_diff`, `web_fetch`, `html_to_markdown`, `head`, `tail`, `cloc`):
    - [ ] `tools/<tool_name>.go` exists and contains the correct tool definition with `package tools`.
    - [ ] `main.go` no longer contains the tool's definition and correctly imports the `tools` package.
    - [ ] Calls to the tool in `main.go` are updated to `tools.<ToolName>`.
- [ ] `main.go` still contains the `todowrite` and `todoread` tool definitions.
- [ ] `go mod tidy` and `go build .` execute successfully (verified by human user).
- [ ] All refactored tools can be invoked and function correctly (verified manually via `manual_verification.md`).
- [ ] A `git diff` shows only the expected changes for this refactoring (single commit).
- [ ] `manual_verification.md` file exists and contains a comprehensive manual verification prompt.