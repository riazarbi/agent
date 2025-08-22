package agent

import (
	"agent/internal/config"
	"agent/internal/session"
	"os"
	"path/filepath"
	"testing"
	"time"
)

// MockSessionManager is a mock implementation of the SessionManager interface
type MockSessionManager struct {
	sessions    map[string]*session.Session
	sessionsDir string
}

// NewMockSessionManager creates a new MockSessionManager
func NewMockSessionManager(sessionsDir string) *MockSessionManager {
	return &MockSessionManager{
		sessions:    make(map[string]*session.Session),
		sessionsDir: sessionsDir,
	}
}

// CreateSession creates a new mock session
func (m *MockSessionManager) CreateSession() (*session.Session, error) {
	sessionID := time.Now().Format("2006-01-02-15-04-05")
	sessionDir := filepath.Join(m.sessionsDir, sessionID)
	if err := os.MkdirAll(sessionDir, 0755); err != nil {
		return nil, err
	}
	s := &session.Session{
		ID:  sessionID,
		Dir: sessionDir,
	}
	m.sessions[sessionID] = s
	return s, nil
}

// LoadSession loads a mock session
func (m *MockSessionManager) LoadSession(id string) (*session.Session, error) {
	if s, ok := m.sessions[id]; ok {
		return s, nil
	}
	return nil, os.ErrNotExist
}

// SaveSession saves a mock session
func (m *MockSessionManager) SaveSession(s *session.Session) error {
	m.sessions[s.ID] = s
	return nil
}

// SelectSessionInteractively is not implemented for the mock
func (m *MockSessionManager) SelectSessionInteractively() (string, error) {
	return "", nil
}

// GetTodos is not implemented for the mock
func (m *MockSessionManager) GetTodos(sessionID string) ([]session.TodoItem, error) {
	return nil, nil
}

// SaveTodos is not implemented for the mock
func (m *MockSessionManager) SaveTodos(sessionID string, todos []session.TodoItem) error {
	return nil
}

// GenerateTodoID is not implemented for the mock
func (m *MockSessionManager) GenerateTodoID() string {
	return "mock-todo-id"
}

func TestNew(t *testing.T) {
	cfg := &config.Config{
		Session: config.SessionConfig{
			Dir: "/tmp/sessions",
		},
	}
	deps := Dependencies{}

	agent, err := New(cfg, deps)

	if err != nil {
		t.Fatalf("New() error = %v, wantErr %v", err, false)
	}

	if agent == nil {
		t.Fatal("New() returned nil agent")
	}

	if agent.config != cfg {
		t.Error("agent.config not set correctly")
	}
}

func TestSetFlags(t *testing.T) {
	agent := &Agent{}
	flags := Flags{
		PromptFile: "test.prompt",
	}

	agent.SetFlags(flags)

	if agent.flags.PromptFile != "test.prompt" {
		t.Error("agent.flags not set correctly")
	}
}

func TestSetTransitionToInteractive(t *testing.T) {
	agent := &Agent{}

	agent.SetTransitionToInteractive(true)

	if !agent.transitionToInteractive {
		t.Error("agent.transitionToInteractive not set correctly")
	}
}

func TestInitializeSession(t *testing.T) {
	sessionsDir := t.TempDir()
	mockSessionManager := NewMockSessionManager(sessionsDir)
	agent := &Agent{
		sessionManager: mockSessionManager,
	}

	t.Run("Create new session", func(t *testing.T) {
		agent.SetFlags(Flags{})
		err := agent.InitializeSession()
		if err != nil {
			t.Fatalf("InitializeSession() error = %v, wantErr %v", err, false)
		}
		if agent.currentSession == nil {
			t.Fatal("currentSession not set")
		}
	})

	t.Run("Resume session", func(t *testing.T) {
		// Create a session to resume
		s, _ := mockSessionManager.CreateSession()

		agent.SetFlags(Flags{ResumeSession: s.ID})
		err := agent.InitializeSession()
		if err != nil {
			t.Fatalf("InitializeSession() error = %v, wantErr %v", err, false)
		}
		if agent.currentSession == nil {
			t.Fatal("currentSession not set")
		}
		if agent.currentSession.ID != s.ID {
			t.Errorf("currentSession.ID = %s, want %s", agent.currentSession.ID, s.ID)
		}
	})

	t.Run("Resume non-existent session", func(t *testing.T) {
		agent.SetFlags(Flags{ResumeSession: "non-existent-session"})
		err := agent.InitializeSession()
		if err == nil {
			t.Fatal("InitializeSession() error = nil, wantErr not nil")
		}
	})
}
