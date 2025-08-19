package tools

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
)

// GitOperations handles git-related tool operations
type GitOperations struct{}

// NewGitTools returns git operation tools
func NewGitTools() []Tool {
	ops := &GitOperations{}
	return []Tool{
		{
			Name:        "git_diff",
			Description: "Returns the output of 'git diff' showing all unstaged changes in the working directory. Use this when you need to see what files have been modified but not yet committed. Do not use this for staged/committed changes.",
			InputSchema: GenerateSchema[GitDiffInput](),
			Handler:     ops.GitDiff,
		},
		{
			Name:        "grep",
			Description: "Search for patterns in files using ripgrep. Supports both literal and regex patterns.",
			InputSchema: GenerateSchema[GrepInput](),
			Handler:     ops.Grep,
		},
		{
			Name:        "glob",
			Description: "Find files matching a glob pattern. Supports standard glob syntax for file discovery.",
			InputSchema: GenerateSchema[GlobInput](),
			Handler:     ops.Glob,
		},
	}
}

// Input type definitions
type GitDiffInput struct {
	// This tool takes no parameters
}

type GrepInput struct {
	Pattern string `json:"pattern" jsonschema_description:"The search pattern to look for (literal or regex)"`
	Args    string `json:"args,omitempty" jsonschema_description:"Optional ripgrep arguments as space-separated string (e.g. '--ignore-case --hidden')"`
}

type GlobInput struct {
	Pattern string `json:"pattern" jsonschema_description:"The glob pattern to match files against (e.g. *.go, **/*.md)"`
}

// GitDiff shows all unstaged changes in the working directory
func (g *GitOperations) GitDiff(input json.RawMessage) (string, error) {
	// Create git diff command
	cmd := exec.Command("git", "diff")

	// Capture both stdout and stderr
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// Execute command
	err := cmd.Run()

	// Check for errors
	if err != nil {
		// If there's stderr output, return that as the error
		if stderr.Len() > 0 {
			return "", fmt.Errorf(stderr.String())
		}
		return "", err
	}

	// If there's no output, it means there are no changes
	if stdout.Len() == 0 {
		return "No changes found in working directory", nil
	}

	return stdout.String(), nil
}

// Grep searches for patterns in files using ripgrep
func (g *GitOperations) Grep(input json.RawMessage) (string, error) {
	grepInput := GrepInput{}
	err := json.Unmarshal(input, &grepInput)
	if err != nil {
		return "", err
	}

	if grepInput.Pattern == "" {
		return "", fmt.Errorf("search pattern cannot be empty")
	}

	// Start with base command and pattern
	args := []string{grepInput.Pattern}

	// Parse space-separated args string if provided
	if grepInput.Args != "" {
		parsedArgs := strings.Fields(grepInput.Args)
		args = append(args, parsedArgs...)
	}

	cmd := exec.Command("rg", args...)

	// Capture both stdout and stderr
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err = cmd.Run()

	// rg exits with status 1 when no matches are found - this is not an error for us
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok && exitError.ExitCode() == 1 {
			return "No matches found", nil
		}
		// For any other error, return stderr output
		if stderr.Len() > 0 {
			return "", fmt.Errorf(stderr.String())
		}
		return "", err
	}

	return stdout.String(), nil
}

// Glob finds files matching a glob pattern
func (g *GitOperations) Glob(input json.RawMessage) (string, error) {
	globInput := GlobInput{}
	err := json.Unmarshal(input, &globInput)
	if err != nil {
		return "", err
	}

	if globInput.Pattern == "" {
		return "", fmt.Errorf("glob pattern cannot be empty")
	}

	matches, err := filepath.Glob(globInput.Pattern)
	if err != nil {
		return "", fmt.Errorf("invalid glob pattern: %w", err)
	}

	if len(matches) == 0 {
		return "No matches found", nil
	}

	result, err := json.Marshal(matches)
	if err != nil {
		return "", err
	}

	return string(result), nil
}