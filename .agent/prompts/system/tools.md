# Agent Tool Specifications

This document provides a comprehensive specification and general usage guidelines for the tools available to the agent. Tools are categorized to clarify their behavior and intended use.

## General Principles

*   **Read the Documentation:** Before using a tool, carefully read its description and input parameters.
*   **Start Simple:** Begin with simple use cases and gradually increase complexity.
*   **Check for Errors:** Always check for errors and handle them appropriately.
*   **Use the Right Tool for the Job:** Choose the most direct tool for the task.
*   **Verify Results:** After using a tool, verify that the results are as expected.

## Simple Passthrough Tools

These tools behave similarly to standard GNU/Linux commands. Users (agents or humans) are assumed to have a basic understanding of their standard command-line behavior. Documentation for these tools is concise.

### `head`

*   **Purpose:** Show first N lines of a file (default 10 lines).
*   **Usage:** `head(args="<optional_arguments>")`
*   **Example:** To show the first 20 lines of `file.txt`, use `head(args="-n 20 file.txt")`

### `tail`

*   **Purpose:** Show last N lines of a file (default 10 lines).
*   **Usage:** `tail(args="<optional_arguments>")`
*   **Example:** To show the last 20 lines of `file.txt`, use `tail(args="-n 20 file.txt")`

### `cloc`

*   **Purpose:** Count lines of code with language breakdown and statistics.
*   **Usage:** `cloc(args="<optional_arguments>")`
*   **Example:** To count lines of code in the current directory, use `cloc(args=".")`

### `wc`

*   **Purpose:** Print newline, word, and byte counts for each file.
*   **Usage:** `wc(args="<optional_arguments>")`
*   **Example:** To count lines in `file.txt`, use `wc(args="-l file.txt")`

### `cp`

*   **Purpose:** Copy files and directories.
*   **Usage:** `cp(args="<arguments>")`
*   **Example:** To copy `source.txt` to `destination.txt`, use `cp(args="source.txt destination.txt")`

### `mv`

*   **Purpose:** Move or rename files or directories.
*   **Usage:** `mv(args="<arguments>")`
*   **Example:** To rename `old.txt` to `new.txt`, use `mv(args="old.txt new.txt")`

### `rm`

*   **Purpose:** Remove files or directories.
*   **Usage:** `rm(args="<arguments>")`
*   **Example:** To delete `file_to_delete.txt`, use `rm(args="file_to_delete.txt")`

### `touch`

*   **Purpose:** Update access/modification times of files, or create them if they don't exist.
*   **Usage:** `touch(args="<arguments>")`
*   **Example:** To create an empty file `new_empty_file.txt`, use `touch(args="new_empty_file.txt")`

### `mkdir`

*   **Purpose:** Create directories.
*   **Usage:** `mkdir(args="<arguments>")`
*   **Example:** To create a directory `new_dir`, use `mkdir(args="new_dir")`

### `task`

*   **Purpose:** Run taskfile commands.
*   **Usage:** `task(args="<arguments>")`
*   **Example:** To run the default task, use `task(args="")`. To run a specific task like `build`, use `task(args="build")`.

## Agent-Specific Tools

These tools have custom functionalities or significant deviations from standard command-line tools. They require more detailed, explicit documentation as they are unique to the agent's environment.

### `read_file`

*   **Purpose:** Reads the contents of a given relative file path. Use this when you want to see what's inside a file. Do not use this with directory names.
*   **Usage:** `read_file(path="<file_path>")`
*   **Notes:**
    *   The file path is relative to the current directory.
    *   Returns an error if the file does not exist or if the path refers to a directory.
*   **Example:** To read the contents of `README.md`, use `read_file(path="README.md")`

### `list_files`

*   **Purpose:** Lists files and directories at a given path. If no path is provided, lists files in the current directory.
*   **Usage:** `list_files(path="<optional_path>")`
*   **Notes:**
    *   The path is relative to the current directory.
    *   Trailing slashes denote directories.
    *   This tool automatically excludes certain directories related to document versioning, agent session history, and agent persona specification.
*   **Example:** To list files in the current directory, use `list_files()`
*   **Example:** To list files in the `agent` directory, use `list_files(path="agent")`

### `write_file`

*   **Purpose:** Overwrites the entire content of a file or creates a new file. Fails if the path refers to an existing directory.
*   **Usage:** `write_file(content="<content_to_write>", path="<file_path>")`
*   **Notes:**
    *   The file path is relative to the current directory.
    *   This operation will completely replace existing content if the file exists.
*   **Example:** To write "Hello World" to `greeting.txt`, use `write_file(content="Hello World", path="greeting.txt")`

### `append_file`

*   **Purpose:** Appends content to a file, or creates the file if it doesn't exist.
*   **Usage:** `append_file(content="<content_to_append>", path="<file_path>")`
*   **Notes:**
    *   The file path is relative to the current directory.
    *   Content is added to the end of the file.
*   **Example:** To append a new line to `log.txt`, use `append_file(content="New log entry.\n", path="log.txt")`

### `edit_file`

*   **Purpose:** Makes atomic edits to a text file by replacing a single occurrence of `old_str` with `new_str`. Can also create new files. Provides detailed feedback and error handling.
*   **Usage:** `edit_file(path="<file_path>", old_str="<text_to_replace>", new_str="<replacement_text>", expected_replacements=<optional_number>)`
*   **Args:**
    *   `path`: The path to the file.
    *   `new_str`: Text to replace `old_str` with. **Will be unescaped as a Go string literal.**
    *   `old_str`: Text to search for - must match exactly and must only have one match exactly. **Will be unescaped as a Go string literal.**
    *   `expected_replacements`: Optional: The expected number of replacements. If actual replacements differ, an error is returned.
*   **Notes:**
    *   `old_str` and `new_str` MUST be different.
    *   If `old_str` is empty (`""`) and `path` does not exist, a new file will be created with `new_str` content.
    *   If `old_str` is empty (`""`) but `path` *already exists*, an error (`ATTEMPT_TO_CREATE_EXISTING_FILE`) is returned to prevent accidental overwrites.
    *   Supports precise literal matching and preserves whitespace.
    *   Returns JSON output with `message`, `actual_replacements`, and `diff`.
    *   Provides granular error reporting for various failure scenarios.
*   **Example 1: Single Replacement**
    ```python
    default_api.edit_file(
        path="src/main.go",
        old_str="func oldFunc()",
        new_str="func newFunc()",
        expected_replacements=1
    )
    ```
*   **Example 2: File Creation**
    ```python
    default_api.edit_file(
        path="new_feature/README.md",
        old_str="",
        new_str="# New Feature\n\nThis is a new feature."
    )
    ```

### `multi_edit`

*   **Purpose:** Perform multiple find-and-replace operations on a single file atomically and sequentially. All edits are applied in the order they are provided, and each subsequent edit operates on the result of the previous edit. The entire set of edits is atomic; either all edits succeed and are applied, or if any edit fails, none are applied.
*   **Usage:** `multi_edit(file_path="<file_path>", edits=[MultiEditEdits(old_string="<old_str>", new_string="<new_str>", replace_all=<optional_bool>)])`
*   **Args:**
    *   `file_path`: Absolute path to the file.
    *   `edits`: Array of edits to apply sequentially. Each edit object contains:
        *   `old_string`: Text to replace.
        *   `new_string`: Text to replace with.
        *   `replace_all`: Optional: true to replace all occurrences, false by default (replaces only the first).
*   **Notes:**
    *   Edits are applied in the order they appear in the `edits` array.
    *   If the first edit has an empty `old_string` (`""`) and the `file_path` does not exist, a new file will be created with the `new_string` content of that first edit. Subsequent edits will then apply to this newly created content.
    *   If any individual edit fails (e.g., `old_string` not found, `old_string` and `new_string` are identical), the entire `multi_edit` operation is rolled back, and no changes are applied to the file.
    *   Error messages will indicate which specific edit caused the failure.
*   **Example 1: Multiple Replacements**
    ```python
    default_api.multi_edit(
        file_path="config/app.json",
        edits=[
            default_api.MultiEditEdits(old_string='"debug": false', new_string='"debug": true'),
            default_api.MultiEditEdits(old_string='version: "1.0"', new_string='version: "2.0"')
        ]
    )
    ```
*   **Example 2: Create file and apply edits**
    ```python
    default_api.multi_edit(
        file_path="new_project/main.go",
        edits=[
            default_api.MultiEditEdits(old_string='', new_string='package main\n\nfunc main() {\n\t// Initial content\n}\n'),
            default_api.MultiEditEdits(old_string='// Initial content', new_string='// Updated content with a new line')
        ]
    )
    ```

### `delete_file`

*   **Purpose:** Deletes a file.
*   **Usage:** `delete_file(path="<file_path>")`
*   **Notes:**
    *   The file path is relative to the current directory.
    *   Use with caution, as this operation cannot be undone.
*   **Example:** To delete the file `new_file.txt`, use `delete_file(path="new_file.txt")`

### `rg`

*   **Purpose:** Search for patterns in files using ripgrep. Supports both literal and regex patterns.
*   **Usage:** `rg(pattern="<search_pattern>", args="<optional_arguments>")`
*   **Notes:**
    *   `args` is a space-separated string of ripgrep arguments (e.g., `--ignore-case --hidden`).
    *   Useful for quickly finding specific code snippets or text within the codebase.
*   **Example:** To search for the word "error" in all files, ignoring case, use `rg(pattern="error", args="--ignore-case")`

### `glob`

*   **Purpose:** Find files matching a glob pattern. Supports standard glob syntax for file discovery.
*   **Usage:** `glob(pattern="<glob_pattern>")`
*   **Notes:**
    *   Supports patterns like `*.go` (all Go files in current directory) or `**/*.md` (all Markdown files in any subdirectory).
*   **Example:** To find all Go files, use `glob(pattern="*.go")`

### `git_diff`

*   **Purpose:** Returns the output of 'git diff' showing all unstaged changes in the working directory. Use this when you need to see what files have been modified but not yet committed. Do not use this for staged/committed changes.
*   **Usage:** `git_diff()`
*   **Notes:**
    *   This tool takes no parameters.
    *   Provides a quick overview of ongoing modifications.

### `web_fetch`

*   **Purpose:** Download and cache web content locally. Accepts text/*, application/json, application/xml, and application/xhtml+xml content types. Returns path to cached file.
*   **Usage:** `web_fetch(url="<url>")`
*   **Notes:**
    *   The URL must start with `http://` or `https://`.
    *   Content is downloaded and cached locally in `.agent/cache/webfetch/` with a generated filename.
    *   Returns a JSON object containing the `path` to the cached file, the HTTP `statusCode`, and the `contentType`.
*   **Example:** To fetch the content of `https://www.example.com`, use `web_fetch(url="https://www.example.com")`

### `html_to_markdown`

*   **Purpose:** Convert an HTML file to clean Markdown format, removing non-text content like images, videos, scripts, and styles. Saves output with same base filename but .md extension.
*   **Usage:** `html_to_markdown(path="<input_html_file_path>")`
*   **Notes:**
    *   The file path is relative to the current directory.
    *   Useful for extracting readable content from web pages.
*   **Example:** To convert `index.html` to Markdown, use `html_to_markdown(path="index.html")`

### `todowrite`

*   **Purpose:** Create and manage structured task lists for complex multi-step operations within the current session. Replaces the entire todo list with the provided JSON.
*   **Usage:** `todowrite(todos_json="<todos_json_string>")`
*   **Notes:**
    *   Each todo requires: `id` (unique), `content` (description), `status` (pending/in_progress/completed/cancelled), `priority` (high/medium/low).
    *   **Constraint: Only one todo item can have the status 'in_progress' at any given time.**
    *   Data persists only within the current session.
*   **Example:** To create a todo list with one item:
    ```python
    default_api.todowrite(todos_json='''
    [
      {"id": "task-1", "content": "Implement feature X", "status": "in_progress", "priority": "high"}
    ]
    ''')
    ```

### `todoread`

*   **Purpose:** Read the current todo list from session state. Returns structured todos with IDs, content, status, and priority.
*   **Usage:** `todoread()`
*   **Notes:**
    *   This tool takes no parameters.
    *   Provides a quick way to review planned and ongoing tasks.
*   **Example:** To read the current todo list, use `todoread()`

### `list_commands`

*   **Purpose:** Lists all currently available development commands from .agent/Taskfile.yml for building, testing, and workflow automation. Commands are dynamically loaded and can change during sessions.
*   **Usage:** `list_commands()`
*   **Notes:**
    *   This tool takes no parameters.
    *   Returns JSON with commands array containing name and description for each command.
    *   Only lists commands that have descriptions (internal commands are hidden).
    *   **Important:** Commands can be added/removed during sessions, so list regularly to see current options.
*   **Example:** `list_commands()` returns available commands like "build", "test", "lint"

### `run_command`

*   **Purpose:** Executes a development command from .agent/Taskfile.yml (e.g., build, test, lint) with optional arguments for code verification workflows. Commands are dynamic and loaded from current Taskfile.
*   **Usage:** `run_command(command="<command_name>", args="<optional_args>")`
*   **Notes:**
    *   The command parameter specifies which command to run (must exist in Taskfile).
    *   Optional args parameter passes arguments via {{.CLI_ARGS}} variable to the command.
    *   Returns detailed output including stdout, stderr, and exit status.
    *   **Critical for development workflows**: Use after code changes to verify builds and tests pass.
    *   Commands execute with 5-minute timeout for safety.
*   **Example:** `run_command(command="build")` to build the project, or `run_command(command="test", args="-v")` to run tests with verbose output
