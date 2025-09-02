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
- [ ] Line count validation (must be under 1000 lines)
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

**Implementation Status**: In Progress (7/10 steps completed)
- [x] **Step 1**: Project setup and structure ✓ DONE
- [x] **Step 2**: Configuration system ✓ DONE
- [x] **Step 3**: CLI foundation ✓ DONE
- [x] **Step 4**: Agent core logic ✓ DONE
- [x] **Step 5**: Bash tool implementation ✓ DONE
- [x] **Step 6**: Session management ✓ DONE
- [x] **Step 7**: CLI integration ✓ DONE
- [ ] **Step 8**: Entry point and packaging
- [ ] **Step 9**: Testing implementation
- [ ] **Step 10**: Final validation and documentation

**Current Step**: Step 8 (Entry point and packaging)
**Next Agent Task**: Create __init__.py and package entry point, configure CLI script in pyproject.toml, test installation with uv tool install, and validate all CLI modes work correctly

**Step 1 Completed**: 
- ✅ Created uv project with Python 3.10+ requirement
- ✅ Set up pyproject.toml with minimal dependencies (Click, LiteLLM, PyYAML)
- ✅ Created src-layout directory structure
- ✅ Configured development tools (ruff, mypy, pytest, black)
- ✅ Verified all tools work correctly
- ✅ Current code: 10 lines total (well within 1000-line limit)

**Step 2 Completed**:
- ✅ Implemented config.py with YAML config loading (156 lines)
- ✅ Added environment variable overrides for all configuration options
- ✅ Implemented basic validation for timeout, max_tokens, and temperature
- ✅ Provided default configuration with all required settings
- ✅ Tested YAML loading, environment overrides, and validation
- ✅ Current code: 166 lines total (well within 1000-line limit)

**Step 3 Completed**:
- ✅ Implemented cli.py with Click-based CLI structure (120 lines)
- ✅ Added all required command-line options: --prompt, --file, --resume
- ✅ Implemented global options: --config, --allow-tools/--no-tools, --confirm/--no-confirm
- ✅ Added verbose and quiet modes with appropriate flags
- ✅ Included comprehensive help text and usage examples
- ✅ Implemented proper error handling with appropriate exit codes
- ✅ Added mode exclusivity validation (prevents conflicting options)
- ✅ Integrated with config.py for configuration loading and CLI overrides
- ✅ Tested CLI help display and basic execution
- ✅ Current code: 286 lines total (well within 1000-line limit)

**Step 4 Completed**:
- ✅ Implemented agent.py with LiteLLM integration (92 lines)
- ✅ Created Agent class with configuration-based initialization
- ✅ Implemented conversation history management with add_message method
- ✅ Added chat_completion method with proper LiteLLM API calls
- ✅ Implemented process_single_prompt for single-shot mode
- ✅ Created interactive_loop for continuous conversation mode
- ✅ Added proper error handling with AgentError and ModelError classes
- ✅ Integrated agent with CLI for all modes (interactive, single-shot, file input)
- ✅ Tested CLI help and basic functionality with proper error messages
- ✅ Current code: 407 lines total (well within 1000-line limit)

**Step 5 Completed**:
- ✅ Implemented bash_tool.py with subprocess execution (123 lines)
- ✅ Created BashTool class with configuration options for confirmation and timeout
- ✅ Added tool enable/disable functionality with set_enabled method
- ✅ Implemented confirmation prompt system with _confirm_execution method
- ✅ Added comprehensive execute_command method with timeout handling
- ✅ Proper exception handling for disabled tool, timeouts, and general errors
- ✅ Returns structured dictionary results with success, output, error, and exit_code
- ✅ Follows Python development standards with type hints and comprehensive docstrings
- ✅ Tested basic functionality including command execution and timeout handling
- ✅ Current code: 519 lines total (well within 1000-line limit)

**Step 6 Completed**:
- ✅ Implemented session.py with file-based session persistence (212 lines)
- ✅ Created Session class for individual conversation data containers
- ✅ Added SessionManager class for session operations (create, save, load, list)
- ✅ Implemented JSON-based session storage in ~/.agent/sessions/ directory
- ✅ Added timestamp-based session ID generation (YYYY-MM-DD-HH-MM-SS format)
- ✅ Implemented comprehensive error handling with SessionError class
- ✅ Added session listing functionality sorted by creation date
- ✅ Included proper type hints and comprehensive docstrings with keywords
- ✅ Follows Python development standards with agent-friendly naming
- ✅ Tested basic functionality including session creation and file operations
- ✅ Current code: 731 lines total (well within 1000-line limit)

**Step 7 Completed**:
- ✅ Integrated session management into CLI for session resume functionality
- ✅ Added SessionManager import and session loading with proper error handling
- ✅ Implemented session resume mode with --resume flag integration
- ✅ Added bash tool integration to Agent class with configuration support
- ✅ Enhanced Agent constructor with BashTool initialization and tool enable/disable logic
- ✅ Added session management methods: start_new_session, resume_from_session, save_current_session
- ✅ Enhanced interactive_loop with automatic session creation and persistence
- ✅ Improved error handling and exit codes with proper cleanup
- ✅ Added comprehensive exception handling including KeyboardInterrupt (Ctrl+C)
- ✅ Implemented proper file input error handling and verbose logging
- ✅ Tested all CLI modes successfully: single-shot, file input, session resume, interactive
- ✅ Verified session persistence and resume functionality works correctly
- ✅ Confirmed total code remains at 809 lines (well within 1000-line limit)

## Validation Criteria

### Code Quality
- [ ] Total lines of code < 1000
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
- [x] Configuration loads from YAML

### Performance
- [ ] Startup time < 200ms
- [ ] Memory usage < 50MB
- [ ] No unnecessary dependencies loaded

### Testing
- [ ] Unit test coverage > 80%
- [ ] Integration tests pass
- [ ] CLI tests pass
- [ ] Real API tests (with skip logic)