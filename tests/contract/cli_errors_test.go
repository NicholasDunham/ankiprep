package contract

import (
	"os/exec"
	"strings"
	"testing"
)

// TestCLIErrorNoArguments tests that command fails with proper exit code when no arguments provided
func TestCLIErrorNoArguments(t *testing.T) {
	cmd := exec.Command("ankiprep")
	output, err := cmd.CombinedOutput()
	
	// Should fail with exit code 1 (invalid arguments)
	if err == nil {
		t.Errorf("Expected command to fail when no arguments provided")
	}
	
	// Check exit code
	if exitError, ok := err.(*exec.ExitError); ok {
		if exitError.ExitCode() != 1 {
			t.Errorf("Expected exit code 1 for invalid arguments, got %d", exitError.ExitCode())
		}
	}
	
	outputStr := string(output)
	
	// Should show usage information
	expectedErrorMessages := []string{
		"required",
		"input",
		"file",
	}
	
	for _, expected := range expectedErrorMessages {
		if !strings.Contains(strings.ToLower(outputStr), expected) {
			t.Errorf("Error output should contain '%s', but got: %s", expected, outputStr)
		}
	}
}

// TestCLIErrorFileNotFound tests proper error handling for non-existent input files
func TestCLIErrorFileNotFound(t *testing.T) {
	cmd := exec.Command("ankiprep", "/nonexistent/file.csv")
	output, err := cmd.CombinedOutput()
	
	// Should fail with exit code 2 (file not found)
	if err == nil {
		t.Errorf("Expected command to fail when input file not found")
	}
	
	// Check exit code
	if exitError, ok := err.(*exec.ExitError); ok {
		if exitError.ExitCode() != 2 {
			t.Errorf("Expected exit code 2 for file not found, got %d", exitError.ExitCode())
		}
	}
	
	outputStr := string(output)
	
	// Should show file not found error
	expectedErrorMessages := []string{
		"not found",
		"nonexistent",
	}
	
	for _, expected := range expectedErrorMessages {
		if !strings.Contains(strings.ToLower(outputStr), expected) {
			t.Errorf("Error output should contain '%s', but got: %s", expected, outputStr)
		}
	}
}

// TestCLIErrorInvalidOption tests handling of invalid command line options
func TestCLIErrorInvalidOption(t *testing.T) {
	cmd := exec.Command("ankiprep", "--invalid-option", "file.csv")
	output, err := cmd.CombinedOutput()
	
	// Should fail with exit code 1 (invalid arguments)
	if err == nil {
		t.Errorf("Expected command to fail with invalid option")
	}
	
	// Check exit code
	if exitError, ok := err.(*exec.ExitError); ok {
		if exitError.ExitCode() != 1 {
			t.Errorf("Expected exit code 1 for invalid option, got %d", exitError.ExitCode())
		}
	}
	
	outputStr := string(output)
	
	// Should show unknown flag error
	expectedErrorMessages := []string{
		"unknown",
		"flag",
		"invalid-option",
	}
	
	for _, expected := range expectedErrorMessages {
		if !strings.Contains(strings.ToLower(outputStr), expected) {
			t.Errorf("Error output should contain '%s', but got: %s", expected, outputStr)
		}
	}
}

// TestCLIErrorPermissionDenied tests handling of permission denied errors
func TestCLIErrorPermissionDenied(t *testing.T) {
	// Try to write output to a directory that should be read-only
	cmd := exec.Command("ankiprep", "-o", "/root/output.csv", "/etc/passwd")
	output, err := cmd.CombinedOutput()
	
	// Should fail with exit code 5 (output error) or 2 (input error)
	if err == nil {
		t.Errorf("Expected command to fail with permission denied")
	}
	
	// Check exit code (could be 2 for input file access or 5 for output write)
	if exitError, ok := err.(*exec.ExitError); ok {
		exitCode := exitError.ExitCode()
		if exitCode != 2 && exitCode != 5 {
			t.Errorf("Expected exit code 2 or 5 for permission errors, got %d", exitCode)
		}
	}
	
	outputStr := string(output)
	
	// Should show permission error
	expectedErrorMessages := []string{
		"permission",
		"denied",
	}
	
	foundPermissionError := false
	for _, expected := range expectedErrorMessages {
		if strings.Contains(strings.ToLower(outputStr), expected) {
			foundPermissionError = true
			break
		}
	}
	
	if !foundPermissionError {
		t.Errorf("Error output should contain permission error, but got: %s", outputStr)
	}
}