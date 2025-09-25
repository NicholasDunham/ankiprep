package contract

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// TestCLIFrenchTypographyFlag tests the --french flag functionality
func TestCLIFrenchTypographyFlag(t *testing.T) {
	// Create temporary input file
	tempDir, err := ioutil.TempDir("", "ankiprep_test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)
	
	inputFile := filepath.Join(tempDir, "input.csv")
	outputFile := filepath.Join(tempDir, "output.csv")
	
	// Write test CSV content with French text requiring typography fixes
	csvContent := `front,back
"Bonjour : comment allez-vous ?","Hello: how are you?"
"Il a dit : « C'est fantastique ! »","He said: \"It's fantastic!\""
"Voulez-vous ; vraiment ?","Do you want; really?"`
	
	err = ioutil.WriteFile(inputFile, []byte(csvContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write input file: %v", err)
	}
	
	// Run ankiprep command with French typography flag
	cmd := exec.Command("ankiprep", "--french", "-o", outputFile, inputFile)
	output, err := cmd.CombinedOutput()
	
	if err != nil {
		t.Fatalf("Expected French processing command to succeed, got error: %v, output: %s", err, string(output))
	}
	
	// Check output file content
	outputContent, err := ioutil.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}
	
	outputStr := string(outputContent)
	
	// Check for French typography transformations
	// NNBSP (U+202F) should be inserted before punctuation
	expectedTransformations := []string{
		"Bonjour\u202F:", // NNBSP before colon
		"vous\u202F?",    // NNBSP before question mark
		"dit\u202F:",     // NNBSP before colon
		"«\u202FC'est",   // NNBSP after opening guillemet
		"fantastique\u202F!", // NNBSP before exclamation mark
		"»",              // Closing guillemet
		"vous\u202F;",    // NNBSP before semicolon
		"vraiment\u202F?", // NNBSP before question mark
	}
	
	for _, expected := range expectedTransformations {
		if !strings.Contains(outputStr, expected) {
			t.Errorf("French typography should transform text to contain '%s', but got: %s", expected, outputStr)
		}
	}
}

// TestCLIFrenchShortFlag tests the -f short flag version
func TestCLIFrenchShortFlag(t *testing.T) {
	// Create temporary input file
	tempDir, err := ioutil.TempDir("", "ankiprep_test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)
	
	inputFile := filepath.Join(tempDir, "input.csv")
	outputFile := filepath.Join(tempDir, "output.csv")
	
	// Write test CSV content
	csvContent := `front,back
"Bonjour : comment allez-vous ?","Hello: how are you?"`
	
	err = ioutil.WriteFile(inputFile, []byte(csvContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write input file: %v", err)
	}
	
	// Run ankiprep command with -f flag
	cmd := exec.Command("ankiprep", "-f", "-o", outputFile, inputFile)
	_, err = cmd.CombinedOutput()
	
	if err != nil {
		t.Fatalf("Expected -f processing command to succeed, got error: %v", err)
	}
	
	// Check output file content
	outputContent, err := ioutil.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}
	
	outputStr := string(outputContent)
	
	// Check for French typography transformation
	if !strings.Contains(outputStr, "Bonjour\u202F:") {
		t.Errorf("Short -f flag should apply French typography, but got: %s", outputStr)
	}
}

// TestCLIFrenchWithoutFlag tests that French typography is NOT applied without flag
func TestCLIFrenchWithoutFlag(t *testing.T) {
	// Create temporary input file
	tempDir, err := ioutil.TempDir("", "ankiprep_test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)
	
	inputFile := filepath.Join(tempDir, "input.csv")
	outputFile := filepath.Join(tempDir, "output.csv")
	
	// Write test CSV content
	csvContent := `front,back
"Bonjour : comment allez-vous ?","Hello: how are you?"`
	
	err = ioutil.WriteFile(inputFile, []byte(csvContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write input file: %v", err)
	}
	
	// Run ankiprep command WITHOUT French flag
	cmd := exec.Command("ankiprep", "-o", outputFile, inputFile)
	_, err = cmd.CombinedOutput()
	
	if err != nil {
		t.Fatalf("Expected processing command to succeed, got error: %v", err)
	}
	
	// Check output file content
	outputContent, err := ioutil.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}
	
	outputStr := string(outputContent)
	
	// Check that French typography is NOT applied (original colon spacing preserved)
	if !strings.Contains(outputStr, "Bonjour : comment") {
		t.Errorf("Without French flag, original spacing should be preserved, but got: %s", outputStr)
	}
	
	// Make sure NNBSP was NOT inserted
	if strings.Contains(outputStr, "Bonjour\u202F:") {
		t.Errorf("Without French flag, NNBSP should NOT be inserted, but got: %s", outputStr)
	}
}