package database

import (
	"context"
	"shop-dashboard/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetFilteredCategories(ctx context.Context, query string, limit int64) ([]models.ProductCategory, error) {
	collection := CategoriesCollection()

	findOptions := options.Find().SetLimit(limit)
	filter := bson.M{}
	if query != "" {
		filter = bson.M{
			"$or": []bson.M{
				{"name": bson.M{"$regex": query, "$options": "i"}},
				{"seName": bson.M{"$regex": query, "$options": "i"}},
			},
		}
	}

	cursor, err := collection.Find(ctx, filter, findOptions)
	if err != nil {
		return []models.ProductCategory{}, err
	}
	defer cursor.Close(ctx)

	var categories []models.ProductCategory
	if err = cursor.All(ctx, &categories); err != nil {
		return nil, err
	}

	return categories, nil
}
