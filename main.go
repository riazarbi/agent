package main

import (
	"bytes"
	"context"
	"embed"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/chzyer/readline"
	"github.com/openai/openai-go/v2"
	"github.com/openai/openai-go/v2/option"
	"github.com/invopop/jsonschema"
)

//go:embed templates/*
var templateFS embed.FS

// Core types
type Agent struct {
	client                *openai.Client
	getUserMessage       func() (string, bool)
	tools                []ToolDefinition
	baseURL              string
	rl                   *readline.Instance
	singleShot           bool
	transitionToInteractive bool
	prePrompts         []string
}

type ToolDefinition struct {
	Name        string                       `json:"name"`
	Description string                       `json:"description"`
	InputSchema openai.FunctionParameters   `json:"input_schema"`
	Function    func(input json.RawMessage) (string, error)
}

// Tool input types
type ReadFileInput struct {
	Path string `json:"path" jsonschema_description:"The relative path of a file in the working directory."`
}

type ListFilesInput struct {
	Path string `json:"path,omitempty" jsonschema_description:"Optional relative path to list files from. Defaults to current directory if not provided."`
}

type EditFileInput struct {
	Path   string `json:"path" jsonschema_description:"The path to the file"`
	OldStr string `json:"old_str" jsonschema_description:"Text to search for - must match exactly and must only have one match exactly"`
	NewStr string `json:"new_str" jsonschema_description:"Text to replace old_str with"`
}

type DeleteFileInput struct {
	Path string `json:"path" jsonschema_description:"The relative path of the file to delete"`
}

type GrepInput struct {
	Pattern string   `json:"pattern" jsonschema_description:"The search pattern to look for (literal or regex)"`
	Args    []string `json:"args,omitempty" jsonschema_description:"Optional ripgrep arguments (e.g. --ignore-case, --hidden)"`
}

// Tool schemas
var ReadFileInputSchema = GenerateSchema[ReadFileInput]()
var ListFilesInputSchema = GenerateSchema[ListFilesInput]()
var EditFileInputSchema = GenerateSchema[EditFileInput]()
var DeleteFileInputSchema = GenerateSchema[DeleteFileInput]()
var GrepInputSchema = GenerateSchema[GrepInput]()

// Tool definitions
var ReadFileDefinition = ToolDefinition{
	Name:        "read_file",
	Description: "Read the contents of a given relative file path. Use this when you want to see what's inside a file. Do not use this with directory names.",
	InputSchema: ReadFileInputSchema,
	Function:    ReadFile,
}

var ListFilesDefinition = ToolDefinition{
	Name:        "list_files",
	Description: "List files and directories at a given path. If no path is provided, lists files in the current directory.",
	InputSchema: ListFilesInputSchema,
	Function:    ListFiles,
}

var EditFileDefinition = ToolDefinition{
	Name: "edit_file",
	Description: `Make edits to a text file.

Replaces 'old_str' with 'new_str' in the given file. 'old_str' and 'new_str' MUST be different from each other.

If the file specified with path doesn't exist, it will be created.
`,
	InputSchema: EditFileInputSchema,
	Function:    EditFile,
}

var DeleteFileDefinition = ToolDefinition{
	Name:        "delete_file",
	Description: "Delete a file at the given relative path. Use with caution as this operation cannot be undone.",
	InputSchema: DeleteFileInputSchema,
	Function:    DeleteFile,
}

var GrepDefinition = ToolDefinition{
	Name:        "grep",
	Description: "Search for patterns in files using ripgrep. Supports both literal and regex patterns.",
	InputSchema: GrepInputSchema,
	Function:    Grep,
}

// Main function
func main() {
	// Parse command line flags
	promptFile := flag.String("f", "", "Path to prompt file for single-shot mode")
	flag.StringVar(promptFile, "prompt-file", "", "Path to prompt file for single-shot mode")
	continueChat := flag.Bool("continue", false, "Continue in interactive mode after processing prompt file")
	timeout := flag.Int("timeout", 60, "Timeout in seconds for non-interactive mode")
	initFlag := flag.Bool("init", false, "Initialize .agent directory")
	prePrompts := flag.String("preprompts", "", "Path to preprompts file (defaults to .agent/preprompts)")
	flag.Parse()

	// Validate flags
	if *continueChat && *promptFile == "" {
		fmt.Println("Error: --continue flag can only be used with --f/--prompt-file")
		os.Exit(1)
	}

	// Handle init command
	if *initFlag {
		if err := copyTemplates(); err != nil {
			fmt.Printf("Error initializing .agent directory: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Successfully initialized .agent directory")
		return
	}

	// Get environment variables with defaults
	baseURL := os.Getenv("AGENT_BASE_URL")
	if baseURL == "" {
		baseURL = "https://api.anthropic.com/v1/"
	}
	
	apiKey := os.Getenv("AGENT_API_KEY")
	if apiKey == "" {
		apiKey = os.Getenv("ANTHROPIC_API_KEY")
	}
	if apiKey == "" {
		fmt.Println("Error: AGENT_API_KEY or ANTHROPIC_API_KEY environment variable must be set")
		os.Exit(1)
	}

	client := openai.NewClient(
		option.WithAPIKey(apiKey),
		option.WithBaseURL(baseURL),
	)

	// Check for .agent directory and offer to create if missing
	if err := checkAndOfferAgentInit(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	tools := []ToolDefinition{ReadFileDefinition, ListFilesDefinition, EditFileDefinition, DeleteFileDefinition, GrepDefinition}
	
	var agent *Agent
	if *promptFile != "" {
		const maxPromptSize = 1024 * 1024 // 1MB max
		
		// Single-shot mode or initial prompt with continue
		fileInfo, err := os.Stat(*promptFile)
		if err != nil {
			fmt.Printf("Error accessing prompt file: %v\n", err)
			return
		}
		if fileInfo.Size() > maxPromptSize {
			fmt.Printf("Error: prompt file too large (max %d bytes)\n", maxPromptSize)
			return
		}
		
		content, err := os.ReadFile(*promptFile)
		if err != nil {
			fmt.Printf("Error reading prompt file: %v\n", err)
			return
		}
		
		promptContent := string(content)
		firstCall := true

		// Initialize readline if we're going to continue
		var rl *readline.Instance
		if *continueChat {
			rl, err = readline.New("")
			if err != nil {
				fmt.Printf("Error initializing readline: %v\n", err)
				return
			}
			defer rl.Close()
		}
		
		initialGetUserMessage := func() (string, bool) {
			if !firstCall {
				if *continueChat {
					// Switch to interactive mode
					rl.SetPrompt("\u001b[94mYou\u001b[0m: ")
					line, err := rl.Readline()
					if err != nil {
						if err == io.EOF {
							return "", false
						}
						fmt.Printf("Error reading input: %v\n", err)
						return "", false
					}
					return line, true
				}
				return "", false
			}
			firstCall = false
			return promptContent, true
		}
		
		prompts, err := getPrePrompts(*prePrompts)
		if err != nil {
			fmt.Printf("Error loading preprompts: %v\n", err)
			os.Exit(1)
		}
		agent = NewAgent(&client, initialGetUserMessage, tools, baseURL, rl, prompts)
		agent.singleShot = !*continueChat
		agent.transitionToInteractive = *continueChat
	} else {
		// Interactive mode
		rl, err := readline.New("")
		if err != nil {
			fmt.Printf("Error initializing readline: %v\n", err)
			return
		}
		defer rl.Close()

		getUserMessage := func() (string, bool) {
			rl.SetPrompt("\u001b[94mYou\u001b[0m: ")
			line, err := rl.Readline()
			if err != nil {
				if err == io.EOF {
					return "", false
				}
				fmt.Printf("Error reading input: %v\n", err)
				return "", false
			}
			return line, true
		}
		
		prompts, err := getPrePrompts(*prePrompts)
		if err != nil {
			fmt.Printf("Error loading preprompts: %v\n", err)
			os.Exit(1)
		}
		agent = NewAgent(&client, getUserMessage, tools, baseURL, rl, prompts)
		agent.singleShot = false
	}
	ctx := context.Background()
	if agent.singleShot {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, time.Duration(*timeout)*time.Second)
		defer cancel()
	}
	if err := agent.Run(ctx); err != nil {
		fmt.Printf("Error: %s\n", err.Error())
	}
}

// Constructor
func NewAgent(client *openai.Client, getUserMessage func() (string, bool), tools []ToolDefinition, baseURL string, rl *readline.Instance, prePrompts []string) *Agent {
	return &Agent{
		client:         client,
		getUserMessage: getUserMessage,
		tools:          tools,
		baseURL:        baseURL,
		rl:             rl,
		prePrompts:   prePrompts,
	}
}

// Agent methods
func (a *Agent) Run(ctx context.Context) error {
	conversation := []openai.ChatCompletionMessageParamUnion{}

	// Add preprompts as user messages in order
	for _, prePrompt := range a.prePrompts {
		if prePrompt != "" {
			conversation = append(conversation, openai.UserMessage(prePrompt))
		}
	}

	if !a.singleShot || a.transitionToInteractive {
		fmt.Printf("Chat with Agent at %s (use 'ctrl-c' to quit)\n", a.baseURL)
	}

	readUserInput := true
	for {
		if readUserInput {
			if len(conversation) > 0 {
				fmt.Println() // Add spacing before user input (except first time)
			}
			userInput, ok := a.getUserMessage()
			if !ok {
				break
			}

			userMessage := openai.UserMessage(userInput)
			conversation = append(conversation, userMessage)
		}

		completion, err := a.runInference(ctx, conversation)
		if err != nil {
			return err
		}
		
		assistantMessage := completion.Choices[0].Message
		conversation = append(conversation, assistantMessage.ToParam())

		toolResults := []openai.ChatCompletionMessageParamUnion{}
		
		// Handle text content
		if assistantMessage.Content != "" {
			fmt.Printf("\u001b[93mAgent\u001b[0m: %s\n", assistantMessage.Content)
		}
		
		// Handle tool calls
		for _, toolCall := range assistantMessage.ToolCalls {
			result := a.executeTool(toolCall.ID, toolCall.Function.Name, json.RawMessage(toolCall.Function.Arguments))
			toolResults = append(toolResults, result)
		}
		
		if len(toolResults) == 0 {
			readUserInput = true
			// Log conversation state after text-only response
			a.logConversation(conversation)
			if a.singleShot && !a.transitionToInteractive {
				break
			}
			continue
		}
		
		readUserInput = false
		conversation = append(conversation, toolResults...)

		// Log conversation state after each cycle
		a.logConversation(conversation)
	}

	return nil
}

func (a *Agent) executeTool(id, name string, input json.RawMessage) openai.ChatCompletionMessageParamUnion {
	var toolDef ToolDefinition
	var found bool
	for _, tool := range a.tools {
		if tool.Name == name {
			toolDef = tool
			found = true
			break
		}
	}
	if !found {
		return openai.ToolMessage("tool not found", id)
	}

	fmt.Printf("\u001b[92mTool\u001b[0m: %s(%s)\n", name, input)
	response, err := toolDef.Function(input)
	if err != nil {
		return openai.ToolMessage(err.Error(), id)
	}
	return openai.ToolMessage(response, id)
}

func (a *Agent) runInference(ctx context.Context, conversation []openai.ChatCompletionMessageParamUnion) (*openai.ChatCompletion, error) {
	openaiTools := []openai.ChatCompletionToolUnionParam{}
	for _, tool := range a.tools {
		openaiTools = append(openaiTools, openai.ChatCompletionFunctionTool(openai.FunctionDefinitionParam{
			Name:        tool.Name,
			Description: openai.String(tool.Description),
			Parameters:  tool.InputSchema,
		}))
	}

	completion, err := a.client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Model:     "claude-3-5-sonnet-20241022",
		MaxTokens: openai.Int(1024),
		Messages:  conversation,
		Tools:     openaiTools,
	})
	return completion, err
}

func (a *Agent) logConversation(conversation []openai.ChatCompletionMessageParamUnion) {
	logFile := os.Getenv("LOG_FILE")
	if logFile == "" {
		logFile = ".agent/agent.log"
		// Ensure .agent directory exists
		if err := os.MkdirAll(".agent", 0755); err != nil {
			fmt.Printf("Warning: Failed to create .agent directory: %v\n", err)
			return
		}
	}
	
	data, err := json.MarshalIndent(conversation, "", "  ")
	if err != nil {
		fmt.Printf("Warning: Failed to marshal conversation for logging: %v\n", err)
		return
	}

	err = os.WriteFile(logFile, data, 0644)
	if err != nil {
		fmt.Printf("Warning: Failed to write conversation log: %v\n", err)
	}
}

// loadPrePrompts reads preprompts list file and loads all prompt contents
func loadPrePrompts(filePath string) ([]string, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read preprompts file %s: %w", filePath, err)
	}

	lines := strings.Split(string(content), "\n")
	var prompts []string

	for lineNum, line := range lines {
		line = strings.TrimSpace(line)
		
		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Validate file exists
		if _, err := os.Stat(line); os.IsNotExist(err) {
			return nil, fmt.Errorf("prompt file %s (line %d) does not exist", line, lineNum+1)
		}

		// Read prompt content
		promptContent, err := os.ReadFile(line)
		if err != nil {
			return nil, fmt.Errorf("failed to read prompt file %s (line %d): %w", line, lineNum+1, err)
		}

		prompts = append(prompts, string(promptContent))
	}

	return prompts, nil
}

// getPrePrompts reads preprompts from a preprompts file or uses override
func getPrePrompts(prepromptsFile string) ([]string, error) {
	if prepromptsFile == "" {
		prepromptsFile = ".agent/preprompts"
	}
	
	// Check if preprompts file exists
	if _, err := os.Stat(prepromptsFile); os.IsNotExist(err) {
		return nil, fmt.Errorf("preprompts file %s does not exist", prepromptsFile)
	}
	
	return loadPrePrompts(prepromptsFile)
}

// getPrePrompt reads the system prompt from file or uses override (legacy - for compatibility)
func getPrePrompt(override string) string {
	if override != "" {
		return override
	}
	
	content, err := os.ReadFile(".agent/prompts/system.md")
	if err != nil {
		return "" // No system prompt if file doesn't exist or can't be read
	}
	
	return string(content)
}

// Utility functions
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

	result := openai.FunctionParameters{
		"type":       "object",
		"properties": properties,
	}
	
	if len(schema.Required) > 0 {
		result["required"] = schema.Required
	}
	
	return result
}

func convertSchemaProperty(prop *jsonschema.Schema) map[string]any {
	result := make(map[string]any)
	
	if prop.Type != "" {
		result["type"] = prop.Type
	}
	if prop.Description != "" {
		result["description"] = prop.Description
	}
	
	return result
}

// Tool implementations
func ReadFile(input json.RawMessage) (string, error) {
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

		// Skip .git and .agent directories
		if info.IsDir() && (relPath == ".git" || strings.HasPrefix(relPath, ".git/") || relPath == ".agent" || strings.HasPrefix(relPath, ".agent/")) {
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

func EditFile(input json.RawMessage) (string, error) {
	editFileInput := EditFileInput{}
	err := json.Unmarshal(input, &editFileInput)
	if err != nil {
		return "", err
	}

	if editFileInput.Path == "" || editFileInput.OldStr == editFileInput.NewStr {
		return "", fmt.Errorf("invalid input parameters")
	}

	content, err := os.ReadFile(editFileInput.Path)
	if err != nil {
		if os.IsNotExist(err) && editFileInput.OldStr == "" {
			return createNewFile(editFileInput.Path, editFileInput.NewStr)
		}
		return "", err
	}

	oldContent := string(content)
	newContent := strings.ReplaceAll(oldContent, editFileInput.OldStr, editFileInput.NewStr)

	if oldContent == newContent && editFileInput.OldStr != "" {
		return "", fmt.Errorf("old_str not found in file")
	}

	err = os.WriteFile(editFileInput.Path, []byte(newContent), 0644)
	if err != nil {
		return "", err
	}

	return "OK", nil
}

func createNewFile(filePath, content string) (string, error) {
	dir := path.Dir(filePath)
	if dir != "." {
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			return "", fmt.Errorf("failed to create directory: %w", err)
		}
	}

	err := os.WriteFile(filePath, []byte(content), 0644)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %w", err)
	}

	return fmt.Sprintf("Successfully created file %s", filePath), nil
}

func DeleteFile(input json.RawMessage) (string, error) {
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

// copyTemplates copies the embedded templates to the .agent directory
func copyTemplates() error {
	// Check if .agent directory already exists
	if _, err := os.Stat(".agent"); err == nil {
		return fmt.Errorf(".agent directory already exists")
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("failed to check .agent directory: %w", err)
	}

	// Walk through the embedded template filesystem
	return fs.WalkDir(templateFS, "templates", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("failed to walk template directory: %w", err)
		}

		// Calculate relative path from templates/ root
		relPath, err := filepath.Rel("templates", path)
		if err != nil {
			return fmt.Errorf("failed to calculate relative path: %w", err)
		}

		// Skip the root templates directory itself
		if relPath == "." {
			return nil
		}

		// Target path in .agent directory
		targetPath := filepath.Join(".agent", relPath)

		if d.IsDir() {
			// Create directory
			if err := os.MkdirAll(targetPath, 0755); err != nil {
				return fmt.Errorf("failed to create directory %s: %w", targetPath, err)
			}
		} else {
			// Read file from embedded filesystem
			content, err := fs.ReadFile(templateFS, path)
			if err != nil {
				return fmt.Errorf("failed to read embedded file %s: %w", path, err)
			}

			// Ensure parent directory exists
			if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
				return fmt.Errorf("failed to create parent directory for %s: %w", targetPath, err)
			}

			// Write file to target location
			if err := os.WriteFile(targetPath, content, 0644); err != nil {
				return fmt.Errorf("failed to write file %s: %w", targetPath, err)
			}
		}

		return nil
	})
}

// checkAndOfferAgentInit checks if .agent directory exists and offers to create it
func checkAndOfferAgentInit() error {
	// Check if .agent directory exists
	if _, err := os.Stat(".agent"); err == nil {
		return nil // Directory exists, continue normally
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("failed to check .agent directory: %w", err)
	}

	// .agent directory doesn't exist, prompt user
	fmt.Print("The .agent directory does not exist. Would you like to create it now? (y/N): ")
	
	var response string
	fmt.Scanln(&response)
	
	response = strings.ToLower(strings.TrimSpace(response))
	if response == "y" || response == "yes" {
		if err := copyTemplates(); err != nil {
			return fmt.Errorf("failed to initialize .agent directory: %w", err)
		}
		fmt.Println("Successfully initialized .agent directory")
		return nil
	}
	
	fmt.Println("Continuing without .agent directory...")
	return nil
}

func Grep(input json.RawMessage) (string, error) {
	grepInput := GrepInput{}
	err := json.Unmarshal(input, &grepInput)
	if err != nil {
		return "", err
	}

	if grepInput.Pattern == "" {
		return "", fmt.Errorf("search pattern cannot be empty")
	}

	// Start with base command and pattern
	args := append([]string{grepInput.Pattern}, grepInput.Args...)
	
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