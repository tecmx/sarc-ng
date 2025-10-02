#!/bin/bash

# Build script for SARC-NG Lambda deployment

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${GREEN}Building SARC-NG Lambda function...${NC}"

# Navigate to the project root
cd "$(dirname "$0")/../.."

# Check if we're in the right directory
if [ ! -f "go.mod" ]; then
    echo -e "${RED}Error: go.mod not found. Please run this script from the project root.${NC}"
    exit 1
fi

# Build the Lambda binary
echo -e "${YELLOW}Building Go binary for Lambda...${NC}"
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o cmd/lambda/bootstrap ./cmd/lambda

# Make sure the binary is executable
chmod +x cmd/lambda/bootstrap

echo -e "${GREEN}Build completed successfully!${NC}"
echo -e "${YELLOW}Binary location: cmd/lambda/bootstrap${NC}"

# Navigate to SAM directory
cd infrastructure/sam

echo -e "${GREEN}Ready for SAM deployment. Run the following commands:${NC}"
echo -e "${YELLOW}  ./sam-docker.sh build${NC}"
echo -e "${YELLOW}  ./sam-docker.sh deploy --guided${NC}"
echo ""
echo -e "${GREEN}Or use native SAM CLI if installed:${NC}"
echo -e "${YELLOW}  sam build${NC}"
echo -e "${YELLOW}  sam deploy --guided${NC}"
