package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const issuerName string = "chirpy-access" 

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256, 
		jwt.RegisteredClaims{
			Issuer: issuerName,
			IssuedAt: jwt.NewNumericDate(time.Now().UTC()),
			ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(expiresIn)),
			Subject: userID.String(),
		},
	)
	return token.SignedString([]byte(tokenSecret)) //this returns a string which is the entire JWT
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	claims := &jwt.RegisteredClaims{}
	//parse token into claims
	token, err := jwt.ParseWithClaims(
		tokenString, 
		claims, 
		jwt.Keyfunc(func (token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method")
			}
			return []byte(tokenSecret), nil
		}),
	)
	if err != nil {
		return uuid.Nil, fmt.Errorf("error parsing token: %w", err)
	}

	if !token.Valid {
		return uuid.Nil, fmt.Errorf("invalid token")
	}
	issuer, err := token.Claims.GetIssuer()
	if err != nil {
		return uuid.Nil, fmt.Errorf("unable to obtain issuer: %v", err)
	}
	if issuer != issuerName {
		return uuid.Nil, fmt.Errorf("invalid issuer")
	}
	userID, err := uuid.Parse(claims.Subject)
	if err != nil {
		return uuid.Nil, fmt.Errorf("error parsing userID: %w", err)
	}
	return userID, nil
}