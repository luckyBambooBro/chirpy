package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256, 
		jwt.RegisteredClaims{
			Issuer: "chirpy-access",
			IssuedAt: time.Now().UTC(), //wrap these in jwts wrapper
			ExpiresAt: time.Now().UTC() + expiresIn, //wrap these in jwts wrapper
			Subject: userID.String(),
		},
	),

}