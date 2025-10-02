# SAM Lambda build configuration
# This Makefile is called by AWS SAM CLI during the build process

# Calculate project root from ARTIFACTS_DIR
# ARTIFACTS_DIR = .../infrastructure/sam/.aws-sam/build/SarcNgFunction
# Project root is 5 directories up
PROJECT_ROOT := $(shell dirname $(shell dirname $(shell dirname $(shell dirname $(shell dirname $(ARTIFACTS_DIR))))))

.PHONY: build-SarcNgFunction

build-SarcNgFunction:
	cd $(PROJECT_ROOT) && GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -buildvcs=false -o $(ARTIFACTS_DIR)/bootstrap ./cmd/lambda

