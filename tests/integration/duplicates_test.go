package integration

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// TestDuplicateDetectionExact tests exact duplicate detection and removal
func TestDuplicateDetectionExact(t *testing.T) {
	// Create temporary directory
	tempDir, err := ioutil.TempDir("", "ankiprep_integration")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)
	
	// Create input file with exact duplicates
	inputFile := filepath.Join(tempDir, "input.csv")
	csvContent := `front,back
hello,bonjour
goodbye,au revoir
hello,bonjour
test,essai
hello,bonjour
good night,bonne nuit`
	
	err = ioutil.WriteFile(inputFile, []byte(csvContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write input file: %v", err)
	}
	
	outputFile := filepath.Join(tempDir, "output.csv")
	
	// Execute ankiprep with verbose output to see duplicate reporting
	cmd := exec.Command("ankiprep", "-v", "-o", outputFile, inputFile)
	output, err := cmd.CombinedOutput()
	
	if err != nil {
		t.Fatalf("ankiprep command failed: %v, output: %s", err, string(output))
	}
	
	// Check verbose output for duplicate reporting
	outputStr := string(output)
	if !strings.Contains(outputStr, "duplicate") {
		t.Errorf("Verbose output should mention duplicate detection: %s", outputStr)
	}
	
	// Read output file
	outputContent, err := ioutil.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}
	
	outputFileStr := string(outputContent)
	
	// Count occurrences of the duplicate entry
	helloCount := strings.Count(outputFileStr, "hello,bonjour")
	if helloCount != 1 {
		t.Errorf("Expected exactly 1 occurrence of 'hello,bonjour' after deduplication, got %d", helloCount)
	}
	
	// Verify other unique entries are still present
	expectedEntries := []string{
		"goodbye,au revoir",
		"test,essai",
		"good night,bonne nuit",
	}
	
	for _, entry := range expectedEntries {
		if !strings.Contains(outputFileStr, entry) {
			t.Errorf("Expected entry missing after deduplication: %s", entry)
		}
	}
}

// TestDuplicateDetectionCaseSensitive tests that duplicate detection is case-sensitive
func TestDuplicateDetectionCaseSensitive(t *testing.T) {
	// Create temporary directory
	tempDir, err := ioutil.TempDir("", "ankiprep_integration")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)
	
	// Create input file with case variations
	inputFile := filepath.Join(tempDir, "input.csv")
	csvContent := `front,back
hello,bonjour
Hello,Bonjour
HELLO,BONJOUR
hello,bonjour`
	
	err = ioutil.WriteFile(inputFile, []byte(csvContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write input file: %v", err)
	}
	
	outputFile := filepath.Join(tempDir, "output.csv")
	
	// Execute ankiprep
	cmd := exec.Command("ankiprep", "-o", outputFile, inputFile)
	_, err = cmd.CombinedOutput()
	
	if err != nil {
		t.Fatalf("ankiprep command failed: %v", err)
	}
	
	// Read output file
	outputContent, err := ioutil.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}
	
	outputStr := string(outputContent)
	
	// All case variations should be treated as different entries except exact duplicates
	expectedEntries := []string{
		"hello,bonjour",  // Only one instance of this exact duplicate
		"Hello,Bonjour",  // Different case, should be preserved
		"HELLO,BONJOUR",  // Different case, should be preserved
	}
	
	for _, entry := range expectedEntries {
		if !strings.Contains(outputStr, entry) {
			t.Errorf("Expected case-sensitive entry: %s", entry)
		}
	}
	
	// Exact duplicate should appear only once
	exactDuplicateCount := strings.Count(outputStr, "hello,bonjour")
	if exactDuplicateCount != 1 {
		t.Errorf("Expected exactly 1 occurrence of exact duplicate 'hello,bonjour', got %d", exactDuplicateCount)
	}
}

// TestDuplicateDetectionAcrossFiles tests duplicate detection across multiple input files
func TestDuplicateDetectionAcrossFiles(t *testing.T) {
	// Create temporary directory
	tempDir, err := ioutil.TempDir("", "ankiprep_integration")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)
	
	// Create first input file
	inputFile1 := filepath.Join(tempDir, "file1.csv")
	csvContent1 := `front,back
hello,bonjour
goodbye,au revoir`
	
	err = ioutil.WriteFile(inputFile1, []byte(csvContent1), 0644)
	if err != nil {
		t.Fatalf("Failed to write input file 1: %v", err)
	}
	
	// Create second input file with some duplicates
	inputFile2 := filepath.Join(tempDir, "file2.csv")
	csvContent2 := `front,back
hello,bonjour
good morning,bonjour
test,essai`
	
	err = ioutil.WriteFile(inputFile2, []byte(csvContent2), 0644)
	if err != nil {
		t.Fatalf("Failed to write input file 2: %v", err)
	}
	
	// Create third input file with more duplicates
	inputFile3 := filepath.Join(tempDir, "file3.csv")
	csvContent3 := `front,back
hello,bonjour
new entry,nouvelle entrée
test,essai`
	
	err = ioutil.WriteFile(inputFile3, []byte(csvContent3), 0644)
	if err != nil {
		t.Fatalf("Failed to write input file 3: %v", err)
	}
	
	outputFile := filepath.Join(tempDir, "merged.csv")
	
	// Execute ankiprep with verbose output
	cmd := exec.Command("ankiprep", "-v", "-o", outputFile, inputFile1, inputFile2, inputFile3)
	output, err := cmd.CombinedOutput()
	
	if err != nil {
		t.Fatalf("ankiprep command failed: %v, output: %s", err, string(output))
	}
	
	// Read output file
	outputContent, err := ioutil.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}
	
	outputStr := string(outputContent)
	
	// Check that duplicates across files are removed
	duplicateEntries := []string{
		"hello,bonjour", // appears in all 3 files
		"test,essai",    // appears in files 2 and 3
	}
	
	for _, entry := range duplicateEntries {
		count := strings.Count(outputStr, entry)
		if count != 1 {
			t.Errorf("Expected exactly 1 occurrence of cross-file duplicate '%s', got %d", entry, count)
		}
	}
	
	// Check that unique entries are preserved
	uniqueEntries := []string{
		"goodbye,au revoir",
		"good morning,bonjour",
		"new entry,nouvelle entrée",
	}
	
	for _, entry := range uniqueEntries {
		if !strings.Contains(outputStr, entry) {
			t.Errorf("Expected unique entry to be preserved: %s", entry)
		}
	}
}

// TestDuplicateDetectionWithEmptyFields tests duplicate detection when fields contain empty values
func TestDuplicateDetectionWithEmptyFields(t *testing.T) {
	// Create temporary directory
	tempDir, err := ioutil.TempDir("", "ankiprep_integration")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)
	
	// Create input file with empty fields
	inputFile := filepath.Join(tempDir, "input.csv")
	csvContent := `front,back,extra
hello,,
hello,,
test,essai,
hello,bonjour,
,empty,
,empty,`
	
	err = ioutil.WriteFile(inputFile, []byte(csvContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write input file: %v", err)
	}
	
	outputFile := filepath.Join(tempDir, "output.csv")
	
	// Execute ankiprep
	cmd := exec.Command("ankiprep", "-o", outputFile, inputFile)
	_, err = cmd.CombinedOutput()
	
	if err != nil {
		t.Fatalf("ankiprep command failed: %v", err)
	}
	
	// Read output file
	outputContent, err := ioutil.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}
	
	outputStr := string(outputContent)
	
	// Entries with empty fields should be treated as distinct from those with values
	// but exact duplicates (including empty fields) should be removed
	entriesWithCounts := map[string]int{
		"hello,,":        1, // Should appear only once (duplicate removed)
		"test,essai,":    1, // Unique
		"hello,bonjour,": 1, // Unique (different from hello,,)
		",empty,":        1, // Should appear only once (duplicate removed)
	}
	
	for entry, expectedCount := range entriesWithCounts {
		actualCount := strings.Count(outputStr, entry)
		if actualCount != expectedCount {
			t.Errorf("Expected %d occurrence(s) of '%s', got %d", expectedCount, entry, actualCount)
		}
	}
}