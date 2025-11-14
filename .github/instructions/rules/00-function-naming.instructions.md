---
applyTo: '**'
name: 00-function-naming
description: Apply when writing functions, refactoring code, or working on naming conventions across the codebase
---

# Universal Naming Conventions for Functions and Assets

> **PRIORITY**: HIGHEST - Consistent naming is critical for maintainability

## **CORE NAMING PRINCIPLES**

### **CONSISTENCY OVER CLEVERNESS**

- Use predictable patterns that developers can easily remember
- Prefer explicit names over abbreviated ones
- Maintain consistency within functional areas

### **SEMANTIC CLARITY**

- Function names should clearly indicate their purpose and return type
- Use domain-specific terminology when appropriate
- Avoid ambiguous or generic terms

### **HIERARCHICAL ORGANIZATION**

- Group related functions using consistent prefixes
- Follow the verb-noun pattern for actions
- Use consistent suffixes for similar operations

## **FUNCTION NAMING PATTERNS**

### **Common Patterns**

- **Calculation**: `calculate_<what>_<type>` (e.g., `calculate_total_score`)
- **Retrieval**: `get_<what>_<type>` (e.g., `get_user_profile`)
- **Validation**: `validate_<what>_<aspect>` (e.g., `validate_email_format`)
- **Generation**: `generate_<what>_<type>` (e.g., `generate_report_summary`)
- **Normalization**: `normalize_<what>_<aspect>` (e.g., `normalize_phone_number`)
- **Creation**: `create_<what>_<type>` (e.g., `create_user_session`)

### **Feature Engineering**

- **Extraction**: `extract_<what>_<type>`
- **Encoding**: `encode_<what>_<method>`

### **Model Operations**

- **Training**: `train_<model_type>_<purpose>` (e.g., `train_classifier_sentiment`)
- **Evaluation**: `evaluate_<model_type>_<metric>` (e.g., `evaluate_regression_rmse`)
- **Prediction**: `predict_<what>_<method>` (e.g., `predict_sales_forecast`)

### **I/O Operations**

- **Pattern**: `<action>_<what>_<object>`
- **Reading**: `read_<format>_<source>` (e.g., `read_csv_file`, `read_json_api`)
- **Writing**: `write_<format>_<destination>` (e.g., `write_parquet_s3`)
- **Loading**: `load_<what>_<from>` (e.g., `load_config_file`)

### **Business Logic Patterns**

- **Processing**: `process_<entity>_<operation>` (e.g., `process_order_payment`)
- **Transformation**: `transform_<from>_to_<to>` (e.g., `transform_raw_to_clean`)
- **Aggregation**: `aggregate_<what>_by_<dimension>` (e.g., `aggregate_sales_by_region`)

## **NAMING CHECKLIST**

Before implementing any function:

- [ ] Follows established verb_noun pattern
- [ ] Uses appropriate prefix (`calculate_`, `get_`, `validate_`, etc.)
- [ ] Consistent with existing naming in the same module
- [ ] Clear and unambiguous purpose
- [ ] No mixing of naming patterns within functional areas

### **Common Anti-Patterns to Avoid**

- AVOID Generic names: `process_data()`, `handle_stuff()`, `do_work()`
- AVOID Abbreviations: `calc_tot_amt()`, `proc_usr_req()`
- AVOID Mixed conventions: `getUserData()` + `validate_email_format()`
- PREFER Specific purpose: `calculate_total_amount()`, `process_user_request()`

**CRITICAL**: Naming consistency is essential for code maintainability and developer productivity!
