#!/bin/bash
# Apply UI Customization to Cognito User Pool
#
# Usage: ./apply-ui-customization.sh [environment]
# Example: ./apply-ui-customization.sh dev

set -e

# Color codes for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Default values
ENV="${1:-dev}"
STACK_NAME="sarc-ng-cognito-${ENV}"
REGION="${AWS_REGION:-us-east-1}"

echo -e "${BLUE}=== SARC-NG Cognito UI Customization ===${NC}"
echo -e "${BLUE}Environment: ${ENV}${NC}"
echo -e "${BLUE}Stack: ${STACK_NAME}${NC}"
echo -e "${BLUE}Region: ${REGION}${NC}"
echo ""

# Check if AWS CLI is installed
if ! command -v aws &> /dev/null; then
    echo -e "${RED}Error: AWS CLI is not installed${NC}"
    echo "Please install AWS CLI: https://aws.amazon.com/cli/"
    exit 1
fi

# Check if files exist
if [ ! -f "custom.css" ]; then
    echo -e "${RED}Error: custom.css not found${NC}"
    exit 1
fi

if [ ! -f "logo.svg" ]; then
    echo -e "${YELLOW}Warning: logo.svg not found. Will apply CSS only.${NC}"
    SKIP_LOGO=true
fi

# Get User Pool ID from CloudFormation stack
echo -e "${YELLOW}Getting User Pool ID from stack...${NC}"
USER_POOL_ID=$(aws cloudformation describe-stacks \
    --stack-name "${STACK_NAME}" \
    --region "${REGION}" \
    --query 'Stacks[0].Outputs[?OutputKey==`CognitoUserPoolId`].OutputValue' \
    --output text 2>/dev/null)

if [ -z "$USER_POOL_ID" ] || [ "$USER_POOL_ID" == "None" ]; then
    echo -e "${RED}Error: Could not find User Pool ID${NC}"
    echo "Make sure the Cognito stack is deployed: make deploy-cognito ENV=${ENV}"
    exit 1
fi

echo -e "${GREEN}✓ User Pool ID: ${USER_POOL_ID}${NC}"

# Get Client ID from CloudFormation stack
echo -e "${YELLOW}Getting Client ID from stack...${NC}"
CLIENT_ID=$(aws cloudformation describe-stacks \
    --stack-name "${STACK_NAME}" \
    --region "${REGION}" \
    --query 'Stacks[0].Outputs[?OutputKey==`CognitoUserPoolClientId`].OutputValue' \
    --output text 2>/dev/null)

if [ -z "$CLIENT_ID" ] || [ "$CLIENT_ID" == "None" ]; then
    echo -e "${RED}Error: Could not find Client ID${NC}"
    exit 1
fi

echo -e "${GREEN}✓ Client ID: ${CLIENT_ID}${NC}"
echo ""

# Check CSS file size (must be < 100KB)
CSS_SIZE=$(wc -c < custom.css)
CSS_SIZE_KB=$((CSS_SIZE / 1024))

if [ $CSS_SIZE -gt 102400 ]; then
    echo -e "${RED}Error: CSS file is too large (${CSS_SIZE_KB}KB)${NC}"
    echo "Maximum size is 100KB"
    exit 1
fi

echo -e "${GREEN}✓ CSS file size OK (${CSS_SIZE_KB}KB)${NC}"

# Check logo file size if exists (must be < 2MB)
if [ "$SKIP_LOGO" != "true" ]; then
    LOGO_SIZE=$(wc -c < logo.svg)
    LOGO_SIZE_KB=$((LOGO_SIZE / 1024))

    if [ $LOGO_SIZE -gt 2097152 ]; then
        echo -e "${RED}Error: Logo file is too large (${LOGO_SIZE_KB}KB)${NC}"
        echo "Maximum size is 2048KB (2MB)"
        exit 1
    fi

    echo -e "${GREEN}✓ Logo file size OK (${LOGO_SIZE_KB}KB)${NC}"
fi

echo ""

# Read CSS content
CSS_CONTENT=$(cat custom.css)

# Apply UI customization
echo -e "${YELLOW}Applying UI customization to Cognito User Pool...${NC}"

if [ "$SKIP_LOGO" == "true" ]; then
    # CSS only
    aws cognito-idp set-ui-customization \
        --user-pool-id "${USER_POOL_ID}" \
        --client-id "${CLIENT_ID}" \
        --css "${CSS_CONTENT}" \
        --region "${REGION}"
else
    # CSS + Logo
    aws cognito-idp set-ui-customization \
        --user-pool-id "${USER_POOL_ID}" \
        --client-id "${CLIENT_ID}" \
        --css "${CSS_CONTENT}" \
        --image-file fileb://logo.svg \
        --region "${REGION}"
fi

echo -e "${GREEN}✓ UI customization applied successfully${NC}"
echo ""

# Get Hosted UI URL
echo -e "${BLUE}=== Testing URL ===${NC}"
DOMAIN=$(aws cloudformation describe-stacks \
    --stack-name "${STACK_NAME}" \
    --region "${REGION}" \
    --query 'Stacks[0].Outputs[?OutputKey==`CognitoUserPoolDomain`].OutputValue' \
    --output text 2>/dev/null)

if [ -n "$DOMAIN" ] && [ "$DOMAIN" != "None" ]; then
    echo -e "${GREEN}Hosted UI URL:${NC}"
    echo "${DOMAIN}/login?client_id=${CLIENT_ID}&response_type=code&redirect_uri=http://localhost:3000/callback"
    echo ""
    echo -e "${YELLOW}Note: Changes may take 2-5 minutes to propagate to the CDN${NC}"
    echo -e "${YELLOW}Clear your browser cache or use incognito mode to see changes${NC}"
else
    echo -e "${YELLOW}Warning: Could not retrieve Hosted UI URL${NC}"
fi

echo ""
echo -e "${GREEN}✓ Done!${NC}"
echo ""
echo -e "${BLUE}Next steps:${NC}"
echo "1. Wait 2-5 minutes for CDN propagation"
echo "2. Open the Hosted UI URL in your browser"
echo "3. Clear browser cache (Ctrl+Shift+R) if you don't see changes"
echo ""
echo "To update the UI, modify custom.css or logo.svg and run this script again."
