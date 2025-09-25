package services

import (
	"fmt"
	"log"
	"strings"
	"time"

	"ankiprep/internal/models"
)

// ProcessFrenchTextRequest represents a request to process French text with typography rules.
type ProcessFrenchTextRequest struct {
	// Source text to process (required)
	Text string `json:"text" validate:"required"`
	// Enable French typography processing
	French bool `json:"french"`
	// Enable smart quotes processing
	SmartQuotes bool `json:"smartQuotes"`
	// Logger for warnings (optional, defaults to stderr)
	Logger models.Logger `json:"-"`
}

// ProcessFrenchTextResponse represents the response from processing French text.
type ProcessFrenchTextResponse struct {
	// Processed text with typography rules applied
	ProcessedText string `json:"processedText"`
	// Number of cloze blocks found and processed
	ClozeCount int `json:"clozeCount"`
	// Number of warnings logged during processing
	WarningCount int `json:"warningCount"`
	// Warning messages for debugging
	Warnings []string `json:"warnings,omitempty"`
	// Processing metadata
	ProcessingTime time.Duration `json:"processingTimeMs"`
}

// Validate checks that the ProcessFrenchTextRequest is valid.
func (req *ProcessFrenchTextRequest) Validate() error {
	// Text must not be empty (and not just whitespace)
	if strings.TrimSpace(req.Text) == "" {
		return fmt.Errorf("text is required and cannot be empty")
	}

	// Text must not exceed 1MB
	if len(req.Text) > 1048576 {
		return fmt.Errorf("text size exceeds 1MB limit (1048576 characters), got %d characters", len(req.Text))
	}

	// At least one processing option must be enabled
	if !req.French && !req.SmartQuotes {
		return fmt.Errorf("at least one of French or SmartQuotes processing must be enabled")
	}

	return nil
}

// TypographyService provides French typography processing with cloze deletion awareness.
type TypographyService struct {
	// Default logger for warnings
	logger models.Logger
	// Legacy processor for backward compatibility
	processor *models.TypographyProcessor
	// French mode setting (for backward compatibility)
	frenchMode bool
	// Smart quotes setting (for backward compatibility)
	smartQuotes bool
}

// NewTypographyService creates a new TypographyService instance with cloze awareness.
func NewTypographyService() *TypographyService {
	return &TypographyService{
		logger:      log.Default(),
		frenchMode:  true,  // Default to French mode for new interface
		smartQuotes: false, // Default to no smart quotes
	}
}

// NewTypographyServiceLegacy creates a TypographyService for backward compatibility.
func NewTypographyServiceLegacy(frenchMode bool, smartQuotes bool) *TypographyService {
	var processor *models.TypographyProcessor
	// Only create processor if we can (models.NewTypographyProcessor might not exist yet)
	// This is for backward compatibility with existing code

	return &TypographyService{
		logger:      log.Default(),
		processor:   processor,
		frenchMode:  frenchMode,
		smartQuotes: smartQuotes,
	}
}

// ProcessFrenchText processes text with French typography rules while preserving cloze deletion blocks.
func (ts *TypographyService) ProcessFrenchText(req ProcessFrenchTextRequest) (*ProcessFrenchTextResponse, error) {
	startTime := time.Now()

	// Validate request
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("invalid request: %w", err)
	}

	// Use request logger or service default
	logger := req.Logger
	if logger == nil {
		logger = ts.logger
	}

	// Create typography context with cloze detection
	context, err := models.NewTypographyContext(req.Text, req.French, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to create typography context: %w", err)
	}

	// Process the text while preserving cloze blocks
	result, err := ts.processWithClozeAwareness(context, req)
	if err != nil {
		return nil, fmt.Errorf("processing failed: %w", err)
	}

	// Create response
	response := &ProcessFrenchTextResponse{
		ProcessedText:  result.ProcessedText,
		ClozeCount:     result.ClozeCount,
		WarningCount:   result.WarningCount,
		Warnings:       result.Warnings,
		ProcessingTime: time.Since(startTime),
	}

	return response, nil
}

// processWithClozeAwareness applies typography rules while preserving cloze deletion blocks.
func (ts *TypographyService) processWithClozeAwareness(context *models.TypographyContext, req ProcessFrenchTextRequest) (*models.TypographyResult, error) {
	sourceText := context.SourceText
	clozeBlocks := context.ClozeBlocks

	if len(clozeBlocks) == 0 {
		// No cloze blocks - apply typography rules normally
		processed := ts.applyTypographyRules(sourceText, req.French, req.SmartQuotes)
		return models.NewTypographyResult(processed, 0, nil), nil
	}

	// Process text in segments:
	// 1. Apply typography to text outside cloze blocks
	// 2. Preserve text inside cloze blocks unchanged
	var result strings.Builder
	warnings := []string{}
	lastEnd := 0

	for _, block := range clozeBlocks {
		// Process text before this cloze block
		if block.StartPos > lastEnd {
			beforeText := sourceText[lastEnd:block.StartPos]
			processedBefore := ts.applyTypographyRules(beforeText, req.French, req.SmartQuotes)
			result.WriteString(processedBefore)
		}

		// Append cloze block unchanged
		result.WriteString(block.FullText)
		lastEnd = block.EndPos
	}

	// Process remaining text after last cloze block
	if lastEnd < len(sourceText) {
		afterText := sourceText[lastEnd:]
		processedAfter := ts.applyTypographyRules(afterText, req.French, req.SmartQuotes)
		result.WriteString(processedAfter)
	}

	return models.NewTypographyResult(result.String(), len(clozeBlocks), warnings), nil
}

// applyTypographyRules applies French typography and smart quotes rules to text.
func (ts *TypographyService) applyTypographyRules(text string, french, smartQuotes bool) string {
	if text == "" {
		return text
	}

	result := text

	// Apply French typography rules
	if french {
		result = ts.applyFrenchTypography(result)
	}

	// Apply smart quotes
	if smartQuotes {
		result = ts.applySmartQuotes(result)
	}

	return result
}

// applyFrenchTypography applies French typography rules (non-breaking spaces before punctuation).
func (ts *TypographyService) applyFrenchTypography(text string) string {
	if text == "" {
		return text
	}

	// Replace space + colon with non-breaking space + colon
	// This is the core French typography rule for colons
	result := strings.ReplaceAll(text, " :", "\u00A0:")

	// Apply other French typography rules as needed:
	// Non-breaking spaces before semicolons, question marks, exclamation marks
	result = strings.ReplaceAll(result, " ;", "\u00A0;")
	result = strings.ReplaceAll(result, " ?", "\u00A0?")
	result = strings.ReplaceAll(result, " !", "\u00A0!")

	// Non-breaking spaces with guillemets (French quotes)
	result = strings.ReplaceAll(result, "« ", "«\u00A0")
	result = strings.ReplaceAll(result, " »", "\u00A0»")

	return result
}

// applySmartQuotes converts straight quotes to curly quotes.
func (ts *TypographyService) applySmartQuotes(text string) string {
	if text == "" {
		return text
	}

	result := text

	// Simple smart quote conversion
	// This is a basic implementation - a production version might be more sophisticated
	result = strings.ReplaceAll(result, `"`, `"`)
	result = strings.ReplaceAll(result, `"`, `"`)
	result = strings.ReplaceAll(result, `'`, `'`)
	result = strings.ReplaceAll(result, `'`, `'`)

	return result
}

// Legacy methods for backward compatibility with existing code

// ProcessEntries applies typography formatting to all text content in data entries (legacy)
func (s *TypographyService) ProcessEntries(entries []*models.DataEntry) []*models.DataEntry {
	if s.processor != nil {
		// Use legacy processor if available
		processedEntries := make([]*models.DataEntry, len(entries))
		for i, entry := range entries {
			processedEntry := s.ProcessEntry(entry)
			processedEntries[i] = processedEntry
		}
		return processedEntries
	}

	// Fallback to new cloze-aware processing
	processedEntries := make([]*models.DataEntry, len(entries))
	for i, entry := range entries {
		processedEntry := s.processEntryWithClozeAwareness(entry)
		processedEntries[i] = processedEntry
	}
	return processedEntries
}

// ProcessEntry applies typography formatting to a single data entry (legacy)
func (s *TypographyService) ProcessEntry(entry *models.DataEntry) *models.DataEntry {
	if entry == nil {
		return nil
	}

	if s.processor != nil {
		// Use legacy processor if available
		processedValues := make(map[string]string)
		for column, value := range entry.Values {
			processedValues[column] = s.processor.ProcessText(value)
		}
		return models.NewDataEntry(processedValues, entry.Source, entry.LineNumber)
	}

	// Fallback to new cloze-aware processing
	return s.processEntryWithClozeAwareness(entry)
}

// ProcessText applies typography formatting to a single text string (legacy)
func (s *TypographyService) ProcessText(text string) string {
	if s.processor != nil {
		return s.processor.ProcessText(text)
	}

	// Fallback to new cloze-aware processing
	request := ProcessFrenchTextRequest{
		Text:        text,
		French:      s.frenchMode,
		SmartQuotes: s.smartQuotes,
	}

	response, err := s.ProcessFrenchText(request)
	if err != nil {
		// On error, return original text
		return text
	}

	return response.ProcessedText
}

// processEntryWithClozeAwareness processes a data entry using cloze-aware typography.
func (s *TypographyService) processEntryWithClozeAwareness(entry *models.DataEntry) *models.DataEntry {
	if entry == nil {
		return nil
	}

	processedValues := make(map[string]string)

	for column, value := range entry.Values {
		// Use new cloze-aware processing
		request := ProcessFrenchTextRequest{
			Text:        value,
			French:      s.frenchMode,
			SmartQuotes: s.smartQuotes,
		}

		response, err := s.ProcessFrenchText(request)
		if err != nil {
			// On error, keep original value and log warning
			s.logger.Printf("Warning: failed to process text in column %s: %v", column, err)
			processedValues[column] = value
		} else {
			processedValues[column] = response.ProcessedText
		}
	}

	return models.NewDataEntry(processedValues, entry.Source, entry.LineNumber)
}

// Legacy compatibility methods

// SetFrenchMode enables or disables French typography rules
func (s *TypographyService) SetFrenchMode(enabled bool) {
	s.frenchMode = enabled
	if s.processor != nil {
		s.processor.FrenchMode = enabled
	}
}

// SetSmartQuotes enables or disables smart quote conversion
func (s *TypographyService) SetSmartQuotes(enabled bool) {
	s.smartQuotes = enabled
	if s.processor != nil {
		s.processor.ConvertSmartQuotes = enabled
	}
}

// IsFrenchModeEnabled returns whether French typography mode is enabled
func (s *TypographyService) IsFrenchModeEnabled() bool {
	if s.processor != nil {
		return s.processor.FrenchMode
	}
	return s.frenchMode
}

// IsSmartQuotesEnabled returns whether smart quote conversion is enabled
func (s *TypographyService) IsSmartQuotesEnabled() bool {
	if s.processor != nil {
		return s.processor.ConvertSmartQuotes
	}
	return s.smartQuotes
}

// GetProcessor returns the underlying typography processor for direct access
func (s *TypographyService) GetProcessor() *models.TypographyProcessor {
	return s.processor
}
