package controllers

import (
	"apost/models"
	"apost/services"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"io"
	"log"
	"net/http"
	"os"
)

func CreateDefaultData(c *gin.Context) {
	jsonFile, err := os.Open("init.json")

	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Successfully Opened users.json")

	byteValue, _ := io.ReadAll(jsonFile)

	var data models.DefaultData
	json.Unmarshal(byteValue, &data)

	users := services.GetUsersWithPagination(bson.M{}, nil, models.Query{})
	if users.Pagination.Count == 0 {
		_, err2 := services.CreateManyUser(data.Users)
		if err2 == nil {
			log.Println("Default User Has Been Created")
		}
	}

	settings := services.GetSettingsWithPagination(bson.M{}, nil, models.Query{})
	if settings.Pagination.Count == 0 {
		_, err2 := services.CreateManySettings(data.Settings)
		if err2 == nil {
			log.Println("Default Setting Has Been Created")
		}
	}

	roles := services.GetRolesWithPagination(bson.M{}, nil, models.Query{})
	if roles.Pagination.Count == 0 {
		_, err2 := services.CreateManyRole(data.Roles)
		if err2 == nil {
			log.Println("Default Role Has Been Created")
		}
	}

	permission := services.GetPermission(bson.M{}, nil)
	if permission == nil {
		_, err2 := services.CreateManyPermissions(data.Permissions)
		if err2 == nil {
			log.Println("Default Permission Has Been Created")
		}
	}

	c.JSON(http.StatusOK, models.Response{Data: "Success"})
}
