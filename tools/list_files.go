package tools

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

// ListFilesInput struct for input to the ListFiles function
type ListFilesInput struct {
	Path string `json:"path,omitempty" jsonschema_description:"Optional relative path to list files from. Defaults to current directory if not provided."`
}

// ListFiles function lists files and directories at a given path.
// If no path is provided, it lists files in the current directory.
// Keywords: list, files, directory, filesystem, ls
func ListFiles(input json.RawMessage) (string, error) {
	listFilesInput := ListFilesInput{}
	err := json.Unmarshal(input, &listFilesInput)
	if err != nil {
		panic(err)
	}

	dir := "."
	if listFilesInput.Path != "" {
		dir = listFilesInput.Path
	}

	var files []string
	err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
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
