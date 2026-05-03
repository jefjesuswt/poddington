package whttp

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestErrorJSON(t *testing.T) {

	tests := []struct {
		name         string
		message      string
		expectedBody string
		status       int
	}{
		{
			name:         "Error 404 Not Found",
			message:      "node not found",
			expectedBody: `{"error":"node not found"}`,
			status:       http.StatusNotFound,
		},
		{
			name:         "Error 401 Unauthorized",
			status:       http.StatusUnauthorized,
			message:      "invalid token",
			expectedBody: `{"error":"invalid token"}`,
		},
		{
			name:         "Error 500 Internal",
			status:       http.StatusInternalServerError,
			message:      "database crashed",
			expectedBody: `{"error":"database crashed"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()

			ErrorJSON(recorder, tt.status, tt.message)

			if recorder.Code != tt.status {
				t.Errorf("expected code: %d, got: %d", tt.status, recorder.Code)
			}

			contentType := recorder.Header().Get("Content-Type")
			if contentType != "application/json" {
				t.Errorf("expected Content-Type: application/json, got: %s", contentType)
			}

			actualBody := strings.TrimSpace(recorder.Body.String())
			if actualBody != tt.expectedBody {
				t.Errorf("expected body: %s, got: %s", tt.expectedBody, actualBody)
			}
		})
	}
}
