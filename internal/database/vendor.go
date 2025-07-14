package database

import (
	"context"
	"regexp"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func FindVendorsBySeName(ctx context.Context, baseSeName string) (*mongo.Cursor, error) {
	collection := VendorsCollection()

	filter := bson.M{"seName": bson.M{"$regex": "^" + regexp.QuoteMeta(baseSeName)}}
	opts := options.Find().SetProjection(bson.M{"seName": 1})
	cursor, err := collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	return cursor, nil
}

func CheckVendorSeName(ctx context.Context, seName string) (int64, error) {
	vendorCollection := VendorsCollection()

	count, err := vendorCollection.CountDocuments(ctx, bson.M{"seName": seName})
	if err != nil {
		return 0, err
	}
	return count, nil
}

func CreateVendor(ctx context.Context, data bson.M) (*mongo.InsertOneResult, error) {
	vendorCollection := VendorsCollection()
	res, err := vendorCollection.InsertOne(ctx, data)
	if err != nil {
		return nil, err
	}

	return res, nil
}
