# Architecture Check

> **COMMAND**: `@agent 00-architecture-check`
> **CATEGORY**: Architecture
> **AUTOFIX**: ❌ None (requires manual refactoring)
> **SCOPE**: Universal (any software architecture)
> **DEPENDENCIES**: Project structure only
> **PORTABILITY**: 🟢 100%

## Purpose

Validate architectural layer separation and dependency rules. Detects cross-layer violations, circular dependencies, and incorrect dependency directions.

**When to Use:**

- ✅ After major refactoring
- ✅ Monthly architecture reviews
- ✅ Before production releases
- ✅ Onboarding new developers

**Expected Duration:** 5-10 minutes

---

## Usage

```bash
# Auto-detect architecture
@agent 00-architecture-check

# Specify architecture type
@agent 00-architecture-check [ARCHITECTURE-TYPE]
```

### Examples

```bash
# Auto-detect architecture
@agent 00-architecture-check auto

# Layered architecture
@agent 00-architecture-check layered

# Hexagonal architecture
@agent 00-architecture-check hexagonal

# Clean architecture
@agent 00-architecture-check clean
```

---

## Parameters

| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `architecture-type` | enum | `auto` | Architecture: `auto` \| `layered` \| `hexagonal` \| `clean` \| `mvc` \| `mvvm` |

**Architecture Types:**

- `auto`: Detect from project structure (recommended)
- `layered`: Traditional N-tier (presentation → business → data)
- `hexagonal`: Ports & adapters
- `clean`: Clean architecture (entities → use cases → adapters → frameworks)
- `mvc`: Model-View-Controller
- `mvvm`: Model-View-ViewModel

---

## What It Does

Validates layer boundaries by checking for cross-layer violations, circular dependencies, dependency direction (should flow inward toward domain), and layer responsibility separation. Reports violations with specific file locations and suggested fixes.

**Important:**

- ❌ **Never auto-applies fixes** - Always asks user for confirmation before making any changes
- ❌ **Never generates summary/md files** - Only reports findings, does not create documentation files

**Standards Applied:**

- Clean Architecture principles
- Dependency Inversion Principle
- Layer separation best practices

---

## Workflow

1. **Scan** – Parse the repository structure (or specified architecture type) to build a dependency graph between layers/modules.
2. **Report** – Produce a detailed list of violations (cross-layer calls, circular deps, outward-facing dependencies) with file paths and suggested fixes.
3. **Prompt** – Not applicable: this command is analysis-only and stops after reporting; it never attempts automated remediation.
4. **Apply** – No automatic changes occur; use the report to refactor manually before rerunning the checker.

---

## Integration

### Makefile

```makefile
arch-check: ## Check architecture boundaries
 @cursor-agent --print --model sonnet-4.5 \
  "@agent 00-architecture-check auto" || exit 1
```

---

## Related Commands

None

---

**Last Updated:** 2025-01-12
**Compatibility:** Cursor Agent CLI + IDE
**Model Tested:** Claude Sonnet 4.5
**Exit Codes:** 0 (pass), 1 (warnings), 2 (violations)
