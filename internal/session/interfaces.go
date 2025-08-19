package session

// SessionManager interface defines the contract for session management
type SessionManager interface {
	CreateSession() (*Session, error)
	LoadSession(id string) (*Session, error)
	SaveSession(session *Session) error
	SelectSessionInteractively() (string, error)
	GetTodos(sessionID string) ([]TodoItem, error)
	SaveTodos(sessionID string, todos []TodoItem) error
	GenerateTodoID() string
}

// Ensure Manager implements SessionManager interface
var _ SessionManager = (*Manager)(nil)
