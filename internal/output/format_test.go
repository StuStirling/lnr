package output

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewFormatter_Table(t *testing.T) {
	f := NewFormatter(false)
	assert.Equal(t, FormatTable, f.format)
}

func TestNewFormatter_JSON(t *testing.T) {
	f := NewFormatter(true)
	assert.Equal(t, FormatJSON, f.format)
}

func TestPrintJSON(t *testing.T) {
	var buf bytes.Buffer
	f := NewFormatter(true)
	f.SetWriter(&buf)

	data := map[string]string{"key": "value"}
	err := f.PrintJSON(data)

	require.NoError(t, err)

	var result map[string]string
	err = json.Unmarshal(buf.Bytes(), &result)
	require.NoError(t, err)
	assert.Equal(t, "value", result["key"])
}

func TestPrintTable_Empty(t *testing.T) {
	var buf bytes.Buffer
	f := NewFormatter(false)
	f.SetWriter(&buf)

	f.PrintTable([]string{"A", "B"}, [][]string{})

	assert.Contains(t, buf.String(), "No results found.")
}

func TestPrintTable_WithData(t *testing.T) {
	var buf bytes.Buffer
	f := NewFormatter(false)
	f.SetWriter(&buf)

	headers := []string{"NAME", "VALUE"}
	rows := [][]string{
		{"foo", "bar"},
		{"baz", "qux"},
	}
	f.PrintTable(headers, rows)

	output := buf.String()
	assert.Contains(t, output, "NAME")
	assert.Contains(t, output, "VALUE")
	assert.Contains(t, output, "foo")
	assert.Contains(t, output, "bar")
}

func TestPrintDetail_Table(t *testing.T) {
	var buf bytes.Buffer
	f := NewFormatter(false)
	f.SetWriter(&buf)

	fields := []DetailField{
		{Label: "Name", Value: "Test"},
		{Label: "ID", Value: "123"},
	}

	err := f.PrintDetail(fields, nil)
	require.NoError(t, err)

	output := buf.String()
	assert.Contains(t, output, "Name:")
	assert.Contains(t, output, "Test")
	assert.Contains(t, output, "ID:")
	assert.Contains(t, output, "123")
}

func TestPrintDetail_JSON(t *testing.T) {
	var buf bytes.Buffer
	f := NewFormatter(true)
	f.SetWriter(&buf)

	data := map[string]string{"name": "Test", "id": "123"}
	fields := []DetailField{} // Not used in JSON mode

	err := f.PrintDetail(fields, data)
	require.NoError(t, err)

	var result map[string]string
	err = json.Unmarshal(buf.Bytes(), &result)
	require.NoError(t, err)
	assert.Equal(t, "Test", result["name"])
}

func TestPriorityLabel(t *testing.T) {
	tests := []struct {
		priority int
		expected string
	}{
		{0, "No priority"},
		{1, "Urgent"},
		{2, "High"},
		{3, "Medium"},
		{4, "Low"},
		{5, "Priority 5"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			assert.Equal(t, tt.expected, PriorityLabel(tt.priority))
		})
	}
}

func TestTruncate(t *testing.T) {
	tests := []struct {
		input    string
		maxLen   int
		expected string
	}{
		{"short", 10, "short"},
		{"this is a long string", 10, "this is..."},
		{"ab", 2, "ab"},
		{"abc", 3, "abc"},
		{"abcd", 3, "abc"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			assert.Equal(t, tt.expected, Truncate(tt.input, tt.maxLen))
		})
	}
}

func TestFormatPercentage(t *testing.T) {
	assert.Equal(t, "50%", FormatPercentage(0.5))
	assert.Equal(t, "100%", FormatPercentage(1.0))
	assert.Equal(t, "0%", FormatPercentage(0.0))
	assert.Equal(t, "75%", FormatPercentage(0.75))
}

func TestEmptyIfNil(t *testing.T) {
	assert.Equal(t, "-", EmptyIfNil(nil))
	s := "value"
	assert.Equal(t, "value", EmptyIfNil(&s))
}

func TestNameOrDash(t *testing.T) {
	assert.Equal(t, "-", NameOrDash(nil))
	empty := ""
	assert.Equal(t, "-", NameOrDash(&empty))
	name := "John"
	assert.Equal(t, "John", NameOrDash(&name))
}
