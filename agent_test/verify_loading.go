package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

// Minimal version of the LoadPromptFile function for verification
func main() {
	if len(os.Args) != 2 {
		log.Fatal("Usage: go run verify_loading.go <path-to-prompt-file>")
	}

	// Change to the test directory so relative paths work
	if err := os.Chdir("agent_test"); err != nil {
		log.Fatal(err)
	}

	filePath := os.Args[1]
	fmt.Printf("=== Verifying prompt file: %s ===\n", filePath)
	
	// This would normally call config.LoadPromptFile, but let's just read and show structure
	content, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}

	ext := filepath.Ext(filePath)
	fmt.Printf("File extension: %s\n", ext)
	fmt.Printf("File size: %d bytes\n", len(content))
	fmt.Printf("Content preview:\n%s\n", string(content))
}