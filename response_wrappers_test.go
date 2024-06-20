package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRespondWithError(t *testing.T) {
	t.Run("5XX error", func(t *testing.T) {
		w := httptest.NewRecorder()
		respondWithError(w, 500, "Internal Server Error")

		if w.Code != 500 {
			t.Errorf("Expected status code 500, but got %d", w.Code)
		}

		expectedBody := `{"error":"Internal Server Error"}`
		if w.Body.String() != expectedBody {
			t.Errorf("Expected body %s, but got %s", expectedBody, w.Body.String())
		}
	})

	t.Run("4XX error", func(t *testing.T) {
		w := httptest.NewRecorder()
		respondWithError(w, 404, "Not Found")

		if w.Code != 404 {
			t.Errorf("Expected status code 404, but got %d", w.Code)
		}

		expectedBody := `{"error":"Not Found"}`
		if w.Body.String() != expectedBody {
			t.Errorf("Expected body %s, but got %s", expectedBody, w.Body.String())
		}
	})
}

func TestRespondWithJSON(t *testing.T) {
	// Create a new HTTP request
	_, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a new HTTP response recorder
	recorder := httptest.NewRecorder()

	// Call the respondWithJSON function with the recorder and a sample payload
	respondWithJSON(recorder, http.StatusOK, map[string]string{"message": "success"})

	// Check the status code of the response
	if recorder.Code != http.StatusOK {
		t.Errorf("Expected status code %d, but got %d", http.StatusOK, recorder.Code)
	}

	// Check the content type of the response
	expectedContentType := "application/json"
	if recorder.Header().Get("Content-Type") != expectedContentType {
		t.Errorf("Expected content type %s, but got %s", expectedContentType, recorder.Header().Get("Content-Type"))
	}

	// Check the content of the response
	expectedBody := `{"message":"success"}`
	if recorder.Body.String() != expectedBody {
		t.Errorf("Expected body %s, but got %s", expectedBody, recorder.Body.String())
	}
}
