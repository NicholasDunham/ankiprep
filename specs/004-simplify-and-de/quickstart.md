# Quickstart: Simplify and De-engineer Codebase

## Objective
Verify that the simplified ankiprep codebase maintains identical functionality while reducing complexity.

## Prerequisites
- Go 1.21+ installed
- Access to ankiprep repository
- Test data files in examples/ directory

## Validation Steps

### Step 1: Pre-Refactoring Baseline
```bash
# Build current version
go build -o ankiprep-before ./cmd/ankiprep

# Create baseline outputs for comparison
./ankiprep-before examples/spanish_vocabulary.tsv -o baseline-simple.csv
./ankiprep-before examples/spanish_vocabulary.tsv -f -q -o baseline-formatted.csv
./ankiprep-before examples/spanish_vocabulary.tsv -s -k -o baseline-dedup.csv
./ankiprep-before examples/spanish_vocabulary.tsv -fqs -v -o baseline-full.csv

# Run full test suite to establish baseline
go test ./... -v > baseline-tests.log 2>&1
```

### Step 2: Complexity Measurement
```bash  
# Measure current code complexity (lines of code, cyclomatic complexity)
find . -name "*.go" -not -path "./vendor/*" | xargs wc -l > complexity-before.txt
gocyclo -over 5 . > cyclomatic-before.txt

# Count abstraction layers
echo "Current package count:" >> complexity-before.txt
find . -name "*.go" -not -path "./vendor/*" -exec dirname {} \; | sort -u | wc -l >> complexity-before.txt
```

### Step 3: Incremental Refactoring Validation
After each refactoring increment:
```bash
# Verify tests still pass
go test ./... -v

# Build and test CLI functionality  
go build -o ankiprep-current ./cmd/ankiprep

# Test basic functionality
./ankiprep-current examples/spanish_vocabulary.tsv -o test-simple.csv
diff baseline-simple.csv test-simple.csv

# Test formatting options
./ankiprep-current examples/spanish_vocabulary.tsv -f -q -o test-formatted.csv  
diff baseline-formatted.csv test-formatted.csv

# Test deduplication and headers
./ankiprep-current examples/spanish_vocabulary.tsv -s -k -o test-dedup.csv
diff baseline-dedup.csv test-dedup.csv

# Test full feature combination
./ankiprep-current examples/spanish_vocabulary.tsv -fqs -v -o test-full.csv
diff baseline-full.csv test-full.csv
```

### Step 4: Final Validation
```bash
# Measure final complexity
find . -name "*.go" -not -path "./vendor/*" | xargs wc -l > complexity-after.txt
gocyclo -over 5 . > cyclomatic-after.txt

echo "Final package count:" >> complexity-after.txt
find . -name "*.go" -not -path "./vendor/*" -exec dirname {} \; | sort -u | wc -l >> complexity-after.txt

# Compare complexity metrics
echo "=== COMPLEXITY COMPARISON ===" > complexity-comparison.txt
echo "BEFORE:" >> complexity-comparison.txt
cat complexity-before.txt >> complexity-comparison.txt
echo "AFTER:" >> complexity-comparison.txt  
cat complexity-after.txt >> complexity-comparison.txt

# Verify all CLI commands work identically
./ankiprep-current --help > help-after.txt
diff help-before.txt help-after.txt

./ankiprep-current --version > version-after.txt  
diff version-before.txt version-after.txt
```

### Step 5: Performance Validation
```bash
# Test with larger files to ensure no performance regression
time ./ankiprep-before large_test_file.csv -o before-large.csv
time ./ankiprep-current large_test_file.csv -o after-large.csv

# Compare outputs
diff before-large.csv after-large.csv

# Memory usage comparison (if available)
/usr/bin/time -v ./ankiprep-before large_test_file.csv -o before-memory.csv 2> memory-before.log
/usr/bin/time -v ./ankiprep-current large_test_file.csv -o after-memory.csv 2> memory-after.log
```

## Success Criteria

### Functional Equivalence
- [ ] All diff commands show no differences in output files
- [ ] All tests pass without modification
- [ ] CLI help and version output unchanged
- [ ] Error handling behavior identical
- [ ] Performance within acceptable variance (±10%)

### Complexity Reduction  
- [ ] Reduced total lines of code (≥10% reduction)
- [ ] Reduced cyclomatic complexity metrics
- [ ] Fewer total packages/abstraction layers
- [ ] Eliminated unnecessary configuration files
- [ ] Consolidated duplicate logic

### Code Quality
- [ ] Improved code readability and maintainability
- [ ] Simplified dependency relationships
- [ ] Reduced interface proliferation
- [ ] Enhanced Go idiom compliance

## Troubleshooting

### If Tests Fail
1. Identify specific failing test
2. Compare expected vs actual behavior
3. Revert last increment and retry with smaller change
4. Update test only if behavior change is intentional and approved

### If Output Differs
1. Check if difference is in formatting/whitespace only
2. Verify UTF-8 encoding preservation
3. Compare character-by-character for subtle differences
4. Check if difference affects Anki import compatibility

### If Performance Regresses
1. Profile before and after implementations
2. Identify performance bottleneck
3. Consider if simplification introduced inefficiency
4. Balance simplicity vs performance needs