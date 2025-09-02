"""YAML configuration loading and management.

This module provides configuration loading from YAML files with environment
variable overrides and basic validation for the minimal AI coding agent.

Keywords: configuration, YAML, config, settings, environment, variables
"""

import os
from pathlib import Path
from typing import Any

import yaml


class ConfigurationError(Exception):
    """Raised when configuration is invalid.

    Keywords: configuration, error, config, settings
    """

    pass


def get_default_config() -> dict[str, Any]:
    """Get default configuration values.

    Keywords: default, configuration, settings, values

    Returns:
        Default configuration dictionary
    """
    return {
        "model": "gpt-3.5-turbo",
        "api_key": None,
        "base_url": None,
        "timeout": 30,
        "max_tokens": 4000,
        "temperature": 0.7,
        "tools_enabled": True,
        "confirmation_required": False,
        "session_dir": str(Path.home() / ".agent" / "sessions"),
    }


def load_config_file(config_path: Path | None = None) -> dict[str, Any]:
    """Load configuration from YAML file.

    Keywords: load, configuration, YAML, file, settings

    Args:
        config_path: Path to YAML config file. If None, uses default location.

    Returns:
        Configuration dictionary from file

    Raises:
        ConfigurationError: If file cannot be loaded or parsed
    """
    if config_path is None:
        config_path = Path.home() / ".agent" / "config.yaml"

    if not config_path.exists():
        return {}

    try:
        with open(config_path) as f:
            data = yaml.safe_load(f)
            return data if data is not None else {}
    except yaml.YAMLError as e:
        raise ConfigurationError(f"Invalid YAML in config file: {e}") from e
    except OSError as e:
        raise ConfigurationError(f"Cannot read config file: {e}") from e


def apply_env_overrides(config: dict[str, Any]) -> None:
    """Apply environment variable overrides to configuration.

    Keywords: environment, variables, overrides, configuration, settings

    Args:
        config: Configuration dictionary to modify in-place
    """
    env_mappings = {
        "API_KEY": "api_key",
        "MODEL": "model",
        "BASE_URL": "base_url",
        "TIMEOUT": "timeout",
        "MAX_TOKENS": "max_tokens",
        "TEMPERATURE": "temperature",
        "TOOLS_ENABLED": "tools_enabled",
        "CONFIRMATION_REQUIRED": "confirmation_required",
        "SESSION_DIR": "session_dir",
    }

    for env_var, config_key in env_mappings.items():
        value = os.getenv(env_var)
        if value is not None:
            # Convert string values to appropriate types
            if config_key in ["timeout", "max_tokens"]:
                config[config_key] = int(value)
            elif config_key == "temperature":
                config[config_key] = float(value)
            elif config_key in ["tools_enabled", "confirmation_required"]:
                config[config_key] = value.lower() in ("true", "1", "yes")
            else:
                config[config_key] = value


def validate_config(config: dict[str, Any]) -> None:
    """Validate configuration values.

    Keywords: validation, configuration, settings, validate

    Args:
        config: Configuration dictionary to validate

    Raises:
        ConfigurationError: If configuration is invalid
    """
    if not isinstance(config.get("timeout"), int) or config["timeout"] <= 0:
        raise ConfigurationError("timeout must be positive integer")

    if not isinstance(config.get("max_tokens"), int) or config["max_tokens"] <= 0:
        raise ConfigurationError("max_tokens must be positive integer")

    temp = config.get("temperature")
    if not isinstance(temp, int | float) or temp < 0 or temp > 2:
        raise ConfigurationError("temperature must be number between 0 and 2")


def load_configuration(config_path: Path | None = None) -> dict[str, Any]:
    """Load complete configuration with defaults, file, and environment overrides.

    Keywords: load, configuration, complete, settings, environment

    Args:
        config_path: Optional path to config file

    Returns:
        Complete configuration dictionary

    Raises:
        ConfigurationError: If configuration is invalid
    """
    config = get_default_config()

    # Load from file
    file_config = load_config_file(config_path)
    config.update(file_config)

    # Apply environment overrides
    apply_env_overrides(config)

    # Validate final configuration
    validate_config(config)

    return config
