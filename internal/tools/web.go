package tools

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	htmltomarkdown "github.com/JohannesKaufmann/html-to-markdown/v2"
)

// WebOperations handles web-related tool operations
type WebOperations struct{}

// NewWebTools returns web operation tools
func NewWebTools() []Tool {
	ops := &WebOperations{}
	return []Tool{
		{
			Name:        "web_fetch",
			Description: "Download and cache web content locally. Accepts text/*, application/json, application/xml, and application/xhtml+xml content types. Returns path to cached file.",
			InputSchema: GenerateSchema[WebFetchInput](),
			Handler:     ops.WebFetch,
		},
		{
			Name:        "html_to_markdown",
			Description: "Convert an HTML file to clean Markdown format, removing non-text content like images, videos, scripts, and styles. Saves output with same base filename but .md extension.",
			InputSchema: GenerateSchema[HtmlToMarkdownInput](),
			Handler:     ops.HtmlToMarkdown,
		},
	}
}

// Input type definitions
type WebFetchInput struct {
	URL string `json:"url" jsonschema_description:"The URL to fetch content from (must start with http:// or https://)"`
}

type HtmlToMarkdownInput struct {
	Path string `json:"path" jsonschema_description:"Input HTML file path"`
}

// Supported content types for WebFetch
var allowedContentTypes = map[string]string{
	"text/plain":            ".txt",
	"text/html":             ".html",
	"text/xml":              ".xml",
	"application/json":      ".json",
	"application/xml":       ".xml",
	"application/xhtml+xml": ".html",
}

// WebFetch downloads and caches web content locally
func (w *WebOperations) WebFetch(input json.RawMessage) (string, error) {
	var webFetchInput WebFetchInput
	err := json.Unmarshal(input, &webFetchInput)
	if err != nil {
		return "", fmt.Errorf("invalid input: %w", err)
	}

	// Validate URL
	if !strings.HasPrefix(webFetchInput.URL, "http://") && !strings.HasPrefix(webFetchInput.URL, "https://") {
		return "", fmt.Errorf("URL must start with http:// or https://")
	}

	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	// Create request
	req, err := http.NewRequest("GET", webFetchInput.URL, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	// Add standard headers
	req.Header.Set("User-Agent", "Mozilla/5.0 WebFetch Tool")
	req.Header.Set("Accept", "text/*, application/json, application/xml, application/xhtml+xml")

	// Make request
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("HTTP GET error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("HTTP request failed with status code: %d", resp.StatusCode)
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %w", err)
	}

	// Generate filename
	baseFilename, err := generateFilename(webFetchInput.URL)
	if err != nil {
		return "", fmt.Errorf("failed to generate filename: %v", err)
	}
	extension, allowed := isAllowedContentType(resp.Header.Get("Content-Type"))
	if !allowed {
		return "", fmt.Errorf("unsupported content type: %s", resp.Header.Get("Content-Type"))
	}
	filename := baseFilename + extension

	// Create cache directory
	cacheDir := ".agent/cache/webfetch"
	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create cache directory: %w", err)
	}

	// Create cache file path
	cacheFilePath := filepath.Join(cacheDir, filename)
	file, err := os.Create(cacheFilePath)
	if err != nil {
		return "", fmt.Errorf("failed to create cache file: %w", err)
	}
	defer file.Close()

	// Write response to cache file
	_, err = file.Write(body)
	if err != nil {
		os.Remove(cacheFilePath) // Clean up on error
		return "", fmt.Errorf("failed to write content: %w", err)
	}

	// Construct result string, including status code
	result := fmt.Sprintf("{\"path\": \"%s\", \"statusCode\": %d, \"contentType\": \"%s\"}", cacheFilePath, resp.StatusCode, resp.Header.Get("Content-Type"))

	return result, nil
}

// HtmlToMarkdown converts an HTML file to Markdown format
func (w *WebOperations) HtmlToMarkdown(input json.RawMessage) (string, error) {
	htmlToMarkdownInput := HtmlToMarkdownInput{}
	if err := json.Unmarshal(input, &htmlToMarkdownInput); err != nil {
		return "", fmt.Errorf("invalid input: %v", err)
	}

	if htmlToMarkdownInput.Path == "" {
		return "", fmt.Errorf("path cannot be empty")
	}

	// Check if input file exists
	if _, err := os.Stat(htmlToMarkdownInput.Path); os.IsNotExist(err) {
		return "", fmt.Errorf("input file does not exist: %s", htmlToMarkdownInput.Path)
	}

	// Read HTML file
	htmlContent, err := os.ReadFile(htmlToMarkdownInput.Path)
	if err != nil {
		return "", fmt.Errorf("failed to read HTML file: %v", err)
	}

	// Convert HTML to Markdown
	markdown, err := htmltomarkdown.ConvertString(string(htmlContent))
	if err != nil {
		return "", fmt.Errorf("failed to convert HTML to markdown: %v", err)
	}

	// Generate output filename: replace .html with .md, or append .md if no extension
	outputPath := htmlToMarkdownInput.Path
	ext := filepath.Ext(outputPath)
	if ext != "" {
		outputPath = strings.TrimSuffix(outputPath, ext) + ".md"
	} else {
		outputPath = outputPath + ".md"
	}

	// Write markdown to output file
	err = os.WriteFile(outputPath, []byte(markdown), 0644)
	if err != nil {
		return "", fmt.Errorf("failed to write markdown file: %v", err)
	}

	// Return result as JSON
	result := map[string]interface{}{
		"inputPath":  htmlToMarkdownInput.Path,
		"outputPath": outputPath,
		"success":    true,
	}

	jsonResult, err := json.Marshal(result)
	if err != nil {
		return "", fmt.Errorf("failed to marshal result: %v", err)
	}

	return string(jsonResult), nil
}

// Helper functions

// generateFilename creates a safe filename from a URL
func generateFilename(inputURL string) (string, error) {
	parsedURL, err := url.Parse(inputURL)
	if err != nil {
		return "", fmt.Errorf("invalid URL: %v", err)
	}

	// Get base components
	host := strings.ToLower(parsedURL.Host)
	path := strings.ToLower(parsedURL.Path)

	// Clean the path
	path = strings.Trim(path, "/")
	path = strings.ReplaceAll(path, "/", "_")

	// Generate hash of full URL
	hasher := sha256.New()
	hasher.Write([]byte(inputURL))
	hash := hex.EncodeToString(hasher.Sum(nil))[:8]

	// Build base filename
	filename := fmt.Sprintf("%s_%s_%s", host, path, hash)

	// Replace invalid characters
	filename = strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '_' || r == '.' {
			return r
		}
		return '_'
	}, filename)

	return filename, nil
}

// isAllowedContentType checks if a content type is allowed and returns appropriate extension
func isAllowedContentType(contentType string) (string, bool) {
	// Extract base content type
	base := strings.Split(contentType, ";")[0]
	base = strings.TrimSpace(base)

	// Check specific allowed types first
	if ext, ok := allowedContentTypes[base]; ok {
		return ext, true
	}

	// Check for other text/* types (default to .txt)
	if strings.HasPrefix(base, "text/") {
		return ".txt", true
	}

	return "", false
}
