package reservations

import (
	"encoding/json"
	"fmt"
	"sarc-ng/pkg/rest/client"
	"strconv"
	"time"

	"github.com/spf13/cobra"
)

// NewCommand creates the reservations command group
func NewCommand(clientFactory func() *client.Client) *cobra.Command {
	reservationsCmd := &cobra.Command{
		Use:   "reservations",
		Short: "Manage reservations",
		Long:  "Create, read, update, and delete reservations in the SARC system.",
	}

	// Add subcommands
	reservationsCmd.AddCommand(newListCommand(clientFactory))
	reservationsCmd.AddCommand(newGetCommand(clientFactory))
	reservationsCmd.AddCommand(newCreateCommand(clientFactory))
	reservationsCmd.AddCommand(newUpdateCommand(clientFactory))
	reservationsCmd.AddCommand(newDeleteCommand(clientFactory))

	return reservationsCmd
}

// List all reservations
func newListCommand(clientFactory func() *client.Client) *cobra.Command {
	var outputFormat string

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all reservations",
		Long:  "Retrieve and display all reservations in the system.",
		RunE: func(cmd *cobra.Command, args []string) error {
			client := clientFactory()
			data, err := client.Reservations().List(1, 100) // Default pagination
			if err != nil {
				return fmt.Errorf("failed to list reservations: %w", err)
			}

			var reservations []Reservation
			if err := json.Unmarshal(data, &reservations); err != nil {
				return fmt.Errorf("failed to parse reservations: %w", err)
			}

			if len(reservations) == 0 {
				fmt.Println("No reservations found.")
				return nil
			}

			return OutputWithFormat(reservations, OutputFormat(outputFormat))
		},
	}

	cmd.Flags().StringVarP(&outputFormat, "output", "o", "table", "Output format (table, json)")
	return cmd
}

// Get a specific reservation
func newGetCommand(clientFactory func() *client.Client) *cobra.Command {
	var outputFormat string

	cmd := &cobra.Command{
		Use:   "get <id>",
		Short: "Get a reservation by ID",
		Long:  "Retrieve and display details for a specific reservation.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.ParseUint(args[0], 10, 32)
			if err != nil {
				return fmt.Errorf("invalid reservation ID: %s", args[0])
			}

			client := clientFactory()
			data, err := client.Reservations().Get(uint(id))
			if err != nil {
				return fmt.Errorf("failed to get reservation: %w", err)
			}

			var reservation Reservation
			if err := json.Unmarshal(data, &reservation); err != nil {
				return fmt.Errorf("failed to parse reservation: %w", err)
			}

			return OutputWithFormat([]Reservation{reservation}, OutputFormat(outputFormat))
		},
	}

	cmd.Flags().StringVarP(&outputFormat, "output", "o", "table", "Output format (table, json)")
	return cmd
}

// Create a new reservation
func newCreateCommand(clientFactory func() *client.Client) *cobra.Command {
	var resourceID, userID uint
	var startTime, endTime string

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new reservation",
		Long:  "Create a new reservation for a resource.",
		RunE: func(cmd *cobra.Command, args []string) error {
			if resourceID == 0 {
				return fmt.Errorf("resource ID is required")
			}
			if userID == 0 {
				return fmt.Errorf("user ID is required")
			}
			if startTime == "" {
				return fmt.Errorf("start time is required")
			}
			if endTime == "" {
				return fmt.Errorf("end time is required")
			}

			// Parse time strings
			start, err := time.Parse(time.RFC3339, startTime)
			if err != nil {
				return fmt.Errorf("invalid start time format: %w", err)
			}
			end, err := time.Parse(time.RFC3339, endTime)
			if err != nil {
				return fmt.Errorf("invalid end time format: %w", err)
			}

			client := clientFactory()
			req := ReservationRequest{
				ResourceID: resourceID,
				UserID:     userID,
				StartTime:  start,
				EndTime:    end,
			}

			data, err := client.Reservations().Create(req)
			if err != nil {
				return fmt.Errorf("failed to create reservation: %w", err)
			}

			var reservation Reservation
			if err := json.Unmarshal(data, &reservation); err != nil {
				return fmt.Errorf("failed to parse created reservation: %w", err)
			}

			fmt.Printf("✅ Reservation created successfully:\n")
			return OutputTable([]Reservation{reservation})
		},
	}

	cmd.Flags().UintVarP(&resourceID, "resource-id", "r", 0, "Resource ID (required)")
	cmd.Flags().UintVarP(&userID, "user-id", "u", 0, "User ID (required)")
	cmd.Flags().StringVarP(&startTime, "start-time", "s", "", "Start time (ISO format, required)")
	cmd.Flags().StringVarP(&endTime, "end-time", "e", "", "End time (ISO format, required)")
	_ = cmd.MarkFlagRequired("resource-id")
	_ = cmd.MarkFlagRequired("user-id")
	_ = cmd.MarkFlagRequired("start-time")
	_ = cmd.MarkFlagRequired("end-time")

	return cmd
}

// Update an existing reservation
func newUpdateCommand(clientFactory func() *client.Client) *cobra.Command {
	var resourceID, userID uint
	var startTime, endTime string

	cmd := &cobra.Command{
		Use:   "update <id>",
		Short: "Update a reservation",
		Long:  "Update an existing reservation's details.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.ParseUint(args[0], 10, 32)
			if err != nil {
				return fmt.Errorf("invalid reservation ID: %s", args[0])
			}

			client := clientFactory()

			// Get current reservation to preserve unchanged fields
			data, err := client.Reservations().Get(uint(id))
			if err != nil {
				return fmt.Errorf("failed to get current reservation: %w", err)
			}

			var current Reservation
			if err := json.Unmarshal(data, &current); err != nil {
				return fmt.Errorf("failed to parse current reservation: %w", err)
			}

			// Use current values if not specified
			if resourceID == 0 {
				resourceID = current.ResourceID
			}
			if userID == 0 {
				userID = current.UserID
			}

			var start, end time.Time
			if startTime == "" {
				start = current.StartTime
			} else {
				start, err = time.Parse(time.RFC3339, startTime)
				if err != nil {
					return fmt.Errorf("invalid start time format: %w", err)
				}
			}
			if endTime == "" {
				end = current.EndTime
			} else {
				end, err = time.Parse(time.RFC3339, endTime)
				if err != nil {
					return fmt.Errorf("invalid end time format: %w", err)
				}
			}

			req := ReservationRequest{
				ResourceID: resourceID,
				UserID:     userID,
				StartTime:  start,
				EndTime:    end,
			}

			updateData, err := client.Reservations().Update(uint(id), req)
			if err != nil {
				return fmt.Errorf("failed to update reservation: %w", err)
			}

			var reservation Reservation
			if err := json.Unmarshal(updateData, &reservation); err != nil {
				return fmt.Errorf("failed to parse updated reservation: %w", err)
			}

			fmt.Printf("✅ Reservation updated successfully:\n")
			return OutputTable([]Reservation{reservation})
		},
	}

	cmd.Flags().UintVarP(&resourceID, "resource-id", "r", 0, "Resource ID")
	cmd.Flags().UintVarP(&userID, "user-id", "u", 0, "User ID")
	cmd.Flags().StringVarP(&startTime, "start-time", "s", "", "Start time (ISO format)")
	cmd.Flags().StringVarP(&endTime, "end-time", "e", "", "End time (ISO format)")

	return cmd
}

// Delete a reservation
func newDeleteCommand(clientFactory func() *client.Client) *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "delete <id>",
		Short: "Delete a reservation",
		Long:  "Delete a reservation by ID. Use --force to skip confirmation.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.ParseUint(args[0], 10, 32)
			if err != nil {
				return fmt.Errorf("invalid reservation ID: %s", args[0])
			}

			client := clientFactory()

			// Get reservation info for confirmation
			data, err := client.Reservations().Get(uint(id))
			if err != nil {
				return fmt.Errorf("failed to get reservation: %w", err)
			}

			var reservation Reservation
			if err := json.Unmarshal(data, &reservation); err != nil {
				return fmt.Errorf("failed to parse reservation: %w", err)
			}

			// Confirm deletion unless forced
			if !force {
				fmt.Printf("Are you sure you want to delete reservation ID %d? [y/N]: ", reservation.ID)
				var response string
				_, _ = fmt.Scanln(&response)
				if response != "y" && response != "Y" && response != "yes" {
					fmt.Println("❌ Deletion cancelled.")
					return nil
				}
			}

			if err := client.Reservations().Delete(uint(id)); err != nil {
				return fmt.Errorf("failed to delete reservation: %w", err)
			}

			fmt.Printf("✅ Reservation %d deleted successfully.\n", reservation.ID)
			return nil
		},
	}

	cmd.Flags().BoolVarP(&force, "force", "f", false, "Skip confirmation prompt")
	return cmd
}
