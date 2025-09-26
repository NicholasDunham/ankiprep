# Typography Processing Contract

**Service**: French Typography Processor  
**Version**: 1.0  
**Date**: September 26, 2025

## Interface Definition

### ProcessFrenchTypography
**Purpose**: Apply French typography rules to text content with intelligent NNBSP handling

```go
type TypographyProcessor interface {
    ProcessFrenchTypography(text string) (string, error)
}
```

**Input**:
- `text` (string): UTF-8 encoded text content requiring French typography processing

**Output**:
- `processedText` (string): Text with French typography rules applied
- `error`: Processing error if any (nil on success)

**Processing Rules**:
1. **Quote Processing**: Apply NNBSP spacing to angled quotes (« »)
   - `«text»` → `«{NNBSP}text{NNBSP}»`
   - `« text »` → `«{NNBSP}text{NNBSP}»` (replace regular space)
   - `«{NNBSP}text{NNBSP}»` → unchanged (preserve existing NNBSP)

2. **Punctuation Processing**: Apply NNBSP spacing before double punctuation (:, ;, !, ?)
   - `text:` → `text{NNBSP}:`
   - `text :` → `text{NNBSP}:` (replace regular space)
   - `text{NNBSP}:` → unchanged (preserve existing NNBSP)

**Error Conditions**:
- Invalid UTF-8 input → `ErrInvalidUTF8`
- Empty input text → `ErrEmptyInput` 
- Processing failure → `ErrProcessingFailed`

## Contract Tests

### Test: Quote Spacing - No Existing Spaces
```go
func TestQuoteSpacing_NoSpaces(t *testing.T) {
    input := "«bonjour»"
    expected := "« bonjour »" // Using NNBSP (U+202F)
    
    result, err := processor.ProcessFrenchTypography(input)
    
    assert.NoError(t, err)
    assert.Equal(t, expected, result)
}
```

### Test: Quote Spacing - Replace Regular Spaces
```go
func TestQuoteSpacing_ReplaceRegularSpaces(t *testing.T) {
    input := "« bonjour »" // Using regular space (U+0020)
    expected := "« bonjour »" // Using NNBSP (U+202F)
    
    result, err := processor.ProcessFrenchTypography(input)
    
    assert.NoError(t, err)
    assert.Equal(t, expected, result)
}
```

### Test: Quote Spacing - Preserve Existing NNBSP
```go
func TestQuoteSpacing_PreserveNNBSP(t *testing.T) {
    input := "« bonjour »" // Already using NNBSP (U+202F)
    expected := "« bonjour »" // Should remain unchanged
    
    result, err := processor.ProcessFrenchTypography(input)
    
    assert.NoError(t, err)
    assert.Equal(t, expected, result)
}
```

### Test: Punctuation Spacing - No Existing Space
```go
func TestPunctuationSpacing_NoSpace(t *testing.T) {
    input := "Bonjour:"
    expected := "Bonjour :" // Using NNBSP (U+202F)
    
    result, err := processor.ProcessFrenchTypography(input)
    
    assert.NoError(t, err)
    assert.Equal(t, expected, result)
}
```

### Test: Punctuation Spacing - Replace Regular Space
```go
func TestPunctuationSpacing_ReplaceRegularSpace(t *testing.T) {
    input := "Bonjour :" // Using regular space (U+0020)
    expected := "Bonjour :" // Using NNBSP (U+202F)
    
    result, err := processor.ProcessFrenchTypography(input)
    
    assert.NoError(t, err)
    assert.Equal(t, expected, result)
}
```

### Test: Punctuation Spacing - Preserve Existing NNBSP
```go
func TestPunctuationSpacing_PreserveNNBSP(t *testing.T) {
    input := "Bonjour :" // Already using NNBSP (U+202F)
    expected := "Bonjour :" // Should remain unchanged
    
    result, err := processor.ProcessFrenchTypography(input)
    
    assert.NoError(t, err)
    assert.Equal(t, expected, result)
}
```

### Test: All Punctuation Types
```go
func TestAllPunctuationTypes(t *testing.T) {
    input := "Bonjour: Comment allez-vous; Très bien! Vous?"
    expected := "Bonjour : Comment allez-vous ; Très bien ! Vous ?" // Using NNBSP
    
    result, err := processor.ProcessFrenchTypography(input)
    
    assert.NoError(t, err)
    assert.Equal(t, expected, result)
}
```

### Test: Error Handling - Invalid UTF-8
```go
func TestErrorHandling_InvalidUTF8(t *testing.T) {
    input := "\xff\xfe" // Invalid UTF-8 sequence
    
    _, err := processor.ProcessFrenchTypography(input)
    
    assert.Error(t, err)
    assert.Equal(t, ErrInvalidUTF8, err)
}
```

### Test: Error Handling - Empty Input
```go
func TestErrorHandling_EmptyInput(t *testing.T) {
    input := ""
    
    _, err := processor.ProcessFrenchTypography(input)
    
    assert.Error(t, err)
    assert.Equal(t, ErrEmptyInput, err)
}
```

## Character Reference

- **Regular Space**: U+0020 ` ` (input from user)
- **Narrow Non-Breaking Space (NNBSP)**: U+202F ` ` (French typography standard)
- **Opening Angle Quote**: U+00AB `«`
- **Closing Angle Quote**: U+00BB `»`
- **Double Punctuation**: `:` (U+003A), `;` (U+003B), `!` (U+0021), `?` (U+003F)