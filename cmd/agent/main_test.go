package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGetPrePrompts(t *testing.T) {
	// Create temporary directory for test files
	tmpDir, err := os.MkdirTemp("", "preprompts_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Create test files
	validYamlFile := filepath.Join(tmpDir, "valid.yml")
	yamlContent := `- message:
    - text: "Test preprompt"`
	if err := os.WriteFile(validYamlFile, []byte(yamlContent), 0644); err != nil {
		t.Fatal(err)
	}

	validMdFile := filepath.Join(tmpDir, "valid.md")
	mdContent := "# Test preprompt"
	if err := os.WriteFile(validMdFile, []byte(mdContent), 0644); err != nil {
		t.Fatal(err)
	}

	validLegacyFile := filepath.Join(tmpDir, "legacy")
	legacyContent := validMdFile // Point to the markdown file
	if err := os.WriteFile(validLegacyFile, []byte(legacyContent), 0644); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name        string
		path        string
		expectError bool
		errorMsg    string
		expectEmpty bool
	}{
		{
			name:        "empty path uses default - file doesn't exist",
			path:        "",
			expectError: false,
			expectEmpty: true,
		},
		{
			name:        "valid YAML file",
			path:        validYamlFile,
			expectError: false,
			expectEmpty: false,
		},
		{
			name:        "valid markdown file",
			path:        validMdFile,
			expectError: false,
			expectEmpty: false,
		},
		{
			name:        "valid legacy file",
			path:        validLegacyFile,
			expectError: false,
			expectEmpty: false,
		},
		{
			name:        "nonexistent file specified",
			path:        filepath.Join(tmpDir, "nonexistent.yml"),
			expectError: true,
			errorMsg:    "preprompt file not found",
		},
		{
			name:        "file without extension",
			path:        filepath.Join(tmpDir, "nonexistent"),
			expectError: true,
			errorMsg:    "preprompt file not found",
		},
		{
			name:        "existing file without proper extension",
			path:        strings.TrimSuffix(validYamlFile, ".yml"),
			expectError: true,
			errorMsg:    "preprompt file not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			messages, err := getPrePrompts(tt.path)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
					return
				}
				if tt.errorMsg != "" && !strings.Contains(err.Error(), tt.errorMsg) {
					t.Errorf("expected error message to contain %q, got %q", tt.errorMsg, err.Error())
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if tt.expectEmpty && len(messages) != 0 {
				t.Errorf("expected empty messages, got %d messages", len(messages))
			}

			if !tt.expectEmpty && len(messages) == 0 {
				t.Errorf("expected non-empty messages, got empty")
			}
		})
	}
}

func TestGetPrePrompts_DefaultPath(t *testing.T) {
	// Test that default path is constructed correctly
	// This test doesn't depend on the file existing
	originalWd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	// Create a temporary directory and change to it
	tmpDir, err := os.MkdirTemp("", "default_path_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)
	defer os.Chdir(originalWd) // Restore working directory

	if err := os.Chdir(tmpDir); err != nil {
		t.Fatal(err)
	}

	// Test with empty path (should use default and not error when file doesn't exist)
	messages, err := getPrePrompts("")
	if err != nil {
		t.Errorf("expected no error with default path when file doesn't exist, got: %v", err)
	}

	if len(messages) != 0 {
		t.Errorf("expected empty messages when default file doesn't exist, got %d messages", len(messages))
	}
}

func TestLoadPromptFile_ErrorHandling(t *testing.T) {
	// Create temporary directory for test files
	tmpDir, err := os.MkdirTemp("", "prompt_file_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Create a valid markdown file
	validMdFile := filepath.Join(tmpDir, "valid.md")
	if err := os.WriteFile(validMdFile, []byte("# Test"), 0644); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name        string
		path        string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "valid markdown file",
			path:        validMdFile,
			expectError: false,
		},
		{
			name:        "nonexistent file",
			path:        filepath.Join(tmpDir, "nonexistent.md"),
			expectError: true,
			errorMsg:    "reading markdown file",
		},
		{
			name:        "unsupported extension",
			path:        filepath.Join(tmpDir, "test.txt"),
			expectError: true,
			errorMsg:    "unsupported file extension",
		},
	}

	// Create the unsupported extension file
	txtFile := filepath.Join(tmpDir, "test.txt")
	if err := os.WriteFile(txtFile, []byte("test content"), 0644); err != nil {
		t.Fatal(err)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			content, err := loadPromptFile(tt.path)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
					return
				}
				if tt.errorMsg != "" && !strings.Contains(err.Error(), tt.errorMsg) {
					t.Errorf("expected error message to contain %q, got %q", tt.errorMsg, err.Error())
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if content == "" {
				t.Errorf("expected non-empty content")
			}
		})
	}
}
