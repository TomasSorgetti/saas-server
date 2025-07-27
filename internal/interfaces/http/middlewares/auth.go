package middlewares

import (
	"context"
	"luthierSaas/internal/infrastructure/security"
	"net/http"
	"strings"
)

type contextKey string

const UserIDKey contextKey = "user_id"

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("access_token")
		if err != nil || cookie.Value == "" {
			http.Error(w, "Acceso no autorizado (token faltante)", http.StatusUnauthorized)
			return
		}

		token := strings.TrimSpace(cookie.Value)

		userID, err := security.ValidateAccessToken(token)
		if err != nil {
			http.Error(w, "Token inválido o expirado", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), UserIDKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}