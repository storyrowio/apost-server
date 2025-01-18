package controllers

import (
	"apost/models"
	"apost/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"net/http"
	"time"
)

func GetTopics(c *gin.Context) {
	var query models.Query

	err := c.ShouldBindQuery(&query)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Response{Data: err.Error()})
		return
	}

	filters := query.GetQueryFind()

	opts := query.GetOptions()
	opts.SetProjection(bson.M{"password": 0})

	results := services.GetTopicsWithPagination(filters, opts, query)

	c.JSON(http.StatusOK, models.Response{Data: results})
	return
}

func CreateTopic(c *gin.Context) {
	request := struct {
		Topics []string `json:"topics"`
	}{}
	err := c.ShouldBindJSON(&request)
	if err != nil {
		c.JSON(400, models.Response{Data: err.Error()})
		return
	}

	currentApp := services.GetCurrentApp(c.Request)
	newData := make([]models.Topic, 0)
	for _, val := range request.Topics {
		data := models.Topic{
			Id:    uuid.New().String(),
			AppId: currentApp.Id,
			Topic: val,
			BasicDate: models.BasicDate{
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		}

		newData = append(newData, data)
	}

	_, err = services.CreateManyTopic(newData)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Response{Data: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.Response{Data: request})
	return
}

func GetTopicById(c *gin.Context) {
	id := c.Param("id")

	result := services.GetTopic(bson.M{"id": id}, nil)
	if result == nil {
		c.JSON(http.StatusNotFound, models.Result{Data: "Data Not Found"})
		return
	}

	c.JSON(http.StatusOK, models.Response{Data: result})
}

func UpdateTopic(c *gin.Context) {
	id := c.Param("id")

	Topic := services.GetTopic(bson.M{"id": id}, nil)
	if Topic == nil {
		c.JSON(http.StatusNotFound, models.Result{Data: "Data Not Found"})
		return
	}

	var request models.Topic

	err := c.ShouldBindJSON(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Response{Data: err.Error()})
		return
	}

	_, err = services.UpdateTopic(id, request)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Response{Data: err.Error()})
		return
	}

	c.JSON(200, models.Response{Data: request})
}

func DeleteTopic(c *gin.Context) {
	id := c.Param("id")

	_, err := services.DeleteTopic(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Response{Data: "Failed Delete Data"})
		return
	}

	c.JSON(http.StatusOK, models.Response{Data: "Success"})
}
