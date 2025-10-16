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

# Environment (dev, staging, prod)
ENV ?= prod
AWS_REGION ?= us-east-1

.PHONY: help build-SarcNgFunction build deploy deploy-cognito deploy-api deploy-all delete delete-cognito delete-api delete-all urls clean validate validate-cognito status ui-upload ui-apply ui-customize

# Default target - show help
help:
	@echo "Available commands:"
	@echo ""
	@echo "Two-Stack Deployment:"
	@echo "  deploy-all       - Deploy both Cognito and API stacks (ENV=prod)"
	@echo "  deploy-cognito   - Deploy Cognito authentication stack"
	@echo "  deploy-api       - Deploy API stack (requires Cognito stack)"
	@echo "  delete-all       - Delete both stacks"
	@echo "  delete-cognito   - Delete Cognito stack"
	@echo "  delete-api       - Delete API stack"
	@echo ""
	@echo "Build & Validate:"
	@echo "  build            - Build the SAM applications"
	@echo "  validate         - Validate both SAM templates"
	@echo "  validate-cognito - Validate Cognito template only"
	@echo ""
	@echo "UI Customization:"
	@echo "  ui-customize     - Apply Cognito Hosted UI customization (ENV=prod)"
	@echo "  ui-upload        - Upload UI assets to S3 (ENV=prod)"
	@echo "  ui-apply         - Apply UI customization from local files"
	@echo ""
	@echo "Utilities:"
	@echo "  status           - Show deployment status for both stacks"
	@echo "  urls             - Show deployed API and Swagger URLs"
	@echo "  clean            - Remove build artifacts"
	@echo "  help             - Show this help message"
	@echo ""
	@echo "Environment:"
	@echo "  Use ENV=dev|staging|prod to specify environment (default: prod)"
	@echo "  Example: make deploy-all ENV=dev"
	@echo "  Example: make ui-customize ENV=dev"

# Build target called by SAM CLI
build-SarcNgFunction:
	cd $(PROJECT_ROOT) && GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -buildvcs=false -o $(ARTIFACTS_DIR)/bootstrap ./cmd/lambda

# Build both SAM applications
build:
	@echo "Building Cognito stack..."
	@sam build --template cognito-stack.yaml
	@echo ""
	@echo "Building API stack..."
	@sam build --template api-stack.yaml

# Deploy Cognito stack
deploy-cognito:
	@echo "=== Deploying Cognito Stack ($(ENV)) ==="
	sam build --template cognito-stack.yaml && \
	sam deploy --template cognito-stack.yaml \
		--stack-name sarc-ng-cognito-$(ENV) \
		--capabilities CAPABILITY_IAM \
		--parameter-overrides \
			Environment=$(ENV) \
			ApplicationName=sarc-ng \
		--no-fail-on-empty-changeset
	@echo "✓ Cognito stack deployed successfully"

# Deploy API stack
deploy-api:
	@echo "=== Deploying API Stack ($(ENV)) ==="
	sam build --template api-stack.yaml && \
	sam package \
		--template-file .aws-sam/build/template.yaml \
		--s3-bucket $$(aws cloudformation describe-stacks --stack-name aws-sam-cli-managed-default --query 'Stacks[0].Outputs[?OutputKey==`SourceBucket`].OutputValue' --output text 2>/dev/null || echo 'aws-sam-cli-managed-default-samclisourcebucket-nbtihuf1byzc') \
		--output-template-file .aws-sam/build/packaged.yaml \
		--region $(AWS_REGION) \
		--force-upload && \
	sam deploy \
		--template-file .aws-sam/build/packaged.yaml \
		--stack-name sarc-ng-$(ENV) \
		--capabilities CAPABILITY_IAM \
		--parameter-overrides Environment=$(ENV) \
		--region $(AWS_REGION) \
		--no-confirm-changeset \
		--no-fail-on-empty-changeset
	@echo "✓ API stack deployed successfully"

# Deploy both stacks in order
deploy-all: deploy-cognito deploy-api
	@echo ""
	@echo "✓ All stacks deployed successfully!"
	@echo "Run 'make urls ENV=$(ENV)' to see your endpoints"

# Legacy deploy command (deploys API only)
deploy: deploy-api

# Delete Cognito stack
delete-cognito:
	sam delete --stack-name sarc-ng-cognito-$(ENV) --no-prompts; \

# Delete API stack
delete-api:
	sam delete --stack-name sarc-ng-$(ENV) --no-prompts

# Delete both stacks
delete-all: delete-api delete-cognito

# Legacy delete command (deletes API only)
delete: delete-api

# Show Swagger documentation and Cognito URLs
urls:
	@echo '=== Deployment URLs ($(ENV)) ==='
	@echo ''
	@echo '📋 Cognito Stack:'
	@COGNITO_DOMAIN=$$(aws cloudformation describe-stacks \
		--stack-name sarc-ng-cognito-$(ENV) \
		--query 'Stacks[0].Outputs[?OutputKey==`CognitoUserPoolDomain`].OutputValue' \
		--output text 2>/dev/null || echo 'Not deployed'); \
	if [ "$$COGNITO_DOMAIN" != "Not deployed" ]; then \
		echo "   Hosted UI: $$COGNITO_DOMAIN"; \
	else \
		echo "   ❌ Cognito stack not deployed"; \
	fi
	@echo ''
	@echo '📋 API Stack:'
	@API_URL=$$(aws cloudformation describe-stacks \
		--stack-name sarc-ng-$(ENV) \
		--query 'Stacks[0].Outputs[?OutputKey==`SarcNgApi`].OutputValue' \
		--output text 2>/dev/null || echo 'Not deployed'); \
	if [ "$$API_URL" != 'Not deployed' ] && [ -n "$$API_URL" ]; then \
		echo "   API Base: $$API_URL"; \
		echo "   Swagger: $${API_URL}swagger/index.html"; \
		echo "   Health: $${API_URL}health"; \
	else \
		echo "   ❌ API stack not deployed"; \
		echo "   💡 Run 'make deploy-all ENV=$(ENV)' to deploy"; \
	fi

# Remove build artifacts
clean:
	rm -rf .aws-sam

# Validate Cognito template
validate-cognito:
	@echo "Validating Cognito template..."
	@sam validate --template cognito-stack.yaml --lint
	@echo "✓ Cognito template is valid"

# Validate both SAM templates
validate: validate-cognito
	@echo ""
	@echo "Validating API template..."
	@sam validate --template api-stack.yaml --lint
	@echo "✓ API template is valid"
	@echo ""
	@echo "✓ All templates are valid"

# Show deployment status for both stacks
status:
	@echo '=== Cognito Stack Status ($(ENV)) ==='
	@aws cloudformation describe-stacks \
		--stack-name sarc-ng-cognito-$(ENV) \
		--query 'Stacks[0].[StackName,StackStatus,LastUpdatedTime]' \
		--output table 2>/dev/null || echo '❌ Cognito stack not deployed'
	@echo ''
	@echo '=== API Stack Status ($(ENV)) ==='
	@aws cloudformation describe-stacks \
		--stack-name sarc-ng-$(ENV) \
		--query 'Stacks[0].[StackName,StackStatus,LastUpdatedTime]' \
		--output table 2>/dev/null || echo '❌ API stack not deployed'

# Upload UI assets to S3
ui-upload:
	@echo "=== Uploading UI Assets to S3 ($(ENV)) ==="
	@cd ui-customization && chmod +x upload-ui-assets.sh && ./upload-ui-assets.sh $(ENV)

# Apply UI customization from local files
ui-apply:
	@echo "=== Applying UI Customization ($(ENV)) ==="
	@cd ui-customization && chmod +x apply-ui-customization.sh && ./apply-ui-customization.sh $(ENV)

# Complete UI customization workflow
ui-customize: ui-apply
	@echo ""
	@echo "✓ UI customization complete!"
	@echo "Visit the Hosted UI URL to see your custom login page."
