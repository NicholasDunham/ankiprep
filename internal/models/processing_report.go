package models

import (
	"fmt"
	"time"
)

// ProcessingReport contains summary of processing actions and statistics
type ProcessingReport struct {
	InputFiles        []string      // List of processed input file paths
	TotalInputRecords int           // Count of records before deduplication
	DuplicatesRemoved int           // Count of duplicate records removed
	OutputRecords     int           // Final count of records in output
	ProcessingTime    time.Duration // Total processing time
	Errors            []string      // List of any processing errors
}

// NewProcessingReport creates a new ProcessingReport instance
func NewProcessingReport() *ProcessingReport {
	return &ProcessingReport{
		InputFiles:        []string{},
		TotalInputRecords: 0,
		DuplicatesRemoved: 0,
		OutputRecords:     0,
		ProcessingTime:    0,
		Errors:            []string{},
	}
}

// Validate checks if the processing report meets all validation requirements
func (r *ProcessingReport) Validate() error {
	// TotalInputRecords >= OutputRecords
	if r.TotalInputRecords < r.OutputRecords {
		return fmt.Errorf("total input records (%d) cannot be less than output records (%d)",
			r.TotalInputRecords, r.OutputRecords)
	}

	// DuplicatesRemoved = TotalInputRecords - OutputRecords
	expectedDuplicates := r.TotalInputRecords - r.OutputRecords
	if r.DuplicatesRemoved != expectedDuplicates {
		return fmt.Errorf("duplicates removed (%d) does not match expected (%d)",
			r.DuplicatesRemoved, expectedDuplicates)
	}

	// ProcessingTime must be positive duration
	if r.ProcessingTime < 0 {
		return fmt.Errorf("processing time cannot be negative")
	}

	return nil
}

// AddInputFile adds an input file path to the report
func (r *ProcessingReport) AddInputFile(path string) {
	r.InputFiles = append(r.InputFiles, path)
}

// AddError adds an error message to the report
func (r *ProcessingReport) AddError(err error) {
	if err != nil {
		r.Errors = append(r.Errors, err.Error())
	}
}

// AddErrorString adds an error message string to the report
func (r *ProcessingReport) AddErrorString(message string) {
	r.Errors = append(r.Errors, message)
}

// SetCounts sets the record counts in the report
func (r *ProcessingReport) SetCounts(totalInput, duplicates, output int) {
	r.TotalInputRecords = totalInput
	r.DuplicatesRemoved = duplicates
	r.OutputRecords = output
}

// SetProcessingTime sets the processing time
func (r *ProcessingReport) SetProcessingTime(duration time.Duration) {
	r.ProcessingTime = duration
}

// HasErrors returns true if the report contains any errors
func (r *ProcessingReport) HasErrors() bool {
	return len(r.Errors) > 0
}

// GetInputFileCount returns the number of input files processed
func (r *ProcessingReport) GetInputFileCount() int {
	return len(r.InputFiles)
}

// GetDuplicationRate returns the percentage of records that were duplicates
func (r *ProcessingReport) GetDuplicationRate() float64 {
	if r.TotalInputRecords == 0 {
		return 0.0
	}
	return float64(r.DuplicatesRemoved) / float64(r.TotalInputRecords) * 100.0
}

// GetSummaryString returns a formatted summary of the processing report
func (r *ProcessingReport) GetSummaryString() string {
	if r.HasErrors() {
		return fmt.Sprintf("Processing failed with %d error(s)", len(r.Errors))
	}

	return fmt.Sprintf("Processed %d unique entries from %d file(s) in %.2f seconds (removed %d duplicates)",
		r.OutputRecords, len(r.InputFiles), r.ProcessingTime.Seconds(), r.DuplicatesRemoved)
}
