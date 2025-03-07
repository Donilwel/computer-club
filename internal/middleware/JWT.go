package middleware

import (
	"computer-club/internal/errors"
	"context"
	"github.com/golang-jwt/jwt/v4"
	"net/http"
	"os"
	"strings"
)

var jwtKey = []byte(os.Getenv("JWT_SECRET"))

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			WriteError(w, http.StatusUnauthorized, errors.ErrMissingToken.Error())
			return
		}

		// Убираем "Bearer " из токена
		tokenString = strings.TrimPrefix(tokenString, "Bearer ")

		claims := jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil || !token.Valid {
			WriteError(w, http.StatusUnauthorized, errors.ErrWrongToken.Error())
			return
		}

		// Проверяем user_id
		userID, ok := claims["user_id"].(float64)
		if !ok {
			WriteError(w, http.StatusUnauthorized, errors.ErrWrongIDFromJWT.Error())
			return
		}

		// Проверяем role
		role, ok := claims["role"].(string)
		if !ok {
			WriteError(w, http.StatusUnauthorized, errors.ErrWrongRoleFromJWT.Error())
			return
		}

		// Добавляем в контекст
		ctx := context.WithValue(r.Context(), "user_id", int64(userID))
		ctx = context.WithValue(ctx, "role", role)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
