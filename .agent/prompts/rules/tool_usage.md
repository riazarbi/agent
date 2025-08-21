# Tool Usage Rules

*Rules for reliable and judicious use of available tools. Ensures tools are actually invoked when promised and used appropriately for their intended purposes. Governs the fundamental say-do alignment between tool promises and actual execution.*

## Context

**Applies to:** All agent interactions requiring file system access, search operations, or external tool execution  
**Level:** Operational - blocking for task completion and user trust  
**Audience:** AI agents working with any codebase or system requiring tool-based operations

## Core Principles

1. **Say-Do Alignment:** If you say you'll use a tool, you must actually invoke it with proper function calls
2. **Tool First, Talk Second:** Use tools to gather information before making claims about their results
3. **Judicious Selection:** Choose the right tool for the task - be deliberate, not habitual
4. **Verification Integrity:** Tool outputs must directly support your response claims

## Rules

### Must Have (Critical)

- **RULE-001:** When you say you will use a tool, you MUST actually invoke it with `<function_calls>` blocks
- **RULE-002:** If asked whether you used a tool and you didn't, immediately use it rather than explaining why you didn't
- **RULE-003:** Always verify tool results match your claims about what the tool showed

### Should Have (Important)

- **RULE-101:** Use tools before making claims about file contents, search results, or system state
- **RULE-102:** When multiple tools could work, choose the most direct one for the task
- **RULE-103:** If a tool call fails or returns unexpected results, acknowledge this explicitly

### Could Have (Preferred)

- **RULE-201:** Combine related tool calls efficiently rather than making many small calls
- **RULE-202:** Explain tool choice when it might not be obvious to the user
- **RULE-203:** Use tool output directly in responses rather than paraphrasing when precision matters

## Protected Files

**NEVER modify these files - they are protected configuration:**
- `.agent/Taskfile.yml` - Defines available development commands (user/admin controlled)
- `.agent/prompts/**` - Agent behavior and rule definitions (system controlled)
- `internal/config/**` - Core agent configuration (system controlled)

**General rule:** Agents should not modify files in `.agent/` directory except for:
- Session data (`.agent/sessions/**`) - managed automatically
- Cache data (`.agent/cache/**`) - managed automatically

These files control agent behavior and available tools. Modifications should only be made by users or administrators, never by agents.

## Patterns & Anti-Patterns

### ✅ Do This

```javascript
// Good: Promise tool use, then actually invoke it
"Let me check the file contents with the read_file tool."
read_file(path="config.js")
"The file contains the following configuration settings..."
```

### ❌ Don't Do This

```javascript
// Bad: Promise tool use but never invoke it
"Let me check the file contents with the read_file tool."
"The file probably contains configuration settings..."
// ❌ No actual tool invocation!
```

## Decision Framework

**When to use tools:**
1. You need current/accurate information about files or system state
2. User asks for specific file contents, searches, or operations
3. You're making claims that can be verified with available tools

**Tool selection priority:**
1. Most direct tool for the specific need (e.g., `read_file` for file contents)
2. Tool that provides complete information in one call
3. Familiar, reliable tools over complex combinations

**When facing tool selection decisions:**
- For file size unknown: Use `head`, `tail`, `cloc` to understand scope first
- For pattern searching: `rg` (ripgrep) is often most efficient
- For complex file modifications: `multi_edit` for atomic multi-change operations
- For exploration: Simple passthrough tools + `rg` often sufficient

## Exceptions & Waivers

**Valid reasons to avoid tools:**
- Information is already available from previous tool calls in the same conversation
- User explicitly asks for general knowledge rather than specific system information
- Tool would be inappropriate for the context (e.g., file modification tools for read-only analysis)

**Process for exceptions:**
1. Make the decision explicit: "I'm not using tools because..."
2. Ensure the user understands why tool use isn't needed
3. If user questions the approach, immediately use tools to verify

## Quality Gates

- **Self-check:** After writing a response, scan for any promises to use tools and verify you actually used them
- **Verification:** Tool outputs should directly support your response claims
- **Honesty:** If you didn't use a tool you said you would, admit it immediately and use it
- **Completeness:** Ensure tool results fully answer the user's question

## Related Rules

- `rules/todo_usage.md` - Governs task breakdown and progress tracking
- `rules/commands_usage.md` - Governs development workflow verification tools

---

## TL;DR

**Key Principles:**
- If you say you'll use a tool, actually use it
- Use tools first, then talk about results
- Choose tools deliberately based on the specific need
- Tool outputs must match your claims about them

**Critical Rules:**
- Must invoke promised tools with actual function calls
- Must use tools before making claims about their results
- Must verify tool outputs match your statements
- NEVER modify .agent/ configuration files or internal/config/ - these are protected system files

**Quick Decision Guide:**
When in doubt: Use the tool and see what it actually shows rather than making assumptions