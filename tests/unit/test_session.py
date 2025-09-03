"""Unit tests for session module.

Keywords: test, session, persistence, file-based, conversation, history
"""

import json
import tempfile
from pathlib import Path
from unittest.mock import patch

import pytest

from python_agent.session import Session, SessionError, SessionManager


class TestSession:
    """Test suite for Session class."""

    def test_session_initialization_with_id(self):
        """Test Session initialization with provided ID."""
        session_id = "test-session-123"
        session = Session(session_id)

        assert session.session_id == session_id
        assert session.messages == []
        assert session.created_at is not None
        assert session.updated_at is not None
        assert session.created_at == session.updated_at

    def test_add_message_updates_session(self):
        """Test adding message updates session properly."""
        session = Session("test-id")
        initial_updated_at = session.updated_at

        # Mock datetime to ensure timestamp differences
        with patch("python_agent.session.datetime") as mock_dt:
            mock_dt.now.return_value.isoformat.side_effect = [
                "2023-01-01T10:30:00",
                "2023-01-01T10:31:00",
            ]

            session.add_message("user", "Hello, world!")

        assert len(session.messages) == 1
        assert session.messages[0]["role"] == "user"
        assert session.messages[0]["content"] == "Hello, world!"
        assert session.messages[0]["timestamp"] == "2023-01-01T10:30:00"
        assert session.updated_at == "2023-01-01T10:31:00"
        assert session.updated_at != initial_updated_at

    def test_add_multiple_messages(self):
        """Test adding multiple messages maintains order."""
        session = Session("test-id")

        session.add_message("user", "First message")
        session.add_message("assistant", "Second message")
        session.add_message("tool", "Third message")

        assert len(session.messages) == 3
        assert session.messages[0]["role"] == "user"
        assert session.messages[0]["content"] == "First message"
        assert session.messages[1]["role"] == "assistant"
        assert session.messages[1]["content"] == "Second message"
        assert session.messages[2]["role"] == "tool"
        assert session.messages[2]["content"] == "Third message"

    def test_to_dict_conversion(self):
        """Test session conversion to dictionary."""
        session = Session("test-session")
        session.add_message("user", "Test message")

        result = session.to_dict()

        assert isinstance(result, dict)
        assert result["session_id"] == "test-session"
        assert len(result["messages"]) == 1
        assert result["messages"][0]["role"] == "user"
        assert result["messages"][0]["content"] == "Test message"
        assert "timestamp" in result["messages"][0]
        assert "created_at" in result
        assert "updated_at" in result

    def test_from_dict_creation(self):
        """Test creating session from dictionary data."""
        data = {
            "session_id": "restored-session",
            "messages": [
                {
                    "role": "user",
                    "content": "Hello",
                    "timestamp": "2023-01-01T10:00:00",
                },
                {
                    "role": "assistant",
                    "content": "Hi there",
                    "timestamp": "2023-01-01T10:01:00",
                },
            ],
            "created_at": "2023-01-01T09:00:00",
            "updated_at": "2023-01-01T10:01:00",
        }

        session = Session.from_dict(data)

        assert session.session_id == "restored-session"
        assert len(session.messages) == 2
        assert session.messages[0]["role"] == "user"
        assert session.messages[0]["content"] == "Hello"
        assert session.messages[1]["role"] == "assistant"
        assert session.messages[1]["content"] == "Hi there"
        assert session.created_at == "2023-01-01T09:00:00"
        assert session.updated_at == "2023-01-01T10:01:00"

    def test_from_dict_with_minimal_data(self):
        """Test creating session from minimal dictionary data."""
        data = {"session_id": "minimal-session"}

        session = Session.from_dict(data)

        assert session.session_id == "minimal-session"
        assert session.messages == []
        assert session.created_at is not None
        assert session.updated_at is not None

    def test_from_dict_with_missing_optional_fields(self):
        """Test creating session from dict with missing optional fields."""
        data = {
            "session_id": "partial-session",
            "messages": [{"role": "user", "content": "test"}],
        }

        session = Session.from_dict(data)

        assert session.session_id == "partial-session"
        assert len(session.messages) == 1
        assert session.created_at is not None
        assert session.updated_at is not None


class TestSessionError:
    """Test suite for SessionError exception."""

    def test_session_error_inheritance(self):
        """Test SessionError is proper Exception subclass."""
        error = SessionError("test message")

        assert isinstance(error, Exception)
        assert str(error) == "test message"

    def test_session_error_with_empty_message(self):
        """Test SessionError with empty message."""
        error = SessionError("")

        assert isinstance(error, Exception)
        assert str(error) == ""


class TestSessionManager:
    """Test suite for SessionManager class."""

    @pytest.fixture
    def temp_sessions_dir(self):
        """Create temporary directory for session tests."""
        with tempfile.TemporaryDirectory() as temp_dir:
            yield Path(temp_dir)

    @pytest.fixture
    def mock_home_path(self, temp_sessions_dir):
        """Mock Path.home() to use temporary directory."""
        with patch("python_agent.session.Path.home") as mock_home:
            # Point home to temp dir so .agent/sessions gets created there
            mock_home.return_value = temp_sessions_dir
            yield temp_sessions_dir

    def test_session_manager_initialization(self, mock_home_path):
        """Test SessionManager initialization creates directory."""
        manager = SessionManager()

        expected_sessions_dir = mock_home_path / ".agent" / "sessions"
        assert manager.sessions_dir == expected_sessions_dir
        assert expected_sessions_dir.exists()

    @patch("python_agent.session.load_configuration")
    def test_session_manager_with_custom_config(self, mock_load_config, mock_home_path):
        """Test SessionManager with custom config path."""
        mock_config = {"test": "value"}
        mock_load_config.return_value = mock_config
        config_path = Path("/custom/config.yaml")

        manager = SessionManager(config_path)

        mock_load_config.assert_called_once_with(config_path)
        assert manager.config == mock_config

    def test_generate_session_id_format(self, mock_home_path):
        """Test session ID generation follows expected format."""
        manager = SessionManager()

        with patch("python_agent.session.datetime") as mock_dt:
            mock_dt.now.return_value.strftime.return_value = "2023-01-01-15-30-45"
            session_id = manager._generate_session_id()

        assert session_id == "2023-01-01-15-30-45"
        mock_dt.now.return_value.strftime.assert_called_once_with("%Y-%m-%d-%H-%M-%S")

    def test_get_session_path(self, mock_home_path):
        """Test session path generation."""
        manager = SessionManager()
        session_id = "2023-01-01-12-00-00"

        path = manager._get_session_path(session_id)

        expected_path = manager.sessions_dir / "2023-01-01-12-00-00.json"
        assert path == expected_path

    def test_create_session(self, mock_home_path):
        """Test session creation."""
        manager = SessionManager()

        with patch.object(manager, "_generate_session_id") as mock_gen_id:
            mock_gen_id.return_value = "test-session-id"
            session = manager.create_session()

        assert isinstance(session, Session)
        assert session.session_id == "test-session-id"
        assert session.messages == []

    def test_save_session_success(self, mock_home_path):
        """Test successful session saving."""
        manager = SessionManager()
        session = Session("test-save")
        session.add_message("user", "Test message")

        manager.save_session(session)

        session_file = manager.sessions_dir / "test-save.json"
        assert session_file.exists()

        with open(session_file) as f:
            saved_data = json.load(f)

        assert saved_data["session_id"] == "test-save"
        assert len(saved_data["messages"]) == 1
        assert saved_data["messages"][0]["content"] == "Test message"

    def test_save_session_file_write_error(self, mock_home_path):
        """Test session save with file write error."""
        manager = SessionManager()
        session = Session("test-error")

        # Mock open to raise OSError
        with (
            patch("builtins.open", side_effect=OSError("Permission denied")),
            pytest.raises(SessionError, match="Failed to save session test-error"),
        ):
                manager.save_session(session)

    def test_load_session_success(self, mock_home_path):
        """Test successful session loading."""
        manager = SessionManager()
        session_data = {
            "session_id": "test-load",
            "messages": [
                {"role": "user", "content": "Hello", "timestamp": "2023-01-01T10:00:00"}
            ],
            "created_at": "2023-01-01T09:00:00",
            "updated_at": "2023-01-01T10:00:00",
        }

        # Create session file
        session_file = manager.sessions_dir / "test-load.json"
        with open(session_file, "w") as f:
            json.dump(session_data, f)

        loaded_session = manager.load_session("test-load")

        assert loaded_session.session_id == "test-load"
        assert len(loaded_session.messages) == 1
        assert loaded_session.messages[0]["content"] == "Hello"
        assert loaded_session.created_at == "2023-01-01T09:00:00"

    def test_load_session_file_not_found(self, mock_home_path):
        """Test loading non-existent session."""
        manager = SessionManager()

        with pytest.raises(SessionError, match="Session not found: non-existent"):
            manager.load_session("non-existent")

    def test_load_session_invalid_json(self, mock_home_path):
        """Test loading session with invalid JSON."""
        manager = SessionManager()
        session_file = manager.sessions_dir / "invalid-json.json"

        # Create file with invalid JSON
        with open(session_file, "w") as f:
            f.write("{ invalid json }")

        with pytest.raises(SessionError, match="Failed to load session invalid-json"):
            manager.load_session("invalid-json")

    def test_load_session_file_read_error(self, mock_home_path):
        """Test loading session with file read error."""
        manager = SessionManager()
        session_file = manager.sessions_dir / "read-error.json"
        session_file.touch()  # Create empty file

        # Mock open to raise OSError on read
        with (
            patch("builtins.open", side_effect=OSError("Permission denied")),
            pytest.raises(SessionError, match="Failed to load session read-error"),
        ):
                manager.load_session("read-error")

    def test_list_sessions_empty(self, mock_home_path):
        """Test listing sessions when directory is empty."""
        manager = SessionManager()

        sessions = manager.list_sessions()

        assert sessions == []

    def test_list_sessions_multiple_files(self, mock_home_path):
        """Test listing multiple session files."""
        manager = SessionManager()

        # Create session files
        session_files = [
            "2023-01-01-10-00-00.json",
            "2023-01-02-11-00-00.json",
            "2023-01-01-09-00-00.json",
        ]
        for filename in session_files:
            (manager.sessions_dir / filename).touch()

        sessions = manager.list_sessions()

        # Should return session IDs (without .json) sorted newest first
        expected = ["2023-01-02-11-00-00", "2023-01-01-10-00-00", "2023-01-01-09-00-00"]
        assert sessions == expected

    def test_list_sessions_ignores_non_json_files(self, mock_home_path):
        """Test listing sessions ignores non-JSON files."""
        manager = SessionManager()

        # Create mixed files
        (manager.sessions_dir / "session1.json").touch()
        (manager.sessions_dir / "session2.json").touch()
        (manager.sessions_dir / "not-session.txt").touch()
        (manager.sessions_dir / "README.md").touch()

        sessions = manager.list_sessions()

        # Should only include JSON files
        expected = ["session2", "session1"]  # Sorted reverse alphabetically
        assert set(sessions) == set(expected)
        assert len(sessions) == 2

    def test_session_exists_true(self, mock_home_path):
        """Test session_exists returns True for existing session."""
        manager = SessionManager()
        session_file = manager.sessions_dir / "existing-session.json"
        session_file.touch()

        result = manager.session_exists("existing-session")

        assert result is True

    def test_session_exists_false(self, mock_home_path):
        """Test session_exists returns False for non-existent session."""
        manager = SessionManager()

        result = manager.session_exists("non-existent")

        assert result is False


class TestSessionIntegration:
    """Integration tests for session functionality."""

    @pytest.fixture
    def temp_sessions_dir(self):
        """Create temporary directory for integration tests."""
        with tempfile.TemporaryDirectory() as temp_dir:
            yield Path(temp_dir)

    @pytest.fixture
    def mock_home_path(self, temp_sessions_dir):
        """Mock Path.home() for integration tests."""
        with patch("python_agent.session.Path.home") as mock_home:
            mock_home.return_value = temp_sessions_dir
            yield temp_sessions_dir

    @pytest.mark.integration
    def test_complete_session_workflow(self, mock_home_path):
        """Test complete session creation, save, and load workflow."""
        manager = SessionManager()

        # Create and populate session
        session = manager.create_session()
        session.add_message("user", "Hello")
        session.add_message("assistant", "Hi there!")

        # Save session
        manager.save_session(session)

        # Verify session exists
        assert manager.session_exists(session.session_id)

        # Load session
        loaded_session = manager.load_session(session.session_id)

        # Verify loaded session matches original
        assert loaded_session.session_id == session.session_id
        assert len(loaded_session.messages) == 2
        assert loaded_session.messages[0]["content"] == "Hello"
        assert loaded_session.messages[1]["content"] == "Hi there!"
        assert loaded_session.created_at == session.created_at

    @pytest.mark.integration
    def test_session_file_persistence(self, mock_home_path):
        """Test that session files persist correctly on filesystem."""
        manager = SessionManager()

        # Create session with known ID
        session = Session("persistence-test")
        session.add_message("user", "Persistent message")
        manager.save_session(session)

        # Create new manager instance (simulating app restart)
        new_manager = SessionManager()

        # Verify session still exists and loads correctly
        assert new_manager.session_exists("persistence-test")
        loaded_session = new_manager.load_session("persistence-test")

        assert loaded_session.session_id == "persistence-test"
        assert len(loaded_session.messages) == 1
        assert loaded_session.messages[0]["content"] == "Persistent message"
