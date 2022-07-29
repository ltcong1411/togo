package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Task struct {
	ID        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	UserID    string             `json:"user_id" bson:"user_id" validate:"required"`
	Content   string             `json:"content" bson:"content" validate:"required"`
	Completed bool               `json:"completed" bson:"completed"`
	CreatedAt int64              `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt int64              `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}
