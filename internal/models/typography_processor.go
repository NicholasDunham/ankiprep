package models

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"
)

// TypographyProcessor handles text formatting transformations
type TypographyProcessor struct {
	FrenchMode         bool // Whether French typography rules are enabled
	ConvertSmartQuotes bool // Whether to convert straight quotes to smart quotes
}

// NewTypographyProcessor creates a new TypographyProcessor instance
func NewTypographyProcessor(frenchMode, smartQuotes bool) *TypographyProcessor {
	return &TypographyProcessor{
		FrenchMode:         frenchMode,
		ConvertSmartQuotes: smartQuotes,
	}
}

// ProcessText applies all typography transformations to the input text
func (tp *TypographyProcessor) ProcessText(text string) string {
	if tp == nil {
		return text
	}

	result := text

	// Apply French typography if enabled
	if tp.FrenchMode {
		result = tp.applyFrenchTypography(result)
		result = tp.applyGuillemetSpacing(result)
	}

	// Apply smart quotes if enabled
	if tp.ConvertSmartQuotes {
		result = tp.convertSmartQuotes(result)
	}

	// FINAL STEP: Ensure all NBSP are converted to NNBSP for consistency
	// This is a final cleanup to catch any NBSP that might have been missed
	if tp.FrenchMode {
		const nbsp = "\u00A0"
		const nnbsp = "\u202F"
		result = strings.ReplaceAll(result, nbsp, nnbsp)
	}

	return result
}

// convertSmartQuotes converts straight quotes to smart quotes
func (tp *TypographyProcessor) convertSmartQuotes(text string) string {
	// Convert double quotes
	text = tp.convertDoubleQuotes(text)

	// Convert single quotes (apostrophes)
	text = tp.convertSingleQuotes(text)

	return text
}

// convertDoubleQuotes converts straight double quotes to smart quotes
func (tp *TypographyProcessor) convertDoubleQuotes(text string) string {
	// Pattern to find quoted text
	re := regexp.MustCompile(`"([^"]*)"`)

	result := re.ReplaceAllStringFunc(text, func(match string) string {
		// Remove the surrounding quotes and replace with smart quotes
		content := match[1 : len(match)-1]
		return "\u201c" + content + "\u201d" // " and "
	})

	return result
}

// convertSingleQuotes converts straight single quotes to smart apostrophes
func (tp *TypographyProcessor) convertSingleQuotes(text string) string {
	// Convert apostrophes in contractions and possessives
	re := regexp.MustCompile(`(\w)'(\w)`)
	text = re.ReplaceAllString(text, `$1\u2019$2`) // '

	// Convert single quotes around text
	re = regexp.MustCompile(`'([^']*)'`)
	text = re.ReplaceAllStringFunc(text, func(match string) string {
		content := match[1 : len(match)-1]
		return "\u2018" + content + "\u2019" // ' and '
	})

	return text
}

// applyFrenchTypography applies French typography rules (NNBSP before punctuation)
func (tp *TypographyProcessor) applyFrenchTypography(text string) string {
	// NNBSP (U+202F) - Narrow No-Break Space
	const nnbsp = "\u202F"
	// NBSP (U+00A0) - Non-Breaking Space (convert all to NNBSP)
	const nbsp = "\u00A0"

	// STEP 1: Convert ALL NBSP to NNBSP first (no exceptions!)
	text = strings.ReplaceAll(text, nbsp, nnbsp)

	// STEP 2: Protect cloze deletion syntax from French typography rules
	// Find all cloze deletions and temporarily replace them with placeholders
	clozePattern := regexp.MustCompile(`\{\{c\d+::[^}]*\}\}`)
	clozeDeletions := clozePattern.FindAllString(text, -1)

	// Replace cloze deletions with numbered placeholders
	for i, cloze := range clozeDeletions {
		placeholder := fmt.Sprintf("__CLOZE_PLACEHOLDER_%d__", i)
		text = strings.Replace(text, cloze, placeholder, 1)
	}

	// STEP 3: Apply NNBSP before French punctuation marks: : ; ! ?
	punctuation := []string{":", ";", "!", "?"}

	for _, punct := range punctuation {
		// Replace regular space + punctuation with NNBSP + punctuation
		text = strings.ReplaceAll(text, " "+punct, nnbsp+punct)

		// Handle cases where there's no space before punctuation
		// Use regex to find word character directly followed by punctuation
		pattern := regexp.MustCompile(`(\w)` + regexp.QuoteMeta(punct))
		text = pattern.ReplaceAllStringFunc(text, func(match string) string {
			// Extract the word character and punctuation
			wordChar := match[:len(match)-1]

			// Check if already has NNBSP (avoid duplicates)
			if strings.HasSuffix(wordChar, nnbsp) {
				return match // Already has NNBSP, don't modify
			}

			return wordChar + nnbsp + punct
		})
	}

	// STEP 4: Restore cloze deletions from placeholders
	for i, cloze := range clozeDeletions {
		placeholder := fmt.Sprintf("__CLOZE_PLACEHOLDER_%d__", i)
		text = strings.Replace(text, placeholder, cloze, 1)
	}

	// Handle French guillemets (quotation marks)
	text = tp.applyGuillemetSpacing(text)

	return text
}

// applyGuillemetSpacing applies proper spacing to French guillemets
func (tp *TypographyProcessor) applyGuillemetSpacing(text string) string {
	const nnbsp = "\u202F"
	const nbsp = "\u00A0"

	// STEP 1: Convert ALL remaining NBSP to NNBSP (should be none, but just in case)
	text = strings.ReplaceAll(text, nbsp, nnbsp)

	// STEP 2: Handle guillemet spacing using only NNBSP
	// Replace regular spaces with NNBSP inside guillemets
	text = strings.ReplaceAll(text, "« ", "«"+nnbsp)
	text = strings.ReplaceAll(text, " »", nnbsp+"»")

	// STEP 3: Add NNBSP where there's no space, but avoid duplicates
	// Only work with NNBSP now since all NBSP should be converted

	// Opening guillemets: « followed by non-NNBSP character (but not space)
	openPattern := regexp.MustCompile("«([^" + regexp.QuoteMeta(nnbsp) + `\s])`)
	text = openPattern.ReplaceAllString(text, "«"+nnbsp+"$1")

	// Closing guillemets: non-NNBSP character followed by » (but not space)
	closePattern := regexp.MustCompile("([^" + regexp.QuoteMeta(nnbsp) + `\s])»`)
	text = closePattern.ReplaceAllString(text, "$1"+nnbsp+"»")

	return text
} // convertLineBreaks converts embedded newlines to HTML line breaks
func (tp *TypographyProcessor) convertLineBreaks(text string) string {
	// Replace \n with <br>
	text = strings.ReplaceAll(text, "\n", "<br>")

	// Replace \r\n with <br> (Windows line endings)
	text = strings.ReplaceAll(text, "\r<br>", "<br>")

	return text
}

// PreserveHTML ensures existing HTML tags are maintained
func (tp *TypographyProcessor) PreserveHTML(text string) string {
	// This function would implement HTML preservation logic
	// For now, we assume HTML is already properly formatted
	return text
}

// IsWhitespace checks if a rune is a whitespace character
func (tp *TypographyProcessor) IsWhitespace(r rune) bool {
	return unicode.IsSpace(r)
}
