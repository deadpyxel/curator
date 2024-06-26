package auth

import (
	"errors"
	"net/http"
	"testing"
)

func TestGetApiKey(t *testing.T) {
	tests := []struct {
		name          string
		headers       http.Header
		expectedKey   string
		expectedError error
	}{
		{
			name: "When it is a valid Authorization header returns key and no error",
			headers: http.Header{
				"Authorization": []string{"ApiKey my-api-key"},
			},
			expectedKey:   "my-api-key",
			expectedError: nil,
		},
		{
			name: "When there is no Authorization header returns no key and matching error",
			headers: http.Header{
				"Authorization": []string{""},
			},
			expectedKey:   "",
			expectedError: errors.New("No authentication information found."),
		},
		{
			name: "When is a malformed Authorization header returns no key and matching error",
			headers: http.Header{
				"Authorization": []string{"token"},
			},
			expectedKey:   "",
			expectedError: errors.New("Malformed auth header"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			key, err := GetApiKey(tt.headers)

			if key != tt.expectedKey {
				t.Errorf("Expected key: %s, got: %s", tt.expectedKey, key)
			}

			if (tt.expectedError != nil || err != nil) && tt.expectedError.Error() != err.Error() {
				t.Errorf("Expected error %v, got %v", tt.expectedError, err)
			}
		})
	}
}
