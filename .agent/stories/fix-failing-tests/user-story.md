# Fix Failing Tests

*This story aims to address the current failure of the `xc test` command, ensuring a stable and reliable test suite. The main objective is to establish a baseline where all relevant tests pass, and to remove any outdated or irrelevant tests.*

## Past Attempts

If this user story has been attempted before, the changes made will appear in the git diff. Our policy is to only make a single commit per user story, so you can always review the git diff to review progress across attempts. 


## Requirements

*Specific, measurable acceptance criteria. These define when the story is complete.*

- The `xc test` command must execute successfully without any errors.
- All relevant existing tests must pass.

## Rules

*Important constraints or business rules that must be followed.*

- **No API Mocking**: Tests must not mock API calls; they should test against real APIs with proper skip logic.
- **Real Environment Testing**: Tests should use actual files, networks, and dependencies.
- **Isolation & Cleanup**: Tests must not affect each other and should properly clean up resources using `t.Cleanup()`.
- **Fast Feedback**: Tests should run quickly and provide clear results.
- **Required Test Coverage**: Aim for at least 80% line coverage, with a target of 90% and 100% for critical paths.
- **Test Naming**: Unit tests should follow the `TestFunctionName_Scenario` naming convention.
- **Test Structure**: Unit tests should be table-driven when testing multiple scenarios.
- **Coverage**: Test both success and error paths.
- **Build Tags**: Ensure appropriate build tags (`unit`, `integration`, `api`, `e2e`, `slow`) are used for different test categories.

## Domain

*Core domain model in pseudo-code if applicable.*

Not applicable

## Extra Considerations

*Edge cases, non-functional requirements, or gotchas.*

- **Removal of Irrelevant/Outdated Tests**: Tests that are identified as being for a very old version of the codebase and are no longer relevant to current functionality may be removed.
- **Identification of "Bad" Tests**: Tests that do not adhere to the testing guide's principles (e.g., flakiness, poor isolation, reliance on mocks) should be refactored or removed.

## Testing Considerations

*What types of tests are needed and what scenarios to cover.*

- **Isolate and Fix Failures**: Identify the first failing test and debug its root cause.
- **Assess Relevance**: Before fixing, determine if the failing test is still relevant to the current codebase.
- **Incremental Fixes**: Address one failing test at a time, re-running `xc test` after each fix.
- **Adherence to Guidelines**: Ensure all new or modified tests strictly follow the `testing_guide.md` principles.
- **Regression Testing**: Verify that fixing one test does not introduce new failures elsewhere.

## Implementation Notes

*Architectural patterns, coding standards, or technology preferences.*

- The `xc test` command is the primary tool for running tests and obtaining test results.
- The `xc test` command can also be used to get test coverage results.

## Specification by Example

*Concrete examples: API samples, user flows, or interaction scenarios.*

**Expected `xc test` output (successful):**

```
go test ./...
ok      github.com/example/project/internal/agent       0.010s
ok      github.com/example/project/internal/config      0.008s
...
```

## Verification

*Actionable checklist to verify story completion.*

- [ ] `xc test` executes successfully without any errors.
- [ ] No new test failures are introduced.
- [ ] Irrelevant or outdated tests have been removed.
- [ ] Remaining tests adhere to the `testing_guide.md` guidelines.
