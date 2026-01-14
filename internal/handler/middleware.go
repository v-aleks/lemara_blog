package handler

import (
	"context"
	"net/http"
	"strings"

	"lemara_blog/internal/utils"
)

type contextKey string

const (
    userIDKey contextKey = "user_id"
    emailKey  contextKey = "email"
)

func AuthMiddleware(jwtSecret string) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            authHeader := r.Header.Get("Authorization")
            if authHeader == "" {
                http.Error(w, "Authorization header required", http.StatusUnauthorized)
                return
            }

            parts := strings.Split(authHeader, " ")
            if len(parts) != 2 || parts[0] != "Bearer" {
                http.Error(w, "Invalid authorization format", http.StatusUnauthorized)
                return
            }

            token := parts[1]
            claims, err := utils.ParseToken(token, jwtSecret)
            if err != nil {
                http.Error(w, "Invalid token", http.StatusUnauthorized)
                return
            }

            // Add user info to context
            ctx := context.WithValue(r.Context(), userIDKey, claims.UserID)
            ctx = context.WithValue(ctx, emailKey, claims.Email)

            next.ServeHTTP(w, r.WithContext(ctx))
        })
    }
}

func GetUserIDFromContext(ctx context.Context) string {
    if val, ok := ctx.Value(userIDKey).(string); ok {
        return val
    }
    return ""
}

func GetEmailFromContext(ctx context.Context) string {
    if val, ok := ctx.Value(emailKey).(string); ok {
        return val
    }
    return ""
}
