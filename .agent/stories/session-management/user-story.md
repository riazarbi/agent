# Session-Based Logging and Resume Feature

*Enable per-session data persistence (logs and todos) with the ability to resume previous sessions, replacing the current global logging approach.*

## Requirements

- Generate human-readable session timestamps in format `YYYY-MM-DD-HH-MM-SS` for each agent run
- Set global `currentSessionID` variable at startup with this timestamp
- Create session directories at `.agent/sessions/[currentSessionID]/` containing:
  - `agent.log` - conversation history in current JSON message format  
  - `todos.json` - persisted todo state from todowrite/todoread tools
- Implement `-resume [sessionID]` CLI flag to resume a specific session
- Optional: Support `-resume` without arguments to show available sessions interactively
- When resuming: 
  - Load conversation array from session's `agent.log` via JSON unmarshal
  - Load todo state from session's `todos.json`
  - Set `currentSessionID` to the resumed session (don't create new session)
  - Continue logging to that session's files
- When not resuming: create new session directory and start fresh logging
- Remove dependency on global `.agent/agent.log` file for new sessions

## Rules

- Session timestamp must be generated at agent startup (not per-message)
- Session directories must be created immediately when agent starts
- Message logging must write to session-specific log file in real-time
- Todo state must be persisted to session-specific JSON file on every todowrite operation
- Resume must fail gracefully if selected session doesn't exist or is corrupted
- No migration of existing global logs - start fresh with session-based approach

## Domain

```go
// Core session management
var currentSessionID string // Set at startup: "2024-12-19-14-30-45"
var currentSessionDir string // Derived: ".agent/sessions/[currentSessionID]/"

type SessionManager struct {
    SessionID   string // Format: "2024-12-19-14-30-45"  
    SessionDir  string // Path: ".agent/sessions/[timestamp]/"
    LogFile     *os.File
    TodosPath   string
    Conversation []openai.ChatCompletionMessageParamUnion // Loaded from previous session or empty
}

// Resume functionality
func loadSessionConversation(sessionID string) ([]openai.ChatCompletionMessageParamUnion, error)
func loadSessionTodos(sessionID string) ([]TodoItem, error)

// Session directory structure:
// .agent/sessions/
//   ├── 2024-12-19-14-30-45/
//   │   ├── agent.log          # JSON array of ChatCompletionMessageParamUnion
//   │   └── todos.json         # JSON: {"todos": [...]}
//   ├── 2024-12-19-10-15-22/
//   │   ├── agent.log
//   │   └── todos.json
```

## Extra Considerations

- Handle filesystem permissions for creating session directories
- **Critical**: Graceful handling of corrupted JSON in session files - don't crash on resume
- Consider disk space - no automatic cleanup implemented in this story  
- Resume selection should be case-insensitive for timestamp input
- Empty sessions (no messages/todos) still count as valid resumable sessions
- **JSON compatibility**: Current `logConversation()` output is already compatible with `[]openai.ChatCompletionMessageParamUnion`
- **Session continuity**: When resuming, `currentSessionID` should be set to the resumed session, not create a new one

## Testing Considerations

- Test session directory creation on agent startup
- Test message logging to correct session file
- Test todo persistence to correct session file
- Test resume flow: list sessions, select valid session, load state
- Test resume failure cases: invalid timestamp, corrupted files, missing directories
- Test concurrent session creation (multiple agent instances)

## Implementation Notes

### Core Implementation Approach
- Set `sessionID` variable at startup using current timestamp in `YYYY-MM-DD-HH-MM-SS` format
- Use `sessionID` for session directory path: `.agent/sessions/[sessionID]/`
- For resume: accept sessionID as argument, no complex session selection UI needed
- Modify `logConversation()` to write to session-specific `agent.log` instead of global log

### Message History Restoration - CRITICAL ASSESSMENT

**Option 1: Reconstruct conversation array from JSON log (RECOMMENDED)**
- The current `logConversation()` already outputs proper JSON: `[]openai.ChatCompletionMessageParamUnion`
- **Feasibility**: HIGH - JSON structure matches Go types exactly
- **Implementation**: Unmarshal JSON back to `[]openai.ChatCompletionMessageParamUnion`
- **Benefits**: Clean, maintains proper OpenAI API structure, no context pollution
- **Risk**: JSON parsing errors on corrupted logs (handle gracefully)

```go
// Pseudocode for loading conversation
func loadConversationFromSession(sessionDir string) ([]openai.ChatCompletionMessageParamUnion, error) {
    data, err := os.ReadFile(filepath.Join(sessionDir, "agent.log"))
    if err != nil { return nil, err }
    
    var conversation []openai.ChatCompletionMessageParamUnion
    err = json.Unmarshal(data, &conversation)
    return conversation, err
}
```

**Option 2: Insert entire logfile as first message (NOT RECOMMENDED)**
- **Feasibility**: EASY but problematic
- **Issues**: 
  - Large logs could exceed token limits
  - Semantically incorrect - breaks conversation structure
  - May confuse model with raw JSON in message content
  - Loses proper message typing (user/assistant/tool)

**RECOMMENDATION**: Use Option 1. The JSON structure is already compatible and maintains clean separation.

## Specification by Example

### CLI Usage
```bash
# Start new session
$ ./agent
# Creates: .agent/sessions/2024-12-19-14-30-45/

# Resume specific session
$ ./agent -resume 2024-12-19-10-15-22
# Loads messages and todos from .agent/sessions/2024-12-19-10-15-22/
# Continues logging to that same session directory

# Resume with discovery (alternative approach)
$ ./agent -resume
Available sessions:
2024-12-19-14-30-45
2024-12-19-10-15-22
2024-12-18-16-45-12
Select session to resume: 2024-12-19-10-15-22
```

### Session Directory Structure
```
.agent/sessions/2024-12-19-14-30-45/
├── agent.log          # Same JSON message format as current global log
└── todos.json         # Persisted todo state: {"todos": [...]}
```

### Todos.json Format
```json
{
  "todos": [
    {
      "id": "task-001",
      "content": "Implement user authentication", 
      "status": "in_progress",
      "priority": "high"
    }
  ]
}
```

## Verification

- [ ] Session timestamp generated in YYYY-MM-DD-HH-MM-SS format on startup
- [ ] Session directory created at .agent/sessions/[timestamp]/ on startup
- [ ] Messages logged to session-specific agent.log file in real-time
- [ ] Todo state persisted to session-specific todos.json on every todowrite
- [ ] `-resume` flag shows available sessions in descending chronological order
- [ ] Resume flow loads both message history and todo state correctly
- [ ] Resume continues logging to selected session's files
- [ ] Graceful error handling for invalid/missing sessions
- [ ] No writes to global .agent/agent.log for new sessions
- [ ] Interactive session selection works correctly

## Next Steps

After creating the story:
1. Save to `.agent/stories/session-management/user-story.md`