# Code Quality Validation Phase

## Purpose
To systematically validate code quality against industry-standard metrics, ensuring the delivered solution meets professional standards for reliability, maintainability, and user value. This phase implements a comprehensive quality assessment framework based on AWS code quality principles.

## Guidelines

**Quality over speed:** Thorough validation prevents technical debt and ensures long-term maintainability. Requirements exist for a reason - partial implementations create future problems.

**Evidence-based assessment:** Provide specific, measurable evidence for all quality judgments. Cite code locations, performance metrics, and test results.

**Context awareness:** Understand the broader system architecture before making quality judgments. Consider the relationship between components and overall system design.

**Professional standards:** Apply consistent quality expectations regardless of implementation source. External providers should meet the same standards as internal teams.

**Read-only evaluation:** This is an evaluation phase only. Do not alter any existing files except for the specified output files. The goal is to assess and report, not to modify the codebase.

**Actually run the code:** try out the siftware and validate claimed functionality.

## Desired Interaction
Act autonomously to complete comprehensive quality validation. Only ask for user input if you lack resources, skills, or tools to complete the assessment.

**Guidelines:**
- Follow-up questions only when additional information is required to complete validation
- Do not ask for clarification unless specifically requested
- **CRITICAL**: Do not modify any existing files except for the specified output files

## Step 2: Context Loading
1. Read the code in this repo to get an idea of what this project does.
2. Check the git diff to see what the developer has done.
3. Look for tests, documentation, and quality indicators throughout the codebase.

## Step 3: Project Purpose Alignment Assessment

**CRITICAL**: Validate that the project still operates within the parameters of the README project purpose.

### Step 3.0: Project Purpose Validation
**Goal:** Ensure the delivered solution aligns with the original project vision and intent.

#### Validation Criteria:
- [ ] **Project Purpose Completeness**: Is the README project purpose section sufficient for evaluation?
- [ ] **Vision Alignment**: Does the solution align with the project's stated vision?
- [ ] **Core Principles Adherence**: Does the implementation follow the project's core principles?
- [ ] **Problem-Solution Fit**: Does the solution address the problems identified in "What We Solve"?
- [ ] **Success Criteria Achievement**: Does the solution meet the "Success Looks Like" criteria?
- [ ] **Scope Compliance**: Does the solution stay within the project's intended scope?

#### Evidence Collection:
- First assess whether the README project purpose section provides sufficient detail for evaluation
- If project purpose is insufficient, document what additional information is needed
- Compare current implementation against README project purpose sections (if sufficient)
- Assess alignment with stated vision and principles
- Evaluate whether delivered value matches intended outcomes
- Document any deviations from original project intent

#### Project Purpose Assessment Guidelines:
**If the project purpose section is insufficient for evaluation:**
- Document what specific information is missing (vision, principles, problems, success criteria)
- Recommend that the project purpose section be fleshed out properly
- Mark project purpose alignment as "INSUFFICIENT_INFO" 
- **CRITICAL**: If project purpose is insufficient, the overall evaluation should be REJECTED
- The failure reason should be "Insufficient project purpose documentation for proper evaluation"
- This is a valid failure that requires the project purpose to be properly defined before re-evaluation

## Step 4: Comprehensive Quality Assessment

Document your findings while working through the following quality dimensions:

### Step 3.1: Reliability Assessment
**Goal:** Code runs as documented each and every time, handling unexpected inputs without crashing.

#### Validation Criteria:
- [ ] **Consistent Execution**: Application starts and runs reliably across multiple executions
- [ ] **Error Handling**: Graceful handling of unexpected inputs and edge cases
- [ ] **Crash Prevention**: No unhandled exceptions or system crashes during normal operation
- [ ] **Resource Management**: Proper cleanup of resources (files, connections, memory)
- [ ] **State Management**: Consistent state handling across application lifecycle

#### Evidence Collection:
- Execute application multiple times with various inputs
- Test error conditions and edge cases
- Monitor resource usage and cleanup
- Document any crashes, exceptions, or unexpected behaviors

### Step 3.2: Extendibility Assessment
**Goal:** Code is easy to update, modify, or use for new functionality.

#### Validation Criteria:
- [ ] **Modular Architecture**: Clear separation of concerns and modular design
- [ ] **Coding Standards**: Compliance with established coding standards and style guides
- [ ] **Code Complexity**: Reasonable cyclomatic complexity and function sizes
- [ ] **Documentation**: Clear documentation of interfaces and APIs
- [ ] **Configuration**: Externalized configuration and environment-specific settings

#### Evidence Collection:
- Analyze code structure and modularity
- Review adherence to coding standards
- Assess code complexity metrics
- Evaluate documentation completeness
- Check configuration management

### Step 3.3: Testability Assessment
**Goal:** Code is easy to develop tests for and run tests on.

#### Validation Criteria:
- [ ] **Test Coverage**: Comprehensive test suite with acceptable coverage percentage
- [ ] **Unit Test Quality**: Meaningful unit tests that validate core functionality
- [ ] **Integration Tests**: Tests that validate component interactions
- [ ] **Test Documentation**: Clear test instructions and expected results
- [ ] **Test Execution**: Tests run successfully and provide clear pass/fail results

#### Evidence Collection:
- Execute test suite and record coverage metrics
- Review test quality and meaningfulness
- Verify test documentation accuracy
- Document test execution results and any failures

### Step 3.4: Portability Assessment
**Goal:** Code can be moved between environments with minimal work.

#### Validation Criteria:
- [ ] **Environment Independence**: Code runs in different environments without modification
- [ ] **Dependency Management**: Clear dependency specifications and version management
- [ ] **Configuration Flexibility**: Environment-specific configuration without code changes
- [ ] **Platform Compatibility**: Works across target platforms and operating systems
- [ ] **Deployment Readiness**: Clear deployment instructions and containerization if applicable

#### Evidence Collection:
- Test in different environments (if possible)
- Review dependency specifications
- Check configuration management
- Validate deployment documentation
- Assess platform compatibility

### Step 3.5: Reusability Assessment
**Goal:** Code is modular and designed for reuse.

#### Validation Criteria:
- [ ] **Component Design**: Functions and classes designed for reuse
- [ ] **API Design**: Clean, well-documented APIs for external consumption
- [ ] **Dependency Injection**: Loose coupling through dependency injection where appropriate
- [ ] **Library Usage**: Proper use of external libraries and frameworks
- [ ] **Code Duplication**: Minimal code duplication and shared utilities

#### Evidence Collection:
- Analyze component design and coupling
- Review API documentation and design
- Check for code duplication and shared utilities
- Assess library usage and dependency management

### Step 3.6: Maintainability Assessment
**Goal:** Code is easy to understand, modify, and maintain over time.

#### Validation Criteria:
- [ ] **Code Readability**: Clear, readable code with meaningful names and structure
- [ ] **Documentation Quality**: Comprehensive documentation of code, APIs, and processes
- [ ] **Code Organization**: Logical file and directory structure
- [ ] **Version Control**: Proper use of version control and meaningful commit messages
- [ ] **Technical Debt**: Minimal technical debt and clear path for future improvements

#### Evidence Collection:
- Review code readability and naming conventions
- Assess documentation completeness and accuracy
- Evaluate code organization and structure
- Check version control practices
- Identify technical debt and improvement opportunities

## Step 5: Requirements Traceability Validation

**CRITICAL**: Validate the story→specification→product chain addresses the user problem.

#### Core Requirements Validation:
- [ ] **User Story Coverage**: Does specification address ALL user story acceptance criteria?
- [ ] **Specification Implementation**: Does product implement ALL specification requirements?
- [ ] **User Value Delivery**: Can end users achieve the user story goals with this product?
- [ ] **Problem-Solution Fit**: Does the final product solve the original user problem?

#### Gap Analysis:
- [ ] **Missing Elements**: What user story requirements are not addressed?
- [ ] **Implementation Gaps**: What specification requirements are not implemented?
- [ ] **Value Gaps**: Does the product deliver the intended user benefit?

## Step 6: Application Validation

### Step 6.1: Basic Functionality Validation
**CRITICAL**: Every deliverable MUST pass basic functionality validation:

- [ ] **Application Startup**: Application starts without errors using documented commands
- [ ] **Basic Connectivity**: Health endpoints respond correctly
- [ ] **Core Functionality**: Primary user workflows function as intended
- [ ] **Graceful Shutdown**: Application stops cleanly without resource leaks

### Step 6.2: Performance Validation
- [ ] **Response Times**: Meets specified performance targets (sub-second for web applications)
- [ ] **Resource Usage**: Reasonable memory and CPU usage
- [ ] **Scalability**: Handles expected load without degradation

### Step 6.3: Security Validation
- [ ] **Input Validation**: Proper validation of user inputs
- [ ] **Error Handling**: Secure error handling without information disclosure
- [ ] **Authentication/Authorization**: Proper security controls where applicable

## Quality Validation Criteria

### Core Quality Checklist
- [ ] **Project Purpose Alignment**: Solution aligns with README project purpose
- [ ] **Reliability**: Code runs consistently and handles errors gracefully
- [ ] **Extendibility**: Modular design with clear interfaces and documentation
- [ ] **Testability**: Comprehensive test suite with good coverage
- [ ] **Portability**: Works across environments with minimal configuration
- [ ] **Reusability**: Components designed for reuse with clean APIs
- [ ] **Maintainability**: Readable, well-documented, and organized code
- [ ] **Requirements Compliance**: All user story and specification requirements met
- [ ] **Application Functionality**: Core application works as intended
- [ ] **Performance Targets**: Meets specified performance requirements
- [ ] **Security Standards**: Proper security controls implemented

### Quality Standards Documentation
Throughout validation, maintain:
1. **Quality Metrics**: Track specific quality indicators and measurements
2. **Evidence-Based Assessment**: Cite specific code locations (file:line) for all claims
3. **Objective Reporting**: Focus on measurable quality criteria
4. **Professional Standards**: Apply consistent quality expectations

### ⚠️ Evaluation Priority
**Proper evaluation requires sufficient project purpose documentation.** A project without a properly defined purpose cannot be meaningfully evaluated for alignment. If the README project purpose section is insufficient, the evaluation should be REJECTED with the reason being insufficient documentation for proper assessment.

## Step 7: Final Evaluation Report Generation

Generate comprehensive evaluation report at the specified output location with both markdown and JSON outputs. Focus on evaluation and recommendations rather than intervention:

### Report Structure Template
```markdown
# Project Evaluation Report: [Project Name]

**Date:** [Date]  
**Evaluator:** Senior Software Engineer (Quality Validation Phase)  
**Status:** ACCEPT/REJECT

## Executive Summary
[Overall evaluation status and key findings]

## Project Purpose Alignment Assessment
- **Project Purpose Completeness**: [Sufficient/Insufficient] - [Evidence]
- **Vision Alignment**: [Excellent/Good/Adequate/Poor/Insufficient_Info] - [Evidence]
- **Core Principles Adherence**: [Excellent/Good/Adequate/Poor/Insufficient_Info] - [Evidence]
- **Problem-Solution Fit**: [Excellent/Good/Adequate/Poor/Insufficient_Info] - [Evidence]
- **Success Criteria Achievement**: [Excellent/Good/Adequate/Poor/Insufficient_Info] - [Evidence]
- **Overall Purpose Alignment**: [Excellent/Good/Adequate/Poor/Insufficient_Info]

## Quality Dimensions Assessment

### Reliability Assessment
- **Consistent Execution**: [Excellent/Good/Adequate/Poor] - [Evidence]
- **Error Handling**: [Excellent/Good/Adequate/Poor] - [Evidence]
- **Crash Prevention**: [Excellent/Good/Adequate/Poor] - [Evidence]
- **Resource Management**: [Excellent/Good/Adequate/Poor] - [Evidence]
- **Overall Reliability**: [Excellent/Good/Adequate/Poor]

### Extendibility Assessment
- **Modular Architecture**: [Excellent/Good/Adequate/Poor] - [Evidence]
- **Coding Standards**: [Excellent/Good/Adequate/Poor] - [Evidence]
- **Code Complexity**: [Excellent/Good/Adequate/Poor] - [Evidence]
- **Documentation**: [Excellent/Good/Adequate/Poor] - [Evidence]
- **Overall Extendibility**: [Excellent/Good/Adequate/Poor]

### Testability Assessment
- **Test Coverage**: [Excellent/Good/Adequate/Poor] - [Percentage and evidence]
- **Unit Test Quality**: [Excellent/Good/Adequate/Poor] - [Evidence]
- **Integration Tests**: [Excellent/Good/Adequate/Poor] - [Evidence]
- **Test Execution**: [Excellent/Good/Adequate/Poor] - [Results]
- **Overall Testability**: [Excellent/Good/Adequate/Poor]

### Portability Assessment
- **Environment Independence**: [Excellent/Good/Adequate/Poor] - [Evidence]
- **Dependency Management**: [Excellent/Good/Adequate/Poor] - [Evidence]
- **Configuration Flexibility**: [Excellent/Good/Adequate/Poor] - [Evidence]
- **Deployment Readiness**: [Excellent/Good/Adequate/Poor] - [Evidence]
- **Overall Portability**: [Excellent/Good/Adequate/Poor]

### Reusability Assessment
- **Component Design**: [Excellent/Good/Adequate/Poor] - [Evidence]
- **API Design**: [Excellent/Good/Adequate/Poor] - [Evidence]
- **Dependency Injection**: [Excellent/Good/Adequate/Poor] - [Evidence]
- **Code Duplication**: [Excellent/Good/Adequate/Poor] - [Evidence]
- **Overall Reusability**: [Excellent/Good/Adequate/Poor]

### Maintainability Assessment
- **Code Readability**: [Excellent/Good/Adequate/Poor] - [Evidence]
- **Documentation Quality**: [Excellent/Good/Adequate/Poor] - [Evidence]
- **Code Organization**: [Excellent/Good/Adequate/Poor] - [Evidence]
- **Technical Debt**: [Excellent/Good/Adequate/Poor] - [Assessment]
- **Overall Maintainability**: [Excellent/Good/Adequate/Poor]

## Requirements Traceability Assessment
- **User Story → Specification**: [Does spec address user problem?]
- **Specification → Product**: [Does product implement spec?]  
- **Product → User Story**: [Does product deliver user value?]
- **Gap Analysis**: [Any missing elements or implementation gaps?]

## Application Validation
- **Startup Test**: [Did application start successfully?]
- **Functionality Test**: [Do core features work as intended?]
- **Performance Test**: [Meets performance targets?]
- **Security Test**: [Proper security controls?]

## Quality Metrics Summary
- **Overall Quality Assessment**: [Excellent/Good/Adequate/Poor]
- **Critical Issues**: [None/Minor/Major/Critical]
- **Technical Debt**: [None/Minor/Major/Significant]
- **Maintenance Readiness**: [Ready/Needs Improvement/Not Ready]

## Recommendations for Next Story Cycle
[Structured recommendations that can be fed into the next user story for improvement]

### Immediate Improvements
[Critical issues that should be addressed in the next story]

### Technical Debt Items
[Technical debt that should be prioritized in future stories]

### Quality Enhancement Opportunities
[Areas where quality can be improved in subsequent stories]

## Final Decision
[ACCEPT/REJECT with key rationale based on evaluation standards]
```

### Evidence Requirements
- **Code Citations**: Include specific file:line references for all quality claims
- **Performance Data**: Qualitative evidence for performance assertions (e.g., "fast response times", "acceptable memory usage")
- **Test Results**: Qualitative test assessment with coverage percentages where available
- **Quality Metrics**: Qualitative assessments for each quality dimension with supporting evidence

## Output

**IMPORTANT**: Only these output files should be created or modified. Do not alter any existing project files.

### Markdown Report
1. File: `check.md`

### JSON Quality Scores
2. File: `@docs/stories/[user-story]/quality_scores.json`

**IMPORTANT**: The JSON must be valid JSON with no trailing commas, proper escaping, and exact string values as specified in the schema. All quality assessments must be one of the exact string values shown in the schema. The recommendations section should provide actionable items that can be directly fed into the next user story planning phase.
