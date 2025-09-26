package contract

import "testing"

// TestCLI_IssuesDocumented documents CLI issues found in typography processor
func TestCLI_IssuesDocumented(t *testing.T) {
	t.Log("CLI issues documented based on typography processor contract tests:")
	t.Log("1. Duplicate NNBSP insertion when NNBSP already exists")
	t.Log("2. Lack of idempotency - same input produces different outputs")
	t.Log("3. Regular spaces not replaced with NNBSP in quotes")
}
