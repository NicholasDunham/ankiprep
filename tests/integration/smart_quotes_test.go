package integration

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// TestSmartQuotesConversionBasic tests basic straight to smart quote conversion
func TestSmartQuotesConversionBasic(t *testing.T) {
	// Create temporary directory
	tempDir, err := ioutil.TempDir("", "ankiprep_integration")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create input file with straight quotes
	inputFile := filepath.Join(tempDir, "input.csv")
	csvContent := `front,back
"He said ""Hello"" to me","Il m'a dit ""Bonjour"""
"The word ""test"" is important","Le mot ""test"" est important"
"She asked ""Why?""","Elle a demandé ""Pourquoi ?"""`

	err = ioutil.WriteFile(inputFile, []byte(csvContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write input file: %v", err)
	}

	outputFile := filepath.Join(tempDir, "output.csv")

	// Execute ankiprep (smart quotes should be enabled by default)
	cmd := exec.Command("ankiprep", "-o", outputFile, inputFile)
	_, err = cmd.CombinedOutput()

	if err != nil {
		t.Fatalf("ankiprep command failed: %v", err)
	}

	// Read output file
	outputContent, err := ioutil.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	outputStr := string(outputContent)

	// Check for smart quote conversions
	// " should become " (U+201C) and " (U+201D)
	expectedTransformations := []string{
		"\u201cHello\u201d",   // "Hello" -> "Hello"
		"\u201cBonjour\u201d", // "Bonjour" -> "Bonjour"
		"\u201ctest\u201d",    // "test" -> "test"
		"\u201cWhy?\u201d",    // "Why?" -> "Why?"
		"\u201cPourquoi",      // "Pourquoi -> "Pourquoi
	}

	for _, expected := range expectedTransformations {
		if !strings.Contains(outputStr, expected) {
			t.Errorf("Smart quotes should convert to contain '%s', but got: %s", expected, outputStr)
		}
	}
}

// TestSmartQuotesApostrophes tests apostrophe conversion from straight to smart
func TestSmartQuotesApostrophes(t *testing.T) {
	// Create temporary directory
	tempDir, err := ioutil.TempDir("", "ankiprep_integration")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create input file with apostrophes
	inputFile := filepath.Join(tempDir, "input.csv")
	csvContent := `front,back
"I'm happy","Je suis content"
"It's working","Ça marche"
"Don't worry","Ne t'inquiète pas"
"We're here","Nous sommes là"`

	err = ioutil.WriteFile(inputFile, []byte(csvContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write input file: %v", err)
	}

	outputFile := filepath.Join(tempDir, "output.csv")

	// Execute ankiprep
	cmd := exec.Command("ankiprep", "-o", outputFile, inputFile)
	_, err = cmd.CombinedOutput()

	if err != nil {
		t.Fatalf("ankiprep command failed: %v", err)
	}

	// Read output file
	outputContent, err := ioutil.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	outputStr := string(outputContent)

	// Check for smart apostrophe conversions
	// ' should become ' (U+2019)
	expectedTransformations := []string{
		"I\u2019m happy",    // I'm -> I'm
		"It\u2019s working", // It's -> It's
		"Don\u2019t worry",  // Don't -> Don't
		"We\u2019re here",   // We're -> We're
		"t\u2019inquiète",   // t'inquiète -> t'inquiète
	}

	for _, expected := range expectedTransformations {
		if !strings.Contains(outputStr, expected) {
			t.Errorf("Smart quotes should convert apostrophes to contain '%s', but got: %s", expected, outputStr)
		}
	}
}

// TestSmartQuotesNested tests proper handling of nested quotes
func TestSmartQuotesNested(t *testing.T) {
	// Create temporary directory
	tempDir, err := ioutil.TempDir("", "ankiprep_integration")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create input file with nested quotes
	inputFile := filepath.Join(tempDir, "input.csv")
	csvContent := `front,back
"He said ""I'm going to the 'store' today""","Il a dit ""Je vais au 'magasin' aujourd'hui"""
"The book ""Tom's Adventure"" is great","Le livre ""L'Aventure de Tom"" est génial"`

	err = ioutil.WriteFile(inputFile, []byte(csvContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write input file: %v", err)
	}

	outputFile := filepath.Join(tempDir, "output.csv")

	// Execute ankiprep
	cmd := exec.Command("ankiprep", "-o", outputFile, inputFile)
	_, err = cmd.CombinedOutput()

	if err != nil {
		t.Fatalf("ankiprep command failed: %v", err)
	}

	// Read output file
	outputContent, err := ioutil.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	outputStr := string(outputContent)

	// Check for proper nested quote handling
	expectedTransformations := []string{
		"\u201cI\u2019m going",  // "I'm -> "I'm
		"\u2018store\u2019",     // 'store' -> 'store'
		"\u201cJe vais",         // "Je -> "Je
		"\u2018magasin\u2019",   // 'magasin' -> 'magasin'
		"\u201cTom\u2019s",      // "Tom's -> "Tom's
		"\u201cL\u2019Aventure", // "L'Aventure -> "L'Aventure
	}

	for _, expected := range expectedTransformations {
		if !strings.Contains(outputStr, expected) {
			t.Errorf("Smart quotes should handle nested quotes to contain '%s', but got: %s", expected, outputStr)
		}
	}
}

// TestSmartQuotesWithFrenchTypography tests smart quotes combined with French typography
func TestSmartQuotesWithFrenchTypography(t *testing.T) {
	// Create temporary directory
	tempDir, err := ioutil.TempDir("", "ankiprep_integration")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create input file combining quotes and French punctuation
	inputFile := filepath.Join(tempDir, "input.csv")
	csvContent := `question,answer
"Il a demandé : ""Comment ça va ?""","He asked: ""How are you?"""
"Elle répondit : ""Très bien, merci !""","She replied: ""Very well, thanks!"""`

	err = ioutil.WriteFile(inputFile, []byte(csvContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write input file: %v", err)
	}

	outputFile := filepath.Join(tempDir, "output.csv")

	// Execute ankiprep with both French typography and smart quotes enabled
	cmd := exec.Command("ankiprep", "--french", "--smart-quotes", "-o", outputFile, inputFile)
	_, err = cmd.CombinedOutput()

	if err != nil {
		t.Fatalf("ankiprep command failed: %v", err)
	}

	// Read output file
	outputContent, err := ioutil.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	outputStr := string(outputContent)

	// Check for both smart quotes and French typography
	expectedTransformations := []string{
		"demandé\u202F:",                   // French NNBSP before colon
		"\u201cComment ça va\u202F?\u201d", // Smart quotes + NNBSP before ?
		"répondit\u202F:",                  // French NNBSP before colon
		"\u201cTrès bien",                  // Smart opening quote
		"merci\u202F!\u201d",               // Smart closing quote + NNBSP before !
	}

	for _, expected := range expectedTransformations {
		if !strings.Contains(outputStr, expected) {
			t.Errorf("Smart quotes with French typography should contain '%s', but got: %s", expected, outputStr)
		}
	}
}

// TestSmartQuotesPreserveExisting tests that existing smart quotes are preserved
func TestSmartQuotesPreserveExisting(t *testing.T) {
	// Create temporary directory
	tempDir, err := ioutil.TempDir("", "ankiprep_integration")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create input file with mixed straight and smart quotes
	inputFile := filepath.Join(tempDir, "input.csv")
	csvContent := "front,back\n\"He said \u201cHello\u201d and \u2018goodbye\u2019\",\"Mixed quotes example\"\n\"She said \"\"Hi\"\" to me\",\"Straight quotes example\""

	err = ioutil.WriteFile(inputFile, []byte(csvContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write input file: %v", err)
	}

	outputFile := filepath.Join(tempDir, "output.csv")

	// Execute ankiprep
	cmd := exec.Command("ankiprep", "-o", outputFile, inputFile)
	_, err = cmd.CombinedOutput()

	if err != nil {
		t.Fatalf("ankiprep command failed: %v", err)
	}

	// Read output file
	outputContent, err := ioutil.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	outputStr := string(outputContent)

	// Check that existing smart quotes are preserved
	if !strings.Contains(outputStr, "\u201cHello\u201d") {
		t.Errorf("Existing smart quotes should be preserved, but got: %s", outputStr)
	}

	if !strings.Contains(outputStr, "\u2018goodbye\u2019") {
		t.Errorf("Existing smart single quotes should be preserved, but got: %s", outputStr)
	}

	// Check that straight quotes are converted
	if !strings.Contains(outputStr, "\u201cHi\u201d") {
		t.Errorf("Straight quotes should be converted to smart quotes, but got: %s", outputStr)
	}
}
