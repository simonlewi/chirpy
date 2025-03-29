package auth_test

import (
	"net/http"
	"testing"

	"chirpy/internal/auth"
)

func TestGetBearerToken(t *testing.T) {
	tests := []struct {
		name          string
		authHeader    string
		expectedToken string
		expectError   bool
	}{
		{
			name:          "Valid bearer token",
			authHeader:    "Bearer abc123.def456.ghi789",
			expectedToken: "abc123.def456.ghi789",
			expectError:   false,
		},
		{
			name:          "Missing Authorization header",
			authHeader:    "",
			expectedToken: "",
			expectError:   true,
		},
		{
			name:          "Missing Bearer prefix",
			authHeader:    "abc123.def456.ghi789",
			expectedToken: "",
			expectError:   true,
		},
		{
			name:          "Empty token after Bearer",
			authHeader:    "Bearer ",
			expectedToken: "",
			expectError:   true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			headers := http.Header{}
			if tc.authHeader != "" {
				headers.Set("Authorization", tc.authHeader)
			}

			token, err := auth.GetBearerToken(headers)

			if tc.expectError && err == nil {
				t.Error("expected error but got none")
				return
			}
			if !tc.expectError && err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}
			if token != tc.expectedToken {
				t.Errorf("got token %q, want %q", token, tc.expectedToken)
			}
		})
	}
}
