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
				"AGENT_API_KEY":       "custom-api-key",
				"AGENT_BASE_URL":      "https://custom.api.com/",
				"AGENT_SESSION_DIR":   "/custom/sessions",
				"AGENT_PREPROMPTS":    "/custom/preprompts",
				"LOG_FILE":            "/custom/agent.log",
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
			name: "missing API key should return error",
			envVars: map[string]string{
				// No API key set
			},
			expectError: true,
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
		{
			name:         "returns empty string when env var is empty",
			key:          "TEST_VAR",
			defaultValue: "default",
			envValue:     "",
			expected:     "default",
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

func TestParseDuration(t *testing.T) {
	tests := []struct {
		name         string
		key          string
		defaultValue time.Duration
		envValue     string
		expected     time.Duration
	}{
		{
			name:         "returns default when env var not set",
			key:          "TEST_DURATION",
			defaultValue: 30 * time.Second,
			envValue:     "",
			expected:     30 * time.Second,
		},
		{
			name:         "parses valid duration",
			key:          "TEST_DURATION",
			defaultValue: 30 * time.Second,
			envValue:     "1m30s",
			expected:     90 * time.Second,
		},
		{
			name:         "returns default for invalid duration",
			key:          "TEST_DURATION",
			defaultValue: 30 * time.Second,
			envValue:     "invalid",
			expected:     30 * time.Second,
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

			result := parseDuration(tt.key, tt.defaultValue)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestParseBool(t *testing.T) {
	tests := []struct {
		name         string
		key          string
		defaultValue bool
		envValue     string
		expected     bool
	}{
		{
			name:         "returns default when env var not set",
			key:          "TEST_BOOL",
			defaultValue: true,
			envValue:     "",
			expected:     true,
		},
		{
			name:         "parses true",
			key:          "TEST_BOOL",
			defaultValue: false,
			envValue:     "true",
			expected:     true,
		},
		{
			name:         "parses false",
			key:          "TEST_BOOL",
			defaultValue: true,
			envValue:     "false",
			expected:     false,
		},
		{
			name:         "returns default for invalid bool",
			key:          "TEST_BOOL",
			defaultValue: true,
			envValue:     "invalid",
			expected:     true,
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

			result := parseBool(tt.key, tt.defaultValue)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// clearEnvVars clears all agent-related environment variables
func clearEnvVars() {
	envVars := []string{
		"AGENT_API_KEY",
		"AGENT_BASE_URL", 
		"AGENT_TIMEOUT",
		"AGENT_REQUEST_DELAY",
		"AGENT_SESSION_DIR",
		"AGENT_LOGGING_DIR",
		"AGENT_LOGGING_ENABLED",
	}

	for _, envVar := range envVars {
		os.Unsetenv(envVar)
	}
}