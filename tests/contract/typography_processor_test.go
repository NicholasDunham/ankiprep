package contract

import (
	"testing"

	"ankiprep/internal/models"
)

// TestTypographyProcessor_FrenchProcessing tests the contract for French typography processing
func TestTypographyProcessor_FrenchProcessing(t *testing.T) {
	processor := models.NewTypographyProcessor(true, false) // French mode enabled

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		// Quote Processing Tests - Contract from typography-interface.md
		{
			name:     "Quote spacing - no existing spaces",
			input:    "«bonjour»",
			expected: "«\u202Fbonjour\u202F»", // NNBSP (U+202F) expected
		},
		{
			name:     "Quote spacing - replace regular spaces",
			input:    "« bonjour »",           // Regular space (U+0020)
			expected: "«\u202Fbonjour\u202F»", // NNBSP (U+202F) expected
		},
		{
			name:     "Quote spacing - preserve existing NNBSP",
			input:    "«\u202Fbonjour\u202F»", // NNBSP (U+202F) already present
			expected: "«\u202Fbonjour\u202F»", // Should remain unchanged
		},

		// Punctuation Processing Tests - Contract from typography-interface.md
		{
			name:     "Punctuation spacing - add NNBSP before colon",
			input:    "Bonjour:",
			expected: "Bonjour\u202F:", // NNBSP (U+202F) expected
		},
		{
			name:     "Punctuation spacing - replace regular space",
			input:    "Bonjour :",
			expected: "Bonjour\u202F:", // NNBSP (U+202F) expected
		},
		{
			name:     "Punctuation spacing - preserve existing NNBSP",
			input:    "Bonjour\u202F:",
			expected: "Bonjour\u202F:", // Should remain unchanged
		},
		{
			name:     "Multiple punctuation marks",
			input:    "Comment? Bien! Merci; Au revoir:",
			expected: "Comment\u202F? Bien\u202F! Merci\u202F; Au revoir\u202F:",
		},

		// Complex scenarios
		{
			name:     "Mixed quotes and punctuation",
			input:    "Il dit «Bonjour: comment allez-vous?»",
			expected: "Il dit «\u202FBonjour\u202F: comment allez-vous\u202F?\u202F»",
		},
		{
			name:     "Already processed text - idempotency test",
			input:    "«\u202FBonjour\u202F: comment allez-vous\u202F?\u202F»",
			expected: "«\u202FBonjour\u202F: comment allez-vous\u202F?\u202F»", // Should remain identical
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := processor.ProcessText(tt.input)
			if result != tt.expected {
				t.Errorf("ProcessText() = %q, expected %q", result, tt.expected)
				// Debug: show actual vs expected character codes
				t.Logf("Actual chars: %+q", []rune(result))
				t.Logf("Expected chars: %+q", []rune(tt.expected))
			}
		})
	}
}

// TestTypographyProcessor_ErrorHandling tests error conditions
func TestTypographyProcessor_ErrorHandling(t *testing.T) {
	processor := models.NewTypographyProcessor(true, false)

	// Test empty input handling
	result := processor.ProcessText("")
	if result != "" {
		t.Errorf("ProcessText(\"\") = %q, expected \"\"", result)
	}

	// Test that processor handles invalid UTF-8 gracefully (Go strings are UTF-8 by default)
	// This is more of a documentation test since Go handles this automatically
}

// TestTypographyProcessor_NonFrenchMode tests that French processing is disabled in non-French mode
func TestTypographyProcessor_NonFrenchMode(t *testing.T) {
	processor := models.NewTypographyProcessor(false, false) // French mode disabled

	input := "«bonjour» Comment:"
	result := processor.ProcessText(input)

	// Should NOT apply French typography rules
	if result != input {
		t.Errorf("ProcessText() with French=false should not modify text, got %q", result)
	}
}
