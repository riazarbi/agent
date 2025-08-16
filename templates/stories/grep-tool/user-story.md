# [Implement Grep Tool Using Ripgrep]

*Enable the agent to perform fast text searches in the local directory using ripgrep (rg) as the underlying implementation.*

## Requirements

- Tool must accept a search pattern and return matching results
- Must implement using ripgrep (rg) for performance
- Must be non-modifying (read-only operations)
- Must follow existing tool pattern in codebase
- Must support both literal and regex patterns
- Must handle standard ripgrep arguments for flexibility

## Rules

- Tool must never modify any files
- Must handle both successful and error cases consistently with other tools
- Must validate input parameters before execution

## Domain

```
type GrepInput struct {
    Pattern string `json:"pattern" jsonschema_description:"The search pattern to look for (literal or regex)"`
    Args    []string `json:"args,omitempty" jsonschema_description:"Optional ripgrep arguments (e.g. --ignore-case, --hidden)"`
}

var GrepDefinition = ToolDefinition{
    Name: "grep",
    Description: "Search for patterns in files using ripgrep. Supports both literal and regex patterns.",
    InputSchema: GrepInputSchema,
    Function: Grep,
}
```

## Implementation Notes

- Use os/exec to call ripgrep
- Follow existing error handling patterns in codebase
- Add appropriate logging consistent with other tools

## Specification by Example

```go
// Example call from agent:
{
    "pattern": "TODO:",
    "args": ["--ignore-case"]
}

// Example response:
main.go:42:TODO: Implement error handling
docs/notes.md:15:todo: update documentation
```

## Next Steps

1. Verify ripgrep is installed and accessible
2. Add GrepInput struct
3. Generate schema using existing pattern
4. Implement Grep function
5. Add tool definition to main tools list
6. Manual testing with various patterns and arguments