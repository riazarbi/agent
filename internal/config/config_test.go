package config

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoad(t *testing.T) {
	tests := []struct {
		name        string
		envVars     map[string]string
		expectError bool
		expected    *Config
		expectedErrorMessage string
	}{
		{
			name: "default config with API key",
			envVars: map[string]string{
				"AGENT_API_KEY": "test-api-key",
			},
			expectError: false,
			expected: &Config{
				API: APIConfig{
					Key:     "test-api-key",
					BaseURL: "https://api.anthropic.com/v1/",
					Timeout: 30 * time.Second,
				},
				Agent: AgentConfig{
					RequestDelay:  0,
					SingleShot:    false,
					PrePrompts:    ".agent/prompts/preprompts",
					ContinueChat:  false,
					Timeout:       60,
					InitFlag:      false,
					ResumeSession: "",
				},
				Session: SessionConfig{
					Dir: ".agent/sessions",
				},
				Logging: LoggingConfig{
					File: ".agent/agent.log",
				},
			},
		},
		{
			name: "custom configuration from environment variables",
			envVars: map[string]string{
				"AGENT_API_KEY":     "custom-api-key",
				"AGENT_BASE_URL":    "https://custom.api.com/",
				"AGENT_SESSION_DIR": "/custom/sessions",
				"AGENT_PREPROMPTS":  "/custom/preprompts",
				"LOG_FILE":          "/custom/agent.log",
			},
			expectError: false,
			expected: &Config{
				API: APIConfig{
					Key:     "custom-api-key",
					BaseURL: "https://custom.api.com/",
					Timeout: 30 * time.Second,
				},
				Agent: AgentConfig{
					RequestDelay:  0,
					SingleShot:    false,
					PrePrompts:    "/custom/preprompts",
					ContinueChat:  false,
					Timeout:       60,
					InitFlag:      false,
					ResumeSession: "",
				},
				Session: SessionConfig{
					Dir: "/custom/sessions",
				},
				Logging: LoggingConfig{
					File: "/custom/agent.log",
				},
			},
		},
		{
			name:    "missing API key should return error",
			envVars: map[string]string{
				// No API key set
			},
			expectError: true,
			expected: nil,
			expectedErrorMessage: "AGENT_API_KEY, ANTHROPIC_API_KEY, or OPENAI_API_KEY environment variable must be set",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clear all environment variables first
			clearEnvVars()

			// Set test environment variables
			for key, value := range tt.envVars {
				t.Setenv(key, value)
			}

			cfg, err := Load()

			if tt.expectError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErrorMessage)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.expected, cfg)
		})
	}
}

func TestGetEnv(t *testing.T) {
	tests := []struct {
		name         string
		key          string
		defaultValue string
		envValue     string
		expected     string
	}{
		{
			name:         "returns default when env var not set",
			key:          "TEST_VAR",
			defaultValue: "default",
			envValue:     "",
			expected:     "default",
		},
		{
			name:         "returns env value when set",
			key:          "TEST_VAR",
			defaultValue: "default",
			envValue:     "env-value",
			expected:     "env-value",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clear environment variable
			os.Unsetenv(tt.key)

			// Set environment variable if provided
			if tt.envValue != "" {
				t.Setenv(tt.key, tt.envValue)
			}

			result := getEnv(tt.key, tt.defaultValue)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// clearEnvVars clears all agent-related environment variables
func clearEnvVars() {
	envVars := []string{
		"AGENT_API_KEY",
		"ANTHROPIC_API_KEY",
		"OPENAI_API_KEY",
		"AGENT_BASE_URL",
		"AGENT_SESSION_DIR",
		"AGENT_PREPROMPTS",
		"LOG_FILE",
	}

	for _, envVar := range envVars {
		os.Unsetenv(envVar)
	}
}
