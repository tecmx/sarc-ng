package resources

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/olekukonko/tablewriter"
)

// OutputFormat represents the output format for displaying data
type OutputFormat string

const (
	// TableFormat displays data in a table
	TableFormat OutputFormat = "table"
	// JSONFormat displays data as JSON
	JSONFormat OutputFormat = "json"
)

// OutputWithFormat displays resources in the specified format
func OutputWithFormat(resources []Resource, format OutputFormat) error {
	switch format {
	case JSONFormat:
		return OutputJSON(resources)
	default:
		return OutputTable(resources)
	}
}

// OutputJSON outputs resources as JSON
func OutputJSON(resources []Resource) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(resources)
}

// OutputTable outputs resources in a formatted table
func OutputTable(resources []Resource) error {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "Name", "Type", "Available", "Created", "Updated"})
	table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
	table.SetCenterSeparator("|")

	for _, resource := range resources {
		available := "Yes"
		if !resource.IsAvailable {
			available = "No"
		}

		table.Append([]string{
			fmt.Sprintf("%d", resource.ID),
			resource.Name,
			resource.Type,
			available,
			formatTime(resource.CreatedAt),
			formatTime(resource.UpdatedAt),
		})
	}

	table.Render()
	return nil
}

// formatTime formats a time.Time for display
func formatTime(t time.Time) string {
	if t.IsZero() {
		return "-"
	}
	return t.Format("2006-01-02 15:04:05")
}
