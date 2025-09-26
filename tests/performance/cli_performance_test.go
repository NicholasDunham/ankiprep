package performance

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// TestCLIPerformance tests ankiprep CLI performance with various file sizes
func TestCLIPerformance(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping performance tests in short mode")
	}

	tests := []struct {
		name           string
		numRows        int
		numColumns     int
		maxTimeSeconds float64
	}{
		{
			name:           "small file (1K rows)",
			numRows:        1000,
			numColumns:     3,
			maxTimeSeconds: 5.0,
		},
		{
			name:           "medium file (10K rows)",
			numRows:        10000,
			numColumns:     5,
			maxTimeSeconds: 15.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temporary directory
			tmpDir, err := os.MkdirTemp("", "perf_test")
			if err != nil {
				t.Fatalf("Failed to create temp directory: %v", err)
			}
			defer os.RemoveAll(tmpDir)

			// Generate test data
			inputFile := filepath.Join(tmpDir, "input.csv")
			outputFile := filepath.Join(tmpDir, "output.csv")

			if err := generateTestCSV(inputFile, tt.numRows, tt.numColumns); err != nil {
				t.Fatalf("Failed to generate test CSV: %v", err)
			}

			// Measure CLI processing time
			startTime := time.Now()
			cmd := exec.Command("ankiprep", "-o", outputFile, inputFile)
			output, err := cmd.CombinedOutput()
			duration := time.Since(startTime)

			if err != nil {
				t.Fatalf("CLI command failed: %v, output: %s", err, output)
			}

			// Check processing time
			durationSeconds := duration.Seconds()
			if durationSeconds > tt.maxTimeSeconds {
				t.Errorf("Processing took %.2fs, expected <= %.2fs", durationSeconds, tt.maxTimeSeconds)
			} else {
				t.Logf("Processing completed in %.2fs (within %.2fs limit)", durationSeconds, tt.maxTimeSeconds)
			}

			// Verify output exists and has content
			outputInfo, err := os.Stat(outputFile)
			if err != nil {
				t.Fatalf("Output file not created: %v", err)
			}

			if outputInfo.Size() == 0 {
				t.Error("Output file is empty")
			}

			// Verify output contains expected structure
			outputContent, err := os.ReadFile(outputFile)
			if err != nil {
				t.Fatalf("Failed to read output file: %v", err)
			}

			outputStr := string(outputContent)
			if !strings.Contains(outputStr, "#separator:comma") {
				t.Error("Output missing Anki metadata")
			}

			// Count data lines (excluding metadata)
			lines := strings.Split(strings.TrimSpace(outputStr), "\n")
			dataLines := 0
			for _, line := range lines {
				if !strings.HasPrefix(line, "#") && strings.TrimSpace(line) != "" {
					dataLines++
				}
			}

			if dataLines != tt.numRows {
				t.Errorf("Expected %d data lines, got %d", tt.numRows, dataLines)
			}
		})
	}
}

// generateTestCSV creates a CSV file with the specified number of rows and columns
func generateTestCSV(filename string, numRows, numColumns int) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// Write header
	var headers []string
	for i := 0; i < numColumns; i++ {
		headers = append(headers, fmt.Sprintf("col%d", i+1))
	}
	if _, err := file.WriteString(strings.Join(headers, ",") + "\n"); err != nil {
		return err
	}

	// Write data rows
	for row := 0; row < numRows; row++ {
		var values []string
		for col := 0; col < numColumns; col++ {
			values = append(values, fmt.Sprintf("value_%d_%d", row, col))
		}
		if _, err := file.WriteString(strings.Join(values, ",") + "\n"); err != nil {
			return err
		}
	}

	return nil
}
