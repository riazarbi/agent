package tools

import (
	"encoding/json"
)

// WebOperations handles web-related tool operations
type WebOperations struct{}

// NewWebTools returns web operation tools
func NewWebTools() []Tool {
	// TODO: Implement web tools during Phase 2.2
	// This will include: web_fetch, html_to_markdown
	return []Tool{}
}

// Placeholder structs for future implementation
type WebFetchInput struct {
	URL string `json:"url"`
}

type HtmlToMarkdownInput struct {
	Path string `json:"path"`
}

// Placeholder functions for future implementation
func (w *WebOperations) WebFetch(input json.RawMessage) (string, error) {
	// TODO: Move WebFetch implementation from main.go
	return "", nil
}

func (w *WebOperations) HtmlToMarkdown(input json.RawMessage) (string, error) {
	// TODO: Move HtmlToMarkdown implementation from main.go
	return "", nil
}