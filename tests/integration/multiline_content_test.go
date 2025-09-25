package integration

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// TestMultilineContentBasic tests handling of multiline content in CSV fields
func TestMultilineContentBasic(t *testing.T) {
	// Create temporary directory
	tempDir, err := ioutil.TempDir("", "ankiprep_integration")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)
	
	// Create input file with multiline content
	inputFile := filepath.Join(tempDir, "input.csv")
	csvContent := `front,back
"Line 1
Line 2
Line 3","Ligne 1
Ligne 2
Ligne 3"
"Single line","Ligne simple"`
	
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
	
	// Check that multiline content is converted to HTML line breaks
	expectedTransformations := []string{
		"Line 1<br>Line 2<br>Line 3",     // Newlines converted to <br>
		"Ligne 1<br>Ligne 2<br>Ligne 3",  // Newlines converted to <br>
		"Single line",                     // Single line preserved
		"Ligne simple",                    // Single line preserved
	}
	
	for _, expected := range expectedTransformations {
		if !strings.Contains(outputStr, expected) {
			t.Errorf("Multiline content should be converted to contain '%s', but got: %s", expected, outputStr)
		}
	}
	
	// Verify Anki HTML flag is set
	if !strings.Contains(outputStr, "#html:true") {
		t.Errorf("Output should contain #html:true flag for multiline content")
	}
}

// TestMultilineContentWithFormatting tests multiline content with existing HTML
func TestMultilineContentWithFormatting(t *testing.T) {
	// Create temporary directory
	tempDir, err := ioutil.TempDir("", "ankiprep_integration")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)
	
	// Create input file with multiline content and HTML formatting
	inputFile := filepath.Join(tempDir, "input.csv")
	csvContent := `front,back
"<b>Bold text</b>
<i>Italic text</i>
Normal text","<b>Texte gras</b>
<i>Texte italique</i>
Texte normal"
"HTML with breaks<br>Already formatted","HTML avec breaks<br>Déjà formaté"`
	
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
	
	// Check that existing HTML is preserved and newlines are converted
	expectedTransformations := []string{
		"<b>Bold text</b><br><i>Italic text</i><br>Normal text",
		"<b>Texte gras</b><br><i>Texte italique</i><br>Texte normal",
		"HTML with breaks<br>Already formatted",  // Existing <br> preserved
		"HTML avec breaks<br>Déjà formaté",      // Existing <br> preserved
	}
	
	for _, expected := range expectedTransformations {
		if !strings.Contains(outputStr, expected) {
			t.Errorf("HTML formatting should be preserved and newlines converted: '%s', but got: %s", expected, outputStr)
		}
	}
}

// TestMultilineContentWithQuotes tests multiline content with quote handling
func TestMultilineContentWithQuotes(t *testing.T) {
	// Create temporary directory
	tempDir, err := ioutil.TempDir("", "ankiprep_integration")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)
	
	// Create input file with multiline content containing quotes
	inputFile := filepath.Join(tempDir, "input.csv")
	csvContent := `question,answer
"He said: ""Hello there!""
How are you doing?
See you later!","Il a dit: ""Salut !""
Comment ça va ?
À bientôt !"
"Simple ""quoted"" text","Texte ""entre guillemets"" simple"`
	
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
	
	// Check that quotes are converted to smart quotes and newlines to <br>
	expectedTransformations := []string{
		"\u201cHello there!\u201d<br>How are you doing?<br>See you later!",
		"\u201cSalut !\u201d<br>Comment ça va ?<br>À bientôt !",
		"\u201cquoted\u201d text",
		"\u201centre guillemets\u201d simple",
	}
	
	for _, expected := range expectedTransformations {
		if !strings.Contains(outputStr, expected) {
			t.Errorf("Multiline content with quotes should contain '%s', but got: %s", expected, outputStr)
		}
	}
}

// TestMultilineContentWithFrench tests multiline French content with typography
func TestMultilineContentWithFrench(t *testing.T) {
	// Create temporary directory
	tempDir, err := ioutil.TempDir("", "ankiprep_integration")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)
	
	// Create input file with multiline French content
	inputFile := filepath.Join(tempDir, "input.csv")
	csvContent := `question,answer
"Comment allez-vous ?
Très bien, merci !
Et vous ?","How are you?
Very well, thanks!
And you?"
"Êtes-vous sûr : oui ou non ?
Répondez s'il vous plaît !","Are you sure: yes or no?
Please answer!"`
	
	err = ioutil.WriteFile(inputFile, []byte(csvContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write input file: %v", err)
	}
	
	outputFile := filepath.Join(tempDir, "output.csv")
	
	// Execute ankiprep with French typography
	cmd := exec.Command("ankiprep", "--french", "-o", outputFile, inputFile)
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
	
	// Check that French typography is applied AND newlines are converted to <br>
	expectedTransformations := []string{
		"vous\u202F?<br>Très bien",        // NNBSP before ? and <br> conversion
		"merci\u202F!<br>Et vous\u202F?",  // NNBSP before ! and ? with <br>
		"sûr\u202F:",                      // NNBSP before :
		"non\u202F?<br>Répondez",          // NNBSP before ? with <br>
		"plaît\u202F!",                    // NNBSP before !
	}
	
	for _, expected := range expectedTransformations {
		if !strings.Contains(outputStr, expected) {
			t.Errorf("French multiline content should contain '%s', but got: %s", expected, outputStr)
		}
	}
}

// TestMultilineContentPreserveSpacing tests that leading/trailing spaces are preserved
func TestMultilineContentPreserveSpacing(t *testing.T) {
	// Create temporary directory
	tempDir, err := ioutil.TempDir("", "ankiprep_integration")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)
	
	// Create input file with multiline content having various spacing
	inputFile := filepath.Join(tempDir, "input.csv")
	csvContent := `front,back
"  Indented line
Normal line
    More indented","  Ligne indentée
Ligne normale
    Plus indentée"
"Line with trailing spaces  
Another line","Ligne avec espaces finaux  
Autre ligne"`
	
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
	
	// Check that spacing is preserved in multiline conversion
	expectedTransformations := []string{
		"  Indented line<br>Normal line<br>    More indented",
		"  Ligne indentée<br>Ligne normale<br>    Plus indentée",
		"trailing spaces  <br>Another line",
		"espaces finaux  <br>Autre ligne",
	}
	
	for _, expected := range expectedTransformations {
		if !strings.Contains(outputStr, expected) {
			t.Errorf("Spacing should be preserved in multiline content '%s', but got: %s", expected, outputStr)
		}
	}
}