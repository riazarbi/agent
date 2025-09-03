"""Unit tests for agent module.

Keywords: test, agent, litellm, conversation, chat, model, AI
"""

from unittest.mock import Mock, patch, MagicMock
from typing import Any

import pytest

from python_agent.agent import Agent, AgentError, ModelError
from python_agent.bash_tool import BashTool
from python_agent.session import Session, SessionManager


class TestAgentError:
    """Test suite for AgentError exception class."""

    def test_agent_error_inheritance(self):
        """Test AgentError inherits from Exception."""
        assert issubclass(AgentError, Exception)

    def test_agent_error_with_message(self):
        """Test AgentError with custom message."""
        message = "Test error message"
        error = AgentError(message)
        assert str(error) == message


class TestModelError:
    """Test suite for ModelError exception class."""

    def test_model_error_inheritance(self):
        """Test ModelError inherits from AgentError."""
        assert issubclass(ModelError, AgentError)
        assert issubclass(ModelError, Exception)

    def test_model_error_with_message(self):
        """Test ModelError with custom message."""
        message = "Model API failed"
        error = ModelError(message)
        assert str(error) == message


class TestAgentInit:
    """Test suite for Agent initialization."""

    def test_agent_initialization_with_minimal_config(self):
        """Test Agent initialization with minimal configuration."""
        config = {"model": "gpt-3.5-turbo", "max_tokens": 100, "temperature": 0.7, "timeout": 30}
        
        agent = Agent(config)
        
        assert agent.config == config
        assert agent.conversation_history == []
        assert isinstance(agent.session_manager, SessionManager)
        assert agent.current_session is None
        assert isinstance(agent.bash_tool, BashTool)

    def test_agent_initialization_with_full_config(self):
        """Test Agent initialization with full configuration."""
        config = {
            "model": "gpt-4",
            "max_tokens": 200,
            "temperature": 0.8,
            "timeout": 60,
            "confirmation_required": True,
            "tools_enabled": False,
            "base_url": "https://custom.api.url"
        }
        
        with patch('python_agent.agent.litellm') as mock_litellm:
            agent = Agent(config)
            
            assert agent.config == config
            assert agent.bash_tool.confirmation_required is True
            assert agent.bash_tool.timeout == 60
            assert agent.bash_tool.enabled is False
            # litellm.api_base should be set to the base_url
            assert mock_litellm.api_base == 'https://custom.api.url'

    def test_agent_initialization_without_base_url(self):
        """Test Agent initialization without base_url in config."""
        config = {"model": "gpt-3.5-turbo", "max_tokens": 100, "temperature": 0.7, "timeout": 30}
        
        # Since no base_url in config, litellm.api_base should not be modified
        agent = Agent(config)
        
        assert agent.config == config


class TestAgentAddMessage:
    """Test suite for Agent.add_message method."""

    @pytest.fixture
    def agent(self):
        """Create Agent instance for testing."""
        config = {"model": "test", "max_tokens": 100, "temperature": 0.7, "timeout": 30}
        return Agent(config)

    def test_add_message_to_conversation_history(self, agent):
        """Test adding message to conversation history."""
        agent.add_message("user", "Hello")
        
        assert len(agent.conversation_history) == 1
        assert agent.conversation_history[0] == {"role": "user", "content": "Hello"}

    def test_add_multiple_messages(self, agent):
        """Test adding multiple messages to conversation history."""
        agent.add_message("user", "Hello")
        agent.add_message("assistant", "Hi there!")
        agent.add_message("user", "How are you?")
        
        assert len(agent.conversation_history) == 3
        assert agent.conversation_history[0] == {"role": "user", "content": "Hello"}
        assert agent.conversation_history[1] == {"role": "assistant", "content": "Hi there!"}
        assert agent.conversation_history[2] == {"role": "user", "content": "How are you?"}

    def test_add_message_with_current_session(self, agent):
        """Test adding message when current session is active."""
        # Create mock session
        mock_session = Mock(spec=Session)
        agent.current_session = mock_session
        
        agent.add_message("user", "Test message")
        
        # Verify message added to conversation history
        assert len(agent.conversation_history) == 1
        assert agent.conversation_history[0] == {"role": "user", "content": "Test message"}
        
        # Verify session.add_message was called
        mock_session.add_message.assert_called_once_with("user", "Test message")


class TestAgentSessionManagement:
    """Test suite for Agent session management methods."""

    @pytest.fixture
    def agent(self):
        """Create Agent instance for testing."""
        config = {"model": "test", "max_tokens": 100, "temperature": 0.7, "timeout": 30}
        return Agent(config)

    def test_start_new_session(self, agent):
        """Test starting a new session."""
        mock_session = Mock(spec=Session)
        agent.session_manager.create_session = Mock(return_value=mock_session)
        
        agent.start_new_session()
        
        assert agent.current_session == mock_session
        agent.session_manager.create_session.assert_called_once()

    def test_resume_from_session(self, agent):
        """Test resuming from an existing session."""
        mock_session = Mock(spec=Session)
        mock_session.messages = [
            {"role": "user", "content": "Previous message"},
            {"role": "assistant", "content": "Previous response"}
        ]
        
        agent.resume_from_session(mock_session)
        
        assert agent.current_session == mock_session
        assert agent.conversation_history == mock_session.messages

    def test_save_current_session_with_session(self, agent):
        """Test saving current session when session exists."""
        mock_session = Mock(spec=Session)
        agent.current_session = mock_session
        agent.session_manager.save_session = Mock()
        
        agent.save_current_session()
        
        agent.session_manager.save_session.assert_called_once_with(mock_session)

    def test_save_current_session_without_session(self, agent):
        """Test saving current session when no session exists."""
        agent.current_session = None
        agent.session_manager.save_session = Mock()
        
        agent.save_current_session()
        
        agent.session_manager.save_session.assert_not_called()


class TestAgentChatCompletion:
    """Test suite for Agent.chat_completion method."""

    @pytest.fixture
    def agent(self):
        """Create Agent instance for testing."""
        config = {
            "model": "gpt-3.5-turbo",
            "max_tokens": 100,
            "temperature": 0.7,
            "timeout": 30
        }
        return Agent(config)

    @patch('python_agent.agent.litellm.completion')
    def test_chat_completion_success(self, mock_completion, agent):
        """Test successful chat completion."""
        # Mock successful response
        mock_response = Mock()
        mock_response.choices = [Mock()]
        mock_response.choices[0].message.content = "Test response"
        mock_completion.return_value = mock_response
        
        messages = [{"role": "user", "content": "Hello"}]
        result = agent.chat_completion(messages)
        
        assert result == "Test response"
        mock_completion.assert_called_once_with(
            model="gpt-3.5-turbo",
            messages=messages,
            max_tokens=100,
            temperature=0.7,
            timeout=30
        )

    @patch('python_agent.agent.litellm.completion')
    def test_chat_completion_with_none_content(self, mock_completion, agent):
        """Test chat completion when response content is None."""
        # Mock response with None content
        mock_response = Mock()
        mock_response.choices = [Mock()]
        mock_response.choices[0].message.content = None
        mock_completion.return_value = mock_response
        
        messages = [{"role": "user", "content": "Hello"}]
        result = agent.chat_completion(messages)
        
        assert result == ""

    @patch('python_agent.agent.litellm.completion')
    def test_chat_completion_api_error(self, mock_completion, agent):
        """Test chat completion when API call fails."""
        mock_completion.side_effect = Exception("API Error")
        
        messages = [{"role": "user", "content": "Hello"}]
        
        with pytest.raises(ModelError, match="Model API call failed: API Error"):
            agent.chat_completion(messages)


class TestAgentProcessSinglePrompt:
    """Test suite for Agent.process_single_prompt method."""

    @pytest.fixture
    def agent(self):
        """Create Agent instance for testing."""
        config = {"model": "test", "max_tokens": 100, "temperature": 0.7, "timeout": 30}
        return Agent(config)

    def test_process_single_prompt(self, agent):
        """Test processing single prompt."""
        with patch.object(agent, 'chat_completion') as mock_chat:
            mock_chat.return_value = "Test response"
            
            result = agent.process_single_prompt("Test prompt")
            
            assert result == "Test response"
            mock_chat.assert_called_once_with([{"role": "user", "content": "Test prompt"}])


class TestAgentInteractiveLoop:
    """Test suite for Agent.interactive_loop method."""

    @pytest.fixture
    def agent(self):
        """Create Agent instance for testing."""
        config = {
            "model": "test",
            "max_tokens": 100,
            "temperature": 0.7,
            "timeout": 30,
            "verbose": False,
            "quiet": False
        }
        return Agent(config)

    @patch('builtins.input')
    @patch('builtins.print')
    def test_interactive_loop_exit_command(self, mock_print, mock_input, agent):
        """Test interactive loop with exit command."""
        mock_input.return_value = "exit"
        
        with patch.object(agent, 'start_new_session') as mock_start:
            with patch.object(agent, 'save_current_session') as mock_save:
                mock_session = Mock()
                mock_session.session_id = "test-session"
                
                agent.interactive_loop()
                
                # start_new_session should be called because current_session is None
                mock_start.assert_called_once()
                mock_save.assert_called_once()

    @patch('builtins.input')
    @patch('builtins.print')
    def test_interactive_loop_with_conversation(self, mock_print, mock_input, agent):
        """Test interactive loop with user input and response."""
        mock_input.side_effect = ["Hello", "exit"]
        
        with patch.object(agent, 'start_new_session'):
            with patch.object(agent, 'save_current_session'):
                with patch.object(agent, 'chat_completion') as mock_chat:
                    with patch.object(agent, 'add_message') as mock_add:
                        mock_chat.return_value = "Hi there!"
                        mock_session = Mock()
                        mock_session.session_id = "test-session"
                        agent.current_session = mock_session
                        
                        agent.interactive_loop()
                        
                        # Verify user message was added
                        mock_add.assert_any_call("user", "Hello")
                        # Verify assistant response was added
                        mock_add.assert_any_call("assistant", "Hi there!")
                        # Verify chat completion was called
                        mock_chat.assert_called_once()

    @patch('builtins.input')
    @patch('builtins.print')
    def test_interactive_loop_with_empty_input(self, mock_print, mock_input, agent):
        """Test interactive loop with empty input."""
        mock_input.side_effect = ["", "  ", "exit"]
        
        with patch.object(agent, 'start_new_session'):
            with patch.object(agent, 'save_current_session'):
                with patch.object(agent, 'chat_completion') as mock_chat:
                    mock_session = Mock()
                    mock_session.session_id = "test-session"
                    agent.current_session = mock_session
                    
                    agent.interactive_loop()
                    
                    # Verify chat completion was not called for empty inputs
                    mock_chat.assert_not_called()

    @patch('builtins.input')
    @patch('builtins.print')
    def test_interactive_loop_keyboard_interrupt(self, mock_print, mock_input, agent):
        """Test interactive loop handles KeyboardInterrupt."""
        mock_input.side_effect = KeyboardInterrupt()
        
        with patch.object(agent, 'start_new_session'):
            with patch.object(agent, 'save_current_session') as mock_save:
                mock_session = Mock()
                mock_session.session_id = "test-session"
                agent.current_session = mock_session
                
                agent.interactive_loop()
                
                mock_save.assert_called()

    @patch('builtins.input')
    @patch('builtins.print')
    def test_interactive_loop_model_error(self, mock_print, mock_input, agent):
        """Test interactive loop handles ModelError."""
        mock_input.side_effect = ["Hello", "exit"]
        
        with patch.object(agent, 'start_new_session'):
            with patch.object(agent, 'save_current_session'):
                with patch.object(agent, 'chat_completion') as mock_chat:
                    mock_chat.side_effect = ModelError("API failed")
                    mock_session = Mock()
                    mock_session.session_id = "test-session"
                    agent.current_session = mock_session
                    
                    agent.interactive_loop()
                    
                    # Verify error message was printed
                    mock_print.assert_any_call("Error: API failed")

    @patch('builtins.input')
    @patch('builtins.print')
    def test_interactive_loop_unexpected_error(self, mock_print, mock_input, agent):
        """Test interactive loop handles unexpected errors."""
        mock_input.side_effect = ["Hello", "exit"]
        
        with patch.object(agent, 'start_new_session'):
            with patch.object(agent, 'save_current_session'):
                with patch.object(agent, 'chat_completion') as mock_chat:
                    mock_chat.side_effect = RuntimeError("Unexpected error")
                    mock_session = Mock()
                    mock_session.session_id = "test-session"
                    agent.current_session = mock_session
                    
                    agent.interactive_loop()
                    
                    # Verify error message was printed
                    mock_print.assert_any_call("Unexpected error: Unexpected error")

    @patch('builtins.input')
    @patch('builtins.print')
    def test_interactive_loop_verbose_mode(self, mock_print, mock_input, agent):
        """Test interactive loop with verbose mode enabled."""
        agent.config["verbose"] = True
        mock_input.side_effect = ["Hello", "exit"]
        
        with patch.object(agent, 'start_new_session') as mock_start:
            with patch.object(agent, 'save_current_session'):
                with patch.object(agent, 'chat_completion') as mock_chat:
                    mock_chat.return_value = "Response"
                    mock_session = Mock()
                    mock_session.session_id = "test-session"
                    
                    # Mock start_new_session to set the current_session
                    def set_session():
                        agent.current_session = mock_session
                    mock_start.side_effect = set_session
                    
                    agent.interactive_loop()
                    
                    # Verify verbose messages were printed
                    mock_print.assert_any_call("Started session: test-session")
                    mock_print.assert_any_call("Agent: Thinking...")

    @patch('builtins.input')
    @patch('builtins.print')
    def test_interactive_loop_quiet_mode(self, mock_print, mock_input, agent):
        """Test interactive loop with quiet mode enabled."""
        agent.config["quiet"] = True
        mock_input.return_value = "exit"
        
        with patch.object(agent, 'start_new_session'):
            with patch.object(agent, 'save_current_session'):
                mock_session = Mock()
                mock_session.session_id = "test-session"
                agent.current_session = mock_session
                
                agent.interactive_loop()
                
                # Verify greeting message was not printed
                print_calls = [call[0][0] for call in mock_print.call_args_list]
                assert "AI Coding Agent (type 'exit' to quit)" not in print_calls


@pytest.mark.integration
class TestAgentIntegration:
    """Integration tests for Agent class."""

    def test_agent_with_real_session_manager(self):
        """Test Agent integration with real SessionManager."""
        config = {"model": "test", "max_tokens": 100, "temperature": 0.7, "timeout": 30}
        agent = Agent(config)
        
        # Test session creation and management
        agent.start_new_session()
        assert agent.current_session is not None
        
        # Test adding messages
        agent.add_message("user", "Test message")
        assert len(agent.conversation_history) == 1
        assert len(agent.current_session.messages) == 1

    def test_agent_session_resume_integration(self):
        """Test Agent session resume integration."""
        config = {"model": "test", "max_tokens": 100, "temperature": 0.7, "timeout": 30}
        agent = Agent(config)
        
        # Create session with messages
        session = Session("test-session-id")
        session.add_message("user", "Previous message")
        session.add_message("assistant", "Previous response")
        
        # Resume session
        agent.resume_from_session(session)
        
        assert agent.current_session == session
        assert len(agent.conversation_history) == 2
        assert agent.conversation_history[0]["content"] == "Previous message"
        assert agent.conversation_history[1]["content"] == "Previous response"