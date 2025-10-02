#!/bin/bash
# Semantic version management script for SARC-NG
# Handles version bumping and git tagging

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Get project root directory
PROJECT_ROOT="$(cd "$(dirname "$0")/.." && pwd)"
VERSION_FILE="${PROJECT_ROOT}/VERSION"

DRY_RUN=false
CREATE_TAG=true
PUSH_TAG=false

show_help() {
    echo -e "${GREEN}Semantic Version Management Tool${NC}"
    echo ""
    echo "Usage: $0 <bump_type> [options]"
    echo ""
    echo "Bump Types:"
    echo "  major       Increment major version (1.0.0 -> 2.0.0)"
    echo "  minor       Increment minor version (1.0.0 -> 1.1.0)"
    echo "  patch       Increment patch version (1.0.0 -> 1.0.1)"
    echo ""
    echo "Options:"
    echo "  --no-tag    Don't create git tag"
    echo "  --push      Push tag to remote"
    echo "  --dry-run   Show what would change without executing"
    echo "  -h, --help  Show this help message"
    echo ""
    echo "Examples:"
    echo "  $0 patch                # 1.0.0 -> 1.0.1"
    echo "  $0 minor --push         # 1.0.0 -> 1.1.0 and push tag"
    echo "  $0 major --dry-run      # Show what would happen"
}

get_current_version() {
    if [ -f "$VERSION_FILE" ]; then
        cat "$VERSION_FILE"
    else
        # Try to get from git tags
        git describe --tags --abbrev=0 2>/dev/null | sed 's/^v//' || echo "0.0.0"
    fi
}

parse_version() {
    local version="$1"
    local major minor patch

    # Remove 'v' prefix if present
    version="${version#v}"

    # Split version into components
    IFS='.' read -r major minor patch <<< "$version"

    # Remove any non-numeric suffix (e.g., -beta, -rc1)
    major="${major%%[^0-9]*}"
    minor="${minor%%[^0-9]*}"
    patch="${patch%%[^0-9]*}"

    echo "$major $minor $patch"
}

bump_version() {
    local bump_type="$1"
    local current_version="$2"

    read -r major minor patch <<< "$(parse_version "$current_version")"

    case "$bump_type" in
        major)
            major=$((major + 1))
            minor=0
            patch=0
            ;;
        minor)
            minor=$((minor + 1))
            patch=0
            ;;
        patch)
            patch=$((patch + 1))
            ;;
        *)
            echo -e "${RED}Error: Invalid bump type '$bump_type'${NC}"
            exit 2
            ;;
    esac

    echo "${major}.${minor}.${patch}"
}

update_version_file() {
    local new_version="$1"

    if $DRY_RUN; then
        echo -e "${YELLOW}[DRY RUN] Would write $new_version to $VERSION_FILE${NC}"
    else
        echo "$new_version" > "$VERSION_FILE"
        echo -e "${GREEN}✓ Updated VERSION file${NC}"
    fi
}

create_git_tag() {
    local version="$1"
    local tag="v${version}"

    if ! command -v git &> /dev/null; then
        echo -e "${YELLOW}⚠ Git not available, skipping tag creation${NC}"
        return
    fi

    if ! git rev-parse --git-dir > /dev/null 2>&1; then
        echo -e "${YELLOW}⚠ Not a git repository, skipping tag creation${NC}"
        return
    fi

    if $DRY_RUN; then
        echo -e "${YELLOW}[DRY RUN] Would create git tag: $tag${NC}"
        if $PUSH_TAG; then
            echo -e "${YELLOW}[DRY RUN] Would push tag to remote${NC}"
        fi
    else
        if git tag -l | grep -q "^${tag}$"; then
            echo -e "${YELLOW}⚠ Tag $tag already exists${NC}"
            return
        fi

        git tag -a "$tag" -m "Release version $version"
        echo -e "${GREEN}✓ Created git tag: $tag${NC}"

        if $PUSH_TAG; then
            git push origin "$tag"
            echo -e "${GREEN}✓ Pushed tag to remote${NC}"
        fi
    fi
}

main() {
    local bump_type=""

    # Parse arguments
    while [ $# -gt 0 ]; do
        case "$1" in
            major|minor|patch)
                bump_type="$1"
                shift
                ;;
            --no-tag)
                CREATE_TAG=false
                shift
                ;;
            --push)
                PUSH_TAG=true
                shift
                ;;
            --dry-run)
                DRY_RUN=true
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

    if [ -z "$bump_type" ]; then
        echo -e "${RED}Error: Bump type required${NC}"
        show_help
        exit 2
    fi

    echo -e "${GREEN}Version Bump Tool${NC}"
    echo ""

    local current_version
    current_version=$(get_current_version)
    echo -e "Current version: ${YELLOW}$current_version${NC}"

    local new_version
    new_version=$(bump_version "$bump_type" "$current_version")
    echo -e "New version:     ${GREEN}$new_version${NC}"
    echo ""

    if $DRY_RUN; then
        echo -e "${YELLOW}DRY RUN MODE - No changes will be made${NC}"
        echo ""
    fi

    # Update version file
    update_version_file "$new_version"

    # Create git tag
    if $CREATE_TAG; then
        create_git_tag "$new_version"
    fi

    echo ""
    if $DRY_RUN; then
        echo -e "${YELLOW}✓ Dry run completed${NC}"
    else
        echo -e "${GREEN}✓ Version bump completed: $current_version -> $new_version${NC}"

        if $CREATE_TAG && ! $PUSH_TAG; then
            echo ""
            echo -e "${YELLOW}Remember to push the tag with: git push origin v$new_version${NC}"
        fi
    fi
}

main "$@"

