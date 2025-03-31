package auth

import (
	"fmt"
	"net/http"
	"strings"
)

func GetAPIKey(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", fmt.Errorf("authorization header not found")
	}

	if !strings.HasPrefix(authHeader, "ApiKey ") {
		return "", fmt.Errorf("authorization header must start with ApiKey")
	}

	// Strip "Autorization " prefix and any whitespace
	apiKey := strings.TrimPrefix(authHeader, "ApiKey ")
	apiKey = strings.TrimSpace(apiKey)

	if apiKey == "" {
		return "", fmt.Errorf("api key cannot be empty")
	}

	return apiKey, nil
}
