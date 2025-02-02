package cMiddleware

import (
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type TokenParams struct {
	SignedString string
	UserID       string
	ExpiresAt    time.Time
}

func MakeJWT(userID, tokenSecret string, expiresIn time.Duration) (*TokenParams, error) {
	now := time.Now()
	expirationTime := now.Add(expiresIn)

	claims := jwt.RegisteredClaims{
		Issuer:    "VidalTM",
		Subject:   userID,
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(expirationTime),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedString, err := token.SignedString([]byte(tokenSecret))
	if err != nil {
		return nil, err
	}

	return &TokenParams{
		SignedString: signedString,
		UserID:       userID,
		ExpiresAt:    expirationTime,
	}, nil
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

func OptionalJWT(secretKey []byte) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader != "" {
				tokenString := strings.TrimPrefix(authHeader, "Bearer ")
				if tokenString != "" {

					token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
						return secretKey, nil
					})
					if err == nil && token.Valid {
						c.Set("user", token)
					}
				}
			}
			return next(c)
		}
	}
}
