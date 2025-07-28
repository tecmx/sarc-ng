package main

import (
	"os"
)

// main is the entry point for the SARC CLI application
func main() {
	// Create the root command
	rootCmd := NewRootCommand()

	// Execute the command and handle any errors
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
