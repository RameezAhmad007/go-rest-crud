package repository

import (
	"context"
	"time"

	"github.com/RameezAhmad007/go-rest-crud/internal/api"
	"github.com/RameezAhmad007/go-rest-crud/internal/config"
	"github.com/RameezAhmad007/go-rest-crud/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var cardsCollection = config.GetCollection("cardscollection")

func CreateCard(ctx context.Context, card api.Card) (api.Card, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	mCard, err := model.FromAPICard(card)
	if err != nil {
		return api.Card{}, err
	}

	var existingCard model.Card
	err = cardsCollection.FindOne(ctx, bson.M{"name": mCard.Name}).Decode(&existingCard)
	if err == nil {
		return api.Card{}, mongo.ErrClientDisconnected // Simulate conflict
	} else if err != mongo.ErrNoDocuments {
		return api.Card{}, err
	}

	mCard.Id = primitive.NewObjectID()
	_, err = cardsCollection.InsertOne(ctx, mCard)
	if err != nil {
		return api.Card{}, err
	}
	return model.ToAPICard(mCard), nil
}

func GetCard(ctx context.Context, id string) (api.Card, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return api.Card{}, err
	}

	var mCard model.Card
	err = cardsCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&mCard)
	if err != nil {
		return api.Card{}, err
	}
	return model.ToAPICard(mCard), nil
}

func GetAllCards(ctx context.Context) ([]api.Card, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var mCards []model.Card
	cursor, err := cardsCollection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	err = cursor.All(ctx, &mCards)
	if err != nil {
		return nil, err
	}

	aCards := make([]api.Card, len(mCards))
	for i, mCard := range mCards {
		aCards[i] = model.ToAPICard(mCard)
	}
	return aCards, nil
}

func UpdateCard(ctx context.Context, id string, updateData map[string]interface{}) (api.Card, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return api.Card{}, err
	}

	update := bson.M{"$set": updateData}
	_, err = cardsCollection.UpdateOne(ctx, bson.M{"_id": objID}, update)
	if err != nil {
		return api.Card{}, err
	}

	var updatedCard model.Card
	err = cardsCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&updatedCard)
	if err != nil {
		return api.Card{}, err
	}
	return model.ToAPICard(updatedCard), nil
}

func DeleteCard(ctx context.Context, id string) (int64, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return 0, err
	}

	result, err := cardsCollection.DeleteOne(ctx, bson.M{"_id": objID})
	return result.DeletedCount, err
}
