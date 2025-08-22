# Unify Command-Line Prompt Handling

*This story aims to unify and enhance the behavior of the command-line flags used to load initial prompts, currently `--preprompts` and `--prompt-file`. The goal is to create a more flexible system that can handle both plain markdown files and structured YAML files for assembling initial agent context, while also improving the command-line interface by renaming the flags.*

## Requirements

- Rename the `--preprompts` flag to `--preprompt`.
- Rename the `--prompt-file` flag to `--prompt`.
- Both `--preprompt` and `--prompt` flags should accept a file path.
- The behavior of the flag should depend on the file extension of the provided file path.

### Markdown File Handling (.md)
- If the file has a `.md` extension, its content should be treated as a single message and added to the agent's conversation history.
- This behavior should be consistent for both `--preprompt` (added before agent activation) and `--prompt` (added after agent activation) flags.

### YAML File Handling (.yml)
- If the file has a `.yml` extension, it should be parsed and processed to generate a series of messages based on a structure that preserves order explicitly.
- The YAML file's top-level structure must be a list of *prompt items*. These items are processed sequentially.
- A prompt item can be one of two types: a `message` constructor or a recursive `file` inclusion.

#### Message Constructor Item
- This item is an object with a single key: `message`.
- The value of `message` must be a list of *content parts*.
- Each content part is an object with a single key, which must be one of: `text`, `command`, or `file`. No other keys are permitted.
- The content parts are concatenated in the order they appear to form a single message.
- Example:
  ```yaml
  - message:
      - text: "The git diff is:"
      - command: "git diff"
  ```

#### Recursive File Inclusion Item
- This item is an object with a single key: `file`.
- The value of `file` must be a path to another YAML file.
- The processor will recursively parse and process this file, inserting the messages it generates at this position in the message sequence.
- The `file` key at this level is exclusively for including other `.yml` files.
- Example:
  ```yaml
  - file: ".agent/prompts/context.yml"
  ```

- A `file` key used *inside* a `message` item refers to a regular text file whose content should be included, not another YAML file to be processed. The processor must differentiate between these two uses of the `file` key based on its location in the structure.

## Rules

- The application must terminate with an informative error message if:
    - A YAML file's structure is invalid.
    - A command specified in a `command` key exits with a non-zero status code.
    - A file path specified in a `file` key does not exist or cannot be read.
- The execution of shell commands via the `command` key is a powerful feature and its security is the responsibility of the operator running the agent. The system will execute the command as provided.

**PROHIBITED** You may not use any git operations or file versioning tools**

## Domain

```
// Represents a single part of a message's content.
type ContentPart struct {
    Text    *string `yaml:"text,omitempty"`
    Command *string `yaml:"command,omitempty"`
    File    *string `yaml:"file,omitempty"` // Path to a content file
}

// Represents a message to be constructed from several parts.
type MessageConstructor struct {
    Message []ContentPart `yaml:"message,omitempty"`
}

// Represents a top-level item in the prompt configuration.
// It can be a message constructor or a recursive file inclusion.
type PromptItem struct {
    MessageConstructor `yaml:",inline"`
    File *string `yaml:"file,omitempty"` // Path to another .yml file
}

// Represents the entire YAML configuration.
type PromptConfig []PromptItem
```

## Extra Considerations

- The order of messages generated from the YAML file must be strictly preserved.
- The file paths for the `file` key should be relative to the agent's working directory.
- Recursive processing of YAML files should handle circular dependencies gracefully (e.g., by detecting them and returning an error).

## Testing Considerations

- Unit tests should be created for the YAML parsing and message compilation logic.
- Test cases should cover:
    - Valid `.md` and `.yml` files.
    - Invalid YAML structures.
    - The specific order of messages generated from the YAML file must be tested to ensure it matches the source file order.
    - The specific order of content concatenated within a single message must be tested to ensure it matches the source file order.
    - Recursive YAML file processing.
    - Error conditions (command failures, missing files).
- Integration tests should verify the end-to-end behavior of the `--preprompt` and `--prompt` flags with both `.md` and `.yml` files.

## Verification

- [ ] The `--preprompts` flag is renamed to `--preprompt`.
- [ ] The `--prompt-file` flag is renamed to `--prompt`.
- [ ] Passing a `.md` file to `--preprompt` or `--prompt` adds its content as a single message.
- [ ] Passing a `.yml` file to `--preprompt` or `--prompt` processes it and adds the compiled messages to the conversation.
- [ ] The YAML processor correctly handles `text`, `command`, and `file` keys.
- [ ] The YAML processor supports recursive file inclusion.
- [ ] The application exits with a clear error message on invalid YAML, command failures, or missing files.
- [ ] All existing tests pass, and new tests for this functionality are added and pass.
- [ ] The agent builds and runs successfully with the new flag behaviors.
