package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"shop-dashboard/internal/database"
	"shop-dashboard/internal/middleware"
	"time"

	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (h *Handler) GetProduct(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "Missing product id", http.StatusBadRequest)
		return
	}

	vendor, ok := middleware.GetUserInfo(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusNotFound)
		return
	}

	productID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		http.Error(w, "Invalid product id", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	collection := database.GetCollection("products")

	pipeline := bson.A{
		bson.M{"$match": bson.M{"_id": productID, "vendor": vendor.ID}},
		bson.M{"$lookup": bson.M{
			"from":         "categories",
			"localField":   "category",
			"foreignField": "_id",
			"as":           "category",
		}},
		bson.M{"$unwind": bson.M{
			"path":                       "$category",
			"preserveNullAndEmptyArrays": true,
		}},
	}

	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}
	defer cursor.Close(ctx)

	var products []bson.M
	if err := cursor.All(ctx, &products); err != nil || len(products) == 0 {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products[0])
}
