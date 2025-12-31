package output

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"text/tabwriter"
)

// Format represents the output format
type Format int

const (
	// FormatTable outputs data as a formatted table
	FormatTable Format = iota
	// FormatJSON outputs data as JSON
	FormatJSON
)

// Formatter handles output formatting
type Formatter struct {
	format Format
	writer io.Writer
}

// NewFormatter creates a new formatter
func NewFormatter(jsonOutput bool) *Formatter {
	format := FormatTable
	if jsonOutput {
		format = FormatJSON
	}
	return &Formatter{
		format: format,
		writer: os.Stdout,
	}
}

// SetWriter sets the output writer (useful for testing)
func (f *Formatter) SetWriter(w io.Writer) {
	f.writer = w
}

// PrintJSON outputs data as formatted JSON
func (f *Formatter) PrintJSON(data interface{}) error {
	encoder := json.NewEncoder(f.writer)
	encoder.SetIndent("", "  ")
	return encoder.Encode(data)
}

// PrintTable outputs data as a formatted table using tabwriter
func (f *Formatter) PrintTable(headers []string, rows [][]string) {
	if len(rows) == 0 {
		_, _ = fmt.Fprintln(f.writer, "No results found.")
		return
	}

	w := tabwriter.NewWriter(f.writer, 0, 0, 2, ' ', 0)

	// Print header
	_, _ = fmt.Fprintln(w, strings.Join(headers, "\t"))

	// Print rows
	for _, row := range rows {
		_, _ = fmt.Fprintln(w, strings.Join(row, "\t"))
	}

	_ = w.Flush()
}

// Print outputs data in the configured format
func (f *Formatter) Print(headers []string, rows [][]string, jsonData interface{}) error {
	switch f.format {
	case FormatJSON:
		return f.PrintJSON(jsonData)
	default:
		f.PrintTable(headers, rows)
		return nil
	}
}

// PrintDetail prints a single item's details
func (f *Formatter) PrintDetail(fields []DetailField, jsonData interface{}) error {
	switch f.format {
	case FormatJSON:
		return f.PrintJSON(jsonData)
	default:
		maxLabelLen := 0
		for _, field := range fields {
			if len(field.Label) > maxLabelLen {
				maxLabelLen = len(field.Label)
			}
		}

		for _, field := range fields {
			if field.Value == "" && !field.ShowEmpty {
				continue
			}
			padding := strings.Repeat(" ", maxLabelLen-len(field.Label))
			_, _ = fmt.Fprintf(f.writer, "%s:%s  %s\n", field.Label, padding, field.Value)
		}
		return nil
	}
}

// DetailField represents a field in detail view
type DetailField struct {
	Label     string
	Value     string
	ShowEmpty bool
}

// PriorityLabel returns a human-readable priority label
func PriorityLabel(priority int) string {
	switch priority {
	case 0:
		return "No priority"
	case 1:
		return "Urgent"
	case 2:
		return "High"
	case 3:
		return "Medium"
	case 4:
		return "Low"
	default:
		return fmt.Sprintf("Priority %d", priority)
	}
}

// Truncate truncates a string to the specified length
func Truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	if maxLen <= 3 {
		return s[:maxLen]
	}
	return s[:maxLen-3] + "..."
}

// FormatPercentage formats a float as a percentage
func FormatPercentage(value float64) string {
	return fmt.Sprintf("%.0f%%", value*100)
}

// EmptyIfNil returns "-" if the string pointer is nil
func EmptyIfNil(s *string) string {
	if s == nil {
		return "-"
	}
	return *s
}

// NameOrDash returns the name or "-" if nil
func NameOrDash(name *string) string {
	if name == nil || *name == "" {
		return "-"
	}
	return *name
}
