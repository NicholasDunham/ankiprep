package services

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"

	"ankiprep/internal/models"
)

// AnkiFormatter handles formatting data for Anki import
type AnkiFormatter struct{}

// NewAnkiFormatter creates a new AnkiFormatter instance
func NewAnkiFormatter() *AnkiFormatter {
	return &AnkiFormatter{}
}

// FormatForAnki creates an OutputFile formatted for Anki import
func (f *AnkiFormatter) FormatForAnki(entries []*models.DataEntry, columns []string, outputPath string) (*models.OutputFile, error) {
	if len(entries) == 0 {
		return nil, fmt.Errorf("no entries to format")
	}

	if len(columns) == 0 {
		return nil, fmt.Errorf("no columns specified")
	}

	// Create output file
	outputFile := models.NewOutputFile(outputPath)

	// Set merged headers
	outputFile.Headers = make([]string, len(columns))
	copy(outputFile.Headers, columns)

	// Add entries to output
	for _, entry := range entries {
		outputFile.AddRecord(*entry)
	}

	// Update Anki headers with column information
	f.updateAnkiHeaders(outputFile)

	return outputFile, nil
}

// WriteToFile writes the formatted output to a CSV file
func (f *AnkiFormatter) WriteToFile(outputFile *models.OutputFile) error {
	if outputFile == nil {
		return fmt.Errorf("output file cannot be nil")
	}

	// Create output file
	file, err := os.Create(outputFile.Path)
	if err != nil {
		return fmt.Errorf("failed to create output file %s: %v", outputFile.Path, err)
	}
	defer file.Close()

	// Write Anki header lines first
	for _, header := range outputFile.GetAnkiHeaderLines() {
		if _, err := file.WriteString(header + "\n"); err != nil {
			return fmt.Errorf("failed to write Anki header: %v", err)
		}
	}

	// Write column headers
	csvWriter := csv.NewWriter(file)
	defer csvWriter.Flush()

	if err := csvWriter.Write(outputFile.Headers); err != nil {
		return fmt.Errorf("failed to write column headers: %v", err)
	}

	// Write data records
	records := outputFile.GetCSVRecords()
	for _, record := range records {
		if err := csvWriter.Write(record); err != nil {
			return fmt.Errorf("failed to write data record: %v", err)
		}
	}

	return nil
}

// ValidateAnkiFormat ensures the output file meets Anki requirements
func (f *AnkiFormatter) ValidateAnkiFormat(outputFile *models.OutputFile) error {
	if outputFile == nil {
		return fmt.Errorf("output file cannot be nil")
	}

	// Check required Anki headers
	ankiHeaders := outputFile.GetAnkiHeaderLines()
	hasSeperator := false
	hasHTML := false
	hasColumns := false

	for _, header := range ankiHeaders {
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

	if !hasSeperator {
		return fmt.Errorf("missing required Anki header: #separator:")
	}
	if !hasHTML {
		return fmt.Errorf("missing required Anki header: #html:")
	}
	if !hasColumns {
		return fmt.Errorf("missing required Anki header: #columns:")
	}

	// Validate column count
	if len(outputFile.Headers) == 0 {
		return fmt.Errorf("output file has no columns")
	}

	// Validate that all records have consistent column count
	expectedColumns := len(outputFile.Headers)
	for i, record := range outputFile.Records {
		if len(record.Values) > expectedColumns {
			return fmt.Errorf("record %d has too many columns (%d, expected %d)",
				i, len(record.Values), expectedColumns)
		}
	}

	return nil
}

// updateAnkiHeaders updates the Anki-specific headers based on current content
func (f *AnkiFormatter) updateAnkiHeaders(outputFile *models.OutputFile) {
	// Build the columns header
	columnsHeader := "#columns:" + strings.Join(outputFile.Headers, ",")

	// Update Anki headers
	outputFile.AnkiHeaders = []string{
		"#separator:comma",
		"#html:true",
		columnsHeader,
	}
}

// GetAnkiCompatibleText ensures text is compatible with Anki's requirements
func (f *AnkiFormatter) GetAnkiCompatibleText(text string) string {
	// Anki accepts HTML, so preserve HTML tags
	// Ensure proper CSV escaping by letting the csv package handle it
	return text
}

// EstimateFileSize estimates the output file size in bytes
func (f *AnkiFormatter) EstimateFileSize(outputFile *models.OutputFile) int64 {
	if outputFile == nil {
		return 0
	}

	// Rough estimation
	headerSize := 0
	for _, header := range outputFile.GetAnkiHeaderLines() {
		headerSize += len(header) + 1 // +1 for newline
	}

	// Column headers
	columnHeaderSize := len(strings.Join(outputFile.Headers, ",")) + 1

	// Data records (rough estimate)
	avgRecordSize := 50 // Average characters per record
	dataSize := len(outputFile.Records) * avgRecordSize

	return int64(headerSize + columnHeaderSize + dataSize)
}

// GetAnkiImportInstructions returns instructions for importing the file into Anki
func (f *AnkiFormatter) GetAnkiImportInstructions() string {
	return `To import this file into Anki:
1. Open Anki and select the deck where you want to add the cards
2. Go to File > Import
3. Select this CSV file
4. Anki should automatically detect the format based on the headers
5. Review the field mapping and click Import
6. The cards will be created based on the column structure`
}

// WriteToFileWithOptions writes the output file with processing options
func (f *AnkiFormatter) WriteToFileWithOptions(outputFile *models.OutputFile, keepHeader bool) error {
	// Create output file
	file, err := os.Create(outputFile.Path)
	if err != nil {
		return fmt.Errorf("failed to create output file %s: %v", outputFile.Path, err)
	}
	defer file.Close()

	// Write Anki header lines first
	for _, header := range outputFile.GetAnkiHeaderLines() {
		if _, err := file.WriteString(header + "\n"); err != nil {
			return fmt.Errorf("failed to write Anki header: %v", err)
		}
	}

	// Write CSV data (no column header row needed - Anki metadata provides column info)
	csvWriter := csv.NewWriter(file)
	defer csvWriter.Flush()

	// Note: We don't write column headers because:
	// - When keepHeader=false: The #columns: metadata line tells Anki the column structure
	// - When keepHeader=true: The original header is preserved as the first data row

	// Write data records (when keepHeader=true, first record is the original header)
	records := outputFile.GetCSVRecords()
	for _, record := range records {
		if err := csvWriter.Write(record); err != nil {
			return fmt.Errorf("failed to write data record: %v", err)
		}
	}

	return nil
}

// generateAnkiColumnHeaders generates standard Anki column names based on count
func (f *AnkiFormatter) generateAnkiColumnHeaders(columnCount int) []string {
	if columnCount == 0 {
		return []string{}
	}

	if columnCount == 1 {
		return []string{"Front"}
	}

	if columnCount == 2 {
		return []string{"Front", "Back"}
	}

	// For 3 or more columns, use Front, Back, Extra1, Extra2, etc.
	headers := []string{"Front", "Back"}
	for i := 3; i <= columnCount; i++ {
		headers = append(headers, fmt.Sprintf("Extra%d", i-2))
	}

	return headers
}
