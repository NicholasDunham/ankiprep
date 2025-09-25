package integration

import (
	"os"
	"testing"
	"path/filepath"
	"strings"
	"ankiprep/internal/app"
)

// TestBackwardCompatibility ensures new --keep-header flag doesn't break existing behavior
// This test MUST FAIL until backward compatibility is properly implemented
func TestBackwardCompatibility(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "backward_compat_test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create test CSV with header
	csvFile := filepath.Join(tmpDir, "test.csv")
	csvContent := `Question,Answer,Extra
"What is 2+2?","4","math"
"What is the capital of France?","Paris","geography"
`
	err = os.WriteFile(csvFile, []byte(csvContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create CSV file: %v", err)
	}

	t.Run("default behavior removes header (backward compatibility)", func(t *testing.T) {
		// Process without --keep-header flag (default behavior)
		config := app.ProcessorConfig{
			Verbose: true,
			// keepHeader: false (default)
		}
		processor := app.NewProcessor(config)

		err := processor.ProcessFiles([]string{csvFile})
		if err != nil {
			t.Fatalf("Processing failed: %v", err)
		}

		// Check output file was created
		ankiFile := strings.TrimSuffix(csvFile, ".csv") + "_anki.csv"
		content, err := os.ReadFile(ankiFile)
		if err != nil {
			t.Fatalf("Failed to read output file: %v", err)
		}

		lines := strings.Split(strings.TrimSpace(string(content)), "\n")
		
		// Should have Anki header + 2 data rows (original header should be gone)
		if len(lines) != 3 {
			t.Errorf("Expected 3 lines (Anki header + 2 data rows), got %d lines", len(lines))
		}

		// First line should be Anki header, not original CSV header
		if !strings.Contains(lines[0], "Front") || !strings.Contains(lines[0], "Back") {
			t.Errorf("Expected Anki header as first line, got: %s", lines[0])
		}

		// Original CSV header should not appear in output
		if strings.Contains(string(content), "Question,Answer,Extra") {
			t.Error("Original CSV header should not appear in output (backward compatibility broken)")
		}

		// Data rows should be preserved
		if !strings.Contains(string(content), "What is 2+2?") {
			t.Error("Expected data rows to be preserved")
		}
		if !strings.Contains(string(content), "What is the capital of France?") {
			t.Error("Expected data rows to be preserved")
		}

		t.Logf("Backward compatibility preserved: %d lines in output", len(lines))
	})

	t.Run("existing CSV processing workflow unchanged", func(t *testing.T) {
		// Test that existing code patterns still work exactly the same
		config := app.ProcessorConfig{Verbose: false} // Typical existing usage
		processor := app.NewProcessor(config)

		// This should work exactly as before the new flag was added
		err := processor.ProcessFiles([]string{csvFile})
		if err != nil {
			t.Fatalf("Existing workflow broken: %v", err)
		}

		ankiFile := strings.TrimSuffix(csvFile, ".csv") + "_anki.csv"
		if _, err := os.Stat(ankiFile); os.IsNotExist(err) {
			t.Fatal("Output file not created - existing workflow broken")
		}

		// Output should match previous behavior exactly
		content, _ := os.ReadFile(ankiFile)
		if strings.Contains(string(content), "Question,Answer,Extra") {
			t.Error("Backward compatibility broken: original header should be removed by default")
		}

		t.Log("Existing workflow preserved")
	})
}