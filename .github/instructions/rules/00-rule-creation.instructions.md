---
applyTo: '**'
name: 00-rule-creation
description: Apply when creating new rules, modifying existing rules, or working with the rule system structure
---

# Rule Creation & Management Guidelines

> **PRIORITY**: HIGHEST - Essential for maintaining rule system integrity and AI agent performance

## **QUICK REFERENCE**

- Target 80-120 lines (sweet spot: ~100 lines) for optimal AI processing
- Use Quick Reference → Core Principles → Essential Patterns → Enforcement Checklist structure
- NEVER modify rule headers (YAML frontmatter) - use `search_replace` tool only
- Show 2-4 focused examples with CORRECT/AVOID contrasts
- Follow template exactly for predictable AI processing

## **CORE PRINCIPLES**

### **OPTIMAL AI PERFORMANCE FIRST**

- Target 80-120 lines total for maximum AI comprehension
- Use consistent structure across all rules for predictable processing
- Focus on actionable guidance over theoretical concepts
- Maintain logical flow from general principles to specific patterns

### **HEADER PROTECTION (CRITICAL)**

- **NEVER modify rule headers** (YAML frontmatter)
- **NEVER change `description` or `alwaysApply` values**
- Use `search_replace` tool ONLY for content changes below header
- Headers control rule discovery and application logic

## **ESSENTIAL PATTERNS**

### **Rule Structure Template**

```markdown
---
description: [Clear context-based trigger description]
alwaysApply: [true for universal rules, false for context-specific]
---

# Rule Title

> **PRIORITY**: [HIGHEST/HIGH/MEDIUM] - [Context trigger phrase]

## **QUICK REFERENCE**
- Key principle 1 (most critical)
- Key principle 2 (essential pattern)
- Key principle 3 (critical rule)

## **CORE PRINCIPLES**
### **PRIMARY PRINCIPLE**
- Clear, actionable guideline
- Supporting detail

## **ESSENTIAL PATTERNS**
### **Pattern Type**
```[language]
# CORRECT approach
example_code_or_pattern

# AVOID this
counter_example
```

## **ENFORCEMENT CHECKLIST**

### **Category**

- [ ] Specific, testable requirement
- [ ] Clear verification criteria

**CRITICAL**: Final emphasis line!

```

### **Safe Modification Process**

```bash
# Step 1: Always read first
read_file .cursor/rules/rule-name.mdc

# Step 2: Use search_replace ONLY (never edit_file)
search_replace --file rule-name.mdc --old "content below header" --new "updated content"

# Step 3: Verify headers remain intact
```

**Rules:**

- Always use `read_file` before making changes
- Only modify content below YAML headers
- Never use `edit_file` tool on rule files

## **WRITING GUIDELINES**

### **Quick References (5-8 bullets)**

- Start with the most critical principle
- Use action-oriented language ("Use X", "Always Y", "Never Z")
- Include specific technical details (not just concepts)
- Limit to essential information only

### **Core Principles (20-40 lines)**

- Group related concepts under clear headings
- Use consistent bullet formatting (not paragraphs)
- Focus on "why" and "what", save "how" for patterns section
- Maintain logical flow from general to specific

### **Essential Patterns (30-50 lines)**

- Show contrast with CORRECT/AVOID patterns
- Include real, applicable code/configurations
- Demonstrate principles in action
- Avoid repetitive or overly similar examples

### **Enforcement Checklists (10-15 lines)**

- Use specific, testable criteria
- Group related items under logical categories
- Focus on verification steps, not just reminders
- Include both technical and process requirements

## **QUALITY TARGETS**

- **80-90 lines**: Focused, specific rules
- **90-110 lines**: Comprehensive domain rules
- **110-120 lines**: Complex architectural rules

## **RULE CREATION CHECKLIST**

### **Structure & Content**

- [ ] Length within 80-120 line target range
- [ ] Quick Reference has 5-8 actionable bullets
- [ ] Core Principles use consistent bullet formatting
- [ ] Essential Patterns show clear CORRECT/AVOID contrasts
- [ ] Enforcement Checklist has testable criteria

### **Header & System Compliance**

- [ ] Header preserved with correct description/alwaysApply
- [ ] Used `search_replace` tool only (never `edit_file`)
- [ ] Verified header integrity after modifications
- [ ] Final CRITICAL line reinforces key message

### **AI Optimization**

- [ ] Follows template exactly for predictable processing
- [ ] Consistent structure with existing rules
- [ ] Actionable guidance over theoretical concepts
- [ ] Logical flow from principles to implementation

**CRITICAL**: Rule headers are the foundation of Cursor's rule system - treat as immutable infrastructure while optimizing content for maximum AI agent performance!
