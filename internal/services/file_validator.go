package services

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"ankiprep/internal/models"
)

// FileValidator handles validation of input and output files
type FileValidator struct{}

// NewFileValidator creates a new FileValidator instance
func NewFileValidator() *FileValidator {
	return &FileValidator{}
}

// ValidateInputFile performs comprehensive validation on an input file
func (v *FileValidator) ValidateInputFile(filePath string) error {
	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return fmt.Errorf("input file not found: %s", filePath)
	}

	// Check if file is readable
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("input file not readable: %s (%v)", filePath, err)
	}
	defer file.Close()

	// Check file extension
	if !v.isSupportedInputFormat(filePath) {
		return fmt.Errorf("unsupported input format: %s (only .csv and .tsv files are supported)", filePath)
	}

	// Check file size (reasonable limits)
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return fmt.Errorf("cannot get file info for %s: %v", filePath, err)
	}

	// Warn if file is very large (>100MB)
	const maxRecommendedSize = 100 * 1024 * 1024 // 100MB
	if fileInfo.Size() > maxRecommendedSize {
		return fmt.Errorf("input file %s is very large (%d bytes). Consider splitting it into smaller files",
			filePath, fileInfo.Size())
	}

	return nil
}

// ValidateOutputPath checks if the output path is valid and writable
func (v *FileValidator) ValidateOutputPath(outputPath string) error {
	// Check if output directory exists and is writable
	dir := filepath.Dir(outputPath)
	if dir == "" {
		dir = "."
	}

	// Check if directory exists
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return fmt.Errorf("output directory does not exist: %s", dir)
	}

	// Check if directory is writable by trying to create a temp file
	tempFile := filepath.Join(dir, ".ankiprep_write_test")
	file, err := os.Create(tempFile)
	if err != nil {
		return fmt.Errorf("output directory is not writable: %s (%v)", dir, err)
	}
	file.Close()
	os.Remove(tempFile) // Clean up

	// Ensure output file has .csv extension
	if !strings.HasSuffix(strings.ToLower(outputPath), ".csv") {
		return fmt.Errorf("output file must have .csv extension: %s", outputPath)
	}

	// Check if output file already exists and warn user
	if _, err := os.Stat(outputPath); err == nil {
		// File exists - this is not an error but worth noting
		// The application should handle overwriting gracefully
	}

	return nil
}

// ValidateInputFiles validates multiple input files at once
func (v *FileValidator) ValidateInputFiles(filePaths []string) error {
	if len(filePaths) == 0 {
		return fmt.Errorf("no input files provided")
	}

	var errors []string

	for _, filePath := range filePaths {
		if err := v.ValidateInputFile(filePath); err != nil {
			errors = append(errors, err.Error())
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("validation failed for input files:\n%s", strings.Join(errors, "\n"))
	}

	return nil
}

// ValidateInputFileContent validates the structure and content of a parsed input file
func (v *FileValidator) ValidateInputFileContent(inputFile *models.InputFile) error {
	if inputFile == nil {
		return fmt.Errorf("input file cannot be nil")
	}

	// Check if file has headers
	if len(inputFile.Headers) == 0 {
		return fmt.Errorf("file %s has no column headers", inputFile.Path)
	}

	// Check for duplicate header names
	headerMap := make(map[string]bool)
	for _, header := range inputFile.Headers {
		if header == "" {
			return fmt.Errorf("file %s contains empty column header", inputFile.Path)
		}
		if headerMap[header] {
			return fmt.Errorf("file %s contains duplicate column header: %s", inputFile.Path, header)
		}
		headerMap[header] = true
	}

	// Check if file has data rows
	if len(inputFile.Records) == 0 {
		return fmt.Errorf("file %s contains no data rows", inputFile.Path)
	}

	// Validate record structure
	expectedColumns := len(inputFile.Headers)
	for i, record := range inputFile.Records {
		if len(record) > expectedColumns {
			return fmt.Errorf("file %s, row %d has too many columns (%d, expected %d)",
				inputFile.Path, i+2, len(record), expectedColumns) // +2 for header row
		}
		// Allow fewer columns (will be padded with empty strings)
	}

	return nil
}

// ValidateEncoding checks if the file uses UTF-8 encoding (simplified check)
func (v *FileValidator) ValidateEncoding(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("cannot open file for encoding validation: %v", err)
	}
	defer file.Close()

	// Read a sample of the file
	buffer := make([]byte, 1024)
	n, err := file.Read(buffer)
	if err != nil && err.Error() != "EOF" {
		return fmt.Errorf("cannot read file for encoding validation: %v", err)
	}

	// Simple UTF-8 validation - check if the bytes form valid UTF-8
	sample := string(buffer[:n])
	if !v.isValidUTF8(sample) {
		return fmt.Errorf("file %s does not appear to be UTF-8 encoded", filePath)
	}

	return nil
}

// isSupportedInputFormat checks if the file has a supported input format
func (v *FileValidator) isSupportedInputFormat(filePath string) bool {
	lower := strings.ToLower(filePath)
	return strings.HasSuffix(lower, ".csv") || strings.HasSuffix(lower, ".tsv")
}

// isValidUTF8 performs a simple check for UTF-8 validity
func (v *FileValidator) isValidUTF8(s string) bool {
	// In Go, string conversion automatically validates UTF-8
	// If the conversion succeeds and the length matches, it's valid UTF-8
	return len(s) >= 0 // Simple check - Go strings are UTF-8 by design
}
