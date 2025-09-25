# Quickstart: Remove Duplicate CSV Header Bug Fix

## Overview
This quickstart validates the fix for the duplicate CSV header bug where ankiprep was incorrectly retaining the original CSV header row alongside the Anki metadata headers.

## Prerequisites
- Go 1.21+ installed
- ankiprep CLI built and available in PATH
- Test CSV files prepared

## Test Scenarios

### Scenario 1: Default Behavior (Remove Header)
Verify that headers are removed by default.

**Input file (`test-input.csv`):**
```csv
Text,Extra,Grammar_Notes
"Hello world","Bonjour monde","Basic greeting"
"Goodbye","Au revoir","Formal farewell"
```

**Commands:**
```bash
# Process with default behavior
./ankiprep test-input.csv test-output.csv

# Verify output
cat test-output.csv
```

**Expected Output (`test-output.csv`):**
```csv
#separator:comma
#html:true
#columns:Text,Extra,Grammar_Notes
"Hello world","Bonjour monde","Basic greeting"
"Goodbye","Au revoir","Formal farewell"
```

**Validation:**
- [ ] Original header row `Text,Extra,Grammar_Notes` is NOT present in output
- [ ] Anki metadata headers are present at the top
- [ ] Both data rows are preserved exactly
- [ ] No duplicate header information

### Scenario 2: Keep Header Flag
Verify that the `--keep-header` flag preserves the first row.

**Input file (`test-input-data.csv`):**
```csv
"Important data","That looks like","A header row"
"Hello world","Bonjour monde","Basic greeting"
"Goodbye","Au revoir","Formal farewell"
```

**Commands:**
```bash
# Process with keep-header flag
./ankiprep --keep-header test-input-data.csv test-output-keep.csv

# Verify output
cat test-output-keep.csv
```

**Expected Output (`test-output-keep.csv`):**
```csv
#separator:comma
#html:true
#columns:Text,Extra,Grammar_Notes
"Important data","That looks like","A header row"
"Hello world","Bonjour monde","Basic greeting"
"Goodbye","Au revoir","Formal farewell"
```

**Validation:**
- [ ] First row is preserved as data (not removed)
- [ ] Anki metadata headers are present at the top
- [ ] All three data rows are included
- [ ] No data loss occurs

### Scenario 3: Short Flag Form
Verify that the `-k` short form works identically.

**Commands:**
```bash
# Process with short flag form
./ankiprep -k test-input-data.csv test-output-short.csv

# Compare outputs
diff test-output-keep.csv test-output-short.csv
```

**Validation:**
- [ ] Short flag produces identical output to long flag
- [ ] No differences in file content

### Scenario 4: Error Handling
Verify proper error handling for various failure modes.

**Commands:**
```bash
# Test missing input file
./ankiprep nonexistent.csv output.csv
echo "Exit code: $?"

# Test permission denied (if applicable)
touch readonly.csv
chmod 444 readonly.csv
./ankiprep test-input.csv readonly.csv
echo "Exit code: $?"

# Test invalid flag
./ankiprep --invalid-flag test-input.csv output.csv
echo "Exit code: $?"
```

**Expected Results:**
- [ ] Missing file error: exit code 1, descriptive error message
- [ ] Permission error: exit code 1, descriptive error message  
- [ ] Invalid flag error: exit code 2, usage message shown
- [ ] All error messages written to stderr, not stdout

### Scenario 5: Backward Compatibility
Verify that existing functionality is unchanged.

**Test with existing Anki files:**
```bash
# Process file that already has some Anki-compatible content
./ankiprep "Pending Flashcards - Cloze.csv" "test-existing.anki.csv"

# Verify output matches expected format
head -n 10 "test-existing.anki.csv"
```

**Validation:**
- [ ] Typography processing (French colon rules) still works
- [ ] Cloze deletion handling remains intact
- [ ] Existing test suite still passes
- [ ] No regression in existing functionality

## Performance Validation

### Large File Test
Test with larger CSV files to ensure performance remains acceptable.

**Commands:**
```bash
# Generate large test file (1000 rows)
for i in {1..1000}; do
  echo "\"Text $i\",\"Extra $i\",\"Notes $i\"" >> large-test.csv
done

# Add header
sed -i '1i Text,Extra,Grammar_Notes' large-test.csv

# Process and time it
time ./ankiprep large-test.csv large-output.csv
```

**Validation:**
- [ ] Processing completes in under 1 second for 1000 rows
- [ ] Memory usage remains reasonable
- [ ] Output file has correct format and row count

## Success Criteria
All scenarios must pass for the bug fix to be considered complete:

- [ ] Default behavior removes headers correctly
- [ ] Keep-header flag preserves first row as data
- [ ] Short flag form works identically to long form
- [ ] Error handling provides appropriate messages and exit codes
- [ ] Backward compatibility maintained
- [ ] Performance remains acceptable for typical use cases

## Cleanup
```bash
# Remove test files
rm -f test-*.csv large-*.csv readonly.csv
```