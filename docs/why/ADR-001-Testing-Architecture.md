# ADR-001: Testing Architecture Standards

## Status
**ACCEPTED** - 2025-08-19

## Context
The code-editing-agent project requires a robust, maintainable testing architecture that ensures reliability while remaining practical for contributors. We need to establish clear standards that can be validated automatically and followed consistently.

## Decision

### Testing Philosophy
We adopt a **"Real Environment First"** testing philosophy with the following core principles:

1. **No API Mocking**: Test against real APIs with skip logic when credentials unavailable
2. **Real Dependencies**: Use actual file systems, networks, and external tools
3. **Comprehensive Coverage**: Unit, integration, and end-to-end testing
4. **Automated Validation**: Architecture standards enforced through tooling

### Architecture Standards

#### Directory Structure
```
├── internal/package/source_test.go  # Unit tests alongside source
├── test/helpers/                    # Shared test utilities
├── test/fixtures/                   # Test data and constants
├── test/integration/                # Integration tests
├── test/e2e/                       # End-to-end tests
└── test/api/                       # API-specific tests
```

#### Test Categories with Build Tags
- `unit` (default): Fast, isolated function tests
- `integration`: Component interaction tests
- `api`: Real API tests with skip logic
- `e2e`: Complete user workflow tests

#### Required Standards
- **80% minimum test coverage**, 90% target
- **Table-driven tests** for multiple scenarios
- **Proper isolation** with cleanup using `t.Cleanup()`
- **Real API testing** with environment-based skipping
- **No mocking of external APIs** ever

## Consequences

### Positive
- **Higher confidence**: Tests validate real behavior, not mocked approximations
- **Better error detection**: Real environment testing catches integration issues
- **Maintainable tests**: No complex mock setup and maintenance
- **Consistent standards**: Automated validation ensures compliance
- **Realistic testing**: Tests reflect actual usage scenarios

### Negative
- **Test complexity**: Real environment tests require more setup
- **External dependencies**: Tests may need API keys, network access
- **Slower feedback**: Some tests inherently slower than unit tests
- **Environment sensitivity**: Tests may fail due to external factors

### Mitigation Strategies
- **Proper categorization**: Use build tags to run fast tests by default
- **Graceful degradation**: Skip tests when dependencies unavailable
- **Clear documentation**: Comprehensive guides for contributors
- **Automated validation**: Tooling to enforce standards

## Implementation

### Phase 1: Documentation & Tooling ✅
- [x] Comprehensive testing guide (`TESTING_GUIDE.md`)
- [x] Automated validation script (`scripts/validate-tests.sh`)
- [x] Pre-commit hooks for standard enforcement
- [x] Architecture decision record (this document)

### Phase 2: Test Infrastructure (In Progress)
- [ ] Enhanced test helpers with API skip logic
- [ ] Build tags for test categorization
- [ ] Test data fixtures and samples
- [ ] CI/CD integration with validation

### Phase 3: Test Implementation
- [ ] Migration of existing tests to new standards
- [ ] Integration tests for all tool categories
- [ ] API tests with real environment logic
- [ ] End-to-end workflow tests

## Validation

### Automated Checks
- Test architecture validation script
- Pre-commit hooks preventing API mocking
- CI/CD pipeline enforcement
- Coverage requirements validation

### Manual Review Requirements
- [ ] Tests follow established patterns
- [ ] No API mocking introduced
- [ ] Proper test isolation implemented
- [ ] Build tags correctly applied
- [ ] Real environment testing used

## References
- `TESTING_GUIDE.md` - Comprehensive implementation guide
- `scripts/validate-tests.sh` - Automated validation tooling
- `test/` directory structure - Implementation examples
- Go testing best practices documentation

## Revision History
- 2025-08-19: Initial decision record created
- Future: Updates will be documented here

---

**This ADR establishes the foundation for reliable, maintainable testing in the code-editing-agent project. All future testing decisions should align with these principles.**