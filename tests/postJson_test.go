package tests

import (
	"bytes"
	"encoding/json"
	sc "github.com/Mans397/eLibrary/serverConnection"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSendJsonHandler(t *testing.T) {
	tests := []struct {
		name         string
		method       string
		body         interface{}
		expectedCode int
	}{
		{
			name:         "Valid POST request",
			method:       http.MethodPost,
			body:         map[string]string{"message": "Hello, world!"},
			expectedCode: http.StatusOK,
		},
		{
			name:         "Invalid method (GET instead of POST)",
			method:       http.MethodGet,
			body:         nil,
			expectedCode: http.StatusMethodNotAllowed,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var bodyBytes []byte
			if tt.body != nil {
				var err error
				bodyBytes, err = json.Marshal(tt.body)
				if err != nil {
					t.Fatal(err)
				}
			}

			req := httptest.NewRequest(tt.method, "/post/json", bytes.NewBuffer(bodyBytes))
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()

			handler := http.HandlerFunc(sc.SendJsonHandler)
			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedCode {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tt.expectedCode)
			}
		})
	}
}
