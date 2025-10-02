#!/bin/bash
# Health check script for SARC-NG services
# Validates service availability and readiness

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Configuration
API_URL="${API_URL:-http://localhost:8080}"
DB_HOST="${DB_HOST:-localhost}"
DB_PORT="${DB_PORT:-3306}"
DB_USER="${DB_USER:-root}"
DB_PASSWORD="${DB_PASSWORD:-example}"
TIMEOUT="${TIMEOUT:-30}"
WAIT_MODE=false
CHECK_ALL=true
CHECK_API=false
CHECK_DB=false

show_help() {
    echo -e "${GREEN}Service Health Check Tool${NC}"
    echo ""
    echo "Usage: $0 [options]"
    echo ""
    echo "Options:"
    echo "  --api           Check API health only"
    echo "  --db            Check database health only"
    echo "  --wait          Wait for services to be ready"
    echo "  --timeout SEC   Timeout in seconds (default: 30)"
    echo "  -h, --help      Show this help message"
    echo ""
    echo "Environment Variables:"
    echo "  API_URL         API base URL (default: http://localhost:8080)"
    echo "  DB_HOST         Database host (default: localhost)"
    echo "  DB_PORT         Database port (default: 3306)"
    echo "  DB_USER         Database user (default: root)"
    echo "  DB_PASSWORD     Database password (default: example)"
    echo "  TIMEOUT         Wait timeout in seconds (default: 30)"
    echo ""
    echo "Exit Codes:"
    echo "  0 - All services healthy"
    echo "  1 - One or more services unhealthy"
    echo "  2 - Timeout waiting for services"
    echo ""
    echo "Examples:"
    echo "  $0                    # Check all services"
    echo "  $0 --api              # Check API only"
    echo "  $0 --db               # Check database only"
    echo "  $0 --wait             # Wait for services to be ready"
    echo "  $0 --wait --timeout 60"
}

check_api_health() {
    local status
    local response

    if ! command -v curl &> /dev/null; then
        echo -e "${RED}✗ curl not installed${NC}"
        return 1
    fi

    response=$(curl -s -o /dev/null -w "%{http_code}" "${API_URL}/health" 2>/dev/null || echo "000")

    if [ "$response" = "200" ]; then
        echo -e "${GREEN}✓ API is healthy${NC} (${API_URL})"
        return 0
    else
        echo -e "${RED}✗ API is unhealthy${NC} (HTTP $response)"
        return 1
    fi
}

check_db_health() {
    if ! command -v mysqladmin &> /dev/null && ! command -v mysql &> /dev/null; then
        echo -e "${YELLOW}⚠ MySQL client not installed - checking port only${NC}"
        check_port "$DB_HOST" "$DB_PORT"
        return $?
    fi

    if mysqladmin ping -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" -p"$DB_PASSWORD" --silent 2>/dev/null; then
        echo -e "${GREEN}✓ Database is healthy${NC} ($DB_HOST:$DB_PORT)"
        return 0
    else
        echo -e "${RED}✗ Database is unhealthy${NC}"
        return 1
    fi
}

check_port() {
    local host="$1"
    local port="$2"

    if command -v nc &> /dev/null; then
        if nc -z -w 2 "$host" "$port" 2>/dev/null; then
            echo -e "${GREEN}✓ Port $port is reachable${NC} on $host"
            return 0
        else
            echo -e "${RED}✗ Port $port is not reachable${NC} on $host"
            return 1
        fi
    elif command -v timeout &> /dev/null; then
        if timeout 2 bash -c "cat < /dev/null > /dev/tcp/$host/$port" 2>/dev/null; then
            echo -e "${GREEN}✓ Port $port is reachable${NC} on $host"
            return 0
        else
            echo -e "${RED}✗ Port $port is not reachable${NC} on $host"
            return 1
        fi
    else
        echo -e "${YELLOW}⚠ Cannot check port - nc or timeout not available${NC}"
        return 1
    fi
}

wait_for_service() {
    local check_func="$1"
    local service_name="$2"
    local elapsed=0

    echo -e "${YELLOW}Waiting for $service_name (timeout: ${TIMEOUT}s)...${NC}"

    while [ $elapsed -lt "$TIMEOUT" ]; do
        if $check_func &>/dev/null; then
            echo -e "${GREEN}✓ $service_name is ready${NC} (after ${elapsed}s)"
            return 0
        fi
        sleep 2
        elapsed=$((elapsed + 2))
    done

    echo -e "${RED}✗ Timeout waiting for $service_name${NC} (${TIMEOUT}s)"
    return 2
}

main() {
    local exit_code=0

    # Parse arguments
    while [ $# -gt 0 ]; do
        case "$1" in
            --api)
                CHECK_API=true
                CHECK_ALL=false
                shift
                ;;
            --db)
                CHECK_DB=true
                CHECK_ALL=false
                shift
                ;;
            --wait)
                WAIT_MODE=true
                shift
                ;;
            --timeout)
                TIMEOUT="$2"
                shift 2
                ;;
            -h|--help)
                show_help
                exit 0
                ;;
            *)
                echo -e "${RED}Error: Unknown option '$1'${NC}"
                show_help
                exit 2
                ;;
        esac
    done

    echo -e "${GREEN}SARC-NG Service Health Check${NC}"
    echo ""

    # Determine what to check
    if $CHECK_ALL; then
        CHECK_API=true
        CHECK_DB=true
    fi

    # Run checks
    if $WAIT_MODE; then
        # Wait mode - keep trying until timeout
        if $CHECK_API; then
            wait_for_service check_api_health "API" || exit_code=$?
        fi
        if $CHECK_DB; then
            wait_for_service check_db_health "Database" || exit_code=$?
        fi
    else
        # Normal mode - check once
        if $CHECK_API; then
            check_api_health || exit_code=1
        fi
        if $CHECK_DB; then
            check_db_health || exit_code=1
        fi
    fi

    echo ""
    if [ $exit_code -eq 0 ]; then
        echo -e "${GREEN}✓ All checks passed${NC}"
    elif [ $exit_code -eq 2 ]; then
        echo -e "${RED}✗ Timeout waiting for services${NC}"
    else
        echo -e "${RED}✗ Some checks failed${NC}"
    fi

    exit $exit_code
}

main "$@"

