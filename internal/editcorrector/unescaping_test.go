package editcorrector

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnescapeGoString(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expected    string
		description string
	}{
		{
			name:        "basic newline",
			input:       "Hello\\nWorld",
			expected:    "Hello\nWorld",
			description: "Should unescape a single newline character.",
		},
		{
			name:        "basic tab",
			input:       "Hello\\tWorld",
			expected:    "Hello\tWorld",
			description: "Should unescape a single tab character.",
		},
		{
			name:        "basic double quote",
			input:       "He said,\\\"Hello\\\"",
			expected:    "He said,\"Hello\"",
			description: "Should unescape a double quote character.",
		},
		{
			name:        "basic backslash",
			input:       "Path\\to\\file",
			expected:    "Path\\to\\file",
			description: "Should unescape a single backslash.",
		},
		{
			name:        "mixed escapes",
			input:       "Line1\\nLine2\\tLine3\\\"quoted\\\"",
			expected:    "Line1\nLine2\tLine3\"quoted\"",
			description: "Should handle multiple types of escapes in one string.",
		},
		{
			name:        "no escapes",
			input:       "Plain text string",
			expected:    "Plain text string",
			description: "Should return the same string if no escapes are present.",
		},
		{
			name:        "empty string",
			input:       "",
			expected:    "",
			description: "Should handle an empty string gracefully.",
		},
		{
			name:        "unicode escape",
			input:       "Euro\\u20AC",
			expected:    "Euro\u20AC", // strconv.Unquote handles \uXXXX
			description: "Should unescape a Unicode character.",
		},
		{
			name:        "already unescaped newline",
			input:       "Hello\nWorld",
			expected:    "Hello\nWorld",
			description: "Should not re-escape already unescaped newlines.",
		},
		{
			name:        "already unescaped tab",
			input:       "Hello\tWorld",
			expected:    "Hello\tWorld",
			description: "Should not re-escape already unescaped tabs.",
		},
		{
			name:        "already unescaped double quote",
			input:       "He said,\"Hello\"",
			expected:    "He said,\"Hello\"",
			description: "Should not re-escape already unescaped double quotes.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := UnescapeGoString(tt.input)
			assert.Equal(t, tt.expected, actual, tt.description)
		})
	}
}
