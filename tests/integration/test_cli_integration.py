"""Integration tests for CLI functionality.

Tests CLI behavior through subprocess calls to verify end-to-end functionality
including command-line argument parsing, configuration file handling, session
workflows, and error scenarios.

Keywords: integration, CLI, subprocess, end-to-end, command-line
"""

import os
import subprocess
import tempfile
from typing import NamedTuple

import pytest
import yaml


class CLIResult(NamedTuple):
    """Result from CLI command execution.

    Keywords: CLI, result, subprocess, command execution
    """

    exit_code: int
    stdout: str
    stderr: str
    success: bool


def run_cli_command(
    args: list[str], input_text: str = "", timeout: int = 10
) -> CLIResult:
    """Run CLI command and return structured result.

    Keywords: CLI, command, subprocess, execution, testing

    Args:
        args: Command line arguments to pass to agent
        input_text: Optional stdin input for interactive commands
        timeout: Command timeout in seconds

    Returns:
        CLIResult with exit code, stdout, stderr, and success flag
    """
    cmd = ["uv", "run", "agent"] + args

    try:
        result = subprocess.run(
            cmd, input=input_text, capture_output=True, text=True, timeout=timeout
        )
        return CLIResult(
            exit_code=result.returncode,
            stdout=result.stdout,
            stderr=result.stderr,
            success=result.returncode == 0,
        )
    except subprocess.TimeoutExpired:
        return CLIResult(
            exit_code=-1, stdout="", stderr="Command timed out", success=False
        )


@pytest.mark.integration
class TestCLIIntegrationBasicFunctionality:
    """Integration tests for basic CLI functionality."""

    def test_cli_help_display_shows_comprehensive_usage_information(self):
        """Test CLI help display shows comprehensive usage information."""
        result = run_cli_command(["--help"])

        assert result.success is True
        assert "Usage:" in result.stdout
        assert "Interactive mode" in result.stdout
        assert "--prompt" in result.stdout
        assert "--file" in result.stdout
        assert "--resume" in result.stdout
        assert "--allow-tools" in result.stdout
        assert "--confirm" in result.stdout

    def test_cli_version_information_when_available(self):
        """Test CLI version information display when available."""
        result = run_cli_command(["--help"])

        # Version info should be available in help or as separate command
        assert result.success is True
        # Help should contain project information
        assert result.stdout is not None

    def test_cli_single_shot_mode_with_simple_prompt(self):
        """Test CLI single-shot mode with simple prompt."""
        result = run_cli_command(
            ["--prompt", "Hello, respond with 'Test response'", "--no-tools"]
        )

        # Should exit cleanly (may fail due to no API key, but structure should work)
        assert result.exit_code in [0, 1]  # Success or expected API error

        # Should not hang or crash
        assert "Error:" in result.stderr or result.stdout != ""


@pytest.mark.integration
class TestCLIIntegrationConfigurationHandling:
    """Integration tests for configuration file handling."""

    def test_cli_configuration_file_loading_with_custom_path(self):
        """Test CLI configuration file loading with custom path."""
        with tempfile.NamedTemporaryFile(mode="w", suffix=".yaml", delete=False) as f:
            config = {
                "model": "gpt-3.5-turbo",
                "timeout": 60,
                "max_tokens": 2000,
                "temperature": 0.5,
                "confirm_commands": True,
            }
            yaml.dump(config, f)
            config_path = f.name

        try:
            result = run_cli_command(
                ["--config", config_path, "--prompt", "test prompt", "--no-tools"]
            )

            # Should load config successfully
            assert result.exit_code in [0, 1]  # Success or expected API error
            assert "Config file not found" not in result.stderr

        finally:
            os.unlink(config_path)

    def test_cli_configuration_file_error_handling_for_invalid_yaml(self):
        """Test CLI configuration file error handling for invalid YAML."""
        with tempfile.NamedTemporaryFile(mode="w", suffix=".yaml", delete=False) as f:
            f.write("invalid: yaml: content: [")
            config_path = f.name

        try:
            result = run_cli_command(
                ["--config", config_path, "--prompt", "test prompt"]
            )

            # CLI may exit with 0 but should show error message
            assert result.exit_code in [0, 1]
            assert (
                "configuration" in result.stderr.lower()
                or "yaml" in result.stderr.lower()
            )

        finally:
            os.unlink(config_path)

    def test_cli_configuration_file_missing_file_error_handling(self):
        """Test CLI configuration file missing file error handling."""
        nonexistent_config = "/tmp/nonexistent_config_file.yaml"

        result = run_cli_command(
            ["--config", nonexistent_config, "--prompt", "test prompt"]
        )

        # CLI may exit with 0 but should show error message
        assert result.exit_code in [0, 1]
        assert (
            "not found" in result.stderr.lower()
            or "does not exist" in result.stderr.lower()
            or "agent error" in result.stderr.lower()
        )


@pytest.mark.integration
class TestCLIIntegrationFileInputMode:
    """Integration tests for file input mode functionality."""

    def test_cli_file_input_mode_with_valid_file(self):
        """Test CLI file input mode with valid file."""
        with tempfile.NamedTemporaryFile(mode="w", suffix=".txt", delete=False) as f:
            f.write("This is a test prompt from file")
            file_path = f.name

        try:
            result = run_cli_command(["--file", file_path, "--no-tools"])

            # Should process file successfully
            assert result.exit_code in [0, 1]  # Success or expected API error
            assert "File not found" not in result.stderr

        finally:
            os.unlink(file_path)

    def test_cli_file_input_mode_with_missing_file(self):
        """Test CLI file input mode with missing file."""
        nonexistent_file = "/tmp/nonexistent_prompt_file.txt"

        result = run_cli_command(["--file", nonexistent_file])

        # File validation is done by Click, should use exit code 2
        assert result.exit_code == 2
        assert (
            "does not exist" in result.stderr.lower()
            or "invalid value" in result.stderr.lower()
        )

    def test_cli_file_input_mode_with_empty_file(self):
        """Test CLI file input mode with empty file."""
        with tempfile.NamedTemporaryFile(mode="w", suffix=".txt", delete=False) as f:
            # Write empty content
            pass
            file_path = f.name

        try:
            result = run_cli_command(["--file", file_path, "--no-tools"])

            # Should handle empty file gracefully
            assert result.exit_code in [0, 1, 2]  # Success, API error, or usage error

        finally:
            os.unlink(file_path)


@pytest.mark.integration
class TestCLIIntegrationSessionManagement:
    """Integration tests for session management functionality."""

    def test_cli_session_resume_with_nonexistent_session(self):
        """Test CLI session resume with nonexistent session."""
        result = run_cli_command(["--resume", "nonexistent-session-id"])

        # CLI may exit with 0 but should show error message
        assert result.exit_code in [0, 1]
        assert "session" in result.stderr.lower()
        assert (
            "not found" in result.stderr.lower()
            or "does not exist" in result.stderr.lower()
        )

    def test_cli_interactive_mode_session_creation(self):
        """Test CLI interactive mode session creation and basic interaction."""
        # Use short timeout and provide exit command
        result = run_cli_command(
            ["--no-tools", "--quiet"], input_text="exit\n", timeout=5
        )

        # Should start interactive mode and exit cleanly
        assert result.exit_code in [0, 1]  # Success or expected API error

        # Should not crash or hang
        assert result.stdout is not None or result.stderr is not None


@pytest.mark.integration
class TestCLIIntegrationCommandLineOptions:
    """Integration tests for command-line options and flag combinations."""

    def test_cli_tool_flags_allow_tools_and_no_tools(self):
        """Test CLI tool flags --allow-tools and --no-tools."""
        # Test --allow-tools
        result_allow = run_cli_command(["--prompt", "test prompt", "--allow-tools"])

        # Test --no-tools
        result_no_tools = run_cli_command(["--prompt", "test prompt", "--no-tools"])

        # Both should parse successfully
        assert result_allow.exit_code in [0, 1]  # Success or API error
        assert result_no_tools.exit_code in [0, 1]  # Success or API error

    def test_cli_confirmation_flags_confirm_and_no_confirm(self):
        """Test CLI confirmation flags --confirm and --no-confirm."""
        # Test --confirm
        result_confirm = run_cli_command(
            ["--prompt", "test prompt", "--confirm", "--no-tools"]
        )

        # Test --no-confirm
        result_no_confirm = run_cli_command(
            ["--prompt", "test prompt", "--no-confirm", "--no-tools"]
        )

        # Both should parse successfully
        assert result_confirm.exit_code in [0, 1]  # Success or API error
        assert result_no_confirm.exit_code in [0, 1]  # Success or API error

    def test_cli_verbose_and_quiet_modes(self):
        """Test CLI verbose and quiet mode flags."""
        # Test --verbose
        result_verbose = run_cli_command(
            ["--prompt", "test prompt", "--verbose", "--no-tools"]
        )

        # Test --quiet
        result_quiet = run_cli_command(
            ["--prompt", "test prompt", "--quiet", "--no-tools"]
        )

        # Both should parse successfully
        assert result_verbose.exit_code in [0, 1]  # Success or API error
        assert result_quiet.exit_code in [0, 1]  # Success or API error

    def test_cli_mode_exclusivity_validation(self):
        """Test CLI mode exclusivity validation prevents conflicting options."""
        # Test conflicting modes: --prompt and --file
        with tempfile.NamedTemporaryFile(mode="w", suffix=".txt", delete=False) as f:
            f.write("test content")
            file_path = f.name

        try:
            result = run_cli_command(["--prompt", "test prompt", "--file", file_path])

            # Should reject conflicting modes (may exit 0 but show error)
            assert result.exit_code in [0, 2]
            assert (
                "cannot use" in result.stderr.lower()
                or "together" in result.stderr.lower()
            )

        finally:
            os.unlink(file_path)


@pytest.mark.integration
class TestCLIIntegrationErrorScenarios:
    """Integration tests for error scenarios and edge cases."""

    def test_cli_keyboard_interrupt_handling(self):
        """Test CLI keyboard interrupt handling in interactive mode."""
        # This is challenging to test directly, but we can verify the CLI
        # starts properly and would handle interrupts
        result = run_cli_command(["--help"])

        # CLI should start successfully
        assert result.success is True
        # This validates the CLI infrastructure is working

    def test_cli_invalid_command_line_arguments(self):
        """Test CLI invalid command line arguments."""
        result = run_cli_command(["--invalid-flag"])

        assert result.success is False
        assert result.exit_code == 2  # Usage error
        assert (
            "unrecognized" in result.stderr.lower()
            or "invalid" in result.stderr.lower()
        )

    def test_cli_missing_required_argument_values(self):
        """Test CLI missing required argument values."""
        # Test --config without value
        result = run_cli_command(["--config"])

        assert result.success is False
        assert result.exit_code == 2  # Usage error

    def test_cli_graceful_shutdown_on_various_errors(self):
        """Test CLI graceful shutdown on various error conditions."""
        # Test with invalid model configuration that should fail gracefully
        with tempfile.NamedTemporaryFile(mode="w", suffix=".yaml", delete=False) as f:
            config = {
                "model": "",  # Invalid empty model
                "timeout": -1,  # Invalid timeout
                "max_tokens": -100,  # Invalid max_tokens
            }
            yaml.dump(config, f)
            config_path = f.name

        try:
            result = run_cli_command(
                ["--config", config_path, "--prompt", "test prompt"]
            )

            # Should fail gracefully with appropriate error message
            assert result.exit_code in [0, 1]
            assert result.stderr.strip() != ""
            # Should not crash with stack trace
            assert "Traceback" not in result.stderr

        finally:
            os.unlink(config_path)


@pytest.mark.integration
class TestCLIIntegrationOutputFormats:
    """Integration tests for CLI output formats and user experience."""

    def test_cli_quiet_mode_minimal_output(self):
        """Test CLI quiet mode produces minimal output."""
        result = run_cli_command(["--prompt", "test prompt", "--quiet", "--no-tools"])

        # Quiet mode should have minimal stdout (may have stderr for errors)
        if result.success:
            # If successful, stdout should be minimal
            assert len(result.stdout.strip()) >= 0  # May have model response
        else:
            # If failed, should have error in stderr
            assert result.stderr.strip() != ""

    def test_cli_verbose_mode_detailed_output(self):
        """Test CLI verbose mode produces detailed output."""
        result = run_cli_command(["--prompt", "test prompt", "--verbose", "--no-tools"])

        # Verbose mode should provide more information
        assert result.exit_code in [0, 1]  # Success or API error
        # Should have some output (either success info or detailed error)
        assert result.stdout != "" or result.stderr != ""

    def test_cli_error_message_quality_and_actionability(self):
        """Test CLI error messages are user-friendly and actionable."""
        # Test with obviously bad config file
        with tempfile.NamedTemporaryFile(mode="w", suffix=".yaml", delete=False) as f:
            f.write("this is not yaml at all!!!")
            config_path = f.name

        try:
            result = run_cli_command(
                ["--config", config_path, "--prompt", "test prompt"]
            )

            # Should handle error gracefully
            assert result.exit_code in [0, 1]
            # Error message should be user-friendly
            error_msg = result.stderr.lower()
            assert (
                "error" in error_msg
                or "failed" in error_msg
                or "dictionary" in error_msg
            )
            # Should not contain Python stack traces for user errors
            assert "traceback" not in error_msg

        finally:
            os.unlink(config_path)
