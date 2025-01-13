package helper

import (
	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken(secretKey string, payload map[string]interface{}) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	for key, value := range payload {
		claims[key] = value
	}

	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
