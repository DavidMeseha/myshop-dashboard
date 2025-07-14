package utils

import (
	"context"
	"shop-dashboard/internal/database"

	"go.mongodb.org/mongo-driver/bson"
)

func ProcessUpdatedTags(ctx context.Context, tags []string, productTags []string) error {
	existingTags, newTags, err := getTags(ctx, tags)
	if err != nil {
		return err
	}

	toDecriment := NonIntersecting(productTags, existingTags)
	if len(toDecriment) > 0 {
		database.DecrimentTags(ctx, toDecriment)
	}
	toIncriment := NonIntersecting(existingTags, productTags)
	if len(toIncriment) > 0 {
		database.IncrimentTags(ctx, toIncriment)
	}

	database.InsetNewTags(ctx, newTags)

	return nil
}

func ProcessTags(ctx context.Context, tags []string) error {
	existingTags, newTags, err := getTags(ctx, tags)
	if err != nil {
		return err
	}

	// Bulk increment productCount for all existing tags
	if len(existingTags) > 0 {
		database.IncrimentTags(ctx, existingTags)
	}

	// Insert new tags and collect their IDs
	if len(newTags) > 0 {
		err := database.InsetNewTags(ctx, newTags)
		if err != nil {
			return err
		}
	}

	return nil
}

func getTags(ctx context.Context, tags []string) ([]string, []string, error) {
	tagCollection := database.GetCollection("tags")

	cursor, err := tagCollection.Find(ctx, bson.M{"seName": bson.M{"$in": tags}})
	if err != nil {
		return nil, nil, err
	}
	defer cursor.Close(ctx)

	existingTags := []string{}
	for cursor.Next(ctx) {
		var tagDoc struct {
			SeName string `bson:"seName"`
		}
		if err := cursor.Decode(&tagDoc); err == nil {
			existingTags = append(existingTags, tagDoc.SeName)
		}
	}
	newTags := NonIntersecting(tags, existingTags)

	return existingTags, newTags, nil
}

func NonIntersecting(a, b []string) []string {
	m := make(map[string]struct{})
	for _, v := range b {
		m[v] = struct{}{}
	}
	var result []string
	for _, v := range a {
		if _, found := m[v]; !found {
			result = append(result, v)
		}
	}
	return result
}
