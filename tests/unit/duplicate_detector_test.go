package unit

import (
	"ankiprep/internal/models"
	"ankiprep/internal/services"
	"testing"
)

func TestDuplicateDetector_NewDuplicateDetector(t *testing.T) {
	detector := services.NewDuplicateDetector()

	if detector == nil {
		t.Fatal("NewDuplicateDetector should return a non-nil DuplicateDetector")
	}

	// Check that seenHashes is initialized (via GetSeenHashes)
	seenHashes := detector.GetSeenHashes()
	if seenHashes == nil {
		t.Fatal("seenHashes should be initialized")
	}

	if len(seenHashes) != 0 {
		t.Errorf("seenHashes should be empty initially, got %d entries", len(seenHashes))
	}
}

func TestDuplicateDetector_DetectDuplicates(t *testing.T) {
	detector := services.NewDuplicateDetector()

	tests := []struct {
		name                   string
		entries                []*models.DataEntry
		expectedUniqueCount    int
		expectedDuplicateCount int
	}{
		{
			name: "no duplicates",
			entries: []*models.DataEntry{
				createTestEntry(map[string]string{"front": "hello", "back": "world"}, "file1.csv", 1),
				createTestEntry(map[string]string{"front": "foo", "back": "bar"}, "file1.csv", 2),
			},
			expectedUniqueCount:    2,
			expectedDuplicateCount: 0,
		},
		{
			name: "exact duplicates",
			entries: []*models.DataEntry{
				createTestEntry(map[string]string{"front": "hello", "back": "world"}, "file1.csv", 1),
				createTestEntry(map[string]string{"front": "hello", "back": "world"}, "file1.csv", 2),
				createTestEntry(map[string]string{"front": "foo", "back": "bar"}, "file1.csv", 3),
			},
			expectedUniqueCount:    2,
			expectedDuplicateCount: 1,
		},
		{
			name: "multiple duplicates",
			entries: []*models.DataEntry{
				createTestEntry(map[string]string{"front": "hello", "back": "world"}, "file1.csv", 1),
				createTestEntry(map[string]string{"front": "hello", "back": "world"}, "file1.csv", 2),
				createTestEntry(map[string]string{"front": "hello", "back": "world"}, "file1.csv", 3),
				createTestEntry(map[string]string{"front": "foo", "back": "bar"}, "file1.csv", 4),
			},
			expectedUniqueCount:    2,
			expectedDuplicateCount: 2,
		},
		{
			name:                   "empty slice",
			entries:                []*models.DataEntry{},
			expectedUniqueCount:    0,
			expectedDuplicateCount: 0,
		},
		{
			name: "nil entries in slice",
			entries: []*models.DataEntry{
				createTestEntry(map[string]string{"front": "hello", "back": "world"}, "file1.csv", 1),
				nil,
				createTestEntry(map[string]string{"front": "foo", "back": "bar"}, "file1.csv", 2),
			},
			expectedUniqueCount:    2,
			expectedDuplicateCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Filter out nil entries for testing
			var validEntries []*models.DataEntry
			for _, entry := range tt.entries {
				if entry != nil {
					validEntries = append(validEntries, entry)
				}
			}

			uniqueEntries, duplicateCount := detector.DetectDuplicates(validEntries)

			if len(uniqueEntries) != tt.expectedUniqueCount {
				t.Errorf("expected %d unique entries, got %d", tt.expectedUniqueCount, len(uniqueEntries))
			}

			if duplicateCount != tt.expectedDuplicateCount {
				t.Errorf("expected %d duplicates, got %d", tt.expectedDuplicateCount, duplicateCount)
			}

			// Verify all returned entries are non-nil and valid
			for i, entry := range uniqueEntries {
				if entry == nil {
					t.Errorf("unique entry %d is nil", i)
				}
			}
		})
	}
}

func TestDuplicateDetector_DetectDuplicatesAcrossFiles(t *testing.T) {
	detector := services.NewDuplicateDetector()

	tests := []struct {
		name                          string
		entries                       []*models.DataEntry
		expectedUniqueCount           int
		expectedDuplicateCount        int
		expectedDuplicateSourcesCount int
	}{
		{
			name: "no duplicates across files",
			entries: []*models.DataEntry{
				createTestEntry(map[string]string{"front": "hello", "back": "world"}, "file1.csv", 1),
				createTestEntry(map[string]string{"front": "foo", "back": "bar"}, "file2.csv", 1),
			},
			expectedUniqueCount:           2,
			expectedDuplicateCount:        0,
			expectedDuplicateSourcesCount: 0,
		},
		{
			name: "duplicates across files",
			entries: []*models.DataEntry{
				createTestEntry(map[string]string{"front": "hello", "back": "world"}, "file1.csv", 1),
				createTestEntry(map[string]string{"front": "hello", "back": "world"}, "file2.csv", 1),
				createTestEntry(map[string]string{"front": "foo", "back": "bar"}, "file1.csv", 2),
			},
			expectedUniqueCount:           2,
			expectedDuplicateCount:        1,
			expectedDuplicateSourcesCount: 1,
		},
		{
			name: "multiple duplicates across multiple files",
			entries: []*models.DataEntry{
				createTestEntry(map[string]string{"front": "hello", "back": "world"}, "file1.csv", 1),
				createTestEntry(map[string]string{"front": "hello", "back": "world"}, "file2.csv", 1),
				createTestEntry(map[string]string{"front": "hello", "back": "world"}, "file3.csv", 1),
				createTestEntry(map[string]string{"front": "foo", "back": "bar"}, "file1.csv", 2),
			},
			expectedUniqueCount:           2,
			expectedDuplicateCount:        2,
			expectedDuplicateSourcesCount: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uniqueEntries, duplicateCount, duplicateSources := detector.DetectDuplicatesAcrossFiles(tt.entries)

			if len(uniqueEntries) != tt.expectedUniqueCount {
				t.Errorf("expected %d unique entries, got %d", tt.expectedUniqueCount, len(uniqueEntries))
			}

			if duplicateCount != tt.expectedDuplicateCount {
				t.Errorf("expected %d duplicates, got %d", tt.expectedDuplicateCount, duplicateCount)
			}

			if len(duplicateSources) != tt.expectedDuplicateSourcesCount {
				t.Errorf("expected %d duplicate source entries, got %d", tt.expectedDuplicateSourcesCount, len(duplicateSources))
			}

			// Verify duplicate sources contain expected file lists
			for hash, sources := range duplicateSources {
				if len(sources) < 2 {
					t.Errorf("duplicate source entry for hash %s should contain at least 2 files, got %d", hash, len(sources))
				}
			}
		})
	}
}

func TestDuplicateDetector_IsExactDuplicate(t *testing.T) {
	detector := services.NewDuplicateDetector()

	entry1 := createTestEntry(map[string]string{"front": "hello", "back": "world"}, "file1.csv", 1)
	entry2 := createTestEntry(map[string]string{"front": "hello", "back": "world"}, "file2.csv", 1)
	entry3 := createTestEntry(map[string]string{"front": "hello", "back": "different"}, "file1.csv", 1)

	tests := []struct {
		name     string
		entry1   *models.DataEntry
		entry2   *models.DataEntry
		expected bool
	}{
		{
			name:     "exact duplicates",
			entry1:   entry1,
			entry2:   entry2,
			expected: true,
		},
		{
			name:     "different content",
			entry1:   entry1,
			entry2:   entry3,
			expected: false,
		},
		{
			name:     "same entry with itself",
			entry1:   entry1,
			entry2:   entry1,
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := detector.IsExactDuplicate(tt.entry1, tt.entry2)
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestDuplicateDetector_GetDuplicateStats(t *testing.T) {
	detector := services.NewDuplicateDetector()

	tests := []struct {
		name               string
		originalCount      int
		uniqueEntries      []*models.DataEntry
		expectedUnique     int
		expectedDuplicates int
		expectedRate       float64
	}{
		{
			name:          "no duplicates",
			originalCount: 3,
			uniqueEntries: []*models.DataEntry{
				createTestEntry(map[string]string{"front": "1"}, "file1.csv", 1),
				createTestEntry(map[string]string{"front": "2"}, "file1.csv", 2),
				createTestEntry(map[string]string{"front": "3"}, "file1.csv", 3),
			},
			expectedUnique:     3,
			expectedDuplicates: 0,
			expectedRate:       0.0,
		},
		{
			name:          "50% duplicates",
			originalCount: 4,
			uniqueEntries: []*models.DataEntry{
				createTestEntry(map[string]string{"front": "1"}, "file1.csv", 1),
				createTestEntry(map[string]string{"front": "2"}, "file1.csv", 2),
			},
			expectedUnique:     2,
			expectedDuplicates: 2,
			expectedRate:       50.0,
		},
		{
			name:               "empty entries",
			originalCount:      0,
			uniqueEntries:      []*models.DataEntry{},
			expectedUnique:     0,
			expectedDuplicates: 0,
			expectedRate:       0.0,
		},
		{
			name:          "100% duplicates (all removed except one)",
			originalCount: 5,
			uniqueEntries: []*models.DataEntry{
				createTestEntry(map[string]string{"front": "1"}, "file1.csv", 1),
			},
			expectedUnique:     1,
			expectedDuplicates: 4,
			expectedRate:       80.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uniqueCount, duplicateCount, duplicateRate := detector.GetDuplicateStats(tt.originalCount, tt.uniqueEntries)

			if uniqueCount != tt.expectedUnique {
				t.Errorf("expected %d unique, got %d", tt.expectedUnique, uniqueCount)
			}

			if duplicateCount != tt.expectedDuplicates {
				t.Errorf("expected %d duplicates, got %d", tt.expectedDuplicates, duplicateCount)
			}

			// Allow small floating point precision differences
			if (duplicateRate-tt.expectedRate) > 0.01 || (duplicateRate-tt.expectedRate) < -0.01 {
				t.Errorf("expected %.2f%% duplicate rate, got %.2f%%", tt.expectedRate, duplicateRate)
			}
		})
	}
}

func TestDuplicateDetector_FindDuplicatesWithinFile(t *testing.T) {
	detector := services.NewDuplicateDetector()

	tests := []struct {
		name                   string
		entries                []*models.DataEntry
		filePath               string
		expectedUniqueCount    int
		expectedDuplicateCount int
	}{
		{
			name: "no duplicates within file",
			entries: []*models.DataEntry{
				createTestEntry(map[string]string{"front": "hello"}, "file1.csv", 1),
				createTestEntry(map[string]string{"front": "world"}, "file1.csv", 2),
				createTestEntry(map[string]string{"front": "hello"}, "file2.csv", 1), // Different file
			},
			filePath:               "file1.csv",
			expectedUniqueCount:    2,
			expectedDuplicateCount: 0,
		},
		{
			name: "duplicates within specified file",
			entries: []*models.DataEntry{
				createTestEntry(map[string]string{"front": "hello"}, "file1.csv", 1),
				createTestEntry(map[string]string{"front": "hello"}, "file1.csv", 2), // Duplicate
				createTestEntry(map[string]string{"front": "world"}, "file1.csv", 3),
				createTestEntry(map[string]string{"front": "hello"}, "file2.csv", 1), // Different file, ignored
			},
			filePath:               "file1.csv",
			expectedUniqueCount:    2,
			expectedDuplicateCount: 1,
		},
		{
			name: "only entries from other files",
			entries: []*models.DataEntry{
				createTestEntry(map[string]string{"front": "hello"}, "file2.csv", 1),
				createTestEntry(map[string]string{"front": "world"}, "file3.csv", 1),
			},
			filePath:               "file1.csv",
			expectedUniqueCount:    0,
			expectedDuplicateCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uniqueEntries, duplicateCount := detector.FindDuplicatesWithinFile(tt.entries, tt.filePath)

			if len(uniqueEntries) != tt.expectedUniqueCount {
				t.Errorf("expected %d unique entries, got %d", tt.expectedUniqueCount, len(uniqueEntries))
			}

			if duplicateCount != tt.expectedDuplicateCount {
				t.Errorf("expected %d duplicates, got %d", tt.expectedDuplicateCount, duplicateCount)
			}

			// Verify all returned entries are from the specified file
			for _, entry := range uniqueEntries {
				if entry.Source != tt.filePath {
					t.Errorf("expected entry from %s, got entry from %s", tt.filePath, entry.Source)
				}
			}
		})
	}
}

func TestDuplicateDetector_Reset(t *testing.T) {
	detector := services.NewDuplicateDetector()

	// Add some entries to populate internal state
	entries := []*models.DataEntry{
		createTestEntry(map[string]string{"front": "hello"}, "file1.csv", 1),
		createTestEntry(map[string]string{"front": "world"}, "file1.csv", 2),
	}
	detector.DetectDuplicates(entries)

	// Verify internal state is populated
	seenHashes := detector.GetSeenHashes()
	if len(seenHashes) != 2 {
		t.Errorf("expected 2 entries in seenHashes before reset, got %d", len(seenHashes))
	}

	// Reset and verify internal state is cleared
	detector.Reset()
	seenHashes = detector.GetSeenHashes()
	if len(seenHashes) != 0 {
		t.Errorf("expected 0 entries in seenHashes after reset, got %d", len(seenHashes))
	}
}

func TestDuplicateDetector_GetSeenHashes(t *testing.T) {
	detector := services.NewDuplicateDetector()

	// Initially should be empty
	seenHashes := detector.GetSeenHashes()
	if len(seenHashes) != 0 {
		t.Errorf("expected 0 entries initially, got %d", len(seenHashes))
	}

	// Add some entries
	entries := []*models.DataEntry{
		createTestEntry(map[string]string{"front": "hello"}, "file1.csv", 1),
		createTestEntry(map[string]string{"front": "world"}, "file2.csv", 1),
	}
	detector.DetectDuplicates(entries)

	// Should now contain entries
	seenHashes = detector.GetSeenHashes()
	if len(seenHashes) != 2 {
		t.Errorf("expected 2 entries after detection, got %d", len(seenHashes))
	}

	// Verify the mapping contains correct source files
	foundFile1 := false
	foundFile2 := false
	for _, source := range seenHashes {
		if source == "file1.csv" {
			foundFile1 = true
		}
		if source == "file2.csv" {
			foundFile2 = true
		}
	}

	if !foundFile1 {
		t.Error("expected file1.csv in seenHashes")
	}
	if !foundFile2 {
		t.Error("expected file2.csv in seenHashes")
	}
}

// Helper function to create test entries
func createTestEntry(values map[string]string, source string, lineNumber int) *models.DataEntry {
	return models.NewDataEntry(values, source, lineNumber)
}
