package models

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// InputFile represents a source CSV/TSV file to be processed
type InputFile struct {
	Path      string     // Absolute file path
	Separator rune       // Field separator (comma or tab)
	Headers   []string   // Column header names
	Records   [][]string // Data rows (excluding header)
	Encoding  string     // Character encoding (UTF-8 only)
}

// NewInputFile creates a new InputFile instance with the given path
func NewInputFile(path string) *InputFile {
	return &InputFile{
		Path:      path,
		Separator: ',', // Default to comma
		Encoding:  "UTF-8",
	}
}

// Validate checks if the input file meets all validation requirements
func (f *InputFile) Validate() error {
	// Check if path exists and is readable
	if _, err := os.Stat(f.Path); os.IsNotExist(err) {
		return fmt.Errorf("file not found: %s", f.Path)
	}

	// Check if file is readable
	file, err := os.Open(f.Path)
	if err != nil {
		return fmt.Errorf("file not readable: %v", err)
	}
	defer file.Close()

	// Validate separator (must be comma or tab)
	if f.Separator != ',' && f.Separator != '\t' {
		return fmt.Errorf("invalid separator: must be comma or tab")
	}

	// Check if encoding is UTF-8 (simplified check)
	if f.Encoding != "UTF-8" {
		return fmt.Errorf("invalid encoding: only UTF-8 supported")
	}

	// Must contain at least one data row
	if len(f.Records) == 0 {
		return fmt.Errorf("file contains no data rows")
	}

	return nil
}

// DetectSeparator attempts to detect the file separator based on file extension
func (f *InputFile) DetectSeparator() {
	ext := strings.ToLower(filepath.Ext(f.Path))
	switch ext {
	case ".tsv":
		f.Separator = '\t'
	case ".csv":
		f.Separator = ','
	default:
		// Default to comma if extension is unclear
		f.Separator = ','
	}
}

// GetSeparatorString returns the separator as a string for display purposes
func (f *InputFile) GetSeparatorString() string {
	if f.Separator == '\t' {
		return "tab"
	}
	return "comma"
}
