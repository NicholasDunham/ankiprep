package integration
package integration_test

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"ankiprep/cmd/ankiprep"
	"ankiprep/internal/app"
)

// TestCLI_ClozeProcessing tests the CLI command with cloze deletion scenarios
func TestCLI_ClozeProcessing(t *testing.T) {
	// Create temporary directory for test files
	tmpDir, err := os.MkdirTemp("", "ankiprep_cli_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	t.Run("process file with cloze deletions", func(t *testing.T) {
		// Create input CSV with cloze deletions
		inputFile := filepath.Join(tmpDir, "input_with_cloze.csv")
		outputFile := filepath.Join(tmpDir, "output_with_cloze.csv")

		inputData := [][]string{
			{"Front", "Back", "Tags"},
			{"Question : What is the capital of France?", "{{c1::Paris : the city of light}}", "geography"},
			{"Cities : Name two major cities", "{{c1::Paris}} and {{c2::Rome : Italy}}", "cities"},
			{"Simple : Regular question", "Regular answer without cloze", "basic"},
		}

		// Write input file
		file, err := os.Create(inputFile)
		if err != nil {
			t.Fatalf("Failed to create input file: %v", err)
		}

		writer := csv.NewWriter(file)
		err = writer.WriteAll(inputData)
		file.Close()
		if err != nil {
			t.Fatalf("Failed to write input CSV: %v", err)
		}

		// Create CLI application
		app := app.NewApp()
		
		// Simulate command line arguments
		// This will fail until we implement the CLI
		args := []string{
			"ankiprep",
			"--input", inputFile,
			"--output", outputFile,
			"--french",
			"--smart-quotes",
		}

		// Capture stdout/stderr
		var stdout, stderr bytes.Buffer
		app.SetOutput(&stdout)
		app.SetErrOutput(&stderr)

		// Execute command
		err = app.Execute(args[1:]) // Skip program name
		if err != nil {
			t.Fatalf("CLI execution failed: %v\nStderr: %s", err, stderr.String())
		}

		// Verify output file was created
		if _, err := os.Stat(outputFile); os.IsNotExist(err) {
			t.Fatal("Output file was not created")
		}

		// Read and verify output
		outputFileHandle, err := os.Open(outputFile)
		if err != nil {
			t.Fatalf("Failed to open output file: %v", err)
		}
		defer outputFileHandle.Close()

		reader := csv.NewReader(outputFileHandle)
		outputRecords, err := reader.ReadAll()
		if err != nil {
			t.Fatalf("Failed to read output CSV: %v", err)
		}

		// Verify header is preserved
		if !equalStringSlice(outputRecords[0], inputData[0]) {
			t.Errorf("Header not preserved: got %v, want %v", outputRecords[0], inputData[0])
		}

		// Verify processing results
		expectedResults := []struct {
			row           int
			expectedFront string
			expectedBack  string
		}{
			{
				row:           1,
				expectedFront: "Question\u00A0: What is the capital of France?", // Colon rule applied
				expectedBack:  "{{c1::Paris : the city of light}}", // Colon inside cloze unchanged
			},
			{
				row:           2,
				expectedFront: "Cities\u00A0: Name two major cities", // Colon rule applied
				expectedBack:  "{{c1::Paris}} and {{c2::Rome : Italy}}", // Colon inside cloze unchanged
			},
			{
				row:           3,
				expectedFront: "Simple\u00A0: Regular question", // Colon rule applied
				expectedBack:  "Regular answer without cloze", // No change expected
			},
		}

		for _, expected := range expectedResults {
			if len(outputRecords) <= expected.row {
				t.Errorf("Missing row %d in output", expected.row)
				continue
			}

			actualRecord := outputRecords[expected.row]
			if len(actualRecord) < 2 {
				t.Errorf("Row %d has insufficient columns", expected.row)
				continue
			}

			if actualRecord[0] != expected.expectedFront {
				t.Errorf("Row %d front: got %q, want %q", 
					expected.row, actualRecord[0], expected.expectedFront)
			}

			if actualRecord[1] != expected.expectedBack {
				t.Errorf("Row %d back: got %q, want %q", 
					expected.row, actualRecord[1], expected.expectedBack)
			}

			// Tags should be unchanged
			if actualRecord[2] != inputData[expected.row][2] {
				t.Errorf("Row %d tags changed: got %q, want %q", 
					expected.row, actualRecord[2], inputData[expected.row][2])
			}
		}
	})

	t.Run("CLI error handling", func(t *testing.T) {
		app := app.NewApp()
		var stdout, stderr bytes.Buffer
		app.SetOutput(&stdout)
		app.SetErrOutput(&stderr)

		// Test missing input file
		args := []string{
			"--input", "/nonexistent/file.csv",
			"--output", filepath.Join(tmpDir, "output.csv"),
			"--french",
		}

		err := app.Execute(args)
		if err == nil {
			t.Error("Expected error for nonexistent input file")
		}

		// Should have helpful error message
		stderrOutput := stderr.String()
		if !strings.Contains(stderrOutput, "input") && !strings.Contains(stderrOutput, "file") {
			t.Errorf("Error message should mention input file issue: %s", stderrOutput)
		}
	})

	t.Run("CLI help and version", func(t *testing.T) {
		app := app.NewApp()

		// Test help command
		var stdout bytes.Buffer
		app.SetOutput(&stdout)

		err := app.Execute([]string{"--help"})
		if err != nil {
			t.Errorf("Help command failed: %v", err)
		}

		helpOutput := stdout.String()
		expectedHelpContent := []string{
			"ankiprep", // Program name
			"input",    // Input flag
			"output",   // Output flag
			"french",   // French flag
		}

		for _, content := range expectedHelpContent {
			if !strings.Contains(helpOutput, content) {
				t.Errorf("Help output should contain %q: %s", content, helpOutput)
			}
		}
	})

	t.Run("CLI progress reporting", func(t *testing.T) {
		// Create larger input file for progress testing
		inputFile := filepath.Join(tmpDir, "large_input.csv")
		outputFile := filepath.Join(tmpDir, "large_output.csv")

		// Create file with many rows
		file, err := os.Create(inputFile)
		if err != nil {
			t.Fatalf("Failed to create large input file: %v", err)
		}

		writer := csv.NewWriter(file)
		writer.Write([]string{"Front", "Back", "Tags"}) // Header

		// Add many rows
		for i := 0; i < 100; i++ {
			writer.Write([]string{
				fmt.Sprintf("Question %d : What is this?", i),
				fmt.Sprintf("{{c1::Answer %d}}", i),
				"test",
			})
		}
		writer.Flush()
		file.Close()

		if writer.Error() != nil {
			t.Fatalf("Failed to write large input CSV: %v", writer.Error())
		}

		// Run CLI with progress reporting
		app := app.NewApp()
		var stdout, stderr bytes.Buffer
		app.SetOutput(&stdout)
		app.SetErrOutput(&stderr)

		args := []string{
			"--input", inputFile,
			"--output", outputFile,
			"--french",
			"--progress", // Enable progress reporting
		}

		err = app.Execute(args)
		if err != nil {
			t.Fatalf("CLI with progress failed: %v", err)
		}

		// Check that progress was reported to stderr
		progressOutput := stderr.String()
		if !strings.Contains(progressOutput, "Processing") {
			t.Error("Expected progress reporting in stderr output")
		}
	})
}

// equalStringSlice compares two string slices for equality
func equalStringSlice(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}