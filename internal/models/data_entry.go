package models

import (
	"crypto/md5"
	"fmt"
	"strings"
)

// DataEntry represents a single row of data with field values
type DataEntry struct {
	Values     map[string]string // Column name to value mapping
	Source     string            // Originating file path
	LineNumber int               // Original line number in source file
}

// NewDataEntry creates a new DataEntry instance
func NewDataEntry(values map[string]string, source string, lineNumber int) *DataEntry {
	return &DataEntry{
		Values:     values,
		Source:     source,
		LineNumber: lineNumber,
	}
}

// Validate checks if the data entry meets all validation requirements
func (e *DataEntry) Validate() error {
	// Values map must not be empty
	if len(e.Values) == 0 {
		return fmt.Errorf("data entry must contain at least one field")
	}

	// All values must be valid UTF-8 strings (Go strings are UTF-8 by default)
	// Source must reference valid input file
	if e.Source == "" {
		return fmt.Errorf("data entry must reference source file")
	}

	return nil
}

// GetValue returns the value for the specified column name
func (e *DataEntry) GetValue(columnName string) string {
	if value, exists := e.Values[columnName]; exists {
		return value
	}
	return "" // Return empty string for missing columns
}

// SetValue sets the value for the specified column name
func (e *DataEntry) SetValue(columnName, value string) {
	if e.Values == nil {
		e.Values = make(map[string]string)
	}
	e.Values[columnName] = value
}

// GetHash returns a hash of all field values for duplicate detection
func (e *DataEntry) GetHash() string {
	// Create a consistent string representation of all values
	var keys []string
	for key := range e.Values {
		keys = append(keys, key)
	}

	// Sort keys for consistent hashing
	for i := 0; i < len(keys); i++ {
		for j := i + 1; j < len(keys); j++ {
			if keys[i] > keys[j] {
				keys[i], keys[j] = keys[j], keys[i]
			}
		}
	}

	var parts []string
	for _, key := range keys {
		parts = append(parts, fmt.Sprintf("%s:%s", key, e.Values[key]))
	}

	content := strings.Join(parts, "|")
	hash := md5.Sum([]byte(content))
	return fmt.Sprintf("%x", hash)
}

// IsExactDuplicate checks if this entry is an exact duplicate of another
func (e *DataEntry) IsExactDuplicate(other *DataEntry) bool {
	// Must have same number of values
	if len(e.Values) != len(other.Values) {
		return false
	}

	// All values must match exactly (case-sensitive)
	for key, value := range e.Values {
		otherValue, exists := other.Values[key]
		if !exists || value != otherValue {
			return false
		}
	}

	return true
}

// ToCSVRecord converts the DataEntry to a CSV record with specified column order
func (e *DataEntry) ToCSVRecord(columns []string) []string {
	record := make([]string, len(columns))
	for i, column := range columns {
		record[i] = e.GetValue(column)
	}
	return record
}
