package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Blog struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserEmail   string             `bson:"user_email" json:"userEmail"`
	Title       string             `bson:"title" json:"title" binding:"required"`
	Content string             `bson:"description" json:"description"`
	Category      string             `bson:"status" json:"status" binding:"required,oneof=pending done"`
	Tags     []string              `bson:"tags" json:"tags"`
	CreatedAt   int64              `bson:"created_at" json:"createdAt"`
	UpdatedAt   int64              `bson:"updated_at" json:"updatedAt"`
}

