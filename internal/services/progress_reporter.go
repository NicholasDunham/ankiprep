package services

import (
	"fmt"
	"os"
	"time"
)

// ProgressReporter handles progress reporting and user feedback
type ProgressReporter struct {
	verbose        bool
	startTime      time.Time
	largeFileMode  bool
	fileSize       int64
	lastUpdate     time.Time
	updateInterval time.Duration
}

const (
	// File size threshold for large file mode (10MB)
	LargeFileThreshold = 10 * 1024 * 1024
	// Processing time threshold for progress updates (5 seconds)
	LongProcessingThreshold = 5 * time.Second
	// Update interval for progress reporting
	DefaultUpdateInterval = 500 * time.Millisecond
)

// NewProgressReporter creates a new ProgressReporter instance
func NewProgressReporter(verbose bool) *ProgressReporter {
	return &ProgressReporter{
		verbose:        verbose,
		updateInterval: DefaultUpdateInterval,
	}
}

// StartProcessing marks the beginning of processing and prints initial message
func (r *ProgressReporter) StartProcessing(fileCount int) {
	r.startTime = time.Now()
	if fileCount == 1 {
		fmt.Printf("Processing 1 input file...\n")
	} else {
		fmt.Printf("Processing %d input files...\n", fileCount)
	}
}

// DetectLargeFiles analyzes input files to determine if large file mode should be enabled
func (r *ProgressReporter) DetectLargeFiles(filePaths []string) {
	var totalSize int64
	for _, path := range filePaths {
		if info, err := os.Stat(path); err == nil {
			totalSize += info.Size()
		}
	}

	r.fileSize = totalSize
	r.largeFileMode = totalSize > LargeFileThreshold

	if r.largeFileMode && r.verbose {
		fmt.Printf("Large file detected (%.1f MB) - enabling detailed progress tracking\n",
			float64(totalSize)/(1024*1024))
	}
}

// ReportHeaderMerging reports on header merging progress
func (r *ProgressReporter) ReportHeaderMerging(columnCount int) {
	if r.verbose {
		fmt.Printf("Merging headers: found %d unique columns\n", columnCount)
	}
}

// ReportRecordProcessing reports on record processing progress
func (r *ProgressReporter) ReportRecordProcessing(totalRecords int) {
	if r.verbose {
		if totalRecords == 1 {
			fmt.Printf("Processing records: 1 total entry\n")
		} else {
			fmt.Printf("Processing records: %s total entries\n", r.formatNumber(totalRecords))
		}
	}

	// Enable large file mode for many records even if file size is small
	if totalRecords > 50000 {
		r.largeFileMode = true
		if r.verbose {
			fmt.Printf("Large dataset detected (%s records) - enabling progress tracking\n",
				r.formatNumber(totalRecords))
		}
	}
}

// ReportDuplicateRemoval reports on duplicate detection and removal
func (r *ProgressReporter) ReportDuplicateRemoval(duplicateCount int) {
	if duplicateCount == 0 {
		if r.verbose {
			fmt.Printf("Removing duplicates: no duplicates found\n")
		}
	} else if duplicateCount == 1 {
		fmt.Printf("Removing duplicates: found 1 duplicate\n")
	} else {
		fmt.Printf("Removing duplicates: found %s duplicates\n", r.formatNumber(duplicateCount))
	}
}

// ReportTypographyProcessing reports on typography processing
func (r *ProgressReporter) ReportTypographyProcessing(frenchMode bool, smartQuotes bool) {
	if r.verbose {
		var features []string
		if frenchMode {
			features = append(features, "French typography")
		}
		if smartQuotes {
			features = append(features, "smart quotes")
		}

		if len(features) > 0 {
			fmt.Printf("Applying typography formatting (%s)...\n", r.joinFeatures(features))
		} else {
			fmt.Printf("Applying typography formatting...\n")
		}
	}
}

// ReportFileWriting reports on output file writing
func (r *ProgressReporter) ReportFileWriting(outputPath string) {
	if r.verbose {
		fmt.Printf("Writing output to %s\n", outputPath)
	}
}

// ReportCompletion reports successful completion with statistics
func (r *ProgressReporter) ReportCompletion(uniqueRecords int, processingTime time.Duration) {
	if uniqueRecords == 1 {
		fmt.Printf("Done. Processed 1 unique entry in %.2f seconds\n", processingTime.Seconds())
	} else {
		fmt.Printf("Done. Processed %s unique entries in %.2f seconds\n",
			r.formatNumber(uniqueRecords), processingTime.Seconds())
	}

	// Show performance stats for large operations
	if r.largeFileMode && r.verbose {
		r.showPerformanceStats(uniqueRecords, processingTime)
	}
}

// ReportError reports an error message
func (r *ProgressReporter) ReportError(err error) {
	fmt.Printf("Error: %v\n", err)
}

// ReportWarning reports a warning message
func (r *ProgressReporter) ReportWarning(message string) {
	if r.verbose {
		fmt.Printf("Warning: %s\n", message)
	}
}

// ReportProgress reports progress for large file processing with time-based updates
func (r *ProgressReporter) ReportProgress(current, total int, operation string) {
	now := time.Now()

	// Always report progress for large files or long operations
	shouldReport := r.largeFileMode ||
		(now.Sub(r.startTime) > LongProcessingThreshold && total > 1000) ||
		(now.Sub(r.lastUpdate) > r.updateInterval)

	// Always report at key milestones
	if current == total || current%1000 == 0 || current == 1 {
		shouldReport = true
	}

	if shouldReport {
		percentage := float64(current) / float64(total) * 100
		elapsed := now.Sub(r.startTime)

		// Calculate ETA for large operations
		if current > 0 && r.largeFileMode {
			avgTimePerItem := elapsed / time.Duration(current)
			remaining := time.Duration(total-current) * avgTimePerItem

			fmt.Printf("%s: %.1f%% complete (%s/%s) - ETA: %s\n",
				operation,
				percentage,
				r.formatNumber(current),
				r.formatNumber(total),
				r.formatDuration(remaining))
		} else {
			fmt.Printf("%s: %.1f%% complete (%s/%s)\n",
				operation,
				percentage,
				r.formatNumber(current),
				r.formatNumber(total))
		}

		r.lastUpdate = now
	}
}

// ReportFileInfo reports information about processed files with size information
func (r *ProgressReporter) ReportFileInfo(filePath string, recordCount int, separator string) {
	if r.verbose {
		sizeInfo := ""
		if info, err := os.Stat(filePath); err == nil {
			if info.Size() > 1024*1024 {
				sizeInfo = fmt.Sprintf(" (%.1f MB)", float64(info.Size())/(1024*1024))
			} else if info.Size() > 1024 {
				sizeInfo = fmt.Sprintf(" (%.1f KB)", float64(info.Size())/1024)
			} else {
				sizeInfo = fmt.Sprintf(" (%d bytes)", info.Size())
			}
		}

		if recordCount == 1 {
			fmt.Printf("File %s: 1 record%s (%s-separated)\n", filePath, sizeInfo, separator)
		} else {
			fmt.Printf("File %s: %s records%s (%s-separated)\n",
				filePath, r.formatNumber(recordCount), sizeInfo, separator)
		}
	}
}

// StartOperation starts timing for a specific operation
func (r *ProgressReporter) StartOperation(operation string) {
	if r.verbose && r.largeFileMode {
		fmt.Printf("Starting %s...\n", operation)
		r.lastUpdate = time.Now()
	}
}

// EndOperation reports completion of a specific operation
func (r *ProgressReporter) EndOperation(operation string, itemCount int) {
	if r.verbose && r.largeFileMode {
		elapsed := time.Since(r.lastUpdate)
		if itemCount > 0 {
			rate := float64(itemCount) / elapsed.Seconds()
			fmt.Printf("Completed %s: %s items in %s (%.0f items/sec)\n",
				operation, r.formatNumber(itemCount), r.formatDuration(elapsed), rate)
		} else {
			fmt.Printf("Completed %s in %s\n", operation, r.formatDuration(elapsed))
		}
	}
}

// ShowMemoryUsage displays current memory usage for large operations
func (r *ProgressReporter) ShowMemoryUsage() {
	if r.verbose && r.largeFileMode {
		// This would require importing runtime and getting memory stats
		// For now, we'll provide a placeholder
		fmt.Printf("Memory usage monitoring available in verbose mode\n")
	}
}

// GetElapsedTime returns the elapsed time since processing started
func (r *ProgressReporter) GetElapsedTime() time.Duration {
	if r.startTime.IsZero() {
		return 0
	}
	return time.Since(r.startTime)
}

// SetVerbose enables or disables verbose output
func (r *ProgressReporter) SetVerbose(verbose bool) {
	r.verbose = verbose
}

// IsVerbose returns whether verbose mode is enabled
func (r *ProgressReporter) IsVerbose() bool {
	return r.verbose
}

// IsLargeFileMode returns whether large file mode is enabled
func (r *ProgressReporter) IsLargeFileMode() bool {
	return r.largeFileMode
}

// GetTotalFileSize returns the total file size detected
func (r *ProgressReporter) GetTotalFileSize() int64 {
	return r.fileSize
}

// SetUpdateInterval sets the minimum interval between progress updates
func (r *ProgressReporter) SetUpdateInterval(interval time.Duration) {
	r.updateInterval = interval
}

// showPerformanceStats displays detailed performance information
func (r *ProgressReporter) showPerformanceStats(recordCount int, totalTime time.Duration) {
	rate := float64(recordCount) / totalTime.Seconds()
	fmt.Printf("Performance: %.0f records/sec", rate)

	if r.fileSize > 0 {
		mbPerSec := float64(r.fileSize) / (1024 * 1024) / totalTime.Seconds()
		fmt.Printf(", %.1f MB/sec", mbPerSec)
	}

	fmt.Printf("\n")
}

// formatNumber formats a number with commas for thousands
func (r *ProgressReporter) formatNumber(n int) string {
	if n < 1000 {
		return fmt.Sprintf("%d", n)
	}

	str := fmt.Sprintf("%d", n)
	result := ""

	for i, digit := range str {
		if i > 0 && (len(str)-i)%3 == 0 {
			result += ","
		}
		result += string(digit)
	}

	return result
}

// formatDuration formats a duration in a human-readable way
func (r *ProgressReporter) formatDuration(d time.Duration) string {
	if d < time.Second {
		return fmt.Sprintf("%.0fms", float64(d)/float64(time.Millisecond))
	}
	if d < time.Minute {
		return fmt.Sprintf("%.1fs", d.Seconds())
	}
	if d < time.Hour {
		return fmt.Sprintf("%.1fm", d.Minutes())
	}
	return fmt.Sprintf("%.1fh", d.Hours())
}

// joinFeatures joins feature names with proper grammar
func (r *ProgressReporter) joinFeatures(features []string) string {
	if len(features) == 0 {
		return ""
	}
	if len(features) == 1 {
		return features[0]
	}
	if len(features) == 2 {
		return features[0] + " and " + features[1]
	}

	result := ""
	for i, feature := range features {
		if i == len(features)-1 {
			result += "and " + feature
		} else if i > 0 {
			result += ", " + feature
		} else {
			result += feature
		}
	}

	return result
}
