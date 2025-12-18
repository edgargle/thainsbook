package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func CreateToken(userId, secret string) (string, string, error) {
	expiryTime := time.Now().Add(time.Hour * 24)
	claims := jwt.MapClaims{
		"userId": userId,
		"exp":    expiryTime.Unix(),
		"iat":    time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", "", err
	}

	expiryString := expiryTime.Format(time.RFC3339)

	return tokenString, expiryString, nil
}

func ValidateToken(tokenString string, secret string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userId, _ := claims["userId"].(string)
		if userId == "" {
			return "", errors.New("token is missing required claims")
		}
		return userId, nil
	}

	return "", errors.New("invalid token")
}
