# SAM Lambda build configuration
# This Makefile is called by AWS SAM CLI during the build process

# Calculate project root
# When SAM CLI calls this: ARTIFACTS_DIR = .../infrastructure/sam/.aws-sam/build/SarcNgFunction
# Go up 5 levels to reach project root
ifdef ARTIFACTS_DIR
PROJECT_ROOT := $(realpath $(ARTIFACTS_DIR)/../../../../..)
else
PROJECT_ROOT := $(realpath ../..)
endif

.PHONY: help build-SarcNgFunction build deploy delete urls clean validate status

# Default target - show help
help:
	@echo "Available commands:"
	@echo ""
	@echo "Build & Deploy:"
	@echo "  build        - Build the SAM application"
	@echo "  validate     - Validate SAM template"
	@echo "  deploy       - Deploy to AWS"
	@echo "  delete       - Delete CloudFormation stack"
	@echo ""
	@echo "Utilities:"
	@echo "  status       - Show deployment status and URLs"
	@echo "  urls         - Show deployed Swagger documentation URL"
	@echo "  clean        - Remove build artifacts"
	@echo "  help         - Show this help message"

# Build target called by SAM CLI
build-SarcNgFunction:
	cd $(PROJECT_ROOT) && GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -buildvcs=false -o $(ARTIFACTS_DIR)/bootstrap ./cmd/lambda

# Build the SAM application
build:
	sam build

# Deploy the application
deploy:
	sam build && sam deploy

# Delete the CloudFormation stack
delete:
	sam delete

# Show Swagger documentation URL
urls:
	@echo 'Getting deployed stack outputs...'
	@API_URL=$$(sam list stack-outputs --stack-name sarc-ng-prod --output json 2>/dev/null | grep -A1 '"OutputKey": "SarcNgApi"' | grep '"OutputValue"' | cut -d'"' -f4 || echo 'Stack not found'); \
	if [ "$$API_URL" != 'Stack not found' ] && [ -n "$$API_URL" ]; then \
		echo '📋 Swagger Documentation: '$$API_URL'swagger/index.html'; \
		echo '🏥 Health Check: '$$API_URL'health'; \
	else \
		echo '❌ Stack sarc-ng-prod not found or not deployed yet.'; \
		echo '💡 Run make deploy first to deploy the application.'; \
	fi

# Remove build artifacts
clean:
	rm -rf .aws-sam

# Validate SAM template
validate:
	@echo "Validating SAM template..."
	@sam validate --lint
	@echo "✓ Template is valid"

# Show deployment status
status:
	@echo '=== Deployment Status ==='
	@sam list stack-outputs --stack-name sarc-ng-prod 2>/dev/null || echo '❌ Stack not deployed'
