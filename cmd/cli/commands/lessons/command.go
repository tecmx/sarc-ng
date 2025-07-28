package lessons

import (
	"encoding/json"
	"fmt"
	"sarc-ng/pkg/rest/client"
	"strconv"
	"time"

	"github.com/spf13/cobra"
)

// NewCommand creates the lessons command group
func NewCommand(clientFactory func() *client.Client) *cobra.Command {
	lessonsCmd := &cobra.Command{
		Use:   "lessons",
		Short: "Manage lessons",
		Long:  "Create, read, update, and delete lessons in the SARC system.",
	}

	// Add subcommands
	lessonsCmd.AddCommand(newListCommand(clientFactory))
	lessonsCmd.AddCommand(newGetCommand(clientFactory))
	lessonsCmd.AddCommand(newCreateCommand(clientFactory))
	lessonsCmd.AddCommand(newUpdateCommand(clientFactory))
	lessonsCmd.AddCommand(newDeleteCommand(clientFactory))

	return lessonsCmd
}

// List all lessons
func newListCommand(clientFactory func() *client.Client) *cobra.Command {
	var outputFormat string

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all lessons",
		Long:  "Retrieve and display all lessons in the system.",
		RunE: func(cmd *cobra.Command, args []string) error {
			client := clientFactory()
			rawResp, err := client.Lessons().List(1, 100) // Get first 100 lessons
			if err != nil {
				return fmt.Errorf("failed to list lessons: %w", err)
			}

			// Parse the raw JSON response
			var response struct {
				Data []Lesson `json:"data"`
			}
			if err := json.Unmarshal(rawResp, &response); err != nil {
				// If it's not paginated, try parsing as direct array
				var lessons []Lesson
				if err := json.Unmarshal(rawResp, &lessons); err != nil {
					return fmt.Errorf("failed to parse response: %w", err)
				}
				response.Data = lessons
			}

			if len(response.Data) == 0 {
				fmt.Println("No lessons found.")
				return nil
			}

			return OutputWithFormat(response.Data, OutputFormat(outputFormat))
		},
	}

	cmd.Flags().StringVarP(&outputFormat, "output", "o", "table", "Output format (table, json)")
	return cmd
}

// Get a specific lesson
func newGetCommand(clientFactory func() *client.Client) *cobra.Command {
	var outputFormat string

	cmd := &cobra.Command{
		Use:   "get <id>",
		Short: "Get a lesson by ID",
		Long:  "Retrieve and display details for a specific lesson.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.ParseUint(args[0], 10, 32)
			if err != nil {
				return fmt.Errorf("invalid lesson ID: %s", args[0])
			}

			client := clientFactory()
			rawResp, err := client.Lessons().Get(uint(id))
			if err != nil {
				return fmt.Errorf("failed to get lesson: %w", err)
			}

			// Parse the raw JSON response
			var lesson Lesson
			if err := json.Unmarshal(rawResp, &lesson); err != nil {
				return fmt.Errorf("failed to parse response: %w", err)
			}

			return OutputWithFormat([]Lesson{lesson}, OutputFormat(outputFormat))
		},
	}

	cmd.Flags().StringVarP(&outputFormat, "output", "o", "table", "Output format (table, json)")
	return cmd
}

// Create a new lesson
func newCreateCommand(clientFactory func() *client.Client) *cobra.Command {
	var title string
	var duration int
	var startTime string

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new lesson",
		Long:  "Create a new lesson with the specified title, duration, and start time.",
		RunE: func(cmd *cobra.Command, args []string) error {
			if title == "" {
				return fmt.Errorf("lesson title is required")
			}
			if duration <= 0 {
				return fmt.Errorf("lesson duration must be greater than 0")
			}

			var parsedTime time.Time
			var err error
			if startTime != "" {
				parsedTime, err = time.Parse("2006-01-02 15:04:05", startTime)
				if err != nil {
					return fmt.Errorf("invalid start time format. Use YYYY-MM-DD HH:MM:SS")
				}
			}

			client := clientFactory()
			req := LessonRequest{
				Title:     title,
				Duration:  duration,
				StartTime: parsedTime,
			}

			rawResp, err := client.Lessons().Create(req)
			if err != nil {
				return fmt.Errorf("failed to create lesson: %w", err)
			}

			// Parse the raw JSON response
			var lesson Lesson
			if err := json.Unmarshal(rawResp, &lesson); err != nil {
				return fmt.Errorf("failed to parse response: %w", err)
			}

			fmt.Printf("✅ Lesson created successfully:\n")
			return OutputTable([]Lesson{lesson})
		},
	}

	cmd.Flags().StringVarP(&title, "title", "t", "", "Lesson title (required)")
	cmd.Flags().IntVarP(&duration, "duration", "d", 0, "Lesson duration in minutes (required)")
	cmd.Flags().StringVarP(&startTime, "start-time", "s", "", "Start time (YYYY-MM-DD HH:MM:SS)")
	_ = cmd.MarkFlagRequired("title")
	_ = cmd.MarkFlagRequired("duration")

	return cmd
}

// Update an existing lesson
func newUpdateCommand(clientFactory func() *client.Client) *cobra.Command {
	var title string
	var duration int
	var startTime string

	cmd := &cobra.Command{
		Use:   "update <id>",
		Short: "Update a lesson",
		Long:  "Update an existing lesson's title, duration, and/or start time.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.ParseUint(args[0], 10, 32)
			if err != nil {
				return fmt.Errorf("invalid lesson ID: %s", args[0])
			}

			client := clientFactory()

			// Get current lesson to preserve unchanged fields
			rawCurrentResp, err := client.Lessons().Get(uint(id))
			if err != nil {
				return fmt.Errorf("failed to get current lesson: %w", err)
			}

			var current Lesson
			if err := json.Unmarshal(rawCurrentResp, &current); err != nil {
				return fmt.Errorf("failed to parse current lesson: %w", err)
			}

			// Use current values if not specified
			if title == "" {
				title = current.Title
			}
			if duration == 0 {
				duration = current.Duration
			}

			parsedTime := current.StartTime
			if startTime != "" {
				parsedTime, err = time.Parse("2006-01-02 15:04:05", startTime)
				if err != nil {
					return fmt.Errorf("invalid start time format. Use YYYY-MM-DD HH:MM:SS")
				}
			}

			req := LessonRequest{
				Title:     title,
				Duration:  duration,
				StartTime: parsedTime,
			}

			rawResp, err := client.Lessons().Update(uint(id), req)
			if err != nil {
				return fmt.Errorf("failed to update lesson: %w", err)
			}

			// Parse the raw JSON response
			var lesson Lesson
			if err := json.Unmarshal(rawResp, &lesson); err != nil {
				return fmt.Errorf("failed to parse response: %w", err)
			}

			fmt.Printf("✅ Lesson updated successfully:\n")
			return OutputTable([]Lesson{lesson})
		},
	}

	cmd.Flags().StringVarP(&title, "title", "t", "", "Lesson title")
	cmd.Flags().IntVarP(&duration, "duration", "d", 0, "Lesson duration in minutes")
	cmd.Flags().StringVarP(&startTime, "start-time", "s", "", "Start time (YYYY-MM-DD HH:MM:SS)")

	return cmd
}

// Delete a lesson
func newDeleteCommand(clientFactory func() *client.Client) *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "delete <id>",
		Short: "Delete a lesson",
		Long:  "Delete a lesson by ID. Use --force to skip confirmation.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.ParseUint(args[0], 10, 32)
			if err != nil {
				return fmt.Errorf("invalid lesson ID: %s", args[0])
			}

			client := clientFactory()

			// Get lesson info for confirmation
			rawResp, err := client.Lessons().Get(uint(id))
			if err != nil {
				return fmt.Errorf("failed to get lesson: %w", err)
			}

			var lesson Lesson
			if err := json.Unmarshal(rawResp, &lesson); err != nil {
				return fmt.Errorf("failed to parse lesson: %w", err)
			}

			if !force {
				fmt.Printf("Are you sure you want to delete lesson '%s' (ID: %d)? [y/N]: ", lesson.Title, lesson.ID)
				var response string
				fmt.Scanln(&response)
				if response != "y" && response != "Y" && response != "yes" && response != "Yes" {
					fmt.Println("Operation cancelled.")
					return nil
				}
			}

			err = client.Lessons().Delete(uint(id))
			if err != nil {
				return fmt.Errorf("failed to delete lesson: %w", err)
			}

			fmt.Printf("✅ Lesson '%s' deleted successfully.\n", lesson.Title)
			return nil
		},
	}

	cmd.Flags().BoolVarP(&force, "force", "f", false, "Force deletion without confirmation")
	return cmd
}
