package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"shop-dashboard/internal/database"
	"shop-dashboard/internal/middleware"
	"time"

	"github.com/go-chi/chi/v5"
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

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	product, err := database.GetProduct(ctx, id, vendor.ID)
	if err != nil {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(product)
}
