package buildings

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

// OutputWithFormat displays buildings in the specified format
func OutputWithFormat(buildings []Building, format OutputFormat) error {
	switch format {
	case JSONFormat:
		return OutputJSON(buildings)
	default:
		return OutputTable(buildings)
	}
}

// OutputJSON outputs buildings as JSON
func OutputJSON(buildings []Building) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(buildings)
}

// OutputTable outputs buildings in a formatted table
func OutputTable(buildings []Building) error {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "Name", "Code", "Created", "Updated"})
	table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
	table.SetCenterSeparator("|")

	for _, building := range buildings {
		table.Append([]string{
			fmt.Sprintf("%d", building.ID),
			building.Name,
			building.Code,
			formatTime(building.CreatedAt),
			formatTime(building.UpdatedAt),
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
