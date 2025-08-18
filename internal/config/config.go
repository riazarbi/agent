package config

import (
	"fmt"
	"os"
	"time"
)

// Config holds all configuration for the agent
type Config struct {
	API     APIConfig
	Agent   AgentConfig
	Session SessionConfig
	Logging LoggingConfig
}

// APIConfig holds API-related configuration
type APIConfig struct {
	Key     string
	BaseURL string
	Timeout time.Duration
}

// AgentConfig holds agent behavior configuration
type AgentConfig struct {
	RequestDelay  time.Duration
	SingleShot    bool
	PrePrompts    string
	PromptFile    string
	ContinueChat  bool
	Timeout       int
	InitFlag      bool
	ResumeSession string
}

// SessionConfig holds session management configuration
type SessionConfig struct {
	Dir string
}

// LoggingConfig holds logging configuration
type LoggingConfig struct {
	File string
}

// Load loads configuration from environment variables and returns a Config struct
func Load() (*Config, error) {
	apiKey := getEnv("AGENT_API_KEY", "")
	if apiKey == "" {
		apiKey = getEnv("ANTHROPIC_API_KEY", "")
	}
	if apiKey == "" {
		return nil, fmt.Errorf("AGENT_API_KEY or ANTHROPIC_API_KEY environment variable must be set")
	}

	return &Config{
		API: APIConfig{
			Key:     apiKey,
			BaseURL: getEnv("AGENT_BASE_URL", "https://api.anthropic.com/v1/"),
			Timeout: 30 * time.Second, // Default timeout
		},
		Agent: AgentConfig{
			RequestDelay:  0, // Will be set from flag
			SingleShot:    false,
			PrePrompts:    getEnv("AGENT_PREPROMPTS", ".agent/prompts/preprompts"),
			ContinueChat:  false,
			Timeout:       60,
			InitFlag:      false,
			ResumeSession: "",
		},
		Session: SessionConfig{
			Dir: getEnv("AGENT_SESSION_DIR", ".agent/sessions"),
		},
		Logging: LoggingConfig{
			File: getEnv("LOG_FILE", ".agent/agent.log"),
		},
	}, nil
}

// getEnv gets an environment variable with a fallback default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}