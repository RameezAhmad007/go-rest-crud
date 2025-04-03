package repository

import (
	"context"
	"time"

	"github.com/RameezAhmad007/go-rest-crud/internal/config"
	"github.com/RameezAhmad007/go-rest-crud/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var cardsCollection = config.GetCollection("cardscollection")

func CreateCard(ctx context.Context, card model.Card) (model.Card, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Check for duplicate name
	var existingCard model.Card
	err := cardsCollection.FindOne(ctx, bson.M{"name": card.Name}).Decode(&existingCard)
	if err == nil {
		return model.Card{}, mongo.ErrClientDisconnected // Simulate conflict
	} else if err != mongo.ErrNoDocuments {
		return model.Card{}, err
	}

	card.ID = primitive.NewObjectID()
	_, err = cardsCollection.InsertOne(ctx, card)
	return card, err
}

func GetCard(ctx context.Context, id string) (model.Card, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return model.Card{}, err
	}

	var card model.Card
	err = cardsCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&card)
	return card, err
}

func GetAllCards(ctx context.Context) ([]model.Card, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var cards []model.Card
	cursor, err := cardsCollection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	err = cursor.All(ctx, &cards)
	return cards, err
}

func UpdateCard(ctx context.Context, id string, updateData map[string]interface{}) (model.Card, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return model.Card{}, err
	}

	update := bson.M{"$set": updateData}
	_, err = cardsCollection.UpdateOne(ctx, bson.M{"_id": objID}, update)
	if err != nil {
		return model.Card{}, err
	}

	var updatedCard model.Card
	err = cardsCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&updatedCard)
	return updatedCard, err
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
