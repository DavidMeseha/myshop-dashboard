package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"shop-dashboard/internal/database"
	"shop-dashboard/internal/utils"
	"strings"
	"sync"
	"time"
)

type UniqueProductRequest struct {
	Name string `json:"name"`
}

type UniqueProductResponse struct {
	SeName string `json:"seName"`
	SKU    string `json:"sku"`
}

func (h *Handler) GenerateProductUniques(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	var req UniqueProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || strings.TrimSpace(req.Name) == "" {
		http.Error(w, "Invalid or missing product name", http.StatusBadRequest)
		return
	}

	baseSeName := utils.GenerateSeName(req.Name)
	baseSKU := utils.GenerateSKU(req.Name)

	cursor, err := database.GetProductsBySKU(ctx, baseSKU)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	existingSeNames := make(map[string]struct{})
	existingSKUs := make(map[string]struct{})
	for cursor.Next(ctx) {
		var doc struct {
			SeName string `bson:"seName"`
			SKU    string `bson:"sku"`
		}
		if err := cursor.Decode(&doc); err == nil {
			existingSeNames[doc.SeName] = struct{}{}
			existingSKUs[doc.SKU] = struct{}{}
		}
	}

	if len(existingSKUs) == 0 {
		res := UniqueProductResponse{
			SeName: baseSeName,
			SKU:    baseSKU,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(res)
		return
	}

	var wg sync.WaitGroup
	seNameCh := make(chan string, 1)
	skuCh := make(chan string, 1)

	wg.Add(1)
	go func() {
		defer wg.Done()
		seName := utils.GenerateUniqueSeName(baseSeName, existingSeNames)
		seNameCh <- seName
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		sku := utils.GenerateUniqueSku(baseSKU, existingSKUs)
		skuCh <- sku
	}()

	wg.Wait()
	close(seNameCh)
	close(skuCh)

	resp := UniqueProductResponse{
		SeName: <-seNameCh,
		SKU:    <-skuCh,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
