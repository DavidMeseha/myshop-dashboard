package handlers

import (
	"context"
	"encoding/json"
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
	category := r.URL.Query().Get("category")
	vendorID := vendor.ID

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	var categoryID primitive.ObjectID
	var err error
	if category != "" {
		categoryID, err = primitive.ObjectIDFromHex(category)
		if err != nil {
			http.Error(w, "Invalid category ID", http.StatusBadRequest)
			return
		}
	} else {
		categoryID = primitive.NilObjectID
	}

	products, totalPages, totalCount, err := database.FilterProducts(ctx, query, limit, page, vendorID, categoryID)
	if err != nil {
		http.Error(w, "Could not fetch products", http.StatusInternalServerError)
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
