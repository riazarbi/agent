"""CLI testing helper functions.

Provides utilities for testing CLI functionality including command execution,
result validation, and test data management for integration and e2e tests.

Keywords: CLI, testing, helpers, utilities, integration, subprocess
"""

import subprocess
import tempfile
from pathlib import Path
from typing import NamedTuple, Optional

import yaml


class CLITestResult(NamedTuple):
    """Structured result from CLI test execution.
    
    Keywords: CLI, test, result, testing utilities
    """
    exit_code: int
    stdout: str
    stderr: str
    success: bool
    timed_out: bool


def execute_cli_command(
    args: list[str], 
    input_text: str = "", 
    timeout: int = 10,
    cwd: Optional[str] = None
) -> CLITestResult:
    """Execute CLI command and return structured result.
    
    Executes the python_agent CLI with given arguments and returns
    comprehensive result information for test validation.
    
    Keywords: CLI, command, execution, testing, subprocess
    
    Args:
        args: Command line arguments to pass to agent CLI
        input_text: Optional stdin input for interactive commands  
        timeout: Command timeout in seconds (default: 10)
        cwd: Optional working directory for command execution
        
    Returns:
        CLITestResult with exit code, output streams, and status flags
    """
    cmd = ["uv", "run", "agent"] + args
    
    try:
        result = subprocess.run(
            cmd,
            input=input_text,
            capture_output=True,
            text=True,
            timeout=timeout,
            cwd=cwd
        )
        return CLITestResult(
            exit_code=result.returncode,
            stdout=result.stdout,
            stderr=result.stderr,
            success=result.returncode == 0,
            timed_out=False
        )
    except subprocess.TimeoutExpired:
        return CLITestResult(
            exit_code=-1,
            stdout="",
            stderr="Command timed out",
            success=False,
            timed_out=True
        )


def create_test_config_file(config_data: dict, temp_dir: Path) -> Path:
    """Create temporary configuration file for testing.
    
    Creates a YAML configuration file in the specified directory
    with the provided configuration data.
    
    Keywords: configuration, test, YAML, file creation, testing utilities
    
    Args:
        config_data: Dictionary containing configuration values
        temp_dir: Directory to create the config file in
        
    Returns:
        Path to the created configuration file
    """
    config_file = temp_dir / "test_config.yaml"
    
    with open(config_file, 'w') as f:
        yaml.dump(config_data, f)
    
    return config_file


def create_test_prompt_file(content: str, temp_dir: Path, filename: str = "test_prompt.txt") -> Path:
    """Create temporary prompt file for testing.
    
    Creates a text file with the specified content for use in
    file input mode testing.
    
    Keywords: prompt, file, test, creation, file input testing
    
    Args:
        content: Text content for the prompt file
        temp_dir: Directory to create the file in  
        filename: Name for the prompt file (default: test_prompt.txt)
        
    Returns:
        Path to the created prompt file
    """
    prompt_file = temp_dir / filename
    prompt_file.write_text(content)
    return prompt_file


def assert_cli_success(result: CLITestResult, allow_api_errors: bool = True) -> None:
    """Assert that CLI command executed successfully.
    
    Validates that CLI command completed successfully, with optional
    allowance for API-related errors that are expected in testing.
    
    Keywords: assertion, CLI, success, validation, testing
    
    Args:
        result: CLITestResult from command execution
        allow_api_errors: Whether to allow exit code 1 for API errors
        
    Raises:
        AssertionError: If command did not complete successfully
    """
    if allow_api_errors:
        assert result.exit_code in [0, 1], f"Command failed with exit code {result.exit_code}: {result.stderr}"
    else:
        assert result.success, f"Command failed with exit code {result.exit_code}: {result.stderr}"
    
    assert not result.timed_out, "Command timed out"


def assert_cli_error(result: CLITestResult, expected_exit_code: int = None) -> None:
    """Assert that CLI command failed as expected.
    
    Validates that CLI command failed with appropriate error handling
    and optionally checks for specific exit code.
    
    Keywords: assertion, CLI, error, validation, failure testing
    
    Args:
        result: CLITestResult from command execution
        expected_exit_code: Optional specific exit code to check for
        
    Raises:
        AssertionError: If command succeeded when it should have failed
    """
    assert not result.success, f"Command should have failed but succeeded with stdout: {result.stdout}"
    assert not result.timed_out, "Command should not have timed out"
    assert result.stderr.strip() != "", "Error command should provide error message"
    
    if expected_exit_code is not None:
        assert result.exit_code == expected_exit_code, f"Expected exit code {expected_exit_code}, got {result.exit_code}"


def assert_help_output_complete(stdout: str) -> None:
    """Assert that help output contains all required elements.
    
    Validates that CLI help output includes all necessary usage information,
    options, and examples for comprehensive user guidance.
    
    Keywords: help, assertion, CLI, documentation, usage validation
    
    Args:
        stdout: Help output from CLI command
        
    Raises:
        AssertionError: If help output is incomplete
    """
    required_elements = [
        "Usage:",
        "--prompt",
        "--file", 
        "--resume",
        "--config",
        "--allow-tools",
        "--no-tools",
        "--confirm",
        "--no-confirm",
        "--verbose",
        "--quiet",
        "--help"
    ]
    
    for element in required_elements:
        assert element in stdout, f"Help output missing required element: {element}"


def create_invalid_yaml_file(temp_dir: Path, filename: str = "invalid.yaml") -> Path:
    """Create invalid YAML file for error testing.
    
    Creates a file with malformed YAML content to test
    configuration error handling.
    
    Keywords: YAML, invalid, error testing, configuration, file creation
    
    Args:
        temp_dir: Directory to create the file in
        filename: Name for the invalid YAML file
        
    Returns:
        Path to the created invalid YAML file
    """
    invalid_file = temp_dir / filename
    invalid_file.write_text("invalid: yaml: content: [missing bracket")
    return invalid_file


def get_default_test_config() -> dict:
    """Get default configuration for CLI testing.
    
    Returns a standard configuration dictionary suitable for
    testing CLI functionality without external API dependencies.
    
    Keywords: configuration, default, testing, CLI setup
    
    Returns:
        Dictionary containing default test configuration values
    """
    return {
        'model': 'gpt-3.5-turbo',
        'timeout': 30,
        'max_tokens': 1000,
        'temperature': 0.7,
        'confirm_commands': False
    }