package main

import (
	"fmt"
	"os"
	"sarc-ng/cmd/cli/commands/buildings"
	"sarc-ng/cmd/cli/commands/classes"
	"sarc-ng/cmd/cli/commands/health"
	"sarc-ng/cmd/cli/commands/lessons"
	"sarc-ng/cmd/cli/commands/reservations"
	"sarc-ng/cmd/cli/commands/resources"
	"sarc-ng/pkg/rest/client"

	"github.com/spf13/cobra"
)

// GlobalConfig holds configuration used across all commands
type GlobalConfig struct {
	APIBaseURL string
	Timeout    int
	Verbose    bool
}

// NewRootCommand creates the root CLI command
func NewRootCommand() *cobra.Command {
	config := &GlobalConfig{}

	rootCmd := &cobra.Command{
		Use:   "sarc",
		Short: "SARC CLI - Resource management and scheduling system",
		Long: `SARC CLI is a command-line interface for the SARC (Schedule and Resource Control) system.
Use this CLI to manage buildings, resources, classes, lessons, and reservations.`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// Validate configuration
			if config.APIBaseURL == "" {
				return fmt.Errorf("API base URL is required. Set it with --api-url or SARC_API_URL environment variable")
			}
			return nil
		},
	}

	// Global flags
	rootCmd.PersistentFlags().StringVar(&config.APIBaseURL, "api-url",
		getEnvWithDefault("SARC_API_URL", "http://localhost:8080"),
		"SARC API base URL")
	rootCmd.PersistentFlags().IntVar(&config.Timeout, "timeout", 30,
		"Request timeout in seconds")
	rootCmd.PersistentFlags().BoolVarP(&config.Verbose, "verbose", "v", false,
		"Enable verbose output")

	// Create API client factory
	clientFactory := func() *client.Client {
		return client.NewClient(client.Config{
			BaseURL: config.APIBaseURL,
		})
	}

	// Add subcommands
	rootCmd.AddCommand(health.NewCommand(clientFactory))
	rootCmd.AddCommand(buildings.NewCommand(clientFactory))
	rootCmd.AddCommand(resources.NewCommand(clientFactory))
	rootCmd.AddCommand(reservations.NewCommand(clientFactory))
	rootCmd.AddCommand(classes.NewCommand(clientFactory))
	rootCmd.AddCommand(lessons.NewCommand(clientFactory))

	return rootCmd
}

// getEnvWithDefault gets an environment variable or returns a default value
func getEnvWithDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
