# Fix Missing Tool Registration

*Several essential tools (todowrite, todoread, list_files) are implemented but not available to the agent because they are not properly registered with the tool registry.*

## Past Attempts

Based on the git status, there has been ongoing refactoring work moving tools from main.go to the internal tools package. The tools exist in both locations but the registration logic in main.go at line 518 (NewAgent function) is not properly configuring the tool registry with session dependencies needed for todo tools.

## Requirements

- The `todowrite` tool must be available and functional for managing todo lists during agent sessions
- The `todoread` tool must be available and functional for reading todo lists from the current session  
- The `list_files` tool must be available and functional for listing files and directories
- All tools must be properly registered with the tool registry and appear in the agent's available tools
- Session-dependent tools (todowrite, todoread) must receive proper session manager configuration
- The fix must not break any existing working tools

## Rules

- Must use the existing internal/tools registry architecture
- Must provide proper session dependencies for todo tools through RegistryConfig
- Must not duplicate tool implementations between main.go and internal/tools
- Must follow the existing pattern of tool factories (NewFileTools, NewWebTools, etc.)

## Domain

```go
// Tool registry manages available tools
type Registry struct {
    tools map[string]Tool
}

// Session-dependent tools need configuration
type RegistryConfig struct {
    SessionManager   session.SessionManager
    CurrentSessionID string
}

// Tools are registered via factory functions
func NewFileTools() []Tool
func NewTodoTools(SessionManager, string) []Tool
```

## Extra Considerations

- The main.go file currently has duplicate tool implementations that should be removed after internal tools are working
- Session management integration needs to be verified to ensure todo tools persist data correctly
- The tool registry initialization happens during agent construction and needs session context
- Missing GenerateSchema function needs to be available in the tools package for schema generation

## Testing Considerations

- Test that all three tools appear in the agent's tool list
- Test todowrite can create and update todos in a session
- Test todoread can retrieve todos from the current session
- Test list_files properly lists files and directories
- Test session isolation for todo tools (different sessions have different todo lists)
- Verify existing tools continue to work after changes

## Implementation Notes

- The issue is in main.go:518 where `internaltools.NewRegistry(nil)` is called without proper RegistryConfig
- Session management is already initialized before agent creation, so dependencies are available  
- Need to ensure GenerateSchema function is available for tool schema generation
- Should move list_files from external tools package to internal file tools

## Specification by Example

After fix, agent should show these tools when asked about available tools:
- `read_file` - Read file contents
- `edit_file` - Edit file with string replacement  
- `delete_file` - Delete a file
- `list_files` - List files and directories
- `todowrite` - Create/manage todo lists
- `todoread` - Read current todo list
- `grep` - Search files with patterns
- `git_diff` - Show git changes
- `web_fetch` - Download web content
- And other existing tools...

## Verification

- [ ] Agent shows all expected tools in its registry (including todowrite, todoread, list_files)
- [ ] todowrite creates and updates todo lists successfully
- [ ] todoread retrieves todo lists from current session
- [ ] list_files lists files and directories correctly
- [ ] Session isolation works for todo tools (different sessions maintain separate todo lists)
- [ ] All existing tools continue to work after changes
- [ ] No duplicate tool implementations remain in codebase
- [ ] Tool registry properly handles session-dependent tool configuration