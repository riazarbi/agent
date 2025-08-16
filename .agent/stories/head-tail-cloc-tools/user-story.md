# Add head, tail, and cloc Tools

*Add three new command-line tools (head, tail, cloc) to the agent's toolkit to improve file analysis and code inspection capabilities.*

## Requirements

**Must have:**
- Add `head` tool that shows first N lines of a file (default 10 lines)
- Add `tail` tool that shows last N lines of a file (default 10 lines)  
- Add `cloc` tool that counts lines of code with language breakdown
- All tools accept flexible arguments like the existing `grep` tool implementation
- Tools are read-only operations that cannot alter files
- Tools return proper error messages for missing executables or invalid arguments
- Add all three tools to the main tools array so they're available to the agent

**Should have:**
- Follow the same implementation pattern as existing tools (GrepInput structure, error handling, etc.)
- Support common arguments for each tool (e.g., head -n 20, tail -f, cloc --exclude-dir)
- Provide helpful descriptions that explain when to use each tool
- Handle cases where the underlying command-line tools are not installed

## Rules

- Tools must be read-only (no file modification capabilities)
- Must follow existing code patterns for consistency
- Error handling should match patterns used in Grep and GitDiff tools
- All tools added to main tools array: `ReadFileDefinition, ListFilesDefinition, ..., HeadDefinition, TailDefinition, ClocDefinition`

## Domain

```go
type HeadInput struct {
    Args []string `json:"args,omitempty" jsonschema_description:"Optional head arguments (e.g. -n 20, filename)"`
}

type TailInput struct {
    Args []string `json:"args,omitempty" jsonschema_description:"Optional tail arguments (e.g. -n 20, -f, filename)"`
}

type ClocInput struct {
    Args []string `json:"args,omitempty" jsonschema_description:"Optional cloc arguments (e.g. --exclude-dir=.git, path)"`
}
```

## Extra Considerations

- `head` and `tail` typically require filenames as arguments
- `cloc` can analyze directories or specific files
- `tail -f` follows files but may not be useful in this context (agent interaction)
- Commands should handle both stdout and stderr appropriately
- Exit codes should be handled (similar to how grep handles exit code 1 for no matches)

## Testing Considerations

- Test with missing executable (command not found)
- Test with invalid arguments 
- Test with non-existent files
- Test basic functionality (head -n 5 file.txt, tail -n 10 file.txt, cloc .)
- Test empty args array (should use defaults)

## Implementation Notes

- Follow the exact same pattern as `Grep` function implementation
- Use `exec.Command` with the tool name as first parameter
- Capture stdout and stderr with bytes.Buffer
- Handle exit errors appropriately 
- Each tool needs: Input struct, Schema generation, ToolDefinition, and implementation function
- Add schema variables at the top with existing ones: `HeadInputSchema`, `TailInputSchema`, `ClocInputSchema`

## Specification by Example

**Head tool usage:**
```json
{
  "args": ["-n", "5", "main.go"]
}
```

**Tail tool usage:**
```json
{
  "args": ["-n", "20", "agent.log"]
}
```

**Cloc tool usage:**
```json
{
  "args": ["--exclude-dir=node_modules", "."]
}
```

## Verification

- [ ] HeadInput, TailInput, and ClocInput structs defined
- [ ] Schema generation variables added for all three tools
- [ ] ToolDefinition variables created for all three tools  
- [ ] Head, Tail, and Cloc functions implemented following Grep pattern
- [ ] All three tools added to main tools array
- [ ] Tools handle missing executables gracefully
- [ ] Tools handle invalid arguments with proper error messages
- [ ] Basic functionality verified (head shows first lines, tail shows last lines, cloc counts code)

## Next Steps

After creating the story, save it and proceed with implementation following the existing patterns in main.go.