# Package development roadmap

## Implement the following tools: 
- [x] edit_file	Modify existing files
- [x] list_files	List directory contents
- [x] read_file	Read file contents
- [ ] touch
- [ ] write_file
- [ ] append_file
- [ ] mv
- [ ] mkdir
- [ ] rm
- [ ] archive
- [ ] expand webfetch
- [ ] multiedit - see https://gist.github.com/wong2/e0f34aac66caf890a332f7b6f9e2ba8f#multiedit
- [x] git_diff Obtain the git diff
- [x] grep	Search file contents
- [x] glob	Find files by pattern
- [x] todowrite	Manage todo lists
- [x] todoread	Read todo lists -  see https://gist.github.com/wong2/e0f34aac66caf890a332f7b6f9e2ba8f#todoread
- [x] webfetch	Fetch web content
- [x] html-to-markdown Convert html content to markdown
- [ ] taskfile executor
- [x] cloc
- [x] head
- [x] tail
- [ ] websearch - see https://gist.github.com/wong2/e0f34aac66caf890a332f7b6f9e2ba8f#websearch

## Features

- [x] Add init command to create a hidden agent directory
- [ ] Implement Taskfile-based tool calling for arbitrary, tightly scoped tool calls.
- [x] Implement preprompt tooling
- [ ] Add testing rules to prompts
- [ ] Add manuals directory for tools and pep8 etc
- [ ] Configurable .agent directory location
- [ ] Allow multiline paste
- [ ] Autocomplete for file paths
- [ ] Track token cost -> maybe generalise to logging everything
- [ ] Anthropic Oauth

## Tech Debt

- [ ] Implement extensive unit testing
- [ ] Implement extensive integration testing
- [ ] Refactor into separate files as per .agent/prompts/rules/go_code_organisation.md

## Nits

- [ ] Implement retry or backoff for too many requests or other API errors. 
- [ ] Parametrise model as well since it differs across backends
- [ ] Clarify use of session todo vs global todo