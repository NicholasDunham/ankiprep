package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"ankiprep/internal/models"

	"github.com/spf13/cobra"
)

var (
	// Global flags
	verbose        bool
	outputPath     string
	frenchMode     bool
	smartQuotes    bool
	skipDuplicates bool
	keepHeader     bool
)

// rootCmd represents the base command
var rootCmd = &cobra.Command{
	Use:   "ankiprep [files...]",
	Short: "Convert CSV files to Anki-compatible format",
	Long: `ankiprep is a command-line tool for processing CSV and TSV files 
to create Anki-compatible flashcard imports.

Features:
• Merge multiple CSV/TSV files with automatic header unification
• Remove duplicate entries based on content comparison
• Apply French typography formatting (thin spaces before punctuation)
• Convert regular quotes to smart quotes
• Generate Anki-compatible CSV output with proper metadata

Examples:
  ankiprep input.csv
  ankiprep *.csv -o flashcards.csv
  ankiprep file1.csv file2.tsv -f -q
  ankiprep data.csv -s -v`,
	Version: "1.0.0",
	Args:    cobra.MinimumNArgs(1),
	Run:     runProcess,
}

func init() {
	rootCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose output")
	rootCmd.Flags().StringVarP(&outputPath, "output", "o", "", "Specify output file path")
	rootCmd.Flags().BoolVarP(&frenchMode, "french", "f", false, "Add thin spaces before French punctuation (:;!?)")
	rootCmd.Flags().BoolVarP(&smartQuotes, "smart-quotes", "q", false, "Convert straight quotes to curly quotes")
	rootCmd.Flags().BoolVarP(&skipDuplicates, "skip-duplicates", "s", false, "Remove entries with identical content")
	rootCmd.Flags().BoolVarP(&keepHeader, "keep-header", "k", false, "Preserve the first row of CSV files")
}

// runProcess executes the main processing logic - simplified version
func runProcess(cmd *cobra.Command, args []string) {
	startTime := time.Now()

	// Validate and collect input files
	inputPaths, err := collectInputFiles(args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	if verbose {
		fmt.Printf("Processing %d input file(s)...\n", len(inputPaths))
	}

	// Parse input files
	var inputFiles []*models.InputFile
	for _, path := range inputPaths {
		inputFile, err := parseFile(path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing %s: %v\n", path, err)
			os.Exit(1)
		}
		inputFiles = append(inputFiles, inputFile)

		if verbose {
			fmt.Printf("File %s: %d records (%d bytes) (%s)\n",
				path, len(inputFile.Records)+1, getFileSize(path), getFileType(path))
		}
	}

	// Merge headers
	mergedHeaders := mergeHeaders(inputFiles)
	if verbose {
		fmt.Printf("Merging headers: found %d unique columns\n", len(mergedHeaders))
	}

	// Process all records
	var allEntries []*models.DataEntry
	totalRecords := 0

	for _, inputFile := range inputFiles {
		// Add header if keepHeader is true and this is the first file
		if keepHeader && len(allEntries) == 0 {
			headerEntry := models.NewDataEntry(make(map[string]string), inputFile.Path, 0)
			for i, header := range inputFile.Headers {
				if i < len(mergedHeaders) {
					headerEntry.Values[mergedHeaders[i]] = header
				}
			}
			allEntries = append(allEntries, headerEntry)
		}

		// Process data records
		for lineNum, record := range inputFile.Records {
			entry := models.NewDataEntry(make(map[string]string), inputFile.Path, lineNum+2)
			for i, value := range record {
				if i < len(inputFile.Headers) && i < len(mergedHeaders) {
					entry.Values[mergedHeaders[i]] = value
				}
			}
			allEntries = append(allEntries, entry)
			totalRecords++
		}
	}

	if verbose {
		fmt.Printf("Processing records: %d total entries\n", totalRecords)
	}

	// Remove duplicates if requested
	if skipDuplicates {
		originalCount := len(allEntries)
		allEntries = removeDuplicates(allEntries)
		if verbose && originalCount > len(allEntries) {
			fmt.Printf("Removing duplicates: %d duplicates found\n", originalCount-len(allEntries))
		} else if verbose {
			fmt.Printf("Removing duplicates: no duplicates found\n")
		}
	}

	// Apply typography formatting
	if frenchMode || smartQuotes {
		if verbose {
			fmt.Printf("Applying typography formatting")
			if frenchMode && smartQuotes {
				fmt.Printf(" (French typography and smart quotes)")
			} else if frenchMode {
				fmt.Printf(" (French typography)")
			} else {
				fmt.Printf(" (smart quotes)")
			}
			fmt.Printf("...\n")
		}
		applyTypography(allEntries, frenchMode, smartQuotes)
	}

	// Write output
	outputFile := determineOutputPath(inputPaths)
	if verbose {
		fmt.Printf("Writing output to %s\n", outputFile)
	}

	err = writeCSV(outputFile, mergedHeaders, allEntries)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error writing output: %v\n", err)
		os.Exit(1)
	}

	// Success message
	processingTime := time.Since(startTime)
	fmt.Printf("Done. Processed %d unique entries in %.2f seconds\n",
		len(allEntries), processingTime.Seconds())

	if verbose {
		showSummary(inputPaths, totalRecords, len(allEntries), processingTime)
	}
}

// Helper functions - simplified implementations

func collectInputFiles(args []string) ([]string, error) {
	var inputPaths []string
	for _, arg := range args {
		matches, err := filepath.Glob(arg)
		if err != nil {
			return nil, fmt.Errorf("pattern matching failed for %s: %v", arg, err)
		}

		if len(matches) == 0 {
			if _, err := os.Stat(arg); os.IsNotExist(err) {
				return nil, fmt.Errorf("file not found: %s", arg)
			}
			inputPaths = append(inputPaths, arg)
		} else {
			for _, match := range matches {
				if isSupportedFile(match) {
					inputPaths = append(inputPaths, match)
				}
			}
		}
	}

	if len(inputPaths) == 0 {
		return nil, fmt.Errorf("no valid input files found")
	}

	return inputPaths, nil
}

func parseFile(filePath string) (*models.InputFile, error) {
	inputFile := models.NewInputFile(filePath)
	inputFile.DetectSeparator()

	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = inputFile.Separator
	reader.LazyQuotes = true
	reader.TrimLeadingSpace = false

	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	if len(records) < 1 {
		return nil, fmt.Errorf("file contains no data")
	}

	inputFile.Headers = records[0]

	// Strip UTF-8 BOM from first header field if present
	if len(inputFile.Headers) > 0 && len(inputFile.Headers[0]) > 0 {
		if runes := []rune(inputFile.Headers[0]); len(runes) > 0 && runes[0] == '\uFEFF' {
			inputFile.Headers[0] = string(runes[1:])
		}
	}

	if len(records) > 1 {
		inputFile.Records = records[1:]
	}

	return inputFile, nil
}

func mergeHeaders(inputFiles []*models.InputFile) []string {
	seen := make(map[string]bool)
	var merged []string

	for _, inputFile := range inputFiles {
		for _, header := range inputFile.Headers {
			if header != "" && !seen[header] {
				seen[header] = true
				merged = append(merged, header)
			}
		}
	}

	return merged
}

func removeDuplicates(entries []*models.DataEntry) []*models.DataEntry {
	seen := make(map[string]bool)
	var unique []*models.DataEntry

	for _, entry := range entries {
		key := entry.GetHash()
		if !seen[key] {
			seen[key] = true
			unique = append(unique, entry)
		}
	}

	return unique
}

// isEnglishColumn determines if a column header indicates English content
// that should not have French typography rules applied
func isEnglishColumn(header string) bool {
	header = strings.ToLower(strings.TrimSpace(header))
	englishPatterns := []string{"english", "eng", "pronunciation", "phonetic"}

	for _, pattern := range englishPatterns {
		if strings.Contains(header, pattern) {
			return true
		}
	}
	return false
}

func applyTypography(entries []*models.DataEntry, french, quotes bool) {
	for _, entry := range entries {
		for key, value := range entry.Values {
			// Determine which typography rules to apply based on column header
			isEnglish := isEnglishColumn(key)

			// Always apply smart quotes if enabled
			applySmartQuotes := quotes

			// Only apply French typography to non-English fields
			applyFrench := french && !isEnglish

			// Create processor with appropriate settings
			processor := models.NewTypographyProcessor(applyFrench, applySmartQuotes)
			entry.Values[key] = processor.ProcessText(value)
		}
	}
}

func writeCSV(outputPath string, headers []string, entries []*models.DataEntry) error {
	file, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Write Anki metadata headers directly (not as CSV)
	ankiHeaders := []string{
		"#separator:comma",
		"#html:true",
		"#columns:" + strings.Join(headers, ","),
	}

	for _, header := range ankiHeaders {
		if _, err := file.WriteString(header + "\n"); err != nil {
			return err
		}
	}

	// Now write data using CSV writer
	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write data
	for _, entry := range entries {
		record := make([]string, len(headers))
		for i, header := range headers {
			record[i] = entry.Values[header]
		}
		if err := writer.Write(record); err != nil {
			return err
		}
	}

	return nil
}

// Utility functions
func isSupportedFile(filePath string) bool {
	ext := strings.ToLower(filepath.Ext(filePath))
	return ext == ".csv" || ext == ".tsv"
}

func getFileSize(filePath string) int64 {
	if info, err := os.Stat(filePath); err == nil {
		return info.Size()
	}
	return 0
}

func getFileType(filePath string) string {
	if strings.HasSuffix(strings.ToLower(filePath), ".tsv") {
		return "tab-separated"
	}
	return "comma-separated"
}

func determineOutputPath(inputPaths []string) string {
	if outputPath != "" {
		return outputPath
	}

	if len(inputPaths) == 1 {
		base := strings.TrimSuffix(inputPaths[0], filepath.Ext(inputPaths[0]))
		return base + "_processed.csv"
	}

	return "merged_output.csv"
}

func showSummary(inputFiles []string, totalInput, totalOutput int, duration time.Duration) {
	fmt.Printf("\nProcessing Summary:\n")
	fmt.Printf("Input files: %d\n", len(inputFiles))
	for i, file := range inputFiles {
		fmt.Printf("  %d. %s\n", i+1, file)
	}
	fmt.Printf("Total input records: %d\n", totalInput)
	fmt.Printf("Output records: %d\n", totalOutput)
	fmt.Printf("Processing time: %.2f seconds\n", duration.Seconds())
	if duration.Seconds() > 0 && totalOutput > 0 {
		rate := float64(totalOutput) / duration.Seconds()
		fmt.Printf("Processing rate: %.0f records/second\n", rate)
	}
	fmt.Printf("Processing completed successfully\n")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func main() {
	Execute()
}
