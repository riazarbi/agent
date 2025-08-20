package tools

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings")

type CpInput struct {
	Source      string `json:"source" jsonschema_description:"The source file or directory."`
	Destination string `json:"destination" jsonschema_description:"The destination file or directory."`
	Recursive   bool   `json:"recursive,omitempty" jsonschema_description:"Whether to copy directories recursively. Defaults to false."`
}

func NewCpTool() Tool {
	cp := &cpTool{}
	return Tool{
		Name:        "cp",
		Description: "Copies files and directories. Supports recursive copy for directories.",
		InputSchema: GenerateSchema[CpInput](),
		Handler: func(input json.RawMessage) (string, error) {
			var args CpInput
			if err := json.Unmarshal(input, &args); err != nil {
				return "", fmt.Errorf("failed to unmarshal input: %w", err)
			}
			return cp.Execute(args)
		},
	}
}

// cpTool implements the cp command.
type cpTool struct{}

// Execute runs the cp command with the given arguments.
func (c *cpTool) Execute(args CpInput) (string, error) {
	cmdArgs := []string{}
	if args.Recursive {
		cmdArgs = append(cmdArgs, "-r")
	}
	cmdArgs = append(cmdArgs, args.Source, args.Destination)

	cmd := exec.Command("cp", cmdArgs...)
	output, err := cmd.CombinedOutput()

	if err != nil {
		// Attempt to parse standard cp error messages for better agent feedback
		errorMsg := strings.TrimSpace(string(output))
		if errorMsg == "" {
			errorMsg = err.Error()
		}
		return "", fmt.Errorf("cp failed: %s", errorMsg)
	}

	// Craft a more informative success message based on the operation
	message := fmt.Sprintf("Successfully copied '%s' to '%s'.", args.Source, args.Destination)
	if args.Recursive {
		message = fmt.Sprintf("Successfully copied '%s' to '%s' recursively.", args.Source, args.Destination)
		// Refine message if copying into an existing directory
		if strings.HasSuffix(args.Destination, "/") || strings.Contains(strings.TrimSuffix(args.Destination, "/"), "/") {
			// This is a heuristic, better to check if destination is an existing dir
			// but for now, assuming if it ends with / or has path components, it's a directory target
			// Or if it just existed as a directory.
			message = fmt.Sprintf("Successfully copied '%s' into '%s' recursively.", args.Source, args.Destination)
		}
	} else if strings.HasSuffix(args.Destination, "/") { // Copying file into existing directory
		fileName := args.Source[strings.LastIndex(args.Source, "/")+1:]
		message = fmt.Sprintf("Successfully copied '%s' to '%s%s'.", args.Source, args.Destination, fileName)
	}

	return fmt.Sprintf(`{"message": "%s"}`, message), nil
}
