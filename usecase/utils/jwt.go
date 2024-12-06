package utils

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

func GenerateTokens(userName string) (string, error) {
	// Generate access token
	accessToken, err := generateJWT(userName, time.Hour*24) // Expires in 24 hours
	if err != nil {
		return "", err
	}

	return accessToken, nil
}

func generateJWT(userName string, expirationTime time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"username": userName,
		"exp":     time.Now().Add(expirationTime).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte("your-secret-key"))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ValidateJWT(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte("your-secret-key"), nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}
