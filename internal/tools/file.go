package tools

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/invopop/jsonschema"
	"github.com/openai/openai-go/v2"
	difflib "github.com/pmezard/go-difflib/difflib"

	"agent/internal/editcorrector"
	"agent/internal/errors"
)

// FileOperations handles file-related tool operations
type FileOperations struct{}

// NewFileTools returns file operation tools
func NewFileTools() []Tool {
	fileOps := &FileOperations{}
	return []Tool{
		{
			Name:        "read_file",
			Description: "Read the contents of a given relative file path. Use this when you want to see what's inside a file. Do not use this with directory names.",
			InputSchema: GenerateSchema[ReadFileInput](),
			Handler:     fileOps.ReadFile,
		},
		{
			Name: "edit_file",
			Description: `Make edits to a text file.

Replaces 'old_str' with 'new_str' in the given file. 'old_str' and 'new_str' MUST be different from each other.

If the file specified with path doesn't exist, it will be created.`,
			InputSchema: GenerateSchema[EditFileInput](),
			Handler:     fileOps.EditFile,
		},
		{
			Name:        "delete_file",
			Description: "Delete a file at the given relative path. Use with caution as this operation cannot be undone.",
			InputSchema: GenerateSchema[DeleteFileInput](),
			Handler:     fileOps.DeleteFile,
		},
		{
			Name:        "head",
			Description: "Show first N lines of a file (default 10 lines). Useful for quickly inspecting the beginning of files without reading the entire content.",
			InputSchema: GenerateSchema[HeadInput](),
			Handler:     fileOps.Head,
		},
		{
			Name:        "tail",
			Description: "Show last N lines of a file (default 10 lines). Useful for checking recent content or log file endings.",
			InputSchema: GenerateSchema[TailInput](),
			Handler:     fileOps.Tail,
		},
		{
			Name:        "cloc",
			Description: "Count lines of code with language breakdown and statistics. Useful for analyzing codebase size and composition.",
			InputSchema: GenerateSchema[ClocInput](),
			Handler:     fileOps.Cloc,
		},
		{
			Name:        "list_files",
			Description: "List files and directories at a given path. If no path is provided, lists files in the current directory.",
			InputSchema: GenerateSchema[ListFilesInput](),
			Handler:     fileOps.ListFiles,
		},
	}
}

// Input structs for file operations
type ReadFileInput struct {
	Path string `json:"path" jsonschema_description:"The relative path of a file in the working directory."`
}

type EditFileInput struct {
	Path                 string `json:"path" jsonschema_description:"The path to the file"`
	OldStr               string `json:"old_str" jsonschema_description:"Text to search for - must match exactly and must only have one match exactly"`
	NewStr               string `json:"new_str" jsonschema_description:"Text to replace old_str with"`
	ExpectedReplacements *int   `json:"expected_replacements,omitempty" jsonschema_description:"Optional: The expected number of replacements. If actual replacements differ, an error is returned."`
}

type AppendFileInput struct {
	Path    string `json:"path" jsonschema_description:"The path to the file"`
	Content string `json:"content" jsonschema_description:"Content to append or write to the file"`
}

type DeleteFileInput struct {
	Path string `json:"path" jsonschema_description:"The relative path of the file to delete."`
}

type HeadInput struct {
	Args string `json:"args,omitempty" jsonschema_description:"Optional head arguments as space-separated string (e.g. '-n 20 filename')"`
}

type TailInput struct {
	Args string `json:"args,omitempty" jsonschema_description:"Optional tail arguments as space-separated string (e.g. '-n 20 -f filename')"`
}

type ClocInput struct {
	Args string `json:"args,omitempty" jsonschema_description:"Optional cloc arguments as space-separated string (e.g. '--exclude-dir=.git path')"`
}

type ListFilesInput struct {
	Path string `json:"path,omitempty" jsonschema_description:"Optional relative path to list files from. Defaults to current directory if not provided."`
}

// File operation implementations
func (f *FileOperations) ReadFile(input json.RawMessage) (string, error) {
	readFileInput := ReadFileInput{}
	err := json.Unmarshal(input, &readFileInput)
	if err != nil {
		panic(err)
	}

	content, err := os.ReadFile(readFileInput.Path)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

func (f *FileOperations) EditFile(input json.RawMessage) (string, error) {
	editFileInput := EditFileInput{}
	err := json.Unmarshal(input, &editFileInput)
	if err != nil {
		return "", fmt.Errorf("invalid input: %w", err)
	}

	if editFileInput.Path == "" {
		return "", errors.NewEditError(errors.EditErrorInvalidPath, editFileInput.Path, "file path cannot be empty", nil)
	}

	// Programmatic unescaping (R1)
	correctedOldStr := editcorrector.UnescapeGoString(editFileInput.OldStr)
	correctedNewStr := editcorrector.UnescapeGoString(editFileInput.NewStr)

	// Read original content if file exists
	originalContentBytes, err := os.ReadFile(editFileInput.Path)
	var originalContent string
	if err == nil {
		originalContent = string(originalContentBytes)
	} else if !os.IsNotExist(err) {
		return "", errors.NewEditError(errors.EditErrorFileReadError, editFileInput.Path, "failed to read file", err)
	}

	// Scenario: old_str and new_str are identical (after correction)
	if correctedOldStr == correctedNewStr {
		// If old_str is empty, and file doesn't exist, this is a create scenario, so it's not "identical" in effect
		if correctedOldStr == "" && os.IsNotExist(err) {
			// This case is handled below by createNewFile, so not a no-op here.
		} else {
			return `{"message": "No changes applied, old_str and new_str are identical.", "actual_replacements": 0, "diff": ""}`, nil
		}
	}

	// Handle file creation with empty old_str (after correction)
	if correctedOldStr == "" {
		if err == nil { // File already exists
			return "", errors.NewEditError(errors.EditErrorCreateExistingFile, editFileInput.Path, "file already exists, cannot create using empty old_str", nil)
		} else if os.IsNotExist(err) { // File does not exist, proceed to create
			return f.createNewFileAtomic(editFileInput.Path, correctedNewStr)
		} else {
			return "", fmt.Errorf("EDIT_FILE_STAT_ERROR: failed to stat file %s: %w", editFileInput.Path, err)
		}
	}

	// For existing file edits
	if os.IsNotExist(err) {
		return "", fmt.Errorf("EDIT_FILE_NOT_FOUND: file not found: %s", editFileInput.Path)
	}

	// Perform replacements and count using corrected old_str
	count := strings.Count(originalContent, correctedOldStr)
	if count == 0 {
		return "", fmt.Errorf("EDIT_NO_OCCURRENCE_FOUND: could not find the string to replace: '%s' (attempted unescaped: '%s')", editFileInput.OldStr, correctedOldStr)
	}

	if editFileInput.ExpectedReplacements != nil && *editFileInput.ExpectedReplacements != count {
		return "", fmt.Errorf("EDIT_EXPECTED_OCCURRENCE_MISMATCH: expected %d occurrences but found %d for '%s' (attempted unescaped: '%s')", *editFileInput.ExpectedReplacements, count, editFileInput.OldStr, correctedOldStr)
	}

	newContent := strings.ReplaceAll(originalContent, correctedOldStr, correctedNewStr)

	// Generate diff
	diff, err := f.generateDiff(editFileInput.Path, originalContent, newContent)
	if err != nil {
		return "", fmt.Errorf("EDIT_DIFF_GENERATION_ERROR: failed to generate diff: %w", err)
	}

	// Atomic write
	tmpFile, err := os.CreateTemp(filepath.Dir(editFileInput.Path), "edit-temp-*.tmp")
	if err != nil {
		return "", fmt.Errorf("EDIT_TEMP_FILE_CREATE_ERROR: failed to create temporary file: %w", err)
	}
	tmpFilePath := tmpFile.Name()
	defer os.Remove(tmpFilePath) // Clean up temp file on exit

	_, err = tmpFile.WriteString(newContent)
	if err != nil {
		tmpFile.Close()
		return "", fmt.Errorf("EDIT_TEMP_FILE_WRITE_ERROR: failed to write to temporary file: %w", err)
	}
	tmpFile.Close()

	err = os.Rename(tmpFilePath, editFileInput.Path)
	if err != nil {
		return "", fmt.Errorf("EDIT_FILE_RENAME_ERROR: failed to rename temporary file to target: %w", err)
	}

	return fmt.Sprintf(`{"message": "Successfully modified file: %s (%d replacement(s)).", "actual_replacements": %d, "diff": %s}`, editFileInput.Path, count, count, strconv.Quote(diff)), nil
}

func (f *FileOperations) DeleteFile(input json.RawMessage) (string, error) {
	deleteFileInput := DeleteFileInput{}
	err := json.Unmarshal(input, &deleteFileInput)
	if err != nil {
		return "", err
	}

	if deleteFileInput.Path == "" {
		return "", fmt.Errorf("path cannot be empty")
	}

	err = os.Remove(deleteFileInput.Path)
	if err != nil {
		if os.IsNotExist(err) {
			return "", fmt.Errorf("file does not exist: %s", deleteFileInput.Path)
		}
		return "", fmt.Errorf("failed to delete file: %w", err)
	}

	return fmt.Sprintf("Successfully deleted file %s", deleteFileInput.Path), nil
}

// Helper methods for EditFile
func (f *FileOperations) createNewFileAtomic(filePath, content string) (string, error) {
	dir := path.Dir(filePath)
	if dir != "." {
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			return "", fmt.Errorf("EDIT_DIR_CREATE_ERROR: failed to create directory for %s: %w", filePath, err)
		}
	}

	tmpFile, err := os.CreateTemp(dir, "create-temp-*.tmp")
	if err != nil {
		return "", fmt.Errorf("EDIT_CREATE_TEMP_FILE_ERROR: failed to create temporary file for new file: %w", err)
	}
	tmpFilePath := tmpFile.Name()
	defer os.Remove(tmpFilePath) // Clean up temp file on exit

	_, err = tmpFile.WriteString(content)
	if err != nil {
		tmpFile.Close()
		return "", fmt.Errorf("EDIT_CREATE_TEMP_WRITE_ERROR: failed to write to temporary new file: %w", err)
	}
	tmpFile.Close()

	err = os.Rename(tmpFilePath, filePath)
	if err != nil {
		return "", fmt.Errorf("EDIT_CREATE_FILE_RENAME_ERROR: failed to rename temporary file to new file: %w", err)
	}

	diff, err := f.generateDiff(filePath, "", content)
	if err != nil {
		return "", fmt.Errorf("EDIT_CREATE_DIFF_GENERATION_ERROR: failed to generate diff for new file: %w", err)
	}

	return fmt.Sprintf(`{"message": "Created new file: %s with provided content.", "actual_replacements": 0, "diff": %s}`, filePath, strconv.Quote(diff)), nil
}

// generateDiff creates a unified diff string between old and new content
func (f *FileOperations) generateDiff(filePath, oldContent, newContent string) (string, error) {
	// Use difflib to generate a unified diff
	diff := difflib.UnifiedDiff{
		A:        difflib.SplitLines(oldContent),
		B:        difflib.SplitLines(newContent),
		FromFile: filePath,
		ToFile:   filePath,
		Context:  3, // Lines of context around changes
	}

	text, err := difflib.GetUnifiedDiffString(diff)
	if err != nil {
		return "", err
	}
	return text, nil
}

// GenerateSchema creates JSON schema for OpenAI function parameters
func GenerateSchema[T any]() openai.FunctionParameters {
	reflector := jsonschema.Reflector{
		AllowAdditionalProperties: false,
		DoNotReference:            true,
	}
	var v T

	schema := reflector.Reflect(v)

	// Convert jsonschema.Schema to openai.FunctionParameters format
	properties := make(map[string]any)
	if schema.Properties != nil {
		for pair := schema.Properties.Oldest(); pair != nil; pair = pair.Next() {
			properties[pair.Key] = convertSchemaProperty(pair.Value)
		}
	}

	required := make([]string, 0, len(schema.Required))
	for _, req := range schema.Required {
		required = append(required, req)
	}

	result := openai.FunctionParameters{
		"type":       "object",
		"properties": properties,
	}

	if len(required) > 0 {
		result["required"] = required
	}

	return result
}

// convertSchemaProperty converts a jsonschema property to the format expected by OpenAI
func convertSchemaProperty(prop *jsonschema.Schema) map[string]any {
	result := make(map[string]any)

	if prop.Type != "" {
		result["type"] = prop.Type
	}
	if prop.Description != "" {
		result["description"] = prop.Description
	}
	if prop.Items != nil {
		result["items"] = convertSchemaProperty(prop.Items)
	}
	if prop.Properties != nil {
		properties := make(map[string]any)
		for pair := prop.Properties.Oldest(); pair != nil; pair = pair.Next() {
			properties[pair.Key] = convertSchemaProperty(pair.Value)
		}
		result["properties"] = properties
	}
	if len(prop.Required) > 0 {
		required := make([]string, len(prop.Required))
		copy(required, prop.Required)
		result["required"] = required
	}

	return result
}

func (f *FileOperations) Head(input json.RawMessage) (string, error) {
	headInput := HeadInput{}
	err := json.Unmarshal(input, &headInput)
	if err != nil {
		return "", err
	}

	// Start with base command
	var args []string

	// Parse space-separated args string if provided
	if headInput.Args != "" {
		args = strings.Fields(headInput.Args)
	}

	cmd := exec.Command("head", args...)

	// Capture both stdout and stderr
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err = cmd.Run()

	// Check for errors
	if err != nil {
		// If there's stderr output, return that as the error
		if stderr.Len() > 0 {
			return "", fmt.Errorf(stderr.String())
		}
		return "", err
	}

	return stdout.String(), nil
}

func (f *FileOperations) Tail(input json.RawMessage) (string, error) {
	tailInput := TailInput{}
	err := json.Unmarshal(input, &tailInput)
	if err != nil {
		return "", err
	}

	// Start with base command
	var args []string

	// Parse space-separated args string if provided
	if tailInput.Args != "" {
		args = strings.Fields(tailInput.Args)
	}

	cmd := exec.Command("tail", args...)

	// Capture both stdout and stderr
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err = cmd.Run()

	// Check for errors
	if err != nil {
		// If there's stderr output, return that as the error
		if stderr.Len() > 0 {
			return "", fmt.Errorf(stderr.String())
		}
		return "", err
	}

	return stdout.String(), nil
}

func (f *FileOperations) Cloc(input json.RawMessage) (string, error) {
	clocInput := ClocInput{}
	err := json.Unmarshal(input, &clocInput)
	if err != nil {
		return "", err
	}

	// Start with base command
	var args []string

	// Parse space-separated args string if provided
	if clocInput.Args != "" {
		args = strings.Fields(clocInput.Args)
	}

	cmd := exec.Command("cloc", args...)

	// Capture both stdout and stderr
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err = cmd.Run()

	// Check for errors
	if err != nil {
		// If there's stderr output, return that as the error
		if stderr.Len() > 0 {
			return "", fmt.Errorf(stderr.String())
		}
		return "", err
	}

	return stdout.String(), nil
}

// ListFiles lists files and directories at a given path
func (f *FileOperations) ListFiles(input json.RawMessage) (string, error) {
	var listFilesInput ListFilesInput
	if err := json.Unmarshal(input, &listFilesInput); err != nil {
		return "", fmt.Errorf("invalid input: %v", err)
	}

	dir := "."
	if listFilesInput.Path != "" {
		dir = listFilesInput.Path
	}

	var files []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(dir, path)
		if err != nil {
			return err
		}

		// Skip .git and .agent/prompts directories
		if info.IsDir() && (relPath == ".git" || strings.HasPrefix(relPath, ".git/") || relPath == ".agent/prompts" || strings.HasPrefix(relPath, ".agent/prompts/") || relPath == ".agent/sessions" || strings.HasPrefix(relPath, ".agent/sessions/")) {
			return filepath.SkipDir
		}

		if relPath != "." {
			if info.IsDir() {
				files = append(files, relPath+"/")
			} else {
				files = append(files, relPath)
			}
		}
		return nil
	})

	if err != nil {
		return "", err
	}

	result, err := json.Marshal(files)
	if err != nil {
		return "", err
	}

	return string(result), nil
}

func (f *FileOperations) AppendFile(input json.RawMessage) (string, error) {
	appendFileInput := AppendFileInput{}
	err := json.Unmarshal(input, &appendFileInput)
	if err != nil {
		return "", fmt.Errorf("invalid input: %w", err)
	}

	if appendFileInput.Path == "" {
		return "", fmt.Errorf("path cannot be empty")
	}

	fileInfo, err := os.Stat(appendFileInput.Path)
	if err == nil && fileInfo.IsDir() {
		return "", fmt.Errorf("Is a directory: '%s'", appendFileInput.Path)
	}

	file, err := os.OpenFile(appendFileInput.Path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return "", fmt.Errorf("failed to open/create file: %w", err)
	}
	defer file.Close()

	_, err = file.WriteString(appendFileInput.Content)
	if err != nil {
		return "", fmt.Errorf("failed to append content to file: %w", err)
	}

	return fmt.Sprintf(`{"message": "Successfully appended to '%s'."}`, appendFileInput.Path), nil
}