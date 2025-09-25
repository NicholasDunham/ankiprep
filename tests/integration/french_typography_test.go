package integration

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// TestFrenchTypographyPunctuation tests NNBSP insertion before French punctuation
func TestFrenchTypographyPunctuation(t *testing.T) {
	// Create temporary directory
	tempDir, err := ioutil.TempDir("", "ankiprep_integration")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)
	
	// Create input file with French text needing typography fixes
	inputFile := filepath.Join(tempDir, "input.csv")
	csvContent := `question,answer
"Comment allez-vous ?","How are you?"
"Voulez-vous du café ?","Would you like coffee?"
"Il a dit : « Bonjour ! »","He said: 'Hello!'"
"Est-ce que vous parlez français ; ou anglais ?","Do you speak French; or English?"`
	
	err = ioutil.WriteFile(inputFile, []byte(csvContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write input file: %v", err)
	}
	
	outputFile := filepath.Join(tempDir, "output.csv")
	
	// Execute ankiprep with French typography flag
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
	
	// Check for NNBSP (U+202F) insertion before French punctuation
	expectedTransformations := []string{
		"vous\u202F?",      // NNBSP before question mark
		"café\u202F?",      // NNBSP before question mark  
		"dit\u202F:",       // NNBSP before colon
		"«\u202FBonjour",   // NNBSP after opening guillemet
		"Bonjour\u202F!",   // NNBSP before exclamation mark
		"français\u202F;",  // NNBSP before semicolon
		"anglais\u202F?",   // NNBSP before question mark
	}
	
	for _, expected := range expectedTransformations {
		if !strings.Contains(outputStr, expected) {
			t.Errorf("French typography should transform text to contain '%s', but got: %s", expected, outputStr)
		}
	}
}

// TestFrenchTypographyGuillemets tests proper spacing with French quotation marks
func TestFrenchTypographyGuillemets(t *testing.T) {
	// Create temporary directory
	tempDir, err := ioutil.TempDir("", "ankiprep_integration")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)
	
	// Create input file with French quotes
	inputFile := filepath.Join(tempDir, "input.csv")
	csvContent := `front,back
"Il a dit : « Bonjour » à Marie","He said 'Hello' to Marie"
"« Comment ça va ? » demanda-t-il","'How are you?' he asked"
"Elle répondit : « Très bien, merci ! »","She replied: 'Very well, thanks!'"
"« Est-ce que tu viens ? » « Oui, j'arrive ! »","'Are you coming?' 'Yes, I'm coming!'"`
	
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
	
	// Check for proper guillemet spacing
	expectedPatterns := []string{
		"«\u202FBonjour\u202F»",    // NNBSP inside guillemets
		"«\u202FComment",           // NNBSP after opening guillemet
		"va\u202F?\u202F»",         // NNBSP before closing guillemet and before ?
		"«\u202FEst-ce",            // NNBSP after opening guillemet
		"viens\u202F?\u202F»",      // NNBSP before ? and closing guillemet
		"«\u202FOui",               // NNBSP after opening guillemet
		"arrive\u202F!\u202F»",     // NNBSP before ! and closing guillemet
	}
	
	for _, expected := range expectedPatterns {
		if !strings.Contains(outputStr, expected) {
			t.Errorf("French typography should contain guillemet pattern '%s', but got: %s", expected, outputStr)
		}
	}
}

// TestFrenchTypographyMixedPunctuation tests complex sentences with multiple punctuation
func TestFrenchTypographyMixedPunctuation(t *testing.T) {
	// Create temporary directory
	tempDir, err := ioutil.TempDir("", "ankiprep_integration")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)
	
	// Create input file with complex French punctuation
	inputFile := filepath.Join(tempDir, "input.csv")
	csvContent := `question,answer,context
"Vraiment ? Vous pensez : « C'est possible ! » ?","Really? You think: 'It's possible!' ?",complex
"Est-ce que Paul ; Marie ; et Jean viennent ?","Are Paul; Marie; and Jean coming?",list
"Attention ! Danger : zone interdite !","Warning! Danger: forbidden zone!",warning`
	
	err = ioutil.WriteFile(inputFile, []byte(csvContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write input file: %v", err)
	}
	
	outputFile := filepath.Join(tempDir, "output.csv")
	
	// Execute ankiprep with French typography
	cmd := exec.Command("ankiprep", "-f", "-o", outputFile, inputFile)
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
	
	// Check for proper spacing with all punctuation types
	expectedTransformations := []string{
		"Vraiment\u202F?",          // NNBSP before ?
		"pensez\u202F:",            // NNBSP before :
		"«\u202FC'est",             // NNBSP after opening guillemet
		"possible\u202F!",          // NNBSP before !
		"Paul\u202F;",              // NNBSP before ;
		"Marie\u202F;",             // NNBSP before ;
		"viennent\u202F?",          // NNBSP before ?
		"Attention\u202F!",         // NNBSP before !
		"Danger\u202F:",            // NNBSP before :
		"interdite\u202F!",         // NNBSP before !
	}
	
	for _, expected := range expectedTransformations {
		if !strings.Contains(outputStr, expected) {
			t.Errorf("French typography should contain transformation '%s', but got: %s", expected, outputStr)
		}
	}
}

// TestFrenchTypographyPreserveExisting tests that existing proper spacing is preserved
func TestFrenchTypographyPreserveExisting(t *testing.T) {
	// Create temporary directory
	tempDir, err := ioutil.TempDir("", "ankiprep_integration")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)
	
	// Create input file with already properly formatted French text
	inputFile := filepath.Join(tempDir, "input.csv")
	// Note: Using NNBSP (U+202F) in the input
	csvContent := "question,answer\n\"Comment allez-vous\u202F?\",\"How are you?\"\n\"Il dit\u202F: « Bonjour\u202F! »\",\"He says: 'Hello!'\""
	
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
	
	// Check that existing proper formatting is preserved
	expectedPreserved := []string{
		"vous\u202F?",      // Already correct NNBSP before ?
		"dit\u202F:",       // Already correct NNBSP before :
		"Bonjour\u202F!",   // Already correct NNBSP before !
	}
	
	for _, expected := range expectedPreserved {
		if !strings.Contains(outputStr, expected) {
			t.Errorf("French typography should preserve existing correct formatting '%s', but got: %s", expected, outputStr)
		}
	}
	
	// Should not introduce double NNBSP
	doubleNNBSP := "\u202F\u202F"
	if strings.Contains(outputStr, doubleNNBSP) {
		t.Errorf("French typography should not create double NNBSP, but got: %s", outputStr)
	}
}