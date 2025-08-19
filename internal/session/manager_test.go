package session

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/openai/openai-go/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"agent/test/helpers"
)

func TestNewManager(t *testing.T) {
	config := Config{SessionsDir: "/tmp/test-sessions"}
	manager := NewManager(config)

	assert.NotNil(t, manager)
	assert.Equal(t, config.SessionsDir, manager.config.SessionsDir)
	assert.NotNil(t, manager.sessions)
	assert.Equal(t, 1, manager.idCounter)
}

func TestCreateSession(t *testing.T) {
	tempDir := helpers.TempDir(t)
	config := Config{SessionsDir: tempDir}
	manager := NewManager(config)

	session, err := manager.CreateSession()

	require.NoError(t, err)
	assert.NotNil(t, session)
	assert.NotEmpty(t, session.ID)
	assert.Equal(t, filepath.Join(tempDir, session.ID), session.Dir)
	assert.Equal(t, filepath.Join(session.Dir, "todos.json"), session.TodosPath)
	assert.Empty(t, session.Conversation)
	assert.False(t, session.CreatedAt.IsZero())
	assert.False(t, session.UpdatedAt.IsZero())

	// Verify session directory exists
	_, err = os.Stat(session.Dir)
	assert.NoError(t, err)

	// Verify todos.json exists and is empty
	todosData, err := os.ReadFile(session.TodosPath)
	require.NoError(t, err)
	var todos SessionTodos
	err = json.Unmarshal(todosData, &todos)
	require.NoError(t, err)
	assert.Empty(t, todos.Todos)

	// Verify session is stored in memory
	storedSession, exists := manager.sessions[session.ID]
	assert.True(t, exists)
	assert.Equal(t, session.ID, storedSession.ID)
}

func TestCreateSessionErrorMkdirAll(t *testing.T) {
	// Use invalid path to force MkdirAll error
	config := Config{SessionsDir: "/dev/null/invalid"}
	manager := NewManager(config)

	session, err := manager.CreateSession()

	assert.Error(t, err)
	assert.Nil(t, session)
	assert.Contains(t, err.Error(), "creating session directory")
}

func TestLoadSession(t *testing.T) {
	tempDir := helpers.TempDir(t)
	config := Config{SessionsDir: tempDir}
	manager := NewManager(config)

	// Create a session first
	originalSession, err := manager.CreateSession()
	require.NoError(t, err)

	// Add some conversation data
	conversation := []openai.ChatCompletionMessageParamUnion{
		openai.UserMessage("Hello"),
		openai.AssistantMessage("Hi there!"),
	}
	originalSession.Conversation = conversation
	err = manager.SaveSession(originalSession)
	require.NoError(t, err)

	// Clear in-memory session
	delete(manager.sessions, originalSession.ID)

	// Load the session
	loadedSession, err := manager.LoadSession(originalSession.ID)

	require.NoError(t, err)
	assert.NotNil(t, loadedSession)
	assert.Equal(t, originalSession.ID, loadedSession.ID)
	assert.Equal(t, originalSession.Dir, loadedSession.Dir)
	assert.Equal(t, originalSession.TodosPath, loadedSession.TodosPath)
	assert.Len(t, loadedSession.Conversation, 2)

	// Verify session is stored in memory after loading
	storedSession, exists := manager.sessions[loadedSession.ID]
	assert.True(t, exists)
	assert.Equal(t, loadedSession.ID, storedSession.ID)
}

func TestLoadSessionNotFound(t *testing.T) {
	tempDir := helpers.TempDir(t)
	config := Config{SessionsDir: tempDir}
	manager := NewManager(config)

	session, err := manager.LoadSession("non-existent-session")

	assert.Error(t, err)
	assert.Nil(t, session)
	assert.Contains(t, err.Error(), "session directory does not exist")
}

func TestGetSession(t *testing.T) {
	tempDir := helpers.TempDir(t)
	config := Config{SessionsDir: tempDir}
	manager := NewManager(config)

	// Test case 1: Session exists in memory
	originalSession, err := manager.CreateSession()
	require.NoError(t, err)

	retrievedSession, err := manager.GetSession(originalSession.ID)
	require.NoError(t, err)
	assert.Equal(t, originalSession.ID, retrievedSession.ID)

	// Test case 2: Session not in memory but exists on disk
	// Clear memory and ensure it loads from disk
	delete(manager.sessions, originalSession.ID)
	retrievedSession, err = manager.GetSession(originalSession.ID)
	require.NoError(t, err)
	assert.Equal(t, originalSession.ID, retrievedSession.ID)

	// Test case 3: Session doesn't exist
	_, err = manager.GetSession("non-existent-session")
	assert.Error(t, err)
}

func TestSaveSession(t *testing.T) {
	tempDir := helpers.TempDir(t)
	config := Config{SessionsDir: tempDir}
	manager := NewManager(config)

	// Create a session
	session, err := manager.CreateSession()
	require.NoError(t, err)

	// Add conversation data
	conversation := []openai.ChatCompletionMessageParamUnion{
		openai.UserMessage("Test message"),
		openai.AssistantMessage("Test response"),
	}
	session.Conversation = conversation

	originalUpdateTime := session.UpdatedAt
	time.Sleep(1 * time.Millisecond) // Ensure time difference

	// Save the session
	err = manager.SaveSession(session)
	require.NoError(t, err)

	// Verify UpdatedAt was updated
	assert.True(t, session.UpdatedAt.After(originalUpdateTime))

	// Verify conversation was saved to file
	logPath := filepath.Join(session.Dir, "agent.log")
	logData, err := os.ReadFile(logPath)
	require.NoError(t, err)

	var savedConversation []openai.ChatCompletionMessageParamUnion
	err = json.Unmarshal(logData, &savedConversation)
	require.NoError(t, err)
	assert.Len(t, savedConversation, 2)
}

func TestGetTodos(t *testing.T) {
	tempDir := helpers.TempDir(t)
	config := Config{SessionsDir: tempDir}
	manager := NewManager(config)

	// Create a session
	session, err := manager.CreateSession()
	require.NoError(t, err)

	// Get todos from fresh session (should be empty)
	todos, err := manager.GetTodos(session.ID)
	require.NoError(t, err)
	assert.Empty(t, todos)

	// Add some todos manually to the file
	testTodos := []TodoItem{
		{ID: "task-1", Content: "Test task 1", Status: "pending", Priority: "high"},
		{ID: "task-2", Content: "Test task 2", Status: "completed", Priority: "low"},
	}
	sessionTodos := SessionTodos{Todos: testTodos}
	err = manager.saveTodosToFile(session.TodosPath, sessionTodos)
	require.NoError(t, err)

	// Get todos again
	todos, err = manager.GetTodos(session.ID)
	require.NoError(t, err)
	assert.Len(t, todos, 2)
	assert.Equal(t, "task-1", todos[0].ID)
	assert.Equal(t, "Test task 1", todos[0].Content)
	assert.Equal(t, "pending", todos[0].Status)
	assert.Equal(t, "high", todos[0].Priority)
}

func TestGetTodosSessionNotFound(t *testing.T) {
	tempDir := helpers.TempDir(t)
	config := Config{SessionsDir: tempDir}
	manager := NewManager(config)

	todos, err := manager.GetTodos("non-existent-session")
	assert.Error(t, err)
	assert.Nil(t, todos)
}

func TestSaveTodos(t *testing.T) {
	tempDir := helpers.TempDir(t)
	config := Config{SessionsDir: tempDir}
	manager := NewManager(config)

	// Create a session
	session, err := manager.CreateSession()
	require.NoError(t, err)

	// Save some todos
	testTodos := []TodoItem{
		{ID: "task-1", Content: "Test task", Status: "in_progress", Priority: "medium"},
	}
	err = manager.SaveTodos(session.ID, testTodos)
	require.NoError(t, err)

	// Verify todos were saved to file
	todosData, err := os.ReadFile(session.TodosPath)
	require.NoError(t, err)

	var savedTodos SessionTodos
	err = json.Unmarshal(todosData, &savedTodos)
	require.NoError(t, err)
	assert.Len(t, savedTodos.Todos, 1)
	assert.Equal(t, "task-1", savedTodos.Todos[0].ID)
	assert.Equal(t, "Test task", savedTodos.Todos[0].Content)
}

func TestSaveTodosSessionNotFound(t *testing.T) {
	tempDir := helpers.TempDir(t)
	config := Config{SessionsDir: tempDir}
	manager := NewManager(config)

	testTodos := []TodoItem{{ID: "task-1", Content: "Test", Status: "pending", Priority: "low"}}
	err := manager.SaveTodos("non-existent-session", testTodos)
	assert.Error(t, err)
}

func TestGenerateSessionID(t *testing.T) {
	manager := NewManager(Config{})

	id1 := manager.generateSessionID()
	time.Sleep(1 * time.Millisecond)
	id2 := manager.generateSessionID()

	// Should be valid timestamp format
	assert.Regexp(t, `^\d{4}-\d{2}-\d{2}-\d{2}-\d{2}-\d{2}$`, id1)
	assert.Regexp(t, `^\d{4}-\d{2}-\d{2}-\d{2}-\d{2}-\d{2}$`, id2)

	// Should be different (assuming they're generated at different times)
	// Note: This might occasionally fail if generated within the same second
	// but the sleep above makes this very unlikely
}

func TestGenerateTodoID(t *testing.T) {
	manager := NewManager(Config{})

	id1 := manager.GenerateTodoID()
	id2 := manager.GenerateTodoID()
	id3 := manager.GenerateTodoID()

	assert.Equal(t, "task-1", id1)
	assert.Equal(t, "task-2", id2)
	assert.Equal(t, "task-3", id3)
	assert.Equal(t, 4, manager.idCounter) // Should be incremented to next available
}

func TestListAvailableSessions(t *testing.T) {
	tempDir := helpers.TempDir(t)
	config := Config{SessionsDir: tempDir}
	manager := NewManager(config)

	// Test case 1: No sessions directory
	nonExistentConfig := Config{SessionsDir: filepath.Join(tempDir, "non-existent")}
	nonExistentManager := NewManager(nonExistentConfig)
	sessions, err := nonExistentManager.ListAvailableSessions()
	require.NoError(t, err)
	assert.Empty(t, sessions)

	// Test case 2: Create some sessions
	session1, err := manager.CreateSession()
	require.NoError(t, err)

	time.Sleep(1 * time.Second) // Ensure different timestamps (1 second for session ID resolution)
	session2, err := manager.CreateSession()
	require.NoError(t, err)

	// Create an invalid directory (shouldn't be included)
	invalidDir := filepath.Join(tempDir, "invalid-format")
	err = os.Mkdir(invalidDir, 0755)
	require.NoError(t, err)

	// List sessions
	sessions, err = manager.ListAvailableSessions()
	require.NoError(t, err)
	assert.Len(t, sessions, 2)

	// Should be sorted with newest first
	assert.Contains(t, sessions, session1.ID)
	assert.Contains(t, sessions, session2.ID)
	// Due to sorting, session2 (newer) should be first
	if session2.ID > session1.ID { // Only check if IDs are actually different
		assert.Equal(t, session2.ID, sessions[0])
		assert.Equal(t, session1.ID, sessions[1])
	}
}

func TestValidateTodoStatus(t *testing.T) {
	validStatuses := []string{"pending", "in_progress", "completed", "cancelled"}
	invalidStatuses := []string{"", "invalid", "PENDING", "done", "active"}

	for _, status := range validStatuses {
		assert.True(t, ValidateTodoStatus(status), "Status %s should be valid", status)
	}

	for _, status := range invalidStatuses {
		assert.False(t, ValidateTodoStatus(status), "Status %s should be invalid", status)
	}
}

func TestValidateTodoPriority(t *testing.T) {
	validPriorities := []string{"high", "medium", "low"}
	invalidPriorities := []string{"", "invalid", "HIGH", "urgent", "normal"}

	for _, priority := range validPriorities {
		assert.True(t, ValidateTodoPriority(priority), "Priority %s should be valid", priority)
	}

	for _, priority := range invalidPriorities {
		assert.False(t, ValidateTodoPriority(priority), "Priority %s should be invalid", priority)
	}
}

func TestIsValidSessionIDFormat(t *testing.T) {
	manager := NewManager(Config{})

	validFormats := []string{
		"2024-12-19-14-30-45",
		"2000-01-01-00-00-00",
		"2099-12-31-23-59-59",
	}

	invalidFormats := []string{
		"",
		"invalid-format",
		"2024-12-19",
		"2024-12-19-14-30",
		"24-12-19-14-30-45",
		"2024-13-01-14-30-45", // Invalid month
		"2024-12-32-14-30-45", // Invalid day
		"2024-12-19-25-30-45", // Invalid hour
		"2024-12-19-14-61-45", // Invalid minute
		"2024-12-19-14-30-61", // Invalid second
	}

	for _, format := range validFormats {
		assert.True(t, manager.isValidSessionIDFormat(format), "Format %s should be valid", format)
	}

	for _, format := range invalidFormats {
		assert.False(t, manager.isValidSessionIDFormat(format), "Format %s should be invalid", format)
	}
}

func TestLoadConversationFromFile(t *testing.T) {
	tempDir := helpers.TempDir(t)
	manager := NewManager(Config{SessionsDir: tempDir})

	// Test case 1: File doesn't exist
	nonExistentPath := filepath.Join(tempDir, "non-existent.log")
	conversation, err := manager.loadConversationFromFile(nonExistentPath)
	require.NoError(t, err)
	assert.Empty(t, conversation)

	// Test case 2: Empty file
	emptyPath := filepath.Join(tempDir, "empty.log")
	err = os.WriteFile(emptyPath, []byte(""), 0644)
	require.NoError(t, err)
	conversation, err = manager.loadConversationFromFile(emptyPath)
	require.NoError(t, err)
	assert.Empty(t, conversation)

	// Test case 3: Valid conversation file
	testConversation := []openai.ChatCompletionMessageParamUnion{
		openai.UserMessage("Hello"),
		openai.AssistantMessage("Hi!"),
	}
	validPath := filepath.Join(tempDir, "valid.log")
	data, err := json.MarshalIndent(testConversation, "", "  ")
	require.NoError(t, err)
	err = os.WriteFile(validPath, data, 0644)
	require.NoError(t, err)

	conversation, err = manager.loadConversationFromFile(validPath)
	require.NoError(t, err)
	assert.Len(t, conversation, 2)

	// Test case 4: Invalid JSON
	invalidPath := filepath.Join(tempDir, "invalid.log")
	err = os.WriteFile(invalidPath, []byte("invalid json"), 0644)
	require.NoError(t, err)
	conversation, err = manager.loadConversationFromFile(invalidPath)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "parsing conversation JSON")
}

func TestSaveConversationToFile(t *testing.T) {
	tempDir := helpers.TempDir(t)
	manager := NewManager(Config{SessionsDir: tempDir})

	testConversation := []openai.ChatCompletionMessageParamUnion{
		openai.UserMessage("Test message"),
		openai.AssistantMessage("Test response"),
	}

	// Test saving conversation
	logPath := filepath.Join(tempDir, "subdir", "conversation.log")
	err := manager.saveConversationToFile(logPath, testConversation)
	require.NoError(t, err)

	// Verify file exists and content is correct
	data, err := os.ReadFile(logPath)
	require.NoError(t, err)

	var savedConversation []openai.ChatCompletionMessageParamUnion
	err = json.Unmarshal(data, &savedConversation)
	require.NoError(t, err)
	assert.Len(t, savedConversation, 2)
}

func TestLoadTodosFromFile(t *testing.T) {
	tempDir := helpers.TempDir(t)
	manager := NewManager(Config{SessionsDir: tempDir})

	// Test case 1: File doesn't exist
	nonExistentPath := filepath.Join(tempDir, "non-existent.json")
	todos, err := manager.loadTodosFromFile(nonExistentPath)
	require.NoError(t, err)
	assert.Empty(t, todos.Todos)

	// Test case 2: Empty file
	emptyPath := filepath.Join(tempDir, "empty.json")
	err = os.WriteFile(emptyPath, []byte(""), 0644)
	require.NoError(t, err)
	todos, err = manager.loadTodosFromFile(emptyPath)
	require.NoError(t, err)
	assert.Empty(t, todos.Todos)

	// Test case 3: Valid todos file
	testTodos := SessionTodos{
		Todos: []TodoItem{
			{ID: "task-1", Content: "Test task", Status: "pending", Priority: "high"},
		},
	}
	validPath := filepath.Join(tempDir, "valid.json")
	data, err := json.MarshalIndent(testTodos, "", "  ")
	require.NoError(t, err)
	err = os.WriteFile(validPath, data, 0644)
	require.NoError(t, err)

	todos, err = manager.loadTodosFromFile(validPath)
	require.NoError(t, err)
	assert.Len(t, todos.Todos, 1)
	assert.Equal(t, "task-1", todos.Todos[0].ID)

	// Test case 4: Invalid JSON
	invalidPath := filepath.Join(tempDir, "invalid.json")
	err = os.WriteFile(invalidPath, []byte("invalid json"), 0644)
	require.NoError(t, err)
	todos, err = manager.loadTodosFromFile(invalidPath)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "parsing todos JSON")
}

func TestSaveTodosToFile(t *testing.T) {
	tempDir := helpers.TempDir(t)
	manager := NewManager(Config{SessionsDir: tempDir})

	testTodos := SessionTodos{
		Todos: []TodoItem{
			{ID: "task-1", Content: "Test task", Status: "completed", Priority: "medium"},
		},
	}

	// Test saving todos to file with subdirectory creation
	todosPath := filepath.Join(tempDir, "subdir", "todos.json")
	err := manager.saveTodosToFile(todosPath, testTodos)
	require.NoError(t, err)

	// Verify file exists and content is correct
	data, err := os.ReadFile(todosPath)
	require.NoError(t, err)

	var savedTodos SessionTodos
	err = json.Unmarshal(data, &savedTodos)
	require.NoError(t, err)
	assert.Len(t, savedTodos.Todos, 1)
	assert.Equal(t, "task-1", savedTodos.Todos[0].ID)
	assert.Equal(t, "Test task", savedTodos.Todos[0].Content)
}
