package services

import (
	"ankiprep/internal/models"
)

// DuplicateDetector handles detection and removal of duplicate entries
type DuplicateDetector struct {
	seenHashes map[string]*models.DataEntry // Map of hash -> first occurrence
}

// NewDuplicateDetector creates a new DuplicateDetector instance
func NewDuplicateDetector() *DuplicateDetector {
	return &DuplicateDetector{
		seenHashes: make(map[string]*models.DataEntry),
	}
}

// DetectDuplicates finds and removes duplicate entries from a slice of DataEntry
func (d *DuplicateDetector) DetectDuplicates(entries []*models.DataEntry) ([]*models.DataEntry, int) {
	var uniqueEntries []*models.DataEntry
	duplicateCount := 0

	// Reset the seen hashes for this detection run
	d.seenHashes = make(map[string]*models.DataEntry)

	for _, entry := range entries {
		hash := entry.GetHash()

		if existingEntry, exists := d.seenHashes[hash]; exists {
			// This is a duplicate - check if it's an exact match
			if entry.IsExactDuplicate(existingEntry) {
				duplicateCount++
				// Skip this entry (don't add to uniqueEntries)
				continue
			}
		}

		// This is a unique entry
		d.seenHashes[hash] = entry
		uniqueEntries = append(uniqueEntries, entry)
	}

	return uniqueEntries, duplicateCount
}

// DetectDuplicatesAcrossFiles finds duplicates across multiple input files
func (d *DuplicateDetector) DetectDuplicatesAcrossFiles(allEntries []*models.DataEntry) ([]*models.DataEntry, int, map[string][]string) {
	var uniqueEntries []*models.DataEntry
	duplicateCount := 0
	duplicateSources := make(map[string][]string) // hash -> list of source files

	// Reset the seen hashes
	d.seenHashes = make(map[string]*models.DataEntry)

	for _, entry := range allEntries {
		hash := entry.GetHash()

		if existingEntry, exists := d.seenHashes[hash]; exists {
			// This is a duplicate
			if entry.IsExactDuplicate(existingEntry) {
				duplicateCount++

				// Track which files contain this duplicate
				if duplicateSources[hash] == nil {
					duplicateSources[hash] = []string{existingEntry.Source}
				}
				duplicateSources[hash] = append(duplicateSources[hash], entry.Source)

				// Skip this entry
				continue
			}
		}

		// This is a unique entry
		d.seenHashes[hash] = entry
		uniqueEntries = append(uniqueEntries, entry)
	}

	return uniqueEntries, duplicateCount, duplicateSources
}

// IsExactDuplicate checks if two entries are exact duplicates
func (d *DuplicateDetector) IsExactDuplicate(entry1, entry2 *models.DataEntry) bool {
	return entry1.IsExactDuplicate(entry2)
}

// GetDuplicateStats returns statistics about duplicates found
func (d *DuplicateDetector) GetDuplicateStats(originalCount int, uniqueEntries []*models.DataEntry) (int, int, float64) {
	uniqueCount := len(uniqueEntries)
	duplicateCount := originalCount - uniqueCount
	duplicateRate := 0.0

	if originalCount > 0 {
		duplicateRate = float64(duplicateCount) / float64(originalCount) * 100.0
	}

	return uniqueCount, duplicateCount, duplicateRate
}

// FindDuplicatesWithinFile finds duplicates within a single file's entries
func (d *DuplicateDetector) FindDuplicatesWithinFile(entries []*models.DataEntry, filePath string) ([]*models.DataEntry, int) {
	var uniqueEntries []*models.DataEntry
	duplicateCount := 0
	localHashes := make(map[string]*models.DataEntry)

	for _, entry := range entries {
		// Only process entries from the specified file
		if entry.Source != filePath {
			continue
		}

		hash := entry.GetHash()

		if existingEntry, exists := localHashes[hash]; exists {
			if entry.IsExactDuplicate(existingEntry) {
				duplicateCount++
				continue
			}
		}

		localHashes[hash] = entry
		uniqueEntries = append(uniqueEntries, entry)
	}

	return uniqueEntries, duplicateCount
}

// GetSeenHashes returns a copy of the currently seen hashes for debugging
func (d *DuplicateDetector) GetSeenHashes() map[string]string {
	result := make(map[string]string)
	for hash, entry := range d.seenHashes {
		result[hash] = entry.Source
	}
	return result
}

// Reset clears the internal state of the duplicate detector
func (d *DuplicateDetector) Reset() {
	d.seenHashes = make(map[string]*models.DataEntry)
}
