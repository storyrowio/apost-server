package services

import (
	"apost/database"
	"apost/lib"
	"apost/models"
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const TopicCollection = "topics"

func GetTopics(filters bson.M, opt *options.FindOptions) []models.Topic {
	results := make([]models.Topic, 0)

	cursor := database.Find(TopicCollection, filters, opt)
	for cursor.Next(context.Background()) {
		var data models.Topic
		if cursor.Decode(&data) == nil {
			results = append(results, data)
		}
	}

	return results
}

func GetTopicsWithPagination(filters bson.M, opt *options.FindOptions, query models.Query) models.Result {
	results := GetTopics(filters, opt)

	count := database.Count(TopicCollection, filters)

	pagination := query.GetPagination(count)

	result := models.Result{
		Data:       results,
		Pagination: pagination,
		Query:      query,
	}

	return result
}

func CreateTopic(params models.Topic) (bool, error) {
	_, err := database.InsertOne(TopicCollection, params)
	if err != nil {
		return false, err
	}

	return true, nil
}

func GetTopic(filter bson.M, opts *options.FindOneOptions) *models.Topic {
	var data models.Topic
	err := database.FindOne(TopicCollection, filter, opts).Decode(&data)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil
		}
		return nil
	}
	return &data
}

func UpdateTopic(id string, Topic interface{}) (*mongo.UpdateResult, error) {
	filters := bson.M{"id": id}

	res, err := database.UpdateOne(TopicCollection, filters, Topic)

	if res == nil {
		return nil, err
	}

	return res, nil
}

func DeleteTopic(id string) (*mongo.DeleteResult, error) {
	filter := bson.M{"id": id}

	res, err := database.DeleteOne(TopicCollection, filter)

	if res == nil {
		return nil, err
	}

	return res, nil
}

func CreateManyTopic(params []models.Topic) (bool, error) {
	data := make([]interface{}, 0)
	for _, val := range params {
		data = append(data, val)
	}

	_, err := database.InsertMany(TopicCollection, data)
	if err != nil {
		return false, err
	}

	return true, nil
}

func GetTopicForAutoPost(appId string) models.Topic {
	topics := GetTopics(bson.M{"appId": appId}, nil)
	numb := lib.RandomNumber(0, len(topics)-1)
	return topics[numb]
}
