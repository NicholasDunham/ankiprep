package app

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"

	"ankiprep/internal/models"
	"ankiprep/internal/services"
)

// Processor handles the main processing logic for CSV to Anki conversion
type Processor struct {
	// Service dependencies
	fileValidator     *services.FileValidator
	csvParser         *services.CSVParser
	columnMerger      *services.ColumnMerger
	duplicateDetector *services.DuplicateDetector
	typographyService *services.TypographyService
	ankiFormatter     *services.AnkiFormatter
	progressReporter  *services.ProgressReporter
	fileService       *services.FileService
	memoryMonitor     *services.MemoryMonitor

	// Configuration
	verbose        bool
	outputPath     string
	frenchMode     bool
	smartQuotes    bool
	skipDuplicates bool

	// Cleanup tracking
	partialFiles []string
}

// ProcessorConfig holds configuration options for the processor
type ProcessorConfig struct {
	OutputPath     string
	FrenchMode     bool
	SmartQuotes    bool
	SkipDuplicates bool
	Verbose        bool
}

// ProcessingOptions holds runtime options that can vary per processing request
type ProcessingOptions struct {
	KeepHeader bool // Whether to keep the original CSV header row
}

// NewProcessor creates a new processor instance with all required services
func NewProcessor(config ProcessorConfig) *Processor {
	progressReporter := services.NewProgressReporter(config.Verbose)

	return &Processor{
		fileValidator:     services.NewFileValidator(),
		csvParser:         services.NewCSVParser(),
		columnMerger:      services.NewColumnMerger(),
		duplicateDetector: services.NewDuplicateDetector(),
		typographyService: services.NewTypographyServiceLegacy(config.FrenchMode, config.SmartQuotes),
		ankiFormatter:     services.NewAnkiFormatter(),
		progressReporter:  progressReporter,
		fileService:       services.NewFileService(),
		memoryMonitor:     services.NewMemoryMonitor(),
		verbose:           config.Verbose,
		outputPath:        config.OutputPath,
		frenchMode:        config.FrenchMode,
		smartQuotes:       config.SmartQuotes,
		skipDuplicates:    config.SkipDuplicates,
		partialFiles:      make([]string, 0),
	}
}

// ProcessFiles processes one or more CSV files and produces an Anki-compatible CSV output
func (p *Processor) ProcessFiles(inputPaths []string, options ProcessingOptions) (*models.ProcessingReport, error) {
	// Ensure cleanup on any failure
	defer func() {
		if err := recover(); err != nil {
			p.cleanupPartialFiles()
			p.memoryMonitor.ResetOptimizations()
			panic(err) // Re-panic after cleanup
		}
	}()

	// Cleanup memory optimizations when done
	defer p.memoryMonitor.ResetOptimizations()

	// Start progress reporting and detect large files
	p.progressReporter.StartProcessing(len(inputPaths))
	p.progressReporter.DetectLargeFiles(inputPaths)

	// Check if we detected large files and optimize memory accordingly
	if p.progressReporter.IsLargeFileMode() {
		// Estimate memory needed for large files (rough calculation)
		fileSize := p.progressReporter.GetTotalFileSize()
		estimatedMemory := p.memoryMonitor.EstimateMemoryRequirement(int(fileSize/100), 100)
		p.memoryMonitor.OptimizeForLargeDataset(estimatedMemory)

		if p.verbose {
			p.progressReporter.ReportWarning(fmt.Sprintf("Large files detected, optimizing for memory usage (estimated: %s)", p.formatBytes(estimatedMemory)))
		}
	} else {
		// Enable monitoring for smaller datasets too
		p.memoryMonitor.Enable()
	}

	// Parse all CSV files
	var allInputFiles []*models.InputFile
	var allEntries []*models.DataEntry
	var totalRecords int

	p.progressReporter.StartOperation("File parsing")
	for _, path := range inputPaths {
		// Validate input file
		if err := p.fileValidator.ValidateInputFile(path); err != nil {
			p.cleanupPartialFiles()
			return nil, fmt.Errorf("validation failed for %s: %w", path, err)
		}

		// Parse the file
		inputFile, err := p.csvParser.ParseFile(path)
		if err != nil {
			p.cleanupPartialFiles()
			return nil, fmt.Errorf("failed to parse %s: %w", path, err)
		}

		// Handle header removal/preservation based on options
		p.handleHeaderProcessing(inputFile, options)

		// Convert to DataEntry objects using the service
		entries, err := p.csvParser.ParseToDataEntries(inputFile)
		if err != nil {
			p.cleanupPartialFiles()
			return nil, fmt.Errorf("failed to convert to data entries for %s: %w", path, err)
		}

		// Report file info
		separatorStr := "comma"
		if inputFile.Separator == '\t' {
			separatorStr = "tab"
		}
		p.progressReporter.ReportFileInfo(inputFile.Path, len(inputFile.Records), separatorStr)

		allInputFiles = append(allInputFiles, inputFile)
		allEntries = append(allEntries, entries...)
		totalRecords += len(entries)

		// Log memory usage after processing each file
		p.memoryMonitor.LogMemoryUsage(p.progressReporter)

		// Trigger GC if memory pressure is high
		if p.memoryMonitor.TriggerGCIfNeeded() && p.verbose {
			p.progressReporter.ReportWarning("Memory pressure detected, triggered garbage collection")
		}
	}
	p.progressReporter.EndOperation("File parsing", len(allInputFiles))

	// Report total record count
	p.progressReporter.ReportRecordProcessing(totalRecords)

	// Merge column headers to create unified structure
	p.progressReporter.StartOperation("Header merging")
	mergedHeaders, err := p.columnMerger.MergeColumns(allInputFiles)
	if err != nil {
		p.cleanupPartialFiles()
		return nil, fmt.Errorf("failed to merge columns: %w", err)
	}
	p.progressReporter.ReportHeaderMerging(len(mergedHeaders))
	p.progressReporter.EndOperation("Header merging", len(mergedHeaders))

	// Normalize entries to use merged headers
	p.progressReporter.StartOperation("Data normalization")
	normalizedEntries := p.columnMerger.NormalizeDataEntries(allEntries, mergedHeaders)
	p.progressReporter.EndOperation("Data normalization", len(normalizedEntries))

	// Remove duplicates if requested
	var uniqueEntries []*models.DataEntry
	var duplicateCount int

	if p.skipDuplicates {
		p.progressReporter.StartOperation("Duplicate detection")
		uniqueEntries, duplicateCount = p.duplicateDetector.DetectDuplicates(normalizedEntries)
		p.progressReporter.ReportDuplicateRemoval(duplicateCount)
		p.progressReporter.EndOperation("Duplicate detection", len(normalizedEntries))
	} else {
		uniqueEntries = normalizedEntries
	}

	// Apply typography processing
	if p.frenchMode || p.smartQuotes {
		p.progressReporter.StartOperation("Typography processing")
		p.progressReporter.ReportTypographyProcessing(p.frenchMode, p.smartQuotes)
		uniqueEntries = p.typographyService.ProcessEntries(uniqueEntries)
		p.progressReporter.EndOperation("Typography processing", len(uniqueEntries))
	}

	// Create output file using the Anki formatter
	outputPath := p.outputPath
	if outputPath == "" {
		outputPath = p.generateDefaultOutputPath(inputPaths)
	}

	// Track the output file for cleanup in case of failure
	p.trackPartialFile(outputPath)

	// Report file writing
	p.progressReporter.StartOperation("Output file creation")
	p.progressReporter.ReportFileWriting(outputPath)

	outputFile, err := p.ankiFormatter.FormatForAnki(uniqueEntries, mergedHeaders, outputPath)
	if err != nil {
		p.cleanupPartialFiles()
		return nil, fmt.Errorf("failed to format for Anki: %w", err)
	}

	// Write output file
	if err := p.ankiFormatter.WriteToFileWithOptions(outputFile, options.KeepHeader); err != nil {
		p.cleanupPartialFiles()
		return nil, fmt.Errorf("failed to write output file: %w", err)
	}
	p.progressReporter.EndOperation("Output file creation", len(uniqueEntries))

	// Clear tracked files since operation was successful
	p.clearPartialFiles()

	// Show memory usage for large operations
	if p.progressReporter.IsLargeFileMode() {
		p.progressReporter.ShowMemoryUsage()
	}

	// Create processing report
	inputFilePathList := make([]string, len(allInputFiles))
	for i, file := range allInputFiles {
		inputFilePathList[i] = file.Path
	}

	report := models.NewProcessingReport()
	report.InputFiles = inputFilePathList
	report.TotalInputRecords = totalRecords
	report.DuplicatesRemoved = duplicateCount
	report.OutputRecords = len(uniqueEntries)
	report.ProcessingTime = p.progressReporter.GetElapsedTime()

	// Report completion
	p.progressReporter.ReportCompletion(len(uniqueEntries), report.ProcessingTime)

	return report, nil
}

// ProcessSingleFile processes a single CSV file (convenience method)
func (p *Processor) ProcessSingleFile(inputPath string, options ProcessingOptions) (*models.ProcessingReport, error) {
	return p.ProcessFiles([]string{inputPath}, options)
}

// SetVerbose enables or disables verbose output
func (p *Processor) SetVerbose(verbose bool) {
	p.verbose = verbose
	p.progressReporter.SetVerbose(verbose)
}

// SetOutputPath sets the output file path
func (p *Processor) SetOutputPath(path string) {
	p.outputPath = path
}

// SetTypographyOptions configures typography processing
func (p *Processor) SetTypographyOptions(frenchMode, smartQuotes bool) {
	p.frenchMode = frenchMode
	p.smartQuotes = smartQuotes
	p.typographyService = services.NewTypographyServiceLegacy(frenchMode, smartQuotes)
}

// SetSkipDuplicates enables or disables duplicate removal
func (p *Processor) SetSkipDuplicates(skip bool) {
	p.skipDuplicates = skip
}

// generateDefaultOutputPath creates a default output path based on input files
func (p *Processor) generateDefaultOutputPath(inputPaths []string) string {
	if len(inputPaths) == 1 {
		// Single file: same directory, replace extension with .anki.csv
		inputPath := inputPaths[0]
		dir := filepath.Dir(inputPath)
		base := filepath.Base(inputPath)

		// Remove existing extension and add .anki.csv
		ext := filepath.Ext(base)
		name := base[:len(base)-len(ext)]

		return filepath.Join(dir, name+".anki.csv")
	}

	// Multiple files: use working directory, create merged filename
	cwd, err := os.Getwd()
	if err != nil {
		cwd = "."
	}

	return filepath.Join(cwd, "merged.anki.csv")
}

// writeOutputFile writes the processed data to the output CSV file
func (p *Processor) writeOutputFile(outputFile *models.OutputFile) error {
	file, err := os.Create(outputFile.Path)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write headers
	if err := writer.Write(outputFile.Headers); err != nil {
		return fmt.Errorf("failed to write headers: %w", err)
	}

	// Write data entries
	for i, entry := range outputFile.Records {
		record := make([]string, len(outputFile.Headers))

		// Map entry data to header positions
		for j, header := range outputFile.Headers {
			record[j] = entry.GetValue(header)
		}

		if err := writer.Write(record); err != nil {
			return fmt.Errorf("failed to write record %d: %w", i+1, err)
		}

		// Progress reporting for large files
		if p.verbose && len(outputFile.Records) > 1000 {
			p.progressReporter.ReportProgress(i+1, len(outputFile.Records), "Writing output")
		}
	}

	return nil
}

// ValidateConfiguration validates the processor configuration
func (p *Processor) ValidateConfiguration() error {
	// Check output path is writable if specified
	if p.outputPath != "" {
		dir := filepath.Dir(p.outputPath)
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			return fmt.Errorf("output directory does not exist: %s", dir)
		}
	}

	return nil
}

// trackPartialFile adds a file path to the cleanup list
func (p *Processor) trackPartialFile(filePath string) {
	p.partialFiles = append(p.partialFiles, filePath)
}

// clearPartialFiles clears the cleanup list (called on successful completion)
func (p *Processor) clearPartialFiles() {
	p.partialFiles = make([]string, 0)
}

// cleanupPartialFiles removes all partial files tracked during processing
func (p *Processor) cleanupPartialFiles() {
	if len(p.partialFiles) == 0 {
		return
	}

	if p.verbose {
		fmt.Printf("Cleaning up %d partial file(s)...\n", len(p.partialFiles))
	}

	cleaned := 0
	for _, filePath := range p.partialFiles {
		if _, err := os.Stat(filePath); err == nil {
			// File exists, try to remove it
			if err := p.fileService.SafeRemoveFile(filePath); err != nil {
				p.progressReporter.ReportWarning(fmt.Sprintf("Failed to cleanup partial file %s: %v", filePath, err))
			} else {
				cleaned++
				if p.verbose {
					fmt.Printf("Removed partial file: %s\n", filePath)
				}
			}
		}
	}

	if p.verbose && cleaned > 0 {
		fmt.Printf("Cleaned up %d partial file(s)\n", cleaned)
	}

	// Clear the list
	p.clearPartialFiles()
}

// Cleanup performs final cleanup of resources and temporary files
func (p *Processor) Cleanup() {
	// Clean up any remaining partial files
	p.cleanupPartialFiles()

	// Clean up temporary files created by the file service
	if err := p.fileService.CleanupTempFiles(); err != nil {
		p.progressReporter.ReportWarning(fmt.Sprintf("Failed to cleanup temp files: %v", err))
	}
}

// formatBytes formats byte count in human-readable format
func (p *Processor) formatBytes(bytes uint64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}

	div, exp := uint64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}

	units := []string{"KB", "MB", "GB", "TB"}
	return fmt.Sprintf("%.1f %s", float64(bytes)/float64(div), units[exp])
}

// handleHeaderProcessing handles header removal or preservation based on options
func (p *Processor) handleHeaderProcessing(inputFile *models.InputFile, options ProcessingOptions) {
	if !options.KeepHeader {
		// Default behavior: header already removed by CSV parser, nothing to do
		return
	}

	// Keep header: prepend the header row back to records as a data row
	if len(inputFile.Headers) > 0 {
		// Create a new slice with header as first row, followed by existing records
		allRecords := make([][]string, 0, len(inputFile.Records)+1)
		allRecords = append(allRecords, inputFile.Headers)
		allRecords = append(allRecords, inputFile.Records...)
		inputFile.Records = allRecords

		// Keep headers for column merging but mark that this file should preserve header
		// We'll use this information later in the output formatting
	}
}
