#!/bin/bash
# Database seeding script for SARC-NG
# Populates database with test data

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
SEED_DIR="${PROJECT_ROOT}/test/fixtures"
ENV="${ENV:-dev}"
CLEAN_MODE=false

show_help() {
    echo -e "${GREEN}Database Seeding Tool${NC}"
    echo ""
    echo "Usage: $0 [options]"
    echo ""
    echo "Options:"
    echo "  --env ENV       Environment (dev, test, prod) (default: dev)"
    echo "  --clean         Clean database before seeding"
    echo "  -h, --help      Show this help message"
    echo ""
    echo "Environment Variables:"
    echo "  DB_HOST         Database host (default: localhost)"
    echo "  DB_PORT         Database port (default: 3306)"
    echo "  DB_USER         Database user (default: root)"
    echo "  DB_PASSWORD     Database password (default: example)"
    echo "  DB_NAME         Database name (default: sarcng)"
    echo "  ENV             Environment (default: dev)"
    echo ""
    echo "Examples:"
    echo "  $0                      # Seed with default dev data"
    echo "  $0 --clean              # Clean and reseed"
    echo "  $0 --env test           # Seed test environment"
}

check_mysql_client() {
    if ! command -v mysql &> /dev/null; then
        echo -e "${RED}Error: mysql client not installed${NC}"
        exit 3
    fi
}

mysql_exec() {
    mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" -p"$DB_PASSWORD" "$DB_NAME" "$@"
}

clean_database() {
    echo -e "${YELLOW}Cleaning database...${NC}"

    # Disable foreign key checks
    mysql_exec -e "SET FOREIGN_KEY_CHECKS = 0;"

    # Get all tables and truncate them
    local tables
    tables=$(mysql_exec -N -e "SHOW TABLES")

    if [ -n "$tables" ]; then
        for table in $tables; do
            echo "  Truncating $table..."
            mysql_exec -e "TRUNCATE TABLE \`$table\`"
        done
    fi

    # Re-enable foreign key checks
    mysql_exec -e "SET FOREIGN_KEY_CHECKS = 1;"

    echo -e "${GREEN}✓ Database cleaned${NC}"
}

seed_buildings() {
    echo -e "${GREEN}Seeding buildings...${NC}"

    mysql_exec <<EOF
INSERT INTO buildings (name, code, address, description, created_at, updated_at) VALUES
    ('Engineering Building A', 'ENG-A', '123 Campus Drive', 'Main engineering building', NOW(), NOW()),
    ('Science Building', 'SCI-B', '456 Research Ave', 'Science and research facilities', NOW(), NOW()),
    ('Library', 'LIB-C', '789 Academic Blvd', 'Central library building', NOW(), NOW())
ON DUPLICATE KEY UPDATE updated_at = NOW();
EOF

    echo -e "${GREEN}✓ Buildings seeded${NC}"
}

seed_resources() {
    echo -e "${GREEN}Seeding resources...${NC}"

    mysql_exec <<EOF
INSERT INTO resources (name, type, capacity, building_id, floor, room_number, description, created_at, updated_at) VALUES
    ('Lab 101', 'laboratory', 30, 1, 1, '101', 'Computer lab', NOW(), NOW()),
    ('Classroom 201', 'classroom', 50, 1, 2, '201', 'Standard classroom', NOW(), NOW()),
    ('Conference Room A', 'meeting_room', 20, 2, 1, 'A-105', 'Conference room with AV equipment', NOW(), NOW()),
    ('Study Room 1', 'study_room', 8, 3, 2, '201', 'Group study room', NOW(), NOW())
ON DUPLICATE KEY UPDATE updated_at = NOW();
EOF

    echo -e "${GREEN}✓ Resources seeded${NC}"
}

seed_classes() {
    echo -e "${GREEN}Seeding classes...${NC}"

    mysql_exec <<EOF
INSERT INTO classes (code, name, description, credits, professor, created_at, updated_at) VALUES
    ('CS101', 'Introduction to Programming', 'Learn programming fundamentals', 4, 'Dr. Smith', NOW(), NOW()),
    ('CS201', 'Data Structures', 'Advanced data structures and algorithms', 4, 'Dr. Johnson', NOW(), NOW()),
    ('MATH101', 'Calculus I', 'Introduction to calculus', 3, 'Prof. Williams', NOW(), NOW())
ON DUPLICATE KEY UPDATE updated_at = NOW();
EOF

    echo -e "${GREEN}✓ Classes seeded${NC}"
}

seed_lessons() {
    echo -e "${GREEN}Seeding lessons...${NC}"

    mysql_exec <<EOF
INSERT INTO lessons (class_id, resource_id, day_of_week, start_time, end_time, created_at, updated_at) VALUES
    (1, 1, 'Monday', '09:00:00', '10:30:00', NOW(), NOW()),
    (1, 1, 'Wednesday', '09:00:00', '10:30:00', NOW(), NOW()),
    (2, 2, 'Tuesday', '14:00:00', '15:30:00', NOW(), NOW()),
    (2, 2, 'Thursday', '14:00:00', '15:30:00', NOW(), NOW()),
    (3, 2, 'Monday', '11:00:00', '12:30:00', NOW(), NOW())
ON DUPLICATE KEY UPDATE updated_at = NOW();
EOF

    echo -e "${GREEN}✓ Lessons seeded${NC}"
}

main() {
    # Parse arguments
    while [ $# -gt 0 ]; do
        case "$1" in
            --env)
                ENV="$2"
                shift 2
                ;;
            --clean)
                CLEAN_MODE=true
                shift
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

    echo -e "${GREEN}SARC-NG Database Seeding${NC}"
    echo -e "Environment: ${YELLOW}$ENV${NC}"
    echo ""

    check_mysql_client

    # Clean database if requested
    if $CLEAN_MODE; then
        clean_database
        echo ""
    fi

    # Seed data
    seed_buildings
    seed_resources
    seed_classes
    seed_lessons

    echo ""
    echo -e "${GREEN}✓ Database seeding completed${NC}"
}

main "$@"

