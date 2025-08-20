# Global TODO Execution Rules

todowrite and todoread are for managing todo lists inside a session. These tools are not accessible in future sessions. 

Sometimes, major tasks can only be accomplished across multiple sessions. In these cases, it is necessary to use .agent/TODO.md. 

## Structuring TODO.md Items

## Workflow for Executing global TODO Items

When working with global TODO lists (TODO.md files), always follow this execution workflow:

0. Review the current state of the TODO list and select the appropriate task to execute.
1. **Execute the selected task**
2. **Implement this task in THE SIMPLEST WAY POSSIBLE**
3. **Run the quality checks**:
   - Provide examples or actions for the user to take to manually verify completion
4. **Ask for review and WAIT FOR APPROVAL**
5. **Mark the TODO item as complete with [X]**

## Key Principles

- Always implement the simplest solution that meets requirements
- Never proceed to the next task without explicit approval
- All quality checks must pass before requesting review
- Wait for user review before marking items complete
- **Constraint: Only one todo item can have the status 'in_progress' at any given time.**
