package auth
import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
)

func MakeRefreshToken() string {
	randBytes := make([]byte, 32)
	_, err := rand.Read(randBytes)
	if err != nil {
		log.Print("error creating randBytes for refresh token")
		return ""
	}
	encodedStr := hex.EncodeToString(randBytes)

}