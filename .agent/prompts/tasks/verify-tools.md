# Your Task

Please perform a comprehensive test of all your available tools to verify they are working correctly. Follow these steps:

1. **Tool Inventory**: List all your available tools and count them.

2. **File Operations Test**:
   - Use `list_files` to show the current directory contents
   - Use `read_file` to read the first 5 lines of main.go
   - Create a test file using `edit_file` with content "Hello World Test"
   - Use `head` to show the first 3 lines of the test file
   - Use `tail` to show the last 3 lines of the test file
   - Delete the test file using `delete_file`

3. **Todo Management Test**:
   - Use `todowrite` to create a todo list with these items:
     * Task 1: "Test file operations" (status: completed, priority: high)
     * Task 2: "Test todo management" (status: in_progress, priority: medium)  
     * Task 3: "Test search operations" (status: pending, priority: low)
   - Use `todoread` to verify the todos were saved correctly
   - Update the todo list using `todowrite` to mark Task 2 as completed
   - Use `todoread` again to confirm the update

4. **Search Operations Test**:
   - Use `grep` to search for "func main" in the current directory
   - Use `glob` to find all .go files in the project
   - Use `git_diff` to show any current changes

5. **Web Operations Test**:
   - Use `web_fetch` to download content from "https://httpbin.org/json"

6. **Code Analysis Test**:
   - Use `cloc` to count lines of code in the current directory

7. **Summary**: Provide a final report showing:
   - Total number of tools tested
   - Which tools worked successfully
   - Any tools that failed or had issues
   - Confirmation that todowrite, todoread, and list_files are fully functional

Please execute each step methodically and report the results for each tool test. You MUST actually use the tools, don't just pretend to. 

Start now!