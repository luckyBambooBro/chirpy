package auth

import (
	"encoding/hex"
	"errors"
	"net/http"
	"crypto/rand"
	"log"
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

func MakeRefreshToken() (string, error) {
	randBytes := make([]byte, 32)
	_, err := rand.Read(randBytes)
	if err != nil {
		log.Print("error creating randBytes for refresh token")
		return "", err
	}
	return hex.EncodeToString(randBytes), nil

}