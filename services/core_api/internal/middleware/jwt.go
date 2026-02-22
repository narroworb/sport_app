package middleware

import (
	"context"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type ctxKey string

const (
	UsernameKey ctxKey = "username"
	UserIDKey   ctxKey = "sub"
)

var secretKey = []byte(os.Getenv("JWT_SECRET"))

func JWTAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "no authorization header", http.StatusUnauthorized)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "invalid authorization header", http.StatusUnauthorized)
			return
		}

		tokenStr := parts[1]

		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (any, error) {
			return secretKey, nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, "invalid token claims", http.StatusUnauthorized)
			return
		}

		username, ok := claims["username"].(string)
		if !ok {
			http.Error(w, "username not found in token", http.StatusUnauthorized)
			return
		}

		userIDFloat, ok := claims["sub"].(float64)
		if !ok {
			http.Error(w, "username not found in token", http.StatusUnauthorized)
			return
		}
		userID := int64(userIDFloat)

		ctx := context.WithValue(r.Context(), UsernameKey, username)
		ctx = context.WithValue(ctx, UserIDKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
