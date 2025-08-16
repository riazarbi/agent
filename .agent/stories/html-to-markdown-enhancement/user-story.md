# HTML-to-Markdown Enhancement

*Enhance web content processing by improving webfetch file extensions and adding HTML-to-Markdown conversion capability to reduce LLM processing overhead from noisy HTML.*

**NOTE TO EXECUTOR**
This user story is very complex, and may need to be done in stages. Focus on working in stages, recording your progress in this user story, which is at .agent/prompts/stories/html-to-markdown-enhancement/user-story.md and in the TODO list, which is at .agent/TODO.md. 

**IT IS HIGHLY LIKELY YOU WON'T COMPLETE THIS IN ONE SHOT, AND THAT IS OK. FOCUS ON QUALITY, PREFER STOPPING OVER THRASHING**


## Requirements

*Specific, measurable acceptance criteria. These define when the story is complete.*

- Enhanced webfetch tool saves downloaded content with appropriate file extensions (.html, .json, .xml, etc.) instead of defaulting to .txt
- New html_to_markdown tool that accepts an HTML file path and converts it to clean Markdown
- html_to_markdown tool saves output with same base filename but .md extension (e.g., page.html → page.md)
- Implementation decision made between reimplementing core conversion logic vs importing the html-to-markdown package
- Both tools integrate seamlessly with existing agent tool architecture
- Conversion produces clean, text-only Markdown suitable for LLM processing (removes images, videos, and other binary data)

## Rules

*Important constraints or business rules that must be followed.*

- New tools must follow existing tool definition patterns (JSON schema, error handling)
- File operations must be safe (no overwriting without explicit intent)
- Must handle common HTML conversion edge cases (malformed HTML, encoding issues)
- Conversion must strip non-text content (images, videos, scripts, styles) and focus on readable text
- Performance should be reasonable for typical web page sizes (under 1MB HTML)

## Domain

*Core domain model in pseudo-code if applicable.*

```
// Enhanced WebFetch
WebFetchInput {
    url: string
}
WebFetchResult {
    path: string           // e.g., ".cache/webfetch/example.com_abc123.html"
    statusCode: int
    contentType: string
}

// New HTML to Markdown Tool
HtmlToMarkdownInput {
    path: string          // input HTML file path
}
HtmlToMarkdownResult {
    inputPath: string     // original HTML file
    outputPath: string    // generated .md file  
    success: boolean
}
```

## Extra Considerations

*Edge cases, non-functional requirements, or gotchas.*

- Handle various HTML content types (text/html, application/xhtml+xml)
- Consider memory usage for large HTML files
- HTML parsing errors should be handled gracefully
- Relative URLs in HTML content (may need base URL context)
- Character encoding detection and handling
- Binary content removal (images, videos, audio, scripts, styles)

## Testing Considerations

*Manual testing scenarios to verify functionality.*

- Test webfetch with various content types to verify correct file extensions
- Test html_to_markdown with sample HTML files from real websites
- Verify HTML-to-Markdown conversion removes non-text content appropriately
- Test error handling with malformed HTML and missing files
- Test with large HTML files to verify reasonable performance
- Test edge cases like empty files or HTML with unusual structures

## Implementation Notes

*Architectural patterns, coding standards, or technology preferences.*

- An existing golang html-to-markdown package exists, the source code is in the html-to-markdown folder. This folder is tempoary and will be deleted after this task - it is simply for evaluation. 
- The source coee should be evaluated to detemrine whether to implenet our own convertor, or to simply use the package as an import.
- Prefer reimplementing core HTML-to-Markdown logic if complexity is manageable (estimate: under 500 lines)
- If reimplementing, focus on essential HTML elements (headings, paragraphs, lists, links, images, code blocks)
- Use existing Go html parser (golang.org/x/net/html) for consistency
- Follow existing tool definition pattern in main.go
- Maintain existing error handling and logging patterns
- Consider extracting shared file operation utilities if code becomes repetitive

## Specification by Example

*Concrete examples: API samples, user flows, or interaction scenarios.*

### Enhanced WebFetch Example
```bash
# Tool call
web_fetch({"url": "https://example.com/page"})

# Returns
{
  "path": ".cache/webfetch/example.com_page_abc123.html",
  "statusCode": 200,
  "contentType": "text/html"
}
```

### HTML-to-Markdown Tool Example
```bash
# Tool call  
html_to_markdown({"path": ".cache/webfetch/example.com_page_abc123.html"})

# Returns
{
  "inputPath": ".cache/webfetch/example.com_page_abc123.html",
  "outputPath": ".cache/webfetch/example.com_page_abc123.md",
  "success": true
}
```

### Typical Workflow
1. Agent calls `web_fetch("https://example.com/article")` → gets `.html` file
2. Agent calls `html_to_markdown(".cache/webfetch/example.com_article_xyz.html")` → gets clean `.md` file
3. Agent calls `read_file(".cache/webfetch/example.com_article_xyz.md")` → gets clean content for LLM processing

## Verification

*Actionable checklist to verify story completion.*

- [ ] webfetch tool saves files with correct extensions based on Content-Type header
- [ ] html_to_markdown tool successfully converts HTML files to readable Markdown
- [ ] html_to_markdown tool generates output files with .md extension and correct base filename
- [ ] Markdown conversion strips non-text content (images, videos, scripts, styles)
- [ ] Both tools integrate with existing agent tool architecture (JSON schema, function definitions)
- [ ] Tools handle error cases gracefully (missing files, invalid HTML)
- [ ] Implementation decision documented (reimplement vs. import) with rationale
- [ ] Manual testing completed with real web pages using example prompts
- [ ] Code follows existing patterns and conventions in main.go

## Next Steps

1. Alter existing web_fetch tool
1. Evaluate html-to-markdown package complexity for reimplement decision
2. Create detailed implementation plan in .agent/TODO.md
3. Execute implementation with progress tracking