package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"ankiprep/internal/app"
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
	keepHeader     bool // New flag for preserving CSV headers

	// Error formatter
	errorFormatter *ErrorFormatter
)

// rootCmd represents the base command when called without any subcommands
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
	// Root command flags
	rootCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose output")
	rootCmd.Flags().StringVarP(&outputPath, "output", "o", "", "Specify output file path. If not provided, generates based on input filename(s)")
	rootCmd.Flags().BoolVarP(&frenchMode, "french", "f", false, "Add thin spaces before French punctuation (:;!?)")
	rootCmd.Flags().BoolVarP(&smartQuotes, "smart-quotes", "q", false, "Convert straight quotes to curly quotes")
	rootCmd.Flags().BoolVarP(&skipDuplicates, "skip-duplicates", "s", false, "Remove entries with identical content")
	rootCmd.Flags().BoolVarP(&keepHeader, "keep-header", "k", false, "Preserve the first row of CSV files (default: remove header)")
}

// runProcess executes the main processing logic
func runProcess(cmd *cobra.Command, args []string) {
	// Initialize error formatter with verbose setting
	errorFormatter = NewErrorFormatter(verbose)

	// Set up panic recovery
	defer errorFormatter.HandlePanic()

	// Validate input files exist
	var inputPaths []string
	for _, arg := range args {
		// Handle glob patterns
		matches, err := filepath.Glob(arg)
		if err != nil {
			errorFormatter.ExitWithError(err, "pattern matching")
		}

		if len(matches) == 0 {
			// No glob match, check if file exists directly
			if _, err := os.Stat(arg); os.IsNotExist(err) {
				errorFormatter.ExitWithError(fmt.Errorf("file not found: %s", arg), "input validation")
			}
			inputPaths = append(inputPaths, arg)
		} else {
			// Add all matches, filtering for supported extensions
			validMatches := 0
			for _, match := range matches {
				if isSupportedFile(match) {
					inputPaths = append(inputPaths, match)
					validMatches++
				}
			}

			if validMatches == 0 {
				errorFormatter.PrintWarning(fmt.Sprintf("No supported files found for pattern: %s", arg))
			}
		}
	}

	if len(inputPaths) == 0 {
		errorFormatter.ExitWithError(fmt.Errorf("no valid input files found"), "input validation")
	}

	// Create processor configuration
	config := app.ProcessorConfig{
		OutputPath:     outputPath,
		FrenchMode:     frenchMode,
		SmartQuotes:    smartQuotes,
		SkipDuplicates: skipDuplicates,
		Verbose:        verbose,
	}

	// Create and configure processor
	processor := app.NewProcessor(config)

	// Set up cleanup on exit
	defer processor.Cleanup()

	// Validate configuration
	if err := processor.ValidateConfiguration(); err != nil {
		errorFormatter.ExitWithError(err, "configuration validation")
	}

	// Process files
	options := app.ProcessingOptions{
		KeepHeader: keepHeader,
	}
	report, err := processor.ProcessFiles(inputPaths, options)
	if err != nil {
		errorFormatter.ExitWithError(err, "file processing")
	}

	// Display summary if verbose or multiple files
	if verbose || len(inputPaths) > 1 {
		displaySummary(report)
	}

	// Success exit
	if verbose {
		errorFormatter.ExitWithSuccess("Processing completed successfully")
	}
}

// displaySummary shows processing summary information
func displaySummary(report *models.ProcessingReport) {
	if report == nil {
		errorFormatter.PrintWarning("Unable to generate processing summary - report is nil")
		return
	}

	fmt.Printf("\nProcessing Summary:\n")
	fmt.Printf("Input files: %d\n", len(report.InputFiles))

	for i, file := range report.InputFiles {
		fmt.Printf("  %d. %s\n", i+1, file)
	}

	fmt.Printf("Total input records: %s\n", formatNumber(report.TotalInputRecords))
	if report.DuplicatesRemoved > 0 {
		fmt.Printf("Duplicates removed: %s\n", formatNumber(report.DuplicatesRemoved))
	}
	fmt.Printf("Output records: %s\n", formatNumber(report.OutputRecords))
	fmt.Printf("Processing time: %.2f seconds\n", report.ProcessingTime.Seconds())

	// Show efficiency metrics
	if report.ProcessingTime.Seconds() > 0 && report.OutputRecords > 0 {
		rate := float64(report.OutputRecords) / report.ProcessingTime.Seconds()
		fmt.Printf("Processing rate: %.0f records/second\n", rate)
	}
}

// isSupportedFile checks if a file has a supported extension
func isSupportedFile(filePath string) bool {
	ext := strings.ToLower(filepath.Ext(filePath))
	return ext == ".csv" || ext == ".tsv"
}

// formatNumber formats a number with commas for thousands
func formatNumber(n int) string {
	if n < 1000 {
		return fmt.Sprintf("%d", n)
	}

	str := fmt.Sprintf("%d", n)
	result := ""

	for i, digit := range str {
		if i > 0 && (len(str)-i)%3 == 0 {
			result += ","
		}
		result += string(digit)
	}

	return result
}

// Execute adds all child commands to the root command and sets flags appropriately
func Execute() {
	// Initialize error formatter for early errors
	errorFormatter = NewErrorFormatter(false) // Will be updated with verbose setting later
	defer errorFormatter.HandlePanic()

	if err := rootCmd.Execute(); err != nil {
		errorFormatter.ExitWithError(err, "command execution")
	}
}

func main() {
	Execute()
}
