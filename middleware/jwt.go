package cMiddleware

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

func MakeJWT(userID, tokenSecret string, expiresIn time.Duration) (string, error) {
	claims := jwt.RegisteredClaims{
		Issuer:    "VidalTM",
		Subject:   userID,
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresIn)),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(tokenSecret)
}

func ValidateJWT(tokenString, tokenSecret string) (string, error) {
	token, err := jwt.ParseWithClaims(tokenString, jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Here we check if it was hashed with the same algorithm that we used, to prevent algorithm switching attacks
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Method.Alg())
		}

		return []byte(tokenSecret), nil
	})
	if err != nil {
		return "", err
	}

	return token.Claims.GetSubject()
}
