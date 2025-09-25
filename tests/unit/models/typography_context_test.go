package models_test

import (
	"bytes"
	"log"
	"testing"

	"ankiprep/internal/models"
)

// TestTypographyContext_Contract verifies the TypographyContext contract
func TestTypographyContext_Contract(t *testing.T) {
	t.Run("valid context creation", func(t *testing.T) {
		var logBuffer bytes.Buffer
		logger := log.New(&logBuffer, "", 0)

		// This will fail until we implement the model
		context := models.TypographyContext{
			SourceText:    "Question : The capital is {{c1::Paris}}",
			ClozeBlocks:   []models.ClozeDeletionBlock{},
			FrenchEnabled: true,
			Logger:        logger,
		}

		if context.SourceText == "" {
			t.Error("SourceText should not be empty")
		}

		if !context.FrenchEnabled {
			t.Error("FrenchEnabled should be true")
		}

		if context.Logger == nil {
			t.Error("Logger should not be nil")
		}
	})

	t.Run("cloze blocks validation - no overlap", func(t *testing.T) {
		// Create overlapping blocks (invalid)
		block1 := models.ClozeDeletionBlock{
			FullText: "{{c1::test}}",
			Number:   1,
			Content:  "test",
			StartPos: 10,
			EndPos:   22,
		}

		block2 := models.ClozeDeletionBlock{
			FullText: "{{c2::overlap}}",
			Number:   2,
			Content:  "overlap",
			StartPos: 20, // Overlaps with block1
			EndPos:   35,
		}

		context := models.TypographyContext{
			SourceText:  "Some text {{c1::test}} and {{c2::overlap}} here",
			ClozeBlocks: []models.ClozeDeletionBlock{block1, block2},
		}

		// This will fail until we implement validation
		err := context.Validate()
		if err == nil {
			t.Error("Expected validation error for overlapping blocks")
		}
	})

	t.Run("cloze blocks validation - proper sorting", func(t *testing.T) {
		// Create blocks in wrong order
		block1 := models.ClozeDeletionBlock{
			FullText: "{{c1::first}}",
			Number:   1,
			Content:  "first",
			StartPos: 30,
			EndPos:   43,
		}

		block2 := models.ClozeDeletionBlock{
			FullText: "{{c2::second}}",
			Number:   2,
			Content:  "second",
			StartPos: 10,
			EndPos:   24,
		}

		// Blocks are not sorted by StartPos
		context := models.TypographyContext{
			SourceText:  "Some text {{c2::second}} and {{c1::first}} here",
			ClozeBlocks: []models.ClozeDeletionBlock{block1, block2}, // Wrong order
		}

		// This will fail until we implement validation
		err := context.Validate()
		if err == nil {
			t.Error("Expected validation error for unsorted blocks")
		}
	})

	t.Run("cloze blocks validation - position bounds", func(t *testing.T) {
		// Create block with positions outside source text
		block := models.ClozeDeletionBlock{
			FullText: "{{c1::test}}",
			Number:   1,
			Content:  "test",
			StartPos: 0,
			EndPos:   50, // Beyond source text length
		}

		context := models.TypographyContext{
			SourceText:  "Short text",
			ClozeBlocks: []models.ClozeDeletionBlock{block},
		}

		// This will fail until we implement validation
		err := context.Validate()
		if err == nil {
			t.Error("Expected validation error for out-of-bounds positions")
		}
	})
}

// TestTypographyResult_Contract verifies the TypographyResult contract
func TestTypographyResult_Contract(t *testing.T) {
	t.Run("result guarantees", func(t *testing.T) {
		// This will fail until we implement the model
		result := models.TypographyResult{
			ProcessedText: "Question\u00A0: The capital is {{c1::Paris}}",
			ClozeCount:    1,
			WarningCount:  0,
			Warnings:      []string{},
		}

		// Test contract guarantees
		if result.ProcessedText == "" {
			t.Error("ProcessedText should not be empty when input was non-empty")
		}

		if result.ClozeCount < 0 {
			t.Error("ClozeCount should be >= 0")
		}

		if result.WarningCount != len(result.Warnings) {
			t.Errorf("WarningCount (%d) should equal len(Warnings) (%d)", 
				result.WarningCount, len(result.Warnings))
		}
	})

	t.Run("warning consistency", func(t *testing.T) {
		warnings := []string{
			"Malformed cloze block at position 10",
			"Invalid cloze number: 0",
		}

		result := models.TypographyResult{
			ProcessedText: "Some processed text",
			ClozeCount:    1,
			WarningCount:  2,
			Warnings:      warnings,
		}

		if result.WarningCount != len(warnings) {
			t.Error("WarningCount must match length of Warnings slice")
		}

		if len(result.Warnings) != 2 {
			t.Errorf("Expected 2 warnings, got %d", len(result.Warnings))
		}
	})
}

// TestTypographyContext_ClozeDetection verifies cloze block detection contract
func TestTypographyContext_ClozeDetection(t *testing.T) {
	tests := []struct {
		name           string
		sourceText     string
		expectedCount  int
		expectedBlocks []models.ClozeDeletionBlock
	}{
		{
			name:          "no cloze blocks",
			sourceText:    "Simple text without any cloze deletions",
			expectedCount: 0,
		},
		{
			name:          "single cloze block",
			sourceText:    "The capital is {{c1::Paris}}",
			expectedCount: 1,
		},
		{
			name:          "multiple cloze blocks",
			sourceText:    "Cities: {{c1::Paris}} and {{c2::Rome}}",
			expectedCount: 2,
		},
		{
			name:          "cloze with colon in content",
			sourceText:    "Answer: {{c1::Paris : capital of France}}",
			expectedCount: 1,
		},
		{
			name:          "complex nested content",
			sourceText:    "Test {{c1::value with {{text}} inside : result}}",
			expectedCount: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// This will fail until we implement ParseClozeBlocks
			blocks, err := models.ParseClozeBlocks(tt.sourceText)
			if err != nil {
				t.Fatalf("ParseClozeBlocks() error = %v", err)
			}

			if len(blocks) != tt.expectedCount {
				t.Errorf("ParseClozeBlocks() found %d blocks, want %d", 
					len(blocks), tt.expectedCount)
			}

			// Verify blocks are properly sorted and non-overlapping
			for i := 1; i < len(blocks); i++ {
				prev := blocks[i-1]
				curr := blocks[i]

				if prev.StartPos >= curr.StartPos {
					t.Error("Blocks should be sorted by StartPos")
				}

				if prev.EndPos > curr.StartPos {
					t.Error("Blocks should not overlap")
				}
			}
		})
	}
}