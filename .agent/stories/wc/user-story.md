# Implement wc Tool

The agent currently lacks a direct and efficient way to copy files and directories. This user story aims to provide a focused, intuitive `wc` tool for counting the number of lines, words, bytes or characters in a file.

## Past Attempts

N/A - This is a new feature set.

## Requirements

*   **Implement `wc` tool:

## Description and Flags from GNU Documentation

Description
By default, the wc command counts the number of lines, words, and bytes in the files specified by the File parameter. The command writes the number of newline characters, words, and bytes to the standard output and keeps a total count for all named files.

When you use the File parameter, the wc command displays the file names as well as the requested counts. If you do not specify a file name for the File parameter, the wc command uses standard input.

The wc command is affected by the LANG, LC_ALL, LC_CTYPE, and LC_MESSAGES environment variables.

The wc command considers a word to be a string of characters of non-zero length which are delimited by a white space (for example SPACE , TAB).

Flags
Item
Description
-c	Counts bytes unless the -k flag is specified. If the -k flag is specified, the wc command counts characters.
-k	Counts characters. Specifying the -k flag is equivalent to specifying the -klwc flag. If you use the -k flag with other flags, then you must include the -c flag. Otherwise, the -k flag is ignored. For more information, see examples 4 and 5.
“Note: This flag is to be withdrawn in a future release.”
-l	Counts lines.
-m	Counts characters. This flag cannot be used with the -c flag.
-w	Counts words. A word is defined as a string of characters delimited by spaces, tabs, or newline characters.


## Rules

*   The `wc` tool MUST be implemented by shelling out to the corresponding system command.
*   All tools must return clear, concise success messages upon completion.
*   All tools must return clear, concise error messages, closely mimicking standard GNU/Linux command error outputs (e.g., "No such file or directory", "-r not specified; omitting directory 'source'").

## Domain

```
// Filesystem operations
type FileSystemTool interface {
    Execute(args map[string]interface{}) (string, error)
}
```

## Extra Considerations

*   Error messages for shelling out tools should capture the underlying system command's `stderr` as accurately as possible to ensure fidelity.

## Testing Considerations

**YOU CANNOT TEST THESE NEW TOOLS, A NEW BINARY MUST BE BUILT FIRST. WRITE INSTRUCTIONS FOR TESTING TO A FILE CALLED check.tct, OVERWRITING PREVIOUS CONTENT**


*   **Integration Tests:** For the `wc` tool, integration tests are crucial. These should:
    *   Run against a real file system.
    *   Verify the exact behavior for all described requirements.
    *   Create and clean up isolated temporary directories for each test case to prevent test interference.
    *   Specifically verify that `wc` produces the same output/errors as its GNU/Linux counterpart. This might involve capturing output of both the tool and a direct shell call and comparing.
    *   Test copying files to files (overwrite), files to directories, and directories to directories (with and without `recursive=true`).
    *   Verify correct error messages are returned for all failure scenarios.

## Implementation Notes

*   All shelled-out commands should use `os/exec` package in Go, ensuring careful handling of arguments to prevent shell injection (i.e., pass arguments as separate strings, not a single command string).
*   Error wrapping should be used to provide context where native Go errors are returned, maintaining clarity for the agent.


## Verification

- [ ] `wc` tool is implemented and available to the agent.
- [ ] `wc`'s behavior, arguments, and error messages are identical to the standard GNU/Linux `wc` command.