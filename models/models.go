package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

//Create Struct
type User struct {
	ID       primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name     string             `json:"name,omitempty" bson:"name,omitempty"`
	Email    string             `json:"email,omitempty" bson:"email,omitempty"`
	Username string             `json:"username,omitempty" bson:"username,omitempty"`
	Password string             `json:"password,omitempty" bson:"password,omitempty"`
	Dob      string             `json:"dob,omitempty" bson:"dob,omitempty"`
	Phone    string             `json:"phone,omitempty" bson:"phone,omitempty"`
}

type Post struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Title     string             `json:"title" bson:"title,omitempty"`
	Body      string             `json:"body" bson:"body,omitempty"`
	PostedOn  time.Time          `json:"postedOn" bson:"postedOn,omitempty"`
	Author    primitive.ObjectID `json:"author" bson:"author,omitempty"`
	Thumbnail string             `json:"thumbnail" bson:"thumbnail,omitempty"`
}
