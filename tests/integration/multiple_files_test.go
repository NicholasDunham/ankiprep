package integration

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// TestMultipleFilesMerging tests merging multiple CSV files with different columns
func TestMultipleFilesMerging(t *testing.T) {
	// Create temporary directory
	tempDir, err := ioutil.TempDir("", "ankiprep_integration")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)
	
	// Create first input file
	inputFile1 := filepath.Join(tempDir, "input1.csv")
	csvContent1 := `front,back
hello,bonjour
goodbye,au revoir`
	
	err = ioutil.WriteFile(inputFile1, []byte(csvContent1), 0644)
	if err != nil {
		t.Fatalf("Failed to write input file 1: %v", err)
	}
	
	// Create second input file with additional column
	inputFile2 := filepath.Join(tempDir, "input2.csv")
	csvContent2 := `front,back,extra
good morning,bonjour,greeting
good night,bonne nuit,farewell`
	
	err = ioutil.WriteFile(inputFile2, []byte(csvContent2), 0644)
	if err != nil {
		t.Fatalf("Failed to write input file 2: %v", err)
	}
	
	// Create third input file with different column order
	inputFile3 := filepath.Join(tempDir, "input3.csv")
	csvContent3 := `extra,front,back,category
noun,cat,chat,animal
noun,dog,chien,animal`
	
	err = ioutil.WriteFile(inputFile3, []byte(csvContent3), 0644)
	if err != nil {
		t.Fatalf("Failed to write input file 3: %v", err)
	}
	
	outputFile := filepath.Join(tempDir, "merged.csv")
	
	// Execute ankiprep with multiple files
	cmd := exec.Command("ankiprep", "-o", outputFile, inputFile1, inputFile2, inputFile3)
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
	
	// Verify merged headers (union of all columns)
	expectedColumns := "#columns:front,back,extra,category"
	if !strings.Contains(outputStr, expectedColumns) {
		t.Errorf("Expected merged columns header: %s\nGot: %s", expectedColumns, outputStr)
	}
	
	// Verify all data is present
	expectedData := []string{
		"hello,bonjour",
		"goodbye,au revoir",
		"good morning,bonjour,greeting",
		"good night,bonne nuit,farewell",
		"cat,chat",
		"dog,chien",
	}
	
	for _, data := range expectedData {
		if !strings.Contains(outputStr, data) {
			t.Errorf("Output missing expected data: %s\nGot: %s", data, outputStr)
		}
	}
}

// TestMultipleFilesHeaderOrder tests that column order is properly maintained
func TestMultipleFilesHeaderOrder(t *testing.T) {
	// Create temporary directory
	tempDir, err := ioutil.TempDir("", "ankiprep_integration")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)
	
	// Create files with different column orders
	inputFile1 := filepath.Join(tempDir, "file1.csv")
	csvContent1 := `alpha,beta,gamma
1,2,3`
	
	err = ioutil.WriteFile(inputFile1, []byte(csvContent1), 0644)
	if err != nil {
		t.Fatalf("Failed to write input file 1: %v", err)
	}
	
	inputFile2 := filepath.Join(tempDir, "file2.csv")
	csvContent2 := `gamma,delta,alpha
30,40,10`
	
	err = ioutil.WriteFile(inputFile2, []byte(csvContent2), 0644)
	if err != nil {
		t.Fatalf("Failed to write input file 2: %v", err)
	}
	
	outputFile := filepath.Join(tempDir, "merged.csv")
	
	// Execute ankiprep
	cmd := exec.Command("ankiprep", "-o", outputFile, inputFile1, inputFile2)
	_, err = cmd.CombinedOutput()
	
	if err != nil {
		t.Fatalf("ankiprep command failed: %v", err)
	}
	
	// Read output content
	outputContent, err := ioutil.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}
	
	outputStr := string(outputContent)
	
	// Should contain all unique columns
	requiredColumns := []string{"alpha", "beta", "gamma", "delta"}
	columnsHeader := ""
	
	// Extract the columns header line
	lines := strings.Split(outputStr, "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "#columns:") {
			columnsHeader = line
			break
		}
	}
	
	if columnsHeader == "" {
		t.Fatalf("Could not find #columns header in output")
	}
	
	// Verify all columns are present
	for _, col := range requiredColumns {
		if !strings.Contains(columnsHeader, col) {
			t.Errorf("Column '%s' missing from header: %s", col, columnsHeader)
		}
	}
}

// TestMultipleFilesMixedSeparators tests handling files with different separators
func TestMultipleFilesMixedSeparators(t *testing.T) {
	// Create temporary directory
	tempDir, err := ioutil.TempDir("", "ankiprep_integration")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)
	
	// Create CSV file
	inputFile1 := filepath.Join(tempDir, "data.csv")
	csvContent := `front,back
hello,world`
	
	err = ioutil.WriteFile(inputFile1, []byte(csvContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write CSV file: %v", err)
	}
	
	// Create TSV file
	inputFile2 := filepath.Join(tempDir, "data.tsv")
	tsvContent := "front\tback\ntest\tessai"
	
	err = ioutil.WriteFile(inputFile2, []byte(tsvContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write TSV file: %v", err)
	}
	
	outputFile := filepath.Join(tempDir, "merged.csv")
	
	// Execute ankiprep with mixed separator files
	cmd := exec.Command("ankiprep", "-o", outputFile, inputFile1, inputFile2)
	_, err = cmd.CombinedOutput()
	
	if err != nil {
		t.Fatalf("ankiprep command failed: %v", err)
	}
	
	// Read output content
	outputContent, err := ioutil.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}
	
	outputStr := string(outputContent)
	
	// Should contain data from both files
	if !strings.Contains(outputStr, "hello,world") {
		t.Errorf("Missing CSV data in merged output: %s", outputStr)
	}
	
	if !strings.Contains(outputStr, "test,essai") {
		t.Errorf("Missing TSV data (should be converted to CSV) in merged output: %s", outputStr)
	}
	
	// Output should always be comma-separated
	if !strings.Contains(outputStr, "#separator:comma") {
		t.Errorf("Output should always use comma separator: %s", outputStr)
	}
}