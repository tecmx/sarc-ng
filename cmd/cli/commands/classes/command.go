package classes

import (
	"encoding/json"
	"fmt"
	"sarc-ng/pkg/rest/client"
	"strconv"

	"github.com/spf13/cobra"
)

// NewCommand creates the classes command group
func NewCommand(clientFactory func() *client.Client) *cobra.Command {
	classesCmd := &cobra.Command{
		Use:   "classes",
		Short: "Manage classes",
		Long:  "Create, read, update, and delete classes in the SARC system.",
	}

	// Add subcommands
	classesCmd.AddCommand(newListCommand(clientFactory))
	classesCmd.AddCommand(newGetCommand(clientFactory))
	classesCmd.AddCommand(newCreateCommand(clientFactory))
	classesCmd.AddCommand(newUpdateCommand(clientFactory))
	classesCmd.AddCommand(newDeleteCommand(clientFactory))

	return classesCmd
}

// List all classes
func newListCommand(clientFactory func() *client.Client) *cobra.Command {
	var outputFormat string

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all classes",
		Long:  "Retrieve and display all classes in the system.",
		RunE: func(cmd *cobra.Command, args []string) error {
			client := clientFactory()
			rawResp, err := client.Classes().List(1, 100) // Get first 100 classes
			if err != nil {
				return fmt.Errorf("failed to list classes: %w", err)
			}

			// Parse the raw JSON response
			var response struct {
				Data []Class `json:"data"`
			}
			if err := json.Unmarshal(rawResp, &response); err != nil {
				// If it's not paginated, try parsing as direct array
				var classes []Class
				if err := json.Unmarshal(rawResp, &classes); err != nil {
					return fmt.Errorf("failed to parse response: %w", err)
				}
				response.Data = classes
			}

			if len(response.Data) == 0 {
				fmt.Println("No classes found.")
				return nil
			}

			return OutputWithFormat(response.Data, OutputFormat(outputFormat))
		},
	}

	cmd.Flags().StringVarP(&outputFormat, "output", "o", "table", "Output format (table, json)")
	return cmd
}

// Get a specific class
func newGetCommand(clientFactory func() *client.Client) *cobra.Command {
	var outputFormat string

	cmd := &cobra.Command{
		Use:   "get <id>",
		Short: "Get a class by ID",
		Long:  "Retrieve and display details for a specific class.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.ParseUint(args[0], 10, 32)
			if err != nil {
				return fmt.Errorf("invalid class ID: %s", args[0])
			}

			client := clientFactory()
			rawResp, err := client.Classes().Get(uint(id))
			if err != nil {
				return fmt.Errorf("failed to get class: %w", err)
			}

			// Parse the raw JSON response
			var class Class
			if err := json.Unmarshal(rawResp, &class); err != nil {
				return fmt.Errorf("failed to parse response: %w", err)
			}

			return OutputWithFormat([]Class{class}, OutputFormat(outputFormat))
		},
	}

	cmd.Flags().StringVarP(&outputFormat, "output", "o", "table", "Output format (table, json)")
	return cmd
}

// Create a new class
func newCreateCommand(clientFactory func() *client.Client) *cobra.Command {
	var name string
	var capacity int

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new class",
		Long:  "Create a new class with the specified name and capacity.",
		RunE: func(cmd *cobra.Command, args []string) error {
			if name == "" {
				return fmt.Errorf("class name is required")
			}
			if capacity <= 0 {
				return fmt.Errorf("class capacity must be greater than 0")
			}

			client := clientFactory()
			req := ClassRequest{
				Name:     name,
				Capacity: capacity,
			}

			rawResp, err := client.Classes().Create(req)
			if err != nil {
				return fmt.Errorf("failed to create class: %w", err)
			}

			// Parse the raw JSON response
			var class Class
			if err := json.Unmarshal(rawResp, &class); err != nil {
				return fmt.Errorf("failed to parse response: %w", err)
			}

			fmt.Printf("✅ Class created successfully:\n")
			return OutputTable([]Class{class})
		},
	}

	cmd.Flags().StringVarP(&name, "name", "n", "", "Class name (required)")
	cmd.Flags().IntVarP(&capacity, "capacity", "c", 0, "Class capacity (required)")
	_ = cmd.MarkFlagRequired("name")
	_ = cmd.MarkFlagRequired("capacity")

	return cmd
}

// Update an existing class
func newUpdateCommand(clientFactory func() *client.Client) *cobra.Command {
	var name string
	var capacity int

	cmd := &cobra.Command{
		Use:   "update <id>",
		Short: "Update a class",
		Long:  "Update an existing class's name and/or capacity.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.ParseUint(args[0], 10, 32)
			if err != nil {
				return fmt.Errorf("invalid class ID: %s", args[0])
			}

			client := clientFactory()

			// Get current class to preserve unchanged fields
			rawCurrentResp, err := client.Classes().Get(uint(id))
			if err != nil {
				return fmt.Errorf("failed to get current class: %w", err)
			}

			var current Class
			if err := json.Unmarshal(rawCurrentResp, &current); err != nil {
				return fmt.Errorf("failed to parse current class: %w", err)
			}

			// Use current values if not specified
			if name == "" {
				name = current.Name
			}
			if capacity == 0 {
				capacity = current.Capacity
			}

			req := ClassRequest{
				Name:     name,
				Capacity: capacity,
			}

			rawResp, err := client.Classes().Update(uint(id), req)
			if err != nil {
				return fmt.Errorf("failed to update class: %w", err)
			}

			// Parse the raw JSON response
			var class Class
			if err := json.Unmarshal(rawResp, &class); err != nil {
				return fmt.Errorf("failed to parse response: %w", err)
			}

			fmt.Printf("✅ Class updated successfully:\n")
			return OutputTable([]Class{class})
		},
	}

	cmd.Flags().StringVarP(&name, "name", "n", "", "Class name")
	cmd.Flags().IntVarP(&capacity, "capacity", "c", 0, "Class capacity")

	return cmd
}

// Delete a class
func newDeleteCommand(clientFactory func() *client.Client) *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "delete <id>",
		Short: "Delete a class",
		Long:  "Delete a class by ID. Use --force to skip confirmation.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.ParseUint(args[0], 10, 32)
			if err != nil {
				return fmt.Errorf("invalid class ID: %s", args[0])
			}

			client := clientFactory()

			// Get class info for confirmation
			rawResp, err := client.Classes().Get(uint(id))
			if err != nil {
				return fmt.Errorf("failed to get class: %w", err)
			}

			var class Class
			if err := json.Unmarshal(rawResp, &class); err != nil {
				return fmt.Errorf("failed to parse class: %w", err)
			}

			if !force {
				fmt.Printf("Are you sure you want to delete class '%s' (ID: %d)? [y/N]: ", class.Name, class.ID)
				var response string
				fmt.Scanln(&response)
				if response != "y" && response != "Y" && response != "yes" && response != "Yes" {
					fmt.Println("Operation cancelled.")
					return nil
				}
			}

			err = client.Classes().Delete(uint(id))
			if err != nil {
				return fmt.Errorf("failed to delete class: %w", err)
			}

			fmt.Printf("✅ Class '%s' deleted successfully.\n", class.Name)
			return nil
		},
	}

	cmd.Flags().BoolVarP(&force, "force", "f", false, "Force deletion without confirmation")
	return cmd
}
