package auth

import (
	"errors"
	"net/http"
	"strings"
)

func GetAPIKey(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("api header is empty")
	}

	parts := strings.Fields(authHeader)
	if len(parts) != 2 || parts[0] != "ApiKey" {
		return "", errors.New("invalid authorisation header")
	}

	return parts[1], nil
}