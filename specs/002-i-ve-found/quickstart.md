# Quickstart: Cloze Deletion Colon Exception

**Feature**: French typography processing with Anki cloze deletion support  
**Prerequisites**: ankiprep CLI tool with French typography feature installed

## Quick Test Scenarios

### Scenario 1: Basic Cloze Processing
**Goal**: Verify cloze syntax colons are preserved while content colons get NNBSP

```bash
# Create test input
echo 'Text,Extra
"Je vais {{c1::essayer}} cette recette : c'"'"'est délicieux !",""' > test-basic.csv

# Process with French typography
ankiprep --french test-basic.csv test-basic-output.csv

# Verify output
cat test-basic-output.csv
# Expected: Cloze colons (::) unchanged, content colon (:) gets NNBSP
```

**Success Criteria**:
- ✅ Cloze syntax `{{c1::essayer}}` preserved exactly  
- ✅ Content colon `:` before "c'est" gets NNBSP
- ✅ No warnings in stderr
- ✅ Processing completes successfully

### Scenario 2: Multiple Cloze Blocks  
**Goal**: Verify independent processing of multiple cloze blocks

```bash
# Create test input with multiple blocks
echo 'Text,Extra
"{{c1::Bonjour}} : comment {{c2::allez-vous}} ?",""' > test-multiple.csv

# Process 
ankiprep -f test-multiple.csv test-multiple-output.csv

# Verify output
cat test-multiple-output.csv
# Expected: Both cloze blocks preserved, content colon gets NNBSP
```

**Success Criteria**:
- ✅ Both `{{c1::Bonjour}}` and `{{c2::allez-vous}}` preserved
- ✅ Content colon `:` between blocks gets NNBSP  
- ✅ Each cloze block processed independently
- ✅ No warnings generated

### Scenario 3: Nested Content Colons
**Goal**: Verify content colons within cloze blocks get proper typography

```bash
# Create test with nested punctuation
echo 'Text,Extra  
"{{c1::phrase with « quoted text : example »}}","More : content"' > test-nested.csv

# Process
ankiprep --french test-nested.csv test-nested-output.csv

# Verify output  
cat test-nested-output.csv
# Expected: Cloze syntax (::) preserved, nested colon (:) gets NNBSP, external colon gets NNBSP
```

**Success Criteria**:
- ✅ Cloze syntax `{{c1::` preserved without NNBSP
- ✅ Nested colon in "text : example" gets NNBSP per French rules
- ✅ External colon in "More : content" gets NNBSP
- ✅ Smart quote processing works within cloze blocks (if enabled)

### Scenario 4: Malformed Cloze Handling
**Goal**: Verify graceful handling of malformed cloze blocks

```bash
# Create test with malformed cloze
echo 'Text,Extra
"{{c1::incomplete block missing close","Normal text : with colon"' > test-malformed.csv

# Process (should show warning but continue)
ankiprep --french test-malformed.csv test-malformed-output.csv

# Check stderr for warnings
# Expected: Warning logged, but processing continues
```

**Success Criteria**:
- ✅ Warning message logged to stderr about malformed block
- ✅ Processing continues successfully  
- ✅ Malformed field processed as regular French text (all colons get NNBSP)
- ✅ Other fields processed normally

### Scenario 5: Complex Cloze with Hints
**Goal**: Verify cloze blocks with hint text are handled correctly

```bash
# Create test with cloze hints
echo 'Text,Extra
"{{c1::essayer::to try}} de faire : c'"'"'est difficile !",""' > test-hints.csv

# Process
ankiprep -f test-hints.csv test-hints-output.csv  

# Verify output
cat test-hints-output.csv
# Expected: All cloze syntax colons preserved, content colon gets NNBSP
```

**Success Criteria**:  
- ✅ Both cloze syntax colons in `{{c1::essayer::to try}}` preserved
- ✅ Content colon in "faire : c'est" gets NNBSP
- ✅ Hint text "to try" processed with French typography rules if needed
- ✅ No warnings generated

## Integration Test

### Full Workflow Test
**Goal**: Process a realistic CSV file with mixed content

```bash
# Create comprehensive test file
cat << 'EOF' > comprehensive-test.csv
Text,Extra,Notes
"Je {{c1::mange}} une pomme : délicieuse !","","Simple cloze"
"{{c1::Bonjour}} : comment {{c2::allez-vous::how are you}} ?","Context","Multiple with hint"  
"Text without cloze : just French content","","No cloze blocks"
"{{c1::malformed missing close","","Malformed block"
"{{c1::phrase « with : punctuation »}} : more text","","Nested content"
EOF

# Process with progress  
ankiprep --french comprehensive-test.csv comprehensive-output.csv

# Verify comprehensive results
echo "=== Processed Output ==="
cat comprehensive-output.csv

echo -e "\n=== Expected Behavior ==="
echo "- Row 1: Cloze preserved, content colon gets NNBSP"
echo "- Row 2: Both cloze blocks preserved, content colons get NNBSP"  
echo "- Row 3: Regular French processing (colon gets NNBSP)"
echo "- Row 4: Warning logged, processed as regular text"
echo "- Row 5: Cloze syntax preserved, nested colon gets NNBSP"
```

**Success Criteria**:
- ✅ All valid cloze blocks identified and preserved
- ✅ Content colons consistently get NNBSP treatment
- ✅ One warning logged for malformed cloze block  
- ✅ Processing summary shows correct counts
- ✅ Output file structure matches input structure

## Performance Verification

### Large File Test
**Goal**: Verify performance with larger datasets

```bash
# Generate larger test file (1000 rows)
python3 -c "
import csv
with open('large-test.csv', 'w', newline='', encoding='utf-8') as f:
    writer = csv.writer(f)
    writer.writerow(['Text', 'Extra'])
    for i in range(1000):
        text = f'Row {i}: Je {{c1::mange}} du pain : c\\'est bon !'
        writer.writerow([text, ''])
"

# Time the processing
time ankiprep --french large-test.csv large-output.csv

# Verify timing is reasonable (should be sub-second)
```

**Success Criteria**:
- ✅ Processing completes in under 1 second
- ✅ All 1000 cloze blocks processed correctly  
- ✅ Memory usage remains reasonable (< 100MB)
- ✅ No performance degradation compared to non-cloze processing

## Cleanup

```bash
# Remove test files
rm -f test-*.csv comprehensive-*.csv large-*.csv
echo "Quickstart test cleanup complete"
```

## Troubleshooting

### Common Issues

**Issue**: Cloze blocks not detected
- **Check**: Verify syntax uses double braces `{{}}` not single `{}`
- **Check**: Ensure cloze number is present: `{{c1::...}}` not `{{::...}}`

**Issue**: Warnings about malformed blocks
- **Check**: Verify closing `}}` brackets are present  
- **Check**: Ensure valid cloze number (digits only)

**Issue**: Content colons not getting NNBSP
- **Check**: Verify `--french` flag is enabled
- **Check**: Ensure colons are outside cloze syntax (not part of `::`)

**Issue**: Performance slower than expected  
- **Check**: File encoding is UTF-8
- **Check**: No extremely long fields (> 10KB per field)

## Next Steps

After successful quickstart:
1. Run full test suite: `go test ./...`
2. Test with your actual Anki CSV files  
3. Verify output in Anki import process
4. Set up automated regression testing if needed