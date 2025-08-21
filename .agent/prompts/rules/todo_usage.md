# Todo Tool Usage Rules

*Rules for using todowrite and todoread tools for micro task management within agent sessions. CRITICAL: Always evaluate instructions and create todos BEFORE starting any non-trivial work. Default to using todos - only skip for genuinely simple operations.*

## Context

**Applies to:** All agent interactions requiring task breakdown, progress tracking, or complex multi-step operations  
**Level:** Operational - enables real-time task orchestration and user visibility  
**Audience:** AI agents working on complex tasks requiring structured progress tracking

## Core Principles

1. **Plan Before Execute:** Always evaluate instructions and create todos for non-trivial tasks before starting work
2. **Proactive Task Management:** Default to using todos - only skip for genuinely trivial operations
3. **Single Focus:** Only one task should be "in_progress" at any time to maintain clear execution context
4. **Real-time Tracking:** Update todo status immediately as work progresses, never batch updates
5. **Clear Distinction:** Micro todos (session-scoped) serve different purposes than macro todos (project-level TODO.md files)

## Rules

### Must Have (Critical)

- **RULE-001:** ALWAYS pause to evaluate any instruction - if it cannot be done trivially (in 1-2 simple steps), create todos BEFORE starting work
- **RULE-002:** Only ONE task can have status "in_progress" at any given time
- **RULE-003:** Always update todo status immediately after completing or changing tasks - never batch updates
- **RULE-004:** Use todoread before reporting progress to ensure accuracy

### Should Have (Important)

- **RULE-101:** When in doubt about todo usage, err on the side of creating todos - over-planning is better than under-planning
- **RULE-102:** Break down user-provided task lists into individual trackable items
- **RULE-103:** Include enough context in task descriptions to be self-explanatory

### Could Have (Preferred)

- **RULE-201:** Use action verbs in task descriptions: "Implement X", "Test Y", "Update Z"
- **RULE-202:** Set appropriate priorities (high/medium/low) based on critical path and user impact
- **RULE-203:** Mark tasks as cancelled rather than deleting when requirements change

## Patterns & Anti-Patterns

### ✅ Do This

```javascript
// Good: Evaluate instruction, then create todos BEFORE starting work
User: "Add error handling to the login function and test it"

// 1. EVALUATE: This involves code changes + testing = NOT trivial
// 2. CREATE TODOS FIRST:
todowrite({
  "todos": [
    {
      "id": "task-1",
      "content": "Analyze current login function implementation",
      "status": "in_progress",
      "priority": "high"
    },
    {
      "id": "task-2", 
      "content": "Add error handling to login function",
      "status": "pending",
      "priority": "high"
    },
    {
      "id": "task-3",
      "content": "Write tests for error handling scenarios",
      "status": "pending",
      "priority": "high"
    },
    {
      "id": "task-4",
      "content": "Run tests to verify implementation",
      "status": "pending",
      "priority": "medium"
    }
  ]
})
// 3. NOW start work on first task
```

### ❌ Don't Do This

```javascript
// Bad: Jump straight into complex work without planning
User: "Add error handling to the login function and test it"

// ❌ WRONG: Start working immediately without todos
read_file("login.js")
edit_file("login.js", ...)
// ❌ No planning, no progress tracking, no user visibility

// Bad: Create todos after already starting work
edit_file("login.js", ...)  // ❌ Already started!
todowrite({...})           // ❌ Too late!
```

## Decision Framework

**MANDATORY EVALUATION PROCESS:**
Before starting ANY task, ask these questions in order:

1. **Can this be done in 1-2 trivial steps?** (e.g., read a single file, answer a simple question)
   - If YES → Proceed without todos
   - If NO → MUST create todos before starting

2. **Does this involve any of the following?**
   - Multiple file changes
   - Code implementation + testing
   - Research + implementation
   - Debugging + fixing
   - Multiple tool invocations
   - Any planning or coordination
   - If YES to ANY → MUST create todos

**Task complexity indicators requiring todos:**
- User says "implement", "fix", "build", "create", "update", "refactor"
- Task involves verification steps (build, test, lint)
- Task requires understanding something before doing it
- Multiple components or files involved
- Any uncertainty about approach or steps

**When facing task breakdown decisions:**
- Default to creating todos unless task is genuinely trivial
- Prefer smaller, focused tasks over large general ones
- Include dependencies in task ordering (high priority for blockers)
- Consider user visibility when setting priorities

## Exceptions & Waivers

**Valid reasons for NOT using todo tools (very limited):**
- Reading a single specific file that user requested
- Answering a direct question about existing code/system
- Simple explanations that require no file access
- Basic search operations with immediate answers

**Process for exceptions:**
1. Evaluate using the MANDATORY EVALUATION PROCESS above
2. If there's ANY doubt, create todos - this is required
3. Remember: Over-planning is better than under-planning
4. User visibility and progress tracking are valuable even for smaller tasks

## Quality Gates

- **Status Validation:** Ensure only one task is ever "in_progress"
- **Progress Accuracy:** Todo status must accurately reflect actual work state
- **User Visibility:** Todos should provide meaningful progress updates to users
- **Task Clarity:** Each task description should be actionable and specific

## Related Rules

- `rules/tool_usage.md` - Governs when and how to use tools in general
- `rules/commands_usage.md` - Governs development workflow verification

---

## TL;DR

**Key Principles:**
- ALWAYS evaluate instructions before starting work
- Default to using todos unless task is genuinely trivial (1-2 simple steps)
- Plan before execute - todos first, then work
- Over-planning is better than under-planning

**Critical Rules:**
- Must evaluate ALL instructions and create todos for non-trivial tasks BEFORE starting
- Must maintain single in_progress task constraint
- Must update status in real-time, not batched

**Quick Decision Guide:**
When in doubt: CREATE TODOS. If you're questioning whether to use todos, the answer is YES.