package classes

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

// OutputWithFormat displays classes in the specified format
func OutputWithFormat(classes []Class, format OutputFormat) error {
	switch format {
	case JSONFormat:
		return OutputJSON(classes)
	default:
		return OutputTable(classes)
	}
}

// OutputJSON outputs classes as JSON
func OutputJSON(classes []Class) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(classes)
}

// OutputTable outputs classes in a formatted table
func OutputTable(classes []Class) error {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "Name", "Capacity", "Created", "Updated"})
	table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
	table.SetCenterSeparator("|")

	for _, class := range classes {
		table.Append([]string{
			fmt.Sprintf("%d", class.ID),
			class.Name,
			fmt.Sprintf("%d", class.Capacity),
			formatTime(class.CreatedAt),
			formatTime(class.UpdatedAt),
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
