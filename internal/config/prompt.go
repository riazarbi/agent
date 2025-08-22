package config

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// ContentPart represents a single part of a message's content.
type ContentPart struct {
	Text    *string `yaml:"text,omitempty"`
	Command *string `yaml:"command,omitempty"`
	File    *string `yaml:"file,omitempty"` // Path to a content file
}

// MessageConstructor represents a message to be constructed from several parts.
type MessageConstructor struct {
	Message []ContentPart `yaml:"message,omitempty"`
}

// PromptItem represents a top-level item in the prompt configuration.
// It can be a message constructor or a recursive file inclusion.
type PromptItem struct {
	MessageConstructor `yaml:",inline"`
	File               *string `yaml:"file,omitempty"` // Path to another .yml file
}

// PromptConfig represents the entire YAML configuration.
type PromptConfig []PromptItem

// LoadPromptFile loads and processes a prompt file based on its extension.
// For .md files, returns a single message with the file content.
// For .yml files, parses and compiles the YAML structure into messages.
func LoadPromptFile(path string) ([]string, error) {
	ext := strings.ToLower(filepath.Ext(path))
	
	switch ext {
	case ".md":
		content, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("reading markdown file %q: %w", path, err)
		}
		return []string{string(content)}, nil
		
	case ".yml", ".yaml":
		return LoadYAMLPromptFile(path)
		
	default:
		return nil, fmt.Errorf("unsupported file extension %q for file %q", ext, path)
	}
}

// LoadYAMLPromptFile loads and processes a YAML prompt file.
func LoadYAMLPromptFile(path string) ([]string, error) {
	// Track visited files to prevent circular dependencies
	visited := make(map[string]bool)
	return loadYAMLPromptFileRecursive(path, visited)
}

// loadYAMLPromptFileRecursive recursively processes YAML files with circular dependency detection.
func loadYAMLPromptFileRecursive(path string, visited map[string]bool) ([]string, error) {
	// Resolve to absolute path for circular dependency detection
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, fmt.Errorf("resolving absolute path for %q: %w", path, err)
	}
	
	if visited[absPath] {
		return nil, fmt.Errorf("circular dependency detected: %q", absPath)
	}
	visited[absPath] = true
	defer delete(visited, absPath) // Remove from visited when done processing this file
	
	// Read and parse YAML file
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading YAML file %q: %w", path, err)
	}
	
	var config PromptConfig
	if err := yaml.Unmarshal(content, &config); err != nil {
		return nil, fmt.Errorf("parsing YAML file %q: %w", path, err)
	}
	
	// Validate the configuration
	if err := validatePromptConfig(config); err != nil {
		return nil, fmt.Errorf("validating YAML file %q: %w", path, err)
	}
	
	// Compile messages from the configuration
	return compilePromptConfig(config, filepath.Dir(path), visited)
}

// validatePromptConfig validates the structure of a PromptConfig.
func validatePromptConfig(config PromptConfig) error {
	for i, item := range config {
		if err := validatePromptItem(item, i); err != nil {
			return err
		}
	}
	return nil
}

// validatePromptItem validates a single PromptItem.
func validatePromptItem(item PromptItem, index int) error {
	hasMessage := len(item.Message) > 0
	hasFile := item.File != nil && *item.File != ""
	
	if hasMessage && hasFile {
		return fmt.Errorf("item %d: cannot have both 'message' and 'file' keys", index)
	}
	if !hasMessage && !hasFile {
		return fmt.Errorf("item %d: must have either 'message' or 'file' key", index)
	}
	
	if hasMessage {
		for j, part := range item.Message {
			if err := validateContentPart(part, index, j); err != nil {
				return err
			}
		}
	}
	
	return nil
}

// validateContentPart validates a single ContentPart.
func validateContentPart(part ContentPart, itemIndex, partIndex int) error {
	count := 0
	if part.Text != nil && *part.Text != "" {
		count++
	}
	if part.Command != nil && *part.Command != "" {
		count++
	}
	if part.File != nil && *part.File != "" {
		count++
	}
	
	if count != 1 {
		return fmt.Errorf("item %d, part %d: must have exactly one of 'text', 'command', or 'file'", itemIndex, partIndex)
	}
	
	return nil
}

// compilePromptConfig compiles a PromptConfig into a list of messages.
func compilePromptConfig(config PromptConfig, baseDir string, visited map[string]bool) ([]string, error) {
	var messages []string
	
	for i, item := range config {
		if item.File != nil && *item.File != "" {
			// Recursive file inclusion
			filePath := *item.File
			if !filepath.IsAbs(filePath) {
				filePath = filepath.Join(baseDir, filePath)
			}
			
			subMessages, err := loadYAMLPromptFileRecursive(filePath, visited)
			if err != nil {
				return nil, fmt.Errorf("processing included file %q from item %d: %w", filePath, i, err)
			}
			messages = append(messages, subMessages...)
		} else {
			// Message constructor
			message, err := compileMessage(item.Message, baseDir)
			if err != nil {
				return nil, fmt.Errorf("compiling message for item %d: %w", i, err)
			}
			messages = append(messages, message)
		}
	}
	
	return messages, nil
}

// compileMessage compiles a list of ContentParts into a single message string.
func compileMessage(parts []ContentPart, baseDir string) (string, error) {
	var result strings.Builder
	
	for i, part := range parts {
		var content string
		var err error
		
		if part.Text != nil && *part.Text != "" {
			content = *part.Text
		} else if part.Command != nil && *part.Command != "" {
			content, err = executeCommand(*part.Command)
			if err != nil {
				return "", fmt.Errorf("executing command %q in part %d: %w", *part.Command, i, err)
			}
		} else if part.File != nil && *part.File != "" {
			filePath := *part.File
			if !filepath.IsAbs(filePath) {
				filePath = filepath.Join(baseDir, filePath)
			}
			
			fileContent, err := os.ReadFile(filePath)
			if err != nil {
				return "", fmt.Errorf("reading file %q in part %d: %w", filePath, i, err)
			}
			content = string(fileContent)
		}
		
		result.WriteString(content)
	}
	
	return result.String(), nil
}

// executeCommand executes a shell command and returns its combined stdout and stderr.
func executeCommand(command string) (string, error) {
	cmd := exec.Command("sh", "-c", command)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("command failed with exit code %v: %s", err, string(output))
	}
	return string(output), nil
}