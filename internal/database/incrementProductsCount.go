package database

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func IncrementVendorProductsCount(ctx context.Context, vendorId primitive.ObjectID) error {
	vendorsCollection := VendorsCollection()
	_, err := vendorsCollection.UpdateOne(ctx, bson.M{"_id": vendorId}, bson.M{"$inc": bson.M{"productCount": 1}})
	if err != nil {
		return err
	}

	return nil
}

func IncrementCategoryProductsCount(ctx context.Context, categoryId primitive.ObjectID) error {
	categoriesCollection := CategoriesCollection()
	_, err := categoriesCollection.UpdateOne(ctx, bson.M{"_id": categoryId}, bson.M{"$inc": bson.M{"productCount": 1}})
	if err != nil {
		return err
	}

	return nil
}
