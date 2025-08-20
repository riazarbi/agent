package tools

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
)

type MkdirInput struct {
	Path    string `json:"path" jsonschema_description:"The path to the directory to create."`
	Parents bool   `json:"parents,omitempty" jsonschema_description:"If true, create parent directories as needed. Defaults to false."`
}

func NewMkdirTool() Tool {
	mkdir := &mkdirTool{}
	return Tool{
		Name:        "mkdir",
		Description: "Create directories.",
		InputSchema: GenerateSchema[MkdirInput](),
		Handler: func(input json.RawMessage) (string, error) {
			var args MkdirInput
			if err := json.Unmarshal(input, &args); err != nil {
				return "", fmt.Errorf("failed to unmarshal input: %w", err)
			}
			return mkdir.Execute(args)
		},
	}
}

// mkdirTool implements the mkdir command.
type mkdirTool struct{}

// Execute runs the mkdir command with the given arguments.
func (m *mkdirTool) Execute(args MkdirInput) (string, error) {
	cmdArgs := []string{}
	if args.Parents {
		cmdArgs = append(cmdArgs, "-p")
	}
	cmdArgs = append(cmdArgs, args.Path)

	cmd := exec.Command("mkdir", cmdArgs...)
	output, err := cmd.CombinedOutput()

	if err != nil {
		errorMessage := strings.TrimSpace(string(output))
		if errorMessage == "" {
			errorMessage = err.Error()
		}
		// Attempt to parse standard mkdir error messages for better agent feedback
		if strings.Contains(errorMessage, "File exists") {
			return "", fmt.Errorf("File exists: '%s'", args.Path)
		}
		if strings.Contains(errorMessage, "No such file or directory") {
			return "", fmt.Errorf("No such file or directory: '%s'", args.Path)
		}
		return "", fmt.Errorf("mkdir failed: %s", errorMessage)
	}

	message := ""
	if args.Parents {
		message = fmt.Sprintf("Successfully ensured directory '%s' exists.", args.Path)
	} else {
		message = fmt.Sprintf("Successfully created directory '%s'.", args.Path)
	}

	return fmt.Sprintf(`{"message": "%s"}`, message), nil
}
