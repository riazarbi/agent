# Package development roadmap

## Implement the following tools: 
- [x] edit_file	Modify existing files
- [x] delete_file	Modify existing files
- [x] list_files	List directory contents
- [x] read_file	Read file contents
- [x] git_diff Obtain the git diff
- [x] grep	Search file contents
- [x] glob	Find files by pattern
- [ ] todowrite	Manage todo lists
- [ ] todoread	Read todo lists
- [x] webfetch	Fetch web content
- [x] html-to-markdown Convert html content to markdown
- [ ] taskfile executor

## Features

- [x] Add init command to create a hidden agent directory
- [ ] Implement Taskfile-based tool calling for arbitrary, tightly scoped tool calls.
- [x] Implement preprompt tooling
- [ ] Add testing rules to prompts

## Tech Debt

- [ ] Implement extensive unit testing
- [ ] Implement extensive integration testing
- [ ] Refactor into separate files as per .agent/prompts/rules/go_code_oragnisation.md