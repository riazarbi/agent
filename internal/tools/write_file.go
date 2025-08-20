package tools

import (
	"fmt"
	"os"
	"encoding/json"
)

// write_file implements the "write_file" tool.
type WriteFileTool struct{}

func (t *WriteFileTool) Name() string {
	return "write_file"
}

func (t *WriteFileTool) Description() string {
	return "Overwrites the entire content of a file or creates a new file. Fails if the path refers to an existing directory."
}

func (t *WriteFileTool) Args() map[string]interface{} {
	return map[string]interface{}{
		"path": map[string]string{
			"type": "string",
			"description": "The path to the file.",
		},
		"content": map[string]string{
			"type": "string",
			"description": "The content to write to the file.",
		},
	}
}

func (t *WriteFileTool) Execute(args map[string]interface{}) (string, error) {
	path, ok := args["path"].(string)
	if !ok || path == "" {
		return "", fmt.Errorf("missing or invalid 'path' argument")
	}

	content, ok := args["content"].(string)
	if !ok {
		// Treat missing content as empty string to allow creating empty files
		content = ""
	}

	fileInfo, err := os.Stat(path)
	if err == nil {
		if fileInfo.IsDir() {
			return "", fmt.Errorf("Is a directory: '%s'", path)
		}
	} else if !os.IsNotExist(err) {
		return "", fmt.Errorf("failed to stat file '%s': %w", path, err)
	}

	err = os.WriteFile(path, []byte(content), 0o644) // 0o644 gives read/write for owner, read for others
	if err != nil {
		return "", fmt.Errorf("failed to write to file '%s': %w", path, err)
	}

	return fmt.Sprintf("Successfully wrote to '%s'.", path), nil
}

// ExecuteJSON implements the Handler interface for the write_file tool.
func (t *WriteFileTool) ExecuteJSON(input json.RawMessage) (string, error) {
	var writeFileInput WriteFileInput
	if err := json.Unmarshal(input, &writeFileInput); err != nil {
		return "", fmt.Errorf("invalid input: %w", err)
	}

	args := map[string]interface{}{
		"path":    writeFileInput.Path,
		"content": writeFileInput.Content,
	}

	return t.Execute(args)
}