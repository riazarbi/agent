
package config

import (
	"os"
	"testing"
)

func TestLoad(t *testing.T) {
	// Unset API key env vars
	os.Unsetenv("AGENT_API_KEY")
	os.Unsetenv("ANTHROPIC_API_KEY")

	// Test for error when no API key is set
	_, err := Load()
	if err == nil {
		t.Errorf("Expected an error when no API key is set, but got nil")
	}

	// Test with AGENT_API_KEY
	os.Setenv("AGENT_API_KEY", "test-agent-key")
	cfg, err := Load()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if cfg.API.Key != "test-agent-key" {
		t.Errorf("Expected AGENT_API_KEY to be 'test-agent-key', but got '%s'", cfg.API.Key)
	}
	os.Unsetenv("AGENT_API_KEY")

	// Test with ANTHROPIC_API_KEY
	os.Setenv("ANTHROPIC_API_KEY", "test-anthropic-key")
	cfg, err = Load()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if cfg.API.Key != "test-anthropic-key" {
		t.Errorf("Expected ANTHROPIC_API_KEY to be 'test-anthropic-key', but got '%s'", cfg.API.Key)
	}
	os.Unsetenv("ANTHROPIC_API_KEY")

	// Test with both set (AGENT_API_KEY should take precedence)
	os.Setenv("AGENT_API_KEY", "test-agent-key")
	os.Setenv("ANTHROPIC_API_KEY", "test-anthropic-key")
	cfg, err = Load()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if cfg.API.Key != "test-agent-key" {
		t.Errorf("Expected AGENT_API_KEY to take precedence, but got '%s'", cfg.API.Key)
	}
	os.Unsetenv("AGENT_API_KEY")
	os.Unsetenv("ANTHROPIC_API_KEY")

	// Test default values
	os.Setenv("AGENT_API_KEY", "test-key")
	cfg, err = Load()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if cfg.API.BaseURL != "https://api.anthropic.com/v1/" {
		t.Errorf("Expected default BaseURL, but got '%s'", cfg.API.BaseURL)
	}
	if cfg.Session.Dir != ".agent/sessions" {
		t.Errorf("Expected default Session.Dir, but got '%s'", cfg.Session.Dir)
	}
	if cfg.Logging.File != ".agent/agent.log" {
		t.Errorf("Expected default Logging.File, but got '%s'", cfg.Logging.File)
	}
	os.Unsetenv("AGENT_API_KEY")

	// Test env var overrides for other values
	os.Setenv("AGENT_BASE_URL", "https://example.com")
	os.Setenv("AGENT_SESSION_DIR", "/tmp/sessions")
	os.Setenv("LOG_FILE", "/tmp/agent.log")
	os.Setenv("AGENT_API_KEY", "test-key")

	cfg, err = Load()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if cfg.API.BaseURL != "https://example.com" {
		t.Errorf("Expected AGENT_BASE_URL to be 'https://example.com', but got '%s'", cfg.API.BaseURL)
	}
	if cfg.Session.Dir != "/tmp/sessions" {
		t.Errorf("Expected AGENT_SESSION_DIR to be '/tmp/sessions', but got '%s'", cfg.Session.Dir)
	}
	if cfg.Logging.File != "/tmp/agent.log" {
		t.Errorf("Expected LOG_FILE to be '/tmp/agent.log', but got '%s'", cfg.Logging.File)
	}
	os.Unsetenv("AGENT_BASE_URL")
	os.Unsetenv("AGENT_SESSION_DIR")
	os.Unsetenv("LOG_FILE")
	os.Unsetenv("AGENT_API_KEY")

	// Test AGENT_MODEL env var
	os.Setenv("AGENT_MODEL", "test-model")
	os.Setenv("AGENT_API_KEY", "test-key")
	cfg, err = Load()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if cfg.API.Model != "test-model" {
		t.Errorf("Expected AGENT_MODEL to be 'test-model', but got '%s'", cfg.API.Model)
	}
	os.Unsetenv("AGENT_MODEL")
	os.Unsetenv("AGENT_API_KEY")
}

