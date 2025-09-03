# Test Quality Evaluation Instructions

## Purpose

These instructions guide an evaluator agent to assess the quality and completeness of tests written for the Minimal AI Coding Agent project. The evaluator MUST conduct a comprehensive review of all test implementations against project standards and requirements.

## Required Reading Before Evaluation

**MANDATORY**: The evaluator MUST read these four documents before conducting any assessment:

1. **PRODUCT_SPECIFICATION.md** - Understand what functionality must be tested
2. **HIGH_LEVEL_PLAN.md** - Review testing requirements and current progress
3. **Python Development Standards** - CONTRIBUTING/python-development-standards.md
4. **Python CLI Development Standards** - CONTRIBUTING/python-cli-development-standards.md

**YOUR EVALUATION SHOULD BE APPROPRIATE FOR THE HIGH LEVEL PLAN LEVEL OF COMPLETION**

## Critical Testing Requirements

### Test Coverage Standards
- **MINIMUM 80% line coverage** (enforced - project will fail below this)
- **TARGET 90% line coverage** (preferred standard)
- **100% coverage for public APIs** and core functionality
- **All CLI commands and options** must be tested (100% CLI coverage)

### Test Architecture Requirements
- **No API Mocking Rule**: Test against real APIs with proper skip logic
- **Real Environment Testing**: Use actual files, networks, and dependencies
- **Test Independence**: Tests must run in any order without dependencies
- **Comprehensive Test Types**: Unit, integration, API, and E2E tests

### Project-Specific Constraints
- **Test code does NOT count** toward 1000-line production code limit
- **All tests must pass** before marking step complete
- **Tests must follow** Python development standards exactly
- **CLI-specific testing** requirements must be met

## Evaluation Framework

### 1. Test Structure Assessment (REQUIRED)

The evaluator MUST verify:

```
✓ Test directory structure follows standards:
  tests/
  ├── __init__.py
  ├── unit/              # Unit tests for each module
  ├── integration/       # Integration tests
  ├── api/              # Real API tests (no mocking)
  ├── cli/              # CLI-specific tests
  ├── e2e/              # End-to-end workflow tests
  ├── helpers/          # Test utilities and assertions
  └── fixtures/         # Test data and samples

✓ Test files match production modules:
  - test_config.py for config.py
  - test_agent.py for agent.py
  - test_session.py for session.py
  - test_bash_tool.py for bash_tool.py
  - test_cli.py for cli.py

✓ Test markers properly implemented:
  - @pytest.mark.unit
  - @pytest.mark.integration
  - @pytest.mark.api
  - @pytest.mark.e2e
  - @pytest.mark.slow
```

### 2. Test Quality Standards (REQUIRED)

#### Unit Test Quality Checklist
- [ ] **Descriptive Test Names**: Use long, descriptive names as documentation
- [ ] **Class-Based Organization**: Tests grouped in logical classes
- [ ] **Comprehensive Coverage**: All public methods and functions tested
- [ ] **Edge Cases**: Invalid inputs, boundary conditions, error scenarios
- [ ] **Parametrized Tests**: Multiple input scenarios where appropriate
- [ ] **Property-Based Testing**: Hypothesis tests for complex functions
- [ ] **Proper Assertions**: Specific, meaningful assertions (not just `assert result`)
- [ ] **Test Independence**: Each test can run in isolation
- [ ] **Fixture Usage**: Proper setup/teardown with pytest fixtures

#### Integration Test Quality Checklist  
- [ ] **Real Environment Setup**: Using actual files, directories, processes
- [ ] **Component Interactions**: Testing how modules work together
- [ ] **Temporary Workspaces**: Proper cleanup with temp directories
- [ ] **Configuration Testing**: Real YAML files and environment variables
- [ ] **Session Management**: File-based persistence testing
- [ ] **Error Propagation**: How errors flow between components
- [ ] **Resource Cleanup**: No test artifacts left behind

#### API Test Quality Checklist
- [ ] **No Mocking**: Tests use real API calls only
- [ ] **Skip Logic**: `pytest.skip()` when API keys not provided
- [ ] **Environment Variables**: Proper API key handling
- [ ] **Rate Limiting**: Considerate of API usage limits
- [ ] **Error Handling**: Network failures, invalid responses
- [ ] **Real Data**: Testing with actual API responses
- [ ] **Timeout Handling**: Proper timeout configuration

#### CLI Test Quality Checklist
- [ ] **Subprocess Testing**: Using `subprocess.run()` for CLI calls
- [ ] **All Commands Tested**: Every CLI command and option covered
- [ ] **Help Text Validation**: All help text accuracy verified
- [ ] **Exit Code Testing**: Correct exit codes for success/failure
- [ ] **Input/Output Testing**: File input, various output formats
- [ ] **Error Message Quality**: User-friendly error messages tested
- [ ] **Cross-Platform**: Tests work on different operating systems
- [ ] **Configuration Integration**: CLI config file handling tested

#### E2E Test Quality Checklist
- [ ] **Complete Workflows**: Full user scenarios from start to finish
- [ ] **Real Installation**: Testing with `uv tool install`
- [ ] **Session Persistence**: Create, save, resume workflows
- [ ] **Interactive Mode**: Full conversation workflows
- [ ] **File Processing**: Complete file-based prompt processing
- [ ] **Tool Integration**: Bash tool execution in real scenarios
- [ ] **Configuration Loading**: End-to-end config file usage

### 3. Coverage Analysis (REQUIRED)

The evaluator MUST run and analyze:

```bash
# Coverage analysis commands
uv run pytest --cov=src --cov-report=html --cov-report=term
uv run pytest --cov=src --cov-report=term --cov-fail-under=80
```

#### Coverage Requirements Verification
- [ ] **Overall Coverage ≥ 80%**: Project fails if below this threshold
- [ ] **Module-Level Coverage**: Each module should have high coverage
- [ ] **Critical Path Coverage**: Core functionality at 100%
- [ ] **CLI Coverage**: All CLI commands and options covered
- [ ] **Error Path Coverage**: Exception handling paths tested
- [ ] **Configuration Coverage**: All config options and overrides tested

### 4. Test Execution Assessment (REQUIRED)

The evaluator MUST verify all test categories pass:

```bash
# Test execution verification
uv run pytest -v                    # All tests
uv run pytest -m unit              # Unit tests only
uv run pytest -m integration       # Integration tests
uv run pytest -m api               # API tests (with skip logic)
uv run pytest -m e2e               # End-to-end tests
uv run pytest -m "not slow"        # Fast tests only
```

#### Execution Quality Checklist
- [ ] **All Tests Pass**: Zero test failures
- [ ] **Fast Execution**: Unit tests complete quickly
- [ ] **Skip Logic Works**: API tests skip gracefully without keys
- [ ] **Parallel Execution**: Tests can run concurrently (if needed)
- [ ] **Clean Output**: Clear test results and failure messages
- [ ] **Resource Management**: No hanging processes or files

### 5. Standards Compliance (REQUIRED)

#### Python Development Standards Compliance
- [ ] **Test File Organization**: Follows mandated directory structure
- [ ] **Naming Conventions**: snake_case for test files and functions
- [ ] **Descriptive Names**: Test names describe behavior being tested
- [ ] **Type Hints**: Test functions have proper type annotations
- [ ] **Docstrings**: Complex test classes have descriptive docstrings
- [ ] **Import Organization**: Standard library, third-party, local imports
- [ ] **Line Length**: 88 character maximum
- [ ] **Code Formatting**: Black formatting applied

#### CLI Development Standards Compliance
- [ ] **CLI Test Coverage**: All CLI commands and options tested
- [ ] **Help Text Testing**: Help display accuracy verified
- [ ] **Exit Code Testing**: All exit codes properly validated
- [ ] **Configuration Testing**: CLI config file handling tested
- [ ] **User Experience**: Error messages and outputs tested
- [ ] **Cross-Platform**: CLI tests work on multiple platforms

### 6. Project-Specific Requirements (REQUIRED)

#### Minimal AI Coding Agent Specific Tests
- [ ] **LiteLLM Integration**: Real API calls with proper skip logic
- [ ] **Bash Tool Testing**: Subprocess execution with confirmation prompts
- [ ] **Session Management**: File-based persistence and resume functionality
- [ ] **Configuration System**: YAML loading and environment overrides
- [ ] **CLI Modes**: Interactive, single-shot, file input, session resume
- [ ] **Tool Control**: --allow-tools/--no-tools flag functionality
- [ ] **Confirmation System**: --confirm/--no-confirm flag behavior
- [ ] **Ultra-Minimal Design**: Tests verify only specified features work

## Evaluation Process

### Phase 1: Document Review (MANDATORY)
```
1. Read all four required documents completely
2. Understand what functionality must be tested
3. Identify current testing progress from HIGH_LEVEL_PLAN.md
4. Note any specific testing requirements or constraints
```

### Phase 2: Structure Analysis
```
1. Verify test directory structure matches requirements
2. Check test file organization and naming
3. Validate test markers and categorization
4. Assess test file completeness against production modules
```

### Phase 3: Quality Assessment
```
1. Review test code quality against standards
2. Verify test independence and isolation
3. Check assertion quality and specificity
4. Assess error handling and edge case coverage
```

### Phase 4: Coverage Validation
```
1. Run coverage analysis commands
2. Verify minimum 80% coverage achieved
3. Identify any critical uncovered code paths
4. Check that CLI commands have 100% coverage
```

### Phase 5: Execution Testing
```
1. Run all test categories independently
2. Verify all tests pass without failures
3. Check skip logic for API tests
4. Validate test performance and resource usage
```

### Phase 6: Standards Compliance
```
1. Verify Python development standards compliance
2. Check CLI development standards adherence
3. Validate project-specific requirements met
4. Assess overall test quality and completeness
```

## Critical Issues That Must Be Flagged

### Immediate Failure Conditions
- **Coverage below 80%**: Project fails quality gates
- **Any test failures**: All tests must pass
- **Missing CLI coverage**: CLI commands not tested
- **API mocking detected**: Violates no-mocking rule
- **Standards violations**: Code style or structure issues

### Warning Conditions  
- **Coverage below 90%**: Below target threshold
- **Slow test execution**: Performance concerns
- **Missing edge cases**: Incomplete test scenarios
- **Poor test organization**: Structure improvements needed
- **Inadequate documentation**: Test purpose unclear

### Quality Concerns
- **Generic test names**: Non-descriptive test names
- **Weak assertions**: Tests don't verify specific outcomes
- **Test dependencies**: Tests that require specific run order
- **Resource leaks**: Tests not cleaning up properly
- **Missing integration**: Components not tested together

## Evaluation Output Requirements

The evaluator MUST provide:

### 1. Executive Summary
- Overall test quality assessment (Pass/Fail)
- Coverage percentage achieved
- Number of tests by category
- Critical issues requiring immediate attention

### 2. Detailed Findings
- Structure compliance assessment
- Quality standards adherence
- Coverage analysis results
- Test execution results
- Standards compliance status

### 3. Issue Classification
- **Critical Issues**: Must be fixed before acceptance
- **Warnings**: Should be addressed for quality
- **Recommendations**: Suggestions for improvement

### 4. Specific Recommendations
- Actionable steps to address issues
- Coverage improvement suggestions
- Test quality enhancement recommendations
- Standards compliance corrections needed

### 5. Acceptance Decision
- **ACCEPT**: All requirements met, tests are production-ready
- **CONDITIONAL**: Minor issues that can be addressed quickly
- **REJECT**: Significant issues requiring substantial rework

## Success Criteria for Test Acceptance

### Minimum Acceptance Requirements
- [ ] **≥80% line coverage** achieved
- [ ] **All tests pass** without failures
- [ ] **All CLI commands tested** (100% CLI coverage)
- [ ] **No API mocking** detected
- [ ] **Python standards** compliance verified
- [ ] **CLI standards** compliance verified
- [ ] **Project requirements** met

### Ideal Acceptance Standards
- [ ] **≥90% line coverage** achieved
- [ ] **Comprehensive test types** implemented
- [ ] **Excellent test organization** and naming
- [ ] **Strong edge case coverage**
- [ ] **Performance within limits**
- [ ] **Clear test documentation**

## Quality Assurance Verification

Before completing evaluation, verify:

- [ ] **All test categories evaluated**: Unit, integration, API, CLI, E2E
- [ ] **Coverage analysis completed**: Detailed coverage report reviewed  
- [ ] **Execution verification done**: All test runs completed successfully
- [ ] **Standards compliance checked**: Both Python and CLI standards verified
- [ ] **Critical issues identified**: All blocking issues documented
- [ ] **Recommendations provided**: Actionable improvement suggestions given
- [ ] **Clear acceptance decision**: Explicit accept/conditional/reject decision made

---

## Final Reminder

**The evaluator's role is critical to project quality. A thorough evaluation ensures:**
1. The tests properly verify all functionality works as specified
2. The code meets professional quality standards
3. The CLI provides the expected user experience  
4. The project is ready for production use
5. Future maintenance and development will be sustainable

**Do not compromise on quality standards - the evaluation must be comprehensive and honest.**