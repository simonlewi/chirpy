package auth_test

import (
	"chirpy/internal/auth"
	"testing"
)

func TestHashPassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
		wantErr  bool
	}{
		{
			name:     "valid password",
			password: "password123",
			wantErr:  false,
		},
		{
			name:     "empty password",
			password: "",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hash, err := auth.HashPassword(tt.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("HashPassword() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && hash == "" {
				t.Errorf("HashPassword() returned empty hash")
			}
		})
	}
}

func TestCheckPasswordHash(t *testing.T) {
	password := "password123"
	hash, err := auth.HashPassword(password)
	if err != nil {
		t.Fatalf("Failed to create hash for testing: %v", err)
	}

	tests := []struct {
		name        string
		hash        string
		password    string
		shouldMatch bool
	}{
		{
			name:        "matching password",
			hash:        hash,
			password:    password,
			shouldMatch: true,
		},
		{
			name:        "wrong password",
			hash:        hash,
			password:    "wrongpassword",
			shouldMatch: false,
		},
		{
			name:        "invalid hash",
			hash:        "invalid_hash",
			password:    password,
			shouldMatch: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := auth.CheckPasswordHash(tt.hash, tt.password)
			if tt.shouldMatch && err != nil {
				t.Errorf("CheckPasswordHash() should match but got error: %v", err)
			}
			if !tt.shouldMatch && err == nil {
				t.Errorf("CheckPasswordHash() should NOT match but did")
			}
		})
	}
}
