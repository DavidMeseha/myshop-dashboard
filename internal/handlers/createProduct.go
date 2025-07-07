package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"shop-dashboard/internal/database"
	"shop-dashboard/internal/middleware"
	"shop-dashboard/internal/models"
	"shop-dashboard/internal/utils"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CreateProductRequest struct {
	Pictures    []string                  `json:"images"`
	Attributes  []models.ProductAttribute `json:"attributes"`
	Category    primitive.ObjectID        `json:"category"`
	Description string                    `json:"fullDescription"`
	Price       models.ProductPrice       `json:"price"`
	Name        string                    `json:"name"`
	SeName      string                    `json:"seName"`
	SKU         string                    `json:"sku"`
	Stock       int64                     `json:"stock"`
	Tags        []string                  `json:"tags" `
	Gender      string                    `json:"gender"`
}

func ProductDto(ctx context.Context, body CreateProductRequest) (models.Product, error) {
	inStock := body.Stock > 0
	hasAttributes := len(body.Attributes) > 0
	err := utils.ProcessTags(ctx, body.Tags)
	if err != nil {
		return models.Product{}, err
	}

	product := models.Product{
		Name:           body.Name,
		SeName:         body.SeName,
		SKU:            body.SKU,
		Price:          body.Price,
		Description:    body.Description,
		HasAttributes:  hasAttributes,
		Attributes:     utils.ProcessAttributes(body.Attributes),
		InStock:        inStock,
		Stock:          body.Stock,
		Likes:          0,
		Carts:          0,
		Saves:          0,
		UsersLiked:     []string{},
		UsersSaved:     []string{},
		UsersReviewed:  []string{},
		UsersCarted:    []string{},
		ProductReviews: []primitive.ObjectID{},
		ReviewOverview: models.ProductReviewOverview{
			RatingSum:    0,
			TotalReviews: 0,
		},
		Pictures:   utils.ProcessPictures(body.Pictures, body.Name),
		Tags:       body.Tags,
		Gender:     body.Gender,
		CategoryID: body.Category,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	return product, nil
}

func (h *Handler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	vendor, ok := middleware.GetUserInfo(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var body CreateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		log.Print(err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	product, err := ProductDto(ctx, body)
	if err != nil {
		http.Error(w, "Could not process product data", http.StatusInternalServerError)
		return
	}
	product.VendorID = vendor.ID

	if product.Name == "" || product.SeName == "" || product.SKU == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	productsCollection := database.GetCollection("products")
	result, err := productsCollection.InsertOne(ctx, product)
	if err != nil {
		log.Print(err)
		http.Error(w, "Failed to create product", http.StatusInternalServerError)
		return
	}

	vendorsCollection := database.GetCollection("vendors")
	_, err = vendorsCollection.UpdateOne(ctx, bson.M{"_id": vendor.ID}, bson.M{"$inc": bson.M{"productCount": 1}})
	if err != nil {
		log.Print(err)
		productsCollection.DeleteOne(ctx, bson.M{"_id": result.InsertedID})
		http.Error(w, "Failed to update vendor product count", http.StatusInternalServerError)
		return
	}

	categoriesCollection := database.GetCollection("categories")
	_, err = categoriesCollection.UpdateOne(ctx, bson.M{"_id": product.CategoryID}, bson.M{"$inc": bson.M{"productCount": 1}})
	if err != nil {
		log.Print(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]any{
		"id":      result.InsertedID,
		"message": "Product created successfully",
	})
}
