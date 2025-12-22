package handlers

import (
	"context"
	"log"
	"net/http"
	"strings"

	"thainsbook/internal/auth"
)

type contextKey string

const UserIdKey contextKey = "userId"

func (a *Application) Authenticate(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			HandleUnauthorized(w, r)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			HandleUnauthorized(w, r)
			return
		}

		tokenString := parts[1]

		userId, err := auth.ValidateToken(tokenString, a.JWT)
		if err != nil {
			HandleUnauthorized(w, r)
			return
		}

		ctx := context.WithValue(r.Context(), UserIdKey, userId)

		next(w, r.WithContext(ctx))
	}
}

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Server Ping: %s %s", r.Method, r.URL.Path)

		next.ServeHTTP(w, r)
	})
}
