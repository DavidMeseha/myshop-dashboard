package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"os"
	"shop-dashboard/internal/database"
	"shop-dashboard/internal/services"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RegisterVendorRequest struct {
	Name   string `json:"name"`
	SeName string `json:"seName"`
	Image  string `json:"image"`
}

func (h *Handler) RegisterVendorHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	token := r.Header.Get("Authorization")
	if token == "" {
		http.Error(w, "Missing token", http.StatusUnauthorized)
		return
	}
	token = strings.TrimPrefix(token, "Bearer ")

	user, err := services.CheckUserToken(token)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if user.IsVendor || !user.IsRegistered {
		http.Error(w, "Forbidden: already a vendor or not registered", http.StatusForbidden)
		return
	}

	var req RegisterVendorRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Name == "" || req.SeName == "" {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate image URL
	image := req.Image
	if image == "" {
		image = os.Getenv("CLIENT_SERVER") + "/images/no_image_placeholder.jpg"
	} else {
		parsedUrl, err := url.ParseRequestURI(image)
		if err != nil || parsedUrl.Scheme == "" || parsedUrl.Host == "" {
			http.Error(w, "Invalid image URL", http.StatusBadRequest)
			return
		}
	}

	count, err := database.CheckVendorSeName(ctx, req.SeName)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	if count > 0 {
		http.Error(w, "seName already exists", http.StatusConflict)
		return
	}

	userObjID, err := primitive.ObjectIDFromHex(user.ID)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusInternalServerError)
		return
	}

	// Insert new vendor
	vendorDoc := bson.M{
		"name":           req.Name,
		"seName":         req.SeName,
		"user":           userObjID,
		"imageUrl":       image,
		"productCount":   0,
		"followersCount": 0,
		"usersFollowed":  []string{},
		"createdAt":      time.Now(),
	}
	res, err := database.CreateVendor(ctx, vendorDoc)
	if err != nil {
		http.Error(w, "Failed to create vendor", http.StatusInternalServerError)
		return
	}

	err = database.UpdateUserIsVendorState(ctx, userObjID)
	if err != nil {
		http.Error(w, "Failed to update user as vendor", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"message":   "Vendor registered successfully",
		"vendor_id": res.InsertedID,
	})
}
