package unit_test

import (
	"os"
	"path/filepath"
	"testing"

	"ankiprep/internal/models"
)

func TestInputFile_NewInputFile(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		wantPath string
	}{
		{
			name:     "valid file path",
			path:     "test.csv",
			wantPath: "test.csv",
		},
		{
			name:     "empty path",
			path:     "",
			wantPath: "",
		},
		{
			name:     "path with spaces",
			path:     "test file.csv",
			wantPath: "test file.csv",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inputFile := models.NewInputFile(tt.path)
			
			if inputFile == nil {
				t.Fatal("NewInputFile returned nil")
			}
			
			if inputFile.Path != tt.wantPath {
				t.Errorf("NewInputFile() path = %v, want %v", inputFile.Path, tt.wantPath)
			}
			
			// Check default values
			if inputFile.Separator != ',' {
				t.Errorf("NewInputFile() separator = %v, want ','", inputFile.Separator)
			}
			
			if inputFile.Encoding != "UTF-8" {
				t.Errorf("NewInputFile() encoding = %v, want 'UTF-8'", inputFile.Encoding)
			}
		})
	}
}

func TestInputFile_DetectSeparator(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		wantSep rune
	}{
		{
			name:    "CSV file extension",
			path:    "test.csv",
			wantSep: ',',
		},
		{
			name:    "TSV file extension",
			path:    "test.tsv",
			wantSep: '\t',
		},
		{
			name:    "uppercase CSV",
			path:    "test.CSV",
			wantSep: ',',
		},
		{
			name:    "uppercase TSV",
			path:    "test.TSV",
			wantSep: '\t',
		},
		{
			name:    "unknown extension defaults to comma",
			path:    "test.txt",
			wantSep: ',',
		},
		{
			name:    "no extension defaults to comma",
			path:    "test",
			wantSep: ',',
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inputFile := models.NewInputFile(tt.path)
			inputFile.DetectSeparator()
			
			if inputFile.Separator != tt.wantSep {
				t.Errorf("DetectSeparator() separator = %v, want %v", inputFile.Separator, tt.wantSep)
			}
		})
	}
}

func TestInputFile_GetSeparatorString(t *testing.T) {
	tests := []struct {
		name      string
		separator rune
		want      string
	}{
		{
			name:      "comma separator",
			separator: ',',
			want:      "comma",
		},
		{
			name:      "tab separator",
			separator: '\t',
			want:      "tab",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inputFile := models.NewInputFile("test.csv")
			inputFile.Separator = tt.separator
			
			got := inputFile.GetSeparatorString()
			if got != tt.want {
				t.Errorf("GetSeparatorString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInputFile_Validate(t *testing.T) {
	// Create a temporary test file
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "test.csv")
	
	// Write test content to file
	content := "header1,header2\nvalue1,value2\n"
	err := os.WriteFile(testFile, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	tests := []struct {
		name        string
		setupFunc   func(*models.InputFile)
		wantErr     bool
		errContains string
	}{
		{
			name: "valid input file",
			setupFunc: func(f *models.InputFile) {
				f.Path = testFile
				f.Headers = []string{"header1", "header2"}
				f.Records = [][]string{{"value1", "value2"}}
				f.Separator = ','
				f.Encoding = "UTF-8"
			},
			wantErr: false,
		},
		{
			name: "file does not exist",
			setupFunc: func(f *models.InputFile) {
				f.Path = "/nonexistent/file.csv"
				f.Headers = []string{"header1", "header2"}
				f.Records = [][]string{{"value1", "value2"}}
			},
			wantErr:     true,
			errContains: "file not found",
		},
		{
			name: "invalid separator",
			setupFunc: func(f *models.InputFile) {
				f.Path = testFile
				f.Headers = []string{"header1", "header2"}
				f.Records = [][]string{{"value1", "value2"}}
				f.Separator = ';' // Invalid separator
				f.Encoding = "UTF-8"
			},
			wantErr:     true,
			errContains: "invalid separator",
		},
		{
			name: "invalid encoding",
			setupFunc: func(f *models.InputFile) {
				f.Path = testFile
				f.Headers = []string{"header1", "header2"}
				f.Records = [][]string{{"value1", "value2"}}
				f.Separator = ','
				f.Encoding = "ISO-8859-1" // Invalid encoding
			},
			wantErr:     true,
			errContains: "invalid encoding",
		},
		{
			name: "no data rows",
			setupFunc: func(f *models.InputFile) {
				f.Path = testFile
				f.Headers = []string{"header1", "header2"}
				f.Records = [][]string{} // Empty records
				f.Separator = ','
				f.Encoding = "UTF-8"
			},
			wantErr:     true,
			errContains: "no data rows",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inputFile := models.NewInputFile("")
			tt.setupFunc(inputFile)
			
			err := inputFile.Validate()
			
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

func TestInputFile_RecordsManipulation(t *testing.T) {
	inputFile := models.NewInputFile("test.csv")
	
	// Test empty records initially
	if len(inputFile.Records) != 0 {
		t.Errorf("NewInputFile() initial records length = %v, want 0", len(inputFile.Records))
	}
	
	// Test adding records directly
	inputFile.Records = append(inputFile.Records, []string{"value1", "value2"})
	inputFile.Records = append(inputFile.Records, []string{"value3", "value4"})
	
	if len(inputFile.Records) != 2 {
		t.Errorf("Records length = %v, want 2", len(inputFile.Records))
	}
	
	// Verify record content
	if inputFile.Records[0][0] != "value1" || inputFile.Records[0][1] != "value2" {
		t.Errorf("First record = %v, want [value1, value2]", inputFile.Records[0])
	}
	
	if inputFile.Records[1][0] != "value3" || inputFile.Records[1][1] != "value4" {
		t.Errorf("Second record = %v, want [value3, value4]", inputFile.Records[1])
	}
}

func TestInputFile_HeadersManipulation(t *testing.T) {
	inputFile := models.NewInputFile("test.csv")
	
	// Test empty headers initially
	if len(inputFile.Headers) != 0 {
		t.Errorf("NewInputFile() initial headers length = %v, want 0", len(inputFile.Headers))
	}
	
	// Test setting headers
	inputFile.Headers = []string{"col1", "col2", "col3"}
	
	if len(inputFile.Headers) != 3 {
		t.Errorf("Headers length = %v, want 3", len(inputFile.Headers))
	}
	
	// Verify header content
	expectedHeaders := []string{"col1", "col2", "col3"}
	for i, header := range inputFile.Headers {
		if header != expectedHeaders[i] {
			t.Errorf("Header[%d] = %v, want %v", i, header, expectedHeaders[i])
		}
	}
}

// Helper function to check if string contains substring
func containsString(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || 
		(len(s) > len(substr) && findSubstring(s, substr)))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}