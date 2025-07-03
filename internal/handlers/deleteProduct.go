package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"shop-dashboard/internal/database"
	"time"

	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (h *Handler) SoftDeleteProduct(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "Missing product id", http.StatusBadRequest)
		return
	}

	productID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		http.Error(w, "Invalid product id", http.StatusBadRequest)
		return
	}

	collection := database.GetCollection("products")
	filter := bson.M{"_id": productID}
	update := bson.M{"$set": bson.M{"deleted": true}}

	res, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		http.Error(w, "Failed to delete product", http.StatusInternalServerError)
		return
	}
	if res.MatchedCount == 0 {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"message": "Product marked as deleted ad will not be avilable to customers/user but will remain avilable for republish",
		"_id":     id,
	})

}

func (h *Handler) RepublishProduct(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "Missing product id", http.StatusBadRequest)
		return
	}

	productID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		http.Error(w, "Invalid product id", http.StatusBadRequest)
		return
	}

	collection := database.GetCollection("products")
	filter := bson.M{"_id": productID}
	update := bson.M{"$set": bson.M{"deleted": false}}

	res, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		http.Error(w, "Failed to re-publish product", http.StatusInternalServerError)
		return
	}
	if res.MatchedCount == 0 {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"message": "Product marked as published and will be visable to customers/users",
		"_id":     id,
	})
}
