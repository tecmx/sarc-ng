# AI Agent Rules for sarc-ng

Reference playbook for agents collaborating on the sarc-ng platform (Go services, Docusaurus docs, and IaC modules).

## **Quick Start**

Rules are organized by priority and apply based on their description conditions.

## **Rule Numbering System**

Rules are organized by **prefix numbers** that indicate their scope and application:

### **00- Universal Development Rules**

Apply to ALL tasks in this repository unless a rule specifies narrower scope:

- **00-code-quality** → SOLID, Clean Code, and pragmatic pattern usage (always applied)
- **00-development-principles** → Senior-level execution, completeness, readable diffs (always applied)
- **00-function-naming** → Shared naming conventions (apply when writing/refactoring functions)
- **00-python-standards** → Modern Python practices (only when touching Python tooling/scripts)
- **00-rule-creation** → Guidance for expanding this rule set


## **Rule Categories**

### **Universal Development (00-)**

#### [00-code-quality.mdc](./00-code-quality.mdc)

- **Scope**: ALL development tasks (always applied)
- **Purpose**: Keep Go services, TypeScript tooling, and IaC modules simple, SOLID, and maintainable

#### [00-development-principles.mdc](./00-development-principles.mdc)

- **Scope**: ALL development tasks
- **Purpose**: Deliver end-to-end slices, wire tests and docs, and remove dead code as you go

#### [00-function-naming.mdc](./00-function-naming.mdc)

- **Scope**: Writing functions or refactoring across Go packages, Node tooling, or scripts
- **Purpose**: Encourage intent-revealing verb-first names (`get`, `build`, `sync`, `validate`)

#### [00-rule-creation.mdc](./00-rule-creation.mdc)

- **Scope**: Creating new rules, modifying existing rules, or working with rule system structure
- **Purpose**: Add project-specific rules (e.g., Go module layering, Docusaurus docs patterns) without breaking format

#### [00-python-standards.mdc](./00-python-standards.mdc)

- **Scope**: Writing Python code, using Python libraries, or working on Python-specific features
- **Purpose**: Keep the occasional Python helper scripts consistent with typing, packaging, and linting best practices

> **Heads-up**: The Go codebase does not currently have a dedicated rule file—lean on the always-on quality and development standards. When Go-specific guidance is needed, add a new `01-go-*` rule via `00-rule-creation`.


______________________________________________________________________

## **Rule Application Levels**

1. **Always Applied** (`alwaysApply: true`) - `00-code-quality` and `00-development-principles` cover all contributions
1. **Request-based** (`alwaysApply: false`) - Opt-in rules like `@00-function-naming` or `@00-python-standards` activate when you call them out

______________________________________________________________________

## **Usage Examples**

```bash
# Trigger specific rules
@00-function-naming when writing new functions
@00-python-standards when working with Python code (rare in this repo)
@00-rule-creation when creating or modifying rules (e.g., adding Go-specific guidance)

# Development principles and code quality are always applied automatically
# No need to request explicitly
```

## **Rule Compliance Checklist**

### **Implementation Checklist**

Before completing any implementation:

- [ ] Core principles applied (always active for all development tasks)
- [ ] Code quality standards followed (SOLID, Clean Code - always active)
- [ ] Function reuse check completed (mandatory for all coding)
- [ ] Naming conventions verified (applied when working with functions)

## **Critical Success Factors**

1. **Minimal always-active rules** - keeps agents fast while enforcing essentials on Go/TypeScript/Terraform work
1. **Description-based application** - trigger optional rules when diving into Python or new rule authoring
1. **Clear rule hierarchy** - universal rules first, then language/domain add-ons
1. **Clear separation** - add new numbered families (e.g., `01-go-services`, `02-infrastructure`) when deeper guidance is required

**Remember**: Rules are living documents that evolve with the platform!

______________________________________________________________________

## **CREATING NEW RULES**

For creating new rules, see **[00-rule-creation.mdc](./00-rule-creation.mdc)** - comprehensive guidelines tailored for this repo's automation flow.

**Quick Standards:**

- **Target Length**: 80-120 lines (sweet spot: ~100 lines)
- **Structure**: Quick Reference → Core Principles → Essential Patterns → Enforcement Checklist
- **Examples**: 2-4 focused examples with ✅/❌ contrasts
- **Consistency**: Follow template exactly for predictable AI processing
