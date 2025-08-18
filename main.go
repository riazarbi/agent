package main

import (
	"bytes"
	"context"
	"crypto/sha256"
	"embed"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/chzyer/readline"
	"github.com/openai/openai-go/v2"
	"github.com/openai/openai-go/v2/option"
	"github.com/invopop/jsonschema"
	htmltomarkdown "github.com/JohannesKaufmann/html-to-markdown/v2"
	difflib "github.com/pmezard/go-difflib/difflib"

	"agent/tools"
	"agent/internal/config"
	"agent/internal/editcorrector"
	"agent/internal/errors"
)

//go:embed templates/*
var templateFS embed.FS

// Session state for todo management and logging
var (
	sessionTodos = make(map[string][]TodoItem)
	sessionMutex sync.RWMutex
	sessionIDCounter = 1
	currentSessionID string // Set at startup: "2024-12-19-14-30-45"
	currentSessionDir string // Derived: ".agent/sessions/[currentSessionID]/"
	currentSessionManager *SessionManager // Global reference to current session manager
)

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
	requestDelay       time.Duration // New field for request delay
	config             *config.Config
}

type ToolDefinition struct {
	Name        string                       `json:"name"`
	Description string                       `json:"description"`
	InputSchema openai.FunctionParameters   `json:"input_schema"`
	Function    func(input json.RawMessage) (string, error)
}

// Session management types
type SessionManager struct {
	SessionID    string                                       `json:"session_id"`    // Format: "2024-12-19-14-30-45"  
	SessionDir   string                                       `json:"session_dir"`   // Path: ".agent/sessions/[timestamp]/"
	LogFile      *os.File                                     `json:"-"`             // File handle for agent.log
	TodosPath    string                                       `json:"todos_path"`    // Path to todos.json
	Conversation []openai.ChatCompletionMessageParamUnion     `json:"conversation"`  // Loaded from previous session or empty
}

type SessionTodos struct {
	Todos []TodoItem `json:"todos"`
}

// Tool input types
type ReadFileInput struct {
	Path string `json:"path" jsonschema_description:"The relative path of a file in the working directory."`
}



type EditFileInput struct {
	Path               string `json:"path" jsonschema_description:"The path to the file"`
	OldStr             string `json:"old_str" jsonschema_description:"Text to search for - must match exactly and must only have one match exactly"`
	NewStr             string `json:"new_str" jsonschema_description:"Text to replace old_str with"`
	ExpectedReplacements *int   `json:"expected_replacements,omitempty" jsonschema_description:"Optional: The expected number of replacements. If actual replacements differ, an error is returned."`
}

type DeleteFileInput struct {
	Path string `json:"path" jsonschema_description:"The relative path of the file to delete"`
}

type GrepInput struct {
	Pattern string `json:"pattern" jsonschema_description:"The search pattern to look for (literal or regex)"`
	Args    string `json:"args,omitempty" jsonschema_description:"Optional ripgrep arguments as space-separated string (e.g. '--ignore-case --hidden')"`
}

type GlobInput struct {
	Pattern string `json:"pattern" jsonschema_description:"The glob pattern to match files against (e.g. *.go, **/*.md)"`
}

type GitDiffInput struct {
	// This tool takes no parameters
}

type WebFetchInput struct {
	URL string `json:"url" jsonschema_description:"The URL to fetch content from (must start with http:// or https://)"`
}

type HtmlToMarkdownInput struct {
	Path string `json:"path" jsonschema_description:"Input HTML file path"`
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

// Todo management types
type TodoItem struct {
	ID       string `json:"id" jsonschema_description:"Unique identifier for the todo item"`
	Content  string `json:"content" jsonschema_description:"Brief description of the task"`
	Status   string `json:"status" jsonschema_description:"Current status: pending, in_progress, completed, cancelled"`
	Priority string `json:"priority" jsonschema_description:"Priority level: high, medium, low"`
}

type TodoWriteInput struct {
	TodosJSON string `json:"todos_json" jsonschema_description:"The updated todo list as JSON string containing array of TodoItem objects"`
}

type TodoReadInput struct {
	// No parameters needed - reads from current session
}

// Result from WebFetch operation
type CacheResult struct {
	Path        string `json:"path"`
	StatusCode  int    `json:"statusCode"`
	ContentType string `json:"contentType"`
}

// Tool schemas
var WebFetchInputSchema = GenerateSchema[WebFetchInput]()
var HtmlToMarkdownInputSchema = GenerateSchema[HtmlToMarkdownInput]()
var ReadFileInputSchema = GenerateSchema[ReadFileInput]()
var ListFilesInputSchema = GenerateSchema[tools.ListFilesInput]()

var EditFileInputSchema = GenerateSchema[EditFileInput]()
var DeleteFileInputSchema = GenerateSchema[DeleteFileInput]()
var GrepInputSchema = GenerateSchema[GrepInput]()
var GlobInputSchema = GenerateSchema[GlobInput]()
var GitDiffInputSchema = GenerateSchema[GitDiffInput]()
var HeadInputSchema = GenerateSchema[HeadInput]()
var TailInputSchema = GenerateSchema[TailInput]()
var ClocInputSchema = GenerateSchema[ClocInput]()
var TodoWriteInputSchema = GenerateSchema[TodoWriteInput]()
var TodoReadInputSchema = GenerateSchema[TodoReadInput]()

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
	Function:    tools.ListFiles,
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

var GlobDefinition = ToolDefinition{
	Name:        "glob",
	Description: "Find files matching a glob pattern. Supports standard glob syntax for file discovery.",
	InputSchema: GlobInputSchema,
	Function:    Glob,
}

var GitDiffDefinition = ToolDefinition{
	Name:        "git_diff",
	Description: "Returns the output of 'git diff' showing all unstaged changes in the working directory. Use this when you need to see what files have been modified but not yet committed. Do not use this for staged/committed changes.",
	InputSchema: GitDiffInputSchema,
	Function:    GitDiff,
}

// Supported content types for WebFetch
var allowedContentTypes = map[string]string{
	"text/plain":             ".txt",
	"text/html":             ".html",
	"text/xml":              ".xml",
	"application/json":      ".json",
	"application/xml":       ".xml",
	"application/xhtml+xml": ".html",
}

func generateFilename(inputURL string) (string, error) {
	parsedURL, err := url.Parse(inputURL)
	if err != nil {
		return "", fmt.Errorf("invalid URL: %v", err)
	}

	// Get base components
	host := strings.ToLower(parsedURL.Host)
	path := strings.ToLower(parsedURL.Path)

	// Clean the path
	path = strings.Trim(path, "/")
	path = strings.ReplaceAll(path, "/", "_")
	
	// Generate hash of full URL
	hasher := sha256.New()
	hasher.Write([]byte(inputURL))
	hash := hex.EncodeToString(hasher.Sum(nil))[:8]

	// Build base filename
	filename := fmt.Sprintf("%s_%s_%s", host, path, hash)

	// Replace invalid characters
	filename = strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '_' || r == '.' {
			return r
		}
		return '_'
	}, filename)

	return filename, nil
}

func isAllowedContentType(contentType string) (string, bool) {
	// Extract base content type
	base := strings.Split(contentType, ";")[0]
	base = strings.TrimSpace(base)
	
	// Check specific allowed types first
	if ext, ok := allowedContentTypes[base]; ok {
		return ext, true
	}
	
	// Check for other text/* types (default to .txt)
	if strings.HasPrefix(base, "text/") {
		return ".txt", true
	}
	
	return "", false
}

var WebFetchDefinition = ToolDefinition{
	Name:        "web_fetch",
	Description: "Download and cache web content locally. Accepts text/*, application/json, application/xml, and application/xhtml+xml content types. Returns path to cached file.",
	InputSchema: WebFetchInputSchema,
	Function:    WebFetch,
}

var HtmlToMarkdownDefinition = ToolDefinition{
	Name:        "html_to_markdown",
	Description: "Convert an HTML file to clean Markdown format, removing non-text content like images, videos, scripts, and styles. Saves output with same base filename but .md extension.",
	InputSchema: HtmlToMarkdownInputSchema,
	Function:    HtmlToMarkdown,
}

var HeadDefinition = ToolDefinition{
	Name:        "head",
	Description: "Show first N lines of a file (default 10 lines). Useful for quickly inspecting the beginning of files without reading the entire content.",
	InputSchema: HeadInputSchema,
	Function:    Head,
}

var TailDefinition = ToolDefinition{
	Name:        "tail",
	Description: "Show last N lines of a file (default 10 lines). Useful for checking recent content or log file endings.",
	InputSchema: TailInputSchema,
	Function:    Tail,
}

var ClocDefinition = ToolDefinition{
	Name:        "cloc",
	Description: "Count lines of code with language breakdown and statistics. Useful for analyzing codebase size and composition.",
	InputSchema: ClocInputSchema,
	Function:    Cloc,
}

var TodoWriteDefinition = ToolDefinition{
	Name:        "todowrite",
	Description: "Create and manage structured task lists for complex multi-step operations within the current session. Each todo requires: 'task' (title), 'content' (description), 'status' (pending/in_progress/completed), 'priority' (high/medium/low). Replaces entire todo list. Data is not persistent across sessions.",
	InputSchema: TodoWriteInputSchema,
	Function:    TodoWrite,
}

var TodoReadDefinition = ToolDefinition{
	Name:        "todoread",
	Description: "Read the current todo list from session state. Returns structured todos with auto-generated IDs, content, status, and priority. Data is session-only and not persistent across invocations.",
	InputSchema: TodoReadInputSchema,
	Function:    TodoRead,
}

// Main function
func main() {
	// Parse command line flags
	promptFile := flag.String("f", "", "Path to prompt file for single-shot mode")
	flag.StringVar(promptFile, "prompt-file", "", "Path to prompt file for single-shot mode")
	continueChat := flag.Bool("continue", false, "Continue in interactive mode after processing prompt file")
	timeout := flag.Int("timeout", 60, "Timeout in seconds for non-interactive mode")
	initFlag := flag.Bool("init", false, "Initialize .agent directory")
	prePrompts := flag.String("preprompts", "", "Path to preprompts file (defaults to .agent/prompts/preprompts)")
	resumeSession := flag.String("resume", "", "Resume a specific session by ID (YYYY-MM-DD-HH-MM-SS), or use 'list' to select interactively")
	requestDelay := flag.Duration("request-delay", 0, "Delay between API requests (e.g., 2s, 500ms)")
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

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("Error loading configuration: %v\n", err)
		os.Exit(1)
	}
	
	// Update config from flags
	cfg.Agent.PromptFile = *promptFile
	cfg.Agent.ContinueChat = *continueChat
	cfg.Agent.Timeout = *timeout
	cfg.Agent.InitFlag = *initFlag
	cfg.Agent.PrePrompts = *prePrompts
	cfg.Agent.ResumeSession = *resumeSession
	cfg.Agent.RequestDelay = *requestDelay

	client := openai.NewClient(
		option.WithAPIKey(cfg.API.Key),
		option.WithBaseURL(cfg.API.BaseURL),
	)

	// Check for .agent directory and offer to create if missing
	if err := checkAndOfferAgentInit(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	// Initialize session management
	_, err = initializeSession(cfg.Agent.ResumeSession)
	if err != nil {
		fmt.Printf("Error initializing session: %v\n", err)
		os.Exit(1)
	}

	// All tools now work with Gemini via OpenAI API - arrays converted to space-separated strings
	tools := []ToolDefinition{ReadFileDefinition, ListFilesDefinition, EditFileDefinition, DeleteFileDefinition, GrepDefinition, GlobDefinition, GitDiffDefinition, WebFetchDefinition, HtmlToMarkdownDefinition, HeadDefinition, TailDefinition, ClocDefinition, TodoWriteDefinition, TodoReadDefinition}
	
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
		
		prompts, err := getPrePrompts(cfg.Agent.PrePrompts)
		if err != nil {
			fmt.Printf("Error loading preprompts: %v\n", err)
			os.Exit(1)
		}
		agent = NewAgent(&client, initialGetUserMessage, tools, cfg.API.BaseURL, rl, prompts, cfg.Agent.RequestDelay, cfg)
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
		
		prompts, err := getPrePrompts(cfg.Agent.PrePrompts)
		if err != nil {
			fmt.Printf("Error loading preprompts: %v\n", err)
			os.Exit(1)
		}
		agent = NewAgent(&client, getUserMessage, tools, cfg.API.BaseURL, rl, prompts, cfg.Agent.RequestDelay, cfg)
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
func NewAgent(client *openai.Client, getUserMessage func() (string, bool), tools []ToolDefinition, baseURL string, rl *readline.Instance, prePrompts []string, requestDelay time.Duration, cfg *config.Config) *Agent {
	return &Agent{
		client:         client,
		getUserMessage: getUserMessage,
		tools:          tools,
		baseURL:        baseURL,
		rl:             rl,
		prePrompts:   prePrompts,
		requestDelay: requestDelay,
		config:       cfg,
	}
}

// Agent methods
func (a *Agent) Run(ctx context.Context) error {
	// Initialize conversation from session if available, otherwise empty
	var conversation []openai.ChatCompletionMessageParamUnion
	if currentSessionDir != "" && currentSessionManager != nil && len(currentSessionManager.Conversation) > 0 {
		// Resume with existing conversation
		conversation = currentSessionManager.Conversation
		fmt.Printf("Resuming session %s with %d previous messages\n", currentSessionID, len(conversation))
	} else {
		// Start with empty conversation
		conversation = []openai.ChatCompletionMessageParamUnion{}
	}

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

	// Use all tools now that problematic grep is excluded
	toolsToUse := openaiTools
	
	// Add delay if configured
	if a.requestDelay > 0 {
		time.Sleep(a.requestDelay)
	}
	
	completion, err := a.client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Model:     "gemini-2.5-flash",
		MaxTokens: openai.Int(4096),
		Messages:  conversation,
		Tools:     toolsToUse,
	})

	// Add verbose error handling for debugging
	if err != nil {
		fmt.Printf("Detailed error: %+v\n", err)
		fmt.Printf("Error type: %T\n", err)
		fmt.Printf("Error string: %s\n", err.Error())

		// Try to extract more details from different error types
		switch e := err.(type) {
		case *openai.Error:
			fmt.Printf("OpenAI API Error Code: %s\n", e.Code)
			fmt.Printf("OpenAI API Error Message: %s\n", e.Message)
			fmt.Printf("OpenAI API Error Type: %s\n", e.Type)
		default:
			fmt.Printf("Other error type: %T\n", e)
		}
	}

	return completion, err
}

func (a *Agent) logConversation(conversation []openai.ChatCompletionMessageParamUnion) {
	// Use session-specific logging if currentSessionDir is set
	var logFile string
	if currentSessionDir != "" {
		logFile = filepath.Join(currentSessionDir, "agent.log")
		// Ensure session directory exists
		if err := os.MkdirAll(currentSessionDir, 0755); err != nil {
			fmt.Printf("Warning: Failed to create session directory: %v\n", err)
			return
		}
	} else {
		// Fallback to global log for backward compatibility
		logFile = a.config.Logging.File
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
		prepromptsFile = ".agent/prompts/preprompts"
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
	
	content, err := os.ReadFile(".agent/prompts/system/system.md")
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



func EditFile(input json.RawMessage) (string, error) {
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
			return createNewFileAtomic(editFileInput.Path, correctedNewStr)
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
	diff, err := generateDiff(editFileInput.Path, originalContent, newContent)
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

func createNewFileAtomic(filePath, content string) (string, error) {
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

	diff, err := generateDiff(filePath, "", content)
	if err != nil {
		return "", fmt.Errorf("EDIT_CREATE_DIFF_GENERATION_ERROR: failed to generate diff for new file: %w", err)
	}

	return fmt.Sprintf(`{"message": "Created new file: %s with provided content.", "actual_replacements": 0, "diff": %s}`, filePath, strconv.Quote(diff)), nil
}

// generateDiff creates a unified diff string between old and new content
func generateDiff(filePath, oldContent, newContent string) (string, error) {
	// Use difflib to generate a unified diff
	diff := difflib.UnifiedDiff{
		A:       difflib.SplitLines(oldContent),
		B:       difflib.SplitLines(newContent),
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

func GitDiff(input json.RawMessage) (string, error) {
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

func WebFetch(input json.RawMessage) (string, error) {
    var webFetchInput WebFetchInput
    err := json.Unmarshal(input, &webFetchInput)
    if err != nil {
        return "", fmt.Errorf("invalid input: %w", err)
    }

	// Validate URL
	if !strings.HasPrefix(webFetchInput.URL, "http://") && !strings.HasPrefix(webFetchInput.URL, "https://") {
		return "", fmt.Errorf("URL must start with http:// or https://")
	}

    // Create HTTP client with timeout
    client := &http.Client{
        Timeout: 30 * time.Second,
    }

    // Create request
    req, err := http.NewRequest("GET", webFetchInput.URL, nil)
    if err != nil {
        return "", fmt.Errorf("failed to create request: %w", err)
    }

	// Add standard headers
	req.Header.Set("User-Agent", "Mozilla/5.0 WebFetch Tool")
	req.Header.Set("Accept", "text/*, application/json, application/xml, application/xhtml+xml")

    // Make request
    resp, err := client.Do(req)
    if err != nil {
        return "", fmt.Errorf("HTTP GET error: %w", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode < 200 || resp.StatusCode >= 300 {
        return "", fmt.Errorf("HTTP request failed with status code: %d", resp.StatusCode)
    }


    // Read the response body
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return "", fmt.Errorf("error reading response body: %w", err)
    }

    // Generate filename
    baseFilename, err := generateFilename(webFetchInput.URL)
    if err != nil {
        return "", fmt.Errorf("failed to generate filename: %v", err)
    }
    extension, allowed := isAllowedContentType(resp.Header.Get("Content-Type"))
    if !allowed {
        return "", fmt.Errorf("unsupported content type: %s", resp.Header.Get("Content-Type"))
    }
    filename := baseFilename + extension

    // Create cache directory
    cacheDir := ".agent/cache/webfetch"
    if err := os.MkdirAll(cacheDir, 0755); err != nil {
        return "", fmt.Errorf("failed to create cache directory: %w", err)
    }

    // Create cache file path
    cacheFilePath := filepath.Join(cacheDir, filename)
    file, err := os.Create(cacheFilePath)
    if err != nil {
        return "", fmt.Errorf("failed to create cache file: %w", err)
    }
    defer file.Close()

    // Write response to cache file
    _, err = file.Write(body)
    if err != nil {
        os.Remove(cacheFilePath) // Clean up on error
        return "", fmt.Errorf("failed to write content: %w", err)
    }

    // Construct result string, including status code
    result := fmt.Sprintf("{\"path\": \"%s\", \"statusCode\": %d, \"contentType\": \"%s\"}", cacheFilePath, resp.StatusCode, resp.Header.Get("Content-Type"))

    return result, nil
}

// generateTodoID creates a simple unique ID for todo items
func generateTodoID() string {
	sessionMutex.Lock()
	defer sessionMutex.Unlock()
	
	id := "task-" + strconv.Itoa(sessionIDCounter)
	sessionIDCounter++
	return id
}

// getCurrentSessionID returns the current session identifier
func getCurrentSessionID() string {
	if currentSessionID == "" {
		return "default" // Fallback for backward compatibility
	}
	return currentSessionID
}

// validateTodoStatus checks if status is valid
func validateTodoStatus(status string) bool {
	validStatuses := map[string]bool{
		"pending":     true,
		"in_progress": true,
		"completed":   true,
		"cancelled":   true,
	}
	return validStatuses[status]
}

// validateTodoPriority checks if priority is valid
func validateTodoPriority(priority string) bool {
	validPriorities := map[string]bool{
		"high":   true,
		"medium": true,
		"low":    true,
	}
	return validPriorities[priority]
}

// TodoWrite manages the todo list for the current session
func TodoWrite(input json.RawMessage) (string, error) {
	todoWriteInput := TodoWriteInput{}
	err := json.Unmarshal(input, &todoWriteInput)
	if err != nil {
		return "", fmt.Errorf("invalid input: %v", err)
	}

	// Parse the JSON string containing the todos array
	var todos []TodoItem
	if todoWriteInput.TodosJSON != "" {
		err = json.Unmarshal([]byte(todoWriteInput.TodosJSON), &todos)
		if err != nil {
			return "", fmt.Errorf("invalid todos JSON: %v", err)
		}
	}

	sessionID := getCurrentSessionID()
	
	// Validate and process todos
	var processedTodos []TodoItem
	inProgressCount := 0
	
	for i, todo := range todos {
		// Generate ID if not provided
		if todo.ID == "" {
			todo.ID = generateTodoID()
		}
		
		// Validate status
		if !validateTodoStatus(todo.Status) {
			return "", fmt.Errorf("invalid status '%s' for todo %d. Must be one of: pending, in_progress, completed, cancelled", todo.Status, i+1)
		}
		
		// Validate priority
		if !validateTodoPriority(todo.Priority) {
			return "", fmt.Errorf("invalid priority '%s' for todo %d. Must be one of: high, medium, low", todo.Priority, i+1)
		}
		
		// Count in_progress todos
		if todo.Status == "in_progress" {
			inProgressCount++
		}
		
		// Validate content is not empty
		if strings.TrimSpace(todo.Content) == "" {
			return "", fmt.Errorf("todo content cannot be empty for todo %d", i+1)
		}
		
		processedTodos = append(processedTodos, todo)
	}
	
	// Enforce only one in_progress task rule
	if inProgressCount > 1 {
		return "", fmt.Errorf("only one task can be 'in_progress' at a time, found %d", inProgressCount)
	}
	
	// Update session state
	sessionMutex.Lock()
	sessionTodos[sessionID] = processedTodos
	sessionMutex.Unlock()
	
	// Persist to file if we have a session directory
	if currentSessionDir != "" {
		todosPath := filepath.Join(currentSessionDir, "todos.json")
		sessionTodosData := SessionTodos{Todos: processedTodos}
		if err := saveTodosToFile(todosPath, sessionTodosData); err != nil {
			fmt.Printf("Warning: Failed to persist todos to file: %v\n", err)
		}
	}
	
	// Count non-completed todos for title
	nonCompletedCount := 0
	for _, todo := range processedTodos {
		if todo.Status != "completed" && todo.Status != "cancelled" {
			nonCompletedCount++
		}
	}
	
	// Return result
	result := map[string]interface{}{
		"title":  fmt.Sprintf("Updated todo list with %d active todos", nonCompletedCount),
		"output": fmt.Sprintf("Successfully updated %d todos", len(processedTodos)),
	}
	
	jsonResult, err := json.Marshal(result)
	if err != nil {
		return "", fmt.Errorf("failed to marshal result: %v", err)
	}
	
	return string(jsonResult), nil
}

// TodoRead retrieves the current todo list from session state
func TodoRead(input json.RawMessage) (string, error) {
	sessionID := getCurrentSessionID()
	
	sessionMutex.RLock()
	todos, exists := sessionTodos[sessionID]
	sessionMutex.RUnlock()
	
	if !exists || len(todos) == 0 {
		result := map[string]interface{}{
			"title":  "0 todos",
			"output": "[]",
		}
		
		jsonResult, err := json.Marshal(result)
		if err != nil {
			return "", fmt.Errorf("failed to marshal result: %v", err)
		}
		
		return string(jsonResult), nil
	}
	
	// Count non-completed todos for title
	nonCompletedCount := 0
	for _, todo := range todos {
		if todo.Status != "completed" && todo.Status != "cancelled" {
			nonCompletedCount++
		}
	}
	
	// Convert todos to JSON
	todosJSON, err := json.Marshal(todos)
	if err != nil {
		return "", fmt.Errorf("failed to marshal todos: %v", err)
	}
	
	result := map[string]interface{}{
		"title":  fmt.Sprintf("%d todos", nonCompletedCount),
		"output": string(todosJSON),
	}
	
	jsonResult, err := json.Marshal(result)
	if err != nil {
		return "", fmt.Errorf("failed to marshal result: %v", err)
	}
	
	return string(jsonResult), nil
}

func Head(input json.RawMessage) (string, error) {
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

func Tail(input json.RawMessage) (string, error) {
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

func Cloc(input json.RawMessage) (string, error) {
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

func Glob(input json.RawMessage) (string, error) {
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

// Session management functions
func generateSessionID() string {
	return time.Now().Format("2006-01-02-15-04-05")
}

func initializeSession(resumeSessionID string) (*SessionManager, error) {
	var sessionID string
	var isResume bool
	
	if resumeSessionID != "" {
		if resumeSessionID == "list" {
			// Interactive session selection
			selectedSession, err := selectSessionInteractively()
			if err != nil {
				return nil, err
			}
			if selectedSession == "" {
				// User cancelled or no sessions available, create new session
				sessionID = generateSessionID()
				isResume = false
			} else {
				sessionID = selectedSession
				isResume = true
			}
		} else {
			// Specific session ID provided - validate format first
			if !isValidSessionIDFormat(resumeSessionID) {
				// Invalid format, offer interactive selection
				fmt.Printf("Invalid session ID format: %s\n", resumeSessionID)
				fmt.Println("Session ID should be in format YYYY-MM-DD-HH-MM-SS")
				fmt.Println("Would you like to select from available sessions? (y/n): ")
				
				var response string
				fmt.Scanln(&response)
				response = strings.ToLower(strings.TrimSpace(response))
				
				if response == "y" || response == "yes" {
					selectedSession, err := selectSessionInteractively()
					if err != nil {
						return nil, err
					}
					if selectedSession == "" {
						sessionID = generateSessionID()
						isResume = false
					} else {
						sessionID = selectedSession
						isResume = true
					}
				} else {
					return nil, fmt.Errorf("invalid session ID format: %s", resumeSessionID)
				}
			} else {
				sessionID = resumeSessionID
				isResume = true
			}
		}
	} else {
		sessionID = generateSessionID()
		isResume = false
	}
	
	// Set global variables
	currentSessionID = sessionID
	currentSessionDir = filepath.Join(".agent", "sessions", sessionID)
	
	sessionManager := &SessionManager{
		SessionID:  sessionID,
		SessionDir: currentSessionDir,
		TodosPath:  filepath.Join(currentSessionDir, "todos.json"),
	}
	
	if isResume {
		// Load existing session
		if err := loadSession(sessionManager); err != nil {
			return nil, fmt.Errorf("failed to load session %s: %w", sessionID, err)
		}
	} else {
		// Create new session
		if err := createNewSession(sessionManager); err != nil {
			return nil, fmt.Errorf("failed to create new session %s: %w", sessionID, err)
		}
	}
	
	// Set global reference
	currentSessionManager = sessionManager
	
	return sessionManager, nil
}

func createNewSession(sm *SessionManager) error {
	// Create session directory
	if err := os.MkdirAll(sm.SessionDir, 0755); err != nil {
		return fmt.Errorf("failed to create session directory: %w", err)
	}
	
	// Create empty agent.log file
	logPath := filepath.Join(sm.SessionDir, "agent.log")
	logFile, err := os.Create(logPath)
	if err != nil {
		return fmt.Errorf("failed to create agent.log: %w", err)
	}
	logFile.Close() // We'll reopen when needed
	
	// Initialize empty todos.json
	emptyTodos := SessionTodos{Todos: []TodoItem{}}
	if err := saveTodosToFile(sm.TodosPath, emptyTodos); err != nil {
		return fmt.Errorf("failed to create todos.json: %w", err)
	}
	
	// Initialize empty conversation
	sm.Conversation = []openai.ChatCompletionMessageParamUnion{}
	
	return nil
}

func loadSession(sm *SessionManager) error {
	// Check if session directory exists
	if _, err := os.Stat(sm.SessionDir); os.IsNotExist(err) {
		return fmt.Errorf("session directory does not exist")
	}
	
	// Load conversation from agent.log
	logPath := filepath.Join(sm.SessionDir, "agent.log")
	conversation, err := loadConversationFromFile(logPath)
	if err != nil {
		return fmt.Errorf("failed to load conversation: %w", err)
	}
	sm.Conversation = conversation
	
	// Load todos from todos.json
	todos, err := loadTodosFromFile(sm.TodosPath)
	if err != nil {
		return fmt.Errorf("failed to load todos: %w", err)
	}
	
	// Update session state with loaded todos
	sessionMutex.Lock()
	sessionTodos[sm.SessionID] = todos.Todos
	sessionMutex.Unlock()
	
	return nil
}

func loadConversationFromFile(logPath string) ([]openai.ChatCompletionMessageParamUnion, error) {
	// Check if file exists
	if _, err := os.Stat(logPath); os.IsNotExist(err) {
		// Return empty conversation if file doesn't exist (new session)
		return []openai.ChatCompletionMessageParamUnion{}, nil
	}
	
	data, err := os.ReadFile(logPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read log file: %w", err)
	}
	
	// Handle empty file
	if len(data) == 0 {
		return []openai.ChatCompletionMessageParamUnion{}, nil
	}
	
	var conversation []openai.ChatCompletionMessageParamUnion
	if err := json.Unmarshal(data, &conversation); err != nil {
		return nil, fmt.Errorf("failed to parse conversation JSON: %w", err)
	}
	
	return conversation, nil
}

func loadTodosFromFile(todosPath string) (SessionTodos, error) {
	// Check if file exists
	if _, err := os.Stat(todosPath); os.IsNotExist(err) {
		// Return empty todos if file doesn't exist
		return SessionTodos{Todos: []TodoItem{}}, nil
	}
	
	data, err := os.ReadFile(todosPath)
	if err != nil {
		return SessionTodos{}, fmt.Errorf("failed to read todos file: %w", err)
	}
	
	// Handle empty file
	if len(data) == 0 {
		return SessionTodos{Todos: []TodoItem{}}, nil
	}
	
	var todos SessionTodos
	if err := json.Unmarshal(data, &todos); err != nil {
		return SessionTodos{}, fmt.Errorf("failed to parse todos JSON: %w", err)
	}
	
	return todos, nil
}

func saveTodosToFile(todosPath string, todos SessionTodos) error {
	data, err := json.MarshalIndent(todos, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal todos: %w", err)
	}
	
	// Ensure directory exists
	if err := os.MkdirAll(filepath.Dir(todosPath), 0755); err != nil {
		return fmt.Errorf("failed to create todos directory: %w", err)
	}
	
	return os.WriteFile(todosPath, data, 0644)
}

func listAvailableSessions() ([]string, error) {
	sessionsDir := ".agent/sessions"
	
	// Check if sessions directory exists
	if _, err := os.Stat(sessionsDir); os.IsNotExist(err) {
		return []string{}, nil
	}
	
	entries, err := os.ReadDir(sessionsDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read sessions directory: %w", err)
	}
	
	var sessions []string
	for _, entry := range entries {
		if entry.IsDir() {
			// Validate session ID format (YYYY-MM-DD-HH-MM-SS)
			sessionID := entry.Name()
			if isValidSessionIDFormat(sessionID) {
				sessions = append(sessions, sessionID)
			}
		}
	}
	
	return sessions, nil
}

func isValidSessionIDFormat(sessionID string) bool {
	// Check if format matches YYYY-MM-DD-HH-MM-SS
	_, err := time.Parse("2006-01-02-15-04-05", sessionID)
	return err == nil
}

func selectSessionInteractively() (string, error) {
	sessions, err := listAvailableSessions()
	if err != nil {
		return "", fmt.Errorf("failed to list available sessions: %w", err)
	}
	
	if len(sessions) == 0 {
		fmt.Println("No previous sessions found. Creating new session...")
		return "", nil
	}
	
	// Sort sessions in descending order (newest first)
	sort.Slice(sessions, func(i, j int) bool {
		return sessions[i] > sessions[j]
	})
	
	fmt.Println("Available sessions:")
	for i, session := range sessions {
		fmt.Printf("  %d) %s\n", i+1, session)
	}
	fmt.Printf("  %d) Create new session\n", len(sessions)+1)
	fmt.Printf("\nSelect a session (1-%d): ", len(sessions)+1)
	
	var choice int
	_, err = fmt.Scanf("%d", &choice)
	if err != nil {
		return "", fmt.Errorf("invalid input: please enter a number")
	}
	
	if choice < 1 || choice > len(sessions)+1 {
		return "", fmt.Errorf("invalid choice: please select a number between 1 and %d", len(sessions)+1)
	}
	
	if choice == len(sessions)+1 {
		// User chose to create new session
		return "", nil
	}
	
	// Return selected session (convert from 1-based to 0-based index)
	return sessions[choice-1], nil
}

func HtmlToMarkdown(input json.RawMessage) (string, error) {
	htmlToMarkdownInput := HtmlToMarkdownInput{}
	if err := json.Unmarshal(input, &htmlToMarkdownInput); err != nil {
		return "", fmt.Errorf("invalid input: %v", err)
	}

	if htmlToMarkdownInput.Path == "" {
		return "", fmt.Errorf("path cannot be empty")
	}

	// Check if input file exists
	if _, err := os.Stat(htmlToMarkdownInput.Path); os.IsNotExist(err) {
		return "", fmt.Errorf("input file does not exist: %s", htmlToMarkdownInput.Path)
	}

	// Read HTML file
	htmlContent, err := os.ReadFile(htmlToMarkdownInput.Path)
	if err != nil {
		return "", fmt.Errorf("failed to read HTML file: %v", err)
	}

	// Convert HTML to Markdown
	markdown, err := htmltomarkdown.ConvertString(string(htmlContent))
	if err != nil {
		return "", fmt.Errorf("failed to convert HTML to markdown: %v", err)
	}

	// Generate output filename: replace .html with .md, or append .md if no extension
	outputPath := htmlToMarkdownInput.Path
	ext := filepath.Ext(outputPath)
	if ext != "" {
		outputPath = strings.TrimSuffix(outputPath, ext) + ".md"
	} else {
		outputPath = outputPath + ".md"
	}

	// Write markdown to output file
	err = os.WriteFile(outputPath, []byte(markdown), 0644)
	if err != nil {
		return "", fmt.Errorf("failed to write markdown file: %v", err)
	}

	// Return result as JSON
	result := map[string]interface{}{
		"inputPath":  htmlToMarkdownInput.Path,
		"outputPath": outputPath,
		"success":    true,
	}

	jsonResult, err := json.Marshal(result)
	if err != nil {
		return "", fmt.Errorf("failed to marshal result: %v", err)
	}

	return string(jsonResult), nil
}


