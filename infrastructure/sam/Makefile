# SAM Lambda build configuration
# This Makefile is called by AWS SAM CLI during the build process

# Calculate project root from ARTIFACTS_DIR
# ARTIFACTS_DIR = .../infrastructure/sam/.aws-sam/build/SarcNgFunction
# Project root is 5 directories up
ifdef ARTIFACTS_DIR
PROJECT_ROOT := $(shell dirname $(shell dirname $(shell dirname $(shell dirname $(shell dirname $(ARTIFACTS_DIR))))))
else
PROJECT_ROOT := $(shell cd ../.. && pwd)
endif

.PHONY: help build-SarcNgFunction build deploy local delete urls

# Default target - show help
help:
	@echo "Available commands:"
	@echo "  build    - Build the SAM application"
	@echo "  deploy   - Deploy the application (sources aws.env)"
	@echo "  local    - Start local API Gateway on port 3000"
	@echo "  delete   - Delete the CloudFormation stack"
	@echo "  urls     - Show Swagger documentation URL"
	@echo "  help     - Show this help message"

# Build target called by SAM CLI
build-SarcNgFunction:
	cd $(PROJECT_ROOT) && GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -buildvcs=false -o $(ARTIFACTS_DIR)/bootstrap ./cmd/lambda

# Build the SAM application
build:
	sam build

# Deploy the application with AWS credentials
deploy:
	@bash -c "source aws.env && sam build && sam deploy"

# Start local API Gateway
local:
	sam local start-api --host 0.0.0.0 --port 3000

# Delete the CloudFormation stack
delete:
	@bash -c "source aws.env && sam delete"

# Show Swagger documentation URL
urls:
	@bash -c "source aws.env && echo 'Getting deployed stack outputs...' && \
	API_URL=\$$(sam list stack-outputs --stack-name sarc-ng-prod --output json 2>/dev/null | grep -A1 '\"OutputKey\": \"SarcNgApi\"' | grep '\"OutputValue\"' | cut -d'\"' -f4 || echo 'Stack not found'); \
	if [ \"\$$API_URL\" != 'Stack not found' ] && [ -n \"\$$API_URL\" ]; then \
		echo '📋 Swagger Documentation: '\$${API_URL}'swagger/index.html'; \
		echo '🏥 Health Check: '\$${API_URL}'health'; \
	else \
		echo '❌ Stack sarc-ng-prod not found or not deployed yet.'; \
		echo '💡 Run make deploy first to deploy the application.'; \
	fi"

