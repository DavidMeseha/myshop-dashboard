package database

import (
	"context"
	"shop-dashboard/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type OriginalProduct struct {
	Category    primitive.ObjectID `bson:"category"`
	ProductTags []string           `bson:"productTags"`
}

func CreateProduct(ctx context.Context, product models.Product) (*mongo.InsertOneResult, error) {
	productsCollection := ProductsCollection()
	result, err := productsCollection.InsertOne(ctx, product)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func HardDeleteProduct(ctx context.Context, productId primitive.ObjectID) {
	productsCollection := ProductsCollection()
	productsCollection.DeleteOne(ctx, bson.M{"_id": productId})
}

func ChangeProductDeleteState(ctx context.Context, productId string, delete bool) (*mongo.UpdateResult, error) {
	productID, err := primitive.ObjectIDFromHex(productId)
	if err != nil {
		return nil, err
	}

	collection := ProductsCollection()
	filter := bson.M{"_id": productID}
	update := bson.M{"$set": bson.M{"deleted": delete}}

	res, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func GetProduct(ctx context.Context, productID string, vendorID primitive.ObjectID) (bson.M, error) {
	ID, err := primitive.ObjectIDFromHex(productID)
	if err != nil {
		return nil, err
	}

	pipeline := bson.A{
		bson.M{"$match": bson.M{"_id": ID, "vendor": vendorID}},
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

	collection := ProductsCollection()
	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var products []bson.M
	if err := cursor.All(ctx, &products); err != nil || len(products) == 0 {
		return nil, err
	}

	return products[0], nil
}

func UpdateProduct(ctx context.Context, productID string, vendorID primitive.ObjectID, updateFields bson.M) (*mongo.UpdateResult, error) {
	ID, err := primitive.ObjectIDFromHex(productID)
	if err != nil {
		return nil, err
	}

	collection := ProductsCollection()
	res, err := collection.UpdateOne(
		ctx,
		bson.M{"_id": ID, "vendor": vendorID},
		bson.M{"$set": updateFields},
	)
	if err != nil {
		return nil, err
	}

	return res, nil
}
