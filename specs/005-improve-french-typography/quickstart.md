# French Typography Feature Quickstart

**Purpose**: Validate French typography processing implementation through targeted test scenarios  
**Target**: Developers implementing and testing the French typography feature  
**Context**: ankiprep CLI application with French NNBSP handling

## Quick Validation Scenarios

### Scenario 1: Basic Quote Processing

**Goal**: Verify NNBSP insertion around French angle quotes

**Setup**:
```bash
# Create test input
echo -e "Front\tBack\n«bonjour»\thello" > test_quotes.tsv
```

**Execute**:
```bash
ankiprep test_quotes.tsv output_quotes.tsv
```

**Expected Result**:
```
Front	Back
« bonjour »	hello
```
*Note: Spaces around "bonjour" are NNBSP (U+202F)*

**Validation**:
- Check NNBSP characters using: `hexdump -C output_quotes.tsv | grep 202f`
- Verify quote characters remain unchanged: `«` (U+00AB), `»` (U+00BB)

### Scenario 2: Punctuation Spacing

**Goal**: Verify NNBSP insertion before French double punctuation

**Setup**:
```bash
# Create test input
echo -e "Question\tAnswer\nComment allez-vous?\tHow are you?" > test_punct.tsv
```

**Execute**:
```bash
ankiprep test_punct.tsv output_punct.tsv
```

**Expected Result**:
```
Question	Answer
Comment allez-vous ?	How are you?
```
*Note: Space before "?" is NNBSP (U+202F)*

**Validation**:
- Confirm NNBSP before punctuation: `hexdump -C output_punct.tsv | grep "20 2f"`
- Test all punctuation: `:`, `;`, `!`, `?`

### Scenario 3: Preserve Existing NNBSP

**Goal**: Ensure existing NNBSP characters are not duplicated

**Setup**:
```bash
# Create input with existing NNBSP (use printf for precise control)
printf "Text\tNote\nCorrect\t« Bonjour : test ? »\n" > test_preserve.tsv
# Note: Use actual NNBSP characters in the input
```

**Execute**:
```bash
ankiprep test_preserve.tsv output_preserve.tsv
```

**Expected Result**:
- File should remain unchanged
- No duplicate NNBSP characters
- Same byte count as input file

**Validation**:
```bash
# Compare file sizes
wc -c test_preserve.tsv output_preserve.tsv

# Compare content
diff test_preserve.tsv output_preserve.tsv
# Should show no differences
```

### Scenario 4: Mixed Typography Elements

**Goal**: Test complex French text with multiple typography rules

**Setup**:
```bash
# Create complex test case
echo -e "Card\tContent\nGreeting\t«Bonjour: comment allez-vous; très bien!»" > test_complex.tsv
```

**Execute**:
```bash
ankiprep test_complex.tsv output_complex.tsv
```

**Expected Result**:
```
Card	Content
Greeting	« Bonjour : comment allez-vous ; très bien ! »
```
*Note: All spaces before punctuation and around quotes are NNBSP*

**Validation**:
- Count NNBSP occurrences: `grep -o $'\u202f' output_complex.tsv | wc -l`
- Expected count: 6 NNBSP characters

### Scenario 5: Error Handling

**Goal**: Verify proper error responses and exit codes

**Test 5a - Missing File**:
```bash
ankiprep nonexistent.tsv output.tsv
echo "Exit code: $?"
```
**Expected**: Exit code 1, error message about file not found

**Test 5b - Permission Error**:
```bash
echo -e "test\tdata" > input.tsv
mkdir readonly && chmod 444 readonly
ankiprep input.tsv readonly/output.tsv
echo "Exit code: $?"
```
**Expected**: Exit code 1, error message about write permissions

## Development Verification Checklist

### Unit Test Coverage
- [ ] Quote processing (with/without existing spaces)
- [ ] Punctuation processing (all four types: `:;!?`)
- [ ] NNBSP preservation detection
- [ ] UTF-8 validation and error handling
- [ ] Empty input handling

### Integration Test Coverage  
- [ ] End-to-end file processing
- [ ] Multi-column TSV handling
- [ ] Header row preservation
- [ ] Large file processing
- [ ] Error condition handling

### Manual Test Cases
- [ ] All quickstart scenarios pass
- [ ] Unicode character verification with hex tools
- [ ] Cross-platform file handling (Windows/macOS/Linux)
- [ ] Performance with typical Anki deck sizes

## Character Reference Guide

For manual verification and debugging:

| Character | Unicode | Hex | Description |
|-----------|---------|-----|-------------|
| (space) | U+0020 | 20 | Regular space (input) |
| (nnbsp) | U+202F | E2 80 AF | Narrow non-breaking space (output) |
| « | U+00AB | C2 AB | Left-pointing double angle quotation mark |
| » | U+00BB | C2 BB | Right-pointing double angle quotation mark |
| : | U+003A | 3A | Colon |
| ; | U+003B | 3B | Semicolon |
| ! | U+0021 | 21 | Exclamation mark |
| ? | U+003F | 3F | Question mark |

## Debugging Commands

**View file in hex**:
```bash
hexdump -C filename.tsv | head -20
```

**Search for NNBSP**:
```bash  
grep -P '\x{202F}' filename.tsv
```

**Count character occurrences**:
```bash
grep -o $'\u202f' filename.tsv | wc -l
```

**Compare Unicode normalization**:
```bash
python3 -c "print(repr(open('filename.tsv').read()))"
```