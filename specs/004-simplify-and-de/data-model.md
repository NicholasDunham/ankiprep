# Data Model: Simplify and De-engineer Codebase

## Entity Overview

Since this is a refactoring project, we're working with existing entities that need to be preserved but potentially simplified in their implementation. The entities represent the core components of the ankiprep CLI tool.

## Core Entities

### CLI Interface
**Purpose**: Command-line interface that users interact with
**Key Attributes**:
- Commands: ankiprep, --help, --version
- Flags: -o/--output, -f/--french, -q/--smart-quotes, -s/--skip-duplicates, -k/--keep-header, -v/--verbose
- Input arguments: file paths (CSV/TSV)
- Exit codes: 0 for success, non-zero for errors

**Relationships**: 
- Invokes File Processor for processing operations
- Uses Typography Engine for text formatting
- Coordinates with Duplicate Detector and Output Generator

**State Transitions**: Command parsing → Validation → Processing → Output generation → Exit

**Validation Rules**:
- At least one input file required
- Output path must be writable if specified
- Input files must exist and be readable
- File extensions must be .csv or .tsv

### File Processor
**Purpose**: Core functionality that transforms input CSV/TSV files
**Key Attributes**:
- Input file paths
- Processing configuration (from CLI flags)
- Parsed CSV data structures
- Processing state and progress

**Relationships**: 
- Receives configuration from CLI Interface
- Uses Typography Engine for text transformations
- Coordinates with Duplicate Detector for deduplication
- Provides data to Output Generator

**State Transitions**: File reading → Parsing → Processing → Data transformation → Ready for output

**Validation Rules**:
- Files must be valid UTF-8 encoded CSV/TSV
- Must have at least two columns
- Header row handling based on --keep-header flag

### Typography Engine  
**Purpose**: Handles French punctuation and smart quotes formatting
**Key Attributes**:
- French punctuation rules (thin spaces before :;!?)
- Smart quotes conversion patterns (" → " and " → ")
- Text transformation state
- Unicode normalization settings

**Relationships**: 
- Invoked by File Processor during text transformation
- Processes individual text fields from CSV data

**State Transitions**: Text input → Pattern matching → Transformation → Normalized output

**Validation Rules**:
- Must preserve text meaning while applying formatting
- Must handle Unicode characters correctly
- Must be reversible/consistent in transformations

### Duplicate Detector
**Purpose**: Logic for identifying and removing duplicate entries  
**Key Attributes**:
- Entry comparison criteria (identical content)
- Deduplication algorithm state
- Duplicate tracking data structures
- Skip flag configuration

**Relationships**: 
- Invoked by File Processor during processing
- Operates on parsed CSV data before output generation

**State Transitions**: Data input → Comparison analysis → Duplicate marking → Filtered output

**Validation Rules**:
- Must maintain stable ordering of non-duplicate entries
- Must preserve first occurrence of duplicate content
- Must respect --skip-duplicates flag setting

### Output Generator
**Purpose**: Creates properly formatted output files
**Key Attributes**:
- Output file path and format
- Anki-compatible CSV formatting rules
- UTF-8 encoding configuration
- File writing state

**Relationships**: 
- Receives processed data from File Processor
- Final component in processing pipeline

**State Transitions**: Data input → Format validation → File writing → Completion confirmation

**Validation Rules**:
- Must produce valid Anki-compatible CSV format
- Must maintain UTF-8 encoding
- Must preserve all data integrity during output
- Must handle file system errors gracefully

## Simplification Targets

### Current Complexity Issues
- Over-abstracted interfaces where simple structs would suffice
- Excessive dependency injection for straightforward operations
- Complex error propagation chains for simple file operations
- Unnecessary separation of concerns for single-user CLI tool

### Simplification Constraints
- Must preserve all entity behaviors and relationships
- Must maintain identical input/output characteristics
- Must keep all validation rules intact
- Must not break existing test expectations

### Success Criteria
- Reduced code complexity while maintaining functionality
- Simplified entity implementations without changing interfaces
- Consolidated duplicate logic across entities
- Improved code readability and maintainability