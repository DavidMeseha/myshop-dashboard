package database

import (
	"context"
	"shop-dashboard/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func InsetNewTags(ctx context.Context, newTags []string) error {
	tagCollection := TagsCollection()
	tagsToAdd := bson.A{}
	for _, tagName := range newTags {
		tagsToAdd = append(tagsToAdd, bson.M{"name": tagName, "seName": tagName, "productCount": 1})
	}

	_, err := tagCollection.InsertMany(ctx, tagsToAdd)
	if err != nil {
		return err
	}

	return nil
}

func DecrimentTags(ctx context.Context, tags []string) {
	tagCollection := TagsCollection()
	_, _ = tagCollection.UpdateMany(
		ctx,
		bson.M{"seName": bson.M{"$in": tags}},
		bson.M{"$inc": bson.M{"productCount": -1}},
	)
}

func IncrimentTags(ctx context.Context, existingTags []string) {
	tagCollection := TagsCollection()
	_, _ = tagCollection.UpdateMany(
		ctx,
		bson.M{"seName": bson.M{"$in": existingTags}},
		bson.M{"$inc": bson.M{"productCount": 1}},
	)
}

func GetFilteredTags(ctx context.Context, query string, limit int64) ([]models.Tag, error) {
	collection := TagsCollection()

	findOptions := options.Find().SetLimit(limit)
	filter := bson.M{}
	if query != "" {
		filter = bson.M{
			"name": bson.M{
				"$regex":   query,
				"$options": "i",
			},
		}
	}

	cursor, err := collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var tags []models.Tag
	if err := cursor.All(ctx, &tags); err != nil {
		return nil, err
	}

	return tags, nil
}
