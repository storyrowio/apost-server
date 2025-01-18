package services

import (
	"apost/database"
	"apost/lib"
	"apost/models"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net/http"
	"time"
)

const AppCollection = "apps"
const ApiKeyCollection = "api_keys"

func GetApps(filters bson.M, opt *options.FindOptions) []models.App {
	results := make([]models.App, 0)

	cursor := database.Find(AppCollection, filters, opt)
	for cursor.Next(context.Background()) {
		var data models.App
		if cursor.Decode(&data) == nil {
			results = append(results, data)
		}
	}

	return results
}

func GetAppsWithPagination(filters bson.M, opt *options.FindOptions, query models.Query) models.Result {
	results := GetApps(filters, opt)

	count := database.Count(AppCollection, filters)

	pagination := query.GetPagination(count)

	result := models.Result{
		Data:       results,
		Pagination: pagination,
		Query:      query,
	}

	return result
}

func CreateApp(params models.App) (bool, error) {
	_, err := database.InsertOne(AppCollection, params)
	if err != nil {
		return false, err
	}

	return true, nil
}

func GetApp(filter bson.M, opts *options.FindOneOptions) *models.App {
	var data models.App
	err := database.FindOne(AppCollection, filter, opts).Decode(&data)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil
		}
		return nil
	}
	return &data
}

func UpdateApp(id string, App interface{}) (*mongo.UpdateResult, error) {
	filters := bson.M{"id": id}

	res, err := database.UpdateOne(AppCollection, filters, App)

	if res == nil {
		return nil, err
	}

	return res, nil
}

func DeleteApp(id string) (*mongo.DeleteResult, error) {
	filter := bson.M{"id": id}

	res, err := database.DeleteOne(AppCollection, filter)

	if res == nil {
		return nil, err
	}

	return res, nil
}

func CreateManyApp(params []models.App) (bool, error) {
	data := make([]interface{}, 0)
	apiKeys := make([]interface{}, 0)

	for _, val := range params {
		val.Id = uuid.New().String()
		val.CreatedAt = time.Now()
		val.UpdatedAt = time.Now()
		data = append(data, val)

		apiKey, _ := lib.GenerateAPIKey(30)
		apiKeyData := models.ApiKey{
			Id:     uuid.New().String(),
			AppId:  val.Id,
			ApiKey: apiKey,
			BasicDate: models.BasicDate{
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		}
		apiKeys = append(apiKeys, apiKeyData)
	}

	_, err := database.InsertMany(AppCollection, data)
	if err != nil {
		return false, err
	}

	_, err = database.InsertMany(ApiKeyCollection, apiKeys)
	if err != nil {
		return false, err
	}

	return true, nil
}

func GetCurrentApp(r *http.Request) *models.App {
	key, err := CheckHeaderApiKey(r)
	if err != nil {
		return nil
	}

	apiKey := GetApiKey(key)
	if apiKey == nil {
		return nil
	}

	app := GetApp(bson.M{"id": apiKey.AppId}, nil)

	return app
}

func SendPostToApp(app models.App, post models.Post) (*map[string]interface{}, error) {
	var err error
	var client = &http.Client{}
	var data map[string]interface{}

	jsonMarshal, err := json.Marshal(post)

	bodyReader := bytes.NewReader(jsonMarshal)

	request, err := http.NewRequest("POST", app.ApiUrl+app.PostEndpoint, bodyReader)
	if err != nil {
		return nil, err
	}

	request.Header.Set("Content-Type", "application/json")

	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	err = json.NewDecoder(response.Body).Decode(&data)
	if err != nil {
		return nil, err
	}

	return &data, nil
}
