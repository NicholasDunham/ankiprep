package acceptance_test

import (
	"bytes"
	"encoding/csv"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// TestAcceptance_ClozeColonException tests the complete user workflow
func TestAcceptance_ClozeColonException(t *testing.T) {
	// This is the ultimate acceptance test - it will fail until everything is implemented

	// Build the CLI binary first
	tmpDir, err := os.MkdirTemp("", "ankiprep_acceptance")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	binaryPath := filepath.Join(tmpDir, "ankiprep")

	// Build the binary
	buildCmd := exec.Command("go", "build", "-o", binaryPath, "../../../cmd/ankiprep")
	buildCmd.Dir = tmpDir

	buildOutput, err := buildCmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Failed to build binary: %v\nOutput: %s", err, buildOutput)
	}

	t.Run("user story: French teacher processing cards", func(t *testing.T) {
		// Scenario: French teacher has Anki cards with cloze deletions
		// They want French typography applied but NOT inside cloze blocks

		inputFile := filepath.Join(tmpDir, "teacher_cards.csv")
		outputFile := filepath.Join(tmpDir, "processed_cards.csv")

		// Create realistic teacher data
		teacherData := [][]string{
			{"Front", "Back", "Tags"},
			{
				"Géographie : Quelle est la capitale de la France ?",
				"La capitale de la France est {{c1::Paris : la ville lumière}}.",
				"géographie::capitales",
			},
			{
				"Histoire : Qui était Napoléon ?",
				"Napoléon était {{c1::empereur des Français : 1804-1814}} et {{c2::général : militaire}}.",
				"histoire::napoléon",
			},
			{
				"Littérature : Citez un auteur français célèbre.",
				"{{c1::Victor Hugo : Les Misérables}} est un auteur français célèbre.",
				"littérature::auteurs",
			},
			{
				"Sciences : Qu'est-ce que H2O ?",
				"H2O est la formule de {{c1::l'eau : composé chimique}}.",
				"sciences::chimie",
			},
		}

		// Write teacher's input file
		file, err := os.Create(inputFile)
		if err != nil {
			t.Fatalf("Failed to create teacher input file: %v", err)
		}

		writer := csv.NewWriter(file)
		err = writer.WriteAll(teacherData)
		file.Close()
		if err != nil {
			t.Fatalf("Failed to write teacher CSV: %v", err)
		}

		// Run the CLI as the teacher would
		cmd := exec.Command(binaryPath,
			"--input", inputFile,
			"--output", outputFile,
			"--french",
			"--smart-quotes",
			"--progress",
		)

		var stdout, stderr bytes.Buffer
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr

		err = cmd.Run()
		if err != nil {
			t.Fatalf("CLI execution failed: %v\nStdout: %s\nStderr: %s",
				err, stdout.String(), stderr.String())
		}

		// Verify the teacher gets what they expect
		outputFileHandle, err := os.Open(outputFile)
		if err != nil {
			t.Fatalf("Failed to open teacher output file: %v", err)
		}
		defer outputFileHandle.Close()

		reader := csv.NewReader(outputFileHandle)
		outputRecords, err := reader.ReadAll()
		if err != nil {
			t.Fatalf("Failed to read teacher output CSV: %v", err)
		}

		// Define expected results for teacher
		expectedResults := []struct {
			description   string
			rowIndex      int
			expectedFront string
			expectedBack  string
			reasoning     string
		}{
			{
				description:   "Geography question with cloze",
				rowIndex:      1,
				expectedFront: "Géographie\u00A0: Quelle est la capitale de la France ?",
				expectedBack:  "La capitale de la France est {{c1::Paris : la ville lumière}}.",
				reasoning:     "Colon in question gets non-breaking space, colon in cloze stays normal",
			},
			{
				description:   "History question with multiple cloze",
				rowIndex:      2,
				expectedFront: "Histoire\u00A0: Qui était Napoléon ?",
				expectedBack:  "Napoléon était {{c1::empereur des Français : 1804-1814}} et {{c2::général : militaire}}.",
				reasoning:     "Front colon processed, both cloze colons preserved",
			},
			{
				description:   "Literature question",
				rowIndex:      3,
				expectedFront: "Littérature\u00A0: Citez un auteur français célèbre.",
				expectedBack:  "{{c1::Victor Hugo : Les Misérables}} est un auteur français célèbre.",
				reasoning:     "Front colon processed, cloze colon preserved",
			},
			{
				description:   "Science question",
				rowIndex:      4,
				expectedFront: "Sciences\u00A0: Qu'est-ce que H2O ?",
				expectedBack:  "H2O est la formule de {{c1::l'eau : composé chimique}}.",
				reasoning:     "Front colon processed, cloze colon preserved",
			},
		}

		// Verify each expected result
		for _, expected := range expectedResults {
			t.Run(expected.description, func(t *testing.T) {
				if len(outputRecords) <= expected.rowIndex {
					t.Fatalf("Missing row %d in teacher output", expected.rowIndex)
				}

				record := outputRecords[expected.rowIndex]
				if len(record) < 3 {
					t.Fatalf("Row %d has insufficient columns", expected.rowIndex)
				}

				// Check front field
				if record[0] != expected.expectedFront {
					t.Errorf("Front field incorrect:\nGot:      %q\nExpected: %q\nReasoning: %s",
						record[0], expected.expectedFront, expected.reasoning)
				}

				// Check back field
				if record[1] != expected.expectedBack {
					t.Errorf("Back field incorrect:\nGot:      %q\nExpected: %q\nReasoning: %s",
						record[1], expected.expectedBack, expected.reasoning)
				}

				// Tags should be unchanged
				originalTags := teacherData[expected.rowIndex][2]
				if record[2] != originalTags {
					t.Errorf("Tags were modified: got %q, want %q", record[2], originalTags)
				}
			})
		}

		// Verify progress reporting
		stderrOutput := stderr.String()
		if !strings.Contains(stderrOutput, "Processing") {
			t.Error("Teacher should see progress updates")
		}

		// Verify summary information
		stdoutOutput := stdout.String()
		if !strings.Contains(stdoutOutput, "cards processed") {
			t.Error("Teacher should see processing summary")
		}
	})

	t.Run("user story: edge cases and error handling", func(t *testing.T) {
		// Test malformed input that a real user might have
		inputFile := filepath.Join(tmpDir, "malformed_cards.csv")
		outputFile := filepath.Join(tmpDir, "processed_malformed.csv")

		malformedData := [][]string{
			{"Front", "Back", "Tags"},
			{
				"Question : Normal question",
				"Answer with {{c1::incomplete cloze", // Malformed cloze
				"test",
			},
			{
				"Another question :",            // Trailing colon
				"{{c0::Invalid number}} answer", // Invalid cloze number
				"test",
			},
			{
				"Empty cloze test:",
				"This has {{c1::}} empty cloze", // Empty cloze content
				"test",
			},
		}

		// Write malformed input
		file, err := os.Create(inputFile)
		if err != nil {
			t.Fatalf("Failed to create malformed input file: %v", err)
		}

		writer := csv.NewWriter(file)
		err = writer.WriteAll(malformedData)
		file.Close()
		if err != nil {
			t.Fatalf("Failed to write malformed CSV: %v", err)
		}

		// Run CLI - should handle errors gracefully
		cmd := exec.Command(binaryPath,
			"--input", inputFile,
			"--output", outputFile,
			"--french",
		)

		var stdout, stderr bytes.Buffer
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr

		err = cmd.Run()
		// CLI should not crash, but may return non-zero for warnings

		// Verify output file exists (processing should continue despite errors)
		if _, err := os.Stat(outputFile); os.IsNotExist(err) {
			t.Fatal("Output file should be created even with malformed input")
		}

		// Should have warnings about malformed content
		stderrOutput := stderr.String()
		if !strings.Contains(stderrOutput, "warning") && !strings.Contains(stderrOutput, "Warning") {
			t.Error("Should warn about malformed cloze blocks")
		}

		// Read output to verify graceful handling
		outputFileHandle, err := os.Open(outputFile)
		if err != nil {
			t.Fatalf("Failed to open malformed output file: %v", err)
		}
		defer outputFileHandle.Close()

		reader := csv.NewReader(outputFileHandle)
		outputRecords, err := reader.ReadAll()
		if err != nil {
			t.Fatalf("Failed to read malformed output CSV: %v", err)
		}

		// Should have same number of rows as input
		if len(outputRecords) != len(malformedData) {
			t.Errorf("Output should have %d rows, got %d", len(malformedData), len(outputRecords))
		}

		// French typography should still be applied where possible
		if len(outputRecords) > 1 {
			firstRow := outputRecords[1]
			if len(firstRow) > 0 && strings.Contains(firstRow[0], " :") {
				t.Error("French colon rule should be applied even with malformed back field")
			}
		}
	})

	t.Run("user story: CLI help and usage", func(t *testing.T) {
		// User runs --help to understand how to use the tool
		cmd := exec.Command(binaryPath, "--help")

		var stdout bytes.Buffer
		cmd.Stdout = &stdout

		err := cmd.Run()
		if err != nil {
			t.Fatalf("Help command failed: %v", err)
		}

		helpOutput := stdout.String()

		// Should contain key information a user needs
		requiredHelpContent := []string{
			"ankiprep",          // Program name
			"CSV",               // File format
			"French typography", // Main feature
			"cloze",             // Key feature
			"--input",           // Required flag
			"--output",          // Required flag
			"--french",          // Main option
			"example",           // Usage example
		}

		for _, content := range requiredHelpContent {
			if !strings.Contains(strings.ToLower(helpOutput), strings.ToLower(content)) {
				t.Errorf("Help output should contain %q for user guidance", content)
			}
		}
	})
}
