# Data Model: Cloze Deletion Colon Exception

**Date**: 2025-09-24  
**Feature**: French typography processing with cloze deletion support

## Core Entities

### ClozeDeletionBlock
Represents a parsed Anki cloze deletion block within text content.

```go
type ClozeDeletionBlock struct {
    // Full matched text including {{}} brackets
    FullText    string
    // Cloze number (e.g., 1 from {{c1::...}})
    Number      int
    // Target content (text between first :: and second :: or closing }})
    Content     string
    // Optional hint text (text after second :: if present)
    Hint        *string
    // Start position in original text
    StartPos    int
    // End position in original text  
    EndPos      int
}
```

**Validation Rules**:
- Number must be positive integer (1-99 typical range)
- Content cannot be empty
- FullText must match pattern `{{c\d+::[^}]*}}`
- StartPos must be < EndPos
- Positions must be valid indices in source text

**State Transitions**: Immutable value object (no state changes)

### TypographyContext
Enhanced context for French typography processing with cloze awareness.

```go
type TypographyContext struct {
    // Source text being processed
    SourceText      string
    // Detected cloze blocks within source text
    ClozeBlocks     []ClozeDeletionBlock
    // Whether French typography rules are enabled
    FrenchEnabled   bool
    // Logging interface for warnings/errors
    Logger          Logger
}
```

**Validation Rules**:
- ClozeBlocks positions must not overlap
- ClozeBlocks must be sorted by StartPos
- All ClozeBlocks StartPos/EndPos must be valid for SourceText length

### TypographyResult
Result of typography processing operation.

```go
type TypographyResult struct {
    // Processed text with typography rules applied
    ProcessedText   string
    // Number of cloze blocks successfully processed
    ClozeCount      int
    // Number of warnings logged (malformed blocks)
    WarningCount    int
    // Processing errors (non-fatal)
    Warnings        []string
}
```

## Entity Relationships

```
TypographyContext
├── contains multiple ClozeDeletionBlock
└── produces TypographyResult

ClozeDeletionBlock
├── parsed from TypographyContext.SourceText  
└── influences typography rule application
```

## Data Volume Assumptions

- **Typical CSV field**: 50-500 characters
- **Cloze blocks per field**: 1-5 blocks typically, up to 20 maximum
- **CSV file size**: Up to 10MB (hundreds to thousands of rows)
- **Memory usage**: < 100MB for largest expected files
- **Processing time**: Sub-second for typical files, < 10 seconds for largest files

## Persistence

**Note**: This feature uses in-memory processing only. No persistent storage of entities required.

- Input: CSV files read from filesystem
- Output: Modified CSV files written to filesystem  
- Intermediate data: Held in memory during processing pipeline
- Logging: Written to stderr/log files as configured

## Integration Patterns

### Input Processing
```
CSV Reader → Field Text → ClozeDeletionBlock Detection → TypographyContext Creation
```

### Typography Processing  
```
TypographyContext → Cloze-Aware Typography Service → TypographyResult
```

### Output Generation
```
TypographyResult → Updated Field Text → CSV Writer
```

## Error Handling

### Malformed Cloze Blocks
- **Detection**: Regex pattern mismatch, incomplete brackets, invalid numbers
- **Response**: Log warning, exclude from ClozeBlocks collection, process field as regular text
- **Recovery**: Continue processing other fields/blocks

### Processing Errors
- **Typography failures**: Log error, return original text unchanged
- **Memory constraints**: Fail fast with clear error message
- **File I/O errors**: Propagate to CLI layer for user-friendly display

## Performance Considerations

### Parsing Optimization
- Use compiled regex patterns (initialize once, reuse)
- Avoid string copying where possible (use slices/indices)
- Process blocks in single pass through text

### Memory Management
- Release large strings promptly after processing
- Use streaming for large CSV files when possible
- Avoid accumulating full file content in memory simultaneously