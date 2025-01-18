package controllers

import (
	"apost/models"
	"apost/services"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"net/http"
)

func GetApps(c *gin.Context) {
	var query models.Query

	err := c.ShouldBindQuery(&query)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Response{Data: err.Error()})
		return
	}

	filters := query.GetQueryFind()

	opts := query.GetOptions()
	opts.SetProjection(bson.M{"password": 0})

	results := services.GetAppsWithPagination(filters, opts, query)

	c.JSON(http.StatusOK, models.Response{Data: results})
	return
}

func CreateManyApp(c *gin.Context) {
	request := struct {
		Apps []models.App `json:"apps"`
	}{}
	err := c.ShouldBindJSON(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Response{Data: err.Error()})
		return
	}

	_, err = services.CreateManyApp(request.Apps)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Response{Data: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.Response{Data: request})
	return
}

func GetAppById(c *gin.Context) {
	id := c.Param("id")

	result := services.GetApp(bson.M{"id": id}, nil)
	if result == nil {
		c.JSON(http.StatusNotFound, models.Result{Data: "Data Not Found"})
		return
	}

	c.JSON(http.StatusOK, models.Response{Data: result})
}

func UpdateApp(c *gin.Context) {
	id := c.Param("id")

	App := services.GetApp(bson.M{"id": id}, nil)
	if App == nil {
		c.JSON(http.StatusNotFound, models.Result{Data: "Data Not Found"})
		return
	}

	var request models.App

	err := c.ShouldBindJSON(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Response{Data: err.Error()})
		return
	}

	_, err = services.UpdateApp(id, request)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Response{Data: err.Error()})
		return
	}

	c.JSON(200, models.Response{Data: request})
}

func DeleteApp(c *gin.Context) {
	id := c.Param("id")

	_, err := services.DeleteApp(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Response{Data: "Failed Delete Data"})
		return
	}

	c.JSON(http.StatusOK, models.Response{Data: "Success"})
}
