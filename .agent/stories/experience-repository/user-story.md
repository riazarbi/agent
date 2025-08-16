# Implement Experience Repository for Cross-Session Knowledge Evolution

*Enable agents to persist, retrieve, and evolve knowledge across sessions through a structured repository of insights, patterns, and solutions, creating a learning memory system.*

## Requirements

- Create `.agent/experience/` directory structure for persistent knowledge storage
- Implement `save_insight` tool for capturing valuable learnings during interactions
- Implement `retrieve_insights` tool for querying accumulated knowledge by topic/pattern
- Implement `update_patterns` tool for evolving recurring solution templates
- All insights must be timestamped, categorized, and human-readable (Markdown format)
- Must support semantic search and retrieval of relevant past experiences
- Must include confidence scoring and validation mechanisms for knowledge quality
- Tools must prevent knowledge pollution through filtering and human oversight options
- Must integrate with existing conversation logging without disruption

## Rules

- Knowledge must be stored in human-readable Markdown files for transparency
- All insights require minimum confidence threshold before persistence
- Contradictory insights must be flagged for human resolution
- Maximum retention limits must prevent unbounded knowledge accumulation
- Must respect existing `.agent/` directory conventions and patterns
- Never automatically modify system prompts - only suggest additions via separate mechanism

## Domain

```
.agent/experience/
├── insights.md              # Timestamped learning entries
├── patterns/
│   ├── problem-X.md        # Recurring problem patterns
│   ├── solution-Y.md       # Proven solution templates
│   └── anti-patterns.md    # What doesn't work
├── context/
│   ├── domain-knowledge.md # Accumulated domain expertise
│   ├── user-preferences.md # Learned user interaction patterns
│   └── system-quirks.md    # Implementation-specific learnings
├── meta/
│   ├── learning-log.md     # What the agent is learning about learning
│   └── knowledge-map.md    # Structure and relationships
└── index.md                # Searchable summary and navigation
```

Core types:
```go
type Insight struct {
    ID          string    `json:"id"`
    Timestamp   time.Time `json:"timestamp"`
    Category    string    `json:"category"`    // problem, solution, pattern, anti-pattern
    Topic       string    `json:"topic"`       // domain area or subject
    Content     string    `json:"content"`     // the actual insight
    Confidence  float64   `json:"confidence"`  // 0.0-1.0 quality score
    Context     string    `json:"context"`     // when/how this was learned
    Tags        []string  `json:"tags"`        // searchable keywords
    Validated   bool      `json:"validated"`   // human confirmation
}

type Pattern struct {
    Name        string   `json:"name"`
    Description string   `json:"description"`
    Instances   []string `json:"instances"`   // references to specific insights
    Effectiveness float64 `json:"effectiveness"`
    LastUsed    time.Time `json:"last_used"`
    UseCount    int      `json:"use_count"`
}
```

## Extra Considerations

- Knowledge retrieval must be fast enough for real-time conversation flow
- Must handle version control conflicts if multiple agents run simultaneously
- Consider implementing knowledge expiry for time-sensitive information
- Need mechanism to identify when insights conflict with each other
- Should provide analytics on knowledge utilization and effectiveness
- Must ensure privacy - no sensitive information inadvertently persisted
- Consider integration with existing TODO.md and story systems for cross-references

## Testing Considerations

- Test insight capture during various types of successful problem-solving
- Test retrieval accuracy when searching for relevant past experiences
- Test pattern evolution when similar problems occur repeatedly
- Test confidence scoring mechanism with high and low-quality insights
- Test knowledge pollution prevention (bad insights don't persist)
- Test cross-session continuity (insights available in new sessions)
- Manual verification that stored knowledge improves future problem-solving

## Implementation Notes

- Follow existing ToolDefinition patterns in main.go
- Use consistent error handling and logging with other tools
- Store everything in Markdown for human readability and version control
- Consider using simple file-based search initially, could evolve to embedded vector search
- Implement proper file locking for concurrent access safety
- Use existing GenerateSchema pattern for tool input validation

## Specification by Example

### Save Insight Example
```json
{
  "category": "solution",
  "topic": "file-parsing",
  "content": "When parsing large Go files, use head/tail/cloc first to understand structure before full read_file. Saves time and provides better context for targeted analysis.",
  "confidence": 0.8,
  "context": "Successfully analyzed main.go efficiently by checking structure first",
  "tags": ["go", "parsing", "efficiency", "file-analysis"]
}
```

### Retrieve Insights Query
```json
{
  "query": "file parsing efficiency",
  "category": "solution",
  "min_confidence": 0.7,
  "limit": 5
}
```

### Pattern Example (auto-generated from insights)
```markdown
# Pattern: Structured File Analysis

**Problem**: Need to understand large code files without reading everything

**Solution Template**:
1. Use `cloc` to get size/complexity overview
2. Use `head` to check imports/structure  
3. Use `tail` to check main functions/exports
4. Use `grep` for specific patterns if needed
5. Only then use `read_file` for targeted sections

**Effectiveness**: 0.85 (based on 12 successful applications)
**Last Used**: 2024-01-15
**Tags**: file-analysis, efficiency, code-review
```

## Verification

- [ ] Experience directory structure created with proper organization
- [ ] save_insight tool implemented with validation and filtering
- [ ] retrieve_insights tool with semantic search capability
- [ ] update_patterns tool for template evolution
- [ ] Confidence scoring system prevents low-quality knowledge persistence
- [ ] Human-readable Markdown format for all stored knowledge
- [ ] Integration with conversation flow doesn't impact performance
- [ ] Knowledge pollution prevention mechanisms working
- [ ] Cross-session knowledge availability verified
- [ ] Analytics and reporting on knowledge utilization
- [ ] Documentation for human oversight and knowledge curation
- [ ] Testing across multiple session types and problem domains

## Next Steps

After creating the story:
1. Save to `.agent/stories/experience-repository/user-story.md`
2. Create initial directory structure in `.agent/experience/`
3. Design insight classification and confidence scoring system
4. Implement MVP with basic save/retrieve functionality
5. Add pattern detection and evolution capabilities
6. Create human oversight and curation interface