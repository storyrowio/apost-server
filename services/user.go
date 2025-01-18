package services

import (
	"apost/database"
	"apost/lib"
	"apost/models"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net/http"
	"time"
)

const UserCollection = "users"
const UserSocialCollection = "usersocials"

func GetUsers(filters bson.M, opt *options.FindOptions) []models.User {
	results := make([]models.User, 0)

	cursor := database.Find(UserCollection, filters, opt)

	if cursor == nil {
		return results
	}
	for cursor.Next(context.Background()) {
		var user models.User

		if cursor.Decode(&user) == nil {
			results = append(results, user)
		}
	}
	return results
}

func GetUsersWithPagination(filters bson.M, opt *options.FindOptions, query models.Query) models.Result {
	results := GetUsers(filters, opt)

	count := database.Count(UserCollection, filters)

	pagination := query.GetPagination(count)

	result := models.Result{
		Data:       results,
		Pagination: pagination,
		Query:      query,
	}

	return result
}

func CreateUser(user models.User) (bool, error) {
	_, err := database.InsertOne(UserCollection, user)
	if err != nil {
		return false, err
	}

	return true, nil
}

func GetUser(filter bson.M, opts *options.FindOneOptions) *models.User {
	cursor := database.FindOne(UserCollection, filter, opts)
	if cursor == nil {
		return nil
	}

	var user models.User
	err := cursor.Decode(&user)
	if err != nil {
		return nil
	}

	return &user
}

func UpdateUser(id string, user interface{}) (*mongo.UpdateResult, error) {
	filters := bson.M{"id": id}

	res, err := database.UpdateOne(UserCollection, filters, user)

	if res == nil {
		return nil, err
	}

	return res, nil
}

func DeleteUser(id string) (*mongo.DeleteResult, error) {
	filter := bson.M{"id": id}

	res, err := database.DeleteOne(UserCollection, filter)

	if res == nil {
		return nil, err
	}

	return res, nil
}

func GetCurrentUser(r *http.Request) *models.User {
	token, e := CheckHeader(r)

	if e != nil {
		return nil
	}

	email, err := VerifyToken(token)

	if err != nil {
		return nil
	}

	user := GetUser(bson.M{"email": email}, options.FindOne().SetProjection(bson.D{{"password", 0}}))

	return user
}

func CreateUserSocial(userSocial models.UserSocial) (bool, error) {
	_, err := database.InsertOne(UserSocialCollection, userSocial)
	if err != nil {
		return false, err
	}

	return true, nil
}

func GetUserSocials(userId string) []models.UserSocial {
	cursor := database.Find(UserCollection, bson.M{"userid": userId}, nil)

	results := make([]models.UserSocial, 0)
	for cursor.Next(context.Background()) {
		var data models.UserSocial
		if cursor.Decode(&data) == nil {
			results = append(results, data)
		}
	}

	return results
}

func CreateManyUser(params []models.User) (bool, error) {
	data := make([]interface{}, 0)
	for _, val := range params {
		val.Password = lib.HashAndSalt(val.Password)
		val.CreatedAt = time.Now()
		val.UpdatedAt = time.Now()
		data = append(data, val)
	}

	_, err := database.InsertMany(UserCollection, data)
	if err != nil {
		return false, err
	}

	return true, nil
}
