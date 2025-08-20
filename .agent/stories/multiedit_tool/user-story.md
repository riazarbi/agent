# MultiEdit Tool for Enhanced File Modification

As an agent, I want a tool to perform multiple find-and-replace operations on a single file atomically and sequentially, so that I can efficiently make complex, related changes without risking intermediate broken states.

## Past Attempts

There was one previous attempt, which implemented a working multi_edit tool. YOu can see what has been done by using git_diff. This attempt did not pass QA because of these findings:

*   **Scenario 1 (Successful Multi-Edit on Existing File):** PASSED. The tool successfully performed all specified edits sequentially.
*   **Scenario 2 (New File Creation with Multi-Edit):** FAILED. Although the tool reported success, the second edit (replacing "improved" with "enhanced") was not applied to the newly created file. This indicates a bug in how `multi_edit` handles sequential operations immediately following file creation.
*   **Scenario 3 (Failed Multi-Edit with Rollback - Occurrence Not Found):** PASSED. The tool correctly identified that the second `old_string` did not exist and rolled back the entire operation, leaving the file unchanged.
*   **Scenario 4 (Failed Multi-Edit - Attempt to Create Existing File):** FAILED. The tool *should* have returned an `ATTEMPT_TO_CREATE_EXISTING_FILE` error because I tried to "create" an already existing file using `old_string=""`. Instead, it appended the `new_string` to the file, which is incorrect behavior.
*   **Scenario 5 (Failed Multi-Edit - Old String and New String Identical):** PASSED. The tool correctly failed with an `EDIT_OLD_NEW_IDENTICAL` error, as expected.

In conclusion, the `multi_edit` tool has critical bugs related to its behavior when creating new files and applying subsequent edits, and when attempting to "create" an existing file. While some scenarios worked as expected, the failures in key atomic operations and error handling indicate that the tool's implementation needs attention.

## Requirements

- The tool must accept a `file_path` (absolute path) and an array of `edits`.
- Each edit in the array must specify an `old_string`, a `new_string`, and an optional `replace_all` boolean (defaulting to `false`).
- All edits must be applied in the order they are provided in the `edits` array.
- Each subsequent edit must operate on the result of the previous edit.
- The entire set of edits must be atomic; either all edits succeed and are applied, or if any edit fails, none are applied.
- The `old_string` must exactly match the file's current content, including whitespace and indentation, for an edit to succeed.
- The `old_string` and `new_string` for any single edit cannot be identical.
- The tool must support creating a new file if the `file_path` is new and the first edit has an empty `old_string` and the new file's content as `new_string`.
- The tool must fail if an attempt is made to create an existing file using an empty `old_string`.

## Rules

- This tool should be preferred over the single `Edit` tool when multiple modifications are needed for the same file.
- The tool should not be used for Jupyter notebooks (.ipynb files); the `NotebookEdit` tool should be used instead.
- All edits must result in idiomatic and correct code; the file must not be left in a broken state.
- All file paths provided to the tool must be absolute.
- Emojis should only be added to files if explicitly requested by the user.
- The `replace_all` parameter can be used for renaming variables or other string replacements across the file.

## Domain

```
// MultiEdit operation structure
interface MultiEditOperation {
  file_path: string; // Absolute path to the file
  edits: Array<{
    old_string: string; // Text to replace
    new_string: string; // Text to replace with
    replace_all?: boolean; // Optional: true to replace all occurrences, false by default
  }>;
}
```

## Extra Considerations

- Consider the performance implications for very large files or a large number of edits.
- The tool should provide clear error messages indicating which specific edit failed in a multi-edit operation.
- How will the tool handle potential race conditions if multiple agents or processes try to modify the same file simultaneously? (Though this might be outside the current scope, it's a good future consideration.)
- Ensure robust handling of special characters and encoding in `old_string` and `new_string`.

## Testing Considerations

**YOU CANNOT TEST THESE NEW TOOLS, A NEW BINARY MUST BE BUILT FIRST. WRITE INSTRUCTIONS FOR TESTING TO A FILE CALLED check.txt, OVERWRITING PREVIOUS CONTENT**

- **Unit Tests:**
    - Test individual edit operations (success, no match, old_string == new_string).
    - Test sequential application of edits.
    - Test atomicity (all or nothing).
    - Test `replace_all` functionality.
    - Test new file creation scenario.
    - Test failure to create an existing file.
    - Test edge cases with whitespace, special characters.
- **Integration Tests:**
    - Test `MultiEdit` in combination with `Read` tool (read, then multi-edit, then read again to verify).
    - Test scenarios where earlier edits affect later edits' `old_string` matches (both success and failure cases).
- **End-to-End Tests:**
    - Simulate an agent workflow involving reading a file, deciding on multiple changes, and applying them using `MultiEdit`.

## Implementation Notes

- The `MultiEdit` tool should leverage the existing `Edit` tool's underlying logic for individual replacements.
- Error handling should be precise, clearly indicating which specific edit failed and why.
- The implementation should ensure that file operations are atomic, possibly by writing to a temporary file and then renaming it.

## Specification by Example

**1. Renaming a variable and adding a comment in a Go file:**

```json
{
  "file_path": "/src/main.go",
  "edits": [
    {
      "old_string": "userCount",
      "new_string": "activeUserCount",
      "replace_all": true
    },
    {
      "old_string": "func calculate() int {",
      "new_string": "// calculate calculates the active user count\nfunc calculate() int {",
      "replace_all": false
    }
  ]
}
```

**2. Creating a new file and adding initial content:**

```json
{
  "file_path": "/new_feature/README.md",
  "edits": [
    {
      "old_string": "",
      "new_string": "# New Feature\n\nThis new feature introduces improved user authentication.",
      "replace_all": false
    }
  ]
}
```

**3. Failing scenario: `old_string` does not match:**

```json
{
  "file_path": "/src/config.js",
  "edits": [
    {
      "old_string": "const API_KEY = \"123\";",
      "new_string": "const API_KEY = \"abc\";",
      "replace_all": false
    },
    {
      "old_string": "const DEBUG_MODE = true;",
      "new_string": "const DEBUG_MODE = false;",
      "replace_all": false
    }
  ]
}
```
*(If `const DEBUG_MODE = true;` does not exist in the file, the entire operation should fail and no changes should be applied.)*

## Verification

- [ ] Verify that `MultiEdit` can successfully apply multiple changes to a single file.
- [ ] Verify that changes are applied in the specified order.
- [ ] Verify that if any individual edit fails, the entire operation is rolled back, leaving the file unchanged.
- [ ] Verify that `replace_all` correctly replaces all occurrences when set to `true`.
- [ ] Verify that new files can be created with initial content.
- [ ] Verify that attempts to create existing files with `old_string=""` fail.
- [ ] Verify that the tool handles exact string matching for `old_string` including whitespace and indentation.
- [ ] Verify that an error is returned if `old_string` and `new_string` are identical for any edit.
