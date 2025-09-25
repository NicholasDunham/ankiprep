package unit_test

import (
	"testing"

	"ankiprep/internal/models"
)

func TestDataEntry_NewDataEntry(t *testing.T) {
	tests := []struct {
		name       string
		values     map[string]string
		source     string
		lineNumber int
		want       map[string]string
	}{
		{
			name:       "empty values",
			values:     map[string]string{},
			source:     "test.csv",
			lineNumber: 1,
			want:       map[string]string{},
		},
		{
			name:       "single field",
			values:     map[string]string{"field1": "value1"},
			source:     "test.csv",
			lineNumber: 2,
			want:       map[string]string{"field1": "value1"},
		},
		{
			name: "multiple fields",
			values: map[string]string{
				"french":  "bonjour",
				"english": "hello",
				"notes":   "greeting",
			},
			source:     "vocabulary.csv",
			lineNumber: 3,
			want: map[string]string{
				"french":  "bonjour",
				"english": "hello",
				"notes":   "greeting",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			entry := models.NewDataEntry(tt.values, tt.source, tt.lineNumber)

			if entry == nil {
				t.Fatal("NewDataEntry returned nil")
			}

			if entry.Source != tt.source {
				t.Errorf("NewDataEntry() source = %v, want %v", entry.Source, tt.source)
			}

			if entry.LineNumber != tt.lineNumber {
				t.Errorf("NewDataEntry() lineNumber = %v, want %v", entry.LineNumber, tt.lineNumber)
			}

			if len(entry.Values) != len(tt.want) {
				t.Errorf("NewDataEntry() values length = %v, want %v", len(entry.Values), len(tt.want))
			}

			for key, expectedValue := range tt.want {
				if actualValue, exists := entry.Values[key]; !exists {
					t.Errorf("NewDataEntry() missing key %v", key)
				} else if actualValue != expectedValue {
					t.Errorf("NewDataEntry() values[%v] = %v, want %v", key, actualValue, expectedValue)
				}
			}
		})
	}
}

func TestDataEntry_GetValue(t *testing.T) {
	values := map[string]string{
		"french":  "bonjour",
		"english": "hello",
		"notes":   "",
	}
	entry := models.NewDataEntry(values, "test.csv", 1)

	tests := []struct {
		name string
		key  string
		want string
	}{
		{
			name: "existing key with value",
			key:  "french",
			want: "bonjour",
		},
		{
			name: "existing key with empty value",
			key:  "notes",
			want: "",
		},
		{
			name: "non-existing key returns empty string",
			key:  "spanish",
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := entry.GetValue(tt.key)

			if got != tt.want {
				t.Errorf("GetValue(%v) = %v, want %v", tt.key, got, tt.want)
			}
		})
	}
}

func TestDataEntry_SetValue(t *testing.T) {
	entry := models.NewDataEntry(map[string]string{
		"existing": "old_value",
	}, "test.csv", 1)

	tests := []struct {
		name  string
		key   string
		value string
	}{
		{
			name:  "set new key",
			key:   "new_key",
			value: "new_value",
		},
		{
			name:  "update existing key",
			key:   "existing",
			value: "updated_value",
		},
		{
			name:  "set empty value",
			key:   "empty_key",
			value: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			entry.SetValue(tt.key, tt.value)

			got := entry.GetValue(tt.key)
			if got != tt.value {
				t.Errorf("SetValue(%v, %v) - got %v, want %v", tt.key, tt.value, got, tt.value)
			}
		})
	}
}

func TestDataEntry_SetValue_NilValues(t *testing.T) {
	// Test setting value when Values map is nil
	entry := &models.DataEntry{
		Values:     nil,
		Source:     "test.csv",
		LineNumber: 1,
	}

	entry.SetValue("test_key", "test_value")

	if entry.Values == nil {
		t.Error("SetValue() did not initialize nil Values map")
	}

	got := entry.GetValue("test_key")
	if got != "test_value" {
		t.Errorf("SetValue() with nil Values - got %v, want test_value", got)
	}
}

func TestDataEntry_Validate(t *testing.T) {
	tests := []struct {
		name        string
		setupFunc   func() *models.DataEntry
		wantErr     bool
		errContains string
	}{
		{
			name: "valid entry",
			setupFunc: func() *models.DataEntry {
				return models.NewDataEntry(
					map[string]string{"french": "bonjour", "english": "hello"},
					"test.csv",
					1,
				)
			},
			wantErr: false,
		},
		{
			name: "empty values map",
			setupFunc: func() *models.DataEntry {
				return models.NewDataEntry(
					map[string]string{},
					"test.csv",
					1,
				)
			},
			wantErr:     true,
			errContains: "must contain at least one field",
		},
		{
			name: "nil values map",
			setupFunc: func() *models.DataEntry {
				return &models.DataEntry{
					Values:     nil,
					Source:     "test.csv",
					LineNumber: 1,
				}
			},
			wantErr:     true,
			errContains: "must contain at least one field",
		},
		{
			name: "empty source",
			setupFunc: func() *models.DataEntry {
				return models.NewDataEntry(
					map[string]string{"field": "value"},
					"",
					1,
				)
			},
			wantErr:     true,
			errContains: "must reference source file",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			entry := tt.setupFunc()
			err := entry.Validate()

			if tt.wantErr {
				if err == nil {
					t.Error("Validate() error = nil, want error")
				} else if tt.errContains != "" && !containsString(err.Error(), tt.errContains) {
					t.Errorf("Validate() error = %v, want error containing %v", err, tt.errContains)
				}
			} else {
				if err != nil {
					t.Errorf("Validate() error = %v, want nil", err)
				}
			}
		})
	}
}

func TestDataEntry_GetHash(t *testing.T) {
	tests := []struct {
		name   string
		entry1 map[string]string
		entry2 map[string]string
		want   bool // true if hashes should be equal
	}{
		{
			name:   "identical entries",
			entry1: map[string]string{"french": "bonjour", "english": "hello"},
			entry2: map[string]string{"french": "bonjour", "english": "hello"},
			want:   true,
		},
		{
			name:   "same entries different order",
			entry1: map[string]string{"french": "bonjour", "english": "hello"},
			entry2: map[string]string{"english": "hello", "french": "bonjour"},
			want:   true,
		},
		{
			name:   "different values",
			entry1: map[string]string{"french": "bonjour", "english": "hello"},
			entry2: map[string]string{"french": "bonjour", "english": "goodbye"},
			want:   false,
		},
		{
			name:   "different keys",
			entry1: map[string]string{"french": "bonjour", "english": "hello"},
			entry2: map[string]string{"french": "bonjour", "spanish": "hola"},
			want:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e1 := models.NewDataEntry(tt.entry1, "test.csv", 1)
			e2 := models.NewDataEntry(tt.entry2, "test.csv", 2)

			hash1 := e1.GetHash()
			hash2 := e2.GetHash()

			if hash1 == "" || hash2 == "" {
				t.Fatal("GetHash() returned empty hash")
			}

			equal := hash1 == hash2
			if equal != tt.want {
				t.Errorf("GetHash() equality = %v, want %v (hash1: %v, hash2: %v)", equal, tt.want, hash1, hash2)
			}
		})
	}
}

func TestDataEntry_IsExactDuplicate(t *testing.T) {
	tests := []struct {
		name   string
		entry1 map[string]string
		entry2 map[string]string
		want   bool
	}{
		{
			name:   "identical entries",
			entry1: map[string]string{"french": "bonjour", "english": "hello"},
			entry2: map[string]string{"french": "bonjour", "english": "hello"},
			want:   true,
		},
		{
			name:   "different values",
			entry1: map[string]string{"french": "bonjour", "english": "hello"},
			entry2: map[string]string{"french": "bonjour", "english": "goodbye"},
			want:   false,
		},
		{
			name:   "different keys",
			entry1: map[string]string{"french": "bonjour", "english": "hello"},
			entry2: map[string]string{"french": "bonjour", "spanish": "hola"},
			want:   false,
		},
		{
			name:   "different lengths",
			entry1: map[string]string{"french": "bonjour"},
			entry2: map[string]string{"french": "bonjour", "english": "hello"},
			want:   false,
		},
		{
			name:   "both empty",
			entry1: map[string]string{},
			entry2: map[string]string{},
			want:   true,
		},
		{
			name:   "case sensitive comparison",
			entry1: map[string]string{"french": "Bonjour"},
			entry2: map[string]string{"french": "bonjour"},
			want:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e1 := models.NewDataEntry(tt.entry1, "test1.csv", 1)
			e2 := models.NewDataEntry(tt.entry2, "test2.csv", 2)

			got := e1.IsExactDuplicate(e2)

			if got != tt.want {
				t.Errorf("IsExactDuplicate() = %v, want %v", got, tt.want)
			}

			// Test symmetry
			got2 := e2.IsExactDuplicate(e1)
			if got2 != tt.want {
				t.Errorf("IsExactDuplicate() symmetry failed: e1.IsExactDuplicate(e2) = %v, e2.IsExactDuplicate(e1) = %v", got, got2)
			}
		})
	}
}

func TestDataEntry_ToCSVRecord(t *testing.T) {
	values := map[string]string{
		"french":  "bonjour",
		"english": "hello",
		"notes":   "greeting",
		"empty":   "",
	}
	entry := models.NewDataEntry(values, "test.csv", 1)

	tests := []struct {
		name    string
		columns []string
		want    []string
	}{
		{
			name:    "all existing columns",
			columns: []string{"french", "english", "notes"},
			want:    []string{"bonjour", "hello", "greeting"},
		},
		{
			name:    "different order",
			columns: []string{"english", "french"},
			want:    []string{"hello", "bonjour"},
		},
		{
			name:    "includes non-existing columns",
			columns: []string{"french", "spanish", "english"},
			want:    []string{"bonjour", "", "hello"},
		},
		{
			name:    "empty columns",
			columns: []string{},
			want:    []string{},
		},
		{
			name:    "includes empty value column",
			columns: []string{"french", "empty", "english"},
			want:    []string{"bonjour", "", "hello"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := entry.ToCSVRecord(tt.columns)

			if len(got) != len(tt.want) {
				t.Errorf("ToCSVRecord() length = %v, want %v", len(got), len(tt.want))
			}

			for i, expectedValue := range tt.want {
				if i >= len(got) {
					t.Errorf("ToCSVRecord() missing index %d", i)
				} else if got[i] != expectedValue {
					t.Errorf("ToCSVRecord()[%d] = %v, want %v", i, got[i], expectedValue)
				}
			}
		})
	}
}
