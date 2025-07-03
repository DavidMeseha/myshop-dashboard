package middleware

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserInfo struct {
	ID             primitive.ObjectID `json:"_id"`
	User           primitive.ObjectID `json:"user"`
	Name           string             `json:"name"`
	SeName         string             `json:"seName"`
	ImageUrl       string             `json:"imageUrl"`
	ProductCount   int                `json:"productCount"`
	FollowersCount int                `json:"followersCount"`
	CreatedAt      string             `json:"createdAt"`
}

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

		url := os.Getenv("CLIENT_SERVER")
		req, err := http.NewRequest("GET", url+"/api/v2/auth/vendor", nil)
		if err != nil {
			http.Error(w, "Internal error", http.StatusInternalServerError)
			return
		}
		req.Header.Set("Authorization", "Bearer "+token)

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil || resp.StatusCode != http.StatusOK {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		defer resp.Body.Close()

		var user UserInfo
		if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
			log.Print(err)
			http.Error(w, "Invalid user data", http.StatusUnauthorized)
			return
		}

		// Store user info in context for handlers to use
		ctx := context.WithValue(r.Context(), userInfoKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Helper to get user info from context in handlers
func GetUserInfo(r *http.Request) (UserInfo, bool) {
	user, ok := r.Context().Value(userInfoKey).(UserInfo)
	return user, ok
}
