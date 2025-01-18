package controllers

import (
	"apost/lib"
	"apost/models"
	"apost/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net/http"
	"time"
)

func GetUsers(c *gin.Context) {
	var query models.Query

	err := c.ShouldBindQuery(&query)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Response{Data: err.Error()})
		return
	}

	filters := query.GetQueryFind()
	opts := query.GetOptions()
	opts.SetProjection(bson.M{"password": 0})

	results := services.GetUsersWithPagination(filters, opts, query)

	c.JSON(http.StatusOK, models.Response{Data: results})
	return
}

func CreateUser(c *gin.Context) {
	var request models.User

	err := c.ShouldBindJSON(&request)
	if err != nil {
		c.JSON(400, models.Response{Data: err.Error()})
		return
	}

	request.Id = uuid.New().String()
	request.Password = lib.HashAndSalt(request.Password)
	request.CreatedAt = time.Now()

	_, err = services.CreateUser(request)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Response{Data: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.Response{Data: request})
	return
}

func GetUserById(c *gin.Context) {
	id := c.Param("id")

	result := services.GetUser(bson.M{"id": id}, options.FindOne().SetProjection(bson.M{"password": 0}))
	if result == nil {
		c.JSON(http.StatusNotFound, models.Result{Data: "Data Not Found"})
		return
	}

	c.JSON(http.StatusOK, models.Response{Data: result})
}

func UpdateUser(c *gin.Context) {
	id := c.Param("id")

	user := services.GetUser(bson.M{"id": id}, nil)
	if user == nil {
		c.JSON(http.StatusNotFound, models.Result{Data: "Data Not Found"})
		return
	}

	var request models.User
	err := c.ShouldBindJSON(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Response{Data: err.Error()})
		return
	}

	if request.Password != "" {
		request.Password = lib.HashAndSalt(request.Password)
	} else {
		request.Password = user.Password
	}

	_, err = services.UpdateUser(id, request)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Response{Data: err.Error()})
		return
	}

	c.JSON(200, models.Response{Data: request})
}

func DeleteUser(c *gin.Context) {
	id := c.Param("id")

	_, err := services.DeleteUser(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Response{Data: "Failed Delete Data"})
		return
	}

	c.JSON(http.StatusOK, models.Response{Data: "Success"})
}
