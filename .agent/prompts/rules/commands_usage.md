# Commands Usage Rules

*Rules for using list_commands and run_command tools for development workflow automation. Governs when and how to execute development commands (build, test, lint) to verify code works after changes. Commands are dynamically loaded from .agent/Taskfile.yml and can change during sessions.*

## Context

**Applies to:** All agent interactions involving code development, testing, building, and quality verification  
**Level:** Operational - critical for development workflow integrity and code quality assurance  
**Audience:** AI agents working on software development tasks requiring build/test/verification cycles

## Core Principles

1. **Verification-Oriented:** Primary purpose is to verify code works correctly after making changes
2. **Progressive Validation:** Follow Build → Test → Lint → Deploy progression for comprehensive verification
3. **Workflow-First:** Commands support development workflows rather than replacing manual verification
4. **Safe Automation:** Commands provide safe, explicit interfaces to development tools

## Rules

### Must Have (Critical)

- **RULE-001:** Always run `build` command after making any code changes to verify compilation
- **RULE-002:** Always run `test` command after code changes to verify functionality still works
- **RULE-003:** Use `list_commands()` regularly to discover available commands, as they can be added/removed during sessions
- **RULE-004:** NEVER modify .agent/Taskfile.yml - this is a protected configuration file that defines available commands

### Should Have (Important)

- **RULE-101:** Run `lint` or quality commands when available to maintain code standards
- **RULE-102:** Use commands before declaring development tasks complete
- **RULE-103:** Parse command output to understand and report failures clearly to users

### Could Have (Preferred)

- **RULE-201:** Use verbose flags (`-v`) when debugging command failures
- **RULE-202:** Run commands in logical order: build → test → lint → specialized
- **RULE-203:** Pass through user-requested flags via the `args` parameter

## Patterns & Anti-Patterns

### ✅ Do This

```javascript
// Good: Standard development verification workflow
// 1. Discover current available commands (they change dynamically)
list_commands()

// 2. After making code changes - verify it builds  
run_command({"command": "build"})

// 3. Run tests to ensure functionality works
run_command({"command": "test"})

// 4. Check code quality if available
run_command({"command": "lint"})
```

### ❌ Don't Do This

```javascript
// Bad: Skip verification after code changes
modify_code()
// ❌ No build verification!
// ❌ No test verification!
"The code changes are complete." // Unverified claim
```

## Decision Framework

**When to use development commands:**
1. Did I modify any code? → Use build/test commands
2. Does user want verification that "it works"? → Use appropriate verification commands  
3. Is this a development task? → Commands likely needed
4. Will user expect functionality to be working? → Verify with commands

**Command selection decision tree:**
1. **Always**: `list_commands()` regularly to discover current available options (commands are dynamic)
2. **Code Changes**: `build` → `test` → `lint` progression
3. **Bug Fixes**: Focus on `test` (possibly with specific test selection args)
4. **New Features**: Full `build` → `test` → `lint` verification
5. **Refactoring**: Comprehensive `build` → `test` → `lint` to ensure no regressions

## Exceptions & Waivers

**Valid reasons for not using commands:**
- Non-development tasks (reading/analyzing existing code, explanations)
- Documentation-only changes that don't affect code execution
- Information gathering or code review activities
- No `.agent/Taskfile.yml` exists (commands not available)

**NEVER modify these protected files:**
- `.agent/Taskfile.yml` - Command definitions (user/admin controlled)
- Commands are configured by users/administrators, not agents

**Process for exceptions:**
1. Clearly distinguish between code modification and analysis tasks
2. For borderline cases (config changes), err on the side of verification
3. Always explain to user when skipping verification and why

## Quality Gates

- **Build Verification:** Code must compile successfully before proceeding
- **Test Verification:** Tests must pass before declaring functionality complete  
- **Command Output Analysis:** Parse and understand command output, don't ignore failures
- **User Communication:** Report command results clearly, including both successes and failures

## Related Rules

- `rules/tool_usage.md` - Governs general tool invocation and say-do alignment
- `rules/todo_usage.md` - Governs task breakdown for complex development workflows

## Command Configuration

Commands are defined in `.agent/Taskfile.yml` using the [Task](https://taskfile.dev/) format:

```yaml
version: '3'

tasks:
  build:
    desc: "Build the project binary"
    cmds:
      - go build -o bin/app ./cmd/app

  test:
    desc: "Run all tests with race detection" 
    cmds:
      - go test -race ./...

  test-verbose:
    desc: "Run tests with verbose output and custom flags"
    cmds:
      - go test {{.CLI_ARGS}} ./...

  lint:
    desc: "Run golangci-lint for code quality"
    cmds:
      - golangci-lint run
```

**Key Requirements:**
- Commands **must** have a `desc` field to be discoverable
- Use `{{.CLI_ARGS}}` to accept arguments from `run_command`
- Commands without descriptions are considered internal and hidden

---

## TL;DR

**Key Principles:**
- Always verify code works after making changes
- Follow progressive validation: build → test → lint
- Use commands to support development workflows
- Check available commands regularly (they can change during sessions)

**Critical Rules:**
- Must run `build` command after code changes
- Must run `test` command after code changes  
- Must use `list_commands()` regularly to discover current available options (commands are dynamic)
- NEVER modify .agent/Taskfile.yml - it's protected configuration

**Quick Decision Guide:**
When in doubt: If you modified code, run build and test commands to verify it works