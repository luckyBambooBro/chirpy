package auth
import (
	"crypto/rand"
	"encoding/hex"
	"log"
)

func MakeRefreshToken() (string, error) {
	randBytes := make([]byte, 32)
	_, err := rand.Read(randBytes)
	if err != nil {
		log.Print("error creating randBytes for refresh token")
		return "", err
	}
	return hex.EncodeToString(randBytes), nil

}