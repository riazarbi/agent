# Implement todowrite and todoread Tools for Micro Task Management

*Enable agents to create and manage structured task lists for complex multi-step operations within a single session, providing real-time task orchestration capabilities.*

## Requirements

- Implement `todowrite` tool that accepts an array of todo items and stores them in session state
- Implement `todoread` tool that returns the current todo list from session state
- Todo items must include: id, content, status, priority fields
- Status must support: pending, in_progress, completed, cancelled
- Priority must support: high, medium, low
- Tools must follow existing ToolDefinition pattern in codebase
- Must use in-memory storage scoped to agent session (no persistence required)
- Must return JSON-formatted output consistent with other tools
- Must include comprehensive usage guidance in dedicated rules file

## Rules

- Tools must never modify files directly - they only manage in-memory state
- Only one task should be "in_progress" at any given time
- Must validate input parameters before execution
- Must handle both successful and error cases consistently with other tools
- Session state must be isolated (no cross-session contamination)

## Domain

```go
type TodoItem struct {
    ID       string `json:"id" jsonschema_description:"Unique identifier for the todo item"`
    Content  string `json:"content" jsonschema_description:"Brief description of the task"`
    Status   string `json:"status" jsonschema_description:"Current status: pending, in_progress, completed, cancelled"`
    Priority string `json:"priority" jsonschema_description:"Priority level: high, medium, low"`
}

type TodoWriteInput struct {
    Todos []TodoItem `json:"todos" jsonschema_description:"The updated todo list"`
}

type TodoReadInput struct {
    // No parameters needed
}

var TodoWriteDefinition = ToolDefinition{
    Name: "todowrite",
    Description: "Create and manage structured task lists for complex multi-step operations within the current session.",
    InputSchema: TodoWriteInputSchema,
    Function: TodoWrite,
}

var TodoReadDefinition = ToolDefinition{
    Name: "todoread", 
    Description: "Read the current todo list from session state.",
    InputSchema: TodoReadInputSchema,
    Function: TodoRead,
}
```

## Intended Usage

### When to Use Todo Tools

**Use proactively for:**
1. **Complex multi-step tasks** - When a task requires 3+ distinct steps or operations
2. **User provides multiple tasks** - When users provide lists of things to be done
3. **Non-trivial operations** - Tasks requiring careful planning or coordination
4. **After receiving new instructions** - Immediately capture user requirements as todos
5. **Progress tracking** - When users need visibility into task completion status

### When NOT to Use Todo Tools

**Skip for:**
1. **Single, straightforward tasks** - Tasks completable in 1-2 simple steps
2. **Trivial operations** - Tasks providing no organizational benefit to track
3. **Purely conversational requests** - Information-only or explanation requests
4. **Quick file operations** - Simple read/edit/search operations

### Task Management Workflow

1. **Create todos** with `todowrite` when starting complex work
2. **Mark in_progress** when actively working (limit to ONE task)
3. **Complete tasks immediately** after finishing - don't batch completions
4. **Update status in real-time** as work progresses
5. **Cancel tasks** that become irrelevant due to changing requirements

### Integration with Existing Systems

- **Micro tasks** (todowrite/todoread): Within-session operational tracking
- **Macro tasks** (TODO.md): Cross-session project management requiring user approval
- Use micro tasks to break down macro tasks into executable steps

## Extra Considerations

- Must implement session-scoped state management (in-memory map by session ID)
- Tool output should include count of non-completed todos in title
- Must generate UUIDs or similar for todo item IDs if not provided
- Should validate status and priority values against allowed enums
- Consider thread safety if agent supports concurrent operations

## Testing Considerations

- Test creating, reading, and updating todo lists
- Test status transitions and validation
- Test priority assignment and validation
- Test session isolation (multiple agent instances)
- Test edge cases: empty lists, invalid status/priority values
- Manual verification of JSON output format

## Implementation Notes

- Follow existing error handling patterns in codebase
- Use consistent logging with other tools
- Add tools to main tools list in main.go
- Consider using sync.RWMutex for thread-safe state access
- Generate schemas using existing GenerateSchema pattern

## Specification by Example

### TodoWrite Example
```json
{
  "todos": [
    {
      "id": "task-001", 
      "content": "Implement user authentication endpoint",
      "status": "in_progress",
      "priority": "high"
    },
    {
      "id": "task-002",
      "content": "Add input validation to registration form", 
      "status": "pending",
      "priority": "medium"
    },
    {
      "id": "task-003",
      "content": "Update API documentation",
      "status": "completed", 
      "priority": "low"
    }
  ]
}
```

### TodoRead Response
```json
{
  "title": "2 todos",
  "output": "[{\"id\":\"task-001\",\"content\":\"Implement user authentication endpoint\",\"status\":\"in_progress\",\"priority\":\"high\"},{\"id\":\"task-002\",\"content\":\"Add input validation to registration form\",\"status\":\"pending\",\"priority\":\"medium\"},{\"id\":\"task-003\",\"content\":\"Update API documentation\",\"status\":\"completed\",\"priority\":\"low\"}]"
}
```

## Verification

- [ ] TodoWriteInput and TodoReadInput structs defined with proper JSON schema tags
- [ ] TodoWrite and TodoRead functions implemented with error handling
- [ ] Schema generation using GenerateSchema pattern
- [ ] Tool definitions added to main tools list
- [ ] Session-scoped state management implemented
- [ ] Status and priority validation implemented
- [ ] JSON output format matches other tools
- [ ] Thread safety considerations addressed
- [ ] Comprehensive usage guidance file created at `.agent/prompts/rules/todo_tool_usage.md`
- [ ] Manual testing with various scenarios completed
- [ ] Integration testing with existing tool ecosystem

## Next Steps

After creating the story:
1. Save to `.agent/stories/todo-tools/user-story.md`
2. Create comprehensive usage guidance at `.agent/prompts/rules/todo_tool_usage.md` 
3. Review reference implementation patterns for behavioral guidance
4. Plan implementation approach for session state management