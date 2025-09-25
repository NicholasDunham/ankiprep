package integration_test

import (
	"ankiprep/internal/app"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestKeepHeaderFlag tests that the --keep-header flag preserves original CSV headers
// This test MUST FAIL until proper keep-header logic is implemented
func TestKeepHeaderFlag(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "keep_header_test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create test CSV with header
	inputFile := filepath.Join(tmpDir, "input.csv")
	csvContent := `Question,Answer,Extra
"What is 2+2?","4","math"
"What is the capital of France?","Paris","geography"
`

	err = os.WriteFile(inputFile, []byte(csvContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test input file: %v", err)
	}

	outputFile := filepath.Join(tmpDir, "output.csv")

	t.Run("keep-header flag preserves original header", func(t *testing.T) {
		// Create processor with output path
		config := app.ProcessorConfig{
			OutputPath: outputFile,
			Verbose:    true,
		}

		processor := app.NewProcessor(config)

		// Process with keep header behavior
		options := app.ProcessingOptions{KeepHeader: true}
		_, err = processor.ProcessFiles([]string{inputFile}, options)
		if err != nil {
			t.Fatalf("Processing failed: %v", err)
		}

		// Check that output file was created
		if _, err := os.Stat(outputFile); os.IsNotExist(err) {
			t.Fatal("Output file was not created")
		}

		// Read and verify output content
		content, err := os.ReadFile(outputFile)
		if err != nil {
			t.Fatalf("Failed to read output file: %v", err)
		}

		lines := strings.Split(strings.TrimSpace(string(content)), "\n")

		// Should have: 3 Anki headers + original header + 2 data rows = 6 lines
		if len(lines) < 6 {
			t.Errorf("Expected at least 6 lines in output, got %d", len(lines))
		}

		// Original CSV header should appear as first data row (after Anki metadata)
		dataStartIndex := -1
		for i, line := range lines {
			if !strings.HasPrefix(line, "#") && strings.Contains(line, ",") {
				dataStartIndex = i
				break
			}
		}

		if dataStartIndex < 0 || dataStartIndex >= len(lines) {
			t.Fatal("Could not find data rows in output")
		}

		firstDataRow := strings.TrimSpace(lines[dataStartIndex])
		if firstDataRow != "Question,Answer,Extra" {
			t.Errorf("Expected original CSV header as first data row, got: %s", firstDataRow)
		}

		// Data should be preserved after the header row
		if !strings.Contains(string(content), "What is 2+2?") {
			t.Error("Expected data row to be preserved")
		}

		t.Log("Keep-header functionality works correctly")
	})
}
