# Typography Service Contract

**Service**: `TypographyService`  
**Method**: `ProcessFrenchText`  
**Description**: Process text with French typography rules, handling Anki cloze deletion blocks specially

## Input Contract

```go
type ProcessFrenchTextRequest struct {
    // Source text to process (required)
    Text        string      `json:"text" validate:"required"`
    // Enable French typography processing
    French      bool        `json:"french"`
    // Enable smart quotes processing  
    SmartQuotes bool        `json:"smartQuotes"`
    // Logger for warnings (optional, defaults to stderr)
    Logger      Logger      `json:"-"`
}
```

**Validation Rules**:
- `Text` must not be empty
- At least one of `French` or `SmartQuotes` must be true
- `Text` length must be ≤ 1MB (1,048,576 characters)

## Output Contract

```go
type ProcessFrenchTextResponse struct {
    // Processed text with typography rules applied
    ProcessedText   string          `json:"processedText"`
    // Number of cloze blocks found and processed
    ClozeCount      int             `json:"clozeCount"`
    // Number of warnings logged during processing
    WarningCount    int             `json:"warningCount"`
    // Warning messages for debugging
    Warnings        []string        `json:"warnings,omitempty"`
    // Processing metadata
    ProcessingTime  time.Duration   `json:"processingTimeMs"`
}
```

**Guarantees**:
- `ProcessedText` will never be empty if `Text` was non-empty
- `ClozeCount` will be ≥ 0
- `WarningCount` equals `len(Warnings)`
- If no cloze blocks found, behavior matches original French typography processing
- Processing is idempotent (applying twice produces same result)

## Behavior Specification

### Cloze Detection Rules
1. Pattern: `{{c\d+::[^}]*}}`
2. Nested blocks not supported (inner `{{}}` treated as literal text)
3. Malformed patterns logged as warnings, processed as regular text
4. Empty cloze content (`{{c1::}}`) is valid

### Colon Processing Rules
1. **Cloze syntax colons** (`::` within `{{c#::...}}`): NO NNBSP added
2. **Content colons** (within cloze content or hint): Apply French NNBSP rules
3. **External colons** (outside any `{{}}`): Apply French NNBSP rules

### Error Conditions

| Condition | Response | HTTP Status | Recovery |
|-----------|----------|-------------|----------|
| Empty text | Return empty ProcessedText | 200 | Continue |
| Text too large | Return error | 400 | User must reduce input |
| Malformed cloze | Log warning, process as regular text | 200 | Continue |
| Processing failure | Return original text | 500 | Log error, fail gracefully |

## Examples

### Basic Cloze Processing
```json
// Input
{
    "text": "Je vais {{c1::essayer}} cette recette : c'est délicieux !",
    "french": true,
    "smartQuotes": false
}

// Output  
{
    "processedText": "Je vais {{c1::essayer}} cette recette : c'est délicieux !",
    "clozeCount": 1,
    "warningCount": 0,
    "warnings": [],
    "processingTimeMs": 2
}
```

Note: The cloze syntax colons (`::`) are preserved without NNBSP, while the content colon (`:`) gets NNBSP treatment.

### Multiple Cloze Blocks
```json
// Input
{
    "text": "{{c1::Bonjour}} : comment {{c2::allez-vous}} ?",
    "french": true,
    "smartQuotes": true  
}

// Output
{
    "processedText": "{{c1::Bonjour}} : comment {{c2::allez-vous}} ?",
    "clozeCount": 2,
    "warningCount": 0,
    "warnings": [],
    "processingTimeMs": 3
}
```

### Malformed Cloze Block
```json  
// Input
{
    "text": "{{c1::incomplete block missing close",
    "french": true,
    "smartQuotes": false
}

// Output
{
    "processedText": "{{c1::incomplete block missing close",
    "clozeCount": 0,
    "warningCount": 1, 
    "warnings": ["Malformed cloze deletion block at position 0: missing closing brackets"],
    "processingTimeMs": 1
}
```

### Nested Content with Colons
```json
// Input  
{
    "text": "{{c1::phrase with « quoted text : example »}}",
    "french": true,
    "smartQuotes": true
}

// Output
{
    "processedText": "{{c1::phrase with « quoted text : example »}}",
    "clozeCount": 1,
    "warningCount": 0,
    "warnings": [],
    "processingTimeMs": 2
}
```

Note: The cloze syntax colons (`::` after `c1`) get no NNBSP, but the content colon (`:` in the quoted text) gets NNBSP per French rules.

## Performance Contract

- **Response time**: ≤ 100ms for texts up to 10KB
- **Memory usage**: ≤ 10MB peak for largest supported inputs  
- **Throughput**: ≥ 1000 operations/second for typical inputs
- **Scalability**: Linear with input text length

## Compatibility

- **Backward compatible**: Existing French typography behavior unchanged for text without cloze blocks
- **Forward compatible**: New cloze patterns can be added without breaking existing functionality
- **Thread safe**: Service methods are safe for concurrent use