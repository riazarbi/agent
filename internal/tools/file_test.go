package tools

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"agent/test/helpers"
)

func TestFileOperations_ReadFile(t *testing.T) {
	tests := []struct {
		name        string
		setup       func(t *testing.T) string // Returns file path
		input       ReadFileInput
		wantContent string
		wantErr     bool
	}{
		{
			name: "read existing file",
			setup: func(t *testing.T) string {
				return helpers.TempFileWithName(t, "test.txt", "Hello, World!")
			},
			input:       ReadFileInput{},
			wantContent: "Hello, World!",
			wantErr:     false,
		},
		{
			name: "read non-existent file",
			setup: func(t *testing.T) string {
				dir := helpers.TempDir(t)
				return filepath.Join(dir, "nonexistent.txt")
			},
			input:   ReadFileInput{},
			wantErr: true,
		},
		{
			name: "read empty file",
			setup: func(t *testing.T) string {
				return helpers.TempFileWithName(t, "empty.txt", "")
			},
			input:       ReadFileInput{},
			wantContent: "",
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &FileOperations{}
			filePath := tt.setup(t)
			tt.input.Path = filePath

			inputBytes, err := json.Marshal(tt.input)
			require.NoError(t, err)

			result, err := f.ReadFile(json.RawMessage(inputBytes))

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.wantContent, result)
		})
	}
}

func TestFileOperations_ReadFile_InvalidJSON(t *testing.T) {
	f := &FileOperations{}

	// Test with invalid JSON - should panic as per implementation
	assert.Panics(t, func() {
		f.ReadFile(json.RawMessage("invalid json"))
	})
}

func TestFileOperations_EditFile(t *testing.T) {
	tests := []struct {
		name          string
		setup         func(t *testing.T) string // Returns file path
		input         EditFileInput
		wantErr       bool
		wantContains  []string // Expected substrings in result
		wantFileExist bool
		wantContent   string
	}{
		{
			name: "create new file with empty old_str",
			setup: func(t *testing.T) string {
				dir := helpers.TempDir(t)
				return filepath.Join(dir, "newfile.txt")
			},
			input: EditFileInput{
				OldStr: "",
				NewStr: "Hello, World!",
			},
			wantErr:       false,
			wantContains:  []string{"Created new file", "newfile.txt"},
			wantFileExist: true,
			wantContent:   "Hello, World!",
		},
		{
			name: "replace text in existing file",
			setup: func(t *testing.T) string {
				return helpers.TempFileWithName(t, "existing.txt", "Hello, World!")
			},
			input: EditFileInput{
				OldStr: "World",
				NewStr: "Go",
			},
			wantErr:       false,
			wantContains:  []string{"Successfully modified file", "1 replacement"},
			wantFileExist: true,
			wantContent:   "Hello, Go!",
		},
		{
			name: "empty file path",
			setup: func(t *testing.T) string {
				return ""
			},
			input: EditFileInput{
				Path:   "",
				OldStr: "test",
				NewStr: "new",
			},
			wantErr: true,
		},
		{
			name: "identical old_str and new_str",
			setup: func(t *testing.T) string {
				return helpers.TempFileWithName(t, "same.txt", "Hello, World!")
			},
			input: EditFileInput{
				OldStr: "World",
				NewStr: "World",
			},
			wantErr:      false,
			wantContains: []string{"No changes applied", "identical"},
		},
		{
			name: "string not found in file",
			setup: func(t *testing.T) string {
				return helpers.TempFileWithName(t, "notfound.txt", "Hello, World!")
			},
			input: EditFileInput{
				OldStr: "NotFound",
				NewStr: "Found",
			},
			wantErr: true,
		},
		{
			name: "expected replacements mismatch",
			setup: func(t *testing.T) string {
				return helpers.TempFileWithName(t, "mismatch.txt", "Hello, World! Hello, Universe!")
			},
			input: EditFileInput{
				OldStr:               "Hello",
				NewStr:               "Hi",
				ExpectedReplacements: intPtr(1), // But there are 2 occurrences
			},
			wantErr: true,
		},
		{
			name: "multiple replacements with correct expected count",
			setup: func(t *testing.T) string {
				return helpers.TempFileWithName(t, "multiple.txt", "foo bar foo baz foo")
			},
			input: EditFileInput{
				OldStr:               "foo",
				NewStr:               "qux",
				ExpectedReplacements: intPtr(3),
			},
			wantErr:       false,
			wantContains:  []string{"Successfully modified file", "3 replacement"},
			wantFileExist: true,
			wantContent:   "qux bar qux baz qux",
		},
		{
			name: "create file when trying to use empty old_str on existing file",
			setup: func(t *testing.T) string {
				return helpers.TempFileWithName(t, "exists.txt", "content")
			},
			input: EditFileInput{
				OldStr: "",
				NewStr: "new content",
			},
			wantErr: true, // Should error because file already exists
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &FileOperations{}
			filePath := tt.setup(t)
			if tt.input.Path == "" && filePath != "" {
				tt.input.Path = filePath
			}

			inputBytes, err := json.Marshal(tt.input)
			require.NoError(t, err)

			result, err := f.EditFile(json.RawMessage(inputBytes))

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)

			// Check result contains expected strings
			for _, contains := range tt.wantContains {
				assert.Contains(t, result, contains)
			}

			// Check file exists and has expected content
			if tt.wantFileExist {
				assert.FileExists(t, tt.input.Path)
				if tt.wantContent != "" {
					content, err := os.ReadFile(tt.input.Path)
					assert.NoError(t, err)
					assert.Equal(t, tt.wantContent, string(content))
				}
			}
		})
	}
}

func TestFileOperations_EditFile_InvalidJSON(t *testing.T) {
	f := &FileOperations{}

	result, err := f.EditFile(json.RawMessage("invalid json"))
	assert.Error(t, err)
	assert.Empty(t, result)
	assert.Contains(t, err.Error(), "invalid input")
}

func TestFileOperations_DeleteFile(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(t *testing.T) string
		input   DeleteFileInput
		wantErr bool
	}{
		{
			name: "delete existing file",
			setup: func(t *testing.T) string {
				return helpers.TempFileWithName(t, "delete.txt", "content")
			},
			input:   DeleteFileInput{},
			wantErr: false,
		},
		{
			name: "delete non-existent file",
			setup: func(t *testing.T) string {
				dir := helpers.TempDir(t)
				return filepath.Join(dir, "nonexistent.txt")
			},
			input:   DeleteFileInput{},
			wantErr: true,
		},
		{
			name: "empty path",
			setup: func(t *testing.T) string {
				return ""
			},
			input: DeleteFileInput{
				Path: "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &FileOperations{}
			filePath := tt.setup(t)
			if tt.input.Path == "" && filePath != "" {
				tt.input.Path = filePath
			}

			inputBytes, err := json.Marshal(tt.input)
			require.NoError(t, err)

			result, err := f.DeleteFile(json.RawMessage(inputBytes))

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Contains(t, result, "Successfully deleted file")

			// Verify file no longer exists
			_, err = os.Stat(tt.input.Path)
			assert.True(t, os.IsNotExist(err))
		})
	}
}

func TestFileOperations_DeleteFile_InvalidJSON(t *testing.T) {
	f := &FileOperations{}

	result, err := f.DeleteFile(json.RawMessage("invalid json"))
	assert.Error(t, err)
	assert.Empty(t, result)
}

func TestFileOperations_Head(t *testing.T) {
	// Create a test file with multiple lines
	content := "line1\nline2\nline3\nline4\nline5\nline6\nline7\nline8\nline9\nline10\nline11\nline12"
	testFile := helpers.TempFileWithName(t, "headtest.txt", content)

	tests := []struct {
		name        string
		input       HeadInput
		wantContain string
		wantErr     bool
	}{
		{
			name: "head with default lines",
			input: HeadInput{
				Args: testFile,
			},
			wantContain: "line1",
			wantErr:     false,
		},
		{
			name: "head with specific line count",
			input: HeadInput{
				Args: "-n 3 " + testFile,
			},
			wantContain: "line3",
			wantErr:     false,
		},
		{
			name: "head with non-existent file",
			input: HeadInput{
				Args: "/path/to/nonexistent/file.txt",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &FileOperations{}

			inputBytes, err := json.Marshal(tt.input)
			require.NoError(t, err)

			result, err := f.Head(json.RawMessage(inputBytes))

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			if tt.wantContain != "" {
				assert.Contains(t, result, tt.wantContain)
			}
		})
	}
}

func TestFileOperations_Tail(t *testing.T) {
	// Create a test file with multiple lines
	content := "line1\nline2\nline3\nline4\nline5\nline6\nline7\nline8\nline9\nline10\nline11\nline12"
	testFile := helpers.TempFileWithName(t, "tailtest.txt", content)

	tests := []struct {
		name        string
		input       TailInput
		wantContain string
		wantErr     bool
	}{
		{
			name: "tail with default lines",
			input: TailInput{
				Args: testFile,
			},
			wantContain: "line12",
			wantErr:     false,
		},
		{
			name: "tail with specific line count",
			input: TailInput{
				Args: "-n 3 " + testFile,
			},
			wantContain: "line10",
			wantErr:     false,
		},
		{
			name: "tail with non-existent file",
			input: TailInput{
				Args: "/path/to/nonexistent/file.txt",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &FileOperations{}

			inputBytes, err := json.Marshal(tt.input)
			require.NoError(t, err)

			result, err := f.Tail(json.RawMessage(inputBytes))

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			if tt.wantContain != "" {
				assert.Contains(t, result, tt.wantContain)
			}
		})
	}
}

func TestFileOperations_Cloc(t *testing.T) {
	// Create a simple test directory with some code files
	dir := helpers.TempDir(t)
	goFile := filepath.Join(dir, "test.go")
	err := os.WriteFile(goFile, []byte("package main\n\nfunc main() {\n\t// Comment\n\tprintln(\"hello\")\n}"), 0644)
	require.NoError(t, err)

	tests := []struct {
		name         string
		input        ClocInput
		wantContain  string
		wantErr      bool
		skipIfNoCloc bool
	}{
		{
			name: "cloc on directory",
			input: ClocInput{
				Args: dir,
			},
			wantContain:  "Go", // Should detect Go files
			wantErr:      false,
			skipIfNoCloc: true,
		},
		{
			name: "cloc with verbose flag",
			input: ClocInput{
				Args: "--include-lang=Go " + dir,
			},
			wantContain:  "Go",
			wantErr:      false,
			skipIfNoCloc: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.skipIfNoCloc {
				// Skip test if cloc is not available
				_, err := os.Stat("/usr/bin/cloc")
				if os.IsNotExist(err) {
					_, err = os.Stat("/usr/local/bin/cloc")
					if os.IsNotExist(err) {
						t.Skip("cloc command not available")
					}
				}
			}

			f := &FileOperations{}

			inputBytes, err := json.Marshal(tt.input)
			require.NoError(t, err)

			result, err := f.Cloc(json.RawMessage(inputBytes))

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			if tt.wantContain != "" {
				assert.Contains(t, result, tt.wantContain)
			}
		})
	}
}

func TestNewFileTools(t *testing.T) {
	tools := NewFileTools()

	// Verify we have the expected number of tools
	assert.Len(t, tools, 6)

	// Verify tool names
	expectedNames := []string{"read_file", "edit_file", "delete_file", "head", "tail", "cloc"}
	var actualNames []string
	for _, tool := range tools {
		actualNames = append(actualNames, tool.Name)
	}

	for _, expected := range expectedNames {
		assert.Contains(t, actualNames, expected)
	}

	// Verify all tools have handlers
	for _, tool := range tools {
		assert.NotNil(t, tool.Handler)
		assert.NotEmpty(t, tool.Description)
		assert.NotNil(t, tool.InputSchema)
	}
}

func TestGenerateSchema(t *testing.T) {
	// Test schema generation for ReadFileInput
	schema := GenerateSchema[ReadFileInput]()

	assert.Equal(t, "object", schema["type"])

	properties, ok := schema["properties"].(map[string]any)
	assert.True(t, ok)
	assert.Contains(t, properties, "path")

	pathProp, ok := properties["path"].(map[string]any)
	assert.True(t, ok)
	assert.Equal(t, "string", pathProp["type"])
}

func TestFileOperations_generateDiff(t *testing.T) {
	f := &FileOperations{}

	oldContent := "line1\nline2\nline3"
	newContent := "line1\nmodified line2\nline3"

	diff, err := f.generateDiff("test.txt", oldContent, newContent)
	assert.NoError(t, err)
	assert.Contains(t, diff, "modified line2")
	assert.Contains(t, diff, "test.txt")
}

func TestFileOperations_createNewFileAtomic(t *testing.T) {
	f := &FileOperations{}
	dir := helpers.TempDir(t)
	filePath := filepath.Join(dir, "subdir", "newfile.txt")
	content := "test content"

	result, err := f.createNewFileAtomic(filePath, content)
	assert.NoError(t, err)
	assert.Contains(t, result, "Created new file")

	// Verify file exists and has correct content
	assert.FileExists(t, filePath)
	actualContent, err := os.ReadFile(filePath)
	assert.NoError(t, err)
	assert.Equal(t, content, string(actualContent))
}

// Helper function to create int pointer
func intPtr(i int) *int {
	return &i
}
