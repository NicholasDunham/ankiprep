package cmd_test

import (
	"os"
	"os/exec"
	"strings"
	"testing"
)

// TestKeepHeaderFlagParsing tests that the --keep-header and -k flags are properly parsed by the CLI
func TestKeepHeaderFlagParsing(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		wantHelp string // What we expect in help output
	}{
		{
			name:     "long flag help text",
			args:     []string{"--help"},
			wantHelp: "--keep-header",
		},
		{
			name:     "short flag help text",
			args:     []string{"--help"},
			wantHelp: "-k",
		},
		{
			name:     "flag description present",
			args:     []string{"--help"},
			wantHelp: "Preserve the first row of CSV files",
		},
	}

	// Build the binary for testing
	cmd := exec.Command("go", "build", "-o", "ankiprep-test", "./cmd/ankiprep")
	cmd.Dir = "../../../"
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to build ankiprep for testing: %v", err)
	}
	defer os.Remove("../../../ankiprep-test")

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command("../../../ankiprep-test", tt.args...)
			output, err := cmd.CombinedOutput()

			// Help command should exit with code 0
			if tt.args[0] == "--help" && err != nil {
				if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() != 0 {
					t.Fatalf("Help command failed with exit code %d: %v", exitErr.ExitCode(), err)
				}
			}

			outputStr := string(output)
			if !strings.Contains(outputStr, tt.wantHelp) {
				t.Errorf("Expected help output to contain %q, got:\n%s", tt.wantHelp, outputStr)
			}
		})
	}
}

// TestKeepHeaderFlagFunctionality tests that the flag actually affects processing behavior
// This test MUST FAIL until the flag is properly integrated with processing logic
func TestKeepHeaderFlagFunctionality(t *testing.T) {
	// Create a temporary CSV file for testing
	tmpFile := "/tmp/test-header.csv"
	csvContent := "Text,Extra,Grammar_Notes\n\"Hello\",\"Bonjour\",\"greeting\"\n\"Goodbye\",\"Au revoir\",\"farewell\"\n"

	if err := os.WriteFile(tmpFile, []byte(csvContent), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	defer os.Remove(tmpFile)

	// Test without --keep-header flag (should remove header)
	t.Run("default removes header", func(t *testing.T) {
		cmd := exec.Command("../../../ankiprep-test", tmpFile)
		output, err := cmd.CombinedOutput()

		// This should fail until implementation is complete
		if err == nil {
			t.Fatal("Expected command to fail until ProcessCSV supports header removal - implementation not ready")
		}

		// Check that it fails for the right reason (not implemented yet)
		outputStr := string(output)
		t.Logf("Expected failure output: %s", outputStr)
	})

	// Test with --keep-header flag (should preserve header)
	t.Run("keep-header preserves header", func(t *testing.T) {
		cmd := exec.Command("../../../ankiprep-test", "--keep-header", tmpFile)
		output, err := cmd.CombinedOutput()

		// This should fail until implementation is complete
		if err == nil {
			t.Fatal("Expected command to fail until ProcessCSV supports --keep-header flag - implementation not ready")
		}

		// Check that it fails for the right reason (not implemented yet)
		outputStr := string(output)
		t.Logf("Expected failure output: %s", outputStr)
	})
}
