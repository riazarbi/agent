package errors

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidationError(t *testing.T) {
	tests := []struct {
		name     string
		field    string
		message  string
		expected string
	}{
		{
			name:     "basic validation error",
			field:    "email",
			message:  "is required",
			expected: "validation failed for email: is required",
		},
		{
			name:     "empty field validation error",
			field:    "",
			message:  "is invalid",
			expected: "validation failed for : is invalid",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidationError{
				Field:   tt.field,
				Message: tt.message,
			}

			assert.Equal(t, tt.expected, err.Error())
		})
	}
}

func TestFileOperationError(t *testing.T) {
	tests := []struct {
		name      string
		operation string
		path      string
		err       error
		expected  string
	}{
		{
			name:      "read file error",
			operation: "read",
			path:      "/path/to/file.txt",
			err:       errors.New("permission denied"),
			expected:  "file read failed for /path/to/file.txt: permission denied",
		},
		{
			name:      "write file error",
			operation: "write",
			path:      "/tmp/test.txt",
			err:       errors.New("disk full"),
			expected:  "file write failed for /tmp/test.txt: disk full",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := FileOperationError{
				Operation: tt.operation,
				Path:      tt.path,
				Wrapped:   tt.err,
			}

			assert.Equal(t, tt.expected, err.Error())
			assert.Equal(t, tt.err, err.Unwrap())
		})
	}
}

func TestEditError(t *testing.T) {
	tests := []struct {
		name      string
		errorType EditErrorType
		path      string
		details   string
		expected  string
	}{
		{
			name:      "no occurrence found error",
			errorType: EditErrorNoOccurrenceFound,
			path:      "test.go",
			details:   "string 'oldText' not found",
			expected:  "EDIT_NO_OCCURRENCE_FOUND: string 'oldText' not found (test.go)",
		},
		{
			name:      "file not found error",
			errorType: EditErrorFileNotFound,
			path:      "config.go",
			details:   "file does not exist",
			expected:  "EDIT_FILE_NOT_FOUND: file does not exist (config.go)",
		},
		{
			name:      "file read error",
			errorType: EditErrorFileReadError,
			path:      "main.go",
			details:   "permission denied",
			expected:  "EDIT_FILE_READ_ERROR: permission denied (main.go)",
		},
		{
			name:      "invalid path error",
			errorType: EditErrorInvalidPath,
			path:      "",
			details:   "path cannot be empty",
			expected:  "EDIT_INVALID_PATH: path cannot be empty ()",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := EditError{
				Type:    tt.errorType,
				Path:    tt.path,
				Details: tt.details,
			}

			assert.Equal(t, tt.expected, err.Error())
		})
	}
}

func TestToolError(t *testing.T) {
	tests := []struct {
		name      string
		tool      string
		operation string
		err       error
		expected  string
	}{
		{
			name:      "tool execution error",
			tool:      "read_file",
			operation: "execute",
			err:       errors.New("file not found"),
			expected:  "tool read_file: execute: file not found",
		},
		{
			name:      "tool validation error",
			tool:      "edit_file",
			operation: "validate",
			err:       ValidationError{Field: "path", Message: "is required"},
			expected:  "tool edit_file: validate: validation failed for path: is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ToolError{
				Tool:    tt.tool,
				Op:      tt.operation,
				Wrapped: tt.err,
			}

			assert.Equal(t, tt.expected, err.Error())
			assert.Equal(t, tt.err, err.Unwrap())
		})
	}
}

func TestSessionError(t *testing.T) {
	tests := []struct {
		name      string
		operation string
		id        string
		err       error
		expected  string
	}{
		{
			name:      "session load error",
			operation: "load",
			id:        "session-123",
			err:       errors.New("file not found"),
			expected:  "session session-123: load: file not found",
		},
		{
			name:      "session save error",
			operation: "save",
			id:        "session-456",
			err:       errors.New("permission denied"),
			expected:  "session session-456: save: permission denied",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := SessionError{
				SessionID: tt.id,
				Op:        tt.operation,
				Wrapped:   tt.err,
			}

			assert.Equal(t, tt.expected, err.Error())
			assert.Equal(t, tt.err, err.Unwrap())
		})
	}
}

func TestSentinelErrors(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected string
	}{
		{
			name:     "ErrToolNotFound",
			err:      ErrToolNotFound,
			expected: "tool not found",
		},
		{
			name:     "ErrSessionNotFound",
			err:      ErrSessionNotFound,
			expected: "session not found",
		},
		{
			name:     "ErrFileNotFound",
			err:      ErrFileNotFound,
			expected: "file not found",
		},
		{
			name:     "ErrInvalidInput",
			err:      ErrInvalidInput,
			expected: "invalid input",
		},
		{
			name:     "ErrInvalidPath",
			err:      ErrInvalidPath,
			expected: "file path cannot be empty",
		},
		{
			name:     "ErrFileAlreadyExists",
			err:      ErrFileAlreadyExists,
			expected: "file already exists",
		},
		{
			name:     "ErrNoOccurrenceFound",
			err:      ErrNoOccurrenceFound,
			expected: "could not find the string to replace",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.err.Error())
		})
	}
}

func TestErrorWrapping(t *testing.T) {
	t.Run("ValidationError with errors.Is", func(t *testing.T) {
		baseErr := ValidationError{Field: "test", Message: "invalid"}
		wrappedErr := ToolError{Tool: "test_tool", Op: "validate", Wrapped: baseErr}

		assert.True(t, errors.Is(wrappedErr, baseErr))
	})

	t.Run("FileOperationError with errors.Is", func(t *testing.T) {
		baseErr := errors.New("permission denied")
		wrappedErr := FileOperationError{Operation: "write", Path: "/test", Wrapped: baseErr}

		assert.True(t, errors.Is(wrappedErr, baseErr))
	})

	t.Run("Sentinel error with errors.Is", func(t *testing.T) {
		wrappedErr := ToolError{Tool: "test", Op: "execute", Wrapped: ErrFileNotFound}

		assert.True(t, errors.Is(wrappedErr, ErrFileNotFound))
	})
}
