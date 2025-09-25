package services

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"

	"ankiprep/internal/models"
)

// CSVParser handles parsing CSV and TSV files
type CSVParser struct{}

// NewCSVParser creates a new CSVParser instance
func NewCSVParser() *CSVParser {
	return &CSVParser{}
}

// ParseFile parses a CSV/TSV file and returns an InputFile with populated data
func (p *CSVParser) ParseFile(filePath string) (*models.InputFile, error) {
	inputFile := models.NewInputFile(filePath)

	// Detect separator based on file extension
	inputFile.DetectSeparator()

	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %s: %v", filePath, err)
	}
	defer file.Close()

	// Create CSV reader with detected separator
	reader := csv.NewReader(file)
	reader.Comma = inputFile.Separator
	reader.LazyQuotes = true        // Allow lazy quotes for better compatibility
	reader.TrimLeadingSpace = false // Preserve leading spaces

	// Read all records
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to parse CSV file %s: %v", filePath, err)
	}

	if len(records) < 1 {
		return nil, fmt.Errorf("file %s contains no data", filePath)
	}

	// First row is headers
	inputFile.Headers = records[0]

	// Remaining rows are data
	if len(records) > 1 {
		inputFile.Records = records[1:]
	} else {
		return nil, fmt.Errorf("file %s contains no data rows", filePath)
	}

	return inputFile, nil
}

// ParseToDataEntries converts an InputFile to a slice of DataEntry objects
func (p *CSVParser) ParseToDataEntries(inputFile *models.InputFile) ([]*models.DataEntry, error) {
	if inputFile == nil {
		return nil, fmt.Errorf("input file cannot be nil")
	}

	var entries []*models.DataEntry

	for lineNumber, record := range inputFile.Records {
		// Skip empty rows
		if p.isEmptyRecord(record) {
			continue
		}

		// Create values map
		values := make(map[string]string)

		for i, value := range record {
			if i < len(inputFile.Headers) {
				values[inputFile.Headers[i]] = value
			}
		}

		// Create DataEntry
		entry := models.NewDataEntry(values, inputFile.Path, lineNumber+2) // +2 because line 1 is headers
		entries = append(entries, entry)
	}

	return entries, nil
}

// ValidateFormat checks if the file format is valid CSV/TSV
func (p *CSVParser) ValidateFormat(filePath string) error {
	// Check file extension
	if !p.isSupportedExtension(filePath) {
		return fmt.Errorf("unsupported file format: only .csv and .tsv files are supported")
	}

	// Try to open and read at least the header
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("cannot access file %s: %v", filePath, err)
	}
	defer file.Close()

	// Create a temporary InputFile to detect separator
	tempFile := models.NewInputFile(filePath)
	tempFile.DetectSeparator()

	// Try to read the first few lines
	reader := csv.NewReader(file)
	reader.Comma = tempFile.Separator
	reader.LazyQuotes = true

	// Read header
	header, err := reader.Read()
	if err != nil {
		if err == io.EOF {
			return fmt.Errorf("file %s is empty", filePath)
		}
		return fmt.Errorf("invalid CSV/TSV format in %s: %v", filePath, err)
	}

	if len(header) == 0 {
		return fmt.Errorf("file %s has no columns", filePath)
	}

	return nil
}

// isSupportedExtension checks if the file has a supported extension
func (p *CSVParser) isSupportedExtension(filePath string) bool {
	lower := strings.ToLower(filePath)
	return strings.HasSuffix(lower, ".csv") || strings.HasSuffix(lower, ".tsv")
}

// isEmptyRecord checks if a record contains only empty values
func (p *CSVParser) isEmptyRecord(record []string) bool {
	for _, value := range record {
		if strings.TrimSpace(value) != "" {
			return false
		}
	}
	return true
}
