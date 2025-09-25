package integration
package integration

import (
	"os"
	"testing"
	"path/filepath"
	"strings"
	"ankiprep/internal/app"
)

// TestErrorHandlingScenarios tests error handling for various failure modes
// This test MUST FAIL until proper error handling is implemented
func TestErrorHandlingScenarios(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "error_handling_test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	t.Run("file not found error", func(t *testing.T) {
		config := app.ProcessorConfig{Verbose: true}
		processor := app.NewProcessor(config)
		
		// Try to process non-existent file
		nonExistentFile := filepath.Join(tmpDir, "doesnotexist.csv")
		err := processor.ProcessFiles([]string{nonExistentFile})
		
		// Should fail with descriptive error
		if err == nil {
			t.Fatal("Expected error for non-existent file")
		}
		
		// Check error message is descriptive
		if !strings.Contains(err.Error(), "file not found") &&
		   !strings.Contains(err.Error(), "no such file") {
			t.Errorf("Expected file not found error, got: %v", err)
		}
		
		t.Logf("File not found error (expected): %v", err)
	})

	t.Run("permission denied error", func(t *testing.T) {
		// Create a file we can't read
		restrictedFile := filepath.Join(tmpDir, "restricted.csv") 
		err := os.WriteFile(restrictedFile, []byte("test"), 0000)
		if err != nil {
			t.Fatalf("Failed to create restricted file: %v", err)
		}
		defer os.Chmod(restrictedFile, 0644) // cleanup
		
		config := app.ProcessorConfig{Verbose: true}
		processor := app.NewProcessor(config)
		
		err = processor.ProcessFiles([]string{restrictedFile})
		
		// Should fail with permission error
		if err == nil {
			t.Fatal("Expected permission error for restricted file")
		}
		
		t.Logf("Permission error (expected): %v", err)
	})

	t.Run("malformed CSV error", func(t *testing.T) {
		// Create malformed CSV file
		malformedFile := filepath.Join(tmpDir, "malformed.csv")
		malformedContent := `"Unclosed quote field
"Normal field","Another field"
`
		
		err := os.WriteFile(malformedFile, []byte(malformedContent), 0644)
		if err != nil {
			t.Fatalf("Failed to create malformed CSV: %v", err)
		}
		
		config := app.ProcessorConfig{Verbose: true}
		processor := app.NewProcessor(config)
		
		err = processor.ProcessFiles([]string{malformedFile})
		
		// Should fail with CSV format error
		if err == nil {
			t.Fatal("Expected CSV format error for malformed file")
		}
		
		t.Logf("CSV format error (expected): %v", err)
	})
}