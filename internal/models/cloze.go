package models

import (
	"fmt"
	"regexp"
	"strings"
)

// ClozeDeletionBlock represents a parsed Anki cloze deletion block within text content.
type ClozeDeletionBlock struct {
	// Full matched text including {{}} brackets
	FullText string
	// Cloze number (e.g., 1 from {{c1::...}})
	Number int
	// Target content (text between first :: and second :: or closing }})
	Content string
	// Optional hint text (text after second :: if present)
	Hint *string
	// Start position in original text
	StartPos int
	// End position in original text
	EndPos int
}

// clozeStartPattern matches the start of cloze deletion patterns
var clozeStartPattern = regexp.MustCompile(`\{\{c(\d+)::`)

// Validate checks that the ClozeDeletionBlock satisfies all validation rules.
func (c *ClozeDeletionBlock) Validate() error {
	// Number must be positive integer (1-99 typical range)
	if c.Number <= 0 {
		return fmt.Errorf("cloze number must be positive, got %d", c.Number)
	}

	// Content cannot be empty
	if strings.TrimSpace(c.Content) == "" {
		return fmt.Errorf("cloze content cannot be empty")
	}

	// StartPos must be < EndPos
	if c.StartPos >= c.EndPos {
		return fmt.Errorf("start position (%d) must be less than end position (%d)",
			c.StartPos, c.EndPos)
	}

	// FullText must match pattern
	if !IsValidClozeDeletionPattern(c.FullText) {
		return fmt.Errorf("full text does not match valid cloze deletion pattern: %q", c.FullText)
	}

	return nil
}

// IsValidClozeDeletionPattern checks if a string matches the cloze deletion pattern.
func IsValidClozeDeletionPattern(text string) bool {
	// Simple check: starts with {{c followed by digits and ::
	return clozeStartPattern.MatchString(text) && strings.HasSuffix(text, "}}")
}

// ParseClozeBlocks extracts all cloze deletion blocks from the given text using brace counting.
// Returns blocks sorted by StartPos with no overlapping positions.
func ParseClozeBlocks(text string) ([]ClozeDeletionBlock, error) {
	var blocks []ClozeDeletionBlock

	// Find all potential cloze starts
	starts := clozeStartPattern.FindAllStringSubmatchIndex(text, -1)

	for _, match := range starts {
		startPos := match[0]
		numberStart := match[2]
		numberEnd := match[3]
		contentStart := match[1] // This is the end of the full match "{{c1::"

		// Parse the cloze number
		numberStr := text[numberStart:numberEnd]
		var number int
		if n, err := fmt.Sscanf(numberStr, "%d", &number); err != nil || n != 1 || number <= 0 {
			continue // Skip invalid numbers
		}

		// Find the matching }} by counting braces
		braceCount := 2 // Start with 2 open braces from {{
		endPos := -1

		for i := contentStart; i < len(text)-1; i++ {
			if text[i] == '{' && text[i+1] == '{' {
				braceCount += 2
				i++ // Skip the next {
			} else if text[i] == '}' && text[i+1] == '}' {
				braceCount -= 2
				if braceCount == 0 {
					endPos = i + 2 // Include both }}
					break
				}
				i++ // Skip the next }
			}
		}

		if endPos == -1 {
			continue // No matching closing braces found
		}

		// Extract the full text and content
		fullText := text[startPos:endPos]
		content := text[contentStart : endPos-2] // Exclude the final }}

		// Parse content and hint
		var hint *string
		if strings.Contains(content, "::") {
			parts := strings.SplitN(content, "::", 2)
			if len(parts) == 2 {
				content = parts[0]
				hintValue := parts[1]
				if strings.TrimSpace(hintValue) != "" {
					hint = &hintValue
				}
			}
		}

		block := ClozeDeletionBlock{
			FullText: fullText,
			Number:   number,
			Content:  content,
			Hint:     hint,
			StartPos: startPos,
			EndPos:   endPos,
		}

		// Validate the block
		if err := block.Validate(); err != nil {
			// Log warning but continue processing
			continue
		}

		blocks = append(blocks, block)
	}

	// Verify no overlaps (blocks are already sorted by StartPos due to search order)
	for i := 1; i < len(blocks); i++ {
		if blocks[i-1].EndPos > blocks[i].StartPos {
			return nil, fmt.Errorf("overlapping cloze blocks detected at positions %d and %d",
				blocks[i-1].StartPos, blocks[i].StartPos)
		}
	}

	return blocks, nil
}
