# Product Specification: Minimal AI Coding Agent

## Overview

An ultra-lightweight Python command-line AI coding agent under 400 lines of code total that uses LiteLLM for model integration and provides a single bash tool for system interaction. Focused on simplicity and core functionality only.

## Core Requirements

### 1. Agent Capabilities
- **Interactive Chat Mode**: Continuous conversation with the AI agent
- **Single-Shot Mode**: Execute single prompts and exit
- **File-based Prompts**: Load prompts from files for automation
- **Session Management**: Save and resume conversation sessions
- **Bash Tool**: Single tool for executing shell commands

### 2. Architecture Components

#### Language Model Integration
- **LiteLLM Backend**: Support for all LiteLLM-compatible providers (OpenAI, Anthropic, local models, etc.)
- Direct API key passthrough to LiteLLM
- Basic error handling

#### Tool System
- **Single Bash Tool**: Execute shell commands only
- No file operations beyond what bash provides
- No specialized development tool integrations

#### Session Management
- Simple file-based conversation history
- Resume sessions by ID
- No metadata tracking beyond basic functionality

#### Configuration System
- **YAML configuration files**
- Environment variable overrides for API keys
- Minimal configuration options

### 3. User Experience

#### Command-Line Interface
- **Interactive Mode**: `agent` (default)
- **Single-Shot**: `agent --prompt "task description"`
- **File Input**: `agent --file path/to/prompt.txt`
- **Session Resume**: `agent --resume SESSION_ID`

#### Output
- Plain text output only
- No colors, no progress indicators, no formatting
- Simple role prefixes (User:, Agent:, Tool:)

### 4. Safety and Control

#### Execution Control
- **Tools Flag**: `--allow-tools` / `--no-tools` (default: allow)
- **Confirmation Flag**: `--confirm` / `--no-confirm` (default: no confirmation)
- No sandbox mode
- No complex safety controls

#### Data Handling
- API keys passed directly to LiteLLM
- No credential exposure protection beyond basic handling
- No session encryption
- No audit logging

## Technical Requirements

### Implementation Language
- **Python 3.10+** using `uv` for package management
- Following Python Development Standards from project documentation
- **Target: Under 400 lines of code total**

### Dependencies (Minimal)
- **Click**: CLI framework
- **LiteLLM**: Model integration
- **PyYAML**: Configuration file parsing
- **Pydantic** (optional): Basic data validation if needed

### Performance Targets
- **Startup Time**: < 200ms
- **Memory Usage**: < 50MB
- **Minimal overhead**: Focus on core functionality only

## Success Criteria

### Complete Feature Set (MVP = Final Product)
1. Interactive chat with LiteLLM-supported models
2. Single bash tool execution
3. Basic session save/resume
4. YAML configuration support
5. File-based prompt input
6. Simple CLI flags for tool/confirmation control

## Non-Requirements

### Explicitly Excluded Features
- File operations beyond bash
- Git integration
- Web requests (beyond bash curl)
- Development tool integration
- Multiple tools
- Plugin/extension system
- Advanced session management
- Colored output
- Progress indicators
- Rich formatting
- Sandbox mode
- Complex safety controls
- Session encryption
- Audit logging
- Tool whitelist/blacklist
- Cost tracking
- Usage limits
- Template rendering
- Complex configuration
- Export/import capabilities

## Design Principles

### Ultra-Minimalism
- Maximum 400 lines of code total
- Single tool only (bash)
- No unnecessary features or complexity
- Plain text output only

### Agent-Friendly Architecture
- Searchable code patterns and naming
- Token-efficient file organization (< 100 lines per module)
- Comprehensive docstrings with keywords
- Simple, readable code structure

### Reliability Through Simplicity
- Basic error handling only
- Fail fast on errors
- No complex state management
- File-based persistence only

## Target Users

### Primary Users
- **Developers**: Who want a simple AI coding assistant
- **System Administrators**: For basic automation tasks
- **Students**: Learning with minimal complexity

## Competitive Analysis

### Differentiators
- **Ultra-lightweight**: Under 400 lines total
- **Single purpose**: Only bash tool execution
- **LiteLLM integration**: Support for all major providers
- **Minimal setup**: YAML config and go