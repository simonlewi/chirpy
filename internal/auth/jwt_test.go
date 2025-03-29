package auth_test

import (
	"strings"
	"testing"
	"time"

	"chirpy/internal/auth"

	"github.com/google/uuid"
)

func TestJWTCreationAndValidation(t *testing.T) {
	// Test setup
	userID := uuid.New()
	tokenSecret := "your-test-secret"
	expiresIn := 24 * time.Hour

	// Test cases
	tests := []struct {
		name        string
		userID      uuid.UUID
		secret      string
		expiresIn   time.Duration
		shouldError bool
	}{
		{
			name:        "Valid token",
			userID:      userID,
			secret:      tokenSecret,
			expiresIn:   expiresIn,
			shouldError: false,
		},
		{
			name:        "Empty secret",
			userID:      userID,
			secret:      "",
			expiresIn:   expiresIn,
			shouldError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			token, err := auth.MakeJWT(tc.userID, tc.secret, tc.expiresIn)
			if tc.shouldError {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("failed to create token: %v", err)
			}

			gotUserID, err := auth.ValidateJWT(token, tc.secret)
			if err != nil {
				t.Fatalf("failed to validate token: %v", err)
			}

			if gotUserID != tc.userID {
				t.Errorf("got user ID %v, want %v", gotUserID, tc.userID)
			}
		})
	}
}

func TestJWTValidation_InvalidCases(t *testing.T) {
	userID := uuid.New()
	tokenSecret := "your-test-secret"
	expiresIn := 24 * time.Hour

	validToken, err := auth.MakeJWT(userID, tokenSecret, expiresIn)
	if err != nil {
		t.Fatalf("failed to create token: %v", err)
	}

	tests := []struct {
		name       string
		token      string
		secret     string
		wantErrMsg string
	}{
		{
			name:       "Invalid token format",
			token:      "not.a.jwt",
			secret:     tokenSecret,
			wantErrMsg: "token is malformed: could not JSON decode header: invalid character '\\u009e' looking for beginning of value",
		},
		{
			name:       "Wrong secret",
			token:      validToken,
			secret:     "wrong-secret",
			wantErrMsg: "token signature is invalid",
		},
		{
			name:       "Empty token",
			token:      "",
			secret:     tokenSecret,
			wantErrMsg: "token is malformed: token contains an invalid number of segments",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, err := auth.ValidateJWT(tc.token, tc.secret)
			if err == nil {
				t.Error("expected error, got nil")
				return
			}
			if !strings.Contains(err.Error(), tc.wantErrMsg) {
				t.Errorf("got error %q, want it to contain %q", err.Error(), tc.wantErrMsg)
			}
		})
	}
}
