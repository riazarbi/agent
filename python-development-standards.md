# Python Development Standards

A comprehensive guide for Python package development using uv, optimized for both human readability and agent-friendly code patterns.

## Purpose & Principles

This document establishes mandatory standards for structuring Python packages using uv. These are not suggestions - they are requirements for all Python projects.

### Core Principles
- **Agent Discoverability**: Structure code for efficient grep/glob exploration
- **Token Efficiency**: Keep files focused and readable in single agent context
- **Search-Friendly Patterns**: Use consistent, searchable naming conventions
- **Focused Responsibility**: Each module has single, clear purpose
- **Real Environment Testing**: No mocking - test against real APIs and dependencies
- **Quality Gates**: Automated enforcement of standards and coverage requirements

## Mandatory Tooling

All Python projects **MUST** use these exact tools with the specified configurations:

- **`uv`** - Package management and environment (Python >= 3.12)
- **`ruff`** - Linting and import sorting
- **`black`** - Code formatting  
- **`mypy`** - Type checking (strict mode)
- **`pytest`** - Testing framework with coverage

For comprehensive uv documentation, see: https://docs.astral.sh/uv/

## Project Structure

**MUST** follow this exact directory structure:

```
project-name/
├── pyproject.toml          # REQUIRED: Project metadata and dependencies
├── uv.lock                 # REQUIRED: Must be committed to version control
├── .venv/                  # Auto-generated, MUST be in .gitignore
├── src/
│   └── project_name/       # REQUIRED: Package must be in src/ layout
│       ├── __init__.py     # REQUIRED
│       ├── models/         # Data structures
│       ├── services/       # Business logic
│       ├── repositories/   # Data access layer (if applicable)
│       ├── validators/     # Input validation
│       └── utils/          # Utilities (if absolutely needed)
└── tests/                  # REQUIRED: Test directory
    ├── __init__.py
    ├── unit/              # Unit tests
    ├── integration/       # Integration tests
    ├── e2e/               # End-to-end tests
    ├── api/               # API tests (real APIs, no mocking)
    ├── helpers/           # Test utilities and assertions
    └── fixtures/          # Test data and samples
```

**MUST NOT** use flat layout (package in project root).

## Environment & Project Setup

### Project Initialization
**MUST** initialize new projects with:
```bash
uv init --python 3.12 project-name
cd project-name
```

**MUST** specify Python version >= 3.12. **MUST NOT** use older versions.

### Dependency Configuration
**MUST** use dependency groups in `pyproject.toml`:

```toml
[project]
name = "project-name"
version = "0.1.0"
description = "Brief description"
authors = [{name = "Author Name", email = "email@domain.com"}]
license = {text = "MIT"}
requires-python = ">=3.12"
keywords = ["keyword1", "keyword2"]
classifiers = [
    "Development Status :: 4 - Beta",
    "Intended Audience :: Developers",
    "License :: OSI Approved :: MIT License",
    "Programming Language :: Python :: 3",
]

# Entry points depend on project type:
# For CLI tools: [project.scripts]
# For libraries: [project.entry-points] or no entry points
# For web applications: typically no entry points in pyproject.toml

[project.optional-dependencies]
dev = [
    "pytest>=7.0", 
    "pytest-cov>=4.0",
    "ruff>=0.1", 
    "mypy>=1.0", 
    "black>=23.0"
]
test = [
    "pytest>=7.0",
    "pytest-cov>=4.0", 
    "pytest-mock>=3.10",
    "pytest-asyncio>=0.21",
    "hypothesis>=6.68",
]
docs = ["mkdocs>=1.5", "mkdocs-material>=9.0"]
```

**MUST NOT** add dependencies without specifying minimum versions.

### Naming Conventions
**MUST** use these naming conventions:
- Package directory: `snake_case`
- Module files: `snake_case.py`
- Entry point modules: descriptive names based on function (e.g., `main.py`, `server.py`, `processor.py`)

### Environment Management
**MUST** run all commands through uv:
```bash
uv sync                      # Environment setup
uv run pytest              # NOT: pytest
uv run ruff check .         # NOT: ruff check .  
uv run mypy src/            # NOT: mypy src/
```

**MUST** use `uv sync` for environment setup. **MUST NOT** use pip directly.

**MUST** commit `uv.lock` to version control.

### Lock File Management
**MUST** regenerate lock file when adding dependencies:
```bash
uv add package-name
uv lock
```

**MUST NOT** manually edit `uv.lock`.

**MUST** use `uv sync --frozen` in CI/CD and production.

## Code Style Standards

### File Size Limits (Agent-Friendly)
- **200 lines**: Preferred maximum for optimal agent processing
- **300 lines**: Warning threshold  
- **400 lines**: Hard limit - split into multiple modules

### Code Layout
- **88 characters maximum** line length (organizational standard)
- **4 spaces** per level indentation (never tabs)
- **2 blank lines** around classes and top-level functions
- **1 blank line** around methods within classes

### Import Organization
Groups separated by blank lines, in this order:
1. **Standard library**
2. **Third-party packages**
3. **Local application imports**

```python
import os
import sys
from pathlib import Path

import click
import requests
from pydantic import BaseModel

from project_name.models import User
from project_name.services import UserService
```

### Agent-Friendly Naming Conventions

#### Modules and Packages
- Use descriptive names that indicate purpose: `user_service.py`, `data_processor.py`
- Avoid generic names: ~~`utils.py`~~, ~~`helpers.py`~~, ~~`common.py`~~
- Include domain keywords for searchability

#### Classes and Types
- **PascalCase** with searchable domain terms
- Include purpose in name: `UserValidator`, `DataProcessor`, `RequestHandler`
- For data classes: `UserData`, `ConfigurationData`

#### Functions and Variables
- **snake_case** with descriptive, searchable names
- Use verb phrases for functions: `validate_user_input`, `process_data`, `handle_request`
- Include domain keywords: `user_email`, `config_path`, `processing_result`

#### Constants
- **SCREAMING_SNAKE_CASE** with clear purpose
- Group related constants: `MAX_USER_COUNT`, `DEFAULT_USER_TIMEOUT`

### Agent-Searchable Patterns
```python
# Good: Includes searchable keywords
class UserAuthenticationService:
    """Handles user authentication and session management.
    
    Keywords: authentication, login, session, user, security
    """
    
def validate_user_credentials(username: str, password: str) -> bool:
    """Validate user credentials against authentication system.
    
    Keywords: validation, credentials, authentication, login
    """

# Bad: Generic, unsearchable names  
class Service:
    """Does stuff."""
    
def check(data):
    """Checks something."""
```

## Type Hints & Documentation

### Type Hints (Mandatory)
**MUST** include type hints for all public functions, methods, and class attributes:

```python
from typing import Optional, List, Dict, Union
from pathlib import Path

class UserService:
    """Service for user management operations.
    
    Keywords: user, service, management, CRUD
    """
    
    def __init__(self, database_url: str, timeout: int = 30) -> None:
        self.database_url = database_url
        self.timeout = timeout
    
    def create_user(
        self, 
        name: str, 
        email: str,
        metadata: Optional[Dict[str, str]] = None
    ) -> User:
        """Create new user with validation.
        
        Args:
            name: User's full name
            email: Valid email address  
            metadata: Optional user metadata
            
        Returns:
            Created user instance
            
        Raises:
            ValidationError: If input data is invalid
            DatabaseError: If creation fails
        """
```

### Documentation Standards

#### Module-Level Docstrings
Include searchable keywords and usage examples:

```python
"""User authentication and session management.

This module provides authentication services including login validation,
session management, and user credential handling.

Keywords: authentication, login, session, user, credentials, security

Basic usage:
    auth_service = AuthenticationService(database_url)
    user = auth_service.authenticate_user(username, password)
    
Main classes:
    - AuthenticationService: Primary authentication interface
    - SessionManager: Handles user sessions
    - CredentialValidator: Validates user credentials
"""
```

#### Function/Method Docstrings
Follow Google/NumPy style with searchable keywords:

```python
def process_user_data(
    user_data: Dict[str, str], 
    validation_rules: List[str]
) -> ProcessedUserData:
    """Process and validate user data according to rules.
    
    Applies validation rules to user data and returns processed result.
    Includes data sanitization and format normalization.
    
    Keywords: validation, processing, user data, sanitization
    
    Args:
        user_data: Raw user data dictionary
        validation_rules: List of validation rule names to apply
        
    Returns:
        ProcessedUserData instance with validated and sanitized data
        
    Raises:
        ValidationError: If data fails validation rules
        ProcessingError: If data processing fails
        
    Example:
        >>> rules = ['email_format', 'name_length']
        >>> data = {'email': 'user@example.com', 'name': 'John Doe'}
        >>> result = process_user_data(data, rules)
    """
```

## Error Handling

### Custom Exception Classes
```python
class ProjectNameError(Exception):
    """Base exception for project-name errors.
    
    Keywords: error, exception, base
    """

class ValidationError(ProjectNameError):
    """Raised when data validation fails.
    
    Keywords: validation, error, data validation
    """
    
    def __init__(self, field: str, message: str, value: str = "") -> None:
        self.field = field
        self.value = value
        super().__init__(f"Validation failed for {field}: {message}")

class ConfigurationError(ProjectNameError):
    """Raised when configuration is invalid.
    
    Keywords: configuration, error, config, settings
    """
```

### Error Handling Patterns
```python
def load_user_configuration(config_path: Path) -> UserConfiguration:
    """Load user configuration from file path.
    
    Keywords: configuration, loading, user config, file
    """
    try:
        with open(config_path, 'r') as f:
            data = json.load(f)
    except FileNotFoundError:
        raise ConfigurationError(
            f"Configuration file not found: {config_path}"
        )
    except json.JSONDecodeError as e:
        raise ConfigurationError(
            f"Invalid JSON in configuration file: {e}"
        ) from e
    
    try:
        return UserConfiguration.from_dict(data)
    except KeyError as e:
        raise ConfigurationError(
            f"Missing required configuration key: {e}"
        ) from e
```

## Testing Architecture

### Testing Philosophy
1. **No API Mocking** - Test against real APIs with proper skip logic
2. **Real Environment Testing** - Use actual files, networks, and dependencies
3. **Comprehensive Coverage** - Unit, integration, and end-to-end testing
4. **Test Independence** - Tests must be fully independent and run in any order
5. **Descriptive Testing** - Use long, descriptive test names as documentation

### Testing Pyramid
```
    🔺 E2E Tests (few, slow, high confidence)
   🔺🔺 Integration Tests (some, medium speed)
  🔺🔺🔺 Unit Tests (many, fast, focused)
```

### Test File Organization

#### Test Markers Usage
Use pytest markers to categorize tests:
```python
import pytest

@pytest.mark.unit
def test_fast_unit_test():
    pass

@pytest.mark.integration  
def test_integration_with_database():
    pass

@pytest.mark.api
def test_real_api_call():
    pass

@pytest.mark.e2e
def test_complete_user_workflow():
    pass

@pytest.mark.slow
def test_long_running_operation():
    pass
```

#### Running Tests by Category
```bash
# Unit tests only (default)
uv run pytest -m unit

# Integration tests
uv run pytest -m integration

# API tests (requires API keys)
uv run pytest -m api

# End-to-end tests
uv run pytest -m e2e

# All tests except slow ones
uv run pytest -m "not slow"
```

### API Testing Standards

#### No Mocking Rule
**NEVER mock API calls.** Instead, test against real APIs with proper skip logic.

#### API Test Template
```python
import pytest
import os

@pytest.mark.api
def test_real_api_call():
    """Test against real API with proper skip logic."""
    api_key = os.getenv("ANTHROPIC_API_KEY")
    if not api_key:
        pytest.skip("Skipping API test: no API key provided (set ANTHROPIC_API_KEY)")
    
    # Test with real API
    client = create_client(api_key)
    result = client.make_request(test_data)
    
    assert result is not None
    assert result.status == "success"
```

### Test Writing Standards

#### Unit Test Template
```python
import pytest
from hypothesis import given, strategies as st

class TestFunctionUnderTest:
    """Test suite for function_under_test."""
    
    @pytest.fixture
    def sample_data(self):
        """Provide sample data for tests."""
        return {"key": "value", "number": 42}
    
    def test_function_under_test_when_valid_input_then_returns_expected_result(self, sample_data):
        """Test function with valid input returns expected result."""
        result = function_under_test(sample_data)
        
        assert result is not None
        assert result["processed"] is True
        assert result["count"] == 1
    
    def test_function_under_test_when_invalid_input_then_raises_value_error(self):
        """Test function with invalid input raises ValueError."""
        with pytest.raises(ValueError, match="Invalid input provided"):
            function_under_test(None)
    
    @pytest.mark.parametrize("input_value,expected", [
        ("valid", True),
        ("", False),
        ("invalid", False),
    ])
    def test_function_under_test_various_inputs(self, input_value, expected):
        """Test function behavior with various input values."""
        result = function_under_test(input_value)
        assert result == expected
    
    @given(st.text())
    def test_function_under_test_property_based(self, text_input):
        """Property-based test using Hypothesis."""
        try:
            result = function_under_test(text_input)
            assert isinstance(result, (str, type(None)))
        except ValueError:
            # Expected for invalid inputs
            pass
```

#### Integration Test Template
```python
import pytest
import tempfile
import shutil
from pathlib import Path

@pytest.mark.integration
class TestComponentIntegration:
    """Integration tests for component interactions."""
    
    @pytest.fixture
    def temp_workspace(self):
        """Create temporary workspace for integration tests."""
        temp_dir = tempfile.mkdtemp()
        yield Path(temp_dir)
        shutil.rmtree(temp_dir, ignore_errors=True)
    
    def test_complete_workflow_integration(self, temp_workspace):
        """Test complete workflow from start to finish."""
        # Setup realistic test environment
        config_file = temp_workspace / "config.json"
        config_file.write_text('{"setting": "test_value"}')
        
        # Test the complete workflow
        result = run_complete_workflow(str(temp_workspace))
        
        # Verify all expected outcomes
        assert result.success is True
        assert (temp_workspace / "output.txt").exists()
        assert "processed" in (temp_workspace / "output.txt").read_text()
```

#### E2E Test Template
```python
import pytest
import subprocess
from pathlib import Path

@pytest.mark.e2e
class TestEndToEndWorkflows:
    """End-to-end tests for complete user workflows."""
    
    def test_complete_application_workflow(self, tmp_path):
        """Test complete application workflow from user perspective."""
        # Create realistic input files
        input_file = tmp_path / "input.txt"
        input_file.write_text("test input data")
        
        # Test the complete workflow (adjust based on your application type)
        # For CLI tools: subprocess.run with uv run
        # For libraries: direct function calls
        # For web apps: HTTP requests or test client
        result = run_application_workflow(
            input_path=str(input_file),
            output_path=str(tmp_path / "output.txt")
        )
        
        # Verify workflow succeeded
        assert result.success is True
        assert result.message == "Processing completed successfully"
        
        # Verify output file was created correctly
        output_file = tmp_path / "output.txt"
        assert output_file.exists()
        assert "processed: test input data" in output_file.read_text()
```

### Test Data Management

#### Required Helper Functions
```python
# tests/helpers/api_helpers.py
def get_test_api_key(env_var_name: str) -> str | None:
    """Get API key from environment with proper handling."""
    return os.getenv(env_var_name)

def skip_if_no_api_key(api_key: str | None, env_var_name: str) -> None:
    """Skip test if no API key is provided."""
    if not api_key:
        pytest.skip(f"Skipping API test: no API key provided (set {env_var_name})")

# tests/helpers/assertions.py
def assert_valid_json(data: str) -> None:
    """Assert that string is valid JSON."""
    try:
        json.loads(data)
    except json.JSONDecodeError:
        pytest.fail(f"Invalid JSON: {data}")

def assert_file_exists(path: str) -> None:
    """Assert that file exists."""
    assert os.path.exists(path), f"File does not exist: {path}"
```

## Tool Configuration

### Ruff Configuration (Mandatory)
```toml
[tool.ruff]
line-length = 88
target-version = "py312"

[tool.ruff.lint]
select = [
    "E",    # pycodestyle errors
    "W",    # pycodestyle warnings  
    "F",    # pyflakes
    "I",    # isort
    "N",    # pep8-naming
    "B",    # flake8-bugbear
    "UP",   # pyupgrade
    "C4",   # flake8-comprehensions
    "SIM",  # flake8-simplify
]
ignore = ["E501"]  # Line too long (handled by black)

[tool.ruff.lint.per-file-ignores]
"tests/*" = ["S101"]  # Allow assert in tests
```

### Black Configuration
```toml
[tool.black]
line-length = 88
target-version = ['py312']
include = '\.pyi?$'
```

### MyPy Configuration (Mandatory)
```toml
[tool.mypy]
strict = true
warn_unreachable = true
warn_unused_ignores = true
python_version = "3.12"
```

### Pytest Configuration
```toml
[tool.pytest.ini_options]
minversion = "7.0"
testpaths = ["tests"]
python_files = ["test_*.py"]
python_classes = ["Test*"]
python_functions = ["test_*"]
markers = [
    "unit: Fast unit tests",
    "integration: Integration tests requiring setup", 
    "api: Tests requiring API keys",
    "e2e: End-to-end workflow tests",
    "slow: Long-running tests",
]
filterwarnings = [
    "error",
    "ignore::UserWarning",
    "ignore::DeprecationWarning",
]
```

### Coverage Configuration
```toml
[tool.coverage.run]
source = ["src"]
omit = [
    "tests/*",
    "*/__pycache__/*",
    "*/migrations/*",
]

[tool.coverage.report]
exclude_lines = [
    "pragma: no cover",
    "def __repr__",
    "raise AssertionError",
    "raise NotImplementedError",
    "if __name__ == .__main__.:",
]
```

## Development Workflow

### Pre-commit Workflow (Required)
**MUST** run this sequence before every commit:
```bash
uv sync
uv run ruff check .
uv run ruff format .
uv run mypy src/
uv run pytest --cov=src --cov-fail-under=80
```

**MUST NOT** commit code that fails any of these checks.

### Test Commands
All tests must run through uv:
```bash
# Install test dependencies
uv sync --group test

# Run all tests
uv run pytest

# Run tests with coverage
uv run pytest --cov=src --cov-report=html --cov-report=term

# Run specific test categories
uv run pytest -m unit
uv run pytest -m "integration and not slow"

# Run tests in parallel (if pytest-xdist installed)
uv run pytest -n auto

# Run tests with verbose output
uv run pytest -v --tb=short
```

### Environment Variables for Testing
- `ANTHROPIC_API_KEY` or `OPENAI_API_KEY` - API keys for testing
- `SKIP_API_TESTS=true` - Skip all API tests
- `SKIP_SLOW_TESTS=true` - Skip long-running tests
- `TEST_TIMEOUT=30` - Override test timeouts (seconds)

## Advanced Features

### Script Dependencies (Standalone Scripts)
For standalone scripts, **MUST** use inline metadata:

```python
#!/usr/bin/env python3
# /// script
# dependencies = [
#     "click>=8.0",
#     "requests>=2.28",
# ]
# ///
```

**MUST** pin major versions for inline dependencies.

### Tool Installation
**MUST** use global tool installation:
```bash
uv tool install package-name
```

**MUST** use `uvx` for one-time tool execution:
```bash
uvx package-name --help
```

**MUST NOT** install development tools into project virtual environments.

### Docker Deployment
**MUST** use this Dockerfile template (adjust CMD based on application type):

```dockerfile
FROM python:3.12-slim
COPY --from=ghcr.io/astral-sh/uv:latest /uv /bin/uv
ENV UV_LINK_MODE=copy
WORKDIR /app
COPY . .
RUN uv sync --frozen
# Adjust CMD based on your application:
# For CLI tools: CMD ["uv", "run", "project-name"]
# For web apps: CMD ["uv", "run", "uvicorn", "project_name.main:app", "--host", "0.0.0.0"]
# For scripts: CMD ["uv", "run", "python", "-m", "project_name"]
CMD ["uv", "run", "python", "-m", "project_name"]
```

## Quality Gates & Enforcement

### Required Coverage
- **Minimum**: 80% line coverage (enforced)
- **Target**: 90% line coverage
- **Critical Paths**: 100% coverage for public APIs and core functionality

### CI/CD Requirements
1. **Unit Tests**: Must pass on all commits
2. **Integration Tests**: Must pass on PRs  
3. **API Tests**: Run with provided API keys (skip if not available)
4. **Coverage Check**: Minimum 80% coverage enforced
5. **Type Checking**: mypy must pass with strict settings
6. **Linting**: ruff must pass with no violations

### Automated Enforcement
These standards are enforced through:
- **Pre-commit hooks** with ruff, black, mypy
- **CI/CD pipeline checks** blocking non-compliant code  
- **Automated code review** tools checking patterns
- **File size limits** preventing overly complex modules
- **Test structure validation** ensuring proper directory layout

## Decision Framework

When writing code, ask:

1. **Can an agent find this functionality via grep without reading unrelated code?**
   → Use searchable keywords in names and comments

2. **Is this file focused enough for an agent to understand completely?**
   → Keep files under 200 lines, split by responsibility

3. **Are the domain concepts clearly expressed in the code structure?**
   → Use descriptive names that match business concepts

4. **Can an agent understand dependencies and relationships quickly?**
   → Clear imports, explicit typing, good documentation

5. **Can this be tested against real systems?**
   → Prefer real API calls with skip logic over mocking

## Violations & Consequences

Breaking any **MUST** rule will result in:
- Immediate build failure
- Code review rejection
- CI/CD pipeline failure

Breaking any **MUST NOT** rule will result in:
- Automatic linting errors
- Test suite failures
- Deployment blocking

**Key Question**: Can an AI agent understand and work with this code efficiently without reading unrelated files?