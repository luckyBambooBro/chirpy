package auth

import (
	"errors"
	"net/http"
	"strings"
)

func GetBearerToken(headers http.Header) (string, error) {
	authText := headers.Get("Authorization") 
	if authText == "" {
		return "", errors.New("no authorization information provided by request header")
	}
	authTextList := strings.Fields(authText)
	if len(authTextList) != 2 {
		return "", errors.New("invalid authorization header length")
	}
	if strings.ToLower(authTextList[0]) != "bearer" {
		return "", errors.New("invalid authorization header format. Need <Bearer TOKEN_STRING>")
	}
	return authTextList[1], nil
}