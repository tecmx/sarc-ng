# Scripts

Operational scripts for SARC-NG.

## Available Scripts

### `db-migrate.sh`
Database migrations management.

```bash
./db-migrate.sh up           # Apply migrations
./db-migrate.sh down         # Rollback
./db-migrate.sh create NAME  # Create new migration
./db-migrate.sh status       # Check status
```

### `db-seed.sh`
Seed database with test data.

```bash
./db-seed.sh                 # Seed default data
./db-seed.sh --env dev       # Environment-specific
./db-seed.sh --clean         # Clean and reseed
```

### `health-check.sh`
Service health validation.

```bash
./health-check.sh            # Check all services
./health-check.sh --api      # API only
./health-check.sh --db       # Database only
./health-check.sh --wait     # Wait for ready
```

**Exit codes:** 0=healthy, 1=unhealthy, 2=timeout

### `version-bump.sh`
Version management.

```bash
./version-bump.sh major      # 1.0.0 -> 2.0.0
./version-bump.sh minor      # 1.0.0 -> 1.1.0
./version-bump.sh patch      # 1.0.0 -> 1.0.1
./version-bump.sh --dry-run  # Preview changes
```

## Common Options

- `-h, --help` - Show help
- `-v, --verbose` - Verbose output
- `-d, --dry-run` - Preview without executing

## Adding Scripts

1. Use standard structure with error handling
2. Make executable: `chmod +x scripts/new-script.sh`
3. Update this README
4. Add to Makefile if needed
5. Test thoroughly
