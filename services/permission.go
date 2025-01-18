package services

import (
	"apost/database"
	"apost/models"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const PermissionCollection = "permissions"

func CreateManyPermissions(params []models.Permission) (bool, error) {
	data := make([]interface{}, 0)
	for _, val := range params {
		data = append(data, val)
	}

	_, err := database.InsertMany(PermissionCollection, data)
	if err != nil {
		return false, err
	}

	return true, nil
}

func GetPermission(filter bson.M, opts *options.FindOneOptions) *models.Permission {
	var data models.Permission
	err := database.FindOne(PermissionCollection, filter, opts).Decode(&data)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil
		}
		return nil
	}
	return &data
}
