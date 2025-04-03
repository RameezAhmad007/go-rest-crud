package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Card struct {
	ID     primitive.ObjectID `json:"id" bson:"_id"`
	Name   string             `json:"name" bson:"name"`
	Number int                `json:"number" bson:"number"`
}
