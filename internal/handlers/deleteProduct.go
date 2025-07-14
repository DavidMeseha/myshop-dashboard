package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"shop-dashboard/internal/database"
	"time"

	"github.com/go-chi/chi/v5"
)

func (h *Handler) SoftDeleteProduct(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "Missing product id", http.StatusBadRequest)
		return
	}

	res, err := database.ChangeProductDeleteState(ctx, id, true)
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

	res, err := database.ChangeProductDeleteState(ctx, id, false)
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
