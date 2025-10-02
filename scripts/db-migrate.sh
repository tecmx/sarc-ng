#!/bin/bash
# Database migration script for SARC-NG
# Handles database schema migrations using golang-migrate or custom tooling

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Get project root directory
PROJECT_ROOT="$(cd "$(dirname "$0")/.." && pwd)"

# Configuration
DB_HOST="${DB_HOST:-localhost}"
DB_PORT="${DB_PORT:-3306}"
DB_USER="${DB_USER:-root}"
DB_PASSWORD="${DB_PASSWORD:-example}"
DB_NAME="${DB_NAME:-sarcng}"
MIGRATIONS_DIR="${PROJECT_ROOT}/migrations"

# Build DSN
DSN="mysql://${DB_USER}:${DB_PASSWORD}@tcp(${DB_HOST}:${DB_PORT})/${DB_NAME}"

show_help() {
    echo -e "${GREEN}Database Migration Tool${NC}"
    echo ""
    echo "Usage: $0 <command> [options]"
    echo ""
    echo "Commands:"
    echo "  up              Apply all pending migrations"
    echo "  down            Rollback last migration"
    echo "  create NAME     Create new migration file"
    echo "  status          Show migration status"
    echo "  version         Show current migration version"
    echo "  force VERSION   Force set version (use with caution)"
    echo ""
    echo "Options:"
    echo "  -h, --help      Show this help message"
    echo ""
    echo "Environment Variables:"
    echo "  DB_HOST         Database host (default: localhost)"
    echo "  DB_PORT         Database port (default: 3306)"
    echo "  DB_USER         Database user (default: root)"
    echo "  DB_PASSWORD     Database password (default: example)"
    echo "  DB_NAME         Database name (default: sarcng)"
    echo ""
    echo "Examples:"
    echo "  $0 up"
    echo "  $0 create add_users_table"
    echo "  $0 status"
}

check_migrate_tool() {
    if ! command -v migrate &> /dev/null; then
        echo -e "${RED}Error: golang-migrate tool not found${NC}"
        echo "Install it with: go install -tags 'mysql' github.com/golang-migrate/migrate/v4/cmd/migrate@latest"
        exit 3
    fi
}

migrate_up() {
    echo -e "${GREEN}Applying all pending migrations...${NC}"
    migrate -path "$MIGRATIONS_DIR" -database "$DSN" up
    echo -e "${GREEN}✓ Migrations applied successfully${NC}"
}

migrate_down() {
    echo -e "${YELLOW}Rolling back last migration...${NC}"
    migrate -path "$MIGRATIONS_DIR" -database "$DSN" down 1
    echo -e "${GREEN}✓ Migration rolled back${NC}"
}

migrate_create() {
    local name="$1"
    if [ -z "$name" ]; then
        echo -e "${RED}Error: Migration name required${NC}"
        echo "Usage: $0 create <migration_name>"
        exit 2
    fi

    echo -e "${GREEN}Creating migration: $name${NC}"
    mkdir -p "$MIGRATIONS_DIR"
    migrate create -ext sql -dir "$MIGRATIONS_DIR" -seq "$name"
    echo -e "${GREEN}✓ Migration files created in $MIGRATIONS_DIR${NC}"
}

migrate_status() {
    echo -e "${GREEN}Migration Status:${NC}"
    migrate -path "$MIGRATIONS_DIR" -database "$DSN" version
}

migrate_version() {
    migrate -path "$MIGRATIONS_DIR" -database "$DSN" version
}

migrate_force() {
    local version="$1"
    if [ -z "$version" ]; then
        echo -e "${RED}Error: Version number required${NC}"
        echo "Usage: $0 force <version>"
        exit 2
    fi

    echo -e "${YELLOW}WARNING: Forcing migration version to $version${NC}"
    read -p "Are you sure? [y/N] " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        migrate -path "$MIGRATIONS_DIR" -database "$DSN" force "$version"
        echo -e "${GREEN}✓ Version forced to $version${NC}"
    else
        echo "Operation cancelled"
    fi
}

main() {
    local command="$1"
    shift || true

    case "$command" in
        -h|--help|help)
            show_help
            exit 0
            ;;
        up)
            check_migrate_tool
            migrate_up
            ;;
        down)
            check_migrate_tool
            migrate_down
            ;;
        create)
            check_migrate_tool
            migrate_create "$@"
            ;;
        status)
            check_migrate_tool
            migrate_status
            ;;
        version)
            check_migrate_tool
            migrate_version
            ;;
        force)
            check_migrate_tool
            migrate_force "$@"
            ;;
        "")
            echo -e "${RED}Error: Command required${NC}"
            show_help
            exit 2
            ;;
        *)
            echo -e "${RED}Error: Unknown command '$command'${NC}"
            show_help
            exit 2
            ;;
    esac
}

main "$@"

