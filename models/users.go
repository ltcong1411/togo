package models

type User struct {
	ID        string `json:"id" bson:"_id,omitempty"`
	Name      string `json:"name" bson:"name"`
	Username  string `json:"username" bson:"username" validate:"required"`
	Password  string `json:"password,omitempty" bson:"password" validate:"required"`
	Token     string `json:"token,omitempty" bson:"-"`
	CreatedAt int64  `json:"created_at,omitempty" bson:"created_at,omitempty"`
}
