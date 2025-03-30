package auth

import (
	"crypto/rand"
	"encoding/hex"
)

// MakeRefreshToken generates a random 256-bit (32-byte) hex-encoded string
func MakeRefreshToken() (string, error) {
	// Create a byte slice of 32 bytes (256 bits)
	randomBytes := make([]byte, 32)

	// Fill the slice with random data using crypto/rand
	if _, err := rand.Read(randomBytes); err != nil {
		return "", err
	}

	// Convert the random bytes to a hex string
	return hex.EncodeToString(randomBytes), nil
}
