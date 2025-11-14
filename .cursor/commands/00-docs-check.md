# Documentation Check

> **COMMAND**: `@agent 00-docs-check`
> **CATEGORY**: Documentation
> **AUTOFIX**: ⚠️ Partial
> **SCOPE**: Universal
> **DEPENDENCIES**: `.cursor/rules/00-readme-standards.mdc`
> **PORTABILITY**: 🟢 100%

## Purpose

Validate markdown documentation files for syntax correctness, link integrity, structure quality, and README standards compliance.

**When to Use:**
- ✅ Before committing markdown changes
- ✅ During documentation reviews
- ✅ As pre-commit hook for markdown files
- ✅ In CI/CD for documentation quality gates
- ❌ Don't use for content proofreading (checks structure, not prose)

**Expected Duration:** 2-5 minutes for full repository

---

## Usage

```bash
# Basic usage
@agent 00-docs-check

# With options
@agent 00-docs-check [path] --strictness=medium --fix-links --check-readme-standards
```

### Examples

```bash
# Check all markdown files
@agent 00-docs-check

# Check specific directory
@agent 00-docs-check docs/

# High strictness validation
@agent 00-docs-check --strictness=high

# Check with autofix for links
@agent 00-docs-check --fix-links --check-readme-standards

# Check single file
@agent 00-docs-check README.md
```

---

## Parameters

| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `path` | path | `.` (all .md files) | File or directory to check |
| `--strictness` | enum | `medium` | Validation level: `low` \| `medium` \| `high` |
| `--fix-links` | flag | `false` | Propose fixes for broken relative links (requires user confirmation) |
| `--check-readme-standards` | flag | `auto` | Apply README.md standards (auto-detected) |

**Strictness Levels:**
- `low`: Critical issues only (broken links, invalid syntax)
- `medium`: Standard checks (default - syntax, links, structure, basic formatting)
- `high`: Strict formatting and style (all checks including line length, emphasis consistency)

---

## What It Does

Validates markdown files for syntax correctness, link integrity, structure quality, and README standards compliance. Provides quality scoring (0-100) with actionable recommendations.

**Important:**
- ❌ **Never auto-applies fixes** - Always asks user for confirmation before making any changes
- ❌ **Never generates summary/md files** - Only reports findings, does not create documentation files

**Standards Applied:**
- `.cursor/rules/00-readme-standards.mdc` - README structure, length, and content guidelines

**Available Fixes (requires explicit `--fix-links` flag and user confirmation):**
- ✅ Heading hierarchy corrections
- ✅ Code block language tags
- ✅ Trailing whitespace removal
- ✅ List marker normalization
- ✅ Broken link fixes (high confidence)
- ⚠️ Manual: README length violations, external links, placeholder content

---

## Workflow

1. **Scan** – Traverse the specified path (default `.`) and parse every Markdown file using the selected strictness level.
2. **Report** – Output a quality score plus detailed findings (missing sections, broken links, structural issues).
3. **Prompt** – Only when `--fix-links` (or future fix flags) is supplied, ask whether to apply the proposed link/format fixes; without fix flags the command stops after reporting.
4. **Apply** – If you confirm, rewrite only the items covered by the enabled fix flags (e.g., fix relative links); decline to exit with no edits.

---

## Integration

### Makefile

```makefile
docs-check: ## Validate markdown documentation
	@cursor-agent --print --model sonnet-4.5 --output-format text \
		"@agent 00-docs-check --strictness=medium --check-readme-standards. \
		Provide quality score (0-100) with actionable recommendations. \
		Exit 0 if score >= 90, exit 1 if score < 90." || exit 1
```

### CI/CD (GitHub Actions)

```yaml
- name: Validate Documentation
  run: |
    cursor-agent --print --model sonnet-4.5 --output-format text \
      "@agent 00-docs-check --strictness=high --check-readme-standards"
```

---

## Related Commands

None

---

**Last Updated:** 2025-01-12
**Compatibility:** Cursor Agent CLI + IDE
**Model Tested:** Claude Sonnet 4.5
**Exit Codes:** 0 (pass), 1 (warnings), 2 (critical), 3 (autofix failed), 4 (invalid params)
