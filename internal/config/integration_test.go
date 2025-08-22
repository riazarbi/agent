package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadPromptFile_Integration(t *testing.T) {
	// Create temporary directory
	tmpDir, err := os.MkdirTemp("", "integration_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Test 1: Load a simple YAML file with command execution
	yamlFile := filepath.Join(tmpDir, "test.yml")
	yamlContent := `- message:
    - text: "Echo test: "
    - command: "echo 'hello world'"
- message:
    - text: "File test follows"
`
	if err := os.WriteFile(yamlFile, []byte(yamlContent), 0644); err != nil {
		t.Fatal(err)
	}

	messages, err := LoadPromptFile(yamlFile)
	if err != nil {
		t.Errorf("LoadPromptFile failed: %v", err)
		return
	}

	if len(messages) != 2 {
		t.Errorf("Expected 2 messages, got %d", len(messages))
		return
	}

	if messages[0] != "Echo test: hello world\n" {
		t.Errorf("Expected 'Echo test: hello world\\n', got %q", messages[0])
	}

	if messages[1] != "File test follows" {
		t.Errorf("Expected 'File test follows', got %q", messages[1])
	}
}

func TestLoadPromptFile_RecursiveInclusion(t *testing.T) {
	// Create temporary directory
	tmpDir, err := os.MkdirTemp("", "recursive_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Create an included file
	includedFile := filepath.Join(tmpDir, "included.yml")
	includedContent := `- message:
    - text: "From included file"`
	if err := os.WriteFile(includedFile, []byte(includedContent), 0644); err != nil {
		t.Fatal(err)
	}

	// Create main file that includes the other
	mainFile := filepath.Join(tmpDir, "main.yml")
	mainContent := `- message:
    - text: "Before include"
- file: "included.yml"
- message:
    - text: "After include"`
	if err := os.WriteFile(mainFile, []byte(mainContent), 0644); err != nil {
		t.Fatal(err)
	}

	messages, err := LoadPromptFile(mainFile)
	if err != nil {
		t.Errorf("LoadPromptFile failed: %v", err)
		return
	}

	expected := []string{
		"Before include",
		"From included file",
		"After include",
	}

	if len(messages) != len(expected) {
		t.Errorf("Expected %d messages, got %d", len(expected), len(messages))
		return
	}

	for i, exp := range expected {
		if messages[i] != exp {
			t.Errorf("Message %d: expected %q, got %q", i, exp, messages[i])
		}
	}
}

func TestLoadPromptFile_CircularDependency(t *testing.T) {
	// Create temporary directory
	tmpDir, err := os.MkdirTemp("", "circular_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Create two files that include each other
	file1 := filepath.Join(tmpDir, "file1.yml")
	file1Content := `- message:
    - text: "From file1"
- file: "file2.yml"`
	if err := os.WriteFile(file1, []byte(file1Content), 0644); err != nil {
		t.Fatal(err)
	}

	file2 := filepath.Join(tmpDir, "file2.yml")
	file2Content := `- message:
    - text: "From file2"
- file: "file1.yml"`
	if err := os.WriteFile(file2, []byte(file2Content), 0644); err != nil {
		t.Fatal(err)
	}

	_, err = LoadPromptFile(file1)
	if err == nil {
		t.Errorf("Expected error for circular dependency, but got none")
	}
}

func TestLoadPromptFile_FileWithContent(t *testing.T) {
	// Create temporary directory
	tmpDir, err := os.MkdirTemp("", "file_content_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Create content file
	contentFile := filepath.Join(tmpDir, "content.txt")
	if err := os.WriteFile(contentFile, []byte("File content here"), 0644); err != nil {
		t.Fatal(err)
	}

	// Create YAML file that includes the content
	yamlFile := filepath.Join(tmpDir, "test.yml")
	yamlContent := `- message:
    - text: "Content: "
    - file: "content.txt"`
	if err := os.WriteFile(yamlFile, []byte(yamlContent), 0644); err != nil {
		t.Fatal(err)
	}

	messages, err := LoadPromptFile(yamlFile)
	if err != nil {
		t.Errorf("LoadPromptFile failed: %v", err)
		return
	}

	if len(messages) != 1 {
		t.Errorf("Expected 1 message, got %d", len(messages))
		return
	}

	if messages[0] != "Content: File content here" {
		t.Errorf("Expected 'Content: File content here', got %q", messages[0])
	}
}