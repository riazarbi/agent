# Python CLI Development Standards

An overlay to the [Python Development Standards](python-development-standards.md) providing specific guidelines for command-line interface tools.

## Prerequisites

This document extends the base [Python Development Standards](python-development-standards.md). All base requirements apply, plus the additional CLI-specific standards below.

## CLI Project Structure

### Required CLI-Specific Additions

CLI projects **MUST** extend the base project structure with these CLI-specific components:

```
project-name/                   # Base structure from python-development-standards.md
├── src/
│   └── project_name/
│       ├── cli.py              # REQUIRED: Main CLI entry point
│       ├── commands/           # REQUIRED for multi-command CLIs
│       │   ├── __init__.py
│       │   ├── command_one.py
│       │   └── command_two.py
│       ├── config/             # CLI configuration handling
│       │   ├── __init__.py
│       │   └── settings.py
│       └── ui/                 # User interface components
│           ├── __init__.py
│           ├── formatting.py   # Output formatting
│           ├── progress.py     # Progress indicators
│           └── prompts.py      # User prompts and interactions
└── tests/
    ├── cli/                    # CLI-specific tests
    │   ├── __init__.py
    │   ├── test_cli_integration.py
    │   └── test_commands.py
    └── fixtures/
        ├── cli_inputs/         # Sample CLI inputs
        └── cli_outputs/        # Expected CLI outputs
```

### CLI Entry Point Configuration

**MUST** define CLI entry points in `pyproject.toml`:

```toml
[project.scripts]
project-name = "project_name.cli:main"
# For multi-command tools, use single entry point with subcommands
```

## CLI Naming Conventions

### Command Naming Standards
- **CLI command name**: `kebab-case` (e.g., `data-processor`, `user-manager`)
- **Package directory**: `snake_case` (e.g., `data_processor`, `user_manager`)  
- **Entry point module**: Always `cli.py`
- **Command modules**: `snake_case` with descriptive names

### CLI-to-Package Mapping
```
CLI Command      → Package Directory → Entry Module
data-processor   → data_processor    → cli.py
user-manager     → user_manager      → cli.py
api-client       → api_client        → cli.py
```

## CLI Architecture Patterns

### Main CLI Module Structure
```python
# src/project_name/cli.py
"""Main CLI entry point for project-name.

This module provides the primary command-line interface for the application,
handling argument parsing, command routing, and global CLI concerns.

Keywords: cli, command line, entry point, main, interface
"""

import sys
from pathlib import Path
from typing import Optional

import click

from project_name.commands import command_one, command_two
from project_name.config.settings import load_cli_config
from project_name.ui.formatting import setup_output_formatting
from project_name.ui.progress import ProgressTracker


@click.group()
@click.option("--config", "-c", type=click.Path(exists=True), 
              help="Path to configuration file")
@click.option("--verbose", "-v", is_flag=True, help="Enable verbose output")
@click.option("--quiet", "-q", is_flag=True, help="Suppress non-error output")
@click.pass_context
def cli(ctx: click.Context, config: Optional[str], verbose: bool, quiet: bool) -> None:
    """Project Name - Brief description of what this tool does.
    
    Longer description explaining the purpose, main use cases, and basic usage
    patterns. Include examples of common workflows.
    
    Examples:
        project-name command-one --input data.txt
        project-name command-two --output results.json --format pretty
    """
    # Ensure context object exists
    ctx.ensure_object(dict)
    
    # Load and store configuration
    ctx.obj["config"] = load_cli_config(config)
    ctx.obj["verbose"] = verbose
    ctx.obj["quiet"] = quiet
    
    # Setup output formatting based on flags
    setup_output_formatting(verbose=verbose, quiet=quiet)


# Register commands
cli.add_command(command_one.command_one)
cli.add_command(command_two.command_two)


def main() -> int:
    """Main entry point for CLI application.
    
    Returns:
        Exit code (0 for success, non-zero for error)
    """
    try:
        cli()
        return 0
    except KeyboardInterrupt:
        click.echo("\nOperation cancelled by user", err=True)
        return 130
    except Exception as e:
        click.echo(f"Error: {e}", err=True)
        return 1


if __name__ == "__main__":
    sys.exit(main())
```

### Command Module Structure
```python
# src/project_name/commands/command_one.py
"""Command One implementation.

Handles [specific functionality] including [key features].

Keywords: command, [domain-specific keywords]
"""

from pathlib import Path
from typing import Optional

import click

from project_name.services.processor import DataProcessor
from project_name.ui.progress import show_progress
from project_name.ui.formatting import format_output


@click.command()
@click.argument("input_file", type=click.Path(exists=True, path_type=Path))
@click.option("--output", "-o", type=click.Path(path_type=Path),
              help="Output file path (default: stdout)")
@click.option("--format", type=click.Choice(["json", "yaml", "csv"]), 
              default="json", help="Output format")
@click.pass_context
def command_one(
    ctx: click.Context, 
    input_file: Path, 
    output: Optional[Path],
    format: str
) -> None:
    """Process input file and generate formatted output.
    
    This command reads data from INPUT_FILE, processes it according to the
    configured rules, and outputs the results in the specified format.
    
    Examples:
        project-name command-one data.txt --output results.json
        project-name command-one data.txt --format yaml
    """
    config = ctx.obj["config"]
    verbose = ctx.obj["verbose"]
    
    try:
        # Initialize processor with configuration
        processor = DataProcessor(config.processing_settings)
        
        # Show progress if verbose mode
        with show_progress("Processing data...", verbose=verbose) as progress:
            result = processor.process_file(input_file, progress_callback=progress.update)
        
        # Format and output results
        formatted_output = format_output(result, format=format)
        
        if output:
            output.write_text(formatted_output)
            click.echo(f"Results written to {output}")
        else:
            click.echo(formatted_output)
            
    except Exception as e:
        raise click.ClickException(f"Command failed: {e}")
```

## CLI User Experience Standards

### Help Text Requirements
**MUST** provide comprehensive help at all levels:

```python
# Command group help
@click.group()
def cli():
    """Tool Name - One-line description.
    
    Longer description explaining:
    - What the tool does
    - Primary use cases  
    - Basic workflow overview
    
    Examples:
        tool-name command --option value
        tool-name --help
        tool-name command --help
    """

# Command help  
@click.command()
def command():
    """Action verb + what it processes/produces.
    
    Detailed description including:
    - What the command does
    - Expected input/output
    - Side effects or requirements
    - Common use cases
    
    Examples:
        tool-name command input.txt --output result.json
        tool-name command --format yaml < input.txt
    """
```

### Error Handling & User Feedback

#### CLI Error Handling Template
```python
# src/project_name/exceptions.py
"""CLI-specific exceptions and error handling.

Keywords: exceptions, cli errors, user feedback
"""

import click


class CLIError(click.ClickException):
    """Base exception for CLI-specific errors."""
    
    def __init__(self, message: str, exit_code: int = 1) -> None:
        super().__init__(message)
        self.exit_code = exit_code


class ConfigurationError(CLIError):
    """Raised when CLI configuration is invalid."""
    
    def __init__(self, message: str, config_path: str = "") -> None:
        if config_path:
            message = f"Configuration error in {config_path}: {message}"
        super().__init__(message)


class InputValidationError(CLIError):
    """Raised when user input is invalid."""
    
    def __init__(self, message: str, suggestion: str = "") -> None:
        full_message = message
        if suggestion:
            full_message += f"\n\nSuggestion: {suggestion}"
        super().__init__(full_message)


def handle_cli_error(e: Exception) -> int:
    """Handle CLI errors with user-friendly messages.
    
    Returns:
        Appropriate exit code
    """
    if isinstance(e, CLIError):
        click.echo(f"Error: {e.message}", err=True)
        return e.exit_code
    elif isinstance(e, FileNotFoundError):
        click.echo(f"Error: File not found: {e.filename}", err=True)
        return 2
    elif isinstance(e, PermissionError):
        click.echo(f"Error: Permission denied: {e.filename}", err=True)
        return 3
    else:
        click.echo(f"Unexpected error: {e}", err=True)
        return 1
```

### Progress Indicators & User Feedback

#### Progress Tracking Implementation
```python
# src/project_name/ui/progress.py
"""Progress indicators and user feedback for CLI operations.

Keywords: progress, ui, user feedback, indicators
"""

import sys
from contextlib import contextmanager
from typing import Iterator, Optional, Callable

import click


class ProgressTracker:
    """Handles progress indication for long-running operations."""
    
    def __init__(self, description: str, total: Optional[int] = None, 
                 verbose: bool = True) -> None:
        self.description = description
        self.total = total
        self.verbose = verbose
        self.current = 0
        self._progress_bar = None
    
    def __enter__(self):
        if self.verbose and self.total:
            self._progress_bar = click.progressbar(
                length=self.total,
                label=self.description,
                show_eta=True,
                show_percent=True
            )
            self._progress_bar.__enter__()
        elif self.verbose:
            click.echo(f"{self.description}...", err=True)
        return self
    
    def __exit__(self, exc_type, exc_val, exc_tb):
        if self._progress_bar:
            self._progress_bar.__exit__(exc_type, exc_val, exc_tb)
        elif self.verbose:
            if exc_type is None:
                click.echo("Done", err=True)
            else:
                click.echo("Failed", err=True)
    
    def update(self, increment: int = 1) -> None:
        """Update progress by increment amount."""
        self.current += increment
        if self._progress_bar:
            self._progress_bar.update(increment)
    
    def set_description(self, description: str) -> None:
        """Update the progress description."""
        self.description = description
        if self.verbose and not self._progress_bar:
            click.echo(f"{description}...", err=True)


@contextmanager
def show_progress(description: str, total: Optional[int] = None, 
                 verbose: bool = True) -> Iterator[ProgressTracker]:
    """Context manager for showing progress during operations."""
    tracker = ProgressTracker(description, total, verbose)
    with tracker:
        yield tracker


def confirm_action(message: str, default: bool = False) -> bool:
    """Get user confirmation for potentially destructive actions."""
    return click.confirm(message, default=default)


def prompt_for_input(message: str, default: str = "", 
                    hide_input: bool = False) -> str:
    """Prompt user for input with validation."""
    return click.prompt(message, default=default, hide_input=hide_input)
```

### Output Formatting Standards

#### Output Formatting Implementation
```python
# src/project_name/ui/formatting.py  
"""Output formatting utilities for CLI.

Keywords: formatting, output, display, json, yaml, table
"""

import json
from typing import Any, Dict, List
from pathlib import Path

import click


class OutputFormatter:
    """Handles different output formats for CLI results."""
    
    def __init__(self, format: str = "json", pretty: bool = True) -> None:
        self.format = format.lower()
        self.pretty = pretty
    
    def format_data(self, data: Any) -> str:
        """Format data according to specified format."""
        if self.format == "json":
            return self._format_json(data)
        elif self.format == "yaml":
            return self._format_yaml(data)
        elif self.format == "table":
            return self._format_table(data)
        elif self.format == "csv":
            return self._format_csv(data)
        else:
            raise ValueError(f"Unsupported format: {self.format}")
    
    def _format_json(self, data: Any) -> str:
        """Format data as JSON."""
        if self.pretty:
            return json.dumps(data, indent=2, ensure_ascii=False)
        return json.dumps(data, separators=(',', ':'))
    
    def _format_yaml(self, data: Any) -> str:
        """Format data as YAML."""
        try:
            import yaml
            return yaml.dump(data, default_flow_style=False, allow_unicode=True)
        except ImportError:
            raise click.ClickException("YAML output requires 'pyyaml' package")
    
    def _format_table(self, data: List[Dict[str, Any]]) -> str:
        """Format data as table (requires tabulate)."""
        try:
            from tabulate import tabulate
            if not data:
                return "No data to display"
            return tabulate(data, headers="keys", tablefmt="grid")
        except ImportError:
            raise click.ClickException("Table output requires 'tabulate' package")
    
    def _format_csv(self, data: List[Dict[str, Any]]) -> str:
        """Format data as CSV."""
        import csv
        import io
        
        if not data:
            return ""
        
        output = io.StringIO()
        writer = csv.DictWriter(output, fieldnames=data[0].keys())
        writer.writeheader()
        writer.writerows(data)
        return output.getvalue()


def setup_output_formatting(verbose: bool = False, quiet: bool = False) -> None:
    """Configure global output formatting based on CLI flags."""
    # Store formatting preferences globally for the CLI session
    click.get_current_context().meta["verbose"] = verbose
    click.get_current_context().meta["quiet"] = quiet


def format_output(data: Any, format: str = "json", pretty: bool = True) -> str:
    """Convenience function for formatting output data."""
    formatter = OutputFormatter(format, pretty)
    return formatter.format_data(data)


def echo_result(data: Any, format: str = "json", output_file: Optional[Path] = None) -> None:
    """Output formatted data to file or stdout."""
    formatted = format_output(data, format)
    
    if output_file:
        output_file.write_text(formatted)
        click.echo(f"Output written to {output_file}", err=True)
    else:
        click.echo(formatted)
```

## CLI Configuration Management

### Configuration File Handling
```python
# src/project_name/config/settings.py
"""CLI configuration management.

Keywords: configuration, settings, config files, cli config
"""

import os
from pathlib import Path
from typing import Optional, Dict, Any
from dataclasses import dataclass, field

import click


@dataclass
class CLIConfig:
    """CLI configuration settings."""
    
    # Global settings
    verbose: bool = False
    output_format: str = "json"
    color_output: bool = True
    
    # Application-specific settings
    processing_settings: Dict[str, Any] = field(default_factory=dict)
    api_settings: Dict[str, Any] = field(default_factory=dict)
    
    @classmethod
    def from_dict(cls, config_data: Dict[str, Any]) -> "CLIConfig":
        """Create CLIConfig from dictionary."""
        return cls(
            verbose=config_data.get("verbose", False),
            output_format=config_data.get("output_format", "json"),
            color_output=config_data.get("color_output", True),
            processing_settings=config_data.get("processing", {}),
            api_settings=config_data.get("api", {})
        )


def get_config_paths() -> List[Path]:
    """Get list of possible configuration file locations in priority order."""
    paths = []
    
    # 1. Environment variable
    if env_config := os.getenv("PROJECT_NAME_CONFIG"):
        paths.append(Path(env_config))
    
    # 2. Current directory
    paths.append(Path.cwd() / "project-name.toml")
    paths.append(Path.cwd() / ".project-name.toml")
    
    # 3. User home directory
    home = Path.home()
    paths.append(home / ".config" / "project-name" / "config.toml")
    paths.append(home / ".project-name.toml")
    
    # 4. System-wide  
    paths.append(Path("/etc/project-name/config.toml"))
    
    return paths


def load_cli_config(config_path: Optional[str] = None) -> CLIConfig:
    """Load CLI configuration from file or defaults."""
    config_data = {}
    
    if config_path:
        # Use specified config file
        config_file = Path(config_path)
        if not config_file.exists():
            raise click.ClickException(f"Configuration file not found: {config_path}")
        config_data = load_config_file(config_file)
    else:
        # Search for config file in standard locations
        for path in get_config_paths():
            if path.exists():
                config_data = load_config_file(path)
                break
    
    return CLIConfig.from_dict(config_data)


def load_config_file(config_path: Path) -> Dict[str, Any]:
    """Load configuration from TOML file."""
    try:
        import tomli
        with open(config_path, "rb") as f:
            return tomli.load(f)
    except ImportError:
        raise click.ClickException("Configuration files require 'tomli' package")
    except Exception as e:
        raise click.ClickException(f"Error loading config file {config_path}: {e}")
```

## CLI Testing Standards

### CLI-Specific Testing Patterns

#### CLI Integration Tests
```python
# tests/cli/test_cli_integration.py
"""Integration tests for CLI functionality.

Keywords: cli tests, integration, command line testing
"""

import subprocess
import json
from pathlib import Path
from typing import List, Dict, Any

import pytest

from tests.helpers.cli_helpers import run_cli_command, CLIResult


class TestCLIIntegration:
    """Integration tests for CLI commands."""
    
    def test_cli_help_displays_correctly(self) -> None:
        """Test that CLI help is displayed correctly."""
        result = run_cli_command(["--help"])
        
        assert result.exit_code == 0
        assert "Project Name" in result.stdout
        assert "Usage:" in result.stdout
        assert "Commands:" in result.stdout
    
    def test_command_help_displays_correctly(self) -> None:
        """Test that command-specific help is displayed."""
        result = run_cli_command(["command-one", "--help"])
        
        assert result.exit_code == 0
        assert "command-one" in result.stdout
        assert "Examples:" in result.stdout
    
    def test_invalid_command_shows_error(self) -> None:
        """Test that invalid commands show helpful error messages."""
        result = run_cli_command(["invalid-command"])
        
        assert result.exit_code != 0
        assert "No such command" in result.stderr
        assert "Did you mean" in result.stderr or "--help" in result.stderr


@pytest.mark.e2e
class TestCLIEndToEnd:
    """End-to-end tests for complete CLI workflows."""
    
    def test_complete_processing_workflow(self, tmp_path: Path) -> None:
        """Test complete processing workflow from CLI perspective."""
        # Setup test data
        input_file = tmp_path / "input.txt"
        input_file.write_text("test data for processing")
        output_file = tmp_path / "output.json"
        
        # Run CLI command
        result = run_cli_command([
            "command-one",
            str(input_file),
            "--output", str(output_file),
            "--format", "json"
        ])
        
        # Verify command succeeded
        assert result.exit_code == 0
        assert "Results written to" in result.stdout
        
        # Verify output file exists and has correct content
        assert output_file.exists()
        output_data = json.loads(output_file.read_text())
        assert "processed" in output_data
        assert output_data["source"] == str(input_file)
    
    def test_configuration_file_usage(self, tmp_path: Path) -> None:
        """Test CLI with configuration file."""
        # Create config file
        config_file = tmp_path / "test-config.toml"
        config_file.write_text("""
        verbose = true
        output_format = "yaml"
        
        [processing]
        option1 = "value1"
        option2 = 42
        """)
        
        # Run command with config
        result = run_cli_command([
            "--config", str(config_file),
            "command-one",
            "test-input"
        ])
        
        # Should succeed with config
        assert result.exit_code == 0
    
    def test_error_handling_and_messages(self, tmp_path: Path) -> None:
        """Test error handling provides helpful messages."""
        nonexistent_file = tmp_path / "does-not-exist.txt"
        
        result = run_cli_command([
            "command-one",
            str(nonexistent_file)
        ])
        
        assert result.exit_code != 0
        assert "Error:" in result.stderr
        assert str(nonexistent_file) in result.stderr
```

#### CLI Test Helpers
```python
# tests/helpers/cli_helpers.py
"""Helper functions for CLI testing.

Keywords: cli testing, test helpers, command line tests
"""

import subprocess
from dataclasses import dataclass
from pathlib import Path
from typing import List, Optional


@dataclass
class CLIResult:
    """Result of CLI command execution."""
    
    exit_code: int
    stdout: str
    stderr: str
    execution_time: float


def run_cli_command(
    args: List[str],
    input_data: Optional[str] = None,
    timeout: float = 30.0,
    env: Optional[dict] = None
) -> CLIResult:
    """Run CLI command and return structured result.
    
    Args:
        args: Command line arguments (without 'uv run project-name')
        input_data: Optional stdin data
        timeout: Command timeout in seconds
        env: Optional environment variables
        
    Returns:
        CLIResult with execution details
    """
    import time
    
    # Build full command
    full_cmd = ["uv", "run", "project-name"] + args
    
    start_time = time.time()
    
    try:
        result = subprocess.run(
            full_cmd,
            input=input_data,
            capture_output=True,
            text=True,
            timeout=timeout,
            env=env
        )
        
        execution_time = time.time() - start_time
        
        return CLIResult(
            exit_code=result.returncode,
            stdout=result.stdout,
            stderr=result.stderr,
            execution_time=execution_time
        )
        
    except subprocess.TimeoutExpired:
        execution_time = time.time() - start_time
        return CLIResult(
            exit_code=-1,
            stdout="",
            stderr=f"Command timed out after {timeout} seconds",
            execution_time=execution_time
        )


def assert_cli_success(result: CLIResult, expected_output: Optional[str] = None) -> None:
    """Assert CLI command succeeded with optional output check."""
    assert result.exit_code == 0, f"CLI command failed: {result.stderr}"
    
    if expected_output:
        assert expected_output in result.stdout, \
            f"Expected '{expected_output}' not found in output: {result.stdout}"


def assert_cli_error(result: CLIResult, expected_error: Optional[str] = None) -> None:
    """Assert CLI command failed with optional error check."""
    assert result.exit_code != 0, f"CLI command unexpectedly succeeded: {result.stdout}"
    
    if expected_error:
        assert expected_error in result.stderr, \
            f"Expected error '{expected_error}' not found in stderr: {result.stderr}"
```

## CLI Deployment & Distribution

### CLI-Specific Docker Template
```dockerfile
# Optimized for CLI tools
FROM python:3.12-slim

# Install uv
COPY --from=ghcr.io/astral-sh/uv:latest /uv /bin/uv

# Set environment
ENV UV_LINK_MODE=copy
ENV PYTHONUNBUFFERED=1

# Create non-root user for CLI execution
RUN useradd --create-home --shell /bin/bash cliuser

WORKDIR /app

# Copy project files
COPY --chown=cliuser:cliuser . .

# Install dependencies
RUN uv sync --frozen

# Switch to non-root user
USER cliuser

# Set CLI as entrypoint
ENTRYPOINT ["uv", "run", "project-name"]

# Default help command
CMD ["--help"]
```

### CLI Installation Instructions Template
```markdown
# Installation Instructions

## Method 1: Install from PyPI (Recommended)
```bash
# Install globally using uv
uv tool install project-name

# Or install using pipx
pipx install project-name
```

## Method 2: Install from Source
```bash
# Clone repository
git clone https://github.com/username/project-name.git
cd project-name

# Install using uv
uv tool install .
```

## Method 3: Development Installation
```bash
# Clone and setup development environment
git clone https://github.com/username/project-name.git
cd project-name
uv sync --group dev

# Run from development environment
uv run project-name --help
```

## Verify Installation
```bash
project-name --help
project-name --version
```
```

## CLI Quality Standards

### CLI-Specific Quality Gates

In addition to base quality requirements, CLI tools **MUST** meet:

- **100% CLI command coverage**: All commands and options must have tests
- **Help text completeness**: Every command must have comprehensive help
- **Error message quality**: All error messages must be actionable
- **Performance standards**: CLI startup time < 500ms, command execution appropriate to task
- **Cross-platform compatibility**: Must work on Linux, macOS, and Windows

### CLI-Specific Enforcement

Additional checks for CLI tools:
- **Command argument validation**: All argument combinations tested
- **Configuration file validation**: All config options tested  
- **Output format consistency**: All output formats produce valid results
- **Help text accuracy**: Help examples must be tested and working
- **Exit code correctness**: All error conditions return appropriate exit codes

## CLI Design Principles

### User Experience Guidelines

1. **Progressive Disclosure**: Start simple, allow complexity through options
2. **Consistent Interface**: Similar commands should work similarly  
3. **Helpful Defaults**: Choose sensible defaults, allow customization
4. **Clear Feedback**: Always indicate what the tool is doing
5. **Error Recovery**: Provide actionable error messages and suggestions
6. **Documentation**: Help text should be complete and include examples

### Performance Expectations

- **Startup time**: < 500ms for simple commands
- **Memory usage**: < 50MB for typical operations
- **Progress indication**: Show progress for operations > 2 seconds
- **Responsiveness**: Respond to Ctrl+C gracefully
- **Output streaming**: Stream output for long-running operations

This overlay provides comprehensive CLI-specific standards while building on the foundation of the base Python development standards.