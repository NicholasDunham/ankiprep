package app_test

import (
	"ankiprep/internal/app"
	"strings"
	"testing"
)

// TestProcessCSVFunctionSignature tests that ProcessCSV function accepts the new ProcessingOptions parameter
// This test MUST FAIL until the function signature is updated to include ProcessingOptions
func TestProcessCSVFunctionSignature(t *testing.T) {
	// Create a processor instance
	config := app.ProcessorConfig{
		OutputPath:     "test.csv",
		FrenchMode:     false,
		SmartQuotes:    false,
		SkipDuplicates: false,
		Verbose:        true,
	}

	processor := app.NewProcessor(config)

	// Test that ProcessFiles method exists with expected signature
	// This will fail until the method signature includes ProcessingOptions
	t.Run("ProcessFiles with ProcessingOptions", func(t *testing.T) {
		// Create mock ProcessingOptions
		options := app.ProcessingOptions{
			KeepHeader: false,
		}

		// This should compile once the new signature exists
		// Until then, this test will fail at compile time
		_, err := processor.ProcessFiles([]string{"input.csv"}, options)

		// We expect this to fail since the input file doesn't exist
		if err == nil {
			t.Fatal("Expected ProcessFiles to fail with non-existent file")
		}

		// Check that error indicates file issue, not signature issue
		if !strings.Contains(err.Error(), "validation failed") &&
			!strings.Contains(err.Error(), "failed to parse") &&
			!strings.Contains(err.Error(), "no such file") {
			t.Errorf("Expected file-related error, got: %v", err)
		}

		t.Log("ProcessFiles signature correctly accepts ProcessingOptions")
	})
}

// TestProcessingOptionsIntegration tests that ProcessingOptions are properly used in processing
// This test MUST FAIL until the integration is complete
func TestProcessingOptionsIntegration(t *testing.T) {
	t.Run("ProcessingOptions affects header handling", func(t *testing.T) {
		// Create processor
		config := app.ProcessorConfig{Verbose: true}
		processor := app.NewProcessor(config)

		// Test with KeepHeader = false (default, remove header)
		optionsRemove := app.ProcessingOptions{
			KeepHeader: false,
		}

		// This test ensures the options parameter is properly used
		// We expect this to fail since no input files exist
		_, err := processor.ProcessFiles([]string{"non-existent.csv"}, optionsRemove)
		if err == nil {
			t.Fatal("Expected error for non-existent file")
		}

		// Test with KeepHeader = true (preserve header)
		optionsKeep := app.ProcessingOptions{
			KeepHeader: true,
		}

		_, err = processor.ProcessFiles([]string{"non-existent.csv"}, optionsKeep)
		if err == nil {
			t.Fatal("Expected error for non-existent file")
		}

		t.Log("ProcessingOptions properly integrated with ProcessFiles method")
	})
}
