package fixtures

// Common test data used across tests

const (
	// Sample file content for testing
	SampleGoCode = `package main

import "fmt"

func main() {
	fmt.Println("Hello, World!")
}
`

	SampleMarkdown = `# Test Document

This is a test markdown file.

## Section 1

Some content here.

## Section 2

More content here.
`

	SampleJSON = `{
	"name": "test",
	"version": "1.0.0",
	"description": "Test JSON file"
}
`

	SampleHTML = `<!DOCTYPE html>
<html>
<head>
    <title>Test Page</title>
</head>
<body>
    <h1>Test Header</h1>
    <p>Test paragraph with <a href="https://example.com">link</a>.</p>
</body>
</html>
`

	// Expected markdown conversion from SampleHTML
	ExpectedMarkdownFromHTML = `# Test Header

Test paragraph with [link](https://example.com).
`
)

// TodoItems represents sample todo data for testing
var TodoItems = []map[string]interface{}{
	{
		"id":      "todo-1",
		"content": "Write unit tests",
		"status":  "pending",
	},
	{
		"id":      "todo-2",
		"content": "Implement feature X",
		"status":  "in_progress",
	},
	{
		"id":      "todo-3",
		"content": "Fix bug Y",
		"status":  "completed",
	},
}
