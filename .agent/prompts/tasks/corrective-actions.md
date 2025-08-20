# Corrective actions

I have generated a driving-shot feature in the application using a single prompt. I have reviewed the feature and found some issues. I have marked each issue with a TODO comment.

I would like you to review the source code of this repository and **read all the TODO comments** that I have added. Group these comments into logical tasks that we can tackle as individual units of work. Sometimes these tasks could span multiple TODO comments across different files. At this point, you've only done a grouping **without writing anything to a file**.

Then I want you to generate a new document called `TODO.md` in the prompts directory of this project. This file should contain the list of logical tasks that you compiled in the previous step. Each entry in this file should contain the following detail:

* The headline of the task
* A prefix checkbox to the headline to indicate if the task is done
* A description of what the task entails
* A short prompt that instructs an LLM to address this task
* A list of affected files

Here is an example of what each item in the TODO list should look like:

- [ ] **Task 1: Update Database Schema Column Types**

**Description**

The type of the database columns was inferred as an outdated type. Modern PostgreSQL schemas tend to use the TEXT type instead of VARCHAR. Since this database migration has not been applied yet, we can still safely make this change.

**Prompt**: Update the Flyway migration script `V1__Initial_audit_table.sql` to change all VARCHAR column definitions to TEXT type. The current schema uses VARCHAR with specific length constraints (VARCHAR(50), VARCHAR(100), etc.), but these should be changed to TEXT for better flexibility and consistency with PostgreSQL best practices.

**Files affected**:

- `src/main/resources/db/migration/V1__Initial_audit_table.sql`

---

Every task should follow the format outlined above.

## Extra Considerations

* Each entry of your output should be structured as a **concise but comprehensive prompt** with **adequate context** so that an LLM can use it to alter or generate code to complete the entire task.
* Each prompt should be followed by the complete list of files that are affected.
* Each task should be followed by `---` to separate it from the next task.
* Each task should be scoped correctly so it will be completed as a **single operation**.
* **Do not make up your own TODOs! Only use the ones provided in the source code!**