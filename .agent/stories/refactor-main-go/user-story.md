# Refactor main.go into Smaller Modules

*Refactor the main.go file into smaller, more maintainable modules to improve code organization and reduce complexity.*

## Requirements

*   Create a new file named `agent.go` containing the `Agent` struct, its methods, and related types.
*   Create a new file named `session.go` containing the `SessionManager` struct, its methods, and related functions for session management.
*   Create a new file named `tools.go` containing the definitions and implementations of the tool functions and their input types.
*   Create a new file named `utils.go` containing utility functions.
*   Modify the existing `main.go` file to import and use the new modules, containing the `main` function, command-line argument parsing, agent initialization, and the main interaction loop.
*   Ensure all existing functionality remains intact after the refactoring.
*   The new `main.go` should be less than 200 lines.
*   All the new files should compile.

## Rules

*   The refactoring must adhere to Go coding standards and best practices.
*   No external dependencies should be introduced.

## Domain

```
// Agent: Represents the main agent with its OpenAI client, tools, and session.
type Agent struct { ... }

// SessionManager: Manages the agent's session, including conversation history and todos.
type SessionManager struct { ... }

// ToolDefinition: Defines a tool that the agent can use.
type ToolDefinition struct { ... }
```

## Extra Considerations

*   Ensure proper error handling in the new modules.
*   Consider adding comments to the new modules to improve readability.
*   The split should follow logical code boundaries

## Testing Considerations

*   Unit tests should be created for the new modules to verify their functionality.
*   Integration tests should be performed to ensure that all components work together correctly.

## Implementation Notes

*   Use Go modules to manage dependencies.

## Specification by Example

*   N/A

## Verification

*   [ ] `agent.go` file exists and contains the `Agent` struct and related code.
*   [ ] `session.go` file exists and contains the `SessionManager` struct and related code.
*   [ ] `tools.go` file exists and contains the tool definitions and implementations.
*   [ ] `utils.go` file exists and contains the utility functions.
*   [ ] `main.go` file exists and contains the `main` function and related code.
*   [ ] All files compile without errors.
*   [ ] All existing functionality works as expected.
*   [ ] `main.go` is less than 200 lines

## Next Steps

1.  Save to `.agent/stories/refactor-main-go/user-story.md`