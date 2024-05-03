package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Book struct {
	Id        primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	Title     string             `bson:"title" json:"title"`
	Author    string             `bson:"author" json:"author"`
	Year      int                `bson:"year" json:"year"`
	CreatedAt primitive.DateTime `bson:"created_at" json:"created_at"`
}
