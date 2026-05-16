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
	testStrings := []struct{
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
	}

	var requestCases []*http.Request
	for _, testString := range testStrings {
		requestCases = append(
			requestCases, httptest.NewRequestWithContext(
				context.Background(), 
				http.MethodPost, 
				"/api/validate_path", 
				strings.NewReader(testString.inputBody),
			),
		) 
		//context.Background better here cos if we do context.WithTimeOut and one fails, the same ctx is passed around and the rest of the test case will fail
	}

	for i, requestCase := range requestCases {
		w := httptest.NewRecorder()

		handlerChirpsValidate(w, requestCase) 
		resp := w.Result()
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}

		bodyText, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("unable to read resp.body: %v", err)
		}
		if string(bodyText) != testStrings[i].expectedJSON {
			t.Errorf("w.Body does not match Request Body:\n w.Body: %v\n request body: %v\n", bodyText, testStrings[i])
		}
			
	}

}