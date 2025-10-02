---
title: Modern Python Language Standards and Best Practices
sidebar_label: Modern Python Language Standards and Best Practices
sidebar_position: 0
tags:
    - ai-rules
agent:
    description: Apply when writing Python code, using Python libraries, or working on Python-specific features and patterns
    alwaysApply: false
---

# Modern Python Language Standards and Best Practices

> **PRIORITY**: HIGH - Apply to all Python development tasks

## **QUICK REFERENCE**

- Use `from __future__ import annotations` + built-in types (`list[str]`, `str | None`)
- Always add type hints with modern union syntax (`str | None` not `Optional[str]`)
- Use f-strings, context managers, and list/dict comprehensions
- Import only `Any` from typing - avoid `Dict`, `List`, `Optional`
- Follow PEP 8: standard library → third-party → local imports

## **CORE PYTHON PRINCIPLES**

### **MODERN PYTHON FEATURES FIRST**

- Use Python 3.9+ features (union types, built-in generics, etc.)
- Prefer modern syntax over legacy patterns
- Eliminate deprecated Python patterns and imports
- Use `from __future__ import annotations` for forward compatibility

### **TYPE SAFETY & CLARITY**

- Always use type hints for function parameters and return values
- Use modern union syntax (`str | None` instead of `Optional[str]`)
- Use built-in collection types (`list[str]` instead of `List[str]`)
- Import `Any` only when actually needed from typing

### **PYTHONIC CODE PATTERNS**

- Follow PEP 8 style guidelines
- Use list/dict comprehensions where appropriate
- Prefer context managers (`with` statements)
- Use f-strings for string formatting (never % or .format())

## **ESSENTIAL SYNTAX PATTERNS**

### **Modern Type Annotations & Imports**

```python
# CORRECT Modern
from __future__ import annotations
from typing import Any
from collections import defaultdict, Counter

def process_data(items: list[dict[str, Any]], key: str | None = None) -> dict[str, int]:
    """Process data with modern type hints."""
    if not items:
        raise ValueError("Items list cannot be empty")

    result: dict[str, int] = {}
    groups = defaultdict(list)
    for item in items:
        groups[item.get(key, 'default')].append(item)

    return {k: len(v) for k, v in groups.items()}

# AVOID Legacy
from typing import Dict, List, Optional
def process_data(items: List[Dict[str, Any]], key: Optional[str] = None) -> Dict[str, int]:
    pass
```

### **Exception Handling & Validation**

```python
# CORRECT Comprehensive patterns
def process_user_data(user_id: int, data: dict[str, Any]) -> dict[str, Any]:
    """Process user data with validation."""
    if user_id <= 0:
        raise ValueError(f"Invalid user_id: {user_id}, must be positive")
    if not isinstance(data, dict):
        raise TypeError(f"Expected dict, got {type(data).__name__}")

    try:
        with open(f"user_{user_id}.json") as f:
            existing_data = json.load(f)
    except (FileNotFoundError, json.JSONDecodeError) as e:
        logger.warning(f"Could not load existing data: {e}")
        existing_data = {}

    return {**existing_data, **data}
```

## **IMPORT ORGANIZATION & PERFORMANCE**

### **PEP 8 Import Order & Efficient Patterns**

```python
# 1. Standard library
import os
from pathlib import Path
from collections import defaultdict, Counter
from dataclasses import dataclass

# 2. Third-party packages
import requests
from dagster import asset

# 3. Local imports
from .utils import helper_function
from ..models import DataModel

# CORRECT Efficient data processing
@dataclass
class DataProcessor:
    """Process data with configurable settings."""
    source_path: str
    batch_size: int = 1000

    def process(self) -> list[dict[str, Any]]:
        """Process data from source with comprehensions and generators."""
        lines = (line.strip() for line in open(self.source_path) if line.strip())
        return [{'id': i, 'data': line.upper()} for i, line in enumerate(lines, 1) if len(line) > 5]
```

**Rules:**

- Use specific imports, never `from module import *`
- Separate imports on individual lines

## **PYTHON STANDARDS CHECKLIST**

### **Modern Syntax**

- [ ] Used `from __future__ import annotations` for forward compatibility
- [ ] Applied modern type hints (`list[str]`, `str | None`)
- [ ] Eliminated legacy typing imports (`List`, `Dict`, `Optional`)
- [ ] Used f-strings for string formatting

### **Code Quality**

- [ ] Applied proper exception handling with specific exceptions
- [ ] Used context managers for resource management
- [ ] Followed PEP 8 import organization (stdlib → third-party → local)
- [ ] Added clear docstrings with examples

### **Performance & Style**

- [ ] Used appropriate data structures (`defaultdict`, `Counter`)
- [ ] Applied comprehensions where beneficial
- [ ] Followed single responsibility principle for functions
- [ ] Used descriptive variable and function names

**CRITICAL**: Modern Python patterns improve readability, maintainability, and performance!
