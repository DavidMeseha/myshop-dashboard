package handlers

import (
	"context"
	"encoding/json"
	"math"
	"net/http"
	"shop-dashboard/internal/database"
	"shop-dashboard/internal/middleware"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (h *Handler) GetProducts(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	vendor, ok := middleware.GetUserInfo(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	page, _ := strconv.ParseInt(r.URL.Query().Get("page"), 10, 64)
	limit, _ := strconv.ParseInt(r.URL.Query().Get("limit"), 10, 64)
	query := r.URL.Query().Get("query")
	vendorID := vendor.ID
	categoryID := r.URL.Query().Get("category")

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	filter := bson.M{}
	if query != "" {
		filter["$or"] = []bson.M{
			{"name": bson.M{"$regex": query, "$options": "i"}},
			{"description": bson.M{"$regex": query, "$options": "i"}},
		}
	}
	filter["vendor"] = vendorID

	if vendorID == primitive.NilObjectID {
		http.Error(w, "Failed to fetch products", http.StatusUnauthorized)
		return
	}

	if categoryID != "" {
		objectID, err := primitive.ObjectIDFromHex(categoryID)
		if err != nil {
			http.Error(w, "Invalid category ID", http.StatusBadRequest)
			return
		}
		filter["category"] = objectID
	}

	collection := database.GetCollection("products")

	// Count total documents for pagination
	totalCount, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		http.Error(w, "Failed to fetch counts", http.StatusInternalServerError)
		return
	}
	totalPages := int64(math.Ceil(float64(totalCount) / float64(limit)))
	skip := (page - 1) * limit

	pipeline := bson.A{
		bson.M{"$match": filter},
		bson.M{"$sort": bson.M{"_id": -1}},
		bson.M{"$skip": skip},
		bson.M{"$limit": limit},
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
		http.Error(w, "Failed to fetch products", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	var products []bson.M
	if err = cursor.All(ctx, &products); err != nil {
		http.Error(w, "Failed to decode products", http.StatusInternalServerError)
		return
	}

	response := struct {
		Data         []bson.M `json:"data"`
		CurrentPage  int64    `json:"currentPage"`
		NextPage     int64    `json:"nextPage"`
		PreviousPage int64    `json:"previousPage"`
		Limit        int64    `json:"limit"`
		TotalPages   int64    `json:"totalPages"`
		TotalCount   int64    `json:"totalCount"`
		HasNext      bool     `json:"hasNext"`
		HasPrevious  bool     `json:"hasPrevious"`
	}{
		Data:         products,
		CurrentPage:  page,
		NextPage:     page + 1,
		PreviousPage: page - 1,
		Limit:        limit,
		TotalPages:   totalPages,
		TotalCount:   totalCount,
		HasNext:      page < totalPages,
		HasPrevious:  page > 1,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
