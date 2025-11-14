# Test Runner

> **COMMAND**: `@agent /test`
> **CATEGORY**: Testing
> **AUTOFIX**: ❌ N/A (executes tests)
> **SCOPE**: Universal (any test framework)
> **DEPENDENCIES**: Test framework, build tool (make/npm/go/etc)
> **PORTABILITY**: 🟢 100%

## Purpose

Run unit and integration tests (standard test suite). Excludes E2E tests which may require external infrastructure deployment.

**When to Use:**

- ✅ Pre-commit validation
- ✅ Standard development testing
- ✅ Quick validation before PR

**Expected Duration:** ~2 minutes

---

## Usage

```bash
@agent /test
```

### Examples

```bash
# Run unit and integration tests
@agent /test
```

---

## Parameters

None (uses test framework markers/tags to filter tests)

**Test Organization:**

- Unit tests - Fast, isolated tests
- Integration tests - Component interaction tests (may use mocks/stubs)

**Note:** E2E tests are excluded. Use project-specific commands (e.g., `make test-e2e`) for end-to-end testing.

---

## What It Does

Executes test framework with markers/tags for unit and integration tests. Uses project's test configuration and mocking strategy. Does not deploy infrastructure or require external services.

**Important:**

- ❌ **Never generates summary/md files** - Only reports test results, does not create documentation files

**Standards Applied:**

- `.cursor/rules/00-testing-standards.mdc` - Test organization, patterns, and quality standards

---

## Workflow

1. **Scan** – Detect the project’s native test runner (Makefile target, `npm test`, `go test`, etc.) and determine which suites (unit + integration) need to run.
2. **Report** – Stream the test output directly, surfacing failures with the original framework formatting.
3. **Prompt** – Not applicable: the command executes immediately with no confirmation step because it never mutates code.
4. **Apply** – Not applicable: there are no automated fixes; resolving test failures is left to the developer before rerunning `/test`.

---

## Integration

### Makefile (Example)

```makefile
test: ## Run unit and integration tests
 @[test-command] [test-args]
```

The Makefile should provide convenience targets:

- `make test-unit` - Unit tests only
- `make test-integration` - Integration tests only
- `make test-e2e` - E2E tests (if applicable)
- `make test-all` - All tests including E2E

---

## Related Commands

- `@agent 00-universal-testing` - Analyze test framework quality

---

**Last Updated:** 2025-01-12
**Compatibility:** Cursor Agent CLI + IDE
**Model Tested:** Claude Sonnet 4.5
**Exit Codes:** 0 (pass), 1 (test failures)
