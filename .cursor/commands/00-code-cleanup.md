# Code Cleanup

> **COMMAND**: `@agent 00-code-cleanup`
> **CATEGORY**: Cleanup
> **AUTOFIX**: 🔧 Manual (requires --confirm flag)
> **SCOPE**: Universal (any language/framework)
> **DEPENDENCIES**: `.cursor/rules/00-dev-principles.mdc`
> **PORTABILITY**: 🟢 100%

## Purpose

Find and report legacy code, unused code, empty folders, backward compatibility code, and auto-generated documentation files in source code only (excludes build artifacts, cache, and virtualenvs).

**When to Use:**

- ✅ Before major releases
- ✅ During code refactoring sessions
- ✅ Monthly/quarterly code health reviews
- ❌ Don't use without reviewing dry-run first

**Expected Duration:** 2-5 minutes

---

## Usage

```bash
# Run scan (always interactive)
@agent 00-code-cleanup
```

### Examples

```bash
# Standard usage - scan and prompt for fixes
@agent 00-code-cleanup

# Expected flow:
# 1. Shows scan results
# 2. "Would you like me to fix these issues? (yes/no)"
# 3. Applies fixes if yes, exits if no
```

---

## Parameters

This command has no parameters - it's always interactive and asks for confirmation before making changes.

---

## What It Does

Scans the codebase for code quality issues and maintenance debt, then **interactively prompts to apply fixes**.

**Scans for:**

1. **Legacy Code** - Deprecated functions, old patterns, commented-out code blocks
2. **Unused Code** - Unused imports, unused functions, unused variables, dead code paths
3. **Empty Folders** - Directories containing no files whatsoever (excluding excluded paths)
4. **Backward Compatibility** - Legacy compatibility layers, version checks, deprecated flags
5. **Auto-generated Docs** - `*_SUMMARY.md`, `*_PLAN.md`, `*_GUIDE.md` files

**Excludes (CRITICAL - Never scan these):**

- **Build Artifacts**: `__pycache__/`, `*.pyc`, `cdk.out/`, `dist/`, `build/`, `node_modules/`, `target/`, compiled binaries
- **Virtual Environments**: `.venv/`, `venv/`, `env/`, `ENV/`, virtualenv directories
- **Coverage/Test Artifacts**: `htmlcov/`, `.coverage`, `coverage.json`, `.pytest_cache/`, `.tox/`
- **Cache Directories**: `.mypy_cache/`, `.ruff_cache/`, `.cache/`, `__pycache__/`
- **Version Control**: `.git/`, `.svn/`, `.hg/`
- **IDE**: `.vscode/`, `.idea/`, `*.swp`

**Important:**

- ✅ **Always interactive** - Shows results, then asks before applying fixes
- ✅ **Proactive** - Offers to fix issues immediately after detection
- ✅ **Source code only** - Scans actual code files, excludes build/cache/virtualenv
- ❌ **Never auto-applies** - Always asks for confirmation first
- ❌ **Never generates summary/md files** - Only reports findings

**Standards Applied:**

- `.cursor/rules/00-dev-principles.mdc` - No backward compatibility code, complete implementations
- User Rule: "Do not maintain or preserve any backward compatibility code"
- User Rule: "Always double-check and complete all implementations"

---

## Workflow

1. **Scan** – Walk the source tree (excluding artifacts/caches) to locate legacy, unused, or auto-generated code.
2. **Report** – Present grouped findings with recommended actions so you can review the dry run.
3. **Prompt** – Ask `Would you like me to fix these issues? (yes/no)` before touching files.
4. **Apply** – When you answer yes, remove the flagged items immediately; answering no exits without changes.

---

## Integration

### Makefile

```makefile
code-cleanup: ## Interactive code cleanup (scan + prompt to fix)
 @cursor-agent --print --model sonnet-4.5 \
  "@agent 00-code-cleanup"
```

### CI/CD (GitHub Actions)

```yaml
# Note: Not recommended for CI/CD - this command is interactive
# For CI/CD, use language-specific linting tools instead
- name: Check Code Quality
  run: |
    make lint
    make type-check
```

---

## Related Commands

- `@agent 00-docs-check` - Validate documentation quality
- `@agent /test` - Run test suite

---

**Last Updated:** 2025-01-12
**Compatibility:** Cursor Agent CLI + IDE
**Model Tested:** Claude Sonnet 4.5
**Exit Codes:** 0 (success), 1 (files found), 2 (errors)
