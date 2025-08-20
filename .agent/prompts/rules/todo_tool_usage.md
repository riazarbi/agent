# Todo Tool Usage Rules

*Guidelines for using the todowrite and todoread tools for micro task management within agent sessions.*

## Context

**Applies to:** All agent interactions requiring task breakdown, progress tracking, or complex multi-step operations  
**Level:** Operational - enables real-time task orchestration  
**Audience:** AI agents working on complex tasks requiring structured progress tracking

## Core Principles

3. **Proactive Task Management:** Use todo tools when tasks become complex or multi-step
4. **Real-time Progress Tracking:** Update todos immediately as work progresses  
5. **Single Focus:** Only one task should be "in_progress" at any time
6. **Clear Distinction:** Micro todos (session) vs macro todos (TODO.md file)

## When to Use Todo Tools

### ✅ Use Todo Tools For:

**Complex Multi-Step Tasks:**
- Tasks requiring 3+ distinct operations
- Tasks with dependencies between steps
- Tasks where progress needs to be visible to user

**User-Provided Task Lists:**
- When users provide multiple tasks to complete
- When breaking down large requests into manageable chunks
- When users need progress visibility

**Planning and Coordination:**
- Before starting complex implementations
- When multiple files/systems need coordination
- When task order matters

**Examples:**
```
User: "Create a new API endpoint with validation, tests, and documentation"
→ Use todowrite to break this into specific steps

User: "Fix the login bug, update the docs, and deploy to staging"  
→ Use todowrite to track these as separate tasks

User: "Implement the user management feature"
→ Use todowrite to plan implementation steps
```

### ❌ Don't Use Todo Tools For:

**Simple Single-Step Tasks:**
- Reading a file
- Making a simple edit
- Basic search operations
- Single function implementations

**Conversational Requests:**
- Questions requiring explanations
- Information lookups
- Code reviews or analysis

**Trivial Operations:**
- Tasks providing no organizational benefit
- Operations that complete in <2 minutes

**Examples:**
```
User: "What's in the main.go file?"
→ Just use read_file, no todo needed

User: "Fix this typo in line 42"
→ Simple edit, no todo needed

User: "How does this function work?"
→ Explanation request, no todo needed
```

## Tool Usage Patterns

### Starting New Tasks

```json
{
  "todos": [
    {
      "id": "",
      "content": "Analyze current codebase structure", 
      "status": "in_progress",
      "priority": "high"
    },
    {
      "id": "",
      "content": "Design API endpoint schema",
      "status": "pending", 
      "priority": "high"
    },
    {
      "id": "",
      "content": "Implement validation logic",
      "status": "pending",
      "priority": "medium"
    },
    {
      "id": "",
      "content": "Write unit tests",
      "status": "pending",
      "priority": "medium"  
    },
    {
      "id": "",
      "content": "Update documentation",
      "status": "pending",
      "priority": "low"
    }
  ]
}
```

### Updating Progress

```json
{
  "todos": [
    {
      "id": "task-1",
      "content": "Analyze current codebase structure",
      "status": "completed", 
      "priority": "high"
    },
    {
      "id": "task-2", 
      "content": "Design API endpoint schema",
      "status": "in_progress",
      "priority": "high"
    },
    {
      "id": "task-3",
      "content": "Implement validation logic", 
      "status": "pending",
      "priority": "medium"
    }
  ]
}
```

## Status Management Rules

### Valid Status Values
- **pending**: Task is planned but not started
- **in_progress**: Task is currently being worked on (ONLY ONE allowed)
- **completed**: Task has been finished successfully  
- **cancelled**: Task was abandoned or became irrelevant

### Status Transition Guidelines
- **pending → in_progress**: When starting work on a task
- **in_progress → completed**: When task is successfully finished
- **in_progress → cancelled**: When task cannot be completed or becomes irrelevant
- **pending → cancelled**: When planned task is no longer needed
- **completed/cancelled**: Terminal states - don't change these

### Priority Guidelines
- **high**: Critical path items, blockers, user-facing issues
- **medium**: Important but not blocking, internal improvements
- **low**: Nice-to-have items, documentation updates, cleanup

## Workflow Integration

### With Existing TODO.md System

**Micro Todos (todowrite/todoread):**
- Session-scoped task tracking
- Real-time progress updates
- Operational task breakdown
- No file persistence required

**Macro Todos (TODO.md file):**
- Project-level task management  
- Cross-session persistence
- Requires user approval for changes
- Strategic planning and roadmaps

**Integration Pattern:**
1. Break down TODO.md items into micro todos for implementation
2. Complete micro todos during session
3. Update TODO.md with results at end of session

### Progress Reporting

Always read todos before reporting progress:
```
Use todoread to check current status, then report:
"Completed 3 of 5 tasks: ✅ Analysis ✅ Schema Design ✅ Validation 
Still pending: Tests and Documentation"
```

## Error Handling

### Common Validation Errors
- **Multiple in_progress tasks**: Only one task can be in_progress
- **Invalid status**: Must be pending/in_progress/completed/cancelled  
- **Invalid priority**: Must be high/medium/low
- **Empty content**: Task description cannot be empty

### Recovery Strategies
- Fix validation errors immediately with corrected todowrite
- If session state is corrupted, start fresh with new todo list
- Always validate before making status changes

## Best Practices

### Task Breakdown
- Keep individual tasks focused and specific
- Use action verbs: "Implement X", "Test Y", "Update Z"
- Avoid vague tasks like "Work on feature"
- Include enough context for future reference

### Progress Updates  
- Update status immediately after completing tasks
- Don't batch completions - update as you go
- Use in_progress status while actively working
- Mark cancelled if requirements change

### Session Management
- Start complex sessions with todowrite to set expectations
- Use todoread to check progress before major transitions
- End sessions by completing or properly transitioning active tasks

## Decision Framework

**Before using todo tools, ask:**
1. Is this task complex enough to benefit from tracking?
2. Are there multiple discrete steps involved?
3. Would the user benefit from seeing progress?
4. Am I likely to context-switch during this work?

**If yes to 2+ questions → Use todo tools**  
**If no → Complete task directly without todos**

## Quality Gates

- **Consistency**: Always update todos when status changes
- **Clarity**: Task descriptions should be self-explanatory
- **Completeness**: Don't leave tasks in ambiguous states
- **User Value**: Focus on todos that provide user visibility

---

## TL;DR

**Use todo tools when:**
- Tasks have 3+ steps or are complex
- Users provide multiple tasks  
- Progress visibility is valuable
- Work requires coordination

**Don't use for:**
- Simple single operations
- Conversational requests
- Trivial tasks

**Key Rules:**
- Only ONE task can be "in_progress" 
- Update status immediately as work progresses
- Use clear, action-oriented task descriptions
- Always read todos before reporting progress
