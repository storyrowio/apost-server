package services

import (
	"apost/database"
	"apost/models"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetApiKey(apiKey string) *models.ApiKey {
	var data models.ApiKey
	err := database.FindOne(ApiKeyCollection, bson.M{"apiKey": apiKey}, nil).Decode(&data)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil
		}
		return nil
	}
	return &data
}
