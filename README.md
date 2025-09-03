# Minimal AI Coding Agent

An ultra-lightweight Python command-line AI coding agent under 1000 lines of code that uses LiteLLM for model integration and provides a single bash tool for system interaction.

## Features

- **Interactive Chat Mode**: Continuous conversation with AI agents
- **Single-Shot Mode**: Execute single prompts and exit  
- **File-based Prompts**: Load prompts from files for automation
- **Session Management**: Save and resume conversation sessions
- **Bash Tool**: Execute shell commands with optional confirmation
- **Multi-Provider Support**: Works with OpenAI, Anthropic, Google Gemini, and other LiteLLM-compatible providers

## Installation

Install using uv (recommended):

```bash
uv tool install .
```

Alternative installation with pipx:

```bash
pipx install .
```

For development:

```bash
uv sync
```

## Quick Start

### Basic Usage

```bash
# Interactive mode (default)
agent

# Single-shot mode
agent --prompt "Explain Python list comprehensions"

# Load prompt from file
agent --file my-prompt.txt

# Resume previous session
agent --resume 2025-09-03-10-30-15
```

### Configuration

Create a config file at `~/.agent/config.yaml`:

```yaml
model: gemini/gemini-2.0-flash-exp
api_key: null  # Use environment variables
timeout: 30
max_tokens: 4000
temperature: 0.7
tools_enabled: true
confirmation_required: false
```

Set API keys via environment variables:

```bash
export GEMINI_API_KEY="your-gemini-key"
export OPENAI_API_KEY="your-openai-key"
export ANTHROPIC_API_KEY="your-anthropic-key"
```

### Command-Line Options

```bash
# Global options
agent --config custom-config.yaml --verbose
agent --quiet --no-tools

# Mode selection (mutually exclusive)
agent --prompt "Your question here"
agent --file input.txt
agent --resume SESSION_ID

# Tool control
agent --allow-tools --confirm    # Enable tools with confirmation
agent --no-tools                 # Disable all tools

# Output control
agent --verbose                  # Show detailed output
agent --quiet                    # Minimal output
```

## Usage Examples

### Interactive Mode

```bash
$ agent
AI Coding Agent (type 'exit' to quit)
Session: 2025-09-03-15-30-45

User: How do I list files in Python?
Agent: You can list files in Python using several methods:

yadda yadda...

User: exit
Session saved: 2025-09-03-15-30-45
```

### Single-Shot Mode

```bash
$ agent --prompt "Write a Python function to calculate factorial"
Agent: Here's a Python function to calculate factorial:

yadda yadda...

```

### File Input Mode

```bash
$ echo "Explain Python decorators with an example" > question.txt
$ agent --file question.txt
Agent: Python decorators are functions that modify or extend the behavior of other functions...
```

### Session Management

```bash
# Resume previous session
$ agent --resume 2025-09-03-15-30-45
AI Coding Agent (type 'exit' to quit)  
Session: 2025-09-03-15-30-45 (resumed)

User: Continue our previous discussion about file listing
Agent: Sure! Let me continue from where we left off about listing files in Python...
```

### Using with Different Models

```bash
# OpenAI GPT
export OPENAI_API_KEY="your-key"
agent --prompt "Hello" # Uses gpt-3.5-turbo (default)

# Google Gemini  
export GEMINI_API_KEY="your-key"
echo "model: gemini/gemini-2.0-flash-exp" > gemini-config.yaml
agent --config gemini-config.yaml --prompt "Hello"

# Anthropic Claude
export ANTHROPIC_API_KEY="your-key"  
echo "model: claude-3-haiku-20240307" > claude-config.yaml
agent --config claude-config.yaml --prompt "Hello"
```

## Project Structure

```
agent/
├── src/python_agent/
│   ├── __init__.py          # Package initialization (10 lines)
│   ├── cli.py              # Main CLI entry point (168 lines)  
│   ├── agent.py            # Core agent logic (138 lines)
│   ├── session.py          # Session management (212 lines)
│   ├── config.py           # Configuration handling (158 lines)
│   └── bash_tool.py        # Bash tool implementation (123 lines)
└── tests/                  # Comprehensive test suite (97% coverage)
```

**Total: 809 lines of code** (well under 1000-line limit)

## Technical Requirements

- Python 3.10+
- Dependencies: Click, LiteLLM, PyYAML only
- Startup time: <100ms (measured: ~60-80ms)
- Memory usage: <50MB
- Test coverage: 97%

## Development

```bash
# Setup development environment
uv sync

# Run tests
uv run pytest --cov=src --cov-report=html

# Run linting
uv run ruff check .
uv run ruff format .

# Type checking  
uv run mypy src/

# Run all quality checks
uv run ruff check . && uv run ruff format . && uv run mypy src/ && uv run pytest --cov=src --cov-fail-under=80
```

## License

MIT License - see LICENSE file for details.