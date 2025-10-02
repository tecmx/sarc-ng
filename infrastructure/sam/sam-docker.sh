#!/bin/bash

# SAM CLI Docker wrapper script

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# SAM CLI Docker image
SAM_IMAGE="public.ecr.aws/sam/build-provided.al2:latest"

# Get the project root directory
PROJECT_ROOT="$(cd "$(dirname "$0")/../.." && pwd)"

echo -e "${GREEN}Using AWS SAM CLI Docker image: ${SAM_IMAGE}${NC}"

# Function to run SAM commands in Docker
run_sam() {
    docker run --rm -it \
        -v "${PROJECT_ROOT}:/workspace" \
        -w "/workspace/infrastructure/sam" \
        "${SAM_IMAGE}" \
        sam "$@"
}

# Check if Docker is available
if ! command -v docker &> /dev/null; then
    echo -e "${RED}Error: Docker is not installed or not in PATH${NC}"
    exit 1
fi

# Check if Docker daemon is running
if ! docker info &> /dev/null; then
    echo -e "${RED}Error: Docker daemon is not running${NC}"
    exit 1
fi

# If no arguments provided, show help
if [ $# -eq 0 ]; then
    echo -e "${YELLOW}AWS SAM CLI Docker Wrapper${NC}"
    echo ""
    echo "Usage: $0 <sam-command> [options]"
    echo ""
    echo "Examples:"
    echo "  $0 build"
    echo "  $0 deploy --guided"
    echo "  $0 local start-api"
    echo "  $0 validate"
    echo ""
    echo "Available SAM commands:"
    run_sam --help
    exit 0
fi

# Run the SAM command
echo -e "${YELLOW}Running: sam $*${NC}"
run_sam "$@"