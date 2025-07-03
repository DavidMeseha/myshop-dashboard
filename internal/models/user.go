package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type UserDateOfBirth struct {
	Day   int `bson:"day" json:"day"`
	Month int `bson:"month" json:"month"`
	Year  int `bson:"year" json:"year"`
}

type City struct {
	ID   primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	Name string             `bson:"name" json:"name"`
	Code string             `bson:"code" json:"code"`
}

type Country struct {
	ID     primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	Name   string             `bson:"name" json:"name"`
	Code   string             `bson:"code" json:"code"`
	Cities []City             `bson:"cities" json:"cities"`
}

type UserAddress struct {
	ID      primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	Address string             `bson:"address" json:"address"`
	City    City               `bson:"city" json:"name"`
	Country Country            `bson:"country" json:"country"`
}

type UserCart struct {
	Product    Product            `bson:"product" json:"product"`
	Quantity   int                `bson:"quantity" json:"quantity"`
	Attributes []ProductAttribute `bson:"attributes" json:"attributes"`
}

type User struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"_id"`

	// Profile information
	ImageUrl    string          `bson:"imageUrl" json:"image_url"`
	FirstName   string          `bson:"firstName" json:"first_name"`
	LastName    string          `bson:"lastName" json:"last_name"`
	Email       string          `bson:"email" json:"email"`
	Password    string          `bson:"password" json:"password"`
	Phone       string          `bson:"phone" json:"phone"`
	Gender      string          `bson:"gender" json:"gender"`
	DateOfBirth UserDateOfBirth `bson:"userDateOfBirth" json:"user_date_of_birth"`

	// User preferences
	Language string `bson:"language" json:"language"`

	// Status flags
	IsLogin      bool `bson:"isLogin" json:"is_login"`
	IsRegistered bool `bson:"isRegistered" json:"is_registered"`
	IsVendor     bool `bson:"isVendor" json:"is_vendor"`

	// Relations and collections
	Saves          []Product       `bson:"saves" json:"saves"`
	Cart           []UserCart      `bson:"cart" json:"cart"`
	Following      []ProductVendor `bson:"following" json:"following"`
	RecentProducts []Product       `bson:"recentProducts" json:"recent_products"`
	Addresses      []UserAddress   `bson:"addesses" json:"addesses"`
	Orders         []string        `bson:"orders" json:"orders"`
}
