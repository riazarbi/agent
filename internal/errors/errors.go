package errors

import (
	"errors"
	"fmt"
)

// Sentinel errors for common conditions
var (
	ErrToolNotFound      = errors.New("tool not found")
	ErrSessionNotFound   = errors.New("session not found")
	ErrFileNotFound      = errors.New("file not found")
	ErrInvalidInput      = errors.New("invalid input")
	ErrInvalidPath       = errors.New("file path cannot be empty")
	ErrFileAlreadyExists = errors.New("file already exists")
	ErrNoOccurrenceFound = errors.New("could not find the string to replace")
)

// ValidationError represents input validation failures
type ValidationError struct {
	Field   string
	Message string
	Wrapped error
}

func (e ValidationError) Error() string {
	if e.Wrapped != nil {
		return fmt.Sprintf("validation failed for %s: %s: %v", e.Field, e.Message, e.Wrapped)
	}
	return fmt.Sprintf("validation failed for %s: %s", e.Field, e.Message)
}

func (e ValidationError) Unwrap() error {
	return e.Wrapped
}

// FileOperationError represents file operation failures
type FileOperationError struct {
	Operation string
	Path      string
	Wrapped   error
}

func (e FileOperationError) Error() string {
	return fmt.Sprintf("file %s failed for %s: %v", e.Operation, e.Path, e.Wrapped)
}

func (e FileOperationError) Unwrap() error {
	return e.Wrapped
}

// EditError represents file editing operation failures
type EditError struct {
	Type    EditErrorType
	Path    string
	Details string
	Wrapped error
}

type EditErrorType string

const (
	EditErrorInvalidPath        EditErrorType = "EDIT_INVALID_PATH"
	EditErrorFileReadError      EditErrorType = "EDIT_FILE_READ_ERROR"
	EditErrorFileNotFound       EditErrorType = "EDIT_FILE_NOT_FOUND"
	EditErrorFileStatError      EditErrorType = "EDIT_FILE_STAT_ERROR"
	EditErrorCreateExistingFile EditErrorType = "ATTEMPT_TO_CREATE_EXISTING_FILE"
	EditErrorNoOccurrenceFound  EditErrorType = "EDIT_NO_OCCURRENCE_FOUND"
	EditErrorOccurrenceMismatch EditErrorType = "EDIT_EXPECTED_OCCURRENCE_MISMATCH"
	EditErrorDiffGeneration     EditErrorType = "EDIT_DIFF_GENERATION_ERROR"
	EditErrorTempFileCreate     EditErrorType = "EDIT_TEMP_FILE_CREATE_ERROR"
	EditErrorTempFileWrite      EditErrorType = "EDIT_TEMP_FILE_WRITE_ERROR"
	EditErrorFileRename         EditErrorType = "EDIT_FILE_RENAME_ERROR"
	EditErrorDirCreate          EditErrorType = "EDIT_DIR_CREATE_ERROR"
	EditErrorCreateTempFile     EditErrorType = "EDIT_CREATE_TEMP_FILE_ERROR"
	EditErrorCreateTempWrite    EditErrorType = "EDIT_CREATE_TEMP_WRITE_ERROR"
)

func (e EditError) Error() string {
	if e.Wrapped != nil {
		return fmt.Sprintf("%s: %s (%s): %v", e.Type, e.Details, e.Path, e.Wrapped)
	}
	return fmt.Sprintf("%s: %s (%s)", e.Type, e.Details, e.Path)
}

func (e EditError) Unwrap() error {
	return e.Wrapped
}

// NewEditError creates a new EditError
func NewEditError(errorType EditErrorType, path, details string, wrapped error) EditError {
	return EditError{
		Type:    errorType,
		Path:    path,
		Details: details,
		Wrapped: wrapped,
	}
}

// ToolError represents tool execution failures
type ToolError struct {
	Tool    string
	Op      string
	Wrapped error
}

func (e ToolError) Error() string {
	return fmt.Sprintf("tool %s: %s: %v", e.Tool, e.Op, e.Wrapped)
}

func (e ToolError) Unwrap() error {
	return e.Wrapped
}

// SessionError represents session management failures
type SessionError struct {
	SessionID string
	Op        string
	Wrapped   error
}

func (e SessionError) Error() string {
	return fmt.Sprintf("session %s: %s: %v", e.SessionID, e.Op, e.Wrapped)
}

func (e SessionError) Unwrap() error {
	return e.Wrapped
}
