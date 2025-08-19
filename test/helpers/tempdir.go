package helpers

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

// TempDir creates a temporary directory for testing and ensures cleanup
func TempDir(t *testing.T) string {
	t.Helper()

	dir, err := os.MkdirTemp("", "agent-test-*")
	require.NoError(t, err)

	t.Cleanup(func() {
		os.RemoveAll(dir)
	})

	return dir
}

// TempFile creates a temporary file with content for testing
func TempFile(t *testing.T, content string) string {
	t.Helper()

	dir := TempDir(t)
	filepath := filepath.Join(dir, "test-file.txt")

	err := os.WriteFile(filepath, []byte(content), 0644)
	require.NoError(t, err)

	return filepath
}

// TempFileWithName creates a temporary file with specific name and content
func TempFileWithName(t *testing.T, filename, content string) string {
	t.Helper()

	dir := TempDir(t)
	filepath := filepath.Join(dir, filename)

	err := os.WriteFile(filepath, []byte(content), 0644)
	require.NoError(t, err)

	return filepath
}
