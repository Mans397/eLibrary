package tests

import (
	"bytes"
	"encoding/json"
	sc "github.com/Mans397/eLibrary/serverConnection"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetJsonHandler(t *testing.T) {
	t.Run("Valid GET request", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/data/json", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()

		handler := http.HandlerFunc(sc.DataJsonHandler)
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("GET /data/json returned wrong status code: got %v want %v", status, http.StatusOK)
		}

		var response []sc.Request
		err = json.NewDecoder(rr.Body).Decode(&response)
		if err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		if len(response) == 0 {
			t.Error("GET /data/json returned empty response")
		}
	})
}

func TestPostJsonHandler(t *testing.T) {
	t.Run("Valid POST request", func(t *testing.T) {
		message := "Hello, World!"
		data := map[string]string{"message": message}
		jsonData, err := json.Marshal(data)
		if err != nil {
			t.Fatal(err)
		}

		req, err := http.NewRequest("POST", "/post/json", bytes.NewBuffer(jsonData))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()

		handler := http.HandlerFunc(sc.SendJsonHandler)
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("POST /post/json returned wrong status code: got %v want %v", status, http.StatusOK)
		}

		var response sc.Response
		err = json.NewDecoder(rr.Body).Decode(&response)
		if err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		if response.Status != "Success" {
			t.Errorf("Expected 'Success' status, got %v", response.Status)
		}
		if response.Message != message {
			t.Errorf("Expected message '%v', got %v", message, response.Message)
		}
	})

	t.Run("Invalid POST request without message", func(t *testing.T) {
		data := map[string]interface{}{}
		jsonData, err := json.Marshal(data)
		if err != nil {
			t.Fatal(err)
		}

		req, err := http.NewRequest("POST", "/post/json", bytes.NewBuffer(jsonData))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()

		handler := http.HandlerFunc(sc.SendJsonHandler)
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusBadRequest {
			t.Errorf("POST /post/json returned wrong status code: got %v want %v", status, http.StatusBadRequest)
		}

		var response sc.Response
		err = json.NewDecoder(rr.Body).Decode(&response)
		if err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		if response.Status != "Fail" {
			t.Errorf("Expected 'Fail' status, got %v", response.Status)
		}
		if response.Message != `"message" field is not found` {
			t.Errorf("Expected error message '\"message\" field is not found', got %v", response.Message)
		}
	})
}
