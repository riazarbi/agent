# Refactor main.go: Extract Tool Definitions

*Reduce the size of `main.go` by extracting the tool definitions into separate files within a `tools/` directory. This will simplify `main.go` and improve its readability.*

## Requirements

- Create a new directory named `tools/`.
- Read the content of `main.go`.
- Identify and extract the code related to each tool's definition (structs, functions, etc.) from `main.go`.
- Create a separate `.go` file for each tool in the `tools/` directory (e.g., `tools/tool1.go`, `tools/tool2.go`).
- Place the corresponding tool definition code into its respective file.
- Modify `main.go` to import the newly created tool packages from the `tools/` directory.
- Ensure all existing functionality remains intact after the refactoring.
- Reduce the line count of `main.go` significantly.
- All the new files should compile.

## Rules

- The refactoring must adhere to Go coding standards and best practices.
- No external dependencies should be introduced.
- Each tool's code must be placed in its own file in the `tools/` directory.
- The main.go file is very large. Be conservative in your read_file and edit_file tool use.

## Domain

```
// Example: ToolDefinition struct (This should be replaced with actual tool definitions)
type ToolDefinition struct { ... }
```

## Extra Considerations

- Consider adding comments to the new tool files to improve readability.
- Ensure proper error handling in the new modules.

## Testing Considerations

- Unit tests may be required for the new tool files to verify their functionality.
- Integration tests should be performed to ensure that all components work together correctly.

## Implementation Notes

- Use Go modules to manage dependencies.

## Specification by Example

*N/A*

## Verification

- [ ] `tools/` directory exists.
- [ ] A `.go` file exists for each tool in the `tools/` directory.
- [ ] Each tool file contains the correct tool definition code.
- [ ] `main.go` imports the tool packages from `tools/`.
- [ ] All files compile without errors.
- [ ] All existing functionality works as expected.
- [ ] `main.go` has a significantly reduced line count.

## Next Steps

1. Save to `.agent/stories/refactor-main-go/user-story.md`