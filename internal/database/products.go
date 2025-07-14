package database

import (
	"context"
	"math"
	"regexp"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetProductsBySKU(ctx context.Context, baseSKU string) (*mongo.Cursor, error) {
	collection := ProductsCollection()

	filter := bson.M{"sku": bson.M{"$regex": "^" + regexp.QuoteMeta(baseSKU)}}
	opts := options.Find().SetProjection(bson.M{"seName": 1, "sku": 1})
	cursor, err := collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	return cursor, nil
}

func FilterProducts(ctx context.Context, query string, limit int64, page int64, vendorID primitive.ObjectID, categoryID primitive.ObjectID) ([]bson.M, int64, int64, error) {
	filter := bson.M{}
	filter["vendor"] = vendorID
	if categoryID != primitive.NilObjectID {
		filter["category"] = categoryID
	}

	if query != "" {
		filter["$or"] = []bson.M{
			{"name": bson.M{"$regex": query, "$options": "i"}},
			{"description": bson.M{"$regex": query, "$options": "i"}},
		}
	}

	collection := ProductsCollection()
	totalCount, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, 0, err
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
		return nil, 0, 0, err
	}
	defer cursor.Close(ctx)

	var products []bson.M
	if err = cursor.All(ctx, &products); err != nil {
		return nil, 0, 0, err
	}

	return products, totalPages, totalCount, nil
}
