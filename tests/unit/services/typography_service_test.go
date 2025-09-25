package services_test

import (
	"bytes"
	"log"
	"testing"

	"ankiprep/internal/services"
)

// TestProcessFrenchTextRequest_Contract verifies the ProcessFrenchTextRequest contract
func TestProcessFrenchTextRequest_Contract(t *testing.T) {
	tests := []struct {
		name    string
		request services.ProcessFrenchTextRequest
		wantErr bool
	}{
		{
			name: "valid request with French and SmartQuotes",
			request: services.ProcessFrenchTextRequest{
				Text:        "Bonjour : comment allez-vous ?",
				French:      true,
				SmartQuotes: true,
			},
			wantErr: false,
		},
		{
			name: "valid request with French only",
			request: services.ProcessFrenchTextRequest{
				Text:   "Bonjour : comment allez-vous ?",
				French: true,
			},
			wantErr: false,
		},
		{
			name: "valid request with SmartQuotes only",
			request: services.ProcessFrenchTextRequest{
				Text:        "Hello \"world\"",
				SmartQuotes: true,
			},
			wantErr: false,
		},
		{
			name: "invalid request with empty text",
			request: services.ProcessFrenchTextRequest{
				Text:   "",
				French: true,
			},
			wantErr: true,
		},
		{
			name: "invalid request with neither French nor SmartQuotes",
			request: services.ProcessFrenchTextRequest{
				Text:        "Some text",
				French:      false,
				SmartQuotes: false,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.request.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("ProcessFrenchTextRequest.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestProcessFrenchText_ClozeExceptionContract verifies the cloze exception contract
func TestProcessFrenchText_ClozeExceptionContract(t *testing.T) {
	service := services.NewTypographyService()

	var logBuffer bytes.Buffer
	logger := log.New(&logBuffer, "", 0)

	tests := []struct {
		name               string
		input              string
		expectedOutput     string
		expectedClozeCount int
	}{
		{
			name:               "colon outside cloze block",
			input:              "Question : What is the capital of France?",
			expectedOutput:     "Question\u00A0: What is the capital of France?",
			expectedClozeCount: 0,
		},
		{
			name:               "colon inside cloze block should not be modified",
			input:              "The capital of France is {{c1::Paris : the city of light}}",
			expectedOutput:     "The capital of France is {{c1::Paris : the city of light}}",
			expectedClozeCount: 1,
		},
		{
			name:               "mixed colons - inside and outside cloze",
			input:              "Question : The capital is {{c1::Paris : beautiful city}}",
			expectedOutput:     "Question\u00A0: The capital is {{c1::Paris : beautiful city}}",
			expectedClozeCount: 1,
		},
		{
			name:               "multiple cloze blocks with colons",
			input:              "Cities : {{c1::Paris : France}} and {{c2::Rome : Italy}}",
			expectedOutput:     "Cities\u00A0: {{c1::Paris : France}} and {{c2::Rome : Italy}}",
			expectedClozeCount: 2,
		},
		{
			name:               "nested cloze patterns should not break parsing",
			input:              "Test : {{c1::value with {{text}} inside : result}}",
			expectedOutput:     "Test\u00A0: {{c1::value with {{text}} inside : result}}",
			expectedClozeCount: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := services.ProcessFrenchTextRequest{
				Text:   tt.input,
				French: true,
				Logger: logger,
			}

			response, err := service.ProcessFrenchText(request)
			if err != nil {
				t.Fatalf("ProcessFrenchText() error = %v", err)
			}

			if response.ProcessedText != tt.expectedOutput {
				t.Errorf("ProcessFrenchText() output = %q, want %q",
					response.ProcessedText, tt.expectedOutput)
			}

			if response.ClozeCount != tt.expectedClozeCount {
				t.Errorf("ProcessFrenchText() cloze count = %d, want %d",
					response.ClozeCount, tt.expectedClozeCount)
			}
		})
	}
}

// TestProcessFrenchText_IdempotencyContract verifies idempotency requirement
func TestProcessFrenchText_IdempotencyContract(t *testing.T) {
	service := services.NewTypographyService()

	input := "Question : What is {{c1::Paris : capital}} of France?"
	request := services.ProcessFrenchTextRequest{
		Text:   input,
		French: true,
	}

	// First processing
	response1, err := service.ProcessFrenchText(request)
	if err != nil {
		t.Fatalf("First ProcessFrenchText() error = %v", err)
	}

	// Second processing with the result of the first
	request.Text = response1.ProcessedText
	response2, err := service.ProcessFrenchText(request)
	if err != nil {
		t.Fatalf("Second ProcessFrenchText() error = %v", err)
	}

	// Results should be identical (idempotent)
	if response1.ProcessedText != response2.ProcessedText {
		t.Errorf("ProcessFrenchText() is not idempotent:\nFirst:  %q\nSecond: %q",
			response1.ProcessedText, response2.ProcessedText)
	}
}
