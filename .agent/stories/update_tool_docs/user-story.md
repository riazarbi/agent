# Update and Refine Agent Tool Documentation

*The existing agent tool documentation files (`.agent/prompts/system/tools.md` and `.agent/prompts/rules/tool_usage.md`) are outdated and lack clear separation in their purpose. This user story aims to update the content of these files, distinctly defining `.agent/prompts/system/tools.md` as the comprehensive tool specification with general guidelines, and `.agent/prompts/rules/tool_usage.md` as the opinionated guide for efficient tool usage. The objective is to enhance clarity and usability for both agents and developers interacting with the tools.*

## Past Attempts

If this user story has been attempted before, the changes made will appear in the git diff. Our policy is to only make a single commit per user story, so you can always review the git diff to review progress across attempts. 


## Requirements

*Specific, measurable acceptance criteria. These define when the story is complete.*

-   **R1: Categorize Tools in `.agent/prompts/system/tools.md`:** The file must be updated to explicitly categorize tools into two distinct sections:
    *   **Simple Passthrough Tools:** These are tools that behave similarly to standard GNU/Linux commands (e.g., `ls`, `cat`, `grep`, `head`, `tail`, `cloc`, `wc`, `cp`, `mv`, `rm`, `touch`, `mkdir`). For these, the documentation should be concise, assuming the user (agent or human) has a basic understanding of their standard command-line behavior.
    *   **Agent-Specific Tools:** These are tools with custom functionalities or significant deviations from standard command-line tools (e.g., `read_file`, `edit_file`, `web_fetch`, `html_to_markdown`, `todowrite`, `todoread`, `git_diff`, `glob`, `rg`). These tools require more detailed, explicit documentation as they are unique to the agent's environment.
-   **R2: Ensure Consistency in `.agent/prompts/rules/tool_usage.md`:** The content of this file must be reviewed to ensure there are no contradictions with the newly defined and categorized tools in `.agent/prompts/system/tools.md`. While the primary goal is not a rewrite, any conflicting or outdated advice must be corrected.

## Rules

*Important constraints or business rules that must be followed.*

-   **Accuracy:** All tool descriptions and usage examples must be factually accurate and reflect the current functionality of the tools.
-   **Clarity & Conciseness:** The language used in both documents should be clear, concise, and easy for both human and agent users to understand. Avoid jargon where possible, or explain it clearly.
-   **No Contradictions:** Ensure that the guidelines and rules presented in `.agent/prompts/rules/tool_usage.md` do not contradict the specifications in `.agent/prompts/system/tools.md`.

## Extra Considerations

*Edge cases, non-functional requirements, or gotchas.*

-   **Audience Nuance:** Remember that the documentation will be consumed by both human developers and AI agents. While the primary goal is clarity for agents, maintain readability and usefulness for humans.
-   **Future Extensibility:** Consider how new tools might be integrated into this documentation structure in the future. The categorization should be robust enough to accommodate additions without requiring major refactors.
-   **Tool Evolution:** The tools themselves may evolve. The documentation should acknowledge that it reflects the current state and may require updates as tools change.

## Testing Considerations

*What types of tests are needed and what scenarios to cover.*

-   **Accuracy Check:** Verify that all tool descriptions, parameters, and examples accurately reflect the tool's current functionality.
-   **Categorization Validation:** Confirm that each tool is placed in the correct category ("Simple Passthrough" or "Agent-Specific") as per the defined criteria.
-   **Clarity Review:** Conduct a review to ensure the language is clear, concise, and unambiguous for both human and agent interpretation.
-   **Consistency Check:** Manually review `.agent/prompts/rules/tool_usage.md` to ensure no contradictions exist with the updated `.agent/prompts/system/tools.md`. Specifically, check if any old advice conflicts with the new tool categorizations or behaviors.
-   **Completeness:** Ensure all available tools are documented in `.agent/prompts/system/tools.md` and that all essential aspects of their usage are covered where appropriate.

## Implementation Notes

*Architectural patterns, coding standards, or technology preferences.*

-   **Markdown Formatting:** Maintain consistent Markdown formatting throughout both documents (e.g., heading levels, code blocks, lists).
-   **Clear Headings:** Use clear and descriptive headings to organize content, especially for tool categorization in `.agent/prompts/system/tools.md`.
-   **Examples:** Ensure code examples provided for "Agent-Specific Tools" are accurate and directly executable where applicable.
-   **Tone and Voice:** Maintain a consistent, informative, and objective tone throughout the documentation.
-   **Internal Cross-referencing:** Where appropriate, use internal links to reference related sections within the documents.

