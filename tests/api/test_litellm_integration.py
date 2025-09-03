"""API tests for LiteLLM integration (real APIs, no mocking).

This module tests the agent's integration with various LiteLLM-supported model 
providers including OpenAI, Anthropic, and other services. Tests are designed
to run against real APIs with proper skip logic when credentials are unavailable.

Keywords: api, litellm, integration, testing, openai, anthropic, model, providers
"""

import os
import pytest
from typing import Any, Dict

from python_agent.agent import Agent, AgentError, ModelError
from python_agent.config import get_default_config


class TestLiteLLMIntegration:
    """Test LiteLLM integration with real model providers."""
    
    def _get_api_key(self, env_var: str) -> str | None:
        """Get API key from environment with proper handling."""
        return os.getenv(env_var)
    
    def _skip_if_no_api_key(self, api_key: str | None, provider: str, env_var: str) -> None:
        """Skip test if no API key is provided."""
        if not api_key:
            pytest.skip(f"Skipping {provider} API test: no API key provided (set {env_var})")
    
    def _create_test_config(self, api_key: str, model: str, base_url: str | None = None) -> Dict[str, Any]:
        """Create test configuration with API credentials."""
        config = get_default_config()
        config.update({
            "api_key": api_key,
            "model": model,
            "max_tokens": 100,  # Keep responses small for testing
            "temperature": 0.1,  # Consistent responses
            "tools_enabled": False,  # Disable tools for API tests
            "confirmation_required": False,
        })
        if base_url:
            config["base_url"] = base_url
        return config

    @pytest.mark.api
    def test_openai_gpt_integration_chat_completion(self):
        """Test successful chat completion with OpenAI GPT models."""
        api_key = self._get_api_key("OPENAI_API_KEY")
        self._skip_if_no_api_key(api_key, "OpenAI", "OPENAI_API_KEY")
        
        config = self._create_test_config(api_key, "gpt-3.5-turbo")
        agent = Agent(config)
        
        # Test simple chat completion
        response = agent.chat_completion([
            {"role": "user", "content": "Say 'Hello, World!' and nothing else."}
        ])
        
        assert response is not None
        assert isinstance(response, str)
        assert len(response) > 0
        assert "Hello" in response or "hello" in response

    @pytest.mark.api
    def test_anthropic_claude_integration_chat_completion(self):
        """Test successful chat completion with Anthropic Claude models."""
        api_key = self._get_api_key("ANTHROPIC_API_KEY")
        self._skip_if_no_api_key(api_key, "Anthropic", "ANTHROPIC_API_KEY")
        
        config = self._create_test_config(api_key, "claude-3-haiku-20240307")
        agent = Agent(config)
        
        # Test simple chat completion
        response = agent.chat_completion([
            {"role": "user", "content": "Respond with exactly 'API test successful'"}
        ])
        
        assert response is not None
        assert isinstance(response, str)
        assert len(response) > 0

    @pytest.mark.api
    def test_google_gemini_integration_chat_completion(self):
        """Test successful chat completion with Google Gemini models."""
        api_key = self._get_api_key("GEMINI_API_KEY")
        self._skip_if_no_api_key(api_key, "Google Gemini", "GEMINI_API_KEY")
        
        config = self._create_test_config(api_key, "gemini/gemini-1.5-flash")
        agent = Agent(config)
        
        # Test simple chat completion
        response = agent.chat_completion([
            {"role": "user", "content": "Say hello in exactly 3 words."}
        ])
        
        assert response is not None
        assert isinstance(response, str)
        assert len(response) > 0

    @pytest.mark.api
    def test_openai_conversation_history_management(self):
        """Test conversation history handling with multiple messages."""
        api_key = self._get_api_key("OPENAI_API_KEY")
        self._skip_if_no_api_key(api_key, "OpenAI", "OPENAI_API_KEY")
        
        config = self._create_test_config(api_key, "gpt-3.5-turbo")
        agent = Agent(config)
        
        # Add multiple messages to conversation
        agent.add_message("user", "My name is TestUser")
        agent.add_message("assistant", "Hello TestUser!")
        agent.add_message("user", "What is my name?")
        
        # Test that conversation history is maintained
        response = agent.chat_completion(agent.conversation_history)
        
        assert response is not None
        assert isinstance(response, str)
        assert len(response) > 0
        # Should reference the name from conversation history
        assert "TestUser" in response or "test" in response.lower()

    @pytest.mark.api
    def test_api_error_handling_invalid_key(self):
        """Test proper error handling for invalid API keys."""
        config = self._create_test_config("invalid_key_12345", "gpt-3.5-turbo")
        agent = Agent(config)
        
        # Should raise ModelError for invalid authentication
        with pytest.raises(ModelError) as exc_info:
            agent.chat_completion([
                {"role": "user", "content": "This should fail"}
            ])
        
        # Verify error message contains relevant information
        error_message = str(exc_info.value).lower()
        assert any(keyword in error_message for keyword in [
            "authentication", "invalid", "key", "unauthorized", "api"
        ])

    @pytest.mark.api
    def test_api_error_handling_invalid_model(self):
        """Test proper error handling for invalid model names."""
        api_key = self._get_api_key("OPENAI_API_KEY")
        self._skip_if_no_api_key(api_key, "OpenAI", "OPENAI_API_KEY")
        
        config = self._create_test_config(api_key, "nonexistent-model-12345")
        agent = Agent(config)
        
        # Should raise ModelError for invalid model
        with pytest.raises(ModelError) as exc_info:
            agent.chat_completion([
                {"role": "user", "content": "This should fail"}
            ])
        
        # Verify error message contains model-related information
        error_message = str(exc_info.value).lower()
        assert any(keyword in error_message for keyword in [
            "model", "invalid", "not found", "unsupported"
        ])

    @pytest.mark.api
    def test_process_single_prompt_integration(self):
        """Test process_single_prompt method with real API."""
        api_key = self._get_api_key("OPENAI_API_KEY")
        self._skip_if_no_api_key(api_key, "OpenAI", "OPENAI_API_KEY")
        
        config = self._create_test_config(api_key, "gpt-3.5-turbo")
        agent = Agent(config)
        
        # Test single prompt processing
        response = agent.process_single_prompt("Respond with 'Single prompt test successful'")
        
        assert response is not None
        assert isinstance(response, str)
        assert len(response) > 0
        
        # Verify conversation history was updated
        assert len(agent.conversation_history) == 2  # user + assistant messages
        assert agent.conversation_history[0]["role"] == "user"
        assert agent.conversation_history[1]["role"] == "assistant"

    @pytest.mark.api
    def test_custom_base_url_configuration(self):
        """Test configuration with custom base URL."""
        api_key = self._get_api_key("OPENAI_API_KEY")
        self._skip_if_no_api_key(api_key, "OpenAI", "OPENAI_API_KEY")
        
        # Test with OpenAI's actual base URL (should work)
        config = self._create_test_config(
            api_key, 
            "gpt-3.5-turbo", 
            base_url="https://api.openai.com/v1"
        )
        agent = Agent(config)
        
        response = agent.chat_completion([
            {"role": "user", "content": "Say 'Custom base URL test'"}
        ])
        
        assert response is not None
        assert isinstance(response, str)
        assert len(response) > 0

    @pytest.mark.api
    def test_temperature_and_max_tokens_configuration(self):
        """Test that temperature and max_tokens configuration is respected."""
        api_key = self._get_api_key("OPENAI_API_KEY")
        self._skip_if_no_api_key(api_key, "OpenAI", "OPENAI_API_KEY")
        
        config = self._create_test_config(api_key, "gpt-3.5-turbo")
        config.update({
            "temperature": 0.0,  # Very deterministic
            "max_tokens": 10,    # Very short response
        })
        agent = Agent(config)
        
        response = agent.chat_completion([
            {"role": "user", "content": "Write a long essay about artificial intelligence."}
        ])
        
        assert response is not None
        assert isinstance(response, str)
        # Response should be short due to max_tokens limit
        assert len(response.split()) <= 15  # Allow some flexibility

    @pytest.mark.api
    def test_empty_content_handling(self):
        """Test handling of responses with empty or None content."""
        api_key = self._get_api_key("OPENAI_API_KEY")
        self._skip_if_no_api_key(api_key, "OpenAI", "OPENAI_API_KEY")
        
        config = self._create_test_config(api_key, "gpt-3.5-turbo")
        config["max_tokens"] = 1  # Very limited to potentially get empty responses
        agent = Agent(config)
        
        # This might produce a very short or empty response
        response = agent.chat_completion([
            {"role": "user", "content": "Just say ok"}
        ])
        
        # Should handle empty responses gracefully
        assert response is not None
        assert isinstance(response, str)


class TestLiteLLMConversationWorkflows:
    """Test complete conversation workflows with LiteLLM."""
    
    def _get_openai_key(self) -> str | None:
        """Get OpenAI API key for workflow tests."""
        return os.getenv("OPENAI_API_KEY")
    
    def _create_test_agent(self, api_key: str) -> Agent:
        """Create test agent with OpenAI configuration."""
        config = get_default_config()
        config.update({
            "api_key": api_key,
            "model": "gpt-3.5-turbo",
            "max_tokens": 150,
            "temperature": 0.1,
            "tools_enabled": False,
            "confirmation_required": False,
        })
        return Agent(config)

    @pytest.mark.api
    def test_multi_turn_conversation_workflow(self):
        """Test multi-turn conversation with context preservation."""
        api_key = self._get_openai_key()
        if not api_key:
            pytest.skip("Skipping conversation workflow test: no OpenAI API key (set OPENAI_API_KEY)")
        
        agent = self._create_test_agent(api_key)
        
        # Turn 1: Initial question
        response1 = agent.process_single_prompt("I have 5 apples.")
        assert response1 is not None
        assert len(agent.conversation_history) == 2
        
        # Turn 2: Follow-up question referencing previous context
        agent.add_message("user", "How many apples do I have?")
        response2 = agent.chat_completion(agent.conversation_history)
        assert response2 is not None
        assert "5" in response2 or "five" in response2.lower()
        
        # Turn 3: Another follow-up
        agent.add_message("assistant", response2)
        agent.add_message("user", "If I eat 2, how many will I have left?")
        response3 = agent.chat_completion(agent.conversation_history)
        assert response3 is not None
        assert "3" in response3 or "three" in response3.lower()

    @pytest.mark.api
    def test_conversation_with_session_management(self):
        """Test conversation workflow with session persistence."""
        api_key = self._get_openai_key()
        if not api_key:
            pytest.skip("Skipping session workflow test: no OpenAI API key (set OPENAI_API_KEY)")
        
        agent = self._create_test_agent(api_key)
        
        # Start a new session
        session = agent.start_new_session()
        assert session is not None
        assert agent.current_session is not None
        
        # Have a conversation
        response = agent.process_single_prompt("Remember that my favorite color is blue.")
        assert response is not None
        
        # Save session
        agent.save_current_session()
        session_id = agent.current_session.session_id
        
        # Create new agent and resume session
        new_agent = self._create_test_agent(api_key)
        resumed_session = new_agent.resume_from_session(session_id)
        assert resumed_session is not None
        assert len(new_agent.conversation_history) > 0
        
        # Test that context is preserved
        new_agent.add_message("user", "What is my favorite color?")
        response = new_agent.chat_completion(new_agent.conversation_history)
        assert response is not None
        assert "blue" in response.lower()

    @pytest.mark.api
    def test_error_recovery_workflow(self):
        """Test conversation recovery after API errors."""
        api_key = self._get_openai_key()
        if not api_key:
            pytest.skip("Skipping error recovery test: no OpenAI API key (set OPENAI_API_KEY)")
        
        agent = self._create_test_agent(api_key)
        
        # Start with successful interaction
        response1 = agent.process_single_prompt("Say hello")
        assert response1 is not None
        history_length = len(agent.conversation_history)
        
        # Simulate error by temporarily changing to invalid model
        original_config = agent.config.copy()
        agent.config["model"] = "invalid-model-name"
        
        # This should fail
        with pytest.raises(ModelError):
            agent.process_single_prompt("This should fail")
        
        # Restore valid configuration
        agent.config = original_config
        
        # Should be able to continue conversation
        response2 = agent.process_single_prompt("Say goodbye")
        assert response2 is not None
        
        # Conversation history should have grown despite the error
        assert len(agent.conversation_history) > history_length