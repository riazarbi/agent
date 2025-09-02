"""Single bash tool implementation for command execution.

Keywords: bash, tool, subprocess, command, execution, shell

This module provides the single bash tool for executing shell commands
with confirmation prompts and timeout handling as required.
"""

import subprocess
from typing import Any


class BashToolError(Exception):
    """Base exception for bash tool errors.

    Keywords: error, exception, bash, tool
    """


class BashTool:
    """Bash command execution tool with confirmation and timeout support.

    Keywords: bash, command, execution, subprocess, shell

    Provides controlled execution of bash commands with optional confirmation
    prompts and configurable timeout handling.
    """

    def __init__(self, confirmation_required: bool = False, timeout: int = 30) -> None:
        """Initialize bash tool with configuration options.

        Args:
            confirmation_required: Whether to prompt for confirmation before execution
            timeout: Command timeout in seconds (default: 30)
        """
        self.confirmation_required = confirmation_required
        self.timeout = timeout
        self.enabled = True

    def set_enabled(self, enabled: bool) -> None:
        """Enable or disable the bash tool.

        Keywords: enable, disable, tool, control

        Args:
            enabled: Whether the tool should be enabled
        """
        self.enabled = enabled

    def _confirm_execution(self, command: str) -> bool:
        """Prompt user to confirm command execution.

        Keywords: confirmation, prompt, user, safety

        Args:
            command: The command to be executed

        Returns:
            True if user confirms, False otherwise
        """
        print(f"Execute command: {command}")
        response = input("Continue? [y/N]: ").strip().lower()
        return response in ("y", "yes")

    def execute_command(self, command: str) -> dict[str, Any]:
        """Execute a bash command with timeout and error handling.

        Keywords: execute, bash, command, subprocess, timeout

        Args:
            command: The bash command to execute

        Returns:
            Dictionary containing execution results with keys:
            - success: bool indicating if command succeeded
            - output: str containing stdout output
            - error: str containing stderr output
            - exit_code: int command exit code

        Raises:
            BashToolError: If tool is disabled or other execution errors occur
        """
        if not self.enabled:
            raise BashToolError("Bash tool is disabled")

        if self.confirmation_required and not self._confirm_execution(command):
            return {
                "success": False,
                "output": "",
                "error": "Command execution cancelled by user",
                "exit_code": 1,
            }

        try:
            result = subprocess.run(
                command,
                shell=True,
                capture_output=True,
                text=True,
                timeout=self.timeout,
            )

            return {
                "success": result.returncode == 0,
                "output": result.stdout,
                "error": result.stderr,
                "exit_code": result.returncode,
            }

        except subprocess.TimeoutExpired:
            return {
                "success": False,
                "output": "",
                "error": f"Command timed out after {self.timeout} seconds",
                "exit_code": 124,
            }
        except Exception as e:
            return {
                "success": False,
                "output": "",
                "error": f"Execution error: {str(e)}",
                "exit_code": 1,
            }
