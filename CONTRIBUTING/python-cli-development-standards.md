# Python CLI Development Standards

An overlay to the [Python Development Standards](python-development-standards.md) providing additional requirements for command-line interface tools.

## Prerequisites

This document extends the base [Python Development Standards](python-development-standards.md). All base requirements apply, plus the CLI-specific standards below.

## CLI Project Structure Requirements

### Additional Directory Structure

CLI projects **MUST** extend the base structure with:

```
project-name/                   # Base structure applies
├── src/
│   └── project_name/
│       ├── cli.py              # REQUIRED: Main CLI entry point
│       ├── commands/           # REQUIRED for multi-command CLIs
│       │   ├── __init__.py
│       │   └── *.py            # Individual command modules
│       └── config/             # CLI configuration handling
└── tests/
    └── cli/                    # CLI-specific tests
        ├── test_cli_integration.py
        └── test_commands.py
```

### Entry Point Configuration

**MUST** define CLI entry points in `pyproject.toml`:

```toml
[project.scripts]
project-name = "project_name.cli:main"
```

## CLI Naming Standards

### Required Naming Conventions
- **CLI command name**: `kebab-case` (e.g., `data-processor`, `user-manager`)
- **Package directory**: `snake_case` (e.g., `data_processor`, `user_manager`)  
- **Entry point module**: Always `cli.py`
- **Command modules**: `snake_case` with descriptive names

### CLI-to-Package Mapping
```
CLI Command      → Package Directory → Entry Module
data-processor   → data_processor    → cli.py
user-manager     → user_manager      → cli.py
```

## CLI Architecture Requirements

### Main CLI Module (`cli.py`)

**MUST** implement:
- Click-based command group with global options
- Context object for sharing configuration
- Proper exception handling with exit codes
- Main entry point function returning int

**MUST** provide global options:
- `--config, -c`: Configuration file path
- `--verbose, -v`: Verbose output flag  
- `--quiet, -q`: Quiet mode flag
- `--help`: Comprehensive help text

### Command Modules

**MUST** implement for each command:
- Click command decorator with proper arguments/options
- Comprehensive docstring with examples
- Type hints for all parameters
- Context-aware configuration access
- Proper error handling with ClickException

## CLI User Experience Standards

### Help Text Requirements

**MUST** provide at all levels:
- **Tool level**: Purpose, common workflows, basic examples
- **Command level**: What it does, input/output, usage examples
- **Option level**: Clear descriptions with default values

### Error Handling Requirements

**MUST** implement:
- User-friendly error messages (not stack traces)
- Actionable suggestions when possible
- Appropriate exit codes (0=success, 1=error, 2=usage error, etc.)
- Graceful Ctrl+C handling

### Output Standards

**MUST** support:
- Multiple output formats (JSON, YAML, table)
- Quiet mode (errors only)
- Verbose mode (progress/debug info)
- Consistent formatting across commands

### Progress Indication

**MUST** provide for operations >2 seconds:
- Progress bars or status indicators
- Clear start/completion messages
- Cancellation support (Ctrl+C)

## CLI Testing Requirements

### Required Test Categories

**MUST** implement beyond base testing:

#### CLI Integration Tests
- Help text display and accuracy
- Command argument validation
- Configuration file handling
- Exit code correctness

#### CLI End-to-End Tests  
- Complete user workflows via subprocess
- Input/output file handling
- Error scenarios and messages
- Cross-platform compatibility

### CLI Test Helpers

**MUST** provide:
- `run_cli_command()` function for subprocess testing
- Structured result objects (exit_code, stdout, stderr)
- Helper assertions for CLI success/failure

## CLI Quality Gates

### Additional Requirements Beyond Base Standards

**MUST** achieve:
- **100% CLI command coverage**: All commands and options tested
- **Help text completeness**: Every command has comprehensive help  
- **Error message quality**: All errors are actionable
- **Performance standards**: CLI startup < 500ms
- **Cross-platform support**: Linux, macOS, Windows

### CLI-Specific Enforcement

**MUST** validate:
- All argument combinations work correctly
- All configuration options are functional
- All output formats produce valid results  
- All help examples are accurate and tested
- Exit codes match documented behavior

## CLI Configuration Standards

### Configuration File Support

**MUST** support (if configuration is needed):
- TOML format configuration files
- Standard config file locations (`~/.config/project-name/`, etc.)
- Environment variable override (`PROJECT_NAME_CONFIG`)
- Command-line config file specification

### Environment Variables

**SHOULD** support standard environment variables:
- `PROJECT_NAME_CONFIG`: Config file path
- `PROJECT_NAME_VERBOSE`: Default verbose mode
- `NO_COLOR`: Disable color output

## CLI Distribution Requirements  

### Installation Methods

**MUST** support:
- `uv tool install project-name` (primary method)
- `pipx install project-name` (alternative)
- Development installation via `uv sync`

### Docker Support

**SHOULD** provide if containerization is needed:
- Dockerfile optimized for CLI execution
- Non-root user configuration
- Proper entrypoint setup

## Performance Standards

### Response Time Requirements
- **Command startup**: < 500ms for simple commands
- **Help display**: < 200ms 
- **Command execution**: Appropriate to task complexity
- **Progress indication**: For operations > 2 seconds

### Resource Usage
- **Memory**: < 50MB for typical operations
- **CPU**: Efficient processing, avoid busy loops
- **I/O**: Stream large outputs, don't buffer unnecessarily

## CLI-Specific Violations

Breaking CLI-specific **MUST** rules results in:
- CLI functionality failures
- Poor user experience
- Failed CLI integration tests
- Distribution/installation issues

**CLI Quality Question**: Does this command-line tool provide a professional, intuitive user experience that follows CLI best practices?