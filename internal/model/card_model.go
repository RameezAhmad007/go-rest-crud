package model

import (
	"github.com/RameezAhmad007/go-rest-crud/internal/api"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Card struct {
	Id     primitive.ObjectID `json:"id" bson:"_id"`
	Name   string             `json:"name" bson:"name"`
	Number int                `json:"number" bson:"number"`
}

// ToAPICard converts model.Card to api.Card
func ToAPICard(mCard Card) api.Card {
	idStr := mCard.Id.Hex()
	return api.Card{
		Id:     &idStr,
		Name:   mCard.Name,
		Number: mCard.Number,
	}
}

// FromAPICard converts api.Card to model.Card
func FromAPICard(aCard api.Card) (Card, error) {
	var id primitive.ObjectID
	if aCard.Id != nil {
		oid, err := primitive.ObjectIDFromHex(*aCard.Id)
		if err != nil {
			return Card{}, err
		}
		id = oid
	}
	return Card{
		Id:     id,
		Name:   aCard.Name,
		Number: aCard.Number,
	}, nil
}
