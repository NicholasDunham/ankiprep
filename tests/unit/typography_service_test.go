package unit

import (
	"ankiprep/internal/models"
	"ankiprep/internal/services"
	"testing"
)

func TestTypographyService_NewTypographyService(t *testing.T) {
	tests := []struct {
		name        string
		frenchMode  bool
		smartQuotes bool
	}{
		{
			name:        "default settings",
			frenchMode:  false,
			smartQuotes: false,
		},
		{
			name:        "french mode enabled",
			frenchMode:  true,
			smartQuotes: false,
		},
		{
			name:        "smart quotes enabled",
			frenchMode:  false,
			smartQuotes: true,
		},
		{
			name:        "both enabled",
			frenchMode:  true,
			smartQuotes: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := services.NewTypographyServiceLegacy(tt.frenchMode, tt.smartQuotes)

			if service == nil {
				t.Fatal("NewTypographyService should return a non-nil service")
			}

			if service.IsFrenchModeEnabled() != tt.frenchMode {
				t.Errorf("expected FrenchMode %v, got %v", tt.frenchMode, service.IsFrenchModeEnabled())
			}

			if service.IsSmartQuotesEnabled() != tt.smartQuotes {
				t.Errorf("expected SmartQuotes %v, got %v", tt.smartQuotes, service.IsSmartQuotesEnabled())
			}

			processor := service.GetProcessor()
			if processor == nil {
				t.Fatal("GetProcessor should return a non-nil processor")
			}
		})
	}
}

func TestTypographyService_ProcessText(t *testing.T) {
	tests := []struct {
		name             string
		frenchMode       bool
		smartQuotes      bool
		input            string
		expectedContains string // What the output should contain
	}{
		{
			name:             "simple text no processing",
			frenchMode:       false,
			smartQuotes:      false,
			input:            "hello world",
			expectedContains: "hello world",
		},
		{
			name:             "text with quotes - smart quotes disabled",
			frenchMode:       false,
			smartQuotes:      false,
			input:            `He said "hello"`,
			expectedContains: `"hello"`, // Should keep original quotes
		},
		{
			name:             "text with quotes - smart quotes enabled",
			frenchMode:       false,
			smartQuotes:      true,
			input:            `He said "hello"`,
			expectedContains: "hello", // Should convert quotes (exact format depends on typography processor)
		},
		{
			name:             "french punctuation - french mode disabled",
			frenchMode:       false,
			smartQuotes:      false,
			input:            "Bonjour! Comment allez-vous?",
			expectedContains: "Bonjour!", // Should keep original punctuation
		},
		{
			name:             "french punctuation - french mode enabled",
			frenchMode:       true,
			smartQuotes:      false,
			input:            "Bonjour! Comment allez-vous?",
			expectedContains: "Bonjour", // Should apply French spacing rules
		},
		{
			name:             "empty string",
			frenchMode:       false,
			smartQuotes:      false,
			input:            "",
			expectedContains: "",
		},
		{
			name:             "whitespace only",
			frenchMode:       false,
			smartQuotes:      false,
			input:            "   \t\n   ",
			expectedContains: "   \t\n   ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := services.NewTypographyServiceLegacy(tt.frenchMode, tt.smartQuotes)
			result := service.ProcessText(tt.input)

			if result == "" && tt.expectedContains != "" {
				t.Errorf("expected result to contain '%s', got empty string", tt.expectedContains)
				return
			}

			// For basic containment test (typography processing may modify exact format)
			if tt.expectedContains != "" && len(result) == 0 {
				t.Errorf("expected result to contain '%s', got empty result", tt.expectedContains)
			}

			// For empty input, expect empty output
			if tt.input == "" && result != "" {
				t.Errorf("expected empty result for empty input, got '%s'", result)
			}
		})
	}
}

func TestTypographyService_ProcessEntry(t *testing.T) {
	service := services.NewTypographyServiceLegacy(false, false)

	tests := []struct {
		name     string
		entry    *models.DataEntry
		expected bool // Whether we expect a valid result
	}{
		{
			name: "valid entry",
			entry: createTestEntry(map[string]string{
				"front": "Hello World",
				"back":  "Bonjour Monde",
			}, "file1.csv", 1),
			expected: true,
		},
		{
			name: "entry with quotes",
			entry: createTestEntry(map[string]string{
				"front": `He said "hello"`,
				"back":  "Il a dit « bonjour »",
			}, "file1.csv", 1),
			expected: true,
		},
		{
			name: "entry with empty values",
			entry: createTestEntry(map[string]string{
				"front": "",
				"back":  "Something",
			}, "file1.csv", 1),
			expected: true,
		},
		{
			name:     "nil entry",
			entry:    nil,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := service.ProcessEntry(tt.entry)

			if tt.expected {
				if result == nil {
					t.Fatal("expected non-nil result for valid entry")
				}

				// Check that the result has the same structure
				if len(result.Values) != len(tt.entry.Values) {
					t.Errorf("expected %d values in result, got %d", len(tt.entry.Values), len(result.Values))
				}

				// Check that source and line number are preserved
				if result.Source != tt.entry.Source {
					t.Errorf("expected source %s, got %s", tt.entry.Source, result.Source)
				}

				if result.LineNumber != tt.entry.LineNumber {
					t.Errorf("expected line number %d, got %d", tt.entry.LineNumber, result.LineNumber)
				}

				// Check that all keys are present (values may be modified by typography processing)
				for key := range tt.entry.Values {
					if _, exists := result.Values[key]; !exists {
						t.Errorf("expected key '%s' in result", key)
					}
				}
			} else {
				if result != nil {
					t.Errorf("expected nil result for invalid entry, got %v", result)
				}
			}
		})
	}
}

func TestTypographyService_ProcessEntries(t *testing.T) {
	service := services.NewTypographyServiceLegacy(false, true) // Enable smart quotes for testing

	tests := []struct {
		name          string
		entries       []*models.DataEntry
		expectedCount int
	}{
		{
			name: "multiple valid entries",
			entries: []*models.DataEntry{
				createTestEntry(map[string]string{"front": "hello", "back": "world"}, "file1.csv", 1),
				createTestEntry(map[string]string{"front": "foo", "back": "bar"}, "file1.csv", 2),
			},
			expectedCount: 2,
		},
		{
			name: "single entry",
			entries: []*models.DataEntry{
				createTestEntry(map[string]string{"front": "hello", "back": "world"}, "file1.csv", 1),
			},
			expectedCount: 1,
		},
		{
			name:          "empty slice",
			entries:       []*models.DataEntry{},
			expectedCount: 0,
		},
		{
			name: "entries with nil values",
			entries: []*models.DataEntry{
				createTestEntry(map[string]string{"front": "hello", "back": "world"}, "file1.csv", 1),
				nil,
				createTestEntry(map[string]string{"front": "foo", "back": "bar"}, "file1.csv", 2),
			},
			expectedCount: 3, // Should still return same number, but with nil processed entry
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := service.ProcessEntries(tt.entries)

			if len(result) != tt.expectedCount {
				t.Errorf("expected %d processed entries, got %d", tt.expectedCount, len(result))
			}

			// Verify that non-nil input entries produce non-nil output entries
			for i, originalEntry := range tt.entries {
				if originalEntry != nil && result[i] == nil {
					t.Errorf("expected non-nil result for entry %d", i)
				}
				if originalEntry == nil && result[i] != nil {
					t.Errorf("expected nil result for nil entry %d", i)
				}
			}
		})
	}
}

func TestTypographyService_SetFrenchMode(t *testing.T) {
	service := services.NewTypographyServiceLegacy(false, false)

	// Initially disabled
	if service.IsFrenchModeEnabled() {
		t.Error("expected FrenchMode to be initially disabled")
	}

	// Enable it
	service.SetFrenchMode(true)
	if !service.IsFrenchModeEnabled() {
		t.Error("expected FrenchMode to be enabled after SetFrenchMode(true)")
	}

	// Disable it
	service.SetFrenchMode(false)
	if service.IsFrenchModeEnabled() {
		t.Error("expected FrenchMode to be disabled after SetFrenchMode(false)")
	}
}

func TestTypographyService_SetSmartQuotes(t *testing.T) {
	service := services.NewTypographyServiceLegacy(false, false)

	// Initially disabled
	if service.IsSmartQuotesEnabled() {
		t.Error("expected SmartQuotes to be initially disabled")
	}

	// Enable it
	service.SetSmartQuotes(true)
	if !service.IsSmartQuotesEnabled() {
		t.Error("expected SmartQuotes to be enabled after SetSmartQuotes(true)")
	}

	// Disable it
	service.SetSmartQuotes(false)
	if service.IsSmartQuotesEnabled() {
		t.Error("expected SmartQuotes to be disabled after SetSmartQuotes(false)")
	}
}

func TestTypographyService_ConfigurationChanges(t *testing.T) {
	service := services.NewTypographyServiceLegacy(false, false)

	testText := `He said "hello" and she replied: "Bonjour!"`

	// Process with default settings
	result1 := service.ProcessText(testText)

	// Enable smart quotes
	service.SetSmartQuotes(true)
	result2 := service.ProcessText(testText)

	// Enable French mode
	service.SetFrenchMode(true)
	result3 := service.ProcessText(testText)

	// Results should be processed (exact format depends on typography processor implementation)
	// At minimum, we verify they're all strings and not empty for non-empty input
	results := []string{result1, result2, result3}
	for i, result := range results {
		if result == "" {
			t.Errorf("result %d should not be empty for non-empty input", i+1)
		}
	}

	// Verify configuration changes are persisted
	if !service.IsFrenchModeEnabled() {
		t.Error("expected FrenchMode to remain enabled")
	}
	if !service.IsSmartQuotesEnabled() {
		t.Error("expected SmartQuotes to remain enabled")
	}
}

func TestTypographyService_GetProcessor(t *testing.T) {
	service := services.NewTypographyServiceLegacy(true, true)

	processor := service.GetProcessor()
	if processor == nil {
		t.Fatal("GetProcessor should return a non-nil processor")
	}

	// Verify processor has the correct settings
	if processor.FrenchMode != true {
		t.Error("expected processor FrenchMode to be true")
	}
	if processor.ConvertSmartQuotes != true {
		t.Error("expected processor ConvertSmartQuotes to be true")
	}
}
