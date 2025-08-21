package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"
	"github.com/chzyer/readline"
	"github.com/openai/openai-go/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"agent/internal/config"
	"agent/internal/session"
	"agent/test/helpers"
)// MockOpenAIClient is a mock implementation of OpenAI client for testing
type MockOpenAIClient struct {
	mock.Mock

}

// MockChatCompletionsService mocks the Chat.Completions service
type MockChatCompletionsService struct {
	mock.Mock

}

// MockChatService mocks the Chat service
type MockChatService struct {
	Completions *MockChatCompletionsService

}

// New mocks the completions creation
func (m *MockChatCompletionsService) New(ctx context.Context, params openai.ChatCompletionNewParams) (*openai.ChatCompletion, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	
}
	return args.Get(0).(*openai.ChatCompletion), args.Error(1)

}

// MockSessionManager is a mock implementation that embeds session.Manager interface
type MockSessionManager struct {
	mock.Mock

}

func (m *MockSessionManager) CreateSession() (*session.Session, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	
}
	return args.Get(0).(*session.Session), args.Error(1)

}

func (m *MockSessionManager) LoadSession(id string) (*session.Session, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	
}
	return args.Get(0).(*session.Session), args.Error(1)

}

func (m *MockSessionManager) SaveSession(session *session.Session) error {
	args := m.Called(session)
	return args.Error(0)

}

func (m *MockSessionManager) SelectSessionInteractively() (string, error) {
	args := m.Called()
	return args.String(0), args.Error(1)

}

func (m *MockSessionManager) GetTodos(sessionID string) ([]session.TodoItem, error) {
	args := m.Called(sessionID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	
}
	return args.Get(0).([]session.TodoItem), args.Error(1)

}

func (m *MockSessionManager) SaveTodos(sessionID string, todos []session.TodoItem) error {
	args := m.Called(sessionID, todos)
	return args.Error(0)

}

func (m *MockSessionManager) GenerateTodoID() string {
	args := m.Called()
	return args.String(0)

}

// Mock getUserMessage function
func mockGetUserMessage(input string, shouldContinue bool) func() (string, bool) {
	callCount := 0
	return func() (string, bool) {
		callCount++
		if callCount == 1 {
			return input, true
		
}
		return "", shouldContinue
	
}

}

// Helper to create a test agent with mocked dependencies
func createTestAgent(t *testing.T, cfg *config.Config, client *openai.Client) (*Agent, *MockSessionManager) {
	t.Helper()	// Create mock session manager
	mockSessionManager := &MockSessionManager{
}
	// Create readline instance (we don't mock this as it's not used in core logic)
	rl, err := readline.New("")
	if err != nil {
		t.Fatal(err)
	
}
	defer rl.Close()
	deps := Dependencies{
		Client:           client,
		GetUserMessage:   mockGetUserMessage("test input", false),
		ReadlineInstance: rl,
		PrePrompts:       []string{
},
		RequestDelay:     0,
		SingleShot:       false,
	
}
	agent, err := New(cfg, deps)
	require.NoError(t, err)

	// Replace the session manager with our mock
	agent.sessionManager = mockSessionManager

	return agent, mockSessionManager

}

func TestNew(t *testing.T) {
	tests := []struct {
		name          string
		config        *config.Config
		deps          Dependencies
		expectedError string
	}{
		{
			name:   "valid configuration and dependencies",
			config: helpers.TestConfig(t),
			deps: Dependencies{
				Client:           &openai.Client{},
				GetUserMessage:   mockGetUserMessage("test", false),
				ReadlineInstance: nil, // Can be nil for testing
				PrePrompts:       []string{"prompt1", "prompt2"},
				RequestDelay:     time.Second,
				SingleShot:       true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			agent, err := New(tt.config, tt.deps)
			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				return
			}

			require.NoError(t, err)
			assert.NotNil(t, agent)
			assert.Equal(t, tt.config, agent.config)
			assert.Equal(t, tt.deps.Client, agent.client)
			assert.Equal(t, tt.deps.PrePrompts, agent.prePrompts)
			assert.Equal(t, tt.deps.RequestDelay, agent.requestDelay)
			assert.Equal(t, tt.deps.SingleShot, agent.singleShot)
			assert.NotNil(t, agent.sessionManager)
		})
	}
}

func TestSetFlags(t *testing.T) {
	cfg := helpers.TestConfig(t)
	agent, _ := createTestAgent(t, cfg, &openai.Client{})
	flags := Flags{
		PromptFile:    "/test/prompt.txt",
		ContinueChat:  true,
		Timeout:       30 * time.Second,
		InitFlag:      true,
		PrePrompts:    "/test/pre",
		ResumeSession: "session-123",
	}
	agent.SetFlags(flags)
	assert.Equal(t, flags, agent.flags)
}

func TestSetTransitionToInteractive(t *testing.T) {
	cfg := helpers.TestConfig(t)
	agent, _ := createTestAgent(t, cfg, &openai.Client{})
	agent.SetTransitionToInteractive(true)
	assert.True(t, agent.transitionToInteractive)
	agent.SetTransitionToInteractive(false)
	assert.False(t, agent.transitionToInteractive)
}

func TestInitializeSession(t *testing.T) {
	tests := []struct {
		name          string
		resumeSession string
		mockSetup     func(*MockSessionManager)
		expectError   string
	}{
		{
			name:          "create new session",
			resumeSession: "",
			mockSetup: func(m *MockSessionManager) {
				testSession := &session.Session{
					ID:  "new-session-12p3",
					Dir: "/tmp/sessions/new-session-123",
				}
				m.On("CreateSession").Return(testSession, nil)
			},
		},
		{
			name:          "load existing session",
			resumeSession: "existing-session-456",
			mockSetup: func(m *MockSessionManager) {
				testSession := &session.Session{
					ID:  "existing-session-456",
					Dir: "/tmp/sessions/existing-session-456",
				}
				m.On("LoadSession", "existing-session-456").Return(testSession, nil)
			},
		},
		{
			name:          "interactive session selection - create new",
			resumeSession: "list",
			mockSetup: func(m *MockSessionManager) {
				// Return empty string to create new session
				m.On("SelectSessionInteractively").Return("", nil)
				testSession := &session.Session{
					ID:  "new-session-789",
					Dir: "/tmp/sessions/new-session-789",
				}
				m.On("CreateSession").Return(testSession, nil)
			},
		},
		{
			name:          "interactive session selection - load existing",
			resumeSession: "list",
			mockSetup: func(m *MockSessionManager) {
				// Return session ID to load
				m.On("SelectSessionInteractively").Return("selected-session-abc", nil)
				testSession := &session.Session{
					ID:  "selected-session-abc",
					Dir: "/tmp/sessions/selected-session-abc",
				}
				m.On("LoadSession", "selected-session-abc").Return(testSession, nil)
			},
		},
		{
			name:          "session creation error",
			resumeSession: "",
			mockSetup: func(m *MockSessionManager) {
				m.On("CreateSession").Return(nil, fmt.Errorf("failed to create session"))
			},
			expectError: "initializing session",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := helpers.TestConfig(t)
			agent, mockSession := createTestAgent(t, cfg, &openai.Client{})
			agent.flags.ResumeSession = tt.resumeSession
			tt.mockSetup(mockSession)
			err := agent.InitializeSession()
			if tt.expectError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectError)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, agent.currentSession)
				assert.NotNil(t, agent.toolRegistry)
			}

			mockSession.AssertExpectations(t)
		})
	}
}

func TestGetCurrentSessionID(t *testing.T) {
	cfg := helpers.TestConfig(t)
	agent, _ := createTestAgent(t, cfg, &openai.Client{})
	// Test with no current session
	assert.Equal(t, "", agent.GetCurrentSessionID())

	// Test with current session
	agent.currentSession = &session.Session{ID: "test-session-123"}
	assert.Equal(t, "test-session-123", agent.GetCurrentSessionID())
}

func TestGetCurrentSessionDir(t *testing.T) {
	cfg := helpers.TestConfig(t)
	agent, _ := createTestAgent(t, cfg, &openai.Client{})
	// Test with no current session
	assert.Equal(t, "", agent.GetCurrentSessionDir())

	// Test with current session
	agent.currentSession = &session.Session{Dir: "/tmp/test-session-dir"}
	assert.Equal(t, "/tmp/test-session-dir", agent.GetCurrentSessionDir())
}

func TestLogConversation(t *testing.T) {
	cfg := helpers.TestConfig(t)
	agent, mockSession := createTestAgent(t, cfg, &openai.Client{})
	conversation := []openai.ChatCompletionMessageParamUnion{
		openai.UserMessage("Hello"),
		openai.AssistantMessage("Hi there"),
	}
	// Test with no current session
	err := agent.LogConversation(conversation)
	assert.NoError(t, err) // Should not error, just not save
	// Test with current session
	testSession := &session.Session{ID: "test-session"}
	agent.currentSession = testSession
	mockSession.On("SaveSession", testSession).Return(nil)
	err = agent.LogConversation(conversation)
	assert.NoError(t, err)
	assert.Equal(t, conversation, agent.currentSession.Conversation)
	mockSession.AssertExpectations(t)
}func TestExecuteTool(t *testing.T) {
	cfg := helpers.TestConfig(t)
	agent, mockSession := createTestAgent(t, cfg, &openai.Client{
})
	// Initialize session so tool registry is set up
	testSession := &session.Session{ID: "test-session", Dir: "/tmp/test"
}
	mockSession.On("CreateSession").Return(testSession, nil)
	err := agent.InitializeSession()
	require.NoError(t, err)	// Test successful tool execution (using a real tool from registry)
	input := json.RawMessage(`{"path": "/nonexistent/file.txt"
}`)
	result := agent.ExecuteTool("call-123", "read_file", input)	// The result is a union type, we can test that it's not nil
	assert.NotNil(t, result)	// Test that the result is a proper tool message by checking it implements the interface
	// We just verify it doesn't panic and produces some result
	assert.IsType(t, openai.ChatCompletionMessageParamUnion{
}, result)

}func TestRunInference(t *testing.T) {
	cfg := helpers.TestConfig(t)

	// Create actual OpenAI client for structure, but we can't test actual API calls
	// This test focuses on the method structure and error handling
	tests := []struct {
		name         string
		conversation []openai.ChatCompletionMessageParamUnion
		expectError  bool
	
}{
		{
			name: "valid conversation",
			conversation: []openai.ChatCompletionMessageParamUnion{
				openai.UserMessage("Hello"),
			
},
			expectError: true, // Will error due to invalid API key in test
		
},
	
}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			agent, mockSession := createTestAgent(t, cfg, &openai.Client{
})			// Initialize session
			testSession := &session.Session{ID: "test-session", Dir: "/tmp/test"
}
			mockSession.On("CreateSession").Return(testSession, nil)
			err := agent.InitializeSession()
			require.NoError(t, err)			ctx := context.Background()
			completion, err := agent.RunInference(ctx, tt.conversation)			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, completion)
			
} else {
				assert.NoError(t, err)
				assert.NotNil(t, completion)
			
}
		
})
	
}

}
func TestRun_InitFlag(t *testing.T) {
	cfg := helpers.TestConfig(t)
	agent, _ := createTestAgent(t, cfg, &openai.Client{
})
	agent.flags.InitFlag = true

	ctx := context.Background()
	err := agent.Run(ctx)	assert.Error(t, err)
	assert.Contains(t, err.Error(), "init functionality not yet implemented")

}func TestRun_SessionInitializationError(t *testing.T) {
	cfg := helpers.TestConfig(t)
	agent, mockSession := createTestAgent(t, cfg, &openai.Client{
})

	// Mock session creation to fail
	mockSession.On("CreateSession").Return(nil, fmt.Errorf("session creation failed"))

	ctx := context.Background()
	err := agent.Run(ctx)	assert.Error(t, err)
	assert.Contains(t, err.Error(), "initializing session")	mockSession.AssertExpectations(t)

}func TestRun_SingleShotMode(t *testing.T) {
	cfg := helpers.TestConfig(t)

	// Create agent with single shot mode
	agent, mockSession := createTestAgent(t, cfg, &openai.Client{
})
	agent.singleShot = true
	agent.getUserMessage = mockGetUserMessage("test message", false)

	// Mock session creation
	testSession := &session.Session{
		ID:           "test-session",
		Dir:          "/tmp/test",
		Conversation: []openai.ChatCompletionMessageParamUnion{
},
	
}
	mockSession.On("CreateSession").Return(testSession, nil)
	// Don't expect SaveSession - agent will error on API call before reaching save	ctx := context.Background()
	err := agent.Run(ctx)	// Will error due to OpenAI API call with test credentials, but that's expected
	assert.Error(t, err)	mockSession.AssertExpectations(t)

}func TestRun_WithPrePrompts(t *testing.T) {
	cfg := helpers.TestConfig(t)
	agent, mockSession := createTestAgent(t, cfg, &openai.Client{
})
	agent.prePrompts = []string{"System prompt", "", "Another prompt"
}
	agent.singleShot = true
	agent.getUserMessage = mockGetUserMessage("user input", false)

	// Mock session creation
	testSession := &session.Session{
		ID:           "test-session",
		Dir:          "/tmp/test",
		Conversation: []openai.ChatCompletionMessageParamUnion{
},
	
}
	mockSession.On("CreateSession").Return(testSession, nil)
	// Don't expect SaveSession - agent will error on API call before reaching save	ctx := context.Background()
	err := agent.Run(ctx)	// Will error due to OpenAI API call, but should have added pre-prompts to conversation
	assert.Error(t, err)	mockSession.AssertExpectations(t)

}func TestRun_ResumeWithExistingConversation(t *testing.T) {
	cfg := helpers.TestConfig(t)
	agent, mockSession := createTestAgent(t, cfg, &openai.Client{
})
	agent.singleShot = true
	agent.getUserMessage = mockGetUserMessage("new input", false)

	// Mock session with existing conversation
	existingConversation := []openai.ChatCompletionMessageParamUnion{
		openai.UserMessage("Previous message"),
		openai.AssistantMessage("Previous response"),
	
}
	testSession := &session.Session{
		ID:           "existing-session",
		Dir:          "/tmp/test",
		Conversation: existingConversation,
	
}
	mockSession.On("CreateSession").Return(testSession, nil)
	// Don't expect SaveSession - agent will error on API call before reaching save	ctx := context.Background()
	err := agent.Run(ctx)	// Will error due to OpenAI API call
	assert.Error(t, err)	mockSession.AssertExpectations(t)

}
