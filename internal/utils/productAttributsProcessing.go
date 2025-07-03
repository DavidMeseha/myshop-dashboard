package utils

import (
	"shop-dashboard/internal/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ProcessAttributes(attris []models.ProductAttribute) []models.ProductAttribute {
	for i, attr := range attris {
		if attris[i].ID == primitive.NilObjectID {
			attris[i].ID = primitive.NewObjectID()
		}
		for j := range attr.Values {
			if attris[i].Values[j].ID == primitive.NilObjectID {
				attris[i].Values[j].ID = primitive.NewObjectID()
			}
		}
	}

	return attris
}
