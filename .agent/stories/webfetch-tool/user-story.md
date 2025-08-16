# Web Content Fetching Tool

*Enable the agent to download and cache web content locally, making it available for analysis through other tools like grep and read_file.*

## Requirements

- Accept URL parameter and fetch content using Go's http.Client
- Store downloaded content in a local cache directory with URL-based filename
- Return path to cached file for further processing
- Return clear error messages for invalid URLs, content types, or fetch failures
- Match existing tool patterns in main.go

## Rules

- Accepted content-types (checked from response header only):
  ```
  text/*
  application/json
  application/xml
  application/xhtml+xml
  ```
- Return error if content-type not in accepted list
- Return error for non-200 status codes
- Default timeout: 30 seconds
- Cache location:
  ```
  .cache/webfetch/<converted-url-filename>
  ```
- Will have to create .cache/webfetch if it does not exist

- URL to filename conversion:
  ```
  Input URL: https://example.com/path/to/page?query=value#fragment
  Output filename: example.com_path_to_page_<hash>.txt
  ```
  
  Conversion rules:
  - Remove protocol (http:// or https://)
  - Replace '/' with '_'
  - Remove query parameters and fragments
  - Add short hash (first 8 chars of SHA-256) of full URL to ensure uniqueness
  - Add appropriate extension based on content-type (.txt, .json, .xml)
  - Convert to lowercase
  - Replace any remaining invalid filename characters with '_'

Example conversions:
```
https://api.github.com/repos/user/repo -> github.com_repos_user_repo_a7f3e1c2.json
https://example.com/article?id=123 -> example.com_article_b8d2f4e6.txt
https://api.service.com/v1/data.json?key=abc -> service.com_v1_data_c9e5d3a1.json
```


## Implementation Steps

1. Parse and validate input URL:
```go
if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
    return "", fmt.Errorf("URL must start with http:// or https://")
}
```

2. Generate cache filename from URL:
```go
filename := generateFilename(url) // Apply URL conversion rules
cachePath := filepath.Join(".cache/webfetch", filename)
```

3. Create HTTP request with timeout:
```go
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()
req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
```

4. Check response status and content-type:
```go
if resp.StatusCode != http.StatusOK {
    return "", fmt.Errorf("HTTP error: %d %s", resp.StatusCode, resp.Status)
}

contentType := resp.Header.Get("Content-Type")
if !isAllowedContentType(contentType) {
    return "", fmt.Errorf("unsupported content type: %s", contentType)

```

5. Download content and save to file:
```go
err = os.MkdirAll(".cache/webfetch", 0755)
if err != nil {
    return "", fmt.Errorf("failed to create cache directory: %v", err)
}

content, err := io.ReadAll(resp.Body)
if err != nil {
    return "", fmt.Errorf("failed to read response body: %v", err)
}

err = os.WriteFile(cachePath, content, 0644)
if err != nil {
    return "", fmt.Errorf("failed to write cache file: %v", err)
}
```

6. Return cache file path and status code:
```go
return CacheResult{
    Path:       cachePath,
    StatusCode: resp.StatusCode,
}, nil
```

## Domain

```go
type CacheResult struct {
    Path       string // Path to cached file
    StatusCode int    // HTTP status code from response
}

func WebFetch(url string) (CacheResult, error)
```

## Extra Considerations

- Large files could cause memory issues since we read entire response into memory
- URL-to-filename conversion must handle all possible URL formats
- Cache directory should be cleared periodically to prevent disk space issues
- Consider adding request headers (User-Agent, Accept) for better compatibility
- May need to handle redirects appropriately

## Testing Considerations

- Test with various URL formats
- Test with all supported content types
- Test error cases:
  - Invalid URLs
  - Network timeouts
  - Non-200 status codes
  - Unsupported content types
  - Invalid characters in URLs
  - Permission issues writing to cache
- Test concurrent access to same URL
- Test very long URLs that might exceed filename length limits

## Verification

- [ ] Tool accepts URL parameter and validates it
- [ ] Tool creates cache directory if not exists
- [ ] Tool correctly converts URLs to filenames
- [ ] Tool properly handles all supported content types
- [ ] Tool returns correct error messages for various failure cases
- [ ] Tool returns both cache path and status code
- [ ] Cache files are readable and contain correct content
- [ ] Tool follows timeout rules
- [ ] Tool properly handles concurrent requests
- [ ] Integration tests pass
- [ ] Unit tests for URL conversion logic pass