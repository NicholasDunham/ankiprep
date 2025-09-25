package main

import (
	"fmt"
	"os"
	"strings"

	"ankiprep/internal/services"
)

// Exit codes for different types of errors
const (
	ExitSuccess            = 0   // Successful execution
	ExitGeneralError       = 1   // General errors
	ExitInputError         = 2   // Input file errors (not found, invalid format)
	ExitOutputError        = 3   // Output file errors (permission, disk space)
	ExitProcessingError    = 4   // Processing errors (parsing, validation)
	ExitConfigurationError = 5   // Configuration errors (invalid flags)
	ExitResourceError      = 6   // Resource errors (memory, disk space)
	ExitInterrupted        = 130 // Process interrupted (Ctrl+C)
)

// ErrorFormatter provides standardized error message formatting
type ErrorFormatter struct {
	verbose bool
}

// NewErrorFormatter creates a new error formatter
func NewErrorFormatter(verbose bool) *ErrorFormatter {
	return &ErrorFormatter{verbose: verbose}
}

// FormatError formats an error message with context and suggestions
func (ef *ErrorFormatter) FormatError(err error, context string) (string, int) {
	if err == nil {
		return "", ExitSuccess
	}

	// Check for specific error types and provide appropriate exit codes
	exitCode := ef.determineExitCode(err)

	// Format the error message
	var message strings.Builder
	message.WriteString(fmt.Sprintf("Error: %v", err))

	// Add context if provided
	if context != "" {
		message.WriteString(fmt.Sprintf(" (%s)", context))
	}

	// Add detailed context for verbose mode
	if ef.verbose {
		if fileErr, isFileError := services.IsFileError(err); isFileError {
			message.WriteString(fmt.Sprintf("\nDetails: %s", services.GetErrorContext(fileErr)))
		}
	}

	// Add suggestions based on error type
	suggestion := ef.getSuggestion(err, exitCode)
	if suggestion != "" {
		message.WriteString(fmt.Sprintf("\nSuggestion: %s", suggestion))
	}

	return message.String(), exitCode
}

// FormatWarning formats a warning message
func (ef *ErrorFormatter) FormatWarning(message string) string {
	return fmt.Sprintf("Warning: %s", message)
}

// FormatSuccess formats a success message
func (ef *ErrorFormatter) FormatSuccess(message string) string {
	return message
}

// determineExitCode determines the appropriate exit code based on error type
func (ef *ErrorFormatter) determineExitCode(err error) int {
	errStr := strings.ToLower(err.Error())

	// File-related errors
	if strings.Contains(errStr, "not found") ||
		strings.Contains(errStr, "no such file") ||
		strings.Contains(errStr, "invalid pattern") {
		return ExitInputError
	}

	// Permission errors
	if strings.Contains(errStr, "permission") ||
		strings.Contains(errStr, "access") {
		return ExitOutputError
	}

	// Disk space errors
	if strings.Contains(errStr, "disk") ||
		strings.Contains(errStr, "space") ||
		strings.Contains(errStr, "no space left") {
		return ExitResourceError
	}

	// Processing errors
	if strings.Contains(errStr, "validation failed") ||
		strings.Contains(errStr, "failed to parse") ||
		strings.Contains(errStr, "failed to convert") ||
		strings.Contains(errStr, "failed to format") {
		return ExitProcessingError
	}

	// Configuration errors
	if strings.Contains(errStr, "invalid") && strings.Contains(errStr, "flag") ||
		strings.Contains(errStr, "configuration") {
		return ExitConfigurationError
	}

	// Memory errors
	if strings.Contains(errStr, "memory") ||
		strings.Contains(errStr, "out of memory") {
		return ExitResourceError
	}

	// Default to general error
	return ExitGeneralError
}

// getSuggestion provides helpful suggestions based on error type
func (ef *ErrorFormatter) getSuggestion(err error, exitCode int) string {
	switch exitCode {
	case ExitInputError:
		return "Check that input files exist and have the correct extensions (.csv or .tsv)"

	case ExitOutputError:
		return "Verify write permissions for the output directory and ensure it exists"

	case ExitProcessingError:
		errStr := strings.ToLower(err.Error())
		if strings.Contains(errStr, "validation") {
			return "Ensure input files are properly formatted CSV/TSV with valid headers"
		}
		if strings.Contains(errStr, "parse") {
			return "Check file encoding (UTF-8 required) and CSV/TSV format"
		}
		if strings.Contains(errStr, "format") {
			return "Verify data compatibility for Anki import format"
		}
		return "Review input file format and content"

	case ExitConfigurationError:
		return "Review command line flags and their values"

	case ExitResourceError:
		return "Free up disk space or memory, or try processing smaller files"

	default:
		return "Run with --verbose flag for more detailed error information"
	}
}

// HandlePanic provides graceful panic recovery with appropriate exit codes
func (ef *ErrorFormatter) HandlePanic() {
	if r := recover(); r != nil {
		message := fmt.Sprintf("Internal error: %v", r)
		if ef.verbose {
			message += "\nThis appears to be an unexpected error. Please report this issue."
		}
		fmt.Fprintf(os.Stderr, "%s\n", message)
		os.Exit(ExitGeneralError)
	}
}

// ExitWithError prints an error message and exits with appropriate code
func (ef *ErrorFormatter) ExitWithError(err error, context string) {
	message, exitCode := ef.FormatError(err, context)
	fmt.Fprintf(os.Stderr, "%s\n", message)
	os.Exit(exitCode)
}

// ExitWithSuccess prints a success message and exits with code 0
func (ef *ErrorFormatter) ExitWithSuccess(message string) {
	if message != "" {
		fmt.Printf("%s\n", ef.FormatSuccess(message))
	}
	os.Exit(ExitSuccess)
}

// PrintWarning prints a formatted warning message
func (ef *ErrorFormatter) PrintWarning(message string) {
	fmt.Fprintf(os.Stderr, "%s\n", ef.FormatWarning(message))
}

// GetExitCodeName returns a human-readable name for an exit code
func GetExitCodeName(code int) string {
	switch code {
	case ExitSuccess:
		return "Success"
	case ExitGeneralError:
		return "General Error"
	case ExitInputError:
		return "Input Error"
	case ExitOutputError:
		return "Output Error"
	case ExitProcessingError:
		return "Processing Error"
	case ExitConfigurationError:
		return "Configuration Error"
	case ExitResourceError:
		return "Resource Error"
	case ExitInterrupted:
		return "Interrupted"
	default:
		return "Unknown Error"
	}
}
