package middleware

import (
	"context"
	"github.com/golang-jwt/jwt/v4"
	"net/http"
	"os"
)

var jwtKey = []byte(os.Getenv("JWT_SECRET"))

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			WriteError(w, http.StatusUnauthorized, "Missing token")
			return
		}

		claims := jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil || !token.Valid {
			WriteError(w, http.StatusUnauthorized, "Invalid token")
			return
		}

		// Добавляем user_id и role в контекст запроса
		ctx := context.WithValue(r.Context(), "user_id", int64(claims["user_id"].(float64)))
		ctx = context.WithValue(ctx, "role", claims["role"].(string))

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
