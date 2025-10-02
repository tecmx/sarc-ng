---
applyTo: '**'
name: 00-development-principles
tags:
    - always-apply
---

# Core Development Principles for Cursor AI Agents

> **PRIORITY**: HIGHEST - These principles apply to ALL development tasks

## **QUICK REFERENCE**

- Write minimal, clean code with single-line function comments
- Always check for existing functions before implementing new ones
- Use real production data only - never sample/mock data
- Eliminate legacy code paths and use modern APIs
- Stay within exact scope - no unrequested features or test structures
- Never use emojis in code, comments, or documentation

## **FUNDAMENTAL PRINCIPLES**

### **THE FEWER LINES OF CODE, THE BETTER**

- Prioritize concise, readable solutions
- Eliminate redundancy and unnecessary complexity
- Choose simple implementations over clever ones

### **PROCEED LIKE A SENIOR DEVELOPER**

- Think before coding - write 2-3 reasoning paragraphs first
- Consider long-term maintainability
- Make architectural decisions with experience in mind
- Anticipate edge cases and future requirements

### **DO NOT STOP UNTIL COMPLETE**

- Ensure all requirements are met
- Test the implementation thoroughly
- Handle edge cases and error conditions
- Provide complete, working solutions

## **DEVELOPMENT PROCESS**

### **REASONING-FIRST APPROACH**

**ALWAYS START WITH THREE REASONING PARAGRAPHS** before implementing:

1. **Initial observations**: "It's possible that..." - start with uncertainty
2. **Deeper analysis**: Consider alternative explanations and approaches
3. **Most likely solution**: Gradually build confidence with supporting evidence

### **MANDATORY FUNCTION REUSE CHECK**

**CRITICAL:** Before implementing ANY functionality, check existing codebase for similar functions.

**Required Questions:**

- Does this functionality already exist in the codebase?
- Can I extend an existing function instead of creating new one?
- Where should this belong in the existing structure?

## **CODE QUALITY STANDARDS**

### **CLEAN CODE REQUIREMENTS**

- Write clean, simple, readable code
- Implement features in the simplest possible way
- Keep files small and focused (<200 lines)
- Use clear, consistent naming
- Never use emojis in any code, comments, or documentation

### **COMMENTING STANDARDS**

- **Every function MUST have a single-line comment** describing its purpose
- **Avoid comments inside function bodies** unless absolutely critical
- **Let code be self-documenting** through clear variable and function names
- **Only add inline comments** for complex business logic or non-obvious algorithms

### **DATA INTEGRITY REQUIREMENTS**

- **NEVER use mock or sample data** - Always use real production data
- **ALWAYS prompt user for real input data** if missing - Never make assumptions
- **STOP and ask for clarification** if essential data is unavailable
- **NO placeholder data** - No "example", "test", or "dummy" data in implementations

## **SCOPE & MODERNIZATION MANAGEMENT**

### **SCOPE MANAGEMENT**

- **NEVER add features not explicitly requested** - even if they seem helpful
- **ASK BEFORE implementing additional functionality** beyond what was requested
- **STICK TO THE EXACT SCOPE** - implement only what the user asked for
- **SUGGEST, DON'T ASSUME** - if you see opportunities for improvement, suggest them and wait for approval
- **DO NOT create test structures** unless explicitly requested - no automatic test scaffolding

### **LEGACY CODE ELIMINATION**

- **ALWAYS eliminate legacy or backward-compatibility code paths**
- **Strip out deprecated APIs, shims, polyfills, and version-specific branches**
- **Replace them with current, officially supported APIs and patterns**
- **Do not keep fallbacks, wrappers, feature flags, or "old-version" conditionals** unless explicitly requested
- **Modernize project dependencies accordingly and delete unused packages**

## **ERROR HANDLING**

### **SYSTEMATIC APPROACH**

- **DO NOT JUMP TO CONCLUSIONS** - Consider multiple possible causes before deciding
- **EXPLAIN IN PLAIN ENGLISH** - Make problems understandable
- **MINIMAL CHANGES ONLY** - Change as few lines as possible
- **TEST EACH CHANGE** - Verify each modification works

### **ERROR FIXING WORKFLOW**

1. **Analyze**: Write reasoning paragraphs about potential causes
2. **Minimal Changes**: Change as few lines as possible
3. **Test**: Verify each change works

## **ENFORCEMENT CHECKLIST**

### **Process & Planning**

- [ ] Followed three-paragraph reasoning process
- [ ] Checked for existing functions in codebase before implementing

### **Code Quality**

- [ ] Used simple, clean implementation
- [ ] Every function has a single-line comment describing its purpose
- [ ] Avoided unnecessary comments in function bodies
- [ ] Eliminated legacy code paths and used modern APIs

### **Data & Testing**

- [ ] Used real production data (NEVER sample/mock data)
- [ ] Tested functionality thoroughly
- [ ] Handled edge cases appropriately

### **Scope Management**

- [ ] Stayed within exact scope requested (no unrequested features)
- [ ] Did NOT create test structures unless explicitly requested

**CRITICAL**: These principles take precedence over all other implementation guidelines!
