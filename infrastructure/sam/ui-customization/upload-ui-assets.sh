#!/bin/bash
# Upload UI Customization Assets to S3 for Cognito Hosted UI
#
# Usage: ./upload-ui-assets.sh [environment] [bucket-name]
# Example: ./upload-ui-assets.sh dev sarc-ng-ui-assets-dev

set -e

# Color codes for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Default values
ENV="${1:-dev}"
BUCKET_NAME="${2:-sarc-ng-ui-assets-${ENV}}"
REGION="${AWS_REGION:-us-east-1}"

echo -e "${BLUE}=== SARC-NG Cognito UI Assets Upload ===${NC}"
echo -e "${BLUE}Environment: ${ENV}${NC}"
echo -e "${BLUE}Bucket: ${BUCKET_NAME}${NC}"
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
    echo -e "${YELLOW}Warning: logo.svg not found. Skipping logo upload.${NC}"
    SKIP_LOGO=true
fi

# Create bucket if it doesn't exist
echo -e "${YELLOW}Checking if bucket exists...${NC}"
if aws s3 ls "s3://${BUCKET_NAME}" 2>&1 | grep -q 'NoSuchBucket'; then
    echo -e "${YELLOW}Bucket does not exist. Creating ${BUCKET_NAME}...${NC}"

    if [ "$REGION" == "us-east-1" ]; then
        aws s3api create-bucket \
            --bucket "${BUCKET_NAME}" \
            --region "${REGION}"
    else
        aws s3api create-bucket \
            --bucket "${BUCKET_NAME}" \
            --region "${REGION}" \
            --create-bucket-configuration LocationConstraint="${REGION}"
    fi

    # Enable versioning
    aws s3api put-bucket-versioning \
        --bucket "${BUCKET_NAME}" \
        --versioning-configuration Status=Enabled

    # Add encryption
    aws s3api put-bucket-encryption \
        --bucket "${BUCKET_NAME}" \
        --server-side-encryption-configuration '{
            "Rules": [{
                "ApplyServerSideEncryptionByDefault": {
                    "SSEAlgorithm": "AES256"
                }
            }]
        }'

    # Block public access (files will be accessed via pre-signed URLs)
    aws s3api put-public-access-block \
        --bucket "${BUCKET_NAME}" \
        --public-access-block-configuration \
            "BlockPublicAcls=true,IgnorePublicAcls=true,BlockPublicPolicy=true,RestrictPublicBuckets=true"

    echo -e "${GREEN}✓ Bucket created successfully${NC}"
else
    echo -e "${GREEN}✓ Bucket exists${NC}"
fi

# Upload CSS file
echo -e "${YELLOW}Uploading custom.css...${NC}"
aws s3 cp custom.css "s3://${BUCKET_NAME}/cognito-ui/custom.css" \
    --content-type "text/css" \
    --metadata "environment=${ENV},asset-type=cognito-ui-css" \
    --region "${REGION}"

echo -e "${GREEN}✓ CSS uploaded successfully${NC}"

# Upload logo if it exists
if [ "$SKIP_LOGO" != "true" ]; then
    echo -e "${YELLOW}Uploading logo.svg...${NC}"
    aws s3 cp logo.svg "s3://${BUCKET_NAME}/cognito-ui/logo.svg" \
        --content-type "image/svg+xml" \
        --metadata "environment=${ENV},asset-type=cognito-ui-logo" \
        --region "${REGION}"

    echo -e "${GREEN}✓ Logo uploaded successfully${NC}"
fi

# Generate pre-signed URLs (valid for 7 days)
echo ""
echo -e "${BLUE}=== Generating Pre-signed URLs ===${NC}"

CSS_URL=$(aws s3 presign "s3://${BUCKET_NAME}/cognito-ui/custom.css" \
    --expires-in 604800 \
    --region "${REGION}")

echo -e "${GREEN}CSS URL:${NC}"
echo "${CSS_URL}"
echo ""

if [ "$SKIP_LOGO" != "true" ]; then
    LOGO_URL=$(aws s3 presign "s3://${BUCKET_NAME}/cognito-ui/logo.svg" \
        --expires-in 604800 \
        --region "${REGION}")

    echo -e "${GREEN}Logo URL:${NC}"
    echo "${LOGO_URL}"
    echo ""
fi

# Save URLs to file for reference
OUTPUT_FILE="ui-assets-urls-${ENV}.txt"
cat > "${OUTPUT_FILE}" << EOF
# SARC-NG Cognito UI Assets URLs
# Environment: ${ENV}
# Generated: $(date)
# Valid for: 7 days

CSS_URL=${CSS_URL}
EOF

if [ "$SKIP_LOGO" != "true" ]; then
    echo "LOGO_URL=${LOGO_URL}" >> "${OUTPUT_FILE}"
fi

echo -e "${GREEN}✓ URLs saved to ${OUTPUT_FILE}${NC}"
echo ""

# Instructions
echo -e "${BLUE}=== Next Steps ===${NC}"
echo "1. Use these URLs to customize your Cognito User Pool Client"
echo "2. In AWS Console:"
echo "   - Go to Cognito > User Pools > sarc-ng-users"
echo "   - Select 'App integration' tab"
echo "   - Click on your app client"
echo "   - Scroll to 'Hosted UI customization'"
echo "   - Upload logo or paste Logo URL"
echo "   - Paste CSS URL or upload CSS file"
echo ""
echo "3. Or use AWS CLI:"
echo "   aws cognito-idp set-ui-customization \\"
echo "     --user-pool-id <USER_POOL_ID> \\"
echo "     --client-id <CLIENT_ID> \\"
echo "     --css \"\$(cat custom.css)\" \\"
echo "     --image-file fileb://logo.svg"
echo ""
echo -e "${GREEN}✓ Upload complete!${NC}"
