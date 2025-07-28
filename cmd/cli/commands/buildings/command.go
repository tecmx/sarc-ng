package buildings

import (
	"encoding/json"
	"fmt"
	"sarc-ng/pkg/rest/client"
	"strconv"

	"github.com/spf13/cobra"
)

// NewCommand creates the buildings command group
func NewCommand(clientFactory func() *client.Client) *cobra.Command {
	buildingsCmd := &cobra.Command{
		Use:   "buildings",
		Short: "Manage buildings",
		Long:  "Create, read, update, and delete buildings in the SARC system.",
	}

	// Add subcommands
	buildingsCmd.AddCommand(newListCommand(clientFactory))
	buildingsCmd.AddCommand(newGetCommand(clientFactory))
	buildingsCmd.AddCommand(newCreateCommand(clientFactory))
	buildingsCmd.AddCommand(newUpdateCommand(clientFactory))
	buildingsCmd.AddCommand(newDeleteCommand(clientFactory))

	return buildingsCmd
}

// List all buildings
func newListCommand(clientFactory func() *client.Client) *cobra.Command {
	var outputFormat string

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all buildings",
		Long:  "Retrieve and display all buildings in the system.",
		RunE: func(cmd *cobra.Command, args []string) error {
			client := clientFactory()
			rawResp, err := client.Buildings().List(1, 100) // Get first 100 buildings
			if err != nil {
				return fmt.Errorf("failed to list buildings: %w", err)
			}

			// Parse the raw JSON response
			var response struct {
				Data []Building `json:"data"`
			}
			if err := json.Unmarshal(rawResp, &response); err != nil {
				return fmt.Errorf("failed to parse response: %w", err)
			}

			if len(response.Data) == 0 {
				fmt.Println("No buildings found.")
				return nil
			}

			return OutputWithFormat(response.Data, OutputFormat(outputFormat))
		},
	}

	cmd.Flags().StringVarP(&outputFormat, "output", "o", "table", "Output format (table, json)")
	return cmd
}

// Get a specific building
func newGetCommand(clientFactory func() *client.Client) *cobra.Command {
	var outputFormat string

	cmd := &cobra.Command{
		Use:   "get <id>",
		Short: "Get a building by ID",
		Long:  "Retrieve and display details for a specific building.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.ParseUint(args[0], 10, 32)
			if err != nil {
				return fmt.Errorf("invalid building ID: %s", args[0])
			}

			client := clientFactory()
			rawResp, err := client.Buildings().Get(uint(id))
			if err != nil {
				return fmt.Errorf("failed to get building: %w", err)
			}

			// Parse the raw JSON response
			var building Building
			if err := json.Unmarshal(rawResp, &building); err != nil {
				return fmt.Errorf("failed to parse response: %w", err)
			}

			return OutputWithFormat([]Building{building}, OutputFormat(outputFormat))
		},
	}

	cmd.Flags().StringVarP(&outputFormat, "output", "o", "table", "Output format (table, json)")
	return cmd
}

// Create a new building
func newCreateCommand(clientFactory func() *client.Client) *cobra.Command {
	var name, code string

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new building",
		Long:  "Create a new building with the specified name and code.",
		RunE: func(cmd *cobra.Command, args []string) error {
			if name == "" {
				return fmt.Errorf("building name is required")
			}
			if code == "" {
				return fmt.Errorf("building code is required")
			}

			client := clientFactory()
			req := BuildingRequest{
				Name: name,
				Code: code,
			}

			rawResp, err := client.Buildings().Create(req)
			if err != nil {
				return fmt.Errorf("failed to create building: %w", err)
			}

			// Parse the raw JSON response
			var building Building
			if err := json.Unmarshal(rawResp, &building); err != nil {
				return fmt.Errorf("failed to parse response: %w", err)
			}

			fmt.Printf("✅ Building created successfully:\n")
			return OutputTable([]Building{building})
		},
	}

	cmd.Flags().StringVarP(&name, "name", "n", "", "Building name (required)")
	cmd.Flags().StringVarP(&code, "code", "c", "", "Building code (required)")
	_ = cmd.MarkFlagRequired("name")
	_ = cmd.MarkFlagRequired("code")

	return cmd
}

// Update an existing building
func newUpdateCommand(clientFactory func() *client.Client) *cobra.Command {
	var name, code string

	cmd := &cobra.Command{
		Use:   "update <id>",
		Short: "Update a building",
		Long:  "Update an existing building's name and/or code.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.ParseUint(args[0], 10, 32)
			if err != nil {
				return fmt.Errorf("invalid building ID: %s", args[0])
			}

			client := clientFactory()

			// Get current building to preserve unchanged fields
			rawCurrentResp, err := client.Buildings().Get(uint(id))
			if err != nil {
				return fmt.Errorf("failed to get current building: %w", err)
			}

			var current Building
			if err := json.Unmarshal(rawCurrentResp, &current); err != nil {
				return fmt.Errorf("failed to parse current building: %w", err)
			}

			// Use current values if not specified
			if name == "" {
				name = current.Name
			}
			if code == "" {
				code = current.Code
			}

			req := BuildingRequest{
				Name: name,
				Code: code,
			}

			rawResp, err := client.Buildings().Update(uint(id), req)
			if err != nil {
				return fmt.Errorf("failed to update building: %w", err)
			}

			// Parse the raw JSON response
			var building Building
			if err := json.Unmarshal(rawResp, &building); err != nil {
				return fmt.Errorf("failed to parse response: %w", err)
			}

			fmt.Printf("✅ Building updated successfully:\n")
			return OutputTable([]Building{building})
		},
	}

	cmd.Flags().StringVarP(&name, "name", "n", "", "Building name")
	cmd.Flags().StringVarP(&code, "code", "c", "", "Building code")

	return cmd
}

// Delete a building
func newDeleteCommand(clientFactory func() *client.Client) *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "delete <id>",
		Short: "Delete a building",
		Long:  "Delete a building by ID. Use --force to skip confirmation.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.ParseUint(args[0], 10, 32)
			if err != nil {
				return fmt.Errorf("invalid building ID: %s", args[0])
			}

			client := clientFactory()

			// Get building info for confirmation
			rawResp, err := client.Buildings().Get(uint(id))
			if err != nil {
				return fmt.Errorf("failed to get building: %w", err)
			}

			var building Building
			if err := json.Unmarshal(rawResp, &building); err != nil {
				return fmt.Errorf("failed to parse building: %w", err)
			}

			if !force {
				fmt.Printf("Are you sure you want to delete building '%s' (ID: %d)? [y/N]: ", building.Name, building.ID)
				var response string
				fmt.Scanln(&response)
				if response != "y" && response != "Y" && response != "yes" && response != "Yes" {
					fmt.Println("Operation cancelled.")
					return nil
				}
			}

			err = client.Buildings().Delete(uint(id))
			if err != nil {
				return fmt.Errorf("failed to delete building: %w", err)
			}

			fmt.Printf("✅ Building '%s' deleted successfully.\n", building.Name)
			return nil
		},
	}

	cmd.Flags().BoolVarP(&force, "force", "f", false, "Force deletion without confirmation")
	return cmd
}
