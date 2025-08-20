# Dynamic System Tools Integration

This story addresses the challenge of integrating new system commands (like `mv`, `rm`, `touch`, `wc`, and future tools like `ssh`) into the agent's capabilities without introducing excessive code complexity or rigidity. The main objective is to establish a dynamic, configuration-driven mechanism for adding these tools, ensuring they behave as direct "passthrough wrappers" to their underlying Linux commands.

## Past Attempts

If this user story has been attempted before, the changes made will appear in the git diff. Our policy is to only make a single commit per user story, so you can always review the git diff to review progress across attempts. 


## Requirements

-   The agent must be able to access `mv`, `rm`, `touch`, and `wc` as callable tools.
-   Adding a new system tool (e.g., `ssh`) to the agent's capabilities must be achievable by modifying a single, central configuration or list or similar data structure within the codebase, without requiring changes to existing tool registration logic or additional tool files.
-   Each integrated system tool must function as a direct "passthrough wrapper," meaning its arguments, behavior, output (stdout/stderr), and error codes must exactly mirror the behavior of the corresponding underlying Linux command.
-   The output of these system tools should be presented to the agent in a clear and usable format, consistent with how other tool outputs are handled.
-   Error conditions for these tools must propagate accurately from the underlying system command to the agent.

## Rules

-   The implementation must prioritize minimizing boilerplate code for each new system tool.
-   The solution should not re-implement the logic of standard Linux commands but rather wrap them directly.
-   Security considerations for executing arbitrary system commands should be taken into account (e.g., preventing shell injection, though the primary focus is on passthrough).

## Domain

```
// Simplified representation of how system tools might be registered
// internal/tools/registry.go (or new system_tools.go)

type SystemTool struct {
    Name    string
    Command string // The underlying Linux command (e.g., "mv", "rm")
}

// A collection of system tools to be registered dynamically
var RegisteredSystemTools = []SystemTool{
    {Name: "mv", Command: "mv"},
    {Name: "rm", Command: "rm"},
    {Name: "touch", Command: "touch"},
    {Name: "wc", Command: "wc"},
    // Future tools like ssh would be added here
}

// Function to dynamically create and register tools based on RegisteredSystemTools
// func NewSystemTools() []ToolDefinition { ... }
```

## Extra Considerations

-   Consider how tool arguments will be parsed and passed to the underlying system command, especially for commands with complex flag structures.
-   Think about potential conflicts with existing tools if a system command name overlaps.
-   The solution should be robust to cases where a system command is not found on the execution environment.

## Testing Considerations

**YOU CANNOT TEST THESE NEW TOOLS, A NEW BINARY MUST BE BUILT FIRST. WRITE INSTRUCTIONS FOR TESTING TO A FILE CALLED check.tct, OVERWRITING PREVIOUS CONTENT**


-   **Unit Tests:** Verify the dynamic registration mechanism correctly processes `RegisteredSystemTools` and creates the expected tool definitions.
-   **Integration Tests:**
    -   Test `mv` for successful file/directory moves, error cases (e.g., no source, destination exists).
    -   Test `rm` for successful file/directory deletion, recursive deletion, error cases (e.g., file not found).
    -   Test `touch` for creating new files and updating timestamps of existing files.
    -   Test `wc` for line, word, and character counts on various file contents.
    -   Verify that the arguments passed to the agent's tool are correctly forwarded to the underlying system command.
    -   Verify that stdout, stderr, and exit codes from the system commands are accurately captured and returned by the agent's tool.
-   **Edge Cases:** Test with unusual characters in filenames, very large files (for `wc`), and non-existent paths.

## Implementation Notes

-   Explore Go's `os/exec` package for executing external commands.
-   Consider how to handle the `args` parameter for the tool function, potentially as a single string that gets parsed or directly passed.
-   The `registry.go` file (or a new `system_tools.go`) would be the logical place to house the `RegisteredSystemTools` list and the function to generate `ToolDefinition` objects from it.

## Specification by Example

### `mv` tool
-   `mv(args="old_file.txt new_file.txt")` -> Should rename `old_file.txt` to `new_file.txt`.
-   `mv(args="my_dir/ new_dir/")` -> Should move `my_dir/` to `new_dir/`.

### `rm` tool
-   `rm(args="file_to_delete.txt")` -> Should delete `file_to_delete.txt`.
-   `rm(args="-r dir_to_delete/")` -> Should recursively delete `dir_to_delete/`.

### `touch` tool
-   `touch(args="new_empty_file.txt")` -> Should create `new_empty_file.txt` if it doesn't exist, or update its timestamp if it does.

### `wc` tool
-   `wc(args="-l file.txt")` -> Should return the line count of `file.txt`.
-   `wc(args="file.txt")` -> Should return line, word, and character counts of `file.txt`.

## Verification

-   [ ] Verify `mv` tool correctly renames files and directories.
-   [ ] Verify `rm` tool correctly deletes files and directories (including recursive).
-   [ ] Verify `touch` tool correctly creates new files and updates timestamps.
-   [ ] Verify `wc` tool correctly returns line, word, and character counts.
-   [ ] Verify that a new system tool (e.g., `ssh` conceptually, without actual `ssh` execution) can be added to the `RegisteredSystemTools` list and is dynamically registered as an available tool.
-   [ ] Verify that arguments and output for `mv`, `rm`, `touch`, `wc` precisely match their native Linux command counterparts.