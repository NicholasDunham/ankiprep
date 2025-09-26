# Data Model: French Typography Processing

**Feature**: Improve French Typography NNBSP Handling  
**Date**: September 26, 2025

## Core Entities (Existing Architecture Enhancement)

### TypographyProcessor (Existing - Enhanced)
**Location**: `internal/models/typography_processor.go`  
**Purpose**: Handles French typography rule application with improved NNBSP detection and management

**Existing Fields**:
- `FrenchMode`: Whether French typography rules are enabled (bool)
- `ConvertSmartQuotes`: Whether to convert straight quotes to smart quotes (bool)

**Existing Methods to Enhance**:
- `ProcessText(text string) string`: Main entry point - already exists, working correctly
- `applyFrenchTypography(text string) string`: **ENHANCE** - currently adds NNBSP but doesn't detect duplicates
- `applyGuillemetSpacing(text string) string`: **ENHANCE** - currently adds NNBSP but doesn't replace regular spaces

**Enhancement Requirements**:
- `applyFrenchTypography`: Must detect existing NNBSP before adding new ones
- `applyGuillemetSpacing`: Must replace regular spaces with NNBSP, preserve existing NNBSP
- Both methods: Must handle duplicate detection to avoid `« text »` → `«  text  »`

**Current Implementation Issues**:
1. **Punctuation spacing**: Always adds NNBSP, doesn't check if NNBSP already exists
2. **Quote spacing**: Always adds NNBSP, doesn't replace regular spaces with NNBSP
3. **No duplicate detection**: Processing same text twice creates double NNBSP

**Enhancement Logic**:
- Detect existing NNBSP characters (U+202F) in text before processing
- Replace regular spaces (U+0020) with NNBSP in French typography contexts
- Preserve existing correct NNBSP spacing without duplication

### Integration Architecture (Existing - No Changes)

**Location**: `cmd/ankiprep/main.go`  
**Function**: `applyTypography(entries []*models.DataEntry, french, quotes bool)`

**Current Flow**:
```go
processor := models.NewTypographyProcessor(french, quotes)
for _, entry := range entries {
    for key, value := range entry.Values {
        entry.Values[key] = processor.ProcessText(value)
    }
}
```

**Enhancement Target**: The existing `applyFrenchTypography` and `applyGuillemetSpacing` methods called within `ProcessText`

## Enhancement Approach (No New Models)

Instead of creating new `QuoteRule`, `PunctuationRule`, and `TextProcessor` models, we enhance the existing methods:

### Method Enhancement Strategy

1. **applyFrenchTypography Enhancement**:
   - Current: Adds NNBSP before punctuation (:, ;, !, ?)
   - Enhancement: Detect existing NNBSP, only add when missing
   - Pattern: `text :` → `text :` (replace), `text :` → unchanged (preserve)

2. **applyGuillemetSpacing Enhancement**:
   - Current: Adds NNBSP inside guillemets « »
   - Enhancement: Handle both addition and replacement of spaces
   - Pattern: `«text»` → `« text »`, `« text »` → `« text »`, `« text »` → unchanged

## State Transitions (Enhanced Methods)

### Quote Processing (applyGuillemetSpacing)
1. `«text»` → `« text »` (add NNBSP when no spacing)
2. `« text »` → `« text »` (replace regular space with NNBSP)  
3. `« text »` → `« text »` (preserve existing NNBSP)

### Punctuation Processing (applyFrenchTypography)
1. `text:` → `text :` (add NNBSP when no space before punctuation)
2. `text :` → `text :` (replace regular space with NNBSP)
3. `text :` → `text :` (preserve existing NNBSP)

## Data Flow (Existing Pipeline)

```
Input Text → DataEntry.Values → applyTypography() → TypographyProcessor.ProcessText() 
    → applyFrenchTypography() [ENHANCED] 
    → applyGuillemetSpacing() [ENHANCED] 
    → Output Text
```

## Validation Constraints

- **Character Preservation**: All original text characters must be preserved (only spacing modifications)
- **Unicode Correctness**: Proper UTF-8 encoding maintained throughout processing
- **Idempotency**: Processing the same text multiple times should yield identical results
- **Rule Consistency**: All French typography rules applied uniformly across text content