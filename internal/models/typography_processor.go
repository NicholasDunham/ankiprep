package models

import (
	"regexp"
	"strings"
	"unicode"

	"golang.org/x/text/unicode/norm"
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
	result := text

	// Apply smart quotes conversion
	if tp.ConvertSmartQuotes {
		result = tp.convertSmartQuotes(result)
	}

	// Apply French typography rules
	if tp.FrenchMode {
		result = tp.applyFrenchTypography(result)
	}

	// Convert embedded newlines to HTML line breaks
	result = tp.convertLineBreaks(result)

	// Normalize Unicode (NFC normalization)
	result = norm.NFC.String(result)

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

	// Apply NNBSP before French punctuation marks: : ; ! ?
	punctuation := []string{":", ";", "!", "?"}

	for _, punct := range punctuation {
		// Pattern: word character followed by space and punctuation
		re := regexp.MustCompile(`(\w)\s*` + regexp.QuoteMeta(punct))
		text = re.ReplaceAllString(text, `$1`+nnbsp+punct)

		// Pattern: word character directly followed by punctuation (no space)
		re = regexp.MustCompile(`(\w)` + regexp.QuoteMeta(punct))
		text = re.ReplaceAllString(text, `$1`+nnbsp+punct)
	}

	// Handle French guillemets (quotation marks)
	text = tp.applyGuillemetSpacing(text)

	return text
}

// applyGuillemetSpacing applies proper spacing to French guillemets
func (tp *TypographyProcessor) applyGuillemetSpacing(text string) string {
	const nnbsp = "\u202F"

	// Opening guillemets: « followed by text
	re := regexp.MustCompile(`«\s*(\S)`)
	text = re.ReplaceAllString(text, `«`+nnbsp+`$1`)

	// Closing guillemets: text followed by »
	re = regexp.MustCompile(`(\S)\s*»`)
	text = re.ReplaceAllString(text, `$1`+nnbsp+`»`)

	return text
}

// convertLineBreaks converts embedded newlines to HTML line breaks
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
