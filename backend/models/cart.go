package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CartItem 表示购物车中的单个商品
type CartItem struct {
	ProductID   primitive.ObjectID `bson:"product_id" json:"product_id"`
	ProductType string             `bson:"product_type" json:"product_type"` // "mouse", "keyboard", 等
	Name        string             `bson:"name" json:"name"`
	Price       float64            `bson:"price" json:"price"`
	Quantity    int                `bson:"quantity" json:"quantity"`
	ImageURL    string             `bson:"image_url" json:"image_url,omitempty"`
}

// Cart 表示用户的购物车
type Cart struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	UserID    primitive.ObjectID `bson:"user_id" json:"user_id"`
	Items     []CartItem         `bson:"items" json:"items"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
}