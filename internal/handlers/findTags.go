package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"shop-dashboard/internal/database"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Tag struct {
	ID           string `json:"_id" bson:"_id"`
	Name         string `json:"name" bson:"name"`
	SeName       string `json:"seName" bson:"seName"`
	ProductCount int    `json:"productCount" bson:"productCount"`
}

func (h *Handler) FindTags(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	query := r.URL.Query().Get("query")
	collection := database.GetCollection("tags")
	findOptions := options.Find().SetLimit(10)

	filter := bson.M{}
	if query != "" {
		filter = bson.M{
			"name": bson.M{
				"$regex":   query,
				"$options": "i", // case-insensitive
			},
		}
	}

	cursor, err := collection.Find(ctx, filter, findOptions)
	if err != nil {
		http.Error(w, "Failed to fetch tags", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	var tags []Tag
	if err := cursor.All(ctx, &tags); err != nil {
		http.Error(w, "Failed to decode tags", http.StatusInternalServerError)
		return
	}

	log.Print(tags)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tags)
}
