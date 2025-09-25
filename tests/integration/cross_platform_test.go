package integration

import (
	"ankiprep/internal/app"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

// TestCrossPlatformFileHandling tests file handling across different platforms
func TestCrossPlatformFileHandling(t *testing.T) {
	tests := []struct {
		name     string
		filePath string
		content  string
	}{
		{
			name:     "unix style path",
			filePath: "test_unix.csv",
			content:  "front,back\n\"Hello\",\"World\"\n\"Test\",\"Data\"\n",
		},
		{
			name:     "file with spaces",
			filePath: "test file with spaces.csv",
			content:  "front,back\n\"Space Test\",\"Another Test\"\n",
		},
		{
			name:     "file with unicode",
			filePath: "test_unicode_文件.csv",
			content:  "front,back\n\"Unicode\",\"测试\"\n\"français\",\"español\"\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test file in temp directory
			tempDir := t.TempDir()
			testFile := filepath.Join(tempDir, tt.filePath)

			// Write test content
			err := os.WriteFile(testFile, []byte(tt.content), 0644)
			if err != nil {
				t.Fatalf("Failed to create test file: %v", err)
			}

			// Process the file
			config := app.ProcessorConfig{
				OutputPath:     filepath.Join(tempDir, "output.csv"),
				FrenchMode:     false,
				SmartQuotes:    false,
				SkipDuplicates: false,
				Verbose:        false,
			}
			processor := app.NewProcessor(config)

			result, err := processor.ProcessFiles([]string{testFile})
			if err != nil {
				t.Fatalf("Processing failed for %s: %v", tt.filePath, err)
			}

			// Verify processing completed
			if result == nil {
				t.Fatal("Expected processing result")
			}

			if result.OutputRecords <= 0 {
				t.Error("Expected processed records")
			}

			// Verify output file exists and is readable
			if _, err := os.Stat(config.OutputPath); os.IsNotExist(err) {
				t.Errorf("Output file not created: %s", config.OutputPath)
			}

			t.Logf("Successfully processed %s: %d records", tt.filePath, result.OutputRecords)
		})
	}
}

// TestLineEndingHandling tests handling of different line endings across platforms
func TestLineEndingHandling(t *testing.T) {
	tests := []struct {
		name        string
		lineEnding  string
		description string
	}{
		{
			name:        "unix_lf",
			lineEnding:  "\n",
			description: "Unix/Linux LF",
		},
		{
			name:        "windows_crlf",
			lineEnding:  "\r\n", 
			description: "Windows CRLF",
		},
		{
			name:        "mac_cr",
			lineEnding:  "\r",
			description: "Classic Mac CR",
		},
	}

	baseContent := []string{
		"front,back",
		"\"Hello\",\"World\"",
		"\"Line 1\",\"Data 1\"",
		"\"Line 2\",\"Data 2\"",
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create content with specific line ending
			content := strings.Join(baseContent, tt.lineEnding)

			// Create test file
			tempDir := t.TempDir()
			testFile := filepath.Join(tempDir, "test.csv")
			
			err := os.WriteFile(testFile, []byte(content), 0644)
			if err != nil {
				t.Fatalf("Failed to create test file: %v", err)
			}

			// Process the file
			config := app.ProcessorConfig{
				OutputPath:     filepath.Join(tempDir, "output.csv"),
				FrenchMode:     false,
				SmartQuotes:    false,
				SkipDuplicates: false,
				Verbose:        false,
			}
			processor := app.NewProcessor(config)

			result, err := processor.ProcessFiles([]string{testFile})
			if err != nil {
				t.Fatalf("Processing failed for %s line endings: %v", tt.description, err)
			}

			// Verify processing completed
			expectedRecords := len(baseContent) - 1 // Minus header
			if result.OutputRecords != expectedRecords {
				t.Errorf("Expected %d records, got %d for %s", expectedRecords, result.OutputRecords, tt.description)
			}

			t.Logf("Successfully processed %s line endings: %d records", tt.description, result.OutputRecords)
		})
	}
}

// TestPathSeparatorHandling tests handling of different path separators
func TestPathSeparatorHandling(t *testing.T) {
	// Skip this test on Windows if we're testing Unix paths, and vice versa
	currentOS := runtime.GOOS

	tests := []struct {
		name           string
		createSubdirs  bool
		skipOn         []string
		expectedOutput bool
	}{
		{
			name:           "nested_directories",
			createSubdirs:  true,
			skipOn:         []string{}, // Run on all platforms
			expectedOutput: true,
		},
	}

	for _, tt := range tests {
		// Skip test on specified platforms
		skip := false
		for _, skipOS := range tt.skipOn {
			if currentOS == skipOS {
				skip = true
				break
			}
		}
		if skip {
			t.Skipf("Skipping %s on %s", tt.name, currentOS)
			continue
		}

		t.Run(tt.name, func(t *testing.T) {
			tempDir := t.TempDir()
			var testFile string

			if tt.createSubdirs {
				// Create nested directories
				subDir := filepath.Join(tempDir, "sub", "directory")
				err := os.MkdirAll(subDir, 0755)
				if err != nil {
					t.Fatalf("Failed to create subdirectories: %v", err)
				}
				testFile = filepath.Join(subDir, "test.csv")
			} else {
				testFile = filepath.Join(tempDir, "test.csv")
			}

			// Create test file
			content := "front,back\n\"Path Test\",\"Directory Test\"\n"
			err := os.WriteFile(testFile, []byte(content), 0644)
			if err != nil {
				t.Fatalf("Failed to create test file: %v", err)
			}

			// Process the file
			outputDir := filepath.Join(tempDir, "output")
			err = os.MkdirAll(outputDir, 0755)
			if err != nil {
				t.Fatalf("Failed to create output directory: %v", err)
			}

			config := app.ProcessorConfig{
				OutputPath:     filepath.Join(outputDir, "result.csv"),
				FrenchMode:     false,
				SmartQuotes:    false,
				SkipDuplicates: false,
				Verbose:        false,
			}
			processor := app.NewProcessor(config)

			result, err := processor.ProcessFiles([]string{testFile})
			if err != nil {
				t.Fatalf("Processing failed: %v", err)
			}

			if tt.expectedOutput {
				if result.OutputRecords <= 0 {
					t.Error("Expected processed records")
				}

				// Verify output file exists
				if _, err := os.Stat(config.OutputPath); os.IsNotExist(err) {
					t.Error("Expected output file to be created")
				}
			}

			t.Logf("Successfully processed nested path: %d records", result.OutputRecords)
		})
	}
}

// TestEncodingHandling tests handling of different text encodings
func TestEncodingHandling(t *testing.T) {
	tests := []struct {
		name        string
		content     string
		description string
	}{
		{
			name:        "ascii_content",
			content:     "front,back\n\"Hello\",\"World\"\n\"Test\",\"ASCII\"\n",
			description: "Pure ASCII content",
		},
		{
			name:        "utf8_content", 
			content:     "front,back\n\"Café\",\"Naïve\"\n\"测试\",\"тест\"\n\"🙂\",\"emoji\"\n",
			description: "UTF-8 with international characters and emoji",
		},
		{
			name:        "accented_content",
			content:     "front,back\n\"résumé\",\"café\"\n\"piñata\",\"jalapeño\"\n\"façade\",\"cliché\"\n",
			description: "Accented Latin characters",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test file
			tempDir := t.TempDir()
			testFile := filepath.Join(tempDir, "encoding_test.csv")

			err := os.WriteFile(testFile, []byte(tt.content), 0644)
			if err != nil {
				t.Fatalf("Failed to create test file: %v", err)
			}

			// Process the file
			config := app.ProcessorConfig{
				OutputPath:     filepath.Join(tempDir, "output.csv"),
				FrenchMode:     false,
				SmartQuotes:    true, // Enable for better text processing
				SkipDuplicates: false,
				Verbose:        false,
			}
			processor := app.NewProcessor(config)

			result, err := processor.ProcessFiles([]string{testFile})
			if err != nil {
				t.Fatalf("Processing failed for %s: %v", tt.description, err)
			}

			// Verify processing completed
			if result.OutputRecords <= 0 {
				t.Errorf("Expected processed records for %s", tt.description)
			}

			// Read output file to verify encoding is preserved
			outputContent, err := os.ReadFile(config.OutputPath)
			if err != nil {
				t.Fatalf("Failed to read output file: %v", err)
			}

			// Basic validation that output contains valid UTF-8
			outputStr := string(outputContent)
			if !strings.Contains(outputStr, "front") || !strings.Contains(outputStr, "back") {
				t.Errorf("Output file doesn't contain expected headers for %s", tt.description)
			}

			t.Logf("Successfully processed %s: %d records", tt.description, result.OutputRecords)
		})
	}
}

// TestFilePermissions tests handling of files with different permissions
func TestFilePermissions(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Skipping file permission tests on Windows")
	}

	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "permission_test.csv")
	content := "front,back\n\"Permission\",\"Test\"\n"

	// Create test file
	err := os.WriteFile(testFile, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	tests := []struct {
		name           string
		fileMode       os.FileMode
		shouldSucceed  bool
		description    string
	}{
		{
			name:          "readable_file",
			fileMode:      0644,
			shouldSucceed: true,
			description:   "Normal readable file",
		},
		{
			name:          "readonly_file",
			fileMode:      0444,
			shouldSucceed: true,
			description:   "Read-only file",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set file permissions
			err := os.Chmod(testFile, tt.fileMode)
			if err != nil {
				t.Fatalf("Failed to set file permissions: %v", err)
			}

			// Try to process the file
			config := app.ProcessorConfig{
				OutputPath:     filepath.Join(tempDir, fmt.Sprintf("output_%s.csv", tt.name)),
				FrenchMode:     false,
				SmartQuotes:    false,
				SkipDuplicates: false,
				Verbose:        false,
			}
			processor := app.NewProcessor(config)

			result, err := processor.ProcessFiles([]string{testFile})

			if tt.shouldSucceed {
				if err != nil {
					t.Errorf("Processing should have succeeded for %s: %v", tt.description, err)
				}
				if result == nil || result.OutputRecords <= 0 {
					t.Errorf("Expected valid result for %s", tt.description)
				}
			} else {
				if err == nil {
					t.Errorf("Processing should have failed for %s", tt.description)
				}
			}

			t.Logf("Test %s completed as expected", tt.description)
		})
	}
}

// TestPlatformSpecificFeatures tests platform-specific functionality
func TestPlatformSpecificFeatures(t *testing.T) {
	currentOS := runtime.GOOS

	t.Run("current_platform_basic_processing", func(t *testing.T) {
		// Create test file
		tempDir := t.TempDir()
		testFile := filepath.Join(tempDir, "platform_test.csv")
		content := fmt.Sprintf("front,back\n\"Platform\",\"%s\"\n\"Architecture\",\"%s\"\n", 
			currentOS, runtime.GOARCH)

		err := os.WriteFile(testFile, []byte(content), 0644)
		if err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}

		// Process the file
		config := app.ProcessorConfig{
			OutputPath:     filepath.Join(tempDir, "platform_output.csv"),
			FrenchMode:     false,
			SmartQuotes:    false,
			SkipDuplicates: false,
			Verbose:        false,
		}
		processor := app.NewProcessor(config)

		result, err := processor.ProcessFiles([]string{testFile})
		if err != nil {
			t.Fatalf("Basic processing failed on %s/%s: %v", currentOS, runtime.GOARCH, err)
		}

		if result.OutputRecords != 2 {
			t.Errorf("Expected 2 records, got %d", result.OutputRecords)
		}

		t.Logf("Platform %s/%s: Successfully processed %d records", 
			currentOS, runtime.GOARCH, result.OutputRecords)
	})
}

// TestResourceLimits tests behavior under resource constraints
func TestResourceLimits(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping resource limit tests in short mode")
	}

	t.Run("many_small_files", func(t *testing.T) {
		tempDir := t.TempDir()
		numFiles := 10
		var inputFiles []string

		// Create multiple small CSV files
		for i := 0; i < numFiles; i++ {
			filename := filepath.Join(tempDir, fmt.Sprintf("file_%d.csv", i))
			content := fmt.Sprintf("front,back\n\"File %d\",\"Content %d\"\n\"Data %d\",\"Test %d\"\n", i, i, i, i)
			
			err := os.WriteFile(filename, []byte(content), 0644)
			if err != nil {
				t.Fatalf("Failed to create file %d: %v", i, err)
			}
			inputFiles = append(inputFiles, filename)
		}

		// Process all files together
		config := app.ProcessorConfig{
			OutputPath:     filepath.Join(tempDir, "combined_output.csv"),
			FrenchMode:     false,
			SmartQuotes:    false,
			SkipDuplicates: true, // Enable duplicate detection across files
			Verbose:        false,
		}
		processor := app.NewProcessor(config)

		result, err := processor.ProcessFiles(inputFiles)
		if err != nil {
			t.Fatalf("Multi-file processing failed: %v", err)
		}

		// Should process 2 records per file = 20 total
		expectedRecords := numFiles * 2
		if result.OutputRecords != expectedRecords {
			t.Errorf("Expected %d records, got %d", expectedRecords, result.OutputRecords)
		}

		t.Logf("Successfully processed %d files with %d total records", numFiles, result.OutputRecords)
	})
}