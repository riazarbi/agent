# High-Level Implementation Plan: Minimal AI Coding Agent

## Project Structure (Following Python Standards)

```
agent/
├── pyproject.toml              # Project metadata, dependencies
├── uv.lock                     # Lockfile
├── src/
│   └── python_agent/
│       ├── __init__.py         # (~10 lines)
│       ├── cli.py              # Main CLI entry point (~80 lines)
│       ├── agent.py            # Core agent logic (~100 lines)
│       ├── session.py          # Session management (~80 lines)
│       ├── config.py           # Configuration handling (~60 lines)
│       └── bash_tool.py        # Bash tool implementation (~50 lines)
└── tests/                      # Test files
    ├── __init__.py
    ├── test_cli.py
    ├── test_agent.py
    ├── test_session.py
    ├── test_config.py
    └── test_bash_tool.py
```

**Total Target: ~380 lines of code across all modules**

## Implementation Phases

### Phase 1: Project Setup and Configuration (Steps 1-3)
**Goal**: Establish project structure and configuration system

#### Step 1: Initialize Project Structure
- [ ] Create `uv` project with Python 3.10+ requirement
- [ ] Set up `pyproject.toml` with minimal dependencies
- [ ] Create src-layout directory structure
- [ ] Configure development tools (ruff, mypy, pytest)

**Files Created**: `pyproject.toml`, directory structure

#### Step 2: Configuration System
- [ ] Implement `config.py` with YAML config loading
- [ ] Support environment variable overrides for API keys
- [ ] Create default config file template
- [ ] Add basic validation

**Files Created**: `src/python_agent/config.py`, default config template

#### Step 3: CLI Foundation
- [ ] Implement basic CLI structure in `cli.py`
- [ ] Add Click commands and options
- [ ] Handle argument parsing for all modes
- [ ] Basic error handling and help text

**Files Created**: `src/python_agent/cli.py`

### Phase 2: Core Agent Implementation (Steps 4-6)
**Goal**: Implement LiteLLM integration and conversation handling

#### Step 4: Agent Core Logic
- [ ] Implement `agent.py` with LiteLLM integration
- [ ] Create conversation loop for interactive mode
- [ ] Handle single-shot and file input modes
- [ ] Implement basic error handling

**Files Created**: `src/python_agent/agent.py`

#### Step 5: Bash Tool Implementation
- [ ] Implement `bash_tool.py` with subprocess execution
- [ ] Add confirmation prompts (when enabled)
- [ ] Handle tool enabling/disabling
- [ ] Basic timeout and error handling

**Files Created**: `src/python_agent/bash_tool.py`

#### Step 6: Session Management
- [ ] Implement `session.py` for conversation persistence
- [ ] File-based session storage with simple JSON format
- [ ] Session creation, saving, and loading
- [ ] Resume session by ID functionality

**Files Created**: `src/python_agent/session.py`

### Phase 3: Integration and CLI Entry (Steps 7-8)
**Goal**: Complete the application and ensure all modes work

#### Step 7: CLI Integration
- [ ] Wire all components together in CLI
- [ ] Implement interactive mode loop
- [ ] Handle all command-line flags properly
- [ ] Add proper exit codes and cleanup

**Files Updated**: `src/python_agent/cli.py`

#### Step 8: Entry Point and Packaging
- [ ] Create `__init__.py` and package entry point
- [ ] Configure CLI script in pyproject.toml
- [ ] Test installation with `uv tool install`
- [ ] Validate all CLI modes work correctly

**Files Updated**: `src/python_agent/__init__.py`, `pyproject.toml`

### Phase 4: Testing and Documentation (Steps 9-10)
**Goal**: Ensure reliability and provide basic documentation

#### Step 9: Testing Implementation
- [ ] Unit tests for each module
- [ ] Integration tests for CLI functionality
- [ ] API tests with real LiteLLM calls (with skip logic)
- [ ] Achieve minimum 80% test coverage

**Files Created**: All test files

#### Step 10: Final Validation and Documentation
- [ ] End-to-end testing of all features
- [ ] Line count validation (must be under 400 lines)
- [ ] Performance testing (startup < 200ms)
- [ ] Create basic README with usage examples

**Files Created**: `README.md`, final validation

## Detailed Module Specifications

### cli.py (~80 lines)
```python
# CLI entry point with Click commands
# Commands: default (interactive), --prompt, --file, --resume
# Flags: --allow-tools/--no-tools, --confirm/--no-confirm
# Error handling and help text
```

### agent.py (~100 lines)
```python  
# Core agent logic with LiteLLM integration
# Conversation loop handling
# Message formatting and role management
# Tool execution coordination
```

### session.py (~80 lines)
```python
# File-based session persistence
# JSON format for conversation history
# Session creation, save, load, resume
# Simple ID-based session management
```

### config.py (~60 lines)
```python
# YAML configuration loading
# Environment variable override handling
# Default configuration values
# Basic validation
```

### bash_tool.py (~50 lines)
```python
# Subprocess execution wrapper
# Confirmation prompt handling
# Tool enable/disable logic
# Basic error handling and timeouts
```

### __init__.py (~10 lines)
```python
# Package initialization
# Version info
# Main entry point
```

## Technical Dependencies

### Required Dependencies
```toml
[project]
dependencies = [
    "click>=8.0",
    "litellm>=1.0", 
    "pyyaml>=6.0",
]

[project.optional-dependencies]
dev = [
    "pytest>=7.0",
    "pytest-cov>=4.0", 
    "ruff>=0.1",
    "mypy>=1.0",
]
```

### Key Design Decisions

#### LiteLLM Integration
- Direct API key passthrough to LiteLLM
- No custom model abstraction layer
- Use LiteLLM's built-in provider support

#### Session Storage
- Simple JSON files in ~/.agent/sessions/
- Session ID format: YYYY-MM-DD-HH-MM-SS
- No metadata beyond conversation history

#### Tool Implementation  
- Single bash tool only
- Use subprocess.run() with shell=True
- Basic timeout handling (30 seconds default)
- Optional confirmation prompts

#### Configuration
- YAML file at ~/.agent/config.yaml
- Environment variables: API_KEY, MODEL, BASE_URL
- Minimal configuration options only

## Progress Tracking

**Implementation Status**: Not Started
- [ ] **Step 1**: Project setup and structure
- [ ] **Step 2**: Configuration system
- [ ] **Step 3**: CLI foundation  
- [ ] **Step 4**: Agent core logic
- [ ] **Step 5**: Bash tool implementation
- [ ] **Step 6**: Session management
- [ ] **Step 7**: CLI integration
- [ ] **Step 8**: Entry point and packaging
- [ ] **Step 9**: Testing implementation
- [ ] **Step 10**: Final validation and documentation

**Current Step**: Step 1 (Project setup and structure)
**Next Agent Task**: Initialize the uv project structure and create pyproject.toml with minimal dependencies

## Validation Criteria

### Code Quality
- [ ] Total lines of code < 400
- [ ] Each module < 100 lines  
- [ ] All functions have docstrings with keywords
- [ ] Type hints on public functions
- [ ] Follows Python development standards

### Functionality
- [ ] Interactive mode works
- [ ] Single-shot mode works
- [ ] File input mode works
- [ ] Session resume works
- [ ] Bash tool executes commands
- [ ] Configuration loads from YAML

### Performance
- [ ] Startup time < 200ms
- [ ] Memory usage < 50MB
- [ ] No unnecessary dependencies loaded

### Testing
- [ ] Unit test coverage > 80%
- [ ] Integration tests pass
- [ ] CLI tests pass
- [ ] Real API tests (with skip logic)