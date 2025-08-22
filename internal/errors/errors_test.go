package errors

import (
	"errors"
	"fmt"
	"testing"
)

func TestValidationError(t *testing.T) {
	err := ValidationError{
		Field:   "testField",
		Message: "test message",
	}
	expected := "validation failed for testField: test message"
	if err.Error() != expected {
		t.Errorf("Expected error message '%s', but got '%s'", expected, err.Error())
	}

	wrappedErr := errors.New("wrapped error")
	errWithWrapped := ValidationError{
		Field:   "testField",
		Message: "test message",
		Wrapped: wrappedErr,
	}
	expectedWithWrapped := fmt.Sprintf("validation failed for testField: test message: %v", wrappedErr)
	if errWithWrapped.Error() != expectedWithWrapped {
		t.Errorf("Expected error message '%s', but got '%s'", expectedWithWrapped, errWithWrapped.Error())
	}
	if !errors.Is(errWithWrapped, wrappedErr) {
		t.Errorf("Expected error to be unwrappable, but it was not")
	}
}

func TestFileOperationError(t *testing.T) {
	wrappedErr := errors.New("wrapped error")
	err := FileOperationError{
		Operation: "read",
		Path:      "/test/path",
		Wrapped:   wrappedErr,
	}
	expected := fmt.Sprintf("file read failed for /test/path: %v", wrappedErr)
	if err.Error() != expected {
		t.Errorf("Expected error message '%s', but got '%s'", expected, err.Error())
	}
	if !errors.Is(err, wrappedErr) {
		t.Errorf("Expected error to be unwrappable, but it was not")
	}
}

func TestEditError(t *testing.T) {
	err := NewEditError(EditErrorFileNotFound, "/test/path", "details", nil)
	expected := "EDIT_FILE_NOT_FOUND: details (/test/path)"
	if err.Error() != expected {
		t.Errorf("Expected error message '%s', but got '%s'", expected, err.Error())
	}

	wrappedErr := errors.New("wrapped error")
	errWithWrapped := NewEditError(EditErrorFileReadError, "/test/path", "details", wrappedErr)
	expectedWithWrapped := fmt.Sprintf("EDIT_FILE_READ_ERROR: details (/test/path): %v", wrappedErr)
	if errWithWrapped.Error() != expectedWithWrapped {
		t.Errorf("Expected error message '%s', but got '%s'", expectedWithWrapped, errWithWrapped.Error())
	}
	if !errors.Is(errWithWrapped, wrappedErr) {
		t.Errorf("Expected error to be unwrappable, but it was not")
	}
}

func TestToolError(t *testing.T) {
	wrappedErr := errors.New("wrapped error")
	err := ToolError{
		Tool:    "testTool",
		Op:      "execute",
		Wrapped: wrappedErr,
	}
	expected := fmt.Sprintf("tool testTool: execute: %v", wrappedErr)
	if err.Error() != expected {
		t.Errorf("Expected error message '%s', but got '%s'", expected, err.Error())
	}
	if !errors.Is(err, wrappedErr) {
		t.Errorf("Expected error to be unwrappable, but it was not")
	}
}

func TestSessionError(t *testing.T) {
	wrappedErr := errors.New("wrapped error")
	err := SessionError{
		SessionID: "testSession",
		Op:        "load",
		Wrapped:   wrappedErr,
	}
	expected := fmt.Sprintf("session testSession: load: %v", wrappedErr)
	if err.Error() != expected {
		t.Errorf("Expected error message '%s', but got '%s'", expected, err.Error())
	}
	if !errors.Is(err, wrappedErr) {
		t.Errorf("Expected error to be unwrappable, but it was not")
	}
}
