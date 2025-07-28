package resources

import (
	"encoding/json"
	"fmt"
	"sarc-ng/pkg/rest/client"
	"strconv"

	"github.com/spf13/cobra"
)

// NewCommand creates the resources command group
func NewCommand(clientFactory func() *client.Client) *cobra.Command {
	resourcesCmd := &cobra.Command{
		Use:   "resources",
		Short: "Manage resources",
		Long:  "Create, read, update, and delete resources in the SARC system.",
	}

	// Add subcommands
	resourcesCmd.AddCommand(newListCommand(clientFactory))
	resourcesCmd.AddCommand(newGetCommand(clientFactory))
	resourcesCmd.AddCommand(newCreateCommand(clientFactory))
	resourcesCmd.AddCommand(newUpdateCommand(clientFactory))
	resourcesCmd.AddCommand(newDeleteCommand(clientFactory))

	return resourcesCmd
}

// List all resources
func newListCommand(clientFactory func() *client.Client) *cobra.Command {
	var outputFormat string

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all resources",
		Long:  "Retrieve and display all resources in the system.",
		RunE: func(cmd *cobra.Command, args []string) error {
			client := clientFactory()
			data, err := client.Resources().List(1, 100)
			if err != nil {
				return fmt.Errorf("failed to list resources: %w", err)
			}

			var resources []Resource
			if err := json.Unmarshal(data, &resources); err != nil {
				return fmt.Errorf("failed to parse resources: %w", err)
			}

			if len(resources) == 0 {
				fmt.Println("No resources found.")
				return nil
			}

			return OutputWithFormat(resources, OutputFormat(outputFormat))
		},
	}

	cmd.Flags().StringVarP(&outputFormat, "output", "o", "table", "Output format (table, json)")
	return cmd
}

// Get a specific resource
func newGetCommand(clientFactory func() *client.Client) *cobra.Command {
	var outputFormat string

	cmd := &cobra.Command{
		Use:   "get <id>",
		Short: "Get a resource by ID",
		Long:  "Retrieve and display details for a specific resource.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.ParseUint(args[0], 10, 32)
			if err != nil {
				return fmt.Errorf("invalid resource ID: %s", args[0])
			}

			client := clientFactory()
			data, err := client.Resources().Get(uint(id))
			if err != nil {
				return fmt.Errorf("failed to get resource: %w", err)
			}

			var resource Resource
			if err := json.Unmarshal(data, &resource); err != nil {
				return fmt.Errorf("failed to parse resource: %w", err)
			}

			return OutputWithFormat([]Resource{resource}, OutputFormat(outputFormat))
		},
	}

	cmd.Flags().StringVarP(&outputFormat, "output", "o", "table", "Output format (table, json)")
	return cmd
}

// Create a new resource
func newCreateCommand(clientFactory func() *client.Client) *cobra.Command {
	var name, resourceType string

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new resource",
		Long:  "Create a new resource with the specified name and type.",
		RunE: func(cmd *cobra.Command, args []string) error {
			if name == "" {
				return fmt.Errorf("resource name is required")
			}
			if resourceType == "" {
				return fmt.Errorf("resource type is required")
			}

			client := clientFactory()
			req := ResourceRequest{
				Name: name,
				Type: resourceType,
			}

			data, err := client.Resources().Create(req)
			if err != nil {
				return fmt.Errorf("failed to create resource: %w", err)
			}

			var resource Resource
			if err := json.Unmarshal(data, &resource); err != nil {
				return fmt.Errorf("failed to parse created resource: %w", err)
			}

			fmt.Printf("✅ Resource created successfully:\n")
			return OutputTable([]Resource{resource})
		},
	}

	cmd.Flags().StringVarP(&name, "name", "n", "", "Resource name (required)")
	cmd.Flags().StringVarP(&resourceType, "type", "t", "", "Resource type (required)")
	_ = cmd.MarkFlagRequired("name")
	_ = cmd.MarkFlagRequired("type")

	return cmd
}

// Update an existing resource
func newUpdateCommand(clientFactory func() *client.Client) *cobra.Command {
	var name, resourceType string

	cmd := &cobra.Command{
		Use:   "update <id>",
		Short: "Update a resource",
		Long:  "Update an existing resource's name and/or type.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.ParseUint(args[0], 10, 32)
			if err != nil {
				return fmt.Errorf("invalid resource ID: %s", args[0])
			}

			client := clientFactory()

			// Get current resource to preserve unchanged fields
			data, err := client.Resources().Get(uint(id))
			if err != nil {
				return fmt.Errorf("failed to get current resource: %w", err)
			}

			var current Resource
			if err := json.Unmarshal(data, &current); err != nil {
				return fmt.Errorf("failed to parse current resource: %w", err)
			}

			// Use current values if not specified
			if name == "" {
				name = current.Name
			}
			if resourceType == "" {
				resourceType = current.Type
			}

			req := ResourceRequest{
				Name: name,
				Type: resourceType,
			}

			updateData, err := client.Resources().Update(uint(id), req)
			if err != nil {
				return fmt.Errorf("failed to update resource: %w", err)
			}

			var resource Resource
			if err := json.Unmarshal(updateData, &resource); err != nil {
				return fmt.Errorf("failed to parse updated resource: %w", err)
			}

			fmt.Printf("✅ Resource updated successfully:\n")
			return OutputTable([]Resource{resource})
		},
	}

	cmd.Flags().StringVarP(&name, "name", "n", "", "Resource name")
	cmd.Flags().StringVarP(&resourceType, "type", "t", "", "Resource type")

	return cmd
}

// Delete a resource
func newDeleteCommand(clientFactory func() *client.Client) *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "delete <id>",
		Short: "Delete a resource",
		Long:  "Delete a resource by ID. Use --force to skip confirmation.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.ParseUint(args[0], 10, 32)
			if err != nil {
				return fmt.Errorf("invalid resource ID: %s", args[0])
			}

			client := clientFactory()

			// Get resource info for confirmation
			data, err := client.Resources().Get(uint(id))
			if err != nil {
				return fmt.Errorf("failed to get resource: %w", err)
			}

			var resource Resource
			if err := json.Unmarshal(data, &resource); err != nil {
				return fmt.Errorf("failed to parse resource: %w", err)
			}

			// Confirm deletion unless forced
			if !force {
				fmt.Printf("Are you sure you want to delete resource '%s' (ID: %d)? [y/N]: ", resource.Name, resource.ID)
				var response string
				_, _ = fmt.Scanln(&response)
				if response != "y" && response != "Y" && response != "yes" {
					fmt.Println("❌ Deletion cancelled.")
					return nil
				}
			}

			if err := client.Resources().Delete(uint(id)); err != nil {
				return fmt.Errorf("failed to delete resource: %w", err)
			}

			fmt.Printf("✅ Resource '%s' (ID: %d) deleted successfully.\n", resource.Name, resource.ID)
			return nil
		},
	}

	cmd.Flags().BoolVarP(&force, "force", "f", false, "Skip confirmation prompt")
	return cmd
}
