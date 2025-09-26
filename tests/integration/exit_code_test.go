package integration

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

// TestExitCodes tests that the CLI returns correct exit codes
// This test MUST FAIL until proper exit code handling is implemented
func TestExitCodes(t *testing.T) {
	// Build the CLI binary for testing
	binPath := filepath.Join(os.TempDir(), "ankiprep_test")
	cmd := exec.Command("go", "build", "-o", binPath, "../../cmd/ankiprep")
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to build CLI binary: %v", err)
	}
	defer os.Remove(binPath)

	t.Run("successful processing returns exit code 0", func(t *testing.T) {
		tmpDir, err := os.MkdirTemp("", "exit_code_test")
		if err != nil {
			t.Fatalf("Failed to create temp directory: %v", err)
		}
		defer os.RemoveAll(tmpDir)

		// Create valid CSV
		csvFile := filepath.Join(tmpDir, "test.csv")
		csvContent := `Question,Answer
"What is 2+2?","4"
"What is the capital of France?","Paris"
`
		err = os.WriteFile(csvFile, []byte(csvContent), 0644)
		if err != nil {
			t.Fatalf("Failed to create CSV file: %v", err)
		}

		// Run CLI
		cmd := exec.Command(binPath, csvFile)
		err = cmd.Run()

		// Should succeed with exit code 0
		if err != nil {
			t.Errorf("Expected exit code 0 for successful processing, got error: %v", err)
		}
	})

	t.Run("file not found returns exit code 1", func(t *testing.T) {
		// Run CLI with non-existent file
		cmd := exec.Command(binPath, "/path/that/does/not/exist.csv")
		err := cmd.Run()

		// Should fail with specific exit code
		if err == nil {
			t.Error("Expected non-zero exit code for file not found")
		} else if exitError, ok := err.(*exec.ExitError); ok {
			exitCode := exitError.ExitCode()
			if exitCode != 1 {
				t.Errorf("Expected exit code 1 for file not found, got %d", exitCode)
			}
		}
	})

	t.Run("invalid arguments return exit code 2", func(t *testing.T) {
		// Run CLI with invalid flag
		cmd := exec.Command(binPath, "--invalid-flag")
		err := cmd.Run()

		// Should fail with specific exit code
		if err == nil {
			t.Error("Expected non-zero exit code for invalid arguments")
		} else if exitError, ok := err.(*exec.ExitError); ok {
			exitCode := exitError.ExitCode()
			if exitCode != 2 {
				t.Errorf("Expected exit code 2 for invalid arguments, got %d", exitCode)
			}
		}
	})

	t.Run("help flag returns exit code 0", func(t *testing.T) {
		// Run CLI with help flag
		cmd := exec.Command(binPath, "--help")
		err := cmd.Run()

		// Should succeed with exit code 0
		if err != nil {
			t.Errorf("Expected exit code 0 for help flag, got error: %v", err)
		}
	})
}
