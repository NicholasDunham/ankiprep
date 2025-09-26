# ankiprep

A lightweight Go CLI tool for processing CSV/TSV files into Anki-compatible format with French typography enhancement.

## Features

- **File Processing**: Convert CSV and TSV files to Anki format
- **French Typography**: Smart quotes and proper French punctuation spacing
- **Multi-file Support**: Process and merge multiple input files  
- **Duplicate Detection**: Automatic duplicate entry removal
- **Header Management**: Flexible header preservation options
- **Memory Efficient**: Optimized for large file processing

## Installation

```bash
# Build from source
go build -o ankiprep ./cmd/ankiprep

# Or use go run for development
go run ./cmd/ankiprep --help
```

## Quick Start

```bash
# Basic processing
./ankiprep input.csv

# French text with typography enhancement  
./ankiprep french_vocab.csv -o anki_ready.csv

# Process multiple files
./ankiprep file1.csv file2.tsv -o combined.csv -v
```

## Usage

```bash
# Basic usage
./ankiprep input.csv

# Multiple files with output path
./ankiprep file1.csv file2.csv -o combined.csv

# French typography and smart quotes
./ankiprep -f -q input.csv

# Skip duplicate detection for faster processing
./ankiprep -s input.csv

# Verbose output
./ankiprep -v input.csv

# Combine multiple short flags
./ankiprep -fqvs input1.csv input2.csv -o output.csv

# Keep original CSV headers in output
./ankiprep --keep-header input.csv
```

### Command Options

- `-o, --output`: Specify output file path
- `-f, --french`: Add thin spaces before French punctuation (:;!?)  
- `-q, --smart-quotes`: Convert straight quotes to curly quotes
- `-s, --skip-duplicates`: Remove entries with identical content
- `-k, --keep-header`: Preserve the first row of CSV files (default: remove header)
- `-v, --verbose`: Enable verbose output

## Input Format

CSV files should have at least two columns with a header row:

```csv
Front,Back
Hello,Bonjour
Goodbye,Au revoir
```

Supports CSV (`.csv`) and TSV (`.tsv`) files with UTF-8 encoding.

## Output

Creates Anki-compatible CSV files with proper escaping and UTF-8 encoding.

## Development

### Project Structure

```text
cmd/ankiprep/        # Main CLI application
  main.go            # All processing logic (377 lines)
  errors.go          # Error handling utilities
internal/models/     # Core data structures  
  *.go               # Data models and typography processing
tests/               # Comprehensive test suite
  unit/              # Unit tests for models and core logic
  integration/       # CLI integration tests  
  performance/       # Performance and memory tests
```

### Architecture

The ankiprep codebase follows a **simplified direct processing architecture**:

- **Single Entry Point**: All CSV→Anki transformation logic in `cmd/ankiprep/main.go`
- **Essential Models**: Core data structures for file processing and typography
- **CLI-Focused Testing**: Tests validate real user workflows via command-line interface

This design prioritizes **simplicity and maintainability** over complex abstractions.

### Building and Testing

```bash
# Run all tests
go test ./...

# Run specific test categories
go test ./tests/unit/...
go test ./tests/integration/...

# Build optimized binary
go build -ldflags "-s -w" -o ankiprep ./cmd/ankiprep

# Development with live rebuilding
go run ./cmd/ankiprep --help
```

### Contributing

The project emphasizes **constitutional simplicity principles**:

1. **Simplicity**: Prefer direct code over abstractions
2. **Directness**: Avoid unnecessary indirection layers  
3. **Efficiency**: Optimize for common use cases
4. **Maintainability**: Single source of truth for business logic
5. **Testability**: Focus on end-user scenarios via CLI interface
