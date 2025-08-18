# Enhance `edit_file` with Self-Correction

*This story addresses the problem of `edit_file` failing due to subtle inaccuracies (e.g., incorrect escaping, minor formatting differences) in `old_str` or `new_str`. The main objective is to make the `edit_file` tool more robust and reliable by enabling it to automatically correct problematic inputs before attempting file modification, thereby reducing `EDIT_NO_OCCURRENCE_FOUND` and `EDIT_EXPECTED_OCCURRENCE_MISMATCH` errors.*

## Requirements

*   **R1: Programmatic Unescaping:** The `edit_file` tool SHALL first attempt to unescape common over-escaping patterns (e.g., `\\n` to `\n`, `\\"` to `"`) on both `old_str` and `new_str` before the primary replacement attempt.
*   **R5: Preserve Atomic Operations:** All file modifications performed by `edit_file` (including after corrections) SHALL remain atomic, leveraging the existing temporary file and rename mechanism.
*   **R6: Maintain Error Granularity:** `edit_file` SHALL continue to return existing specific error codes (e.g., `EDIT_INVALID_PATH`, `EDIT_FILE_READ_ERROR`, `ATTEMPT_TO_CREATE_EXISTING_FILE`) for scenarios where correction is not applicable or fails to resolve the issue.
*   **R7: Respect `expected_replacements`:** The `expected_replacements` parameter SHALL be validated against the actual occurrences of the *final, corrected* `old_str`. If the corrected `old_str` leads to multiple occurrences when `expected_replacements` is 1, an error should still be returned, unless the correction strategy specifically handles this (e.g. by focusing on a unique match).
*   **R8: No Regression:** All existing, correctly functioning `edit_file` use cases SHALL continue to work as before, without degradation or introduction of new errors.

## Rules

*   **RL1: Correction Precedence:** Programmatic unescaping (R1) MUST be attempted before any further corrections.
*   **RL2: Internal Correction:** Correction mechanisms MUST be internal functions called by `EditFile` and SHALL NOT be exposed as separate callable tools.
*   **RL3: Unique Match for `old_str` Correction:** For `old_str` correction, the system MUST explicitly aim to find an "exact literal, unique match." If a unique match cannot be confidently provided, it should result in the original `EDIT_NO_OCCURRENCE_FOUND` error.## Domain

```go
// Existing EditFileInput
type EditFileInput struct {
	Path               string `json:"path" jsonschema_description:"The path to the file"`
	OldStr             string `json:"old_str" jsonschema_description:"Text to search for - must match exactly and must only have one match exactly"`
	NewStr             string `json:"new_str" jsonschema_description:"Text to replace old_str with"`
	ExpectedReplacements *int   `json:"expected_replacements,omitempty" jsonschema_description:"Optional: The expected number of replacements. If actual replacements differ, an error is returned."`
}

// Internal structure to hold corrected parameters (similar to CorrectedEditParams in TS)
type CorrectedEditResult struct {
    FilePath string
    OldString string
    NewString string
    Occurrences int
}

// Proposed new internal function signatures for LLM correction
// These will encapsulate the LLM prompting and parsing


// Proposed new utility for programmatic unescaping
func unescapeGoString(inputString string) string
```

## Extra Considerations

*   **Performance Impact:** Each correction step introduces additional latency. While beneficial for correctness, monitor overall agent execution time.
*   **Complexity Management:** Integrating multiple correction layers adds complexity to `EditFile`. Ensure the logic is clearly structured and well-commented.
*   **Go Context:** Ensure that calls within the `edit_file` function respect the `context.Context` for timeouts and cancellation.
*   **Error Bubbling:** Clearly define how errors from correction attempts should propagate. If a correction attempt fails, `edit_file` should fall back to reporting the original, uncorrectable error.

## Testing Considerations

*   **Unit Tests for `unescapeGoString`:**
    *   Test cases: `\\n` -> `\n`, `\\t` -> `\t`, `\\\"` -> `\"`, `\\\\` -> `\\`, mixed escaped and unescaped strings, strings with no escaping, empty string.
*   **Integration Tests for `EditFile`:**
    *   **Baseline Success:** Verify existing simple replacements work unchanged.
    *   **Programmatic Unescaping Success (R1):**
        *   `old_str` with `\\n` replaced with `\n` in file.
        *   `new_str` with `\\n` successfully written as `\n`.
    *   **LLM `old_str` Correction Success (R2):**
        *   Test a scenario where `old_str` has minor differences (e.g., extra space, missing escape) that the LLM can correct to a unique match.
    *   **LLM `old_str` and `new_str` Adjustment Success (R3):**
        *   Test `old_str` corrected, and `new_str` adjusted accordingly (e.g., if `old_str` gained a newline, `new_str` also gains one to match intent).
    *   **LLM `new_str` Escaping Correction Success (R4):
        *   Test `edit_file` with `new_str` containing `\\n` which should be written as `\n`.
    *   **Correction Failure - `old_str` not found (0 occurrences):**
        *   Test a problematic `old_str` that cannot be corrected by programmatic unescaping *or* LLM, leading to `EDIT_NO_OCCURRENCE_FOUND`.
    *   **Correction Failure - `expected_replacements` mismatch:**
        *   Test a corrected `old_str` that leads to more than one match when `expected_replacements` is set to 1.
    *   **Edge Cases:**
        *   Multi-line `old_str` and `new_str`.
        *   `old_str` or `new_str` being very long.
        *   File creation (where `old_str` is empty).
        *   Empty files.


## Implementation Notes

*   **New Package/File:** Create a new Go package, e.g., `internal/editcorrector`, to house the `unescapeGoString` function and the correction logic. This keeps `main.go` cleaner.
*   **Error Handling in Correction Chain:** Design the flow within `EditFile` to gracefully handle errors returned by the correction functions. If a correction fails, it should default back to the original `EditFileInput` values and then allow the standard `EditFile` error handling to take over.
*   **Pass Context:** Ensure `context.Context` is passed down to all new functions to enable proper timeout and cancellation.
*   **Logging:** Add detailed internal logging for when corrections are attempted, whether they succeed, and what the corrected values are. This will be invaluable for debugging.

## Specification by Example

**Example 1: Programmatic Unescaping of `old_str`**

**Initial `edit_file` call:**
```python
default_api.edit_file(
    path="my_file.txt",
    old_str="print(\"Hello\\nWorld\")",
    new_str="print(\"Hello New World\")",
    expected_replacements=1
)
```
**`my_file.txt` content:**
```
print("Hello
World")
```
**Expected Behavior:**
1.  `edit_file` attempts to find `"print(\"Hello\\nWorld\")"` in `my_file.txt`. Fails.
2.  `edit_file` calls `unescapeGoString("print(\"Hello\\nWorld\")")` which returns `"print(\"Hello\nWorld\")"`.
3.  `edit_file` re-attempts to find `"print(\"Hello\nWorld\")"`. Succeeds, 1 occurrence found.
4.  `new_str` `"print(\"Hello New World\")"` is used.
5.  File is updated.

**Example 2: Correction for `old_str` and `new_str` Adjustment**

**Initial `edit_file` call:**
```python
default_api.edit_file(
    path="config.js",
    old_str="const API_URL = 'https://old.api/v1';",
    new_str="const API_URL = 'https://new.api/v2';",
    expected_replacements=1
)
```
**`config.js` content:**
```javascript
// Existing content in config.js
const API_URL = "https://old.api/v1"; // Note: double quotes, missing semicolon
```
**Expected Behavior:**
1.  `edit_file` attempts to find `"const API_URL = 'https://old.api/v1';"` in `config.js`. Fails (due to double quotes and missing semicolon in file).
2.  Programmatic unescaping might not fully resolve this.
3.  `edit_file` attempts to identify and correct `"const API_URL = 'https://old.api/v1';"` to an exact match.
4.  System finds and uses `{"corrected_target_snippet": "const API_URL = \"https://old.api/v1\""}`.
5.  `edit_file` re-attempts to find `"const API_URL = \"https://old.api/v1\""`. Succeeds, 1 occurrence found.
6.  `edit_file` adjusts `"const API_URL = 'https://new.api/v2';"` based on the corrected `old_str`.
7.  System uses `{"corrected_new_string": "const API_URL = \"https://new.api/v2\""}`.
8.  File is updated using the corrected `old_str` and `new_str`.

## Verification

- [ ] The `unescapeGoString` function is implemented and unit tested for all specified cases.
- [ ] The `EditFile` function incorporates the programmatic unescaping (R1).
- [ ] All file operations in `EditFile` remain atomic (R5).
- [ ] Existing `EditFile` error handling (e.g., `EDIT_FILE_NOT_FOUND`, `ATTEMPT_TO_CREATE_EXISTING_FILE`) remains functional (R6).
- [ ] The `expected_replacements` parameter is correctly honored with corrected strings (R7).
- [ ] Comprehensive integration tests are added to cover all successful correction scenarios and defined failure modes.
- [ ] No regressions are introduced to existing `edit_file` functionalities (R8).