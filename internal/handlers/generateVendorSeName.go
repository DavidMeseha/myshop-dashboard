package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"shop-dashboard/internal/database"
	"shop-dashboard/internal/utils"
	"strings"
	"time"
)

type VendorSeNameRequest struct {
	Name string `json:"name"`
}

type VendorSeNameResponse struct {
	SeName string `json:"seName"`
}

func (h *Handler) GenerateVendorSeName(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	var req VendorSeNameRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || strings.TrimSpace(req.Name) == "" {
		http.Error(w, "Invalid or missing product name", http.StatusBadRequest)
		return
	}

	baseSeName := utils.GenerateSeName(req.Name)

	cursor, err := database.FindVendorsBySeName(ctx, baseSeName)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	existingSeNames := make(map[string]struct{})
	for cursor.Next(ctx) {
		var doc struct {
			SeName string `bson:"seName"`
		}
		if err := cursor.Decode(&doc); err == nil {
			existingSeNames[doc.SeName] = struct{}{}
		}
	}

	if len(existingSeNames) == 0 {
		res := VendorSeNameResponse{
			SeName: baseSeName,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(res)
		return
	}

	seName := utils.GenerateUniqueSeName(baseSeName, existingSeNames)

	resp := UniqueProductResponse{
		SeName: seName,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
