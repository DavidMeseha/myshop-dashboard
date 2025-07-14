package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AttributeControlType string

const (
	DropdownList AttributeControlType = "DropdownList"
	RadioList    AttributeControlType = "RadioList"
	Checkboxes   AttributeControlType = "Checkboxes"
	TextBox      AttributeControlType = "TextBox"
	ColorSquares AttributeControlType = "Color"
)

type ProductPicture struct {
	ImageUrl      string `bson:"imageUrl" json:"imageUrl"`
	Title         string `bson:"title" json:"title"`
	AlternateText string `bson:"alternateText" json:"alternateText"`
}

type Tag struct {
	ID           string `json:"_id" bson:"_id"`
	Name         string `json:"name" bson:"name"`
	SeName       string `json:"seName" bson:"seName"`
	ProductCount int    `json:"productCount" bson:"productCount"`
}

type ProductPrice struct {
	OldPrice float64 `bson:"oldPrice" json:"oldPrice"`
	Price    float64 `bson:"price" json:"price"`
}

type ProductAttributeValue struct {
	ID                   primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Name                 string             `bson:"name" json:"name"`
	ColorSquaresRgb      string             `bson:"colorRgb,omitempty" json:"ColorRgb,omitempty"`
	PriceAdjustmentValue float64            `bson:"priceAdjustmentValue" json:"priceAdjustmentValue"`
}

type ProductAttribute struct {
	ID                   primitive.ObjectID      `bson:"_id,omitempty" json:"_id"`
	Name                 string                  `bson:"name" json:"name"`
	AttributeControlType AttributeControlType    `bson:"attributeControlType" json:"attributeControlType"`
	Values               []ProductAttributeValue `bson:"values" json:"values"`
}

type ProductReviewOverview struct {
	RatingSum    int `bson:"ratingSum" json:"ratingSum"`
	TotalReviews int `bson:"totalReviews" json:"totalReviews"`
}

type ProductCategory struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Name          string             `bson:"name" json:"name"`
	SeName        string             `bson:"seName" json:"seName"`
	ProductsCount int                `bson:"productsCount" json:"productsCount"`
}

type ProductVendor struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Name           string             `bson:"name" json:"name"`
	SeName         string             `bson:"seName" json:"seName"`
	ImageUrl       string             `bson:"imageUrl" json:"imageUrl"`
	ProductCount   int                `bson:"productCount" json:"ProductCount"`
	FollowersCount int                `bson:"followersCount" json:"followersCount"`
	User           primitive.ObjectID `bson:"user" json:"user"`
}

type ProductReview struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	Product    Product            `bson:"product" json:"product"`
	Customer   User               `bson:"customer" json:"customer"`
	ReviewText string             `bson:"reviewText" json:"reviewText"`
	Rating     int                `bson:"rating" json:"rating"`
}

type Product struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`

	// Basic Fields
	Name   string `json:"name" bson:"name"`
	SeName string `json:"se_name" bson:"seName"`
	SKU    string `json:"sku" bson:"sku"`

	// Media Fields
	Pictures []ProductPicture `json:"pictures" bson:"pictures"`

	// Pricing Fields
	Price ProductPrice `json:"price" bson:"price"`

	// Meta Fields
	Description     string `json:"fullDescription" bson:"fullDescription"`
	MetaDescription string `json:"metaDescription" bson:"metaDescription"`
	MetaKeywords    string `json:"metaKeywords" bson:"metaKeywords"`
	MetaTitle       string `json:"metaTitle" bson:"metaTitle"`

	// Attribute Fields
	HasAttributes bool               `json:"hasAttributes" bson:"hasAttributes"`
	Attributes    []ProductAttribute `json:"productAttributes" bson:"productAttributes"`

	// Status Fields
	InStock bool `json:"in_stock" bson:"inStock"`

	// Stats Fields
	Stock          int64                 `json:"stock" bson:"stock"`
	Likes          int64                 `json:"likes" bson:"likes"`
	Carts          int64                 `json:"carts" bson:"carts"`
	Saves          int64                 `json:"saves" bson:"saves"`
	UsersLiked     []string              `json:"usersLiked" bson:"usersLiked"`
	UsersSaved     []string              `json:"usersSaved" bson:"usersSaved"`
	UsersCarted    []string              `json:"usersCarted" bson:"usersCarted"`
	UsersReviewed  []string              `json:"usersReviewed" bson:"usersReviewed"`
	ReviewOverview ProductReviewOverview `json:"productReviewOverview" bson:"productReviewOverview"`

	// Relation Fields
	Gender         string               `json:"gender" bson:"gender"`
	CategoryID     primitive.ObjectID   `json:"category" bson:"category"`
	VendorID       primitive.ObjectID   `json:"vendor" bson:"vendor"`
	ProductReviews []primitive.ObjectID `json:"productReviews" bson:"productReviews"`
	Tags           []string             `json:"productTags" bson:"productTags"`
	CreatedAt      time.Time            `json:"createdAt" bson:"createdAt"`
	UpdatedAt      time.Time            `json:"updatedAt" bson:"updatedAt"`
}
