package integration_test

import (
	"ankiprep/internal/app"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestDefaultHeaderRemoval tests that the original CSV header is removed by default
// This test MUST FAIL until proper header removal logic is implemented
func TestDefaultHeaderRemoval(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "header_removal_test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create test CSV with header
	inputFile := filepath.Join(tmpDir, "input.csv")
	csvContent := `Text,Extra,Grammar_Notes
"Hello world","Bonjour monde","Basic greeting"
"Goodbye","Au revoir","Formal farewell"
`

	err = os.WriteFile(inputFile, []byte(csvContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test input file: %v", err)
	}

	t.Run("default behavior removes header", func(t *testing.T) {
		// Create processor with default settings (no specific output path)
		config := app.ProcessorConfig{
			Verbose: true,
		}

		processor := app.NewProcessor(config)

		// Process the file - should remove header by default
		options := app.ProcessingOptions{KeepHeader: false}
		_, err = processor.ProcessFiles([]string{inputFile}, options)
		if err != nil {
			t.Fatalf("Processing failed: %v", err)
		}

		// Check that output file was created with default naming
		expectedOutput := strings.TrimSuffix(inputFile, ".csv") + ".anki.csv"
		if _, err := os.Stat(expectedOutput); os.IsNotExist(err) {
			t.Fatalf("Output file was not created at expected location: %s", expectedOutput)
		}

		// Read and verify output content
		content, err := os.ReadFile(expectedOutput)
		if err != nil {
			t.Fatalf("Failed to read output file: %v", err)
		}

		lines := strings.Split(strings.TrimSpace(string(content)), "\n")

		// Should have: 3 Anki headers + 1 column header + 2 data rows = 6 lines
		if len(lines) < 5 {
			t.Errorf("Expected at least 5 lines in output, got %d", len(lines))
		}

		// Original CSV header should NOT appear as a data row
		// It's OK if it appears in metadata (#columns: line), but not as an actual CSV row
		dataStartIndex := -1
		for i, line := range lines {
			if !strings.HasPrefix(line, "#") && strings.Contains(line, ",") {
				dataStartIndex = i
				break
			}
		}

		if dataStartIndex > 0 && dataStartIndex < len(lines) {
			// Check if the first data row is the original header
			firstDataRow := lines[dataStartIndex]
			if strings.TrimSpace(firstDataRow) == "Text,Extra,Grammar_Notes" {
				t.Error("Original CSV header should be removed but appears as first data row")
			}
		}

		// Data should be preserved
		if !strings.Contains(string(content), "Hello world") {
			t.Error("Expected data row to be preserved")
		}

		t.Log("Header removal works correctly")
	})
}
