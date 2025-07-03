package handlers

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"shop-dashboard/internal/database"
	"shop-dashboard/internal/middleware"
	"shop-dashboard/internal/models"
	"shop-dashboard/internal/utils"
	"time"

	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type EditProductRequest struct {
	Pictures    []string                  `json:"images,omitempty" bson:"pictures,omitempty"`
	Attributes  []models.ProductAttribute `json:"attributes,omitempty" bson:"productAttributes,omitempty"`
	Category    primitive.ObjectID        `json:"category,omitempty" bson:"category,omitempty"`
	Description string                    `json:"fullDescription,omitempty" bson:"fullDescription,omitempty"`
	Price       models.ProductPrice       `json:"price,omitempty" bson:"price,omitempty"`
	Name        string                    `json:"name,omitempty" bson:"name,omitempty"`
	SeName      string                    `json:"seName,omitempty" bson:"seName,omitempty"`
	SKU         string                    `json:"sku,omitempty" bson:"sku,omitempty"`
	Stock       int64                     `json:"stock,omitempty" bson:"stock,omitempty"`
	Tags        []string                  `json:"tags,omitempty" bson:"productTags,omitempty"`
	Gender      string                    `json:"gender,omitempty" bson:"gender,omitempty"`
}

func (h *Handler) EditProductData(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	vendor, ok := middleware.GetUserInfo(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

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
	productsCollection := database.GetCollection("products")

	var originalProduct struct {
		Category    primitive.ObjectID `bson:"category"`
		ProductTags []string           `bson:"productTags"`
	}
	err = productsCollection.FindOne(ctx, bson.M{"_id": productID, "vendor": vendor.ID}).Decode(&originalProduct)
	if err != nil {
		http.Error(w, "Product not found or not owned by vendor", http.StatusNotFound)
		return
	}

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	var fields EditProductRequest
	if err := json.Unmarshal(bodyBytes, &fields); err != nil {
		http.Error(w, "Invalid request body", http.StatusNoContent)
		return
	}

	var updateFields bson.M
	if err := json.Unmarshal(bodyBytes, &updateFields); err != nil {
		http.Error(w, "Invalid body", http.StatusBadRequest)
		return
	}

	if len(updateFields) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	if fields.Tags != nil {
		updateFields["productTags"] = fields.Tags
		delete(updateFields, "tags")
		err = utils.ProcessUpdatedTags(ctx, fields.Tags, originalProduct.ProductTags)
		if err != nil {
			http.Error(w, "Failed to Process product tags", http.StatusInternalServerError)
			return
		}
	}

	if fields.Attributes != nil {
		delete(updateFields, "attributes")
		updateFields["productAttributes"] = utils.ProcessAttributes(fields.Attributes)
	}

	if fields.Pictures != nil {
		delete(updateFields, "pictures")
		updateFields["pictures"] = utils.ProcessPictures(fields.Pictures, fields.Name)
	}

	updateFields["updatedAt"] = time.Now()

	res, err := productsCollection.UpdateOne(
		ctx,
		bson.M{"_id": productID, "vendor": vendor.ID},
		bson.M{"$set": updateFields},
	)
	if err != nil {
		http.Error(w, "Failed to update product", http.StatusInternalServerError)
		return
	}
	if res.MatchedCount == 0 {
		http.Error(w, "Product not found or not owned by vendor", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"_id":     id,
		"message": "Product updated successfully",
	})
}
