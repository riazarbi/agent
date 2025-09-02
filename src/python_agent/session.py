"""File-based session management for conversation persistence.

This module provides session management functionality for saving and resuming
conversation history. Sessions are stored as JSON files with timestamp-based IDs.

Keywords: session, persistence, file-based, conversation, history, JSON

Main classes:
    - SessionManager: Primary interface for session operations
    - Session: Individual session data container

Basic usage:
    manager = SessionManager()
    session = manager.create_session()
    session.add_message("user", "Hello")
    manager.save_session(session)
    loaded = manager.load_session(session.session_id)
"""

import json
from datetime import datetime
from pathlib import Path
from typing import Any

from python_agent.config import load_configuration


class SessionError(Exception):
    """Base exception for session-related errors.

    Keywords: error, exception, session
    """


class Session:
    """Individual session data container for conversation history.

    Keywords: session, conversation, history, messages

    Attributes:
        session_id: Unique session identifier (timestamp-based)
        messages: List of conversation messages
        created_at: Session creation timestamp
        updated_at: Last update timestamp
    """

    def __init__(self, session_id: str) -> None:
        """Initialize session with ID and empty message history.

        Args:
            session_id: Unique session identifier
        """
        self.session_id = session_id
        self.messages: list[dict[str, Any]] = []
        self.created_at = datetime.now().isoformat()
        self.updated_at = self.created_at

    def add_message(self, role: str, content: str) -> None:
        """Add message to conversation history.

        Args:
            role: Message role (user, assistant, tool)
            content: Message content
        """
        self.messages.append(
            {"role": role, "content": content, "timestamp": datetime.now().isoformat()}
        )
        self.updated_at = datetime.now().isoformat()

    def to_dict(self) -> dict[str, Any]:
        """Convert session to dictionary for JSON serialization.

        Returns:
            Session data as dictionary
        """
        return {
            "session_id": self.session_id,
            "messages": self.messages,
            "created_at": self.created_at,
            "updated_at": self.updated_at,
        }

    @classmethod
    def from_dict(cls, data: dict[str, Any]) -> "Session":
        """Create session from dictionary data.

        Args:
            data: Session data dictionary

        Returns:
            Session instance
        """
        session = cls(data["session_id"])
        session.messages = data.get("messages", [])
        session.created_at = data.get("created_at", session.created_at)
        session.updated_at = data.get("updated_at", session.updated_at)
        return session


class SessionManager:
    """File-based session persistence manager.

    Keywords: session, manager, persistence, file-based, JSON

    Manages session creation, saving, loading, and directory structure.
    Sessions are stored in ~/.agent/sessions/ as JSON files.
    """

    def __init__(self, config_path: Path | None = None) -> None:
        """Initialize session manager with configuration.

        Args:
            config_path: Path to configuration file (optional)
        """
        self.config = load_configuration(config_path)
        self.sessions_dir = Path.home() / ".agent" / "sessions"
        self._ensure_sessions_directory()

    def _ensure_sessions_directory(self) -> None:
        """Create sessions directory if it doesn't exist."""
        self.sessions_dir.mkdir(parents=True, exist_ok=True)

    def _generate_session_id(self) -> str:
        """Generate unique session ID based on current timestamp.

        Returns:
            Session ID in format: YYYY-MM-DD-HH-MM-SS
        """
        return datetime.now().strftime("%Y-%m-%d-%H-%M-%S")

    def _get_session_path(self, session_id: str) -> Path:
        """Get file path for session ID.

        Args:
            session_id: Session identifier

        Returns:
            Path to session file
        """
        return self.sessions_dir / f"{session_id}.json"

    def create_session(self) -> Session:
        """Create new session with generated ID.

        Returns:
            New session instance
        """
        session_id = self._generate_session_id()
        return Session(session_id)

    def save_session(self, session: Session) -> None:
        """Save session to JSON file.

        Args:
            session: Session to save

        Raises:
            SessionError: If save operation fails
        """
        try:
            session_path = self._get_session_path(session.session_id)
            with open(session_path, "w", encoding="utf-8") as f:
                json.dump(session.to_dict(), f, indent=2, ensure_ascii=False)
        except OSError as e:
            raise SessionError(
                f"Failed to save session {session.session_id}: {e}"
            ) from e

    def load_session(self, session_id: str) -> Session:
        """Load session from JSON file.

        Args:
            session_id: Session identifier to load

        Returns:
            Loaded session instance

        Raises:
            SessionError: If session file not found or invalid
        """
        session_path = self._get_session_path(session_id)

        if not session_path.exists():
            raise SessionError(f"Session not found: {session_id}")

        try:
            with open(session_path, encoding="utf-8") as f:
                data = json.load(f)
            return Session.from_dict(data)
        except (OSError, json.JSONDecodeError) as e:
            raise SessionError(f"Failed to load session {session_id}: {e}") from e

    def list_sessions(self) -> list[str]:
        """List all available session IDs.

        Returns:
            List of session IDs sorted by creation date (newest first)
        """
        session_files = list(self.sessions_dir.glob("*.json"))
        session_ids = [f.stem for f in session_files]
        return sorted(session_ids, reverse=True)  # Newest first

    def session_exists(self, session_id: str) -> bool:
        """Check if session exists.

        Args:
            session_id: Session identifier to check

        Returns:
            True if session exists, False otherwise
        """
        return self._get_session_path(session_id).exists()
