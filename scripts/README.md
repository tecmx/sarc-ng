# Scripts Directory

This directory contains operational scripts for the SARC-NG project.

## Scripts Overview

### Database Management

#### `db-migrate.sh`

**Purpose**: Run database migrations.

**Usage**:

```bash
./db-migrate.sh up           # Apply all pending migrations
./db-migrate.sh down         # Rollback last migration
./db-migrate.sh create NAME  # Create new migration
./db-migrate.sh status       # Show migration status
```

---

#### `db-seed.sh`

**Purpose**: Seed database with test data.

**Usage**:

```bash
./db-seed.sh                 # Seed with default test data
./db-seed.sh --env dev       # Seed for specific environment
./db-seed.sh --clean         # Clean and reseed database
```

---

### Health & Validation

#### `health-check.sh`

**Purpose**: Validate service health and readiness.

**Usage**:

```bash
./health-check.sh            # Check all services
./health-check.sh --api      # Check API only
./health-check.sh --db       # Check database only
./health-check.sh --wait     # Wait for services to be ready
```

**Exit codes**:

- `0` - All services healthy
- `1` - One or more services unhealthy
- `2` - Timeout waiting for services

---

### Development Utilities

#### `version-bump.sh`

**Purpose**: Semantic version management.

**Usage**:

```bash
./version-bump.sh major      # 1.0.0 -> 2.0.0
./version-bump.sh minor      # 1.0.0 -> 1.1.0
./version-bump.sh patch      # 1.0.0 -> 1.0.1
./version-bump.sh --dry-run  # Show what would change
```

**Features**:

- Updates `VERSION` file
- Creates git tag
- Updates relevant config files

---

## Script Conventions

### Common Options

All scripts support these standard options:

- `-h, --help` - Show help message
- `-v, --verbose` - Enable verbose output
- `-d, --dry-run` - Show what would happen without executing

### Exit Codes

- `0` - Success
- `1` - General error
- `2` - Usage error (invalid arguments)
- `3` - Environment error (missing dependencies)

### Error Handling

All scripts:

- Use `set -e` to exit on errors
- Validate dependencies before execution
- Provide clear error messages
- Clean up on failure

## Adding New Scripts

When adding a new script:

1. **Use the template structure**:

   ```bash
   #!/bin/bash
   set -e

   # Script description
   # Usage: ./script-name.sh [options]

   # Colors for output
   RED='\033[0;31m'
   GREEN='\033[0;32m'
   YELLOW='\033[1;33m'
   NC='\033[0m' # No Color

   # Main function
   main() {
       # Script logic
   }

   main "$@"
   ```

2. **Make it executable**:

   ```bash
   chmod +x scripts/new-script.sh
   ```

3. **Update this README** with documentation

4. **Add to Makefile** if it should be called from make targets

5. **Test thoroughly** before committing

## Best Practices

1. **Always validate inputs** before executing commands
2. **Provide helpful error messages** with context
3. **Use colored output** for better readability
4. **Check for required tools** before running
5. **Make scripts idempotent** when possible
6. **Document exit codes** and their meanings
7. **Support dry-run mode** for safety
8. **Clean up temporary files** on exit

## Integration with Makefile

Scripts are called from the infrastructure Makefile:

```makefile
# In infrastructure/Makefile
SCRIPTS_DIR := ../scripts

health-check:
    @$(SCRIPTS_DIR)/health-check.sh --all

db-migrate:
    @$(SCRIPTS_DIR)/db-migrate.sh up
```

This provides:

- Single source of truth for script locations
- Easy testing of scripts independently
- Integration with make-based workflows
