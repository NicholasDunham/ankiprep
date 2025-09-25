package integration

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// TestSingleFileProcessingBasic tests processing a single CSV file with basic data
func TestSingleFileProcessingBasic(t *testing.T) {
	// Create temporary directory
	tempDir, err := ioutil.TempDir("", "ankiprep_integration")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)
	
	// Create input file
	inputFile := filepath.Join(tempDir, "input.csv")
	outputFile := filepath.Join(tempDir, "output.csv")
	
	csvContent := `front,back,extra
hello,bonjour,greeting
goodbye,au revoir,farewell
yes,oui,affirmation`
	
	err = ioutil.WriteFile(inputFile, []byte(csvContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write input file: %v", err)
	}
	
	// Execute ankiprep
	cmd := exec.Command("ankiprep", "-o", outputFile, inputFile)
	output, err := cmd.CombinedOutput()
	
	if err != nil {
		t.Fatalf("ankiprep command failed: %v, output: %s", err, string(output))
	}
	
	// Verify output file exists
	if _, err := os.Stat(outputFile); os.IsNotExist(err) {
		t.Fatalf("Output file was not created")
	}
	
	// Read and verify output content
	outputContent, err := ioutil.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}
	
	outputStr := string(outputContent)
	
	// Verify Anki headers are present
	requiredHeaders := []string{
		"#separator:comma",
		"#html:true",
		"#columns:front,back,extra",
	}
	
	for _, header := range requiredHeaders {
		if !strings.Contains(outputStr, header) {
			t.Errorf("Output missing required Anki header: %s\nGot: %s", header, outputStr)
		}
	}
	
	// Verify data content is preserved
	expectedData := []string{
		"hello,bonjour,greeting",
		"goodbye,au revoir,farewell",
		"yes,oui,affirmation",
	}
	
	for _, data := range expectedData {
		if !strings.Contains(outputStr, data) {
			t.Errorf("Output missing expected data: %s\nGot: %s", data, outputStr)
		}
	}
}

// TestSingleFileProcessingTSV tests processing a tab-separated file
func TestSingleFileProcessingTSV(t *testing.T) {
	// Create temporary directory
	tempDir, err := ioutil.TempDir("", "ankiprep_integration")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)
	
	// Create TSV input file
	inputFile := filepath.Join(tempDir, "input.tsv")
	outputFile := filepath.Join(tempDir, "output.csv")
	
	tsvContent := "front\tback\nexample\texemple\ntest\tessai"
	
	err = ioutil.WriteFile(inputFile, []byte(tsvContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write input file: %v", err)
	}
	
	// Execute ankiprep
	cmd := exec.Command("ankiprep", "-o", outputFile, inputFile)
	_, err = cmd.CombinedOutput()
	
	if err != nil {
		t.Fatalf("ankiprep command failed: %v", err)
	}
	
	// Read and verify output content
	outputContent, err := ioutil.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}
	
	outputStr := string(outputContent)
	
	// TSV input should be converted to CSV output
	if !strings.Contains(outputStr, "#separator:comma") {
		t.Errorf("TSV input should be converted to CSV output format")
	}
	
	// Data should be converted from tabs to commas
	if !strings.Contains(outputStr, "example,exemple") {
		t.Errorf("TSV data should be converted to CSV format, got: %s", outputStr)
	}
}

// TestSingleFileProcessingWithSmartQuotes tests smart quote conversion
func TestSingleFileProcessingWithSmartQuotes(t *testing.T) {
	// Create temporary directory
	tempDir, err := ioutil.TempDir("", "ankiprep_integration")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)
	
	// Create input file with straight quotes
	inputFile := filepath.Join(tempDir, "input.csv")
	outputFile := filepath.Join(tempDir, "output.csv")
	
	csvContent := `front,back
"He said ""Hello"" to me","Il m'a dit ""Bonjour"""`
	
	err = ioutil.WriteFile(inputFile, []byte(csvContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write input file: %v", err)
	}
	
	// Execute ankiprep
	cmd := exec.Command("ankiprep", "-o", outputFile, inputFile)
	_, err = cmd.CombinedOutput()
	
	if err != nil {
		t.Fatalf("ankiprep command failed: %v", err)
	}
	
	// Read and verify output content
	outputContent, err := ioutil.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}
	
	outputStr := string(outputContent)
	
	// Check for smart quotes conversion (straight " to smart quotes)
	if strings.Contains(outputStr, `"Hello"`) && !strings.Contains(outputStr, "\u201cHello\u201d") {
		t.Errorf("Smart quotes should be converted from straight quotes, got: %s", outputStr)
	}
}