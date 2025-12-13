package util

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken(id string, extraClaims ...map[string]any) (string, error) {

	secret := os.Getenv("JWT_SECRET")

	token_expiration := os.Getenv("TOKEN_EXPIRATION")
	if token_expiration == "" {
		token_expiration = "1h"
	}

	token_expiration_duration, err := time.ParseDuration(token_expiration)
	if err != nil {
		return "", err
	}

	jwt_claims := jwt.MapClaims{
		"id":  id,
		"exp": time.Now().Add(token_expiration_duration).Unix(),
		"iat": time.Now().Unix(),
	}

	// Merge extra claims if provided
	if len(extraClaims) > 0 {
		for k, v := range extraClaims[0] {
			jwt_claims[k] = v
		}
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt_claims)

	return token.SignedString([]byte(secret))
}
