package database

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func UpdateUserIsVendorState(ctx context.Context, userID primitive.ObjectID) error {
	usersCollection := UsersCollection()
	_, err := usersCollection.UpdateOne(
		ctx,
		bson.M{"_id": userID},
		bson.M{"$set": bson.M{"isVendor": true}},
	)

	return err
}
