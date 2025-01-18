package controllers

import (
	"apost/models"
	"apost/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"time"
)

func GetRoles(c *gin.Context) {
	var query models.Query

	err := c.ShouldBindQuery(&query)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Response{Data: err.Error()})
		return
	}

	filters := query.GetQueryFind()

	keyword, keywordExist := c.GetQuery("keyword")
	if keywordExist {
		regex := primitive.Regex{
			Pattern: "^" + keyword,
			Options: "i",
		}
		filters["$or"] = bson.A{
			bson.M{"name": regex},
			bson.M{"email": regex},
		}
	}

	opts := query.GetOptions()
	opts.SetProjection(bson.M{"password": 0})

	results := services.GetRolesWithPagination(filters, opts, query)

	c.JSON(http.StatusOK, models.Response{Data: results})
	return
}

func CreateRole(c *gin.Context) {
	var request models.Role

	err := c.ShouldBindJSON(&request)
	if err != nil {
		c.JSON(400, models.Response{Data: err.Error()})
		return
	}

	request.Id = uuid.New().String()
	request.CreatedAt = time.Now()
	request.UpdatedAt = time.Now()

	_, err = services.CreateRole(request)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Response{Data: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.Response{Data: request})
	return
}

func GetRoleById(c *gin.Context) {
	id := c.Param("id")

	result := services.GetRole(bson.M{"id": id}, nil)
	if result == nil {
		c.JSON(http.StatusNotFound, models.Result{Data: "Data Not Found"})
		return
	}

	c.JSON(http.StatusOK, models.Response{Data: result})
}

func UpdateRole(c *gin.Context) {
	id := c.Param("id")

	Role := services.GetRole(bson.M{"id": id}, nil)
	if Role == nil {
		c.JSON(http.StatusNotFound, models.Result{Data: "Data Not Found"})
		return
	}

	var request models.Role

	err := c.ShouldBindJSON(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Response{Data: err.Error()})
		return
	}

	_, err = services.UpdateRole(id, request)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Response{Data: err.Error()})
		return
	}

	c.JSON(200, models.Response{Data: request})
}

func DeleteRole(c *gin.Context) {
	id := c.Param("id")

	_, err := services.DeleteRole(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Response{Data: "Failed Delete Data"})
		return
	}

	c.JSON(http.StatusOK, models.Response{Data: "Success"})
}

func AttachPermissionsToRole(c *gin.Context) {
	request := struct {
		RoleId      string   `json:"roleId"`
		Permissions []string `json:"permissions"`
	}{}
	err := c.ShouldBindJSON(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Response{Data: err.Error()})
		return
	}

	role := services.GetRole(bson.M{"id": request.RoleId}, nil)
	if role == nil {
		c.JSON(http.StatusNotFound, models.Response{Data: err.Error()})
		return
	}

	role.UpdatedAt = time.Now()
	role.Permissions = append(role.Permissions, request.Permissions...)
	_, err = services.UpdateRole(role.Id, role)
	if err != nil {
		c.JSON(http.StatusNotFound, models.Response{Data: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.Response{Data: role})
}
