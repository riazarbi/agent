package helpers

import (
	"testing"
	"time"

	"agent/internal/config"
)

// TestConfig creates a test configuration with sensible defaults for testing
func TestConfig(t *testing.T) *config.Config {
	t.Helper()

	return &config.Config{
		API: config.APIConfig{
			Key:     "test-api-key",
			BaseURL: "https://api.test.com/v1/",
			Timeout: 10 * time.Second,
		},
		Agent: config.AgentConfig{
			RequestDelay:  0,
			SingleShot:    false,
			PrePrompts:    "/tmp/test-preprompts",
			PromptFile:    "",
			ContinueChat:  false,
			Timeout:       60,
			InitFlag:      false,
			ResumeSession: "",
		},
		Session: config.SessionConfig{
			Dir: "/tmp/test-sessions",
		},
		Logging: config.LoggingConfig{
			File: "/tmp/test-logs/agent.log",
		},
	}
}

// TestConfigWithOverrides creates a test config with custom overrides
func TestConfigWithOverrides(t *testing.T, overrides map[string]interface{}) *config.Config {
	t.Helper()

	cfg := TestConfig(t)

	if val, ok := overrides["api_key"]; ok {
		cfg.API.Key = val.(string)
	}
	if val, ok := overrides["session_dir"]; ok {
		cfg.Session.Dir = val.(string)
	}
	if val, ok := overrides["single_shot"]; ok {
		cfg.Agent.SingleShot = val.(bool)
	}

	return cfg
}
