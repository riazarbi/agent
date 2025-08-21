package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/chzyer/readline"
	"github.com/openai/openai-go/v2"

	"agent/internal/config"
	"agent/internal/session"
	"agent/internal/tools"
)

// Flags represents command-line flags passed to the agent
type Flags struct {
	PromptFile    string
	ContinueChat  bool
	Timeout       time.Duration
	InitFlag      bool
	PrePrompts    string
	ResumeSession string
}

// Agent represents the main application agent
type Agent struct {
	config                  *config.Config
	flags                   Flags
	client                  *openai.Client
	sessionManager          session.SessionManager
	currentSession          *session.Session
	toolRegistry            *tools.Registry
	getUserMessage          func() (string, bool)
	rl                      *readline.Instance
	singleShot              bool
	transitionToInteractive bool
	prePrompts              []string
	requestDelay            time.Duration
}

// Dependencies holds the dependencies needed to create an Agent
type Dependencies struct {
	Client           *openai.Client
	GetUserMessage   func() (string, bool)
	ReadlineInstance *readline.Instance
	PrePrompts       []string
	RequestDelay     time.Duration
	SingleShot       bool
}

// New creates a new Agent instance with the provided configuration and dependencies
func New(cfg *config.Config, deps Dependencies) (*Agent, error) {
	// Initialize session manager
	sessionConfig := session.Config{
		SessionsDir: cfg.Session.Dir,
	}
	sessionMgr := session.NewManager(sessionConfig)

	agent := &Agent{
		config:         cfg,
		client:         deps.Client,
		sessionManager: sessionMgr,
		getUserMessage: deps.GetUserMessage,
		rl:             deps.ReadlineInstance,
		prePrompts:     deps.PrePrompts,
		requestDelay:   deps.RequestDelay,
		singleShot:     deps.SingleShot,
	}

	return agent, nil
}

// SetFlags sets the command-line flags for the agent
func (a *Agent) SetFlags(flags Flags) {
	a.flags = flags
}

// SetTransitionToInteractive sets whether the agent should transition to interactive mode
func (a *Agent) SetTransitionToInteractive(transition bool) {
	a.transitionToInteractive = transition
}

// InitializeSession initializes or resumes a session based on flags
func (a *Agent) InitializeSession() error {
	var currentSession *session.Session
	var err error

	if a.flags.ResumeSession != "" {
		if a.flags.ResumeSession == "list" {
			// Interactive session selection
			selectedSessionID, err := a.sessionManager.SelectSessionInteractively()
			if err != nil {
				return fmt.Errorf("selecting session: %w", err)
			}

			if selectedSessionID == "" {
				// Create new session
				currentSession, err = a.sessionManager.CreateSession()
			} else {
				// Load selected session
				currentSession, err = a.sessionManager.LoadSession(selectedSessionID)
			}
		} else {
			// Load specific session
			currentSession, err = a.sessionManager.LoadSession(a.flags.ResumeSession)
		}
	} else {
		// Create new session
		currentSession, err = a.sessionManager.CreateSession()
	}

	if err != nil {
		return fmt.Errorf("initializing session: %w", err)
	}

	a.currentSession = currentSession

	// Initialize tool registry with session dependencies
	registryConfig := &tools.RegistryConfig{
		SessionManager:   a.sessionManager,
		CurrentSessionID: a.currentSession.ID,
	}
	a.toolRegistry = tools.NewRegistry(registryConfig)

	return nil
}

// GetCurrentSessionID returns the current session ID
func (a *Agent) GetCurrentSessionID() string {
	if a.currentSession != nil {
		return a.currentSession.ID
	}
	return ""
}

// GetCurrentSessionDir returns the current session directory
func (a *Agent) GetCurrentSessionDir() string {
	if a.currentSession != nil {
		return a.currentSession.Dir
	}
	return ""
}

// LogConversation saves the conversation to the current session
func (a *Agent) LogConversation(conversation []openai.ChatCompletionMessageParamUnion) error {
	if a.currentSession != nil {
		a.currentSession.Conversation = conversation
		return a.sessionManager.SaveSession(a.currentSession)
	}
	return nil
}

// ExecuteTool executes a tool by name with the given input
func (a *Agent) ExecuteTool(id, name string, input json.RawMessage) openai.ChatCompletionMessageParamUnion {
	fmt.Printf("\u001b[92mTool\u001b[0m: %s(%s)\n", name, input)
	response, err := a.toolRegistry.Execute(name, input)
	if err != nil {
		return openai.ToolMessage(err.Error(), id)
	}
	return openai.ToolMessage(response, id)
}

// RunInference runs model inference with the current conversation
func (a *Agent) RunInference(ctx context.Context, conversation []openai.ChatCompletionMessageParamUnion) (*openai.ChatCompletion, error) {
	openaiTools := []openai.ChatCompletionToolUnionParam{}
	for _, tool := range a.toolRegistry.List() {
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
		Model:     "gemini-2.5-pro",
		MaxTokens: openai.Int(8192),
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

// Run starts the agent with the given context
func (a *Agent) Run(ctx context.Context) error {
	// Handle init command first
	if a.flags.InitFlag {
		return fmt.Errorf("init functionality not yet implemented in refactored agent")
	}

	// Initialize session first
	if err := a.InitializeSession(); err != nil {
		return fmt.Errorf("initializing session: %w", err)
	}

	// Initialize conversation from session if available, otherwise empty
	var conversation []openai.ChatCompletionMessageParamUnion
	if len(a.currentSession.Conversation) > 0 {
		// Resume with existing conversation
		conversation = a.currentSession.Conversation
		fmt.Printf("Resuming session %s with %d previous messages\n", a.currentSession.ID, len(conversation))
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
		fmt.Printf("Chat with Agent at %s\n", a.config.API.BaseURL)
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

		completion, err := a.RunInference(ctx, conversation)
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
			result := a.ExecuteTool(toolCall.ID, toolCall.Function.Name, json.RawMessage(toolCall.Function.Arguments))
			toolResults = append(toolResults, result)
		}

		if len(toolResults) == 0 {
			readUserInput = true
			// Log conversation state after text-only response
			if err := a.LogConversation(conversation); err != nil {
				fmt.Printf("Warning: Failed to log conversation: %v\n", err)
			}
			if a.singleShot && !a.transitionToInteractive {
				break
			}
			continue
		}

		readUserInput = false
		conversation = append(conversation, toolResults...)

		// Log conversation state after each cycle
		if err := a.LogConversation(conversation); err != nil {
			fmt.Printf("Warning: Failed to log conversation: %v\n", err)
		}
	}

	return nil
}
