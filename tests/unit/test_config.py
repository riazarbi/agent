"""Unit tests for configuration module.

Keywords: test, config, configuration, YAML, environment
"""

import os
import tempfile
from pathlib import Path
from unittest.mock import patch

import pytest
import yaml

from python_agent.config import (
    ConfigurationError,
    apply_env_overrides,
    get_default_config,
    load_config_file,
    load_configuration,
    validate_config,
)


class TestGetDefaultConfig:
    """Test suite for get_default_config function."""

    def test_returns_expected_default_values(self):
        """Test that default config contains expected values."""
        config = get_default_config()

        assert config["model"] == "gpt-3.5-turbo"
        assert config["api_key"] is None
        assert config["base_url"] is None
        assert config["timeout"] == 30
        assert config["max_tokens"] == 4000
        assert config["temperature"] == 0.7
        assert config["tools_enabled"] is True
        assert config["confirmation_required"] is False
        assert "sessions" in config["session_dir"]


class TestLoadConfigFile:
    """Test suite for load_config_file function."""

    def test_returns_empty_dict_when_file_not_exists(self):
        """Test loading non-existent config file returns empty dict."""
        non_existent = Path("/tmp/does_not_exist.yaml")
        result = load_config_file(non_existent)
        assert result == {}

    def test_loads_valid_yaml_file(self):
        """Test loading valid YAML config file."""
        with tempfile.NamedTemporaryFile(mode="w", suffix=".yaml", delete=False) as f:
            yaml.dump({"model": "custom-model", "timeout": 60}, f)
            temp_path = Path(f.name)

        try:
            result = load_config_file(temp_path)
            assert result["model"] == "custom-model"
            assert result["timeout"] == 60
        finally:
            temp_path.unlink()

    def test_raises_configuration_error_on_invalid_yaml(self):
        """Test that invalid YAML raises ConfigurationError."""
        with tempfile.NamedTemporaryFile(mode="w", suffix=".yaml", delete=False) as f:
            f.write("invalid: yaml: content:")
            temp_path = Path(f.name)

        try:
            with pytest.raises(ConfigurationError, match="Invalid YAML"):
                load_config_file(temp_path)
        finally:
            temp_path.unlink()

    def test_returns_empty_dict_for_empty_yaml_file(self):
        """Test that empty YAML file returns empty dict."""
        with tempfile.NamedTemporaryFile(mode="w", suffix=".yaml", delete=False) as f:
            f.write("")
            temp_path = Path(f.name)

        try:
            result = load_config_file(temp_path)
            assert result == {}
        finally:
            temp_path.unlink()


class TestApplyEnvOverrides:
    """Test suite for apply_env_overrides function."""

    def test_applies_string_environment_variables(self):
        """Test applying string environment variable overrides."""
        config = {"api_key": None, "model": "default"}

        with patch.dict(os.environ, {"API_KEY": "test-key", "MODEL": "gpt-4"}):
            apply_env_overrides(config)

        assert config["api_key"] == "test-key"
        assert config["model"] == "gpt-4"

    def test_applies_integer_environment_variables(self):
        """Test applying integer environment variable overrides."""
        config = {"timeout": 30, "max_tokens": 4000}

        with patch.dict(os.environ, {"TIMEOUT": "60", "MAX_TOKENS": "8000"}):
            apply_env_overrides(config)

        assert config["timeout"] == 60
        assert config["max_tokens"] == 8000

    def test_applies_float_environment_variables(self):
        """Test applying float environment variable overrides."""
        config = {"temperature": 0.7}

        with patch.dict(os.environ, {"TEMPERATURE": "0.5"}):
            apply_env_overrides(config)

        assert config["temperature"] == 0.5

    def test_applies_boolean_environment_variables(self):
        """Test applying boolean environment variable overrides."""
        config = {"tools_enabled": True, "confirmation_required": False}

        with patch.dict(
            os.environ, {"TOOLS_ENABLED": "false", "CONFIRMATION_REQUIRED": "true"}
        ):
            apply_env_overrides(config)

        assert config["tools_enabled"] is False
        assert config["confirmation_required"] is True

    @pytest.mark.parametrize(
        "env_value,expected",
        [
            ("true", True),
            ("1", True),
            ("yes", True),
            ("false", False),
            ("0", False),
            ("no", False),
            ("invalid", False),
        ],
    )
    def test_boolean_conversion_various_values(self, env_value, expected):
        """Test boolean conversion with various environment values."""
        config = {"tools_enabled": True}

        with patch.dict(os.environ, {"TOOLS_ENABLED": env_value}):
            apply_env_overrides(config)

        assert config["tools_enabled"] is expected


class TestValidateConfig:
    """Test suite for validate_config function."""

    def test_validates_valid_configuration(self):
        """Test that valid configuration passes validation."""
        config = {"timeout": 30, "max_tokens": 4000, "temperature": 0.7}
        # Should not raise any exception
        validate_config(config)

    def test_raises_error_for_invalid_timeout(self):
        """Test validation error for invalid timeout values."""
        with pytest.raises(
            ConfigurationError, match="timeout must be positive integer"
        ):
            validate_config({"timeout": 0, "max_tokens": 4000, "temperature": 0.7})

        with pytest.raises(
            ConfigurationError, match="timeout must be positive integer"
        ):
            validate_config(
                {"timeout": "invalid", "max_tokens": 4000, "temperature": 0.7}
            )

    def test_raises_error_for_invalid_max_tokens(self):
        """Test validation error for invalid max_tokens values."""
        with pytest.raises(
            ConfigurationError, match="max_tokens must be positive integer"
        ):
            validate_config({"timeout": 30, "max_tokens": 0, "temperature": 0.7})

        with pytest.raises(
            ConfigurationError, match="max_tokens must be positive integer"
        ):
            validate_config({"timeout": 30, "max_tokens": -1, "temperature": 0.7})

    def test_raises_error_for_invalid_temperature(self):
        """Test validation error for invalid temperature values."""
        with pytest.raises(ConfigurationError, match="temperature must be number"):
            validate_config({"timeout": 30, "max_tokens": 4000, "temperature": -0.1})

        with pytest.raises(ConfigurationError, match="temperature must be number"):
            validate_config({"timeout": 30, "max_tokens": 4000, "temperature": 2.1})

        with pytest.raises(ConfigurationError, match="temperature must be number"):
            validate_config(
                {"timeout": 30, "max_tokens": 4000, "temperature": "invalid"}
            )


class TestLoadConfiguration:
    """Test suite for load_configuration function."""

    def test_loads_complete_configuration_with_defaults(self):
        """Test loading complete configuration with default values."""
        config = load_configuration()

        # Verify all expected keys are present
        expected_keys = [
            "model",
            "api_key",
            "base_url",
            "timeout",
            "max_tokens",
            "temperature",
            "tools_enabled",
            "confirmation_required",
            "session_dir",
        ]
        for key in expected_keys:
            assert key in config

    def test_merges_file_config_with_defaults(self):
        """Test that file config overrides defaults."""
        with tempfile.NamedTemporaryFile(mode="w", suffix=".yaml", delete=False) as f:
            yaml.dump({"model": "custom-model", "timeout": 120}, f)
            temp_path = Path(f.name)

        try:
            config = load_configuration(temp_path)
            assert config["model"] == "custom-model"
            assert config["timeout"] == 120
            assert config["temperature"] == 0.7  # Should keep default
        finally:
            temp_path.unlink()

    def test_environment_overrides_file_and_defaults(self):
        """Test that environment variables override both file and defaults."""
        with tempfile.NamedTemporaryFile(mode="w", suffix=".yaml", delete=False) as f:
            yaml.dump({"model": "file-model", "timeout": 60}, f)
            temp_path = Path(f.name)

        try:
            with patch.dict(os.environ, {"MODEL": "env-model", "TEMPERATURE": "0.9"}):
                config = load_configuration(temp_path)
                assert config["model"] == "env-model"  # Env override
                assert config["timeout"] == 60  # From file
                assert config["temperature"] == 0.9  # Env override
        finally:
            temp_path.unlink()

    def test_raises_configuration_error_on_validation_failure(self):
        """Test that validation errors are propagated."""
        with tempfile.NamedTemporaryFile(mode="w", suffix=".yaml", delete=False) as f:
            yaml.dump({"timeout": -1}, f)
            temp_path = Path(f.name)

        try:
            with pytest.raises(ConfigurationError, match="timeout must be positive"):
                load_configuration(temp_path)
        finally:
            temp_path.unlink()
