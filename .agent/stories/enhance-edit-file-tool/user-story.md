# Enhance `edit_file` Tool Reliability and Functionality

This story aims to significantly improve the robustness and usability of the `edit_file` tool. By implementing atomic writes, explicit replacement count validation, and granular error reporting, we can prevent data corruption, provide clearer feedback to the user, and enable more confident and precise code modifications. This enhancement will be informed by research into existing, robust `edit` tool implementations.

## Requirements

*   The `edit_file` tool should perform atomic writes or use a similar mechanism (e.g., write to a temp file then rename) to prevent partial writes and file corruption. Rollback or clear error states must be handled if an operation fails mid-way.
*   The `edit_file` tool should accept an optional `expected_replacements` parameter. If provided, the tool must verify that the actual number of replacements matches this expectation. If not, it should fail gracefully with an informative error.
*   The `edit_file` tool should provide clear and specific error messages for various failure scenarios, including:
    *   `old_str` not found.
    *   Multiple unexpected matches for `old_str` when only one was expected (or `expected_replacements` was not specified).
    *   Mismatch between `expected_replacements` and actual occurrences.
    *   Attempting to create a file that already exists when `old_str` is empty.
    *   File not found when `old_str` is not empty.
    *   Permission denied or other file system errors during read/write.
*   If `old_str` is an empty string (`""`) and the file does not exist, the tool should create a new file with `new_str` content.
*   If `old_str` is an empty string (`""`) and the file *does* exist, the tool should result in a specific error (e.g., "File already exists, cannot create using empty `old_str`").
*   If `old_str` and `new_str` are identical, the tool should report that no changes were applied and return a success status, effectively being a no-op but clearly communicating it.
*   The tool must continue to support exact string replacement, treating `old_str` as a literal string.
*   The tool must accurately preserve whitespace, indentation, and newlines for the surrounding content during replacements.
*   Upon successful completion of an `edit_file` operation (including file creation), the tool should generate and return a clear, human-readable, pretty-printed diff (e.g., in a unified diff format, similar to `git diff` output) showing the changes made between the original content and the new content. This diff will replace the previous terse output of the tool.

## Rules

*   The `edit_file` operation must be atomic to ensure data integrity.
*   The `edit_file` tool should validate the `expected_replacements` parameter against actual replacements if provided.
*   Error messages should be user-friendly and diagnostic, clearly indicating the cause of failure.
*   When `old_str` is empty, the tool's behavior is exclusively for file creation and must not overwrite existing files.
*   `old_str` and `new_str` being identical should result in a no-op with a success confirmation.
*   `old_str` will always be treated as a literal string for exact matching.

## Domain

```
// Key entities and relationships

// edit_file Parameters
interface EditFileArgs {
  path: string;
  old_str: string;
  new_str: string;
  expected_replacements?: number; // Optional
}

// EditFile Result (Success)
interface EditFileSuccessResult {
  message: string; // e.g., "File edited successfully", "No changes applied"
  actual_replacements: number;
  diff: string; // A human-readable, pretty-printed diff showing the changes
}

// EditFile Result (Error)
interface EditFileErrorResult {
  error: string; // Specific error message
  code: string; // A machine-readable error code (e.g., "NOT_FOUND", "MULTIPLE_MATCHES", "FILE_EXISTS")
}
```

## Extra Considerations

*   **Performance for Large Files:** The tool must perform efficiently with large files (e.g., 1500+ lines), aiming to mitigate performance degradation observed in the current tool.
*   **Encoding:** The tool should primarily support UTF-8 encoding.
*   **Line Endings:** The implementation should define a clear strategy for handling line endings (CRLF vs. LF), either preserving them or normalizing them consistently.
*   **Concurrency:** Given the single-threaded nature of the agent, concurrent `edit_file` operations on the same file are not a concern.
*   **Symlinks:** The tool should follow best practices for handling symbolic links.
*   **Permissions:** While system permissions will govern access, the tool should provide informative error messages if a file cannot be read or written due to permission issues.

## Testing Considerations

*   **Positive Scenarios:**
    *   Successful single replacement of `old_str` with `new_str`.
    *   Successful multiple replacements when `expected_replacements` is specified and matches.
    *   Successful file creation when `old_str` is empty and the file does not exist.
    *   No-op scenario: `old_str` and `new_str` are identical, resulting in a success message with no actual file change.
    *   Verification that whitespace, indentation, and newlines are perfectly preserved around the replaced content.
    *   Verification that the returned diff accurately reflects the changes for all successful scenarios (single replacement, multiple replacements, file creation, no-op). The diff should be human-readable and pretty-printed.

*   **Negative Scenarios (Error Handling):**
    *   `old_str` not found in the file.
    *   Multiple unexpected matches for `old_str` when a single replacement was implied (i.e., `expected_replacements` not specified or set to 1).
    *   Mismatch between `expected_replacements` and actual occurrences found (e.g., `expected_replacements=2` but only 1 found, or 3 found).
    *   Attempting to create a file (`old_str=""`) when the file *already exists*.
    *   File not found when `old_str` is not empty (attempting to edit a non-existent file).
    *   Permission denied errors during file read or write operations.
    *   Attempting to edit a directory instead of a file.
    *   Invalid file paths (e.g., not an absolute path if that's a requirement).

*   **Performance Testing:**
    *   Measure the **bytes/characters read and written** by the `edit_file` tool when operating on files of varying sizes, specifically including large files (e.g., 1500 lines, 5000 lines, 10000 lines). The goal is to ensure efficient I/O operations and avoid unnecessary reads/writes.
    *   Compare the I/O efficiency (bytes/characters read/written) against the current `edit_file` tool for large files to demonstrate improvement.

*   **Edge Cases/Specific Conditions:**
    *   Files with unusual characters or non-UTF-8 content (if support for such is considered).
    *   Files with mixed line endings (CRLF/LF) to ensure consistent behavior based on the chosen line ending strategy.
    *   Testing with very long lines or very short lines.

## Implementation Notes

*   **Atomic File Operations:** Prioritize implementing atomic file writes (e.g., write to a temporary file and then rename/move) to ensure data integrity and prevent file corruption during failures.
*   **Error Handling Structure:** Implement a clear and consistent error handling mechanism that allows for granular error types and messages, as outlined in the "Requirements" and "Domain" sections. Consider using custom error classes or structured error objects.
*   **Inspiration from Existing Implementations:** Review and draw inspiration from the error handling, atomic write strategies, and replacement logic found in the provided TypeScript `edit.ts` examples (e.g., `gemini-cli` and `sst/opencode`). Pay particular attention to how they manage different replacement scenarios and error conditions.
*   **Efficient File I/O for Large Files:** When handling large files, consider strategies to optimize file reading and writing to minimize memory usage and improve performance. Avoid reading the entire file into memory unnecessarily if chunked processing is feasible for certain operations.
*   **Literal String Matching:** Ensure that `old_str` is treated as a literal string for exact matching, as specified in the requirements. Avoid any regex interpretation unless explicitly added as a future feature.
*   **Maintain Whitespace and Indentation:** The implementation must preserve the surrounding context's whitespace and indentation during replacements.
*   **Integrated Diff Generation for Readability:** Upon successful completion, the `edit_file` tool should generate a human-readable, pretty-printed diff between the original and new file content and include this diff in its success response. This directly replaces the existing `edit_file` output format and aims for a visual confirmation similar to the `gemini-cli`'s `fileDiff` output.

## Specification by Example

**Scenario 1: Successful Single Replacement**
*   **Action:** User calls `edit_file` to replace a unique string.
*   **Expected `edit_file` parameters:**
    ```
    path: "src/main.go"
    old_str: "func oldFunc()"
    new_str: "func newFunc()"
    ```
*   **Assumed initial `src/main.go` content:**
    ```go
    package main

    func init() {
        // initialization code
    }

    func oldFunc() {
        // old logic
    }

    func anotherFunc() {
        // unrelated code
    }
    ```
*   **Expected `edit_file` tool output (success):**
    ```
    Successfully modified file: src/main.go (1 replacement).
    --- a/src/main.go
    +++ b/src/main.go
    @@ -3,7 +3,7 @@
         // initialization code
     }

    -func oldFunc() {
    +func newFunc() {
         // old logic
     }

    ```

**Scenario 2: Successful File Creation**
*   **Action:** User calls `edit_file` to create a new file.
*   **Expected `edit_file` parameters:**
    ```
    path: "new_feature/README.md"
    old_str: ""
    new_str: "# New Feature\n\nThis is a new feature."
    ```
*   **Expected `edit_file` tool output (success):**
    ```
    Created new file: new_feature/README.md with provided content.
    --- /dev/null
    +++ b/new_feature/README.md
    @@ -0,0 +1,3 @@
    +# New Feature
    +
    +This is a new feature.
    ```

**Scenario 3: `old_str` Not Found Error**
*   **Action:** User calls `edit_file` with `old_str` that doesn't exist in the file.
*   **Expected `edit_file` parameters:**
    ```
    path: "src/utils.js"
    old_str: "nonExistentFunction()"
    new_str: "newFunction()"
    ```
*   **Expected `edit_file` tool output (error):**
    ```
    Error: Failed to edit, could not find the string to replace.
    ```
    (Internal error code: `EDIT_NO_OCCURRENCE_FOUND`)

**Scenario 4: `expected_replacements` Mismatch Error**
*   **Action:** User calls `edit_file` expecting 2 replacements, but only 1 is found.
*   **Expected `edit_file` parameters:**
    ```
    path: "config/settings.yaml"
    old_str: "debug: false"
    new_str: "debug: true"
    expected_replacements: 2
    ```
*   **Assumed initial `config/settings.yaml` content:**
    ```yaml
    app_name: MyWebApp
    environment: development
    debug: false # Only one occurrence
    logging: info
    ```
*   **Expected `edit_file` tool output (error):**
    ```
    Error: Failed to edit, expected 2 occurrences but found 1.
    ```
    (Internal error code: `EDIT_EXPECTED_OCCURRENCE_MISMATCH`)

**Scenario 5: Attempting to Create Existing File Error**
*   **Action:** User calls `edit_file` to create a file (`old_str=""`) but the file already exists.
*   **Expected `edit_file` parameters:**
    ```
    path: "existing_file.txt"
    old_str: ""
    new_str: "New content for existing file."
    ```
*   **Assumed `existing_file.txt` already exists:**
    ```
    This is some old content.
    ```
*   **Expected `edit_file` tool output (error)::**
    ```
    Error: File already exists, cannot create using empty old_str.
    ```
    (Internal error code: `ATTEMPT_TO_CREATE_EXISTING_FILE`)

## Verification

- [ ] All success scenarios outlined in "Specification by Example" are tested and pass.
- [ ] All negative/error scenarios outlined in "Specification by Example" are tested and produce the expected error messages and codes.
- [ ] The `edit_file` tool produces a human-readable, pretty-printed diff in its success response for all successful operations (replacements and file creation).
- [ ] Performance measurements (bytes/characters read and written) for large files meet the efficiency goals and show improvement over the previous tool.
- [ ] Whitespace, indentation, and newlines are correctly preserved during edits.
- [ ] The tool correctly handles file creation when `old_str` is empty and the file does not exist.
- [ ] The tool correctly rejects file creation when `old_str` is empty but the file already exists.
- [ ] The tool correctly performs a no-op when `old_str` and `new_str` are identical.
- [ ] The tool treats `old_str` as a literal string for exact matching.
- [ ] All error messages are clear, specific, and provide diagnostic information.

