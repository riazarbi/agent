# Agent Implementation Instructions

## Purpose

These instructions guide a coding agent to incrementally implement the Minimal AI Coding Agent project. Each agent MUST complete exactly one step from the HIGH_LEVEL_PLAN.md before stopping to allow the next agent to continue.

## Required Reading Before Starting

**MANDATORY**: Every agent MUST read these four documents before taking any action:

1. **PRODUCT_SPECIFICATION.md** - Understand the complete project requirements
2. **HIGH_LEVEL_PLAN.md** - Review the implementation plan and current progress
3. **Python Development Standards** - CONTRIBUTING/python-development-standards.md
4. **Python CLI Development Standards** - CONTRIBUTING/python-cli-development-standards.md

## Critical Constraints

### Code Quality Requirements
- **MAXIMUM 1000 lines of code total** across all modules
- **MAXIMUM 100 lines per module** (target <100, hard limit 100)
- **Python 3.10+ compatibility** required
- **Ultra-minimal dependencies**: Click, LiteLLM, PyYAML only
- **No rich formatting, colors, or progress indicators**
- **Agent-friendly naming**: Use searchable keywords in function/class names

### Implementation Rules
- **ONE STEP ONLY**: Complete exactly one step from HIGH_LEVEL_PLAN.md then STOP
- **UPDATE PROGRESS**: Update progress tracking in HIGH_LEVEL_PLAN.md after completing step
- **VALIDATE ADHERENCE**: Ensure implementation matches PRODUCT_SPECIFICATION.md requirements
- **FOLLOW STANDARDS**: Adhere to Python development standards documents
- **THINK FIRST**: Analyze the task and plan approach before writing any code

### Prohibited Actions
- **NO multiple steps** - complete one step only
- **NO deviating** from the specification
- **NO adding features** not in the specification
- **NO complex abstractions** - keep it simple
- **NO dependencies** beyond Click, LiteLLM, PyYAML

## Agent Workflow

### 1. Assessment Phase (REQUIRED)

Before writing any code, the agent MUST:

```
1. Read all four required documents completely
2. Identify current implementation status from HIGH_LEVEL_PLAN.md progress tracking
3. Determine which step to implement next
4. Verify the step aligns with PRODUCT_SPECIFICATION.md requirements
5. Plan the specific implementation approach
```

### 2. Implementation Phase

After assessment, proceed with:

```
1. Implement exactly one step from the plan
2. Follow Python development standards strictly
3. Keep modules under 100 lines each
4. Use agent-friendly naming conventions
5. Add comprehensive docstrings with keywords
```

### 3. Validation Phase (REQUIRED)

After implementation, the agent MUST:

```
1. Count total lines of code across all modules
2. Verify no prohibited features were added
3. Confirm adherence to PRODUCT_SPECIFICATION.md
4. Test the implemented functionality works
5. Update progress in HIGH_LEVEL_PLAN.md
```

### 4. Documentation Phase (REQUIRED)

Before stopping, the agent MUST:

```
1. Update progress tracking section in HIGH_LEVEL_PLAN.md
2. Mark completed step as [x] Done
3. Update "Current Step" and "Next Agent Task" 
4. Document any discoveries or implementation notes
5. Ensure next agent knows exactly what to do
```

## Step-by-Step Implementation Guide

### Current Project Status
Check `HIGH_LEVEL_PLAN.md` progress tracking section to determine current status.

### Implementation Order
Follow the exact order from HIGH_LEVEL_PLAN.md:

1. **Project Setup** - Initialize uv project, pyproject.toml
2. **Configuration** - YAML config loading, environment variables  
3. **CLI Foundation** - Click-based CLI structure
4. **Agent Core** - LiteLLM integration and conversation handling
5. **Bash Tool** - Single tool implementation
6. **Session Management** - File-based persistence
7. **CLI Integration** - Wire components together
8. **Entry Point** - Package configuration and installation
9. **Testing** - Unit and integration tests
10. **Validation** - Final checks and documentation

### Module-Specific Guidelines

#### cli.py (~80 lines max)
```python
"""Main CLI entry point for the minimal AI coding agent.

Keywords: CLI, command-line, interface, click, main, entry-point

Implements:
- Interactive mode (default)
- Single-shot mode (--prompt)
- File input mode (--file) 
- Session resume (--resume)
- Tool flags (--allow-tools/--no-tools)
- Confirmation flags (--confirm/--no-confirm)
"""
```

#### agent.py (~100 lines max)
```python
"""Core agent logic with LiteLLM integration.

Keywords: agent, litellm, conversation, chat, model, AI

Implements:
- LiteLLM client initialization
- Conversation loop management
- Message role handling (user, assistant, tool)
- Tool execution coordination
- Error handling and recovery
"""
```

#### session.py (~80 lines max) 
```python
"""File-based session management for conversation persistence.

Keywords: session, persistence, file-based, conversation, history

Implements:
- Session creation with timestamp IDs
- JSON-based conversation storage
- Session loading and resuming
- Simple file-based persistence
"""
```

#### config.py (~60 lines max)
```python
"""YAML configuration loading and management.

Keywords: configuration, YAML, config, settings, environment

Implements:
- YAML config file loading
- Environment variable overrides  
- Default configuration values
- Basic configuration validation
"""
```

#### bash_tool.py (~50 lines max)
```python
"""Single bash tool implementation for command execution.

Keywords: bash, tool, subprocess, command, execution, shell

Implements:
- Subprocess command execution
- Confirmation prompt handling
- Tool enable/disable logic
- Timeout and error handling
"""
```

## Error Handling Guidelines

### When Implementation Fails
If the agent cannot complete a step:

1. **Document the issue** in HIGH_LEVEL_PLAN.md
2. **Do not attempt multiple approaches** - report and stop
3. **Update progress tracking** to indicate the blocker
4. **Provide specific next steps** for resolution

### When Requirements Conflict
If specifications seem to conflict:

1. **PRODUCT_SPECIFICATION.md takes precedence** over implementation details
2. **Python standards documents** take precedence for code style
3. **1000-line limit** is absolute and cannot be exceeded
4. **Ask for clarification** rather than making assumptions

## Quality Assurance Checklist

Before marking a step complete, verify:

- [ ] **Line Count**: Total project lines < 1000, current module < 100
- [ ] **Dependencies**: Only Click, LiteLLM, PyYAML used
- [ ] **Functionality**: Implemented feature works as specified
- [ ] **Standards**: Follows Python development standards
- [ ] **Documentation**: Docstrings with keywords added
- [ ] **Progress**: HIGH_LEVEL_PLAN.md updated with progress
- [ ] **Specification**: Matches PRODUCT_SPECIFICATION.md requirements
- [ ] **Testing**: Basic functionality verified
- [ ] **Next Steps**: Clear instructions for next agent

## Success Criteria

### Step Completion
A step is complete when:
- All functionality for the step works correctly
- Code follows all quality requirements  
- Progress is documented in HIGH_LEVEL_PLAN.md
- Next agent has clear instructions

### Project Completion  
The project is complete when:
- All 10 steps in HIGH_LEVEL_PLAN.md are marked complete
- Total code is under 1000 lines
- All CLI modes work (interactive, single-shot, file, resume)
- Basic tests pass
- Installation via `uv tool install` works

## Common Pitfalls to Avoid

1. **Implementing multiple steps** - Only do one step per agent
2. **Adding extra features** - Stick strictly to the specification  
3. **Over-engineering** - Keep implementations simple
4. **Exceeding line limits** - Monitor code size constantly
5. **Skipping documentation** - Always update progress tracking
6. **Ignoring standards** - Follow Python development rules exactly
7. **Complex abstractions** - Prefer simple, direct implementations
8. **Missing validation** - Test functionality before marking complete

## Emergency Procedures

### If Line Count Approaches 1000
1. **Stop implementation immediately**
2. **Refactor existing code** to reduce lines
3. **Simplify implementations** where possible
4. **Document the constraint** in HIGH_LEVEL_PLAN.md

### If Specification Requires Clarification
1. **Document the ambiguity** clearly
2. **Do not make assumptions** or add features
3. **Update HIGH_LEVEL_PLAN.md** with the question
4. **Stop and wait** for clarification

### If Dependencies Don't Support Requirements
1. **Document the limitation** in HIGH_LEVEL_PLAN.md
2. **Propose minimal workarounds** within spec limits
3. **Do not add new dependencies** without approval
4. **Stop and request guidance** on resolution

---

## Final Reminder

**Each agent implementing this project MUST:**
1. Read all four required documents first
2. Complete exactly one step from HIGH_LEVEL_PLAN.md  
3. Update progress tracking before stopping
4. Ensure next agent has clear instructions
5. Never exceed 1000 total lines of code
6. Follow ultra-minimal design principles

**The goal is a working, ultra-lightweight AI coding agent under 1000 lines that does exactly what the specification requires - nothing more, nothing less.**