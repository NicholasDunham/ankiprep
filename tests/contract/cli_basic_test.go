package contract

import (
	"os/exec"
	"strings"
	"testing"
)

// TestCLIBasicHelp tests the --help flag functionality
func TestCLIBasicHelp(t *testing.T) {
	cmd := exec.Command("ankiprep", "--help")
	output, err := cmd.CombinedOutput()
	
	if err != nil {
		t.Fatalf("Expected help command to succeed, got error: %v", err)
	}
	
	outputStr := string(output)
	
	// Check for expected help content
	expectedStrings := []string{
		"ankiprep",
		"CSV to Anki",
		"--help",
		"--version",
		"--output",
		"--french",
		"--verbose",
	}
	
	for _, expected := range expectedStrings {
		if !strings.Contains(outputStr, expected) {
			t.Errorf("Help output should contain '%s', but got: %s", expected, outputStr)
		}
	}
}

// TestCLIBasicVersion tests the --version flag functionality
func TestCLIBasicVersion(t *testing.T) {
	cmd := exec.Command("ankiprep", "--version")
	output, err := cmd.CombinedOutput()
	
	if err != nil {
		t.Fatalf("Expected version command to succeed, got error: %v", err)
	}
	
	outputStr := string(output)
	
	// Check for version information
	if !strings.Contains(outputStr, "1.0.0") {
		t.Errorf("Version output should contain '1.0.0', but got: %s", outputStr)
	}
}

// TestCLIBasicShortFlags tests short flag versions
func TestCLIBasicShortFlags(t *testing.T) {
	// Test -h flag
	cmd := exec.Command("ankiprep", "-h")
	output, err := cmd.CombinedOutput()
	
	if err != nil {
		t.Fatalf("Expected -h command to succeed, got error: %v", err)
	}
	
	outputStr := string(output)
	if !strings.Contains(outputStr, "ankiprep") {
		t.Errorf("-h should show help, but got: %s", outputStr)
	}
	
	// Test -V flag
	cmd = exec.Command("ankiprep", "-V")
	output, err = cmd.CombinedOutput()
	
	if err != nil {
		t.Fatalf("Expected -V command to succeed, got error: %v", err)
	}
	
	outputStr = string(output)
	if !strings.Contains(outputStr, "1.0.0") {
		t.Errorf("-V should show version, but got: %s", outputStr)
	}
}