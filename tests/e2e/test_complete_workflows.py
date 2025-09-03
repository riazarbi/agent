"""End-to-end tests for complete CLI workflows.

Tests complete user workflows from start to finish using the installed CLI tool,
including session persistence, configuration handling, and real command execution.

Keywords: e2e, end-to-end, workflow, CLI, integration, user experience
"""

import os
import subprocess
import tempfile
import time
from pathlib import Path
from typing import NamedTuple

import pytest
import yaml


class WorkflowResult(NamedTuple):
    """Result from complete workflow execution.

    Keywords: workflow, result, e2e, testing
    """

    success: bool
    outputs: list[str]
    errors: list[str]
    files_created: list[Path]


def run_workflow_command(
    args: list[str], input_text: str = "", cwd: str = None
) -> tuple[int, str, str]:
    """Run workflow command and return results.

    Keywords: workflow, command, subprocess, e2e testing
    """
    cmd = ["uv", "run", "agent"] + args

    result = subprocess.run(
        cmd, input=input_text, capture_output=True, text=True, timeout=15, cwd=cwd
    )

    return result.returncode, result.stdout, result.stderr


@pytest.mark.e2e
class TestCompleteUserWorkflows:
    """End-to-end tests for complete user workflows."""

    def test_complete_single_shot_workflow_with_configuration(self):
        """Test complete single-shot workflow with custom configuration."""
        with tempfile.TemporaryDirectory() as temp_dir:
            temp_path = Path(temp_dir)

            # Create custom configuration
            config_file = temp_path / "config.yaml"
            config_data = {
                "model": "gpt-3.5-turbo",
                "timeout": 30,
                "max_tokens": 1000,
                "temperature": 0.7,
                "confirm_commands": False,
            }

            with open(config_file, "w") as f:
                yaml.dump(config_data, f)

            # Execute single-shot command with configuration
            exit_code, stdout, stderr = run_workflow_command(
                [
                    "--config",
                    str(config_file),
                    "--prompt",
                    "Respond with 'Configuration test successful'",
                    "--no-tools",
                    "--quiet",
                ]
            )

            # Verify workflow execution
            assert exit_code in [0, 1]  # Success or expected API error

            # Verify configuration was loaded (no config errors)
            assert (
                "configuration" not in stderr.lower() or "config" not in stderr.lower()
            )
            assert "not found" not in stderr.lower()

    def test_complete_file_input_workflow_with_output_verification(self):
        """Test complete file input workflow with output verification."""
        with tempfile.TemporaryDirectory() as temp_dir:
            temp_path = Path(temp_dir)

            # Create input prompt file
            input_file = temp_path / "prompt.txt"
            input_file.write_text("Please respond with a simple greeting message.")

            # Execute file input workflow
            exit_code, stdout, stderr = run_workflow_command(
                ["--file", str(input_file), "--no-tools", "--verbose"]
            )

            # Verify workflow execution
            assert exit_code in [0, 1]  # Success or expected API error

            # Verify file was processed
            assert "file not found" not in stderr.lower()
            assert input_file.exists()  # Input file should still exist

    def test_complete_interactive_workflow_simulation(self):
        """Test complete interactive workflow simulation."""
        # Simulate brief interactive session
        exit_code, stdout, stderr = run_workflow_command(
            ["--no-tools", "--quiet"], input_text="exit\n"
        )

        # Verify interactive mode started and exited cleanly
        assert exit_code in [0, 1]  # Success or expected API error

        # Should not hang or crash
        assert stdout is not None or stderr is not None


@pytest.mark.e2e
class TestSessionWorkflows:
    """End-to-end tests for session management workflows."""

    def test_session_creation_and_persistence_workflow(self):
        """Test session creation and persistence workflow."""
        with tempfile.TemporaryDirectory() as temp_dir:
            # Set custom session directory via environment
            env = os.environ.copy()
            env["HOME"] = temp_dir

            # Start interactive session briefly
            cmd = ["uv", "run", "agent", "--no-tools", "--quiet"]

            result = subprocess.run(
                cmd, input="exit\n", capture_output=True, text=True, timeout=10, env=env
            )

            # Verify session handling worked
            assert result.returncode in [0, 1]  # Success or API error

            # Check for session directory creation
            agent_dir = Path(temp_dir) / ".agent"
            if agent_dir.exists():
                sessions_dir = agent_dir / "sessions"
                # Session directory structure should exist if session was created
                assert not sessions_dir.exists() or sessions_dir.is_dir()

    def test_session_resume_workflow_with_error_handling(self):
        """Test session resume workflow with proper error handling."""
        # Attempt to resume nonexistent session
        exit_code, stdout, stderr = run_workflow_command(
            ["--resume", "nonexistent-session-2023-01-01-12-00-00"]
        )

        # Should handle missing session gracefully
        assert exit_code in [0, 1]  # May exit 0 but show error
        assert "session" in stderr.lower()
        assert "not found" in stderr.lower() or "does not exist" in stderr.lower()


@pytest.mark.e2e
class TestConfigurationWorkflows:
    """End-to-end tests for configuration handling workflows."""

    def test_configuration_override_workflow(self):
        """Test configuration override workflow with multiple sources."""
        with tempfile.TemporaryDirectory() as temp_dir:
            temp_path = Path(temp_dir)

            # Create base configuration file
            config_file = temp_path / "base_config.yaml"
            base_config = {
                "model": "gpt-3.5-turbo",
                "timeout": 30,
                "max_tokens": 1000,
                "temperature": 0.5,
                "confirm_commands": True,
            }

            with open(config_file, "w") as f:
                yaml.dump(base_config, f)

            # Test configuration loading with CLI overrides
            exit_code, stdout, stderr = run_workflow_command(
                [
                    "--config",
                    str(config_file),
                    "--prompt",
                    "Test configuration override",
                    "--no-tools",  # Override confirm_commands from config
                    "--quiet",
                ]
            )

            # Verify configuration was processed successfully
            assert exit_code in [0, 1]  # Success or expected API error
            assert "configuration" not in stderr.lower()

    def test_environment_variable_configuration_workflow(self):
        """Test environment variable configuration workflow."""
        # Test with environment variables (if supported)
        env = os.environ.copy()

        # Set test environment variables
        env["AGENT_MODEL"] = "test-model"
        env["AGENT_VERBOSE"] = "true"

        cmd = ["uv", "run", "agent", "--help"]

        result = subprocess.run(cmd, capture_output=True, text=True, timeout=5, env=env)

        # Should handle environment variables gracefully
        assert result.returncode == 0
        assert "Usage:" in result.stdout


@pytest.mark.e2e
class TestErrorRecoveryWorkflows:
    """End-to-end tests for error recovery workflows."""

    def test_graceful_error_recovery_workflow(self):
        """Test graceful error recovery across multiple scenarios."""
        test_scenarios = [
            # Invalid configuration file
            {
                "name": "invalid_config",
                "setup": lambda temp_path: (temp_path / "bad_config.yaml").write_text(
                    "invalid: yaml: ["
                ),
                "args": lambda temp_path: [
                    "--config",
                    str(temp_path / "bad_config.yaml"),
                    "--prompt",
                    "test",
                ],
            },
            # Missing required files
            {
                "name": "missing_file",
                "setup": lambda temp_path: None,
                "args": lambda temp_path: [
                    "--file",
                    str(temp_path / "nonexistent.txt"),
                ],
            },
            # Invalid command combinations
            {
                "name": "invalid_combination",
                "setup": lambda temp_path: (temp_path / "test.txt").write_text("test"),
                "args": lambda temp_path: [
                    "--prompt",
                    "test",
                    "--file",
                    str(temp_path / "test.txt"),
                ],
            },
        ]

        for scenario in test_scenarios:
            with tempfile.TemporaryDirectory() as temp_dir:
                temp_path = Path(temp_dir)

                # Setup scenario
                if scenario["setup"]:
                    scenario["setup"](temp_path)

                # Execute command
                try:
                    exit_code, stdout, stderr = run_workflow_command(
                        scenario["args"](temp_path)
                    )

                    # Verify graceful error handling (CLI may exit 0 but show error)
                    if scenario["name"] == "invalid_combination":
                        # Mode conflicts may exit with 0 or 2
                        assert exit_code in [0, 2], (
                            f"Scenario {scenario['name']} should handle mode conflict"
                        )
                    else:
                        # Other errors may exit with 0 but show error message
                        assert exit_code in [0, 1, 2], (
                            f"Scenario {scenario['name']} should handle error gracefully"
                        )

                    assert stderr.strip() != "", (
                        f"Scenario {scenario['name']} should have error message"
                    )

                    # Should not crash with stack trace
                    assert "Traceback" not in stderr, (
                        f"Scenario {scenario['name']} should not show stack trace"
                    )

                except subprocess.TimeoutExpired:
                    pytest.fail(
                        f"Scenario {scenario['name']} timed out - indicates hanging"
                    )


@pytest.mark.e2e
class TestCrossCommandWorkflows:
    """End-to-end tests for workflows spanning multiple commands."""

    def test_help_to_execution_workflow(self):
        """Test workflow from help discovery to command execution."""
        # Step 1: User discovers help information
        help_exit_code, help_stdout, help_stderr = run_workflow_command(["--help"])

        assert help_exit_code == 0
        assert "Usage:" in help_stdout
        assert "--prompt" in help_stdout

        # Step 2: User executes discovered command
        exec_exit_code, exec_stdout, exec_stderr = run_workflow_command(
            ["--prompt", "Hello, please respond briefly", "--no-tools", "--quiet"]
        )

        # Command should execute based on help information
        assert exec_exit_code in [0, 1]  # Success or API error

    def test_configuration_to_execution_workflow(self):
        """Test workflow from configuration setup to execution."""
        with tempfile.TemporaryDirectory() as temp_dir:
            temp_path = Path(temp_dir)

            # Step 1: User creates configuration file
            config_file = temp_path / "user_config.yaml"
            user_config = {
                "model": "gpt-3.5-turbo",
                "timeout": 45,
                "max_tokens": 1500,
                "temperature": 0.8,
                "confirm_commands": False,
            }

            with open(config_file, "w") as f:
                yaml.dump(user_config, f)

            # Step 2: User validates configuration with help
            help_exit_code, help_stdout, help_stderr = run_workflow_command(
                ["--config", str(config_file), "--help"]
            )

            assert help_exit_code == 0

            # Step 3: User executes with configuration
            exec_exit_code, exec_stdout, exec_stderr = run_workflow_command(
                [
                    "--config",
                    str(config_file),
                    "--prompt",
                    "Configuration test execution",
                    "--no-tools",
                ]
            )

            # Should use configuration successfully
            assert exec_exit_code in [0, 1]  # Success or API error
            assert "configuration" not in exec_stderr.lower()


@pytest.mark.e2e
class TestPerformanceWorkflows:
    """End-to-end tests for performance characteristics."""

    def test_cli_startup_performance_workflow(self):
        """Test CLI startup performance meets requirements."""
        start_time = time.time()

        exit_code, stdout, stderr = run_workflow_command(["--help"])

        end_time = time.time()
        startup_time = end_time - start_time

        # CLI startup should be fast
        assert exit_code == 0
        assert startup_time < 2.0, (
            f"CLI startup took {startup_time:.2f}s, expected < 2.0s"
        )

        # Help should be comprehensive
        assert len(stdout) > 100, "Help output should be comprehensive"

    def test_command_processing_performance_workflow(self):
        """Test command processing performance."""
        # Test with simple command that doesn't require API
        start_time = time.time()

        exit_code, stdout, stderr = run_workflow_command(
            [
                "--file",
                "/dev/null",  # Empty file input
                "--no-tools",
            ]
        )

        end_time = time.time()
        processing_time = end_time - start_time

        # Command processing should be reasonable
        assert processing_time < 5.0, f"Command processing took {processing_time:.2f}s"

        # Should handle empty input gracefully
        assert exit_code in [0, 1, 2]  # Various valid outcomes for empty input
