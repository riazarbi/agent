# [Implement Git Diff Tool]

*Enable the agent to see changes made since the last commit, helping it understand what has been modified in the working directory.*

## Requirements

- Tool must return the git diff output for unstaged changes
- Must be non-modifying (read-only operations)
- Must follow existing tool pattern in codebase
- Must handle error cases gracefully and return meaningful messages
- Must not exit chat on errors

## Rules

- Tool must never modify any files
- Must be read-only operation
- Must handle both successful and error cases consistently with other tools
- Must return clear error messages when:
  - Not in a git repository
  - No changes exist
  - Git command fails

## Domain

```
type GitDiffInput struct {
    // This tool takes no parameters
}

var GitDiffDefinition = ToolDefinition{
    Name: "git_diff",
    Description: "Returns the output of 'git diff' showing all unstaged changes in the working directory",
    InputSchema: GitDiffInputSchema,
    Function: GitDiff,
}
```

## Extra Considerations

- Return full diff output regardless of size - agent can handle parsing large content
- For binary files, return the standard git message "Binary files a/file.bin and b/file.bin differ"
- Use UTF-8 encoding for all output

## Implementation Notes

- Use os/exec to execute "git diff" command with no additional flags
- Return complete, unmodified diff output exactly as git produces it
- Implement using standard error handling: return error string on failure, diff output on success
- Add logging for command execution and error cases matching grep tool implementation

## Specification by Example

```go
// Example call from agent:
<invoke name="git_diff">
</invoke>

// Example successful response with changes:
"diff --git a/main.go b/main.go
index abc123..def456 100644
--- a/main.go
+++ b/main.go
@@ -10,6 +10,7 @@ func main() {
     // New line added
+    fmt.Println("Hello")
 }"

// Example response when no changes:
"No changes found in working directory"

// Example error response:
"Error: not a git repository"
```

## Verification

- [ ] Tool returns correct diff output when there are changes
- [ ] Tool handles "no changes" case gracefully
- [ ] Tool handles "not a git repository" error appropriately
- [ ] Output format is useful for agent comprehension
- [ ] Error messages are clear and informative
- [ ] Implementation follows existing codebase patterns

## Next Steps

1. Create GitDiff function
2. Generate schema using existing pattern
3. Add tool definition to main tools list
4. Add error handling for common cases
5. Manual testing with various scenarios