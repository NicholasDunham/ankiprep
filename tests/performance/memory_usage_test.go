package performance

import (
	"ankiprep/internal/app"
	"ankiprep/internal/services"
	"os"
	"runtime"
	"testing"
)

// TestMemoryUsageScaling tests how memory usage scales with data size
func TestMemoryUsageScaling(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping memory usage scaling tests in short mode")
	}

	tests := []struct {
		name              string
		numRows           int
		maxMemoryMB       float64
		expectedLinearMB  float64 // Expected memory usage for linear scaling baseline
	}{
		{
			name:             "baseline - 1K rows",
			numRows:          1000,
			maxMemoryMB:      20.0,
			expectedLinearMB: 2.0,
		},
		{
			name:             "scaling - 10K rows",
			numRows:          10000,
			maxMemoryMB:      100.0,
			expectedLinearMB: 20.0,
		},
		{
			name:             "scaling - 50K rows",
			numRows:          50000,
			maxMemoryMB:      400.0,
			expectedLinearMB: 100.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test file
			testFile := createLargeTempCSV(t, tt.numRows, 8)
			defer os.Remove(testFile)

			// Measure baseline memory
			runtime.GC()
			runtime.GC()
			var startMem runtime.MemStats
			runtime.ReadMemStats(&startMem)

			// Process the file
			config := app.ProcessorConfig{
				FrenchMode:     false,
				SmartQuotes:    true,
				SkipDuplicates: true,
				Verbose:        false,
			}
			processor := app.NewProcessor(config)
			
			_, err := processor.ProcessFiles([]string{testFile})
			if err != nil {
				t.Fatalf("Processing failed: %v", err)
			}

			// Measure peak memory
			runtime.GC()
			runtime.GC()
			var endMem runtime.MemStats
			runtime.ReadMemStats(&endMem)

			memoryUsedMB := float64(endMem.Alloc-startMem.Alloc) / (1024 * 1024)

			// Check memory usage is within acceptable bounds
			if memoryUsedMB > tt.maxMemoryMB {
				t.Errorf("Memory usage too high: %.2fMB (max: %.2fMB)", memoryUsedMB, tt.maxMemoryMB)
			}

			// Log memory efficiency compared to expected linear scaling
			efficiency := tt.expectedLinearMB / memoryUsedMB * 100
			if efficiency > 100 {
				efficiency = 100 // Cap at 100% for better results
			}

			t.Logf("Memory usage: %.2fMB (expected: %.2fMB, efficiency: %.1f%%, %d rows)", 
				memoryUsedMB, tt.expectedLinearMB, efficiency, tt.numRows)
		})
	}
}

// TestMemoryUsageByComponent tests memory usage of individual components
func TestMemoryUsageByComponent(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping component memory usage tests in short mode")
	}

	const numRows = 10000
	const numCols = 6
	testFile := createLargeTempCSV(t, numRows, numCols)
	defer os.Remove(testFile)

	t.Run("CSV parsing memory usage", func(t *testing.T) {
		runtime.GC()
		runtime.GC()
		var startMem runtime.MemStats
		runtime.ReadMemStats(&startMem)

		parser := services.NewCSVParser()
		inputFile, err := parser.ParseFile(testFile)
		if err != nil {
			t.Fatalf("ParseFile failed: %v", err)
		}

		_, err = parser.ParseToDataEntries(inputFile)
		if err != nil {
			t.Fatalf("ParseToDataEntries failed: %v", err)
		}

		runtime.GC()
		runtime.GC()
		var endMem runtime.MemStats
		runtime.ReadMemStats(&endMem)

		memoryUsedMB := float64(endMem.Alloc-startMem.Alloc) / (1024 * 1024)
		t.Logf("CSV parsing memory usage: %.2fMB for %d rows", memoryUsedMB, numRows)

		// Should be reasonable for CSV parsing
		if memoryUsedMB > 100.0 {
			t.Errorf("CSV parsing uses too much memory: %.2fMB", memoryUsedMB)
		}
	})

	t.Run("duplicate detection memory usage", func(t *testing.T) {
		// Create test entries first
		entries := createTestEntriesWithDuplicates(numRows, 0.2)

		runtime.GC()
		runtime.GC()
		var startMem runtime.MemStats
		runtime.ReadMemStats(&startMem)

		detector := services.NewDuplicateDetector()
		_, _ = detector.DetectDuplicates(entries)

		runtime.GC()
		runtime.GC()
		var endMem runtime.MemStats
		runtime.ReadMemStats(&endMem)

		memoryUsedMB := float64(endMem.Alloc-startMem.Alloc) / (1024 * 1024)
		t.Logf("Duplicate detection memory usage: %.2fMB for %d entries", memoryUsedMB, numRows)

		// Should be reasonable for duplicate detection
		if memoryUsedMB > 80.0 {
			t.Errorf("Duplicate detection uses too much memory: %.2fMB", memoryUsedMB)
		}
	})

	t.Run("typography processing memory usage", func(t *testing.T) {
		// Create test entries with text content
		entries := createTestEntriesWithText(numRows/4, 200) // Fewer entries but longer text

		runtime.GC()
		runtime.GC()
		var startMem runtime.MemStats
		runtime.ReadMemStats(&startMem)

		service := services.NewTypographyService(true, true)
		_ = service.ProcessEntries(entries)

		runtime.GC()
		runtime.GC()
		var endMem runtime.MemStats
		runtime.ReadMemStats(&endMem)

		memoryUsedMB := float64(endMem.Alloc-startMem.Alloc) / (1024 * 1024)
		t.Logf("Typography processing memory usage: %.2fMB for %d entries", memoryUsedMB, len(entries))

		// Should be reasonable for typography processing
		if memoryUsedMB > 60.0 {
			t.Errorf("Typography processing uses too much memory: %.2fMB", memoryUsedMB)
		}
	})
}

// TestMemoryGrowthPattern tests memory usage patterns over time
func TestMemoryGrowthPattern(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping memory growth pattern tests in short mode")
	}

	const batchSize = 5000
	const numBatches = 6
	var memoryReadings []float64

	detector := services.NewDuplicateDetector()

	for batch := 1; batch <= numBatches; batch++ {
		// Create test entries for this batch
		entries := createTestEntriesWithDuplicates(batchSize, 0.15)

		// Process entries
		_, _ = detector.DetectDuplicates(entries)

		// Measure memory after this batch
		runtime.GC()
		runtime.GC()
		var mem runtime.MemStats
		runtime.ReadMemStats(&mem)

		memoryMB := float64(mem.Alloc) / (1024 * 1024)
		memoryReadings = append(memoryReadings, memoryMB)

		t.Logf("Batch %d: Memory usage %.2fMB (total entries processed: %d)", 
			batch, memoryMB, batch*batchSize)
	}

	// Check that memory growth is not excessive
	startMemory := memoryReadings[0]
	endMemory := memoryReadings[len(memoryReadings)-1]
	memoryGrowthRatio := endMemory / startMemory

	if memoryGrowthRatio > 3.0 {
		t.Errorf("Memory growth too high: %.2fx over %d batches", memoryGrowthRatio, numBatches)
		t.Logf("Memory readings: %v", memoryReadings)
	}

	t.Logf("Memory growth pattern: start=%.2fMB, end=%.2fMB, ratio=%.2fx", 
		startMemory, endMemory, memoryGrowthRatio)
}

// TestMemoryMonitorIntegration tests the memory monitor service
func TestMemoryMonitorIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping memory monitor integration tests in short mode")
	}

	monitor := services.NewMemoryMonitor()

	// Enable monitoring
	monitor.Enable()

	// Create some processing load
	entries := createTestEntriesWithText(5000, 150)
	service := services.NewTypographyService(true, true)
	_ = service.ProcessEntries(entries)

	// Get stats
	stats := monitor.GetCurrentStats()

	// Verify stats are collected
	if stats.Allocated <= 0 {
		t.Error("Expected positive allocated memory")
	}

	t.Logf("Memory monitor stats: Allocated=%.2fMB, System=%.2fMB", 
		float64(stats.Allocated)/(1024*1024), float64(stats.System)/(1024*1024))

	// Memory should be reasonable for the test workload
	allocatedMB := float64(stats.Allocated) / (1024 * 1024)
	if allocatedMB > 200.0 {
		t.Errorf("Allocated memory too high: %.2fMB", allocatedMB)
	}
}

// TestMemoryEfficiencyWithLargeFiles tests memory efficiency with realistic file sizes
func TestMemoryEfficiencyWithLargeFiles(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping large file memory efficiency tests in short mode")
	}

	tests := []struct {
		name               string
		numRows            int
		estimatedFileSizeMB float64
		maxMemoryMultiplier float64 // Maximum memory usage as multiple of file size
	}{
		{
			name:                "medium file - 25K rows",
			numRows:             25000,
			estimatedFileSizeMB: 5.0,
			maxMemoryMultiplier: 4.0,
		},
		{
			name:                "large file - 75K rows",
			numRows:             75000,
			estimatedFileSizeMB: 15.0,
			maxMemoryMultiplier: 3.5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testFile := createLargeTempCSV(t, tt.numRows, 10)
			defer os.Remove(testFile)

			// Get actual file size
			fileInfo, err := os.Stat(testFile)
			if err != nil {
				t.Fatalf("Failed to stat test file: %v", err)
			}
			actualFileSizeMB := float64(fileInfo.Size()) / (1024 * 1024)

			// Process with memory monitoring
			monitor := services.NewMemoryMonitor()
			monitor.Enable()

			config := app.ProcessorConfig{
				FrenchMode:     true,
				SmartQuotes:    true,
				SkipDuplicates: true,
				Verbose:        false,
			}
			processor := app.NewProcessor(config)
			
			_, err = processor.ProcessFiles([]string{testFile})
			if err != nil {
				t.Fatalf("Processing failed: %v", err)
			}

			stats := monitor.GetCurrentStats()

			// Check memory efficiency
			allocatedMB := float64(stats.Allocated) / (1024 * 1024)
			memoryMultiplier := allocatedMB / actualFileSizeMB
			if memoryMultiplier > tt.maxMemoryMultiplier {
				t.Errorf("Memory efficiency poor: %.2fx file size (max: %.2fx)", 
					memoryMultiplier, tt.maxMemoryMultiplier)
			}

			t.Logf("Memory efficiency: File=%.2fMB, Allocated=%.2fMB, Multiplier=%.2fx", 
				actualFileSizeMB, allocatedMB, memoryMultiplier)
		})
	}
}