package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestValidateContentPart(t *testing.T) {
	tests := []struct {
		name        string
		part        ContentPart
		expectError bool
	}{
		{
			name: "valid text part",
			part: ContentPart{Text: strPtr("hello")},
			expectError: false,
		},
		{
			name: "valid command part",
			part: ContentPart{Command: strPtr("echo hello")},
			expectError: false,
		},
		{
			name: "valid file part",
			part: ContentPart{File: strPtr("test.txt")},
			expectError: false,
		},
		{
			name: "empty part",
			part: ContentPart{},
			expectError: true,
		},
		{
			name: "multiple fields",
			part: ContentPart{Text: strPtr("hello"), Command: strPtr("echo")},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateContentPart(tt.part, 0, 0)
			if tt.expectError && err == nil {
				t.Errorf("expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

func TestValidatePromptItem(t *testing.T) {
	tests := []struct {
		name        string
		item        PromptItem
		expectError bool
	}{
		{
			name: "valid message item",
			item: PromptItem{
				MessageConstructor: MessageConstructor{
					Message: []ContentPart{
						{Text: strPtr("hello")},
					},
				},
			},
			expectError: false,
		},
		{
			name: "valid file item",
			item: PromptItem{
				File: strPtr("test.yml"),
			},
			expectError: false,
		},
		{
			name: "empty item",
			item: PromptItem{},
			expectError: true,
		},
		{
			name: "both message and file",
			item: PromptItem{
				MessageConstructor: MessageConstructor{
					Message: []ContentPart{
						{Text: strPtr("hello")},
					},
				},
				File: strPtr("test.yml"),
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validatePromptItem(tt.item, 0)
			if tt.expectError && err == nil {
				t.Errorf("expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

func TestCompileMessage(t *testing.T) {
	// Create a temporary directory for test files
	tmpDir, err := os.MkdirTemp("", "prompt_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Create a test file
	testFile := filepath.Join(tmpDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("file content"), 0644); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name     string
		parts    []ContentPart
		expected string
	}{
		{
			name: "single text part",
			parts: []ContentPart{
				{Text: strPtr("hello")},
			},
			expected: "hello",
		},
		{
			name: "multiple parts",
			parts: []ContentPart{
				{Text: strPtr("The git diff is:\n")},
				{Command: strPtr("echo 'diff output'")},
			},
			expected: "The git diff is:\ndiff output\n",
		},
		{
			name: "file content",
			parts: []ContentPart{
				{File: strPtr("test.txt")},
			},
			expected: "file content",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := compileMessage(tt.parts, tmpDir)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestLoadPromptFile_Markdown(t *testing.T) {
	// Create a temporary directory
	tmpDir, err := os.MkdirTemp("", "prompt_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Create a markdown file
	mdFile := filepath.Join(tmpDir, "test.md")
	mdContent := "# Hello World\nThis is a test."
	if err := os.WriteFile(mdFile, []byte(mdContent), 0644); err != nil {
		t.Fatal(err)
	}

	messages, err := LoadPromptFile(mdFile)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}

	if len(messages) != 1 {
		t.Errorf("expected 1 message, got %d", len(messages))
		return
	}

	if messages[0] != mdContent {
		t.Errorf("expected %q, got %q", mdContent, messages[0])
	}
}

func TestLoadPromptFile_YAML_Simple(t *testing.T) {
	// Create a temporary directory
	tmpDir, err := os.MkdirTemp("", "prompt_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Create a simple YAML file
	yamlFile := filepath.Join(tmpDir, "test.yml")
	yamlContent := `- message:
    - text: "Hello from YAML"
- message:
    - text: "Second message"
    - text: " with multiple parts"`
	
	if err := os.WriteFile(yamlFile, []byte(yamlContent), 0644); err != nil {
		t.Fatal(err)
	}

	messages, err := LoadPromptFile(yamlFile)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}

	if len(messages) != 2 {
		t.Errorf("expected 2 messages, got %d", len(messages))
		return
	}

	if messages[0] != "Hello from YAML" {
		t.Errorf("expected %q, got %q", "Hello from YAML", messages[0])
	}

	if messages[1] != "Second message with multiple parts" {
		t.Errorf("expected %q, got %q", "Second message with multiple parts", messages[1])
	}
}

func TestExecuteCommand(t *testing.T) {
	// Test a simple command
	result, err := executeCommand("echo 'hello world'")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}

	expected := "hello world\n"
	if result != expected {
		t.Errorf("expected %q, got %q", expected, result)
	}

	// Test command failure
	_, err = executeCommand("false")
	if err == nil {
		t.Errorf("expected error for failing command")
	}
}

// Helper function to create string pointers
func strPtr(s string) *string {
	return &s
}