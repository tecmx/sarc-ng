package lessons

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

// OutputWithFormat displays lessons in the specified format
func OutputWithFormat(lessons []Lesson, format OutputFormat) error {
	switch format {
	case JSONFormat:
		return OutputJSON(lessons)
	default:
		return OutputTable(lessons)
	}
}

// OutputJSON outputs lessons as JSON
func OutputJSON(lessons []Lesson) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(lessons)
}

// OutputTable outputs lessons in a formatted table
func OutputTable(lessons []Lesson) error {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "Title", "Duration", "Start Time", "End Time", "Created", "Updated"})
	table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
	table.SetCenterSeparator("|")

	for _, lesson := range lessons {
		table.Append([]string{
			fmt.Sprintf("%d", lesson.ID),
			lesson.Title,
			fmt.Sprintf("%d min", lesson.Duration),
			formatTime(lesson.StartTime),
			formatTime(lesson.EndTime),
			formatTime(lesson.CreatedAt),
			formatTime(lesson.UpdatedAt),
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
