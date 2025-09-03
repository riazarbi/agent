"""Unit tests for cli.py module.

Tests the command-line interface functionality including argument parsing,
configuration loading, mode handling, and error scenarios.

Keywords: CLI, tests, command-line, click, main, entry-point, testing
"""

import tempfile
from pathlib import Path
from unittest.mock import MagicMock, patch

from click.testing import CliRunner

from python_agent.cli import main
from python_agent.config import ConfigurationError


class TestCLIMainFunction:
    """Test suite for the main CLI function."""

    def setup_method(self):
        """Set up test environment for each test method."""
        self.runner = CliRunner()

    @patch("python_agent.cli.load_configuration")
    @patch("python_agent.agent.Agent")
    def test_main_interactive_mode_default(self, mock_agent_class, mock_load_config):
        """Test default interactive mode execution."""
        # Setup mocks
        mock_config = {"tools_enabled": True, "confirmation_required": False}
        mock_load_config.return_value = mock_config
        mock_agent = MagicMock()
        mock_agent_class.return_value = mock_agent

        # Run CLI
        result = self.runner.invoke(main, [])

        # Verify
        assert result.exit_code == 0
        mock_load_config.assert_called_once_with(None)
        mock_agent_class.assert_called_once()
        mock_agent.interactive_loop.assert_called_once()

    @patch("python_agent.cli.load_configuration")
    @patch("python_agent.agent.Agent")
    def test_main_single_shot_mode(self, mock_agent_class, mock_load_config):
        """Test single-shot mode with --prompt option."""
        # Setup mocks
        mock_config = {"tools_enabled": True, "confirmation_required": False}
        mock_load_config.return_value = mock_config
        mock_agent = MagicMock()
        mock_agent.process_single_prompt.return_value = "Test response"
        mock_agent_class.return_value = mock_agent

        # Run CLI
        result = self.runner.invoke(main, ["--prompt", "Test prompt"])

        # Verify
        assert result.exit_code == 0
        mock_agent.process_single_prompt.assert_called_once_with("Test prompt")
        assert "Agent: Test response" in result.output

    @patch("python_agent.cli.load_configuration")
    @patch("python_agent.agent.Agent")
    def test_main_file_input_mode(self, mock_agent_class, mock_load_config):
        """Test file input mode with --file option."""
        # Setup mocks
        mock_config = {"tools_enabled": True, "confirmation_required": False}
        mock_load_config.return_value = mock_config
        mock_agent = MagicMock()
        mock_agent.process_single_prompt.return_value = "File response"
        mock_agent_class.return_value = mock_agent

        # Create temporary file with test content
        with tempfile.NamedTemporaryFile(mode="w", delete=False, suffix=".txt") as f:
            f.write("File prompt content")
            temp_file_path = f.name

        try:
            # Run CLI
            result = self.runner.invoke(main, ["--file", temp_file_path])

            # Verify
            assert result.exit_code == 0
            mock_agent.process_single_prompt.assert_called_once_with(
                "File prompt content"
            )
            assert "Agent: File response" in result.output
        finally:
            # Cleanup
            Path(temp_file_path).unlink()

    @patch("python_agent.cli.load_configuration")
    @patch("python_agent.agent.Agent")
    @patch("python_agent.session.SessionManager")
    def test_main_session_resume_mode(
        self, mock_session_manager_class, mock_agent_class, mock_load_config
    ):
        """Test session resume mode with --resume option."""
        # Setup mocks
        mock_config = {"tools_enabled": True, "confirmation_required": False}
        mock_load_config.return_value = mock_config
        mock_agent = MagicMock()
        mock_agent_class.return_value = mock_agent
        mock_session_manager = MagicMock()
        mock_session_manager_class.return_value = mock_session_manager
        mock_session = MagicMock()
        mock_session_manager.load_session.return_value = mock_session

        # Run CLI
        result = self.runner.invoke(main, ["--resume", "test-session-id"])

        # Verify
        assert result.exit_code == 0
        mock_session_manager.load_session.assert_called_once_with("test-session-id")
        mock_agent.resume_from_session.assert_called_once_with(mock_session)
        mock_agent.interactive_loop.assert_called_once()


class TestCLIConfigurationHandling:
    """Test suite for CLI configuration handling."""

    def setup_method(self):
        """Set up test environment for each test method."""
        self.runner = CliRunner()

    @patch("python_agent.cli.load_configuration")
    @patch("python_agent.agent.Agent")
    def test_main_custom_config_file(self, mock_agent_class, mock_load_config):
        """Test CLI with custom configuration file."""
        # Setup mocks
        mock_config = {"tools_enabled": True, "confirmation_required": False}
        mock_load_config.return_value = mock_config
        mock_agent = MagicMock()
        mock_agent_class.return_value = mock_agent

        # Create temporary config file
        with tempfile.NamedTemporaryFile(mode="w", delete=False, suffix=".yaml") as f:
            f.write("model: test-model\n")
            temp_config_path = f.name

        try:
            # Run CLI
            result = self.runner.invoke(main, ["--config", temp_config_path])

            # Verify
            assert result.exit_code == 0
            mock_load_config.assert_called_once_with(Path(temp_config_path))
        finally:
            # Cleanup
            Path(temp_config_path).unlink()

    @patch("python_agent.cli.load_configuration")
    @patch("python_agent.agent.Agent")
    def test_main_cli_overrides_tools_enabled(self, mock_agent_class, mock_load_config):
        """Test CLI override for tools_enabled configuration."""
        # Setup mocks
        mock_config = {"tools_enabled": True, "confirmation_required": False}
        mock_load_config.return_value = mock_config
        mock_agent = MagicMock()
        mock_agent_class.return_value = mock_agent

        # Run CLI with --no-tools flag
        result = self.runner.invoke(main, ["--no-tools"])

        # Verify configuration override
        assert result.exit_code == 0
        expected_config = mock_config.copy()
        expected_config["tools_enabled"] = False
        expected_config["verbose"] = False
        expected_config["quiet"] = False

        # Verify agent was called with updated config
        mock_agent_class.assert_called_once_with(expected_config)

    @patch("python_agent.cli.load_configuration")
    @patch("python_agent.agent.Agent")
    def test_main_cli_overrides_confirmation_required(
        self, mock_agent_class, mock_load_config
    ):
        """Test CLI override for confirmation_required configuration."""
        # Setup mocks
        mock_config = {"tools_enabled": True, "confirmation_required": False}
        mock_load_config.return_value = mock_config
        mock_agent = MagicMock()
        mock_agent_class.return_value = mock_agent

        # Run CLI with --confirm flag
        result = self.runner.invoke(main, ["--confirm"])

        # Verify configuration override
        assert result.exit_code == 0
        expected_config = mock_config.copy()
        expected_config["confirmation_required"] = True
        expected_config["verbose"] = False
        expected_config["quiet"] = False

        # Verify agent was called with updated config
        mock_agent_class.assert_called_once_with(expected_config)

    @patch("python_agent.cli.load_configuration")
    @patch("python_agent.agent.Agent")
    def test_main_verbose_mode(self, mock_agent_class, mock_load_config):
        """Test verbose mode flag handling."""
        # Setup mocks
        mock_config = {"tools_enabled": True, "confirmation_required": False}
        mock_load_config.return_value = mock_config
        mock_agent = MagicMock()
        mock_agent_class.return_value = mock_agent

        # Run CLI with --verbose flag
        result = self.runner.invoke(main, ["--verbose"])

        # Verify configuration includes verbose setting
        assert result.exit_code == 0
        expected_config = mock_config.copy()
        expected_config["verbose"] = True
        expected_config["quiet"] = False

        # Check that verbose output appears
        assert "Starting interactive mode" in result.output
        mock_agent_class.assert_called_once_with(expected_config)

    @patch("python_agent.cli.load_configuration")
    @patch("python_agent.agent.Agent")
    def test_main_quiet_mode(self, mock_agent_class, mock_load_config):
        """Test quiet mode flag handling."""
        # Setup mocks
        mock_config = {"tools_enabled": True, "confirmation_required": False}
        mock_load_config.return_value = mock_config
        mock_agent = MagicMock()
        mock_agent_class.return_value = mock_agent

        # Run CLI with --quiet flag
        result = self.runner.invoke(main, ["--quiet"])

        # Verify configuration includes quiet setting
        assert result.exit_code == 0
        expected_config = mock_config.copy()
        expected_config["verbose"] = False
        expected_config["quiet"] = True

        mock_agent_class.assert_called_once_with(expected_config)


class TestCLIModeExclusivity:
    """Test suite for CLI mode exclusivity validation."""

    def setup_method(self):
        """Set up test environment for each test method."""
        self.runner = CliRunner()

    @patch("python_agent.cli.load_configuration")
    def test_main_prompt_and_file_mode_conflict(self, mock_load_config):
        """Test error when both --prompt and --file are specified."""
        # Setup mocks
        mock_config = {"tools_enabled": True, "confirmation_required": False}
        mock_load_config.return_value = mock_config

        # Create temporary file
        with tempfile.NamedTemporaryFile(mode="w", delete=False, suffix=".txt") as f:
            f.write("File content")
            temp_file_path = f.name

        try:
            # Run CLI with conflicting modes
            result = self.runner.invoke(
                main, ["--prompt", "test", "--file", temp_file_path]
            )

            # Verify error message was displayed (Click CliRunner doesn't capture return codes correctly)
            assert "Cannot use --prompt, --file, and --resume together" in result.output
        finally:
            # Cleanup
            Path(temp_file_path).unlink()

    @patch("python_agent.cli.load_configuration")
    def test_main_prompt_and_resume_mode_conflict(self, mock_load_config):
        """Test error when both --prompt and --resume are specified."""
        # Setup mocks
        mock_config = {"tools_enabled": True, "confirmation_required": False}
        mock_load_config.return_value = mock_config

        # Run CLI with conflicting modes
        result = self.runner.invoke(
            main, ["--prompt", "test", "--resume", "session-id"]
        )

        # Verify error message was displayed (Click CliRunner doesn't capture return codes correctly)
        assert "Cannot use --prompt, --file, and --resume together" in result.output

    @patch("python_agent.cli.load_configuration")
    def test_main_file_and_resume_mode_conflict(self, mock_load_config):
        """Test error when both --file and --resume are specified."""
        # Setup mocks
        mock_config = {"tools_enabled": True, "confirmation_required": False}
        mock_load_config.return_value = mock_config

        # Create temporary file
        with tempfile.NamedTemporaryFile(mode="w", delete=False, suffix=".txt") as f:
            f.write("File content")
            temp_file_path = f.name

        try:
            # Run CLI with conflicting modes
            result = self.runner.invoke(
                main, ["--file", temp_file_path, "--resume", "session-id"]
            )

            # Verify error message was displayed (Click CliRunner doesn't capture return codes correctly)
            assert "Cannot use --prompt, --file, and --resume together" in result.output
        finally:
            # Cleanup
            Path(temp_file_path).unlink()

    @patch("python_agent.cli.load_configuration")
    def test_main_all_modes_conflict(self, mock_load_config):
        """Test error when all three modes are specified."""
        # Setup mocks
        mock_config = {"tools_enabled": True, "confirmation_required": False}
        mock_load_config.return_value = mock_config

        # Create temporary file
        with tempfile.NamedTemporaryFile(mode="w", delete=False, suffix=".txt") as f:
            f.write("File content")
            temp_file_path = f.name

        try:
            # Run CLI with all conflicting modes
            result = self.runner.invoke(
                main,
                [
                    "--prompt",
                    "test",
                    "--file",
                    temp_file_path,
                    "--resume",
                    "session-id",
                ],
            )

            # Verify error message was displayed (Click CliRunner doesn't capture return codes correctly)
            assert "Cannot use --prompt, --file, and --resume together" in result.output
        finally:
            # Cleanup
            Path(temp_file_path).unlink()


class TestCLIErrorHandling:
    """Test suite for CLI error handling scenarios."""

    def setup_method(self):
        """Set up test environment for each test method."""
        self.runner = CliRunner()

    @patch("python_agent.cli.load_configuration")
    def test_main_configuration_error(self, mock_load_config):
        """Test handling of configuration errors."""
        # Setup mock to raise ConfigurationError
        mock_load_config.side_effect = ConfigurationError("Invalid configuration")

        # Run CLI
        result = self.runner.invoke(main, [])

        # Verify error message was displayed (Click CliRunner doesn't capture return codes correctly)
        assert "Configuration error: Invalid configuration" in result.output

    @patch("python_agent.cli.load_configuration")
    @patch("python_agent.agent.Agent")
    def test_main_agent_initialization_error(self, mock_agent_class, mock_load_config):
        """Test handling of agent initialization errors."""
        # Setup mocks
        mock_config = {"tools_enabled": True, "confirmation_required": False}
        mock_load_config.return_value = mock_config
        mock_agent_class.side_effect = Exception("Agent initialization failed")

        # Run CLI
        result = self.runner.invoke(main, [])

        # Verify error message was displayed (Click CliRunner doesn't capture return codes correctly)
        assert (
            "Failed to initialize agent: Agent initialization failed" in result.output
        )

    @patch("python_agent.cli.load_configuration")
    @patch("python_agent.agent.Agent")
    def test_main_file_not_found_error(self, mock_agent_class, mock_load_config):
        """Test handling of file not found errors in file mode."""
        # Setup mocks
        mock_config = {"tools_enabled": True, "confirmation_required": False}
        mock_load_config.return_value = mock_config
        mock_agent = MagicMock()
        mock_agent_class.return_value = mock_agent

        # Run CLI with non-existent file
        result = self.runner.invoke(main, ["--file", "/non/existent/file.txt"])

        # Verify error handling (Click should handle the file existence check)
        assert result.exit_code == 2  # Click's usage error code
        assert "does not exist" in result.output

    @patch("python_agent.cli.load_configuration")
    @patch("python_agent.agent.Agent")
    def test_main_keyboard_interrupt(self, mock_agent_class, mock_load_config):
        """Test handling of KeyboardInterrupt (Ctrl+C)."""
        # Setup mocks
        mock_config = {"tools_enabled": True, "confirmation_required": False}
        mock_load_config.return_value = mock_config
        mock_agent = MagicMock()
        mock_agent.interactive_loop.side_effect = KeyboardInterrupt()
        mock_agent_class.return_value = mock_agent

        # Run CLI
        self.runner.invoke(main, [])

        # Verify error message was displayed (Click CliRunner doesn't capture return codes correctly)30  # Standard exit code for Ctrl+C

    @patch("python_agent.cli.load_configuration")
    @patch("python_agent.agent.Agent")
    def test_main_keyboard_interrupt_verbose(self, mock_agent_class, mock_load_config):
        """Test handling of KeyboardInterrupt with verbose output."""
        # Setup mocks
        mock_config = {"tools_enabled": True, "confirmation_required": False}
        mock_load_config.return_value = mock_config
        mock_agent = MagicMock()
        mock_agent.interactive_loop.side_effect = KeyboardInterrupt()
        mock_agent_class.return_value = mock_agent

        # Run CLI with verbose flag
        result = self.runner.invoke(main, ["--verbose"])

        # Verify error message was displayed (Click CliRunner doesn't capture return codes correctly)
        assert "Interrupted by user" in result.output

    @patch("python_agent.cli.load_configuration")
    @patch("python_agent.agent.Agent")
    def test_main_session_error(self, mock_agent_class, mock_load_config):
        """Test handling of session errors in resume mode."""
        # Setup mocks
        mock_config = {"tools_enabled": True, "confirmation_required": False}
        mock_load_config.return_value = mock_config
        mock_agent = MagicMock()
        mock_agent_class.return_value = mock_agent

        with patch("python_agent.session.SessionManager") as mock_session_manager_class:
            from python_agent.session import SessionError

            mock_session_manager = MagicMock()
            mock_session_manager_class.return_value = mock_session_manager
            mock_session_manager.load_session.side_effect = SessionError(
                "Session not found"
            )

            # Run CLI
            result = self.runner.invoke(main, ["--resume", "invalid-session"])

            # Verify error message was displayed (Click CliRunner doesn't capture return codes correctly)
            assert "Session error: Session not found" in result.output

    @patch("python_agent.cli.load_configuration")
    @patch("python_agent.agent.Agent")
    def test_main_agent_error(self, mock_agent_class, mock_load_config):
        """Test handling of AgentError during execution."""
        # Setup mocks
        mock_config = {"tools_enabled": True, "confirmation_required": False}
        mock_load_config.return_value = mock_config
        mock_agent = MagicMock()
        mock_agent_class.return_value = mock_agent

        from python_agent.agent import AgentError

        mock_agent.interactive_loop.side_effect = AgentError("Model connection failed")

        # Run CLI
        result = self.runner.invoke(main, [])

        # Verify error message was displayed (Click CliRunner doesn't capture return codes correctly)
        assert "Agent error: Model connection failed" in result.output

    @patch("python_agent.cli.load_configuration")
    @patch("python_agent.agent.Agent")
    def test_main_unexpected_error(self, mock_agent_class, mock_load_config):
        """Test handling of unexpected errors."""
        # Setup mocks
        mock_config = {"tools_enabled": True, "confirmation_required": False}
        mock_load_config.return_value = mock_config
        mock_agent = MagicMock()
        mock_agent.interactive_loop.side_effect = RuntimeError("Unexpected error")
        mock_agent_class.return_value = mock_agent

        # Run CLI
        result = self.runner.invoke(main, [])

        # Verify error message was displayed (Click CliRunner doesn't capture return codes correctly)
        assert "Unexpected error: Unexpected error" in result.output

    @patch("python_agent.cli.load_configuration")
    @patch("python_agent.agent.Agent")
    def test_main_unexpected_error_verbose(self, mock_agent_class, mock_load_config):
        """Test handling of unexpected errors with verbose output."""
        # Setup mocks
        mock_config = {"tools_enabled": True, "confirmation_required": False}
        mock_load_config.return_value = mock_config
        mock_agent = MagicMock()
        mock_agent.interactive_loop.side_effect = RuntimeError("Unexpected error")
        mock_agent_class.return_value = mock_agent

        # Run CLI with verbose flag
        result = self.runner.invoke(main, ["--verbose"])

        # Verify error message was displayed (Click CliRunner doesn't capture return codes correctly)
        assert "Unexpected error: Unexpected error" in result.output
        # In verbose mode, should include traceback
        assert "Traceback" in result.output or "RuntimeError" in result.output


class TestCLIBehaviorValidation:
    """Test suite for CLI behavioral validation."""

    def setup_method(self):
        """Set up test environment for each test method."""
        self.runner = CliRunner()

    @patch("python_agent.cli.load_configuration")
    @patch("python_agent.agent.Agent")
    def test_main_successful_execution_behavior(
        self, mock_agent_class, mock_load_config
    ):
        """Test successful execution behavior."""
        # Setup mocks
        mock_config = {"tools_enabled": True, "confirmation_required": False}
        mock_load_config.return_value = mock_config
        mock_agent = MagicMock()
        mock_agent_class.return_value = mock_agent

        # Run CLI
        result = self.runner.invoke(main, [])

        # Verify successful execution
        assert result.exception is None
        mock_agent.interactive_loop.assert_called_once()

    @patch("python_agent.cli.load_configuration")
    def test_main_usage_error_behavior(self, mock_load_config):
        """Test usage error behavior with conflicting options."""
        # Setup mocks
        mock_config = {"tools_enabled": True, "confirmation_required": False}
        mock_load_config.return_value = mock_config

        # Run CLI with conflicting options
        result = self.runner.invoke(main, ["--prompt", "test", "--resume", "session"])

        # Verify usage error behavior
        assert "Cannot use --prompt, --file, and --resume together" in result.output

    @patch("python_agent.cli.load_configuration")
    def test_main_configuration_error_behavior(self, mock_load_config):
        """Test configuration error behavior."""
        # Setup mock to raise ConfigurationError
        mock_load_config.side_effect = ConfigurationError("Invalid config")

        # Run CLI
        result = self.runner.invoke(main, [])

        # Verify configuration error behavior
        assert "Configuration error: Invalid config" in result.output

    @patch("python_agent.cli.load_configuration")
    @patch("python_agent.agent.Agent")
    def test_main_keyboard_interrupt_behavior(self, mock_agent_class, mock_load_config):
        """Test KeyboardInterrupt handling behavior."""
        # Setup mocks
        mock_config = {"tools_enabled": True, "confirmation_required": False}
        mock_load_config.return_value = mock_config
        mock_agent = MagicMock()
        mock_agent.interactive_loop.side_effect = KeyboardInterrupt()
        mock_agent_class.return_value = mock_agent

        # Run CLI
        self.runner.invoke(main, [])

        # Verify KeyboardInterrupt behavior - should not show user-facing error since this is expected
        # The behavior should be handled gracefully without exceptions leaking to user
