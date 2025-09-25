package models_test
package models_test

import (
	"testing"

	"ankiprep/internal/models"
)

// TestClozeDeletionBlock_Contract verifies the ClozeDeletionBlock contract
func TestClozeDeletionBlock_Contract(t *testing.T) {
	t.Run("valid cloze block creation", func(t *testing.T) {
		// This will fail until we implement the model
		block := models.ClozeDeletionBlock{
			FullText: "{{c1::Paris : capital of France}}",
			Number:   1,
			Content:  "Paris : capital of France",
			Hint:     nil,
			StartPos: 0,
			EndPos:   33,
		}

		// Test validation rules
		if block.Number <= 0 {
			t.Error("Number must be positive integer")
		}

		if block.Content == "" {
			t.Error("Content cannot be empty")
		}

		if block.StartPos >= block.EndPos {
			t.Error("StartPos must be < EndPos")
		}

		if len(block.FullText) != block.EndPos-block.StartPos {
			t.Error("FullText length should match position range")
		}
	})

	t.Run("cloze block with hint", func(t *testing.T) {
		hint := "city"
		block := models.ClozeDeletionBlock{
			FullText: "{{c2::Rome::city}}",
			Number:   2,
			Content:  "Rome",
			Hint:     &hint,
			StartPos: 10,
			EndPos:   28,
		}

		if block.Hint == nil {
			t.Error("Hint should not be nil when provided")
		}

		if *block.Hint != "city" {
			t.Errorf("Hint = %q, want %q", *block.Hint, "city")
		}
	})

	t.Run("validation edge cases", func(t *testing.T) {
		tests := []struct {
			name      string
			block     models.ClozeDeletionBlock
			wantValid bool
		}{
			{
				name: "valid minimal block",
				block: models.ClozeDeletionBlock{
					FullText: "{{c1::A}}",
					Number:   1,
					Content:  "A",
					StartPos: 0,
					EndPos:   9,
				},
				wantValid: true,
			},
			{
				name: "invalid zero number",
				block: models.ClozeDeletionBlock{
					FullText: "{{c0::test}}",
					Number:   0,
					Content:  "test",
					StartPos: 0,
					EndPos:   12,
				},
				wantValid: false,
			},
			{
				name: "invalid negative number",
				block: models.ClozeDeletionBlock{
					FullText: "{{c-1::test}}",
					Number:   -1,
					Content:  "test",
					StartPos: 0,
					EndPos:   13,
				},
				wantValid: false,
			},
			{
				name: "invalid empty content",
				block: models.ClozeDeletionBlock{
					FullText: "{{c1::}}",
					Number:   1,
					Content:  "",
					StartPos: 0,
					EndPos:   8,
				},
				wantValid: false,
			},
			{
				name: "invalid position order",
				block: models.ClozeDeletionBlock{
					FullText: "{{c1::test}}",
					Number:   1,
					Content:  "test",
					StartPos: 10,
					EndPos:   5,
				},
				wantValid: false,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				// This will fail until we implement validation
				err := tt.block.Validate()
				isValid := err == nil

				if isValid != tt.wantValid {
					t.Errorf("ClozeDeletionBlock.Validate() valid = %v, want %v, error = %v", 
						isValid, tt.wantValid, err)
				}
			})
		}
	})
}

// TestClozeDeletionBlock_PatternMatching verifies pattern matching contract
func TestClozeDeletionBlock_PatternMatching(t *testing.T) {
	tests := []struct {
		name        string
		fullText    string
		shouldMatch bool
	}{
		{
			name:        "basic cloze pattern",
			fullText:    "{{c1::Paris}}",
			shouldMatch: true,
		},
		{
			name:        "cloze with hint",
			fullText:    "{{c2::Rome::city}}",
			shouldMatch: true,
		},
		{
			name:        "cloze with colon in content",
			fullText:    "{{c1::Paris : capital}}",
			shouldMatch: true,
		},
		{
			name:        "cloze with complex content",
			fullText:    "{{c3::Value with {{nested}} brackets : result}}",
			shouldMatch: true,
		},
		{
			name:        "invalid - no closing brackets",
			fullText:    "{{c1::Paris",
			shouldMatch: false,
		},
		{
			name:        "invalid - no opening brackets", 
			fullText:    "c1::Paris}}",
			shouldMatch: false,
		},
		{
			name:        "invalid - wrong format",
			fullText:    "{c1::Paris}",
			shouldMatch: false,
		},
		{
			name:        "invalid - no colon separators",
			fullText:    "{{c1Paris}}",
			shouldMatch: false,
		},
		{
			name:        "invalid - no number",
			fullText:    "{{c::Paris}}",
			shouldMatch: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// This will fail until we implement pattern matching
			matches := models.IsValidClozeDeletionPattern(tt.fullText)

			if matches != tt.shouldMatch {
				t.Errorf("IsValidClozeDeletionPattern(%q) = %v, want %v", 
					tt.fullText, matches, tt.shouldMatch)
			}
		})
	}
}