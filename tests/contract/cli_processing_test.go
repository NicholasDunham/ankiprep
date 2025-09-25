package contract

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// TestCLIFileProcessingBasic tests basic file processing functionality
func TestCLIFileProcessingBasic(t *testing.T) {
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
hello,bonjour
goodbye,au revoir`
	
	err = ioutil.WriteFile(inputFile, []byte(csvContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write input file: %v", err)
	}
	
	// Run ankiprep command
	cmd := exec.Command("ankiprep", "-o", outputFile, inputFile)
	output, err := cmd.CombinedOutput()
	
	if err != nil {
		t.Fatalf("Expected processing command to succeed, got error: %v, output: %s", err, string(output))
	}
	
	// Check that output file was created
	if _, err := os.Stat(outputFile); os.IsNotExist(err) {
		t.Errorf("Expected output file %s to be created", outputFile)
	}
	
	// Check output file content
	outputContent, err := ioutil.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}
	
	outputStr := string(outputContent)
	
	// Check for Anki-specific headers
	expectedHeaders := []string{
		"#separator:comma",
		"#html:true",
		"#columns:front,back",
	}
	
	for _, expected := range expectedHeaders {
		if !strings.Contains(outputStr, expected) {
			t.Errorf("Output should contain Anki header '%s', but got: %s", expected, outputStr)
		}
	}
	
	// Check for data content
	if !strings.Contains(outputStr, "hello,bonjour") {
		t.Errorf("Output should contain data 'hello,bonjour', but got: %s", outputStr)
	}
}

// TestCLIFileProcessingMultiple tests processing multiple input files
func TestCLIFileProcessingMultiple(t *testing.T) {
	// Create temporary input files
	tempDir, err := ioutil.TempDir("", "ankiprep_test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)
	
	inputFile1 := filepath.Join(tempDir, "input1.csv")
	inputFile2 := filepath.Join(tempDir, "input2.csv")
	outputFile := filepath.Join(tempDir, "output.csv")
	
	// Write test CSV content
	csvContent1 := `front,back
hello,bonjour`
	
	csvContent2 := `front,back,extra
goodbye,au revoir,additional`
	
	err = ioutil.WriteFile(inputFile1, []byte(csvContent1), 0644)
	if err != nil {
		t.Fatalf("Failed to write input file 1: %v", err)
	}
	
	err = ioutil.WriteFile(inputFile2, []byte(csvContent2), 0644)
	if err != nil {
		t.Fatalf("Failed to write input file 2: %v", err)
	}
	
	// Run ankiprep command with multiple files
	cmd := exec.Command("ankiprep", "-o", outputFile, inputFile1, inputFile2)
	output, err := cmd.CombinedOutput()
	
	if err != nil {
		t.Fatalf("Expected processing command to succeed, got error: %v, output: %s", err, string(output))
	}
	
	// Check output file content
	outputContent, err := ioutil.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}
	
	outputStr := string(outputContent)
	
	// Check for merged headers (union of all columns)
	if !strings.Contains(outputStr, "#columns:front,back,extra") {
		t.Errorf("Output should contain merged columns header, but got: %s", outputStr)
	}
	
	// Check for both data entries
	if !strings.Contains(outputStr, "hello,bonjour") {
		t.Errorf("Output should contain first file data, but got: %s", outputStr)
	}
	if !strings.Contains(outputStr, "goodbye,au revoir") {
		t.Errorf("Output should contain second file data, but got: %s", outputStr)
	}
}

// TestCLIProcessingProgress tests verbose output and progress reporting
func TestCLIProcessingProgress(t *testing.T) {
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
hello,bonjour
hello,bonjour
goodbye,au revoir`
	
	err = ioutil.WriteFile(inputFile, []byte(csvContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write input file: %v", err)
	}
	
	// Run ankiprep command with verbose flag
	cmd := exec.Command("ankiprep", "-v", "-o", outputFile, inputFile)
	output, err := cmd.CombinedOutput()
	
	if err != nil {
		t.Fatalf("Expected verbose processing command to succeed, got error: %v", err)
	}
	
	outputStr := string(output)
	
	// Check for verbose output indicators
	expectedMessages := []string{
		"Processing",
		"input files",
		"Merging headers",
		"Processing records",
		"Removing duplicates",
		"Writing output",
		"Done",
	}
	
	for _, expected := range expectedMessages {
		if !strings.Contains(outputStr, expected) {
			t.Errorf("Verbose output should contain '%s', but got: %s", expected, outputStr)
		}
	}
}