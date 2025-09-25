package services

import (
	"fmt"

	"ankiprep/internal/models"
)

// ColumnMerger handles merging columns from multiple input files
type ColumnMerger struct{}

// NewColumnMerger creates a new ColumnMerger instance
func NewColumnMerger() *ColumnMerger {
	return &ColumnMerger{}
}

// MergeColumns combines column headers from multiple input files using union operation
func (m *ColumnMerger) MergeColumns(inputFiles []*models.InputFile) ([]string, error) {
	if len(inputFiles) == 0 {
		return nil, fmt.Errorf("no input files provided for column merging")
	}

	// Use a map to track unique column names and preserve order
	columnMap := make(map[string]bool)
	var mergedColumns []string

	// Process each input file in order
	for _, inputFile := range inputFiles {
		if inputFile == nil {
			continue
		}

		for _, column := range inputFile.Headers {
			// Skip empty column names
			if column == "" {
				continue
			}

			// Add column if not already seen
			if !columnMap[column] {
				columnMap[column] = true
				mergedColumns = append(mergedColumns, column)
			}
		}
	}

	if len(mergedColumns) == 0 {
		return nil, fmt.Errorf("no valid column headers found in input files")
	}

	return mergedColumns, nil
}

// NormalizeDataEntries ensures all data entries have values for all merged columns
func (m *ColumnMerger) NormalizeDataEntries(entries []*models.DataEntry, mergedColumns []string) []*models.DataEntry {
	normalizedEntries := make([]*models.DataEntry, len(entries))

	for i, entry := range entries {
		// Create a new entry with normalized values
		normalizedValues := make(map[string]string)

		// Copy existing values
		for key, value := range entry.Values {
			normalizedValues[key] = value
		}

		// Ensure all merged columns exist (fill missing with empty string)
		for _, column := range mergedColumns {
			if _, exists := normalizedValues[column]; !exists {
				normalizedValues[column] = ""
			}
		}

		// Create normalized entry
		normalizedEntry := models.NewDataEntry(normalizedValues, entry.Source, entry.LineNumber)
		normalizedEntries[i] = normalizedEntry
	}

	return normalizedEntries
}

// GetColumnMapping returns a mapping of original columns to merged position
func (m *ColumnMerger) GetColumnMapping(originalColumns []string, mergedColumns []string) map[string]int {
	mapping := make(map[string]int)

	for _, originalCol := range originalColumns {
		for i, mergedCol := range mergedColumns {
			if originalCol == mergedCol {
				mapping[originalCol] = i
				break
			}
		}
	}

	return mapping
}

// ValidateColumnCompatibility checks if columns from different files can be merged
func (m *ColumnMerger) ValidateColumnCompatibility(inputFiles []*models.InputFile) error {
	if len(inputFiles) <= 1 {
		return nil // Single file or no files - no compatibility issues
	}

	// Check for any obvious incompatibilities
	// For CSV merging, any column name is compatible with any other
	// This method is mainly for future extensibility

	columnCounts := make(map[string]int)

	for _, inputFile := range inputFiles {
		for _, column := range inputFile.Headers {
			if column != "" {
				columnCounts[column]++
			}
		}
	}

	// For now, just validate that we have some columns
	if len(columnCounts) == 0 {
		return fmt.Errorf("no valid columns found across input files")
	}

	return nil
}

// GetColumnStatistics returns statistics about column merging
func (m *ColumnMerger) GetColumnStatistics(inputFiles []*models.InputFile, mergedColumns []string) map[string]interface{} {
	stats := make(map[string]interface{})

	// Count unique columns per file
	fileColumnCounts := make(map[string]int)
	allColumns := make(map[string]int) // column -> number of files containing it

	for _, inputFile := range inputFiles {
		uniqueColumns := make(map[string]bool)
		for _, column := range inputFile.Headers {
			if column != "" && !uniqueColumns[column] {
				uniqueColumns[column] = true
				fileColumnCounts[inputFile.Path]++
				allColumns[column]++
			}
		}
	}

	// Calculate statistics
	totalFiles := len(inputFiles)
	commonColumns := 0
	fileSpecificColumns := 0

	for _, fileCount := range allColumns {
		if fileCount == totalFiles {
			commonColumns++
		} else {
			fileSpecificColumns++
		}
	}

	stats["total_merged_columns"] = len(mergedColumns)
	stats["common_columns"] = commonColumns
	stats["file_specific_columns"] = fileSpecificColumns
	stats["file_column_counts"] = fileColumnCounts
	stats["column_file_distribution"] = allColumns

	return stats
}

// ReorderColumns reorders columns according to a specified order preference
func (m *ColumnMerger) ReorderColumns(columns []string, preferredOrder []string) []string {
	if len(preferredOrder) == 0 {
		return columns
	}

	// Create a map for quick lookup
	columnSet := make(map[string]bool)
	for _, col := range columns {
		columnSet[col] = true
	}

	var reordered []string
	used := make(map[string]bool)

	// First, add columns in preferred order if they exist
	for _, preferred := range preferredOrder {
		if columnSet[preferred] && !used[preferred] {
			reordered = append(reordered, preferred)
			used[preferred] = true
		}
	}

	// Then add remaining columns in their original order
	for _, col := range columns {
		if !used[col] {
			reordered = append(reordered, col)
			used[col] = true
		}
	}

	return reordered
}
