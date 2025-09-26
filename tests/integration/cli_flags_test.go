package integration

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// TestKeepHeaderFlag tests the --keep-header flag functionality via CLI
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

	t.Run("keep header with -k flag", func(t *testing.T) {
		outputFile := filepath.Join(tmpDir, "output_short.csv")

		// Test with -k (short flag)
		cmd := exec.Command("ankiprep", "-k", "-o", outputFile, inputFile)
		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("Command failed: %v, output: %s", err, output)
		}

		// Read the output
		result, err := os.ReadFile(outputFile)
		if err != nil {
			t.Fatalf("Failed to read output file: %v", err)
		}

		resultStr := string(result)

		// Should contain the original header as a data row (first data line after metadata)
		lines := strings.Split(strings.TrimSpace(resultStr), "\n")
		dataLines := []string{}
		for _, line := range lines {
			if !strings.HasPrefix(line, "#") && strings.TrimSpace(line) != "" {
				dataLines = append(dataLines, line)
			}
		}

		// Should have 3 data lines (header + 2 data)
		if len(dataLines) != 3 {
			t.Errorf("Expected 3 data lines (header + data), got %d lines: %v", len(dataLines), dataLines)
		}

		// First data line should be the header
		if len(dataLines) > 0 && !strings.Contains(dataLines[0], "Question") {
			t.Errorf("Expected first data line to contain header, got: %s", dataLines[0])
		}

		// Should still have Anki metadata headers
		if !strings.Contains(resultStr, "#separator:comma") {
			t.Errorf("Expected Anki metadata headers, got: %s", resultStr)
		}
	})

	t.Run("keep header with --keep-header flag", func(t *testing.T) {
		outputFile := filepath.Join(tmpDir, "output_long.csv")

		// Test with --keep-header (long flag)
		cmd := exec.Command("ankiprep", "--keep-header", "-o", outputFile, inputFile)
		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("Command failed: %v, output: %s", err, output)
		}

		// Read the output
		result, err := os.ReadFile(outputFile)
		if err != nil {
			t.Fatalf("Failed to read output file: %v", err)
		}

		resultStr := string(result)

		// Should contain the original header as a data row (first data line after metadata)
		lines := strings.Split(strings.TrimSpace(resultStr), "\n")
		dataLines := []string{}
		for _, line := range lines {
			if !strings.HasPrefix(line, "#") && strings.TrimSpace(line) != "" {
				dataLines = append(dataLines, line)
			}
		}

		// Should have 3 data lines (header + 2 data)
		if len(dataLines) != 3 {
			t.Errorf("Expected 3 data lines (header + data), got %d lines: %v", len(dataLines), dataLines)
		}
	})

	t.Run("without keep header flag (default behavior)", func(t *testing.T) {
		outputFile := filepath.Join(tmpDir, "output_default.csv")

		// Test without --keep-header (default behavior should not include header as data)
		cmd := exec.Command("ankiprep", "-o", outputFile, inputFile)
		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("Command failed: %v, output: %s", err, output)
		}

		// Read the output
		result, err := os.ReadFile(outputFile)
		if err != nil {
			t.Fatalf("Failed to read output file: %v", err)
		}

		resultStr := string(result)

		// Should contain the column names in the #columns metadata line
		if !strings.Contains(resultStr, "#columns:Question,Answer,Extra") {
			t.Errorf("Expected column metadata, got: %s", resultStr)
		}

		// Should contain the actual data
		if !strings.Contains(resultStr, "What is 2+2?") {
			t.Errorf("Expected data content, got: %s", resultStr)
		}

		// Should NOT contain the header as a data row (count data lines)
		lines := strings.Split(strings.TrimSpace(resultStr), "\n")
		dataLines := 0
		for _, line := range lines {
			if !strings.HasPrefix(line, "#") && strings.TrimSpace(line) != "" {
				dataLines++
			}
		}
		// Should have only 2 data lines (not 3 with header)
		if dataLines != 2 {
			t.Errorf("Expected 2 data lines (without header), got %d lines in: %s", dataLines, resultStr)
		}

		// Should still have Anki metadata headers
		if !strings.Contains(resultStr, "#separator:comma") {
			t.Errorf("Expected Anki metadata headers, got: %s", resultStr)
		}
	})
}

// TestVerboseFlag tests the --verbose flag functionality
func TestVerboseFlag(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "verbose_test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create simple test CSV
	inputFile := filepath.Join(tmpDir, "input.csv")
	csvContent := `front,back
hello,world
`
	err = os.WriteFile(inputFile, []byte(csvContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test input file: %v", err)
	}

	t.Run("verbose output with -v", func(t *testing.T) {
		outputFile := filepath.Join(tmpDir, "output.csv")

		cmd := exec.Command("ankiprep", "-v", "-o", outputFile, inputFile)
		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("Command failed: %v, output: %s", err, output)
		}

		outputStr := string(output)

		// Verbose mode should include more detailed processing information
		// At minimum, it should still show the "Done. Processed X entries" message
		if !strings.Contains(outputStr, "Processed") || !strings.Contains(outputStr, "entries") {
			t.Errorf("Expected verbose processing information, got: %s", outputStr)
		}
	})
}

// TestCLIContract validates essential CLI interface contracts
func TestCLIContract(t *testing.T) {
	t.Run("help flag displays usage", func(t *testing.T) {
		cmd := exec.Command("ankiprep", "--help")
		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("Help command should succeed, got error: %v", err)
		}

		outputStr := string(output)

		// Must contain basic CLI information
		requiredStrings := []string{
			"ankiprep",
			"Usage:",
			"Flags:",
			"--help",
			"--version",
			"--output",
			"--french",
			"--verbose",
		}

		for _, required := range requiredStrings {
			if !strings.Contains(outputStr, required) {
				t.Errorf("Help output missing required string '%s', got: %s", required, outputStr)
			}
		}
	})

	t.Run("version flag displays version", func(t *testing.T) {
		cmd := exec.Command("ankiprep", "--version")
		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("Version command should succeed, got error: %v", err)
		}

		outputStr := string(output)

		// Should contain version information
		if !strings.Contains(outputStr, "1.0.0") {
			t.Errorf("Version output should contain version number, got: %s", outputStr)
		}
	})

	t.Run("no arguments shows error", func(t *testing.T) {
		cmd := exec.Command("ankiprep")
		output, err := cmd.CombinedOutput()

		// Should fail when no arguments provided
		if err == nil {
			t.Errorf("Command should fail when no arguments provided")
		}

		outputStr := string(output)
		// Should provide some usage guidance
		if !strings.Contains(outputStr, "required") && !strings.Contains(outputStr, "Usage:") {
			t.Errorf("Error output should provide usage guidance, got: %s", outputStr)
		}
	})
}
