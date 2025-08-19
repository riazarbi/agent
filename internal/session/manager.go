package session

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"sync"
	"time"

	"github.com/openai/openai-go/v2"
)

// TodoItem represents a single todo item with all required fields
type TodoItem struct {
	ID       string `json:"id"`
	Content  string `json:"content"`
	Status   string `json:"status"`   // pending, in_progress, completed, cancelled
	Priority string `json:"priority"` // high, medium, low
}

// Session represents a single agent session with conversation and todos
type Session struct {
	ID           string                                   `json:"session_id"`   // Format: "2024-12-19-14-30-45"
	Dir          string                                   `json:"session_dir"`  // Path: ".agent/sessions/[timestamp]/"
	TodosPath    string                                   `json:"todos_path"`   // Path to todos.json
	Conversation []openai.ChatCompletionMessageParamUnion `json:"conversation"` // Conversation history
	CreatedAt    time.Time                                `json:"created_at"`
	UpdatedAt    time.Time                                `json:"updated_at"`
}

// SessionTodos is the file format for persisting todos
type SessionTodos struct {
	Todos []TodoItem `json:"todos"`
}

// Manager handles session creation, loading, and persistence
type Manager struct {
	config       Config
	sessions     map[string]*Session
	sessionMutex sync.RWMutex
	idCounter    int
}

// Config contains session management configuration
type Config struct {
	SessionsDir string
}

// NewManager creates a new session manager
func NewManager(config Config) *Manager {
	return &Manager{
		config:       config,
		sessions:     make(map[string]*Session),
		sessionMutex: sync.RWMutex{},
		idCounter:    1,
	}
}

// CreateSession creates a new session with a generated ID
func (m *Manager) CreateSession() (*Session, error) {
	sessionID := m.generateSessionID()
	sessionDir := filepath.Join(m.config.SessionsDir, sessionID)

	if err := os.MkdirAll(sessionDir, 0755); err != nil {
		return nil, fmt.Errorf("creating session directory: %w", err)
	}

	session := &Session{
		ID:           sessionID,
		Dir:          sessionDir,
		TodosPath:    filepath.Join(sessionDir, "todos.json"),
		Conversation: []openai.ChatCompletionMessageParamUnion{},
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	// Create empty todos.json
	emptyTodos := SessionTodos{Todos: []TodoItem{}}
	if err := m.saveTodosToFile(session.TodosPath, emptyTodos); err != nil {
		return nil, fmt.Errorf("creating todos.json: %w", err)
	}

	// Store in memory
	m.sessionMutex.Lock()
	m.sessions[sessionID] = session
	m.sessionMutex.Unlock()

	return session, nil
}

// LoadSession loads an existing session by ID
func (m *Manager) LoadSession(sessionID string) (*Session, error) {
	sessionDir := filepath.Join(m.config.SessionsDir, sessionID)

	// Check if session directory exists
	if _, err := os.Stat(sessionDir); os.IsNotExist(err) {
		return nil, fmt.Errorf("session directory does not exist: %s", sessionID)
	}

	session := &Session{
		ID:        sessionID,
		Dir:       sessionDir,
		TodosPath: filepath.Join(sessionDir, "todos.json"),
		CreatedAt: time.Now(), // Will be overwritten if we add metadata file later
		UpdatedAt: time.Now(),
	}

	// Load conversation from agent.log
	logPath := filepath.Join(sessionDir, "agent.log")
	conversation, err := m.loadConversationFromFile(logPath)
	if err != nil {
		return nil, fmt.Errorf("loading conversation: %w", err)
	}
	session.Conversation = conversation

	// Store in memory
	m.sessionMutex.Lock()
	m.sessions[sessionID] = session
	m.sessionMutex.Unlock()

	return session, nil
}

// GetSession retrieves a session from memory or loads it
func (m *Manager) GetSession(sessionID string) (*Session, error) {
	m.sessionMutex.RLock()
	session, exists := m.sessions[sessionID]
	m.sessionMutex.RUnlock()

	if exists {
		return session, nil
	}

	// Try to load from disk
	return m.LoadSession(sessionID)
}

// SaveSession persists session data to disk
func (m *Manager) SaveSession(session *Session) error {
	session.UpdatedAt = time.Now()

	// Save conversation to agent.log
	logPath := filepath.Join(session.Dir, "agent.log")
	if err := m.saveConversationToFile(logPath, session.Conversation); err != nil {
		return fmt.Errorf("saving conversation: %w", err)
	}

	return nil
}

// GetTodos retrieves todos for a session
func (m *Manager) GetTodos(sessionID string) ([]TodoItem, error) {
	session, err := m.GetSession(sessionID)
	if err != nil {
		return nil, err
	}

	todos, err := m.loadTodosFromFile(session.TodosPath)
	if err != nil {
		return nil, fmt.Errorf("loading todos: %w", err)
	}

	return todos.Todos, nil
}

// SaveTodos persists todos for a session
func (m *Manager) SaveTodos(sessionID string, todos []TodoItem) error {
	session, err := m.GetSession(sessionID)
	if err != nil {
		return err
	}

	sessionTodos := SessionTodos{Todos: todos}
	if err := m.saveTodosToFile(session.TodosPath, sessionTodos); err != nil {
		return fmt.Errorf("saving todos: %w", err)
	}

	return nil
}

// GenerateSessionID creates a new session ID with timestamp format
func (m *Manager) generateSessionID() string {
	return time.Now().Format("2006-01-02-15-04-05")
}

// GenerateTodoID creates a unique ID for todo items
func (m *Manager) GenerateTodoID() string {
	m.sessionMutex.Lock()
	defer m.sessionMutex.Unlock()

	id := "task-" + strconv.Itoa(m.idCounter)
	m.idCounter++
	return id
}

// ListAvailableSessions returns all available session IDs
func (m *Manager) ListAvailableSessions() ([]string, error) {
	// Check if sessions directory exists
	if _, err := os.Stat(m.config.SessionsDir); os.IsNotExist(err) {
		return []string{}, nil
	}

	entries, err := os.ReadDir(m.config.SessionsDir)
	if err != nil {
		return nil, fmt.Errorf("reading sessions directory: %w", err)
	}

	var sessions []string
	for _, entry := range entries {
		if entry.IsDir() {
			// Validate session ID format (YYYY-MM-DD-HH-MM-SS)
			sessionID := entry.Name()
			if m.isValidSessionIDFormat(sessionID) {
				sessions = append(sessions, sessionID)
			}
		}
	}

	// Sort sessions in descending order (newest first)
	sort.Slice(sessions, func(i, j int) bool {
		return sessions[i] > sessions[j]
	})

	return sessions, nil
}

// SelectSessionInteractively prompts user to select from available sessions
func (m *Manager) SelectSessionInteractively() (string, error) {
	sessions, err := m.ListAvailableSessions()
	if err != nil {
		return "", fmt.Errorf("listing available sessions: %w", err)
	}

	if len(sessions) == 0 {
		fmt.Println("No previous sessions found. Creating new session...")
		return "", nil
	}

	fmt.Println("Available sessions:")
	for i, session := range sessions {
		fmt.Printf("  %d) %s\n", i+1, session)
	}
	fmt.Printf("  %d) Create new session\n", len(sessions)+1)
	fmt.Printf("\nSelect a session (1-%d): ", len(sessions)+1)

	var choice int
	if _, err := fmt.Scanf("%d", &choice); err != nil {
		return "", fmt.Errorf("invalid input: %w", err)
	}

	if choice < 1 || choice > len(sessions)+1 {
		return "", fmt.Errorf("invalid choice: %d", choice)
	}

	if choice == len(sessions)+1 {
		// Create new session
		return "", nil
	}

	return sessions[choice-1], nil
}

// ValidateTodoStatus checks if status is valid
func ValidateTodoStatus(status string) bool {
	validStatuses := map[string]bool{
		"pending":     true,
		"in_progress": true,
		"completed":   true,
		"cancelled":   true,
	}
	return validStatuses[status]
}

// ValidateTodoPriority checks if priority is valid
func ValidateTodoPriority(priority string) bool {
	validPriorities := map[string]bool{
		"high":   true,
		"medium": true,
		"low":    true,
	}
	return validPriorities[priority]
}

// Private helper methods

func (m *Manager) isValidSessionIDFormat(sessionID string) bool {
	// Check if format matches YYYY-MM-DD-HH-MM-SS
	_, err := time.Parse("2006-01-02-15-04-05", sessionID)
	return err == nil
}

func (m *Manager) loadConversationFromFile(logPath string) ([]openai.ChatCompletionMessageParamUnion, error) {
	// Check if file exists
	if _, err := os.Stat(logPath); os.IsNotExist(err) {
		// Return empty conversation if file doesn't exist (new session)
		return []openai.ChatCompletionMessageParamUnion{}, nil
	}

	data, err := os.ReadFile(logPath)
	if err != nil {
		return nil, fmt.Errorf("reading conversation file: %w", err)
	}

	// Handle empty file
	if len(data) == 0 {
		return []openai.ChatCompletionMessageParamUnion{}, nil
	}

	var conversation []openai.ChatCompletionMessageParamUnion
	if err := json.Unmarshal(data, &conversation); err != nil {
		return nil, fmt.Errorf("parsing conversation JSON: %w", err)
	}

	return conversation, nil
}

func (m *Manager) saveConversationToFile(logPath string, conversation []openai.ChatCompletionMessageParamUnion) error {
	// Ensure directory exists
	if err := os.MkdirAll(filepath.Dir(logPath), 0755); err != nil {
		return fmt.Errorf("creating log directory: %w", err)
	}

	data, err := json.MarshalIndent(conversation, "", "  ")
	if err != nil {
		return fmt.Errorf("marshaling conversation: %w", err)
	}

	return os.WriteFile(logPath, data, 0644)
}

func (m *Manager) loadTodosFromFile(todosPath string) (SessionTodos, error) {
	// Check if file exists
	if _, err := os.Stat(todosPath); os.IsNotExist(err) {
		// Return empty todos if file doesn't exist
		return SessionTodos{Todos: []TodoItem{}}, nil
	}

	data, err := os.ReadFile(todosPath)
	if err != nil {
		return SessionTodos{}, fmt.Errorf("reading todos file: %w", err)
	}

	// Handle empty file
	if len(data) == 0 {
		return SessionTodos{Todos: []TodoItem{}}, nil
	}

	var todos SessionTodos
	if err := json.Unmarshal(data, &todos); err != nil {
		return SessionTodos{}, fmt.Errorf("parsing todos JSON: %w", err)
	}

	return todos, nil
}

func (m *Manager) saveTodosToFile(todosPath string, todos SessionTodos) error {
	data, err := json.MarshalIndent(todos, "", "  ")
	if err != nil {
		return fmt.Errorf("marshaling todos: %w", err)
	}

	// Ensure directory exists
	if err := os.MkdirAll(filepath.Dir(todosPath), 0755); err != nil {
		return fmt.Errorf("creating todos directory: %w", err)
	}

	return os.WriteFile(todosPath, data, 0644)
}
