package models

import (
	"fmt"
	"path/filepath"
	"strings"
)

// OutputFile represents the final merged and formatted CSV output
type OutputFile struct {
	Path        string      // Output file path (always .csv extension)
	Headers     []string    // Union of all input file headers
	Records     []DataEntry // Deduplicated and merged data entries
	AnkiHeaders []string    // Anki-specific header lines
}

// NewOutputFile creates a new OutputFile instance
func NewOutputFile(path string) *OutputFile {
	// Ensure output file has .csv extension
	if !strings.HasSuffix(strings.ToLower(path), ".csv") {
		path = strings.TrimSuffix(path, filepath.Ext(path)) + ".csv"
	}

	return &OutputFile{
		Path:    path,
		Headers: []string{},
		Records: []DataEntry{},
		AnkiHeaders: []string{
			"#separator:comma",
			"#html:true",
		},
	}
}

// Validate checks if the output file meets all validation requirements
func (f *OutputFile) Validate() error {
	// Path must be writable location
	dir := filepath.Dir(f.Path)
	if dir == "" {
		return fmt.Errorf("invalid output path")
	}

	// Headers must be union of all input headers
	if len(f.Headers) == 0 {
		return fmt.Errorf("output file must have at least one header")
	}

	// Records must be deduplicated (this is enforced by the processing logic)

	// AnkiHeaders must include required directives
	hasSeperator := false
	hasHTML := false
	hasColumns := false

	for _, header := range f.AnkiHeaders {
		if strings.HasPrefix(header, "#separator:") {
			hasSeperator = true
		}
		if strings.HasPrefix(header, "#html:") {
			hasHTML = true
		}
		if strings.HasPrefix(header, "#columns:") {
			hasColumns = true
		}
	}

	if !hasSeperator || !hasHTML || !hasColumns {
		return fmt.Errorf("output file missing required Anki headers")
	}

	return nil
}

// MergeHeaders combines headers from multiple input files (union operation)
func (f *OutputFile) MergeHeaders(inputFiles []*InputFile) {
	headerSet := make(map[string]bool)

	// Add headers from each input file
	for _, inputFile := range inputFiles {
		for _, header := range inputFile.Headers {
			if header != "" && !headerSet[header] {
				headerSet[header] = true
				f.Headers = append(f.Headers, header)
			}
		}
	}

	// Update the columns header for Anki
	f.updateColumnsHeader()
}

// AddRecord adds a data entry to the output file
func (f *OutputFile) AddRecord(entry DataEntry) {
	f.Records = append(f.Records, entry)
}

// GetRecordCount returns the number of records in the output
func (f *OutputFile) GetRecordCount() int {
	return len(f.Records)
}

// updateColumnsHeader updates the #columns: header with current column list
func (f *OutputFile) updateColumnsHeader() {
	columnsHeader := "#columns:" + strings.Join(f.Headers, ",")

	// Remove any existing columns header
	filtered := []string{}
	for _, header := range f.AnkiHeaders {
		if !strings.HasPrefix(header, "#columns:") {
			filtered = append(filtered, header)
		}
	}

	// Add the new columns header
	f.AnkiHeaders = append(filtered, columnsHeader)
}

// GetAnkiHeaderLines returns the Anki header lines as strings
func (f *OutputFile) GetAnkiHeaderLines() []string {
	return f.AnkiHeaders
}

// GetCSVRecords returns all records as CSV-compatible string arrays
func (f *OutputFile) GetCSVRecords() [][]string {
	records := make([][]string, len(f.Records))
	for i, record := range f.Records {
		records[i] = record.ToCSVRecord(f.Headers)
	}
	return records
}
