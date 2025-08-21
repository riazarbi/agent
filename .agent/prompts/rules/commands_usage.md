# Commands Usage Rules

*Rules for using `xc` as the primary tool for user-defined command execution and development workflow automation. Governs when and how to execute development commands (build, test, lint) to verify code works after changes. Commands are dynamically loaded from .agent/Taskfile.yml and can change during sessions.*

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

- **RULE-001:** Always run `xc build` after making any code changes to verify compilation
- **RULE-002:** Always run `xc test` after code changes to verify functionality still works
- **RULE-003:** Use `xc(args="")` regularly to discover available commands, as they can be added/removed during sessions
- **RULE-004:** NEVER modify .agent/Taskfile.yml - this is a protected configuration file that defines available xc commands

### Should Have (Important)

- **RULE-101:** Run `xc lint` or other quality commands when available to maintain code standards
- **RULE-102:** Use `xc` commands before declaring development tasks complete
- **RULE-103:** Parse `xc` command output to understand and report failures clearly to users

### Could Have (Preferred)

- **RULE-201:** Use verbose flags (`-v`) when debugging command failures
- **RULE-202:** Run commands in logical order: build → test → lint → specialized
- **RULE-203:** Pass through user-requested flags via the `args` parameter

## Patterns & Anti-Patterns

### ✅ Do This

```javascript
// Good: Standard development verification workflow
// 1. Discover current available commands (they change dynamically)
xc(args="")

// 2. After making code changes - verify it builds  
xc(args="build")

// 3. Run tests to ensure functionality works
xc(args="test")

// 4. Check code quality if available
xc(args="lint")
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
1. **Always**: `xc(args="")` regularly to discover current available options (commands are dynamic)
2. **Code Changes**: `xc build` → `xc test` → `xc lint` progression
3. **Bug Fixes**: Focus on `xc test` (possibly with specific test selection args)
4. **New Features**: Full `xc build` → `xc test` → `xc lint` verification
5. **Refactoring**: Comprehensive `xc build` → `xc test` → `xc lint` to ensure no regressions

## Exceptions & Waivers

**Valid reasons for not using commands:**
- Non-development tasks (reading/analyzing existing code, explanations)
- Documentation-only changes that don't affect code execution
- Information gathering or code review activities
- No `README.md` exists (xc commands not available)

**NEVER modify these protected files:**
- `README.md` - Command definitions (user/admin controlled)
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

Commands are defined in `README.md` using a format compatible with `xc`. For full documentation on defining commands, refer to the [xcfile.dev](https://xcfile.dev) website.

The task specifications can be found at https://github.com/joerdav/xc/tree/main/doc/content/task-syntax.

--- BEGIN EXAMPLE ---

## Tasks

### deploy

Requires: test
Directory: ./deployment
Env: ENVIRONMENT=STAGING

```
sh deploy.sh
```

--- END EXAMPLE ---


## TL;DR

**Key Principles:**
- Always verify code works after making changes
- Follow progressive validation: build → test → lint
- Use commands to support development workflows
- Check available commands regularly (they can change during sessions)

**Critical Rules:**
- Must run `xc build` command after code changes
- Must run `xc test` command after code changes
- Must use `xc(args="")` regularly to discover current available options (commands are dynamic)
- NEVER modify README.md - it's protected configuration

**Quick Decision Guide:**
When in doubt: If you modified code, run build and test commands to verify it works