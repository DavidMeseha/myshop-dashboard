package middleware

import (
	"context"
	"net/http"
	"shop-dashboard/internal/models"
	"shop-dashboard/internal/services"
	"strings"
)

type contextKey string

const userInfoKey contextKey = "userInfo"

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token == "" {
			http.Error(w, "Missing token", http.StatusUnauthorized)
			return
		}
		token = strings.TrimPrefix(token, "Bearer ")

		user, err := services.CheckVendorToken(token)
		if err != nil {
			http.Error(w, "Internal error", http.StatusInternalServerError)
			return
		}

		ctx := context.WithValue(r.Context(), userInfoKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Helper to get user info from context in handlers
func GetUserInfo(r *http.Request) (models.UserInfo, bool) {
	user, ok := r.Context().Value(userInfoKey).(models.UserInfo)
	return user, ok
}
