# Enhance Tool Descriptions for Clarity and Completeness

Several existing tool descriptions (`list_files`, `todowrite`, `web_fetch`, `edit_file`) currently lack crucial details regarding their specific behaviors, constraints, or return values. This can lead to agent confusion, incorrect assumptions, and unexpected outcomes when using these tools. The main objective is to enhance these tool descriptions to provide comprehensive and accurate information, enabling agents to use them more effectively and predictably.

## Past Attempts

N/A

## Requirements

*   **`list_files` description update**: The tool description for `list_files` in `.agent/prompts/system/tools.md` shall be updated to include the note: "Note: This tool automatically excludes certain directories that are related to document versioning, agent session history and agent persona specification."
*   **`todowrite` description update**: The tool description for `todowrite` in `.agent/prompts/system/tools.md` shall be updated to include the constraint: "Constraint: Only one todo item can have the status 'in_progress' at any given time." Additionally, `.agent/prompts/rules/todo_execution.md` and `.agent/prompts/rules/todo_tool_usage.md` shall be reviewed and updated if necessary to reflect this constraint.
*   **`web_fetch` description update**: The tool description for `web_fetch` in `.agent/prompts/system/tools.md` shall be updated to include the details: "Content is downloaded and cached locally in `.agent/cache/webfetch/` with a generated filename based on the URL and a hash. Returns a JSON object containing the `path` to the cached file, the HTTP `statusCode`, and the `contentType`."
*   **`edit_file` description update**: The descriptions for the `old_str` and `new_str` parameters within the `edit_file` tool in `.agent/prompts/system/tools.md` shall be updated to include a note indicating that the strings "will be unescaped as a Go string literal."
*   General review and update of `.agent/prompts/rules/tool_usage.md` to ensure consistency with the updated tool descriptions.

## Rules

N/A

## Domain

N/A

## Extra Considerations

N/A

## Testing Considerations

*   **Documentation Content Verification**: Confirm that the specified new text is present in the correct documentation files for each tool (`.agent/prompts/system/tools.md`, `.agent/prompts/rules/todo_execution.md`, `.agent/prompts/rules/todo_tool_usage.md`, `.agent/prompts/rules/tool_usage.md`).
*   **Accuracy Check**: Ensure that the updated descriptions accurately reflect the actual behavior and constraints of the tools.
*   **Formatting Compliance**: Verify that the changes maintain the existing formatting and readability of the markdown files.

## Implementation Notes

*   Directly edit the specified markdown files:
    *   `.agent/prompts/system/tools.md`
    *   `.agent/prompts/rules/tool_usage.md`
    *   `.agent/prompts/rules/todo_execution.md`
    *   `.agent/prompts/rules/todo_tool_usage.md`
*   Locate the relevant sections for each tool's description and parameters and insert/modify the text as required by the "Requirements" section.

## Specification by Example

### `list_files` in `.agent/prompts/system/tools.md` (excerpt)

```markdown
### `list_files`

*   **Purpose:** List files and directories at a given path. If no path is provided, lists files in the current directory.
*   **Usage:** `list_files(path="<optional_path>")`
*   **Notes:**
    *   If no path is provided, lists files in the current directory.
    *   The path is relative to the current directory.
    *   Trailing slashes denote directories.
    *   **Note: This tool automatically excludes certain directories that are related to document versioning, agent session history and agent persona specification.**
```

### `todowrite` in `.agent/prompts/system/tools.md` (excerpt)

```markdown
### `todowrite`

*   **Purpose:** Manages the todo list for the current session.
*   **Usage:** `todowrite(todos_json="<todos_json>")`
*   **Notes:**
    *   `todos_json` is a JSON string containing an array of TodoItem objects.
    *   Replaces the entire todo list.
    *   **Constraint: Only one todo item can have the status 'in_progress' at any given time.**
```

### `web_fetch` in `.agent/prompts/system/tools.md` (excerpt)

```markdown
### `web_fetch`

*   **Purpose:** Downloads and caches web content locally. Accepts text/*, application/json, application/xml, and application/xhtml+xml content types. Returns path to cached file.
*   **Usage:** `web_fetch(url="<url>")`
*   **Notes:**
    *   The URL must start with `http://` or `https://`.
    *   Accepts `text/*`, `application/json`, `application/xml`, and `application/xhtml+xml` content types.
    *   **Content is downloaded and cached locally in `.agent/cache/webfetch/` with a generated filename based on the URL and a hash. Returns a JSON object containing the `path` to the cached file, the HTTP `statusCode`, and the `contentType`.**
```

### `edit_file` in `.agent/prompts/system/tools.md` (excerpt - partial)

```markdown
### `edit_file`

...

*   **Args:**
    *   `new_str`: Text to replace old_str with **- will be unescaped as a Go string literal.**
    *   `old_str`: Text to search for - must match exactly and must only have one match exactly **- will be unescaped as a Go string literal.**
    *   `path`: The path to the file
    *   `expected_replacements`: Optional: The expected number of replacements. If actual replacements differ, an error is returned.
```

## Verification

- [ ] Confirm that `.agent/prompts/system/tools.md` has the updated `list_files` description regarding excluded directories.
- [ ] Confirm that `.agent/prompts/system/tools.md` has the updated `todowrite` description regarding the 'in_progress' constraint.
- [ ] Confirm that `.agent/prompts/rules/todo_execution.md` and `.agent/prompts/rules/todo_tool_usage.md` are reviewed and updated if necessary for the 'in_progress' constraint.
- [ ] Confirm that `.agent/prompts/system/tools.md` has the updated `web_fetch` description detailing caching and full return value.
- [ ] Confirm that `.agent/prompts/system/tools.md` has the updated `edit_file` parameter descriptions for `old_str` and `new_str` regarding Go string literal unescaping.
- [ ] Confirm that `.agent/prompts/rules/tool_usage.md` is reviewed and consistent with the updated tool descriptions.
