"""Unit tests for bash tool module.

Keywords: test, bash, tool, subprocess, command, execution
"""

import subprocess
from unittest.mock import Mock, patch

import pytest

from python_agent.bash_tool import BashTool, BashToolError


class TestBashToolInit:
    """Test suite for BashTool initialization."""

    def test_default_initialization(self):
        """Test BashTool with default parameters."""
        tool = BashTool()

        assert tool.confirmation_required is False
        assert tool.timeout == 30
        assert tool.enabled is True

    def test_custom_initialization(self):
        """Test BashTool with custom parameters."""
        tool = BashTool(confirmation_required=True, timeout=60)

        assert tool.confirmation_required is True
        assert tool.timeout == 60
        assert tool.enabled is True


class TestBashToolSetEnabled:
    """Test suite for BashTool.set_enabled method."""

    def test_enable_tool(self):
        """Test enabling the bash tool."""
        tool = BashTool()
        tool.enabled = False
        
        tool.set_enabled(True)
        
        assert tool.enabled is True

    def test_disable_tool(self):
        """Test disabling the bash tool."""
        tool = BashTool()
        
        tool.set_enabled(False)
        
        assert tool.enabled is False


class TestBashToolConfirmExecution:
    """Test suite for BashTool._confirm_execution method."""

    def test_confirm_execution_with_y_input(self):
        """Test confirmation with 'y' input returns True."""
        tool = BashTool()
        
        with patch("builtins.input", return_value="y"):
            result = tool._confirm_execution("echo test")
            
        assert result is True

    def test_confirm_execution_with_yes_input(self):
        """Test confirmation with 'yes' input returns True."""
        tool = BashTool()
        
        with patch("builtins.input", return_value="yes"):
            result = tool._confirm_execution("echo test")
            
        assert result is True

    def test_confirm_execution_with_n_input(self):
        """Test confirmation with 'n' input returns False."""
        tool = BashTool()
        
        with patch("builtins.input", return_value="n"):
            result = tool._confirm_execution("echo test")
            
        assert result is False

    def test_confirm_execution_with_no_input(self):
        """Test confirmation with 'no' input returns False."""
        tool = BashTool()
        
        with patch("builtins.input", return_value="no"):
            result = tool._confirm_execution("echo test")
            
        assert result is False

    def test_confirm_execution_with_empty_input(self):
        """Test confirmation with empty input returns False (default)."""
        tool = BashTool()
        
        with patch("builtins.input", return_value=""):
            result = tool._confirm_execution("echo test")
            
        assert result is False

    def test_confirm_execution_with_arbitrary_input(self):
        """Test confirmation with arbitrary input returns False."""
        tool = BashTool()
        
        with patch("builtins.input", return_value="maybe"):
            result = tool._confirm_execution("echo test")
            
        assert result is False

    def test_confirm_execution_case_insensitive(self):
        """Test confirmation is case insensitive."""
        tool = BashTool()
        
        with patch("builtins.input", return_value="Y"):
            result = tool._confirm_execution("echo test")
            
        assert result is True

    def test_confirm_execution_strips_whitespace(self):
        """Test confirmation strips leading/trailing whitespace."""
        tool = BashTool()
        
        with patch("builtins.input", return_value="  yes  "):
            result = tool._confirm_execution("echo test")
            
        assert result is True


class TestBashToolExecuteCommand:
    """Test suite for BashTool.execute_command method."""

    def test_execute_command_when_disabled_raises_error(self):
        """Test executing command when tool is disabled raises error."""
        tool = BashTool()
        tool.set_enabled(False)
        
        with pytest.raises(BashToolError, match="Bash tool is disabled"):
            tool.execute_command("echo test")

    def test_execute_successful_command(self):
        """Test executing successful command returns expected result."""
        tool = BashTool()
        mock_result = Mock()
        mock_result.returncode = 0
        mock_result.stdout = "test output"
        mock_result.stderr = ""
        
        with patch("subprocess.run", return_value=mock_result):
            result = tool.execute_command("echo test")
            
        assert result["success"] is True
        assert result["output"] == "test output"
        assert result["error"] == ""
        assert result["exit_code"] == 0

    def test_execute_failed_command(self):
        """Test executing failed command returns expected result."""
        tool = BashTool()
        mock_result = Mock()
        mock_result.returncode = 1
        mock_result.stdout = ""
        mock_result.stderr = "command not found"
        
        with patch("subprocess.run", return_value=mock_result):
            result = tool.execute_command("invalid_command")
            
        assert result["success"] is False
        assert result["output"] == ""
        assert result["error"] == "command not found"
        assert result["exit_code"] == 1

    def test_execute_command_with_timeout(self):
        """Test executing command that times out."""
        tool = BashTool(timeout=5)
        
        with patch("subprocess.run", side_effect=subprocess.TimeoutExpired("sleep 10", 5)):
            result = tool.execute_command("sleep 10")
            
        assert result["success"] is False
        assert result["output"] == ""
        assert result["error"] == "Command timed out after 5 seconds"
        assert result["exit_code"] == 124

    def test_execute_command_with_subprocess_exception(self):
        """Test executing command that raises subprocess exception."""
        tool = BashTool()
        
        with patch("subprocess.run", side_effect=OSError("Process failed")):
            result = tool.execute_command("echo test")
            
        assert result["success"] is False
        assert result["output"] == ""
        assert result["error"] == "Execution error: Process failed"
        assert result["exit_code"] == 1

    def test_execute_command_calls_subprocess_with_correct_parameters(self):
        """Test that subprocess.run is called with correct parameters."""
        tool = BashTool(timeout=60)
        mock_result = Mock()
        mock_result.returncode = 0
        mock_result.stdout = "output"
        mock_result.stderr = ""
        
        with patch("subprocess.run", return_value=mock_result) as mock_run:
            tool.execute_command("echo test")
            
            mock_run.assert_called_once_with(
                "echo test",
                shell=True,
                capture_output=True,
                text=True,
                timeout=60
            )

    def test_execute_command_with_confirmation_accepted(self):
        """Test executing command with confirmation required and accepted."""
        tool = BashTool(confirmation_required=True)
        mock_result = Mock()
        mock_result.returncode = 0
        mock_result.stdout = "test output"
        mock_result.stderr = ""
        
        with patch("builtins.input", return_value="y"), \
             patch("subprocess.run", return_value=mock_result):
            result = tool.execute_command("echo test")
            
        assert result["success"] is True
        assert result["output"] == "test output"
        assert result["error"] == ""
        assert result["exit_code"] == 0

    def test_execute_command_with_confirmation_rejected(self):
        """Test executing command with confirmation required and rejected."""
        tool = BashTool(confirmation_required=True)
        
        with patch("builtins.input", return_value="n"):
            result = tool.execute_command("echo test")
            
        assert result["success"] is False
        assert result["output"] == ""
        assert result["error"] == "Command execution cancelled by user"
        assert result["exit_code"] == 1

    @pytest.mark.parametrize("command,timeout", [
        ("echo hello", 30),
        ("ls -la", 10),
        ("pwd", 5),
    ])
    def test_execute_various_commands(self, command, timeout):
        """Test executing various commands with different timeouts."""
        tool = BashTool(timeout=timeout)
        mock_result = Mock()
        mock_result.returncode = 0
        mock_result.stdout = f"output for {command}"
        mock_result.stderr = ""
        
        with patch("subprocess.run", return_value=mock_result):
            result = tool.execute_command(command)
            
        assert result["success"] is True
        assert result["output"] == f"output for {command}"


class TestBashToolErrorHandling:
    """Test suite for BashTool error handling scenarios."""

    def test_bash_tool_error_inheritance(self):
        """Test that BashToolError inherits from Exception."""
        error = BashToolError("test error")
        assert isinstance(error, Exception)
        assert str(error) == "test error"

    def test_bash_tool_error_with_empty_message(self):
        """Test BashToolError with empty message."""
        error = BashToolError("")
        assert str(error) == ""

    def test_execute_command_handles_general_exception(self):
        """Test that general exceptions are handled properly."""
        tool = BashTool()
        
        with patch("subprocess.run", side_effect=RuntimeError("Runtime error")):
            result = tool.execute_command("echo test")
            
        assert result["success"] is False
        assert result["output"] == ""
        assert result["error"] == "Execution error: Runtime error"
        assert result["exit_code"] == 1


class TestBashToolIntegration:
    """Integration tests for BashTool with real subprocess calls."""

    @pytest.mark.integration
    def test_execute_real_echo_command(self):
        """Test executing real echo command (integration test)."""
        tool = BashTool()
        
        result = tool.execute_command("echo 'integration test'")
        
        assert result["success"] is True
        assert "integration test" in result["output"]
        assert result["error"] == ""
        assert result["exit_code"] == 0

    @pytest.mark.integration  
    def test_execute_real_failed_command(self):
        """Test executing real failed command (integration test)."""
        tool = BashTool()
        
        result = tool.execute_command("command_that_does_not_exist_12345")
        
        assert result["success"] is False
        assert result["output"] == ""
        assert "not found" in result["error"] or "not recognized" in result["error"]
        assert result["exit_code"] != 0

    @pytest.mark.integration
    def test_execute_real_command_with_timeout(self):
        """Test executing real command with very short timeout (integration test)."""
        tool = BashTool(timeout=1)
        
        # Use a command that should timeout on most systems
        result = tool.execute_command("sleep 2")
        
        assert result["success"] is False
        assert result["output"] == ""
        assert "timed out after 1 seconds" in result["error"]
        assert result["exit_code"] == 124