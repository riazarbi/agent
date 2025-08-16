# [Implement Glob Pattern File Search Tool]

*Enable the agent to find files matching glob patterns in the local directory using Go's filepath.Glob, providing a safe way for the LLM to discover files without modification capabilities.*

## Requirements

- Must accept a glob pattern as input (e.g., "*.go", "**/*.md")
- Must be non-modifying (read-only operations)
- Must return matching file paths
- Must follow existing tool pattern in codebase
- Must handle invalid patterns gracefully
- Must be restricted to the working directory scope

## Rules

- Tool must never modify any files
- Must be read-only operations
- Must handle both successful and error cases consistently with other tools
- Must validate input parameters before execution
- Must not allow traversal outside working directory for security

## Domain

```
type GlobInput struct {
    Pattern string `json:"pattern" jsonschema_description:"The glob pattern to match files against (e.g. *.go, **/*.md)"`
}

var GlobDefinition = ToolDefinition{
    Name: "glob",
    Description: "Find files matching a glob pattern. Supports standard glob syntax for file discovery.",
    InputSchema: GlobInputSchema,
    Function: Glob,
}
```

## Extra Considerations

- Pattern validation to prevent malicious patterns
- Handling of hidden files and directories
- Relative path handling
- Performance with large directory structures
- Proper error messages for invalid patterns

## Testing Considerations

- Test with various glob patterns (*.go, **/*.md, etc.)
- Test with invalid patterns
- Test directory boundary cases
- Test with empty directories
- Test with hidden files
- Test with non-existent patterns

## Implementation Notes

- Use `filepath.Glob` from Go standard library
- Follow existing error handling patterns in codebase
- Add appropriate logging consistent with other tools
- Consider adding pattern validation/sanitization

## Specification by Example

```go
// Example call from agent:
{
    "pattern": "*.go"
}

// Example response:
[
    "main.go",
    "tool.go",
    "utils.go"
]

// Example call with nested pattern:
{
    "pattern": "**/*.md"
}

// Example response:
[
    "README.md",
    "docs/guide.md",
    "templates/prompts/system.md"
]
```

## Verification

- [ ] Tool accepts glob patterns and returns matching files
- [ ] Tool properly handles invalid patterns
- [ ] Tool remains within working directory scope
- [ ] Tool follows consistent error handling pattern
- [ ] Documentation includes clear examples of pattern usage
- [ ] Tool performs efficiently on various directory sizes
- [ ] Integration with existing tool framework is complete

## Next Steps

1. Create tool implementation in main.go
2. Add GlobInput struct
3. Generate schema using existing pattern
4. Implement Glob function
5. Add tool definition to main tools list
6. Add validation for patterns
7. Manual testing with various pattern types