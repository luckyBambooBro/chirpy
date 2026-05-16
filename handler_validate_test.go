package main

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandlerChirpsValidate(t *testing.T) {

	//requestCases
	testCases := []struct{
		name string
		inputBody string
		expectedStatus int
		expectedJSON string
	}{
		{
			name: "standard string",
			inputBody: `{"body": "whatever trevor dawg"}`,
			expectedStatus: http.StatusOK,
			expectedJSON: `{"cleaned_body": "whatever trevor dawg"}`,
		},
			{
			name: "standard string",
			inputBody: `{"body": "what in the kerfuffle is that!"}`,
			expectedStatus: http.StatusOK,
			expectedJSON: `{"cleaned_body": "what in the **** is that!"}`,
		},
	}

	for _, tc := range testCases {
			r := httptest.NewRequestWithContext(
				context.Background(), 
				http.MethodPost, 
				"/api/validate_chirp",
				strings.NewReader(tc.inputBody),
				)
			w := httptest.NewRecorder()
		handlerChirpsValidate(w, r)
		resp := w.Result()
		if resp.StatusCode != http.StatusOK {
			t.Errorf("expected status: %v\nactual status: %v\n", http.StatusOK, resp.StatusCode)
		}
		bodyText, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("unable to read resp.body: %v", err)
		}
		resp.Body.Close()
		actualStr := strings.TrimSpace(string(bodyText))
		expectedStr := strings.TrimSpace(tc.expectedJSON)
		
		if actualStr != expectedStr {
			t.Errorf("JSON body mismatch\nactual: %v, expected: %v", actualStr, expectedStr)
		}
	}
}