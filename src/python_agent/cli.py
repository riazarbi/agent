"""Main CLI entry point for the minimal AI coding agent.

Keywords: CLI, command-line, interface, click, main, entry-point

Implements:
- Interactive mode (default)
- Single-shot mode (--prompt)
- File input mode (--file)
- Session resume (--resume)
- Tool flags (--allow-tools/--no-tools)
- Confirmation flags (--confirm/--no-confirm)
"""

import sys
from pathlib import Path

import click

from python_agent.config import ConfigurationError, load_configuration


@click.command()
@click.option("--prompt", "-p", help="Single-shot mode: execute prompt and exit")
@click.option(
    "--file",
    "-f",
    "input_file",
    type=click.Path(exists=True, path_type=Path),
    help="File input mode: load prompt from file",
)
@click.option("--resume", "-r", "session_id", help="Resume session by ID")
@click.option(
    "--config", "-c", type=click.Path(path_type=Path), help="Configuration file path"
)
@click.option(
    "--allow-tools/--no-tools",
    default=None,
    help="Enable or disable tool execution (overrides config)",
)
@click.option(
    "--confirm/--no-confirm",
    default=None,
    help="Enable or disable confirmation prompts (overrides config)",
)
@click.option("--verbose", "-v", is_flag=True, help="Enable verbose output")
@click.option("--quiet", "-q", is_flag=True, help="Enable quiet mode (errors only)")
def main(
    prompt: str | None,
    input_file: Path | None,
    session_id: str | None,
    config: Path | None,
    allow_tools: bool | None,
    confirm: bool | None,
    verbose: bool,
    quiet: bool,
) -> int:
    """Minimal AI coding agent with single bash tool execution.

    Provides interactive chat mode by default, with options for single-shot
    prompts, file input, and session resumption.

    Examples:
        agent                           # Interactive mode
        agent --prompt "List files"     # Single-shot mode
        agent --file prompt.txt         # File input mode
        agent --resume 2024-01-01-12-00 # Resume session
    """
    try:
        # Load configuration
        agent_config = load_configuration(config)

        # Override config with CLI flags
        if allow_tools is not None:
            agent_config["tools_enabled"] = allow_tools
        if confirm is not None:
            agent_config["confirmation_required"] = confirm

        # Set verbosity
        agent_config["verbose"] = verbose
        agent_config["quiet"] = quiet

        # Validate mode exclusivity
        modes = [prompt, input_file, session_id]
        active_modes = [m for m in modes if m is not None]
        if len(active_modes) > 1:
            click.echo(
                "Error: Cannot use --prompt, --file, and --resume together", err=True
            )
            return 2

        # Import agent after successful configuration loading
        from python_agent.agent import Agent, AgentError

        # Initialize agent with configuration
        try:
            agent = Agent(agent_config)
        except Exception as e:
            click.echo(f"Failed to initialize agent: {e}", err=True)
            return 1

        try:
            # Handle different modes
            if prompt:
                # Single-shot mode
                if verbose:
                    click.echo("Running single-shot mode")
                response = agent.process_single_prompt(prompt)
                click.echo(f"Agent: {response}")
            elif input_file:
                # File input mode
                if verbose:
                    click.echo(f"Reading prompt from file: {input_file}")
                try:
                    file_content = input_file.read_text()
                except (FileNotFoundError, PermissionError, UnicodeDecodeError) as e:
                    click.echo(f"Error reading file {input_file}: {e}", err=True)
                    return 1
                response = agent.process_single_prompt(file_content)
                click.echo(f"Agent: {response}")
            elif session_id:
                # Session resume mode
                if verbose:
                    click.echo(f"Resuming session: {session_id}")
                from python_agent.session import SessionError, SessionManager

                session_manager = SessionManager()
                try:
                    session = session_manager.load_session(session_id)
                    agent.resume_from_session(session)
                    agent.interactive_loop()
                except SessionError as e:
                    click.echo(f"Session error: {e}", err=True)
                    return 1
            else:
                # Interactive mode (default)
                if verbose:
                    click.echo("Starting interactive mode")
                agent.interactive_loop()

        except AgentError as e:
            click.echo(f"Agent error: {e}", err=True)
            return 1
        except KeyboardInterrupt:
            if verbose:
                click.echo("\nInterrupted by user", err=True)
            return 130  # Standard exit code for Ctrl+C
        except Exception as e:
            click.echo(f"Unexpected error: {e}", err=True)
            if verbose:
                import traceback

                click.echo(traceback.format_exc(), err=True)
            return 1

        return 0

    except ConfigurationError as e:
        click.echo(f"Configuration error: {e}", err=True)
        return 1
    except Exception as e:
        if verbose:
            click.echo(f"Unexpected error: {e}", err=True)
        else:
            click.echo(f"Error: {e}", err=True)
        return 1


if __name__ == "__main__":
    sys.exit(main())
