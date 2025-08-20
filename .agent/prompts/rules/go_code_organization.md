# Go Code Organization for Agent-Friendly Codebases

*Rules for organizing Go code to optimize for LLM agent interaction patterns. These rules prioritize discoverability, modularity, and efficient navigation using tools like grep, glob, and list_files before requiring expensive read_file operations.*

## Context

*Code organization optimized for agents that interact with codebases through tool-based exploration rather than IDE navigation.*

**Applies to:** Go projects where LLM agents are primary contributors, refactoring efforts, new feature development  
**Level:** Strategic - influences overall codebase structure and agent workflow efficiency  
**Audience:** Developers and LLM agents working on Go codebases

## Core Principles

1. **Agent Discoverability:** Code structure should allow agents to locate relevant functionality using grep, glob, and list_files without reading entire files
2. **Token Efficiency:** Minimize the need for expensive read_file operations by making file purposes clear from names and structure
3. **Focused Responsibility:** Each file should have a single, clear purpose that can be inferred from its name and location
4. **Search-Friendly Patterns:** Use consistent naming conventions that enable precise grep patterns to locate specific functionality

## Rules

### Must Have (Critical)
*Non-negotiable rules that must always be followed. Violation blocks efficient agent interaction.*

- **ORG-001:** No single file should exceed 300 lines - forces focused responsibility and enables complete file reading
- **ORG-002:** File names must clearly indicate their primary responsibility (e.g., `tool_definitions.go`, `agent_runner.go`, `cli_flags.go`)  
- **ORG-003:** All exported functions, types, and variables must have comments that include searchable keywords for their domain
- **ORG-004:** Related functionality must be grouped into packages with clear boundaries (tools, agent, cli, etc.)

### Should Have (Important)
*Strong recommendations that significantly improve agent workflow efficiency.*

- **ORG-101:** Function and type names should include domain keywords that agents can grep for (e.g., `ToolDefinition`, `ExecuteTool`, `AgentConversation`)
- **ORG-102:** Package-level documentation should list the main types and functions with one-line descriptions
- **ORG-103:** Constants and variables should be grouped by purpose with clear section comments
- **ORG-104:** Interface definitions should be in separate files from their implementations to enable focused reading

### Could Have (Preferred)
*Best practices that further optimize agent interaction.*

- **ORG-201:** Include `// Keywords: ` comments on complex functions listing searchable terms
- **ORG-202:** Use consistent file naming patterns within packages (`*_definition.go`, `*_implementation.go`, `*_test.go`)
- **ORG-203:** Group imports by standard library, third-party, and local packages with blank line separation

## Patterns & Anti-Patterns

### ✅ Do This
*File structure that enables efficient agent navigation*

```go
// File: tools/definitions.go
// Purpose: Tool definitions and schemas
// Keywords: tool, definition, schema, input, validation

package tools

// ToolDefinition defines the structure for all agent tools
// Keywords: tool, definition, interface, schema
type ToolDefinition struct {
    Name        string
    Description string
    InputSchema openai.FunctionParameters
    Function    func(input json.RawMessage) (string, error)
}

// ReadFileDefinition provides file reading capabilities
// Keywords: read, file, filesystem, content
var ReadFileDefinition = ToolDefinition{...}
```

```go
// File: agent/runner.go  
// Purpose: Core agent execution logic
// Keywords: agent, run, conversation, inference

package agent

// Agent handles conversation flow and tool execution
// Keywords: agent, conversation, tools, execution
type Agent struct {...}
```

### ❌ Don't Do This
*Monolithic organization that requires full file reading*

```go
// File: main.go (800+ lines)
// Contains: types, schemas, tools, agent logic, CLI parsing, utilities
// Problem: Agent must read entire file to understand any single component

package main

// Mixed concerns make it impossible to grep for specific functionality
type Agent struct {...}          // Line 45
var ReadFileDefinition = ...     // Line 123  
func main() {...}                // Line 234
func ReadFile(...) {...}         // Line 567
func GenerateSchema[T any]() ... // Line 678
```

## Decision Framework

*Guidance for organizing code when facing structural decisions*

**When adding new functionality:**
1. Can this be found via grep without reading unrelated code? If no, create focused file
2. Does this file exceed 300 lines? If yes, split by responsibility
3. Would an agent need to read multiple files to understand this feature? If yes, consider consolidation

**When refactoring existing code:**
1. Group by agent workflow - what would an agent need to read together?
2. Split mixed concerns into separate, greppable files  
3. Ensure file names clearly indicate contents for efficient navigation

## Exceptions & Waivers

*When these rules may be relaxed*

**Valid reasons for exceptions:**
- Generated code that shouldn't be manually modified
- Temporary refactoring states (document timeline for cleanup)
- Integration with third-party libraries requiring specific organization

**Process for exceptions:**
1. Document the exception reason and timeline in package comment
2. Include explicit keywords in comments to maintain discoverability
3. Plan remediation if exception creates agent workflow friction

## Quality Gates

*Verification of agent-friendly organization*

- **Automated checks:** Line count limits, file naming patterns, required package comments
- **Agent testing:** Can sample workflows be completed without reading >3 files per feature?  
- **Grep coverage:** All major functionality should be findable via targeted grep patterns

## Related Rules

*Complementary rule sets for complete agent-codebase interaction*

- `rules/tool_usage.md` - How agents should interact with available tools
- `rules/testing_patterns.md` - Test organization that supports agent-driven development
- `rules/api_design.md` - Interface design for agent-friendly APIs

## References

*Resources informing these agent-centric organization patterns*

- [Go Package Layout Standards](https://github.com/golang-standards/project-layout) - Traditional Go organization  
- [Effective Go](https://golang.org/doc/effective_go.html) - Go language conventions
- [Agent Tool Usage Patterns](internal) - Analysis of efficient agent workflows

---

## TL;DR

*Essential rules for agent-friendly Go code organization*

**Key Principles:**
- Structure code for discoverability via grep/glob, not just human reading
- Keep files focused and small enough for complete agent reading
- Use predictable naming that enables efficient agent navigation
- Optimize for agent workflow efficiency over traditional organization

**Critical Rules:**
- Max 300 lines per file to enable complete reading
- File names must clearly indicate purpose and contents  
- All exports must have searchable keyword comments
- Group related functionality into focused packages

**Quick Decision Guide:**
When in doubt: Can an agent find and understand this functionality without reading unrelated code?