# Rename `grep` tool to `rg` for `ripgrep` consistency

*Brief description of the problem this solves and main objective.*

The `grep` tool is currently named inconsistently with its underlying `ripgrep` (`rg`) implementation. This naming mismatch can lead to agent confusion regarding the expected syntax and arguments. The main objective is to rename the `grep` tool to `rg` to accurately reflect its `ripgrep` backend, thereby improving clarity and ensuring agents use the correct `ripgrep`-specific arguments.

## Past Attempts

If this user story has been attempted before, the changes made will appear in the git diff. Our policy is to only make a single commit per user story, so you can always review the git diff to review progress across attempts. 


## Requirements

*Specific, measurable acceptance criteria. These define when the story is complete.*

*   The `grep` tool, as exposed to the agent, shall be renamed to `rg`. This means when an agent lists its available tools, "rg" will appear instead of "grep`".
*   All internal code references to the `grep` tool's function name, definition, and usage shall be updated to `rg`.
*   The tool documentation file `.agent/prompts/system/tools.md` shall be updated to reflect the `rg` tool, replacing all instances and descriptions of `grep`.
*   The tool documentation file `.agent/prompts/rules/tool_usage.md` shall be updated to reflect the `rg` tool, replacing all instances and descriptions of `grep`.
*   The function signature in the tool library (`default_api`) for the search tool shall be updated from `grep(pattern, args)` to `rg(pattern, args)`.

## Rules

N/A

## Domain

N/A

## Extra Considerations

*   No specific handling is required for potential agent "memories" that might reference the old `grep` tool name.
*   No external integrations or scripts are expected to be hard-coded to use `grep` that would require updates beyond the agent's internal tool exposure.

## Testing Considerations

*   **Unit Tests:** Verify that the `rg` function can be called successfully with valid `ripgrep` arguments and that `grep` can no longer be called.
*   **Integration Tests:**
    *   Test that the agent correctly identifies and calls the `rg` tool when prompted to perform a search.
    *   Verify that searching with `rg` produces expected results.
*   **Documentation Tests:** Confirm that both `.agent/prompts/system/tools.md` and `.agent/prompts/rules/tool_usage.md` correctly refer to `rg` and no longer contain references to `grep` in the context of the tool.

## Implementation Notes

*   Identify all files that define or reference the `grep` tool, including its function signature in `default_api`, its internal implementation, and its exposure to the agent.
*   Perform a global find-and-replace for `grep` with `rg` in relevant code and documentation files, ensuring to respect context to avoid unintended changes.

## Specification by Example

**Scenario: Agent lists available tools**

*   **Given** an agent is asked to list its available tools.
*   **When** the agent executes the internal command to retrieve tool definitions.
*   **Then** the list of tools returned to the agent will include `rg` with its `ripgrep` description and argument structure, and `grep` will not be present.

**Scenario: Agent uses the search tool**

*   **Given** an agent needs to search for a pattern in files.
*   **When** the agent invokes `rg(pattern="some_pattern", args="--ignore-case")`.
*   **Then** the `ripgrep` command is executed successfully, and results are returned as expected.

## Verification

- [ ] Confirm that `default_api.grep` is removed and `default_api.rg` is callable.
- [ ] Confirm that all internal code references (e.g., tool registration, internal calls) have been updated from `grep` to `rg`.
- [ ] Verify that `.agent/prompts/system/tools.md` correctly describes `rg` and its `ripgrep` arguments.
- [ ] Verify that `.agent/prompts/rules/tool_usage.md` correctly references `rg` in its guidelines.
- [ ] Conduct a test execution where the agent successfully uses the `rg` tool to perform a search.
