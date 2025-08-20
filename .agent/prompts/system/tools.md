# Tool Usage Guidance

This document provides guidance on using the available tools effectively.

## General Principles

*   **Read the Documentation:** Before using a tool, carefully read its description and input parameters.
*   **Start Simple:** Begin with simple use cases and gradually increase complexity.
*   **Check for Errors:** Always check for errors and handle them appropriately.
*   **Use the Right Tool for the Job:** Choose the most direct tool for the task.
*   **Verify Results:** After using a tool, verify that the results are as expected.

## Tool-Specific Guidance

### `read_file`

*   **Purpose:** Reads the contents of a file.
*   **Usage:** `read_file(path="<file_path>")`
*   **Notes:**
    *   The file path is relative to the current directory.
    *   Do not use this tool with directory names.
    *   Returns an error if the file does not exist.
*   **Example:** To read the contents of `README.md`, use `read_file(path="README.md")`

### `list_files`

*   **Purpose:** Lists files and directories.
*   **Usage:** `list_files(path="<optional_path>")`
*   **Notes:**
    *   If no path is provided, lists files in the current directory.
    *   The path is relative to the current directory.
    *   Trailing slashes denote directories.
    *   **Note: This tool automatically excludes certain directories that are related to document versioning, agent session history and agent persona specification.**
*   **Example:** To list files in the current directory, use `list_files()`
*   **Example:** To list files in the `agent` directory, use `list_files(path="agent")`

### `edit_file`

*   **Purpose:** Edits a text file, creates new files, and provides detailed feedback on changes.
*   **Usage:** `edit_file(path="<file_path>", old_str="<text_to_replace>", new_str="<replacement_text>", expected_replacements=<optional_number>)`
*   **Args:**
    *   `path`: The path to the file
    *   `new_str`: Text to replace old_str with **- will be unescaped as a Go string literal.**
    *   `old_str`: Text to search for - must match exactly and must only have one match exactly **- will be unescaped as a Go string literal.**
    *   `expected_replacements`: Optional: The expected number of replacements. If actual replacements differ, an error is returned.
*   **Notes:**
    *   The file path is relative to the current directory.
    *   **Atomic Operation:** Writes are atomic to prevent data corruption (writes to a temporary file then renames).
    *   **New File Creation:** If `old_str` is an empty string (`""`) and the `path` does not exist, a new file will be created with `new_str` content.
    *   **Preventing Overwrite on Creation:** If `old_str` is empty (`""`) but the `path` *already exists*, the tool will return an error (`ATTEMPT_TO_CREATE_EXISTING_FILE`).
    *   **Literal Matching:** `old_str` is treated as a literal string for exact matching.
    *   **Whitespace Preservation:** Accurately preserves whitespace, indentation, and newlines for surrounding content.
    *   **Optional `expected_replacements`:**
        *   If provided, the tool verifies that the actual number of replacements matches this expectation.
        *   If the count does not match, it fails gracefully with an `EDIT_EXPECTED_OCCURRENCE_MISMATCH` error.
    *   **No-op Scenario:** If `old_str` and `new_str` are identical, no changes are applied to the file, and a success message indicating `0 replacements` is returned.
    *   **Detailed Output:** Upon successful completion, the tool returns a JSON object containing:
        *   `message`: A human-readable message indicating success (e.g., "Successfully modified file:...", "Created new file:...", "No changes applied:...").
        *   `actual_replacements`: The number of replacements made.
        *   `diff`: A human-readable, pretty-printed unified diff showing the changes.
    *   **Granular Error Reporting:** Provides specific error messages and codes for various failure scenarios:
        *   `EDIT_INVALID_PATH`: File path is empty.
        *   `EDIT_FILE_READ_ERROR`: Failed to read the file.
        *   `EDIT_FILE_STAT_ERROR`: Failed to get file status.
        *   `EDIT_FILE_NOT_FOUND`: File not found (when `old_str` is not empty).
        *   `EDIT_NO_OCCURRENCE_FOUND`: `old_str` was not found in the file.
        *   `EDIT_EXPECTED_OCCURRENCE_MISMATCH`: Actual replacements do not match `expected_replacements`.
        *   `ATTEMPT_TO_CREATE_EXISTING_FILE`: Attempting to create a file that already exists with `old_str=""`.
        *   `EDIT_DIR_CREATE_ERROR`: Failed to create directory for the file.
        *   `EDIT_TEMP_FILE_CREATE_ERROR`: Failed to create a temporary file.
        *   `EDIT_TEMP_FILE_WRITE_ERROR`: Failed to write to a temporary file.
        *   `EDIT_FILE_RENAME_ERROR`: Failed to rename temporary file to target.
        *   `EDIT_CREATE_TEMP_FILE_ERROR`: Failed to create temp file for new file.
        *   `EDIT_CREATE_TEMP_WRITE_ERROR`: Failed to write to temp new file.
        *   `EDIT_CREATE_FILE_RENAME_ERROR`: Failed to rename temp file to new file.
        *   `EDIT_DIFF_GENERATION_ERROR`: Failed to generate the diff.

*   **Example 1: Successful Single Replacement**
    ```python
    default_api.edit_file(
        path="src/main.go",
        old_str="func oldFunc()",
        new_str="func newFunc()",
        expected_replacements=1
    )
    ```
    (Expected output will include a success message and a unified diff.)

*   **Example 2: Successful File Creation**
    ```python
    default_api.edit_file(
        path="new_feature/README.md",
        old_str="",
        new_str="# New Feature\n\nThis is a new feature."
    )
    ```
    (Expected output will include a success message and a unified diff for creation.)

*   **Example 3: `old_str` Not Found Error**
    ```python
    default_api.edit_file(
        path="src/utils.js",
        old_str="nonExistentFunction()",
        new_str="newFunction()"
    )
    ```
    (Expected error: `Error: EDIT_NO_OCCURRENCE_FOUND: could not find the string to replace: 'nonExistentFunction()'`)

*   **Example 4: `expected_replacements` Mismatch Error**
    ```python
    default_api.edit_file(
        path="config/settings.yaml",
        old_str="debug: false",
        new_str="debug: true",
        expected_replacements=2
    )
    ```
    (Expected error: `Error: EDIT_EXPECTED_OCCURRENCE_MISMATCH: expected 2 occurrences but found 1 for 'debug: false'`)

*   **Example 5: Attempting to Create Existing File Error**
    ```python
    default_api.edit_file(
        path="existing_file.txt",
        old_str="",
        new_str="New content for existing file."
    )
    ```
    (Expected error: `Error: ATTEMPT_TO_CREATE_EXISTING_FILE: File already exists, cannot create using empty old_str: existing_file.txt`)


### `delete_file`

*   **Purpose:** Deletes a file.
*   **Usage:** `delete_file(path="<file_path>")`
*   **Notes:**
    *   The file path is relative to the current directory.
    *   Use with caution, as this operation cannot be undone.
*   **Example:** To delete the file `new_file.txt`, use `delete_file(path="new_file.txt")`

### `rg`

*   **Purpose:** Searches for patterns in files using ripgrep. Supports both literal and regex patterns.
*   **Usage:** `rg(pattern="<search_pattern>", args="<optional_arguments>")`
*   **Notes:**
    *   Supports both literal and regex patterns.
    *   `args` is a space-separated string of ripgrep arguments (e.g., `--ignore-case --hidden`).
*   **Example:** To search for the word "error" in all files, use `rg(pattern="error")`
*   **Example:** To search for the word "error" in all files, ignoring case, use `rg(pattern="error", args="--ignore-case")`

### `glob`

*   **Purpose:** Finds files matching a glob pattern.
*   **Usage:** `glob(pattern="<glob_pattern>")`
*   **Notes:**
    *   Supports standard glob syntax (e.g., `*.go`, `**/*.md`).
*   **Example:** To find all Go files, use `glob(pattern="*.go")`
*   **Example:** To find all Markdown files in all subdirectories, use `glob(pattern="**/*.md")`

### `git_diff`

*   **Purpose:** Shows unstaged changes in the working directory.
*   **Usage:** `git_diff()`
*   **Notes:**
    *   This tool takes no parameters.
    *   It only shows unstaged changes.

### `web_fetch`

*   **Purpose:** Downloads and caches web content.
*   **Usage:** `web_fetch(url="<url>")`
*   **Notes:**
    *   The URL must start with `http://` or `https://`.
    *   Accepts `text/*`, `application/json`, `application/xml`, and `application/xhtml+xml` content types.
    *   **Content is downloaded and cached locally in `.agent/cache/webfetch/` with a generated filename based on the URL and a hash. Returns a JSON object containing the `path` to the cached file, the HTTP `statusCode`, and the `contentType`.**   **Example:** To fetch the content of `https://www.example.com`, use `web_fetch(url="https://www.example.com")`

### `html_to_markdown`

*   **Purpose:** Converts an HTML file to Markdown.
*   **Usage:** `html_to_markdown(path="<file_path>")`
*   **Notes:**
    *   The file path is relative to the current directory.
    *   Saves the output with the same base filename but a `.md` extension.
*   **Example:** To convert `index.html` to Markdown, use `html_to_markdown(path="index.html")`

### `head`

*   **Purpose:** Shows the first N lines of a file.
*   **Usage:** `head(args="<optional_arguments>")`
*   **Notes:**
    *   `args` is a space-separated string of head arguments (e.g., `-n 20 filename`).
    *   Defaults to 10 lines if no arguments are provided.
*   **Example:** To show the first 20 lines of `file.txt`, use `head(args="-n 20 file.txt")`
*   **Example:** To show the first 10 lines of `file.txt`, use `head(args="file.txt")`

### `tail`

*   **Purpose:** Shows the last N lines of a file.
*   **Usage:** `tail(args="<optional_arguments>")`
*   **Notes:**
    *   `args` is a space-separated string of tail arguments (e.g., `-n 20 -f filename`).
    *   Defaults to 10 lines if no arguments are provided.
*   **Example:** To show the last 20 lines of `file.txt`, use `tail(args="-n 20 file.txt")`
*   **Example:** To show the last 10 lines of `file.txt`, use `tail(args="file.txt")`

### `cloc`

*   **Purpose:** Counts lines of code with language breakdown.
*   **Usage:** `cloc(args="<optional_arguments>")`
*   **Notes:**
    *   `args` is a space-separated string of cloc arguments (e.g., `--exclude-dir=.git path`).
    *   If no arguments are provided, it will return usage.
*   **Example:** To count lines of code in the current directory, use `cloc(args=".")`
*   **Example:** To exclude the `.git` directory, use `cloc(args="--exclude-dir=.git .")`

### `todowrite`

*   **Purpose:** Manages the todo list for the current session.
*   **Usage:** `todowrite(todos_json="<todos_json>")`
*   **Notes:**
    *   `todos_json` is a JSON string containing an array of TodoItem objects.
    *   Replaces the entire todo list.
    *   **Constraint: Only one todo item can have the status 'in_progress' at any given time.**
*   **Example:** To create a todo list with one item, use `todowrite(todos_json="[{\"task\": \"Do something\", \"content\": \"Do this task\", \"status\": \"pending\", \"priority\": \"medium\"}]")`

### `todoread`

*   **Purpose:** Retrieves the current todo list.
*   **Usage:** `todoread()`
*   **Notes:**
    *   This tool takes no parameters.
    *   Returns the todo list as a JSON string.
