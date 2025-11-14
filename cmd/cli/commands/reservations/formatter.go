package reservations

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

// OutputWithFormat displays reservations in the specified format
func OutputWithFormat(reservations []Reservation, format OutputFormat) error {
	switch format {
	case JSONFormat:
		return OutputJSON(reservations)
	default:
		return OutputTable(reservations)
	}
}

// OutputJSON outputs reservations as JSON
func OutputJSON(reservations []Reservation) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(reservations)
}

// OutputTable outputs reservations in a formatted table
func OutputTable(reservations []Reservation) error {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "Resource ID", "User ID", "Start Time", "End Time", "Status", "Created", "Updated"})
	table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
	table.SetCenterSeparator("|")

	for _, reservation := range reservations {
		table.Append([]string{
			fmt.Sprintf("%d", reservation.ID),
			fmt.Sprintf("%d", reservation.ResourceID),
			fmt.Sprintf("%d", reservation.UserID),
			formatTime(reservation.StartTime),
			formatTime(reservation.EndTime),
			reservation.Status,
			formatTime(reservation.CreatedAt),
			formatTime(reservation.UpdatedAt),
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
