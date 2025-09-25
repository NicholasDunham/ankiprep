package performance

import (
	"ankiprep/internal/app"
	"ankiprep/internal/models"
	"ankiprep/internal/services"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"
)

// TestLargeFileProcessing tests performance with large CSV files
func TestLargeFileProcessing(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping large file performance tests in short mode")
	}

	tests := []struct {
		name            string
		numRows         int
		numColumns      int
		maxTimeSeconds  float64
		maxMemoryMB     float64
	}{
		{
			name:            "small file (1K rows)",
			numRows:         1000,
			numColumns:      5,
			maxTimeSeconds:  2.0,
			maxMemoryMB:     50.0,
		},
		{
			name:            "medium file (10K rows)",
			numRows:         10000,
			numColumns:      10,
			maxTimeSeconds:  10.0,
			maxMemoryMB:     100.0,
		},
		{
			name:            "large file (100K rows)",
			numRows:         100000,
			numColumns:      10,
			maxTimeSeconds:  60.0,
			maxMemoryMB:     500.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temporary test file
			testFile := createLargeTempCSV(t, tt.numRows, tt.numColumns)
			defer os.Remove(testFile)

			// Measure initial memory
			var startMem runtime.MemStats
			runtime.GC()
			runtime.ReadMemStats(&startMem)

			// Start timing
			startTime := time.Now()

			// Process the file
			config := app.ProcessorConfig{
				FrenchMode:     false,
				SmartQuotes:    false,
				SkipDuplicates: false,
				Verbose:        false,
			}
			processor := app.NewProcessor(config)
			
			inputPaths := []string{testFile}
			result, err := processor.ProcessFiles(inputPaths)

			// End timing
			elapsed := time.Since(startTime).Seconds()

			// Measure final memory
			var endMem runtime.MemStats
			runtime.GC()
			runtime.ReadMemStats(&endMem)

			// Calculate memory usage in MB
			memoryUsedMB := float64(endMem.Alloc-startMem.Alloc) / (1024 * 1024)

			// Check for errors
			if err != nil {
				t.Fatalf("Processing failed: %v", err)
			}

			if result == nil {
				t.Fatal("Expected non-nil processing result")
			}

			// Performance assertions
			if elapsed > tt.maxTimeSeconds {
				t.Errorf("Processing took too long: %.2f seconds (max: %.2f)", elapsed, tt.maxTimeSeconds)
			}

			if memoryUsedMB > tt.maxMemoryMB {
				t.Errorf("Memory usage too high: %.2f MB (max: %.2f)", memoryUsedMB, tt.maxMemoryMB)
			}

			// Verify expected number of processed entries (minus header)
			if result.OutputRecords < tt.numRows-1 {
				t.Errorf("Expected at least %d processed entries, got %d", tt.numRows-1, result.OutputRecords)
			}

			t.Logf("Performance stats - Time: %.2fs, Memory: %.2fMB, Entries: %d", 
				elapsed, memoryUsedMB, result.OutputRecords)
		})
	}
}

// TestCSVParserPerformance tests CSV parsing performance specifically
func TestCSVParserPerformance(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping CSV parser performance tests in short mode")
	}

	tests := []struct {
		name           string
		numRows        int
		numColumns     int
		maxTimeSeconds float64
	}{
		{
			name:           "CSV parsing - 50K rows",
			numRows:        50000,
			numColumns:     8,
			maxTimeSeconds: 30.0,
		},
		{
			name:           "CSV parsing - 100K rows",
			numRows:        100000,
			numColumns:     12,
			maxTimeSeconds: 60.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temporary test file
			testFile := createLargeTempCSV(t, tt.numRows, tt.numColumns)
			defer os.Remove(testFile)

			parser := services.NewCSVParser()

			startTime := time.Now()

			// Parse file to InputFile model
			inputFile, err := parser.ParseFile(testFile)
			if err != nil {
				t.Fatalf("ParseFile failed: %v", err)
			}

			// Convert to data entries
			entries, err := parser.ParseToDataEntries(inputFile)
			if err != nil {
				t.Fatalf("ParseToDataEntries failed: %v", err)
			}

			elapsed := time.Since(startTime).Seconds()

			// Performance assertions
			if elapsed > tt.maxTimeSeconds {
				t.Errorf("CSV parsing took too long: %.2f seconds (max: %.2f)", elapsed, tt.maxTimeSeconds)
			}

			// Verify reasonable number of entries (header row excluded)
			expectedEntries := tt.numRows - 1
			if len(entries) < expectedEntries-100 { // Allow small variance for empty rows
				t.Errorf("Expected around %d entries, got %d", expectedEntries, len(entries))
			}

			t.Logf("CSV parsing stats - Time: %.2fs, Records: %d, Entries: %d", 
				elapsed, len(inputFile.Records), len(entries))
		})
	}
}

// TestDuplicateDetectionPerformance tests performance of duplicate detection
func TestDuplicateDetectionPerformance(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping duplicate detection performance tests in short mode")
	}

	tests := []struct {
		name              string
		numEntries        int
		duplicatePercent  float64
		maxTimeSeconds    float64
	}{
		{
			name:             "small dataset with 10% duplicates",
			numEntries:       10000,
			duplicatePercent: 0.1,
			maxTimeSeconds:   5.0,
		},
		{
			name:             "large dataset with 25% duplicates",
			numEntries:       50000,
			duplicatePercent: 0.25,
			maxTimeSeconds:   20.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test entries with known duplicates
			entries := createTestEntriesWithDuplicates(tt.numEntries, tt.duplicatePercent)

			detector := services.NewDuplicateDetector()
			startTime := time.Now()

			// Detect duplicates
			uniqueEntries, duplicateCount := detector.DetectDuplicates(entries)

			elapsed := time.Since(startTime).Seconds()

			// Performance assertions
			if elapsed > tt.maxTimeSeconds {
				t.Errorf("Duplicate detection took too long: %.2f seconds (max: %.2f)", elapsed, tt.maxTimeSeconds)
			}

			// Verify reasonable duplicate detection
			expectedDuplicates := int(float64(tt.numEntries) * tt.duplicatePercent)
			if duplicateCount < expectedDuplicates-100 || duplicateCount > expectedDuplicates+100 {
				t.Errorf("Expected around %d duplicates, got %d", expectedDuplicates, duplicateCount)
			}

			expectedUnique := tt.numEntries - duplicateCount
			if len(uniqueEntries) != expectedUnique {
				t.Errorf("Expected %d unique entries, got %d", expectedUnique, len(uniqueEntries))
			}

			t.Logf("Duplicate detection stats - Time: %.2fs, Total: %d, Unique: %d, Duplicates: %d", 
				elapsed, tt.numEntries, len(uniqueEntries), duplicateCount)
		})
	}
}

// TestTypographyPerformance tests typography processing performance
func TestTypographyPerformance(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping typography performance tests in short mode")
	}

	tests := []struct {
		name           string
		numEntries     int
		avgTextLength  int
		maxTimeSeconds float64
		frenchMode     bool
		smartQuotes    bool
	}{
		{
			name:           "basic typography - 10K entries",
			numEntries:     10000,
			avgTextLength:  100,
			maxTimeSeconds: 10.0,
			frenchMode:     false,
			smartQuotes:    true,
		},
		{
			name:           "french typography - 10K entries",
			numEntries:     10000,
			avgTextLength:  150,
			maxTimeSeconds: 15.0,
			frenchMode:     true,
			smartQuotes:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test entries with varied text content
			entries := createTestEntriesWithText(tt.numEntries, tt.avgTextLength)

			service := services.NewTypographyService(tt.frenchMode, tt.smartQuotes)
			startTime := time.Now()

			// Process entries
			processedEntries := service.ProcessEntries(entries)

			elapsed := time.Since(startTime).Seconds()

			// Performance assertions
			if elapsed > tt.maxTimeSeconds {
				t.Errorf("Typography processing took too long: %.2f seconds (max: %.2f)", elapsed, tt.maxTimeSeconds)
			}

			// Verify all entries were processed
			if len(processedEntries) != len(entries) {
				t.Errorf("Expected %d processed entries, got %d", len(entries), len(processedEntries))
			}

			// Verify processing actually occurred (entries should be different objects)
			for i, original := range entries {
				if processedEntries[i] == original {
					t.Errorf("Entry %d was not processed (same reference)", i)
				}
			}

			t.Logf("Typography processing stats - Time: %.2fs, Entries: %d, Mode: french=%v, quotes=%v", 
				elapsed, len(processedEntries), tt.frenchMode, tt.smartQuotes)
		})
	}
}

// TestMemoryLeaks tests for memory leaks in large dataset processing
func TestMemoryLeaks(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping memory leak tests in short mode")
	}

	const iterations = 10
	const rowsPerIteration = 5000

	var memStats []uint64

	for i := 0; i < iterations; i++ {
		// Create test file
		testFile := createLargeTempCSV(t, rowsPerIteration, 5)

		// Process file
		config := app.ProcessorConfig{
			FrenchMode:     false,
			SmartQuotes:    true,
			SkipDuplicates: true,
			Verbose:        false,
		}
		processor := app.NewProcessor(config)
		inputPaths := []string{testFile}
		
		_, err := processor.ProcessFiles(inputPaths)
		if err != nil {
			t.Fatalf("Processing iteration %d failed: %v", i, err)
		}

		// Clean up file
		os.Remove(testFile)

		// Force garbage collection and measure memory
		runtime.GC()
		runtime.GC() // Call twice to be thorough

		var m runtime.MemStats
		runtime.ReadMemStats(&m)
		memStats = append(memStats, m.Alloc)
	}

	// Check for memory growth trend
	firstHalf := memStats[:iterations/2]
	secondHalf := memStats[iterations/2:]

	var firstHalfAvg, secondHalfAvg uint64
	for _, mem := range firstHalf {
		firstHalfAvg += mem
	}
	firstHalfAvg /= uint64(len(firstHalf))

	for _, mem := range secondHalf {
		secondHalfAvg += mem
	}
	secondHalfAvg /= uint64(len(secondHalf))

	// Allow some growth but not excessive
	growthRatio := float64(secondHalfAvg) / float64(firstHalfAvg)
	if growthRatio > 1.5 {
		t.Errorf("Potential memory leak detected: memory grew by %.2fx over %d iterations", growthRatio, iterations)
		t.Logf("Memory stats: %v", memStats)
	}

	t.Logf("Memory leak test completed - Growth ratio: %.2fx", growthRatio)
}

// Helper function to create large temporary CSV files
func createLargeTempCSV(t *testing.T, numRows, numColumns int) string {
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "large_test.csv")

	file, err := os.Create(testFile)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	defer file.Close()

	// Write header
	header := ""
	for i := 0; i < numColumns; i++ {
		if i > 0 {
			header += ","
		}
		header += fmt.Sprintf("column_%d", i)
	}
	header += "\n"
	file.WriteString(header)

	// Write data rows
	for row := 0; row < numRows; row++ {
		line := ""
		for col := 0; col < numColumns; col++ {
			if col > 0 {
				line += ","
			}
			line += fmt.Sprintf("data_%d_%d", row, col)
		}
		line += "\n"
		file.WriteString(line)
	}

	return testFile
}

// Helper function to create test entries with known duplicates
func createTestEntriesWithDuplicates(numEntries int, duplicatePercent float64) []*models.DataEntry {
	entries := make([]*models.DataEntry, numEntries)
	numDuplicates := int(float64(numEntries) * duplicatePercent)

	// Create unique entries
	for i := 0; i < numEntries-numDuplicates; i++ {
		values := map[string]string{
			"front": fmt.Sprintf("unique_front_%d", i),
			"back":  fmt.Sprintf("unique_back_%d", i),
		}
		entries[i] = models.NewDataEntry(values, "test.csv", i+1)
	}

	// Create duplicates by copying some of the unique entries
	dupIndex := 0
	for i := numEntries - numDuplicates; i < numEntries; i++ {
		// Copy from earlier entries
		sourceIndex := dupIndex % (numEntries - numDuplicates)
		if sourceIndex < len(entries) && entries[sourceIndex] != nil {
			values := map[string]string{
				"front": entries[sourceIndex].Values["front"],
				"back":  entries[sourceIndex].Values["back"],
			}
			entries[i] = models.NewDataEntry(values, "test.csv", i+1)
		}
		dupIndex++
	}

	return entries
}

// Helper function to create test entries with varied text content
func createTestEntriesWithText(numEntries, avgTextLength int) []*models.DataEntry {
	entries := make([]*models.DataEntry, numEntries)

	for i := 0; i < numEntries; i++ {
		front := generateTestText(avgTextLength, i)
		back := generateTestText(avgTextLength, i+numEntries)
		
		values := map[string]string{
			"front": front,
			"back":  back,
		}
		entries[i] = models.NewDataEntry(values, "test.csv", i+1)
	}

	return entries
}

// Helper function to generate test text with quotes and punctuation
func generateTestText(length, seed int) string {
	texts := []string{
		"This is a sample text with \"quoted content\" and punctuation!",
		"Here's another example: \"How are you?\" she asked.",
		"French example: « Bonjour ! Comment allez-vous ? »",
		"Mixed quotes: 'single' and \"double\" quotes together.",
		"Punctuation test: Hello! How are you? Fine, thanks.",
		"Numbers and symbols: 123, $45.67, 89% complete.",
	}

	baseText := texts[seed%len(texts)]
	
	// Extend text to approximate desired length
	result := baseText
	for len(result) < length {
		result += " " + texts[(seed+len(result))%len(texts)]
	}

	// Truncate to desired length
	if len(result) > length {
		result = result[:length]
	}

	return result
}