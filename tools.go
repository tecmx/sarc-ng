//go:build tools

// Package tools imports tool dependencies
// This file ensures they're tracked in go.mod without being included in the main build
package tools

import (
	// Import tools that can be managed by go modules
	_ "github.com/air-verse/air"                            // Hot-reload for development
	_ "github.com/go-delve/delve/cmd/dlv"                   // Debugger for Go
	_ "github.com/golangci/golangci-lint/cmd/golangci-lint" // Linting
	_ "github.com/google/wire/cmd/wire"                     // Dependency injection
	_ "github.com/swaggo/swag/cmd/swag"                     // API documentation generator
)

// This file tracks tool dependencies using go modules.
//
// To install these tools, run:
//   go install github.com/air-verse/air@latest
//   go install github.com/go-delve/delve/cmd/dlv@latest
//   go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
//   go install github.com/google/wire/cmd/wire@latest
//   go install github.com/swaggo/swag/cmd/swag@latest
//   go install golang.org/x/tools/gopls@latest
//
// Or install all at once:
//   go list -f '{{range .Imports}}{{.}} {{end}}' tools.go | xargs -n1 go install
