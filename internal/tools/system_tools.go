package tools

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
)

// SystemTool defines a system command to be exposed as an agent tool.
type SystemTool struct {
	Name    string
	Command string // The underlying Linux command (e.g., "mv", "rm")
	Description string
	ArgsSchema map[string]any // Optional: schema for arguments
}

// RegisteredSystemTools is a collection of system tools to be registered dynamically.
var RegisteredSystemTools = []SystemTool{
	{
		Name:    "mv",
		Command: "mv",
		Description: "Moves or renames files or directories.",
		ArgsSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"args": map[string]any{
					"type": "string",
					"description": "Arguments to pass to the mv command (e.g., 'old_file.txt new_file.txt')",
				},
			},
			"required": []string{"args"},
		},
	},
	{
		Name:    "rm",
		Command: "rm",
		Description: "Removes files or directories.",
		ArgsSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"args": map[string]any{
					"type": "string",
					"description": "Arguments to pass to the rm command (e.g., 'file_to_delete.txt' or '-r dir_to_delete/')",
				},
			},
			"required": []string{"args"},
		},
	},
	{
		Name:    "touch",
		Command: "touch",
		Description: "Updates the access and modification times of files, or creates them if they don't exist.",
		ArgsSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"args": map[string]any{
					"type": "string",
					"description": "Arguments to pass to the touch command (e.g., 'new_empty_file.txt')",
				},
			},
			"required": []string{"args"},
		},
	},
	{
		Name:    "wc",
		Command: "wc",
		Description: "Prints newline, word, and byte counts for each file.",
		ArgsSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"args": map[string]any{
					"type": "string",
					"description": "Arguments to pass to the wc command (e.g., '-l file.txt')",
				},
			},
			"required": []string{"args"},
		},
	},
	{
		Name:    "cp",
		Command: "cp",
		Description: "Copies files and directories.",
		ArgsSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"args": map[string]any{
					"type": "string",
					"description": "Arguments to pass to the cp command (e.g., 'source_file.txt destination_file.txt' or '-r source_dir/ destination_dir/')",
				},
			},
			"required": []string{"args"},
		},
	},
	{
		Name:    "mkdir",
		Command: "mkdir",
		Description: "Create directories.",
		ArgsSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"args": map[string]any{
					"type": "string",
					"description": "Arguments to pass to the mkdir command (e.g., 'new_dir' or '-p new_parent_dir/new_dir')",
				},
			},
			"required": []string{"args"},
		},
	},
	{
		Name:    "find",
		Command: "find",
		Description: "Searches for files and directories in a directory hierarchy.",
		ArgsSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"args": map[string]any{
					"type": "string",
					"description": "Arguments to pass to the find command (e.g., \". -name *.txt\")",
				},
			},
			"required": []string{"args"},
		},
	},
	{
		Name:    "task",
		Command: "task",
		Description: "Run taskfile commands.",
		ArgsSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"args": map[string]any{
					"type": "string",
					"description": "Arguments to pass to the task command (e.g., 'build' or '')",
				},
			},
			"required": []string{"args"},
		},
	},
	{
		Name:    "xc",
		Command: "xc",
		Description: "Run xc commands.",
		ArgsSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"args": map[string]any{
					"type": "string",
					"description": "Arguments to pass to the xc command (e.g., 'build' or '')",
				},
			},
			"required": []string{"args"},
		},
	},
}

// NewSystemTools creates a slice of Tool objects from RegisteredSystemTools.
func NewSystemTools() []Tool {
	var tools []Tool
	for _, sysTool := range RegisteredSystemTools {
		tool := Tool{
			Name:        sysTool.Name,
			Description: sysTool.Description,
			InputSchema: sysTool.ArgsSchema,
			Handler:     createSystemToolHandler(sysTool.Command),
		}
		tools = append(tools, tool)
	}
	return tools
}

// createSystemToolHandler returns a Handler function for a given system command.
func createSystemToolHandler(command string) func(input json.RawMessage) (string, error) {
	return func(input json.RawMessage) (string, error) {
		var args struct {
			Args string `json:"args"`
		}
		if err := json.Unmarshal(input, &args); err != nil {
			return "", fmt.Errorf("invalid input for %s: %w", command, err)
		}

        var allArgs []string
        if command == "xc" {
                allArgs = append(allArgs, "-no-tty") // Always add -no-tty for xc
        }
        allArgs = append(allArgs, splitArgs(args.Args)...) // Add user-provided arguments

        cmd := exec.Command(command, allArgs...) // Use the combined arguments

		var stdout, stderr bytes.Buffer
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr

		err := cmd.Run()
		if err != nil {
			return "", fmt.Errorf("error executing %s: %s (stderr: %s)", command, err, stderr.String())
		}

		if stderr.Len() > 0 {
			return stderr.String(), nil // Return stderr as output if present, even on success
		}
		return stdout.String(), nil
	}
}

// splitArgs splits a string of arguments into a slice,
// attempting to handle quoted arguments. This is a simplified split
// and might not cover all edge cases of shell argument parsing.
func splitArgs(s string) []string {
    if s == "" {
        return []string{}
    }
    // A basic split by space, not robust for complex shell arguments with quotes
    // This will need improvement if complex arguments are common.
    // For now, it's a direct split to match the basic passthrough requirement.
    var args []string
    inQuote := false
    currentArg := []rune{}
    for _, r := range s {
        if r == '"' {
            inQuote = !inQuote
        } else if r == ' ' && !inQuote {
            args = append(args, string(currentArg))
            currentArg = []rune{}
        } else {
            currentArg = append(currentArg, r)
        }
    }
    args = append(args, string(currentArg)) // Add the last argument
    return args
}
