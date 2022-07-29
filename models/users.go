package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID             primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name           string             `json:"name" bson:"name"`
	Username       string             `json:"username" bson:"username" validate:"required"`
	Password       string             `json:"password,omitempty" bson:"password" validate:"required"`
	DailyTaskLimit int                `json:"daily_task_limit,omitempty" bson:"daily_task_limit"`
	CreatedAt      int64              `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt      int64              `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}
