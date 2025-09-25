package unit_test

import (
	"os"
	"path/filepath"
	"testing"

	"ankiprep/internal/services"
)

func TestCSVParser_NewCSVParser(t *testing.T) {
	parser := services.NewCSVParser()

	if parser == nil {
		t.Fatal("NewCSVParser returned nil")
	}
}

func TestCSVParser_ValidateFormat(t *testing.T) {
	parser := services.NewCSVParser()
	tempDir := t.TempDir()

	tests := []struct {
		name        string
		setupFunc   func() string // Returns file path
		wantErr     bool
		errContains string
	}{
		{
			name: "valid CSV file",
			setupFunc: func() string {
				testFile := filepath.Join(tempDir, "test.csv")
				content := "header1,header2\nvalue1,value2\n"
				os.WriteFile(testFile, []byte(content), 0644)
				return testFile
			},
			wantErr: false,
		},
		{
			name: "valid TSV file",
			setupFunc: func() string {
				testFile := filepath.Join(tempDir, "test.tsv")
				content := "header1\theader2\nvalue1\tvalue2\n"
				os.WriteFile(testFile, []byte(content), 0644)
				return testFile
			},
			wantErr: false,
		},
		{
			name: "unsupported extension",
			setupFunc: func() string {
				testFile := filepath.Join(tempDir, "test.txt")
				content := "header1,header2\nvalue1,value2\n"
				os.WriteFile(testFile, []byte(content), 0644)
				return testFile
			},
			wantErr:     true,
			errContains: "unsupported file format",
		},
		{
			name: "file does not exist",
			setupFunc: func() string {
				return "/nonexistent/file.csv"
			},
			wantErr:     true,
			errContains: "cannot access file",
		},
		{
			name: "empty file",
			setupFunc: func() string {
				testFile := filepath.Join(tempDir, "empty.csv")
				os.WriteFile(testFile, []byte(""), 0644)
				return testFile
			},
			wantErr:     true,
			errContains: "is empty",
		},
		{
			name: "file with no columns (just newlines)",
			setupFunc: func() string {
				testFile := filepath.Join(tempDir, "nocols.csv")
				os.WriteFile(testFile, []byte("\n\n"), 0644)
				return testFile
			},
			wantErr:     true,
			errContains: "is empty", // The actual error from CSV parser for empty content
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filePath := tt.setupFunc()
			err := parser.ValidateFormat(filePath)

			if tt.wantErr {
				if err == nil {
					t.Error("ValidateFormat() error = nil, want error")
				} else if tt.errContains != "" && !containsString(err.Error(), tt.errContains) {
					t.Errorf("ValidateFormat() error = %v, want error containing %v", err, tt.errContains)
				}
			} else {
				if err != nil {
					t.Errorf("ValidateFormat() error = %v, want nil", err)
				}
			}
		})
	}
}

func TestCSVParser_ParseFile(t *testing.T) {
	parser := services.NewCSVParser()
	tempDir := t.TempDir()

	tests := []struct {
		name            string
		setupFunc       func() string // Returns file path
		wantHeaders     []string
		wantRecordCount int
		wantSeparator   rune
		wantErr         bool
		errContains     string
	}{
		{
			name: "simple CSV file",
			setupFunc: func() string {
				testFile := filepath.Join(tempDir, "simple.csv")
				content := "french,english\nbonjour,hello\nau revoir,goodbye\n"
				os.WriteFile(testFile, []byte(content), 0644)
				return testFile
			},
			wantHeaders:     []string{"french", "english"},
			wantRecordCount: 2,
			wantSeparator:   ',',
			wantErr:         false,
		},
		{
			name: "TSV file",
			setupFunc: func() string {
				testFile := filepath.Join(tempDir, "simple.tsv")
				content := "french\tenglish\nbonjour\thello\nau revoir\tgoodbye\n"
				os.WriteFile(testFile, []byte(content), 0644)
				return testFile
			},
			wantHeaders:     []string{"french", "english"},
			wantRecordCount: 2,
			wantSeparator:   '\t',
			wantErr:         false,
		},
		{
			name: "CSV with quoted fields",
			setupFunc: func() string {
				testFile := filepath.Join(tempDir, "quoted.csv")
				content := "\"Name\",\"Description\"\n\"John, Jr.\",\"A person with comma in name\"\n"
				os.WriteFile(testFile, []byte(content), 0644)
				return testFile
			},
			wantHeaders:     []string{"Name", "Description"},
			wantRecordCount: 1,
			wantSeparator:   ',',
			wantErr:         false,
		},
		{
			name: "CSV with empty cells",
			setupFunc: func() string {
				testFile := filepath.Join(tempDir, "empty_cells.csv")
				content := "col1,col2,col3\nvalue1,,value3\n,value2,\n"
				os.WriteFile(testFile, []byte(content), 0644)
				return testFile
			},
			wantHeaders:     []string{"col1", "col2", "col3"},
			wantRecordCount: 2,
			wantSeparator:   ',',
			wantErr:         false,
		},
		{
			name: "file with only headers",
			setupFunc: func() string {
				testFile := filepath.Join(tempDir, "headers_only.csv")
				content := "header1,header2\n"
				os.WriteFile(testFile, []byte(content), 0644)
				return testFile
			},
			wantErr:     true,
			errContains: "contains no data rows",
		},
		{
			name: "completely empty file",
			setupFunc: func() string {
				testFile := filepath.Join(tempDir, "completely_empty.csv")
				content := ""
				os.WriteFile(testFile, []byte(content), 0644)
				return testFile
			},
			wantErr:     true,
			errContains: "contains no data",
		},
		{
			name: "nonexistent file",
			setupFunc: func() string {
				return "/nonexistent/path/file.csv"
			},
			wantErr:     true,
			errContains: "failed to open file",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filePath := tt.setupFunc()
			inputFile, err := parser.ParseFile(filePath)

			if tt.wantErr {
				if err == nil {
					t.Error("ParseFile() error = nil, want error")
				} else if tt.errContains != "" && !containsString(err.Error(), tt.errContains) {
					t.Errorf("ParseFile() error = %v, want error containing %v", err, tt.errContains)
				}
				return
			}

			if err != nil {
				t.Fatalf("ParseFile() error = %v, want nil", err)
			}

			if inputFile == nil {
				t.Fatal("ParseFile() returned nil InputFile")
			}

			// Check headers
			if len(inputFile.Headers) != len(tt.wantHeaders) {
				t.Errorf("ParseFile() headers length = %v, want %v", len(inputFile.Headers), len(tt.wantHeaders))
			} else {
				for i, expectedHeader := range tt.wantHeaders {
					if inputFile.Headers[i] != expectedHeader {
						t.Errorf("ParseFile() headers[%d] = %v, want %v", i, inputFile.Headers[i], expectedHeader)
					}
				}
			}

			// Check record count
			if len(inputFile.Records) != tt.wantRecordCount {
				t.Errorf("ParseFile() record count = %v, want %v", len(inputFile.Records), tt.wantRecordCount)
			}

			// Check separator
			if inputFile.Separator != tt.wantSeparator {
				t.Errorf("ParseFile() separator = %v, want %v", inputFile.Separator, tt.wantSeparator)
			}

			// Check path is set
			if inputFile.Path != filePath {
				t.Errorf("ParseFile() path = %v, want %v", inputFile.Path, filePath)
			}
		})
	}
}

func TestCSVParser_ParseToDataEntries(t *testing.T) {
	parser := services.NewCSVParser()
	tempDir := t.TempDir()

	tests := []struct {
		name           string
		setupFunc      func() string // Returns file path
		wantEntryCount int
		wantFirstEntry map[string]string
		wantErr        bool
		errContains    string
	}{
		{
			name: "normal CSV conversion",
			setupFunc: func() string {
				testFile := filepath.Join(tempDir, "normal.csv")
				content := "french,english,notes\nbonjour,hello,greeting\nau revoir,goodbye,farewell\n"
				os.WriteFile(testFile, []byte(content), 0644)
				return testFile
			},
			wantEntryCount: 2,
			wantFirstEntry: map[string]string{
				"french":  "bonjour",
				"english": "hello",
				"notes":   "greeting",
			},
			wantErr: false,
		},
		{
			name: "CSV with empty rows (should be skipped)",
			setupFunc: func() string {
				testFile := filepath.Join(tempDir, "empty_rows.csv")
				content := "col1,col2\nvalue1,value2\n,\n   ,   \nvalue3,value4\n"
				os.WriteFile(testFile, []byte(content), 0644)
				return testFile
			},
			wantEntryCount: 2, // Empty rows should be skipped
			wantFirstEntry: map[string]string{
				"col1": "value1",
				"col2": "value2",
			},
			wantErr: false,
		},
		{
			name: "CSV with consistent columns and empty values",
			setupFunc: func() string {
				testFile := filepath.Join(tempDir, "consistent.csv")
				content := "col1,col2,col3\nvalue1,,value3\n,value2,\n"
				os.WriteFile(testFile, []byte(content), 0644)
				return testFile
			},
			wantEntryCount: 2,
			wantFirstEntry: map[string]string{
				"col1": "value1",
				"col2": "",
				"col3": "value3",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filePath := tt.setupFunc()

			// First parse the file
			inputFile, err := parser.ParseFile(filePath)
			if err != nil {
				t.Fatalf("ParseFile() error = %v, want nil for setup", err)
			}

			// Then convert to data entries
			entries, err := parser.ParseToDataEntries(inputFile)

			if tt.wantErr {
				if err == nil {
					t.Error("ParseToDataEntries() error = nil, want error")
				} else if tt.errContains != "" && !containsString(err.Error(), tt.errContains) {
					t.Errorf("ParseToDataEntries() error = %v, want error containing %v", err, tt.errContains)
				}
				return
			}

			if err != nil {
				t.Fatalf("ParseToDataEntries() error = %v, want nil", err)
			}

			// Check entry count
			if len(entries) != tt.wantEntryCount {
				t.Errorf("ParseToDataEntries() entry count = %v, want %v", len(entries), tt.wantEntryCount)
			}

			if len(entries) > 0 && tt.wantFirstEntry != nil {
				firstEntry := entries[0]

				// Check source and line number
				if firstEntry.Source != filePath {
					t.Errorf("ParseToDataEntries() first entry source = %v, want %v", firstEntry.Source, filePath)
				}

				if firstEntry.LineNumber != 2 { // Should be 2 (header is line 1)
					t.Errorf("ParseToDataEntries() first entry line number = %v, want 2", firstEntry.LineNumber)
				}

				// Check values
				for key, expectedValue := range tt.wantFirstEntry {
					actualValue := firstEntry.GetValue(key)
					if actualValue != expectedValue {
						t.Errorf("ParseToDataEntries() first entry[%s] = %v, want %v", key, actualValue, expectedValue)
					}
				}
			}
		})
	}
}

func TestCSVParser_ParseToDataEntries_NilInput(t *testing.T) {
	parser := services.NewCSVParser()

	_, err := parser.ParseToDataEntries(nil)

	if err == nil {
		t.Error("ParseToDataEntries(nil) error = nil, want error")
	} else if !containsString(err.Error(), "cannot be nil") {
		t.Errorf("ParseToDataEntries(nil) error = %v, want error containing 'cannot be nil'", err)
	}
}

func TestCSVParser_Integration(t *testing.T) {
	// Test the full workflow: validate -> parse -> convert to entries
	parser := services.NewCSVParser()
	tempDir := t.TempDir()

	// Create a test file with various scenarios
	testFile := filepath.Join(tempDir, "integration.csv")
	content := `"French","English","Notes","Category"
"bonjour","hello","greeting","basics"
"au revoir","goodbye","farewell","basics"
"merci","thank you","gratitude","politeness"
"s'il vous plaît","please","request","politeness"`

	err := os.WriteFile(testFile, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Step 1: Validate format
	err = parser.ValidateFormat(testFile)
	if err != nil {
		t.Errorf("ValidateFormat() error = %v, want nil", err)
	}

	// Step 2: Parse file
	inputFile, err := parser.ParseFile(testFile)
	if err != nil {
		t.Fatalf("ParseFile() error = %v, want nil", err)
	}

	expectedHeaders := []string{"French", "English", "Notes", "Category"}
	if len(inputFile.Headers) != len(expectedHeaders) {
		t.Fatalf("Headers length = %v, want %v", len(inputFile.Headers), len(expectedHeaders))
	}

	// Step 3: Convert to data entries
	entries, err := parser.ParseToDataEntries(inputFile)
	if err != nil {
		t.Fatalf("ParseToDataEntries() error = %v, want nil", err)
	}

	if len(entries) != 4 {
		t.Errorf("Entry count = %v, want 4", len(entries))
	}

	// Validate first entry
	firstEntry := entries[0]
	if firstEntry.GetValue("French") != "bonjour" {
		t.Errorf("First entry French = %v, want 'bonjour'", firstEntry.GetValue("French"))
	}

	if firstEntry.GetValue("English") != "hello" {
		t.Errorf("First entry English = %v, want 'hello'", firstEntry.GetValue("English"))
	}

	// Validate that entries have proper source information
	for i, entry := range entries {
		if entry.Source != testFile {
			t.Errorf("Entry[%d] source = %v, want %v", i, entry.Source, testFile)
		}

		expectedLineNumber := i + 2 // +2 because header is line 1, data starts at line 2
		if entry.LineNumber != expectedLineNumber {
			t.Errorf("Entry[%d] line number = %v, want %v", i, entry.LineNumber, expectedLineNumber)
		}
	}
}
