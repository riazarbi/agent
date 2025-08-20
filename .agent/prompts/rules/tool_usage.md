# Tool Usage Rules

*Rules for reliable and judicious use of available tools. Ensures tools are actually used when promised and used appropriately.*

## Context

**Applies to:** All agent interactions requiring file system access, search, or external operations  
**Level:** Operational - blocking for task completion  
**Audience:** AI agents working with this codebase

## Core Principles

1. **Say-Do Alignment:** If you say you'll use a tool, you must actually invoke it
2. **Tool First, Talk Second:** Use tools before describing what you found
3. **Judicious Selection:** Choose the right tool for the task, don't over-tool

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
- **RULE-202:** Explain tool choice when it might not be obvious
- **RULE-203:** Use tool output directly in responses rather than paraphrasing when precision matters

## Decision Framework

**When to use tools:**
1. You need current/accurate information about files or system state
2. User asks for specific file contents, searches, or operations
3. You're making claims that can be verified with available tools

**Tool selection priority:**
1. Most direct tool for the specific need
2. Tool that provides complete information in one call
3. Familiar tools over complex combinations

## Quality Gates

- **Self-check:** After writing a response, scan for any promises to use tools and verify you actually used them
- **Verification:** Tool outputs should directly support your response claims
- **Honesty:** If you didn't use a tool you said you would, admit it immediately and use it

---

## TL;DR

**Key Principles:**
- If you say you'll use a tool, actually use it
- Use tools first, then talk about results
- Choose tools deliberately, not habitually

**Critical Rules:**
- Must invoke promised tools with actual function calls
- Must use tools before making claims about their results
- Must verify tool outputs match your statements

**Judicious Use**
- When you don't know the size of a file, first use simple passthrough tools like `head`, `tail`, and `cloc` to understand its content and size.
- For searching patterns, `rg` is often the most efficient tool.
- For complex, atomic file modifications involving multiple changes, use `multi_edit`.
- Often you can answer your question using simple passthrough tools and `rg`.

**Quick Decision Guide:**
When in doubt: Use the tool and see what it actually shows
