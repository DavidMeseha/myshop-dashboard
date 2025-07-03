package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"shop-dashboard/internal/database"
	"shop-dashboard/internal/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (h *Handler) FindVendors(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	query := r.URL.Query().Get("query")

	collection := database.GetCollection("vendors")
	findOptions := options.Find().SetLimit(8)

	filter := bson.M{}
	if query != "" {
		filter = bson.M{
			"$or": []bson.M{
				{"name": bson.M{"$regex": query, "$options": "i"}},
				{"seName": bson.M{"$regex": query, "$options": "i"}},
			},
		}
	}

	cursor, err := collection.Find(ctx, filter, findOptions)
	if err != nil {
		http.Error(w, "Failed to fetch products", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	var vendors []models.ProductVendor
	if err = cursor.All(ctx, &vendors); err != nil {
		http.Error(w, "Failed to decode products", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(vendors)
}
