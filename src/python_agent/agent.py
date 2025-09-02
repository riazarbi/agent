"""Core agent logic with LiteLLM integration.

Keywords: agent, litellm, conversation, chat, model, AI
"""

from typing import Any

import litellm

from python_agent.bash_tool import BashTool
from python_agent.session import Session, SessionManager


class AgentError(Exception):
    """Base exception for agent-related errors."""

    pass


class ModelError(AgentError):
    """Raised when model API calls fail."""

    pass


class Agent:
    """Core agent with LiteLLM integration and conversation management."""

    def __init__(self, config: dict[str, Any]) -> None:
        """Initialize agent with configuration."""
        self.config = config
        self.conversation_history: list[dict[str, Any]] = []
        self.session_manager = SessionManager()
        self.current_session: Session | None = None
        
        # Initialize bash tool
        self.bash_tool = BashTool(
            confirmation_required=config.get("confirmation_required", False),
            timeout=config.get("timeout", 30)
        )
        self.bash_tool.set_enabled(config.get("tools_enabled", True))

        # Configure LiteLLM base URL only
        if config.get("base_url"):
            litellm.api_base = config["base_url"]

    def add_message(self, role: str, content: str) -> None:
        """Add message to conversation history."""
        self.conversation_history.append({"role": role, "content": content})
        if self.current_session:
            self.current_session.add_message(role, content)
    
    def start_new_session(self) -> None:
        """Start a new conversation session."""
        self.current_session = self.session_manager.create_session()
        
    def resume_from_session(self, session: Session) -> None:
        """Resume conversation from an existing session."""
        self.current_session = session
        self.conversation_history = session.messages.copy()
        
    def save_current_session(self) -> None:
        """Save current session to disk."""
        if self.current_session:
            self.session_manager.save_session(self.current_session)

    def chat_completion(self, messages: list[dict[str, Any]]) -> str:
        """Get chat completion from LiteLLM model."""
        try:
            response = litellm.completion(
                model=self.config["model"],
                messages=messages,
                max_tokens=self.config["max_tokens"],
                temperature=self.config["temperature"],
                timeout=self.config["timeout"],
            )
            content = response.choices[0].message.content
            return content if content is not None else ""
        except Exception as e:
            raise ModelError(f"Model API call failed: {e}") from e

    def process_single_prompt(self, prompt: str) -> str:
        """Process single prompt and return response."""
        messages = [{"role": "user", "content": prompt}]
        return self.chat_completion(messages)

    def interactive_loop(self) -> None:
        """Run interactive conversation loop."""
        verbose = self.config.get("verbose", False)

        # Start new session if not resuming
        if not self.current_session:
            self.start_new_session()
            if verbose:
                print(f"Started session: {self.current_session.session_id}")

        if not self.config.get("quiet", False):
            print("AI Coding Agent (type 'exit' to quit)")
            if self.current_session:
                print(f"Session: {self.current_session.session_id}\n")

        while True:
            try:
                user_input = input("User: ").strip()

                if user_input.lower() in ("exit", "quit", "bye"):
                    self.save_current_session()
                    if not self.config.get("quiet", False):
                        print(f"Session saved: {self.current_session.session_id if self.current_session else 'none'}")
                    break

                if not user_input:
                    continue

                self.add_message("user", user_input)

                if verbose:
                    print("Agent: Thinking...")

                response = self.chat_completion(self.conversation_history)
                self.add_message("assistant", response)

                print(f"Agent: {response}\n")
                
                # Save session after each interaction
                self.save_current_session()

            except KeyboardInterrupt:
                print("\nSaving session...")
                self.save_current_session()
                print("Goodbye!")
                break
            except ModelError as e:
                print(f"Error: {e}")
            except Exception as e:
                print(f"Unexpected error: {e}")
                if verbose:
                    print("Use --verbose for more details")
