package jwt

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var secretKey = []byte(os.Getenv("JWT_SECRET"))

func GenerateToken(username string, userID int64) (string, error) {
	claims := jwt.MapClaims{
		"username": username,
		"sub":      userID,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(secretKey)
}

func ParseToken(tokenStr string) (string, int64, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (any, error) {
		return secretKey, nil
	})

	if err != nil || !token.Valid {
		return "", -1, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", -1, err
	}

	username := claims["username"].(string)
	idFloat, ok := claims["sub"].(float64)
	if !ok {
		return "", -1, fmt.Errorf("invalid id claim: %+v", claims)
	}

	userID := int64(idFloat)

	return username, userID, nil
}
