package services

import (
	"fmt"
	"runtime"
	"runtime/debug"
	"time"
)

// MemoryMonitor provides memory usage tracking and optimization for large datasets
type MemoryMonitor struct {
	initialMemory uint64
	peakMemory    uint64
	enabled       bool
	gcThreshold   uint64 // Memory threshold to trigger GC (in bytes)
	lastGC        time.Time
	gcInterval    time.Duration
}

// MemoryStats represents memory usage statistics
type MemoryStats struct {
	Allocated  uint64        // Currently allocated memory
	TotalAlloc uint64        // Total allocated memory over time
	System     uint64        // System memory obtained from OS
	NumGC      uint32        // Number of GC cycles
	PauseTotal time.Duration // Total GC pause time
	LastGC     time.Time     // Time of last GC
	HeapInUse  uint64        // Heap memory in use
	StackInUse uint64        // Stack memory in use
}

// NewMemoryMonitor creates a new memory monitor
func NewMemoryMonitor() *MemoryMonitor {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	return &MemoryMonitor{
		initialMemory: m.Alloc,
		peakMemory:    m.Alloc,
		enabled:       false,
		gcThreshold:   100 * 1024 * 1024, // 100MB default threshold
		gcInterval:    30 * time.Second,  // Minimum interval between forced GCs
	}
}

// Enable turns on memory monitoring
func (mm *MemoryMonitor) Enable() {
	mm.enabled = true
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	mm.initialMemory = m.Alloc
	mm.peakMemory = m.Alloc
}

// Disable turns off memory monitoring
func (mm *MemoryMonitor) Disable() {
	mm.enabled = false
}

// SetGCThreshold sets the memory threshold for triggering garbage collection
func (mm *MemoryMonitor) SetGCThreshold(bytes uint64) {
	mm.gcThreshold = bytes
}

// SetGCInterval sets the minimum interval between forced garbage collections
func (mm *MemoryMonitor) SetGCInterval(interval time.Duration) {
	mm.gcInterval = interval
}

// GetCurrentStats returns current memory statistics
func (mm *MemoryMonitor) GetCurrentStats() MemoryStats {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	// Update peak memory
	if m.Alloc > mm.peakMemory {
		mm.peakMemory = m.Alloc
	}

	return MemoryStats{
		Allocated:  m.Alloc,
		TotalAlloc: m.TotalAlloc,
		System:     m.Sys,
		NumGC:      m.NumGC,
		PauseTotal: time.Duration(m.PauseTotalNs),
		LastGC:     time.Unix(0, int64(m.LastGC)),
		HeapInUse:  m.HeapInuse,
		StackInUse: m.StackInuse,
	}
}

// GetMemoryUsage returns formatted memory usage information
func (mm *MemoryMonitor) GetMemoryUsage() string {
	if !mm.enabled {
		return "Memory monitoring disabled"
	}

	stats := mm.GetCurrentStats()

	return fmt.Sprintf("Memory: %s allocated, %s heap, %s system, %d GC cycles",
		formatBytes(stats.Allocated),
		formatBytes(stats.HeapInUse),
		formatBytes(stats.System),
		stats.NumGC)
}

// GetMemoryDelta returns the change in memory usage since monitoring started
func (mm *MemoryMonitor) GetMemoryDelta() string {
	if !mm.enabled {
		return "Memory monitoring disabled"
	}

	stats := mm.GetCurrentStats()
	delta := int64(stats.Allocated) - int64(mm.initialMemory)

	if delta >= 0 {
		return fmt.Sprintf("+%s from start", formatBytes(uint64(delta)))
	} else {
		return fmt.Sprintf("-%s from start", formatBytes(uint64(-delta)))
	}
}

// GetPeakMemoryUsage returns the peak memory usage observed
func (mm *MemoryMonitor) GetPeakMemoryUsage() string {
	if !mm.enabled {
		return "Memory monitoring disabled"
	}

	return fmt.Sprintf("Peak memory: %s", formatBytes(mm.peakMemory))
}

// CheckMemoryPressure checks if memory usage is high and suggests optimizations
func (mm *MemoryMonitor) CheckMemoryPressure() (bool, string) {
	if !mm.enabled {
		return false, ""
	}

	stats := mm.GetCurrentStats()

	// Check if allocated memory exceeds threshold
	if stats.Allocated > mm.gcThreshold {
		return true, fmt.Sprintf("High memory usage: %s allocated (threshold: %s)",
			formatBytes(stats.Allocated), formatBytes(mm.gcThreshold))
	}

	// Check if heap usage is very high relative to system memory
	if stats.HeapInUse > stats.System/2 {
		return true, fmt.Sprintf("High heap usage: %s of %s system memory",
			formatBytes(stats.HeapInUse), formatBytes(stats.System))
	}

	return false, ""
}

// TriggerGCIfNeeded performs garbage collection if memory pressure is high
func (mm *MemoryMonitor) TriggerGCIfNeeded() bool {
	if !mm.enabled {
		return false
	}

	// Check if enough time has passed since last forced GC
	if time.Since(mm.lastGC) < mm.gcInterval {
		return false
	}

	pressure, _ := mm.CheckMemoryPressure()
	if pressure {
		beforeStats := mm.GetCurrentStats()
		runtime.GC()
		runtime.GC() // Call twice to ensure cleanup
		mm.lastGC = time.Now()

		afterStats := mm.GetCurrentStats()
		freed := beforeStats.Allocated - afterStats.Allocated

		if freed > 1024*1024 { // Only report if freed > 1MB
			return true
		}
	}

	return false
}

// OptimizeForLargeDataset configures memory settings for large dataset processing
func (mm *MemoryMonitor) OptimizeForLargeDataset(expectedDataSize uint64) {
	mm.Enable()

	// Set GC threshold based on expected data size
	// Use 2x the expected size as threshold, minimum 50MB
	threshold := expectedDataSize * 2
	if threshold < 50*1024*1024 {
		threshold = 50 * 1024 * 1024
	}
	mm.SetGCThreshold(threshold)

	// More aggressive GC for large datasets
	mm.SetGCInterval(10 * time.Second)

	// Configure runtime for better memory management
	debug.SetGCPercent(50) // Trigger GC more frequently

	// Set memory limit if we can estimate it
	if expectedDataSize > 500*1024*1024 { // > 500MB
		// Set a reasonable memory limit to prevent OOM
		memLimit := expectedDataSize * 3 // Allow 3x the data size
		if memLimit > 2*1024*1024*1024 { // Cap at 2GB
			memLimit = 2 * 1024 * 1024 * 1024
		}
		debug.SetMemoryLimit(int64(memLimit))
	}
}

// ResetOptimizations resets memory optimizations to defaults
func (mm *MemoryMonitor) ResetOptimizations() {
	debug.SetGCPercent(100)            // Default GC percentage
	debug.SetMemoryLimit(-1)           // Remove memory limit
	mm.gcThreshold = 100 * 1024 * 1024 // Reset to 100MB
	mm.gcInterval = 30 * time.Second   // Reset to 30s
}

// GetDetailedStats returns comprehensive memory statistics
func (mm *MemoryMonitor) GetDetailedStats() map[string]interface{} {
	stats := mm.GetCurrentStats()

	return map[string]interface{}{
		"enabled":         mm.enabled,
		"allocated_mb":    float64(stats.Allocated) / (1024 * 1024),
		"heap_inuse_mb":   float64(stats.HeapInUse) / (1024 * 1024),
		"system_mb":       float64(stats.System) / (1024 * 1024),
		"peak_mb":         float64(mm.peakMemory) / (1024 * 1024),
		"initial_mb":      float64(mm.initialMemory) / (1024 * 1024),
		"num_gc":          stats.NumGC,
		"pause_total_ms":  float64(stats.PauseTotal) / float64(time.Millisecond),
		"last_gc":         stats.LastGC.Format(time.RFC3339),
		"gc_threshold_mb": float64(mm.gcThreshold) / (1024 * 1024),
	}
}

// EstimateMemoryRequirement estimates memory needed for processing
func (mm *MemoryMonitor) EstimateMemoryRequirement(recordCount int, avgRecordSize int) uint64 {
	// Estimate total data size
	dataSize := uint64(recordCount * avgRecordSize)

	// Account for:
	// - Original data in memory
	// - Processed data structures
	// - Intermediate processing buffers
	// - Go runtime overhead

	multiplier := 3.5 // Conservative multiplier for overhead
	return uint64(float64(dataSize) * multiplier)
}

// ShouldStreamProcess determines if streaming processing should be used
func (mm *MemoryMonitor) ShouldStreamProcess(estimatedMemory uint64) bool {
	// Get available memory (simplified approach)
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	// If estimated memory is more than 50% of system memory, use streaming
	return estimatedMemory > m.Sys/2
}

// LogMemoryUsage logs current memory usage to the progress reporter
func (mm *MemoryMonitor) LogMemoryUsage(reporter *ProgressReporter) {
	if !mm.enabled || reporter == nil {
		return
	}

	usage := mm.GetMemoryUsage()
	delta := mm.GetMemoryDelta()

	if reporter.IsVerbose() {
		fmt.Printf("Memory status: %s (%s)\n", usage, delta)

		// Check for memory pressure and warn if needed
		if pressure, message := mm.CheckMemoryPressure(); pressure {
			reporter.ReportWarning(message)
		}
	}
}

// formatBytes formats byte count in human-readable format
func formatBytes(bytes uint64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}

	div, exp := uint64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}

	units := []string{"KB", "MB", "GB", "TB"}
	return fmt.Sprintf("%.1f %s", float64(bytes)/float64(div), units[exp])
}
