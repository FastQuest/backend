package auth

import (
	"context"
	"net/http"
	"os"
	"strings"

	"flashquest/pkg/apiresp"
	jwtsec "flashquest/pkg/security/jwt"
)

const (
	ContextKeyUserID = "user_id"
	ContextKeyRole   = "role"
)

func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			apiresp.WriteError(w, http.StatusUnauthorized, "MISSING_TOKEN", "Authorization header is required")
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			apiresp.WriteError(w, http.StatusUnauthorized, "INVALID_TOKEN_FORMAT", "Authorization header must be Bearer <token>")
			return
		}

		token := parts[1]
		publicKeyPEM := os.Getenv("JWT_PUBLIC_KEY")
		claims, err := jwtsec.ParseAccessToken(token, publicKeyPEM)
		if err != nil {
			apiresp.WriteError(w, http.StatusUnauthorized, "INVALID_TOKEN", "Invalid or expired token")
			return
		}

		ctx := context.WithValue(r.Context(), ContextKeyUserID, claims.UserID)
		ctx = context.WithValue(ctx, ContextKeyRole, claims.Role)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
