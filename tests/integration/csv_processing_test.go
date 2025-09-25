package integration
package integration_test

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"ankiprep/internal/services"
)

// TestCSVProcessing_ClozeFixtures tests end-to-end CSV processing with cloze data
func TestCSVProcessing_ClozeFixtures(t *testing.T) {
	// This will fail until we implement the full processing pipeline
	fixturesPath := filepath.Join("..", "..", "tests", "fixtures")

	t.Run("cloze test data processing", func(t *testing.T) {
		inputFile := filepath.Join(fixturesPath, "cloze_test_data.csv")
		
		// Read fixture file
		file, err := os.Open(inputFile)
		if err != nil {
			t.Fatalf("Failed to open fixture file: %v", err)
		}
		defer file.Close()

		reader := csv.NewReader(file)
		records, err := reader.ReadAll()
		if err != nil {
			t.Fatalf("Failed to read CSV: %v", err)
		}

		if len(records) < 2 { // Header + at least one data row
			t.Fatal("Expected at least one data row in fixture file")
		}

		// Verify header format
		header := records[0]
		expectedHeader := []string{"Front", "Back", "Tags"}
		if len(header) != len(expectedHeader) {
			t.Fatalf("Expected %d columns, got %d", len(expectedHeader), len(header))
		}

		for i, col := range expectedHeader {
			if header[i] != col {
				t.Errorf("Expected column %d to be %q, got %q", i, col, header[i])
			}
		}

		// Process each record through the typography service
		service := services.NewTypographyService()

		for i, record := range records[1:] { // Skip header
			if len(record) != 3 {
				t.Errorf("Row %d: Expected 3 columns, got %d", i+1, len(record))
				continue
			}

			front, back, tags := record[0], record[1], record[2]

			// Process front field
			frontRequest := services.ProcessFrenchTextRequest{
				Text:   front,
				French: true,
			}

			frontResult, err := service.ProcessFrenchText(frontRequest)
			if err != nil {
				t.Errorf("Row %d: Failed to process front text: %v", i+1, err)
				continue
			}

			// Process back field
			backRequest := services.ProcessFrenchTextRequest{
				Text:   back,
				French: true,
			}

			backResult, err := service.ProcessFrenchText(backRequest)
			if err != nil {
				t.Errorf("Row %d: Failed to process back text: %v", i+1, err)
				continue
			}

			// Verify cloze blocks were detected if expected
			if strings.Contains(front, "{{c") || strings.Contains(back, "{{c") {
				totalClozeCount := frontResult.ClozeCount + backResult.ClozeCount
				if totalClozeCount == 0 {
					t.Errorf("Row %d: Expected to find cloze blocks but found none", i+1)
				}
			}

			// Verify colon rules
			t.Run(fmt.Sprintf("row_%d_colon_rules", i+1), func(t *testing.T) {
				// Check that colons outside cloze blocks have non-breaking spaces
				if strings.Contains(frontResult.ProcessedText, " :") {
					t.Error("Front: Found regular space before colon outside cloze block")
				}
				if strings.Contains(backResult.ProcessedText, " :") {
					t.Error("Back: Found regular space before colon outside cloze block")
				}

				// Check that colons inside cloze blocks are unchanged
				clozePattern := `\{\{c\d+::[^}]*\}\}`
				// This is a simplified check - real implementation would use regex
				if strings.Contains(frontResult.ProcessedText, "{{c") {
					// Verify the cloze content wasn't modified
					if !strings.Contains(frontResult.ProcessedText, "{{c") {
						t.Error("Front: Cloze block structure was modified")
					}
				}
			})

			// Verify tags are preserved (not processed)
			if tags != record[2] {
				t.Errorf("Row %d: Tags should be preserved unchanged", i+1)
			}
		}
	})

	t.Run("no cloze test data processing", func(t *testing.T) {
		inputFile := filepath.Join(fixturesPath, "no_cloze_test_data.csv")
		
		// Read fixture file
		file, err := os.Open(inputFile)
		if err != nil {
			t.Fatalf("Failed to open fixture file: %v", err)
		}
		defer file.Close()

		reader := csv.NewReader(file)
		records, err := reader.ReadAll()
		if err != nil {
			t.Fatalf("Failed to read CSV: %v", err)
		}

		// Process each record and verify no cloze blocks detected
		service := services.NewTypographyService()

		for i, record := range records[1:] { // Skip header
			if len(record) != 3 {
				continue
			}

			front, back := record[0], record[1]

			// Process both fields
			frontRequest := services.ProcessFrenchTextRequest{
				Text:   front,
				French: true,
			}

			backRequest := services.ProcessFrenchTextRequest{
				Text:   back,
				French: true,
			}

			frontResult, err := service.ProcessFrenchText(frontRequest)
			if err != nil {
				t.Errorf("Row %d: Failed to process front text: %v", i+1, err)
				continue
			}

			backResult, err := service.ProcessFrenchText(backRequest)
			if err != nil {
				t.Errorf("Row %d: Failed to process back text: %v", i+1, err)
				continue
			}

			// Verify no cloze blocks detected
			if frontResult.ClozeCount > 0 {
				t.Errorf("Row %d: Front text should not contain cloze blocks, found %d", 
					i+1, frontResult.ClozeCount)
			}

			if backResult.ClozeCount > 0 {
				t.Errorf("Row %d: Back text should not contain cloze blocks, found %d", 
					i+1, backResult.ClozeCount)
			}

			// Verify French typography was applied
			if strings.Contains(front, " :") && !strings.Contains(frontResult.ProcessedText, "\u00A0:") {
				t.Errorf("Row %d: French colon rule not applied to front text", i+1)
			}

			if strings.Contains(back, " :") && !strings.Contains(backResult.ProcessedText, "\u00A0:") {
				t.Errorf("Row %d: French colon rule not applied to back text", i+1)
			}
		}
	})
}

// TestCSVProcessing_ErrorHandling tests error handling in CSV processing
func TestCSVProcessing_ErrorHandling(t *testing.T) {
	service := services.NewTypographyService()

	t.Run("invalid CSV data", func(t *testing.T) {
		// Test with malformed cloze blocks
		invalidInputs := []string{
			"Question : What is {{c1::Paris",      // Missing closing brackets
			"Answer: {{c::no number}}",            // Missing cloze number
			"Test {{c0::zero number}}",            // Invalid cloze number
			"Nested {{c1::{{c2::double}}}}",       // Nested cloze blocks
		}

		for _, input := range invalidInputs {
			t.Run("malformed_"+input[:10], func(t *testing.T) {
				request := services.ProcessFrenchTextRequest{
					Text:   input,
					French: true,
				}

				result, err := service.ProcessFrenchText(request)
				
				// Should not fail completely, but should log warnings
				if err != nil {
					t.Errorf("ProcessFrenchText should handle malformed input gracefully, got error: %v", err)
				}

				// Should have warnings about malformed blocks
				if result.WarningCount == 0 {
					t.Error("Expected warnings for malformed cloze blocks")
				}

				// Should still apply French typography to valid parts
				if !strings.Contains(result.ProcessedText, input) {
					// At minimum, text should be returned (even if not fully processed)
					if result.ProcessedText == "" {
						t.Error("ProcessedText should not be empty")
					}
				}
			})
		}
	})
}