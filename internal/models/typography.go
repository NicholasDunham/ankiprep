package models

import (
	"fmt"
	"log"
)

// Logger interface for logging warnings and errors during processing.
type Logger interface {
	Printf(format string, v ...interface{})
	Print(v ...interface{})
}

// TypographyContext provides enhanced context for French typography processing with cloze awareness.
type TypographyContext struct {
	// Source text being processed
	SourceText string
	// Detected cloze blocks within source text
	ClozeBlocks []ClozeDeletionBlock
	// Whether French typography rules are enabled
	FrenchEnabled bool
	// Logging interface for warnings/errors
	Logger Logger
}

// Validate checks that the TypographyContext satisfies all validation rules.
func (tc *TypographyContext) Validate() error {
	// Check that all cloze blocks have valid positions for the source text
	sourceLen := len(tc.SourceText)
	for i, block := range tc.ClozeBlocks {
		if block.StartPos < 0 || block.EndPos > sourceLen {
			return fmt.Errorf("cloze block %d has invalid positions [%d:%d] for text of length %d",
				i, block.StartPos, block.EndPos, sourceLen)
		}

		// Validate the block itself
		if err := block.Validate(); err != nil {
			return fmt.Errorf("cloze block %d validation failed: %w", i, err)
		}
	}

	// Check that blocks are sorted by StartPos
	for i := 1; i < len(tc.ClozeBlocks); i++ {
		if tc.ClozeBlocks[i-1].StartPos >= tc.ClozeBlocks[i].StartPos {
			return fmt.Errorf("cloze blocks must be sorted by StartPos, but block %d (%d) >= block %d (%d)",
				i-1, tc.ClozeBlocks[i-1].StartPos, i, tc.ClozeBlocks[i].StartPos)
		}
	}

	// Check that blocks do not overlap
	for i := 1; i < len(tc.ClozeBlocks); i++ {
		prev := tc.ClozeBlocks[i-1]
		curr := tc.ClozeBlocks[i]
		if prev.EndPos > curr.StartPos {
			return fmt.Errorf("overlapping cloze blocks: block %d ends at %d, block %d starts at %d",
				i-1, prev.EndPos, i, curr.StartPos)
		}
	}

	return nil
}

// NewTypographyContext creates a new TypographyContext with cloze blocks detected from the source text.
func NewTypographyContext(sourceText string, frenchEnabled bool, logger Logger) (*TypographyContext, error) {
	if logger == nil {
		logger = log.Default()
	}

	// Parse cloze blocks from source text
	blocks, err := ParseClozeBlocks(sourceText)
	if err != nil {
		return nil, fmt.Errorf("failed to parse cloze blocks: %w", err)
	}

	context := &TypographyContext{
		SourceText:    sourceText,
		ClozeBlocks:   blocks,
		FrenchEnabled: frenchEnabled,
		Logger:        logger,
	}

	// Validate the context
	if err := context.Validate(); err != nil {
		return nil, fmt.Errorf("invalid typography context: %w", err)
	}

	return context, nil
}

// TypographyResult represents the result of typography processing operation.
type TypographyResult struct {
	// Processed text with typography rules applied
	ProcessedText string
	// Number of cloze blocks successfully processed
	ClozeCount int
	// Number of warnings logged (malformed blocks)
	WarningCount int
	// Processing errors (non-fatal)
	Warnings []string
}

// NewTypographyResult creates a new TypographyResult with validation.
func NewTypographyResult(processedText string, clozeCount int, warnings []string) *TypographyResult {
	if warnings == nil {
		warnings = []string{}
	}

	return &TypographyResult{
		ProcessedText: processedText,
		ClozeCount:    clozeCount,
		WarningCount:  len(warnings),
		Warnings:      warnings,
	}
}

// AddWarning adds a warning to the result and updates the warning count.
func (tr *TypographyResult) AddWarning(warning string) {
	tr.Warnings = append(tr.Warnings, warning)
	tr.WarningCount = len(tr.Warnings)
}
