package services

import (
	"apost/database"
	"apost/models"
	"context"
	"encoding/json"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const SettingCollection = "settings"

func GetSettings(filters bson.M, opt *options.FindOptions) []models.Setting {
	results := make([]models.Setting, 0)

	cursor := database.Find(SettingCollection, filters, opt)
	for cursor.Next(context.Background()) {
		var data models.Setting
		if cursor.Decode(&data) == nil {
			results = append(results, data)
		}
	}

	return results
}

func GetSettingsWithPagination(filters bson.M, opt *options.FindOptions, query models.Query) models.Result {
	results := GetSettings(filters, opt)

	count := database.Count(SettingCollection, filters)

	pagination := query.GetPagination(count)

	result := models.Result{
		Data:       results,
		Pagination: pagination,
		Query:      query,
	}

	return result
}

func CreateSetting(Setting models.Setting) (bool, error) {
	_, err := database.InsertOne(SettingCollection, Setting)
	if err != nil {
		return false, err
	}

	return true, nil
}

func CreateManySettings(params []models.Setting) (bool, error) {
	data := make([]interface{}, 0)
	for _, val := range params {
		data = append(data, val)
	}

	_, err := database.InsertMany(SettingCollection, data)
	if err != nil {
		return false, err
	}

	return true, nil
}

func GetSetting(filter bson.M, opts *options.FindOneOptions) *models.Setting {
	var Setting models.Setting
	err := database.FindOne(SettingCollection, filter, opts).Decode(&Setting)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil
		}
		return nil
	}
	return &Setting
}

func UpdateSetting(id string, Setting interface{}) (*mongo.UpdateResult, error) {
	filters := bson.M{"id": id}

	res, err := database.UpdateOne(SettingCollection, filters, Setting)

	if res == nil {
		return nil, err
	}

	return res, nil
}

func DeleteSetting(ids []string) (*mongo.DeleteResult, error) {
	filter := bson.M{"id": ids}

	res, err := database.DeleteMany(SettingCollection, filter)

	if res == nil {
		return nil, err
	}

	return res, nil
}

func GetCloudinaryConfig() *models.CloudinaryInstanceSetting {
	setting := GetSetting(bson.M{"code": "cloudinary"}, nil)
	if setting == nil {
		return nil
	}

	var res models.CloudinaryInstanceSetting
	a, _ := json.Marshal(setting)
	err := json.Unmarshal(a, &res)
	if err != nil {
		return nil
	}

	if res.CloudinaryCloudName == "" && res.CloudinarySecret == "" {
		return nil
	}

	return &res
}
