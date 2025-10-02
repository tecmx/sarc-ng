.PHONY: build-SarcNgFunction
build-SarcNgFunction:
	cd /workspace && GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o $(ARTIFACTS_DIR)/bootstrap ./cmd/lambda