package controllers

import (
	"apost/lib"
	"apost/models"
	"apost/services"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net/http"
	"os"
	"time"
)

func SignUp(c *gin.Context) {
	var request models.AuthRequest
	err := c.ShouldBindJSON(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Response{Data: err.Error()})
		return
	}

	user := services.GetUser(bson.M{"email": request.Email}, nil)

	if request.IsSocial {
		err := services.AuthWithSocial(user, request)
		if err != nil {
			c.JSON(http.StatusBadRequest, models.Response{Data: err.Error()})
			return
		}
	} else {
		if user != nil {
			c.JSON(http.StatusBadRequest, models.Response{Data: errors.New("Email already exist")})
			return
		}
		roleId := ""
		role := services.GetRole(bson.M{"code": models.AdminRole}, nil)
		if role != nil {
			roleId = role.Id
		}

		user = &models.User{
			Id:         uuid.New().String(),
			Active:     false,
			Email:      request.Email,
			Name:       request.Name,
			Password:   lib.HashAndSalt(request.Password),
			RoleId:     roleId,
			Status:     false,
			LastActive: time.Now(),
			BasicDate: models.BasicDate{
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		}

		_, err = services.CreateUser(*user)
		if err != nil {
			c.JSON(http.StatusBadRequest, models.Response{Data: err.Error()})
			return
		}
	}

	token, err := lib.GenerateToken(request.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Response{Data: err.Error()})
		return
	}

	result := models.LoginResult{
		Id:    user.Id,
		Name:  user.Name,
		Email: user.Email,
		Token: *token,
	}

	c.JSON(http.StatusOK, models.Response{Data: result})
	return
}

func SignIn(c *gin.Context) {
	var request models.AuthRequest
	err := c.ShouldBindJSON(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Response{Data: err.Error()})
		return
	}

	user := services.GetUser(bson.M{"email": request.Email}, nil)
	if user == nil {
		if request.IsSocial {
			roleId := ""
			role := services.GetRole(bson.M{"code": models.AdminRole}, nil)
			if role != nil {
				roleId = role.Id
			}

			newUserId := uuid.New().String()
			userSocial := models.UserSocial{
				Id:       uuid.New().String(),
				UserId:   newUserId,
				Type:     request.Social,
				SocialId: request.SocialId,
			}

			newUser := models.User{
				Id:            newUserId,
				RoleId:        roleId,
				Name:          request.Name,
				Email:         request.Email,
				Image:         request.Image,
				Active:        true,
				LastActive:    time.Now(),
				CurrentSocial: userSocial.Id,
				BasicDate:     models.BasicDate{},
			}
			_, err := services.CreateUser(newUser)
			if err != nil {
				c.JSON(http.StatusBadRequest, models.Response{Data: err.Error()})
				return
			}

			_, err = services.CreateUserSocial(userSocial)
			if err != nil {
				c.JSON(http.StatusBadRequest, models.Response{Data: err.Error()})
				return
			}
		} else {
			c.JSON(http.StatusBadRequest, models.Response{Message: "Invalid Credentials"})
			return
		}
	}

	if !request.IsSocial && request.Password != "" {
		pass := lib.ComparePassword(user.Password, []byte(request.Password))
		if !pass {
			c.JSON(http.StatusBadRequest, models.Response{Message: "Password does not match"})
			return
		}
	}

	token, err := lib.GenerateToken(request.Email)
	result := models.LoginResult{
		Id:    user.Id,
		Name:  user.Name,
		Email: user.Email,
		Token: *token,
	}

	c.JSON(http.StatusOK, models.Response{Data: result})
	return
}

func RefreshToken(c *gin.Context) {
	token, err := services.CheckHeader(c.Request)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Response{Data: err.Error()})
		return
	}

	email, err := services.VerifyToken(token)
	if err != nil {
		newToken, errGenerate := services.GenerateToken(email)
		if errGenerate != nil {
			c.JSON(http.StatusBadRequest, models.Response{Data: err.Error()})
			return
		}

		token = *newToken
	}

	user := services.GetUser(bson.M{"email": email}, options.FindOne().SetProjection(bson.M{"password": 0}))
	if user == nil {
		c.JSON(http.StatusBadRequest, models.Response{Data: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.Response{Data: token})
}

func Activate(c *gin.Context) {
	token := c.Param("token")

	_, err := services.Activate(token)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Response{Data: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.Response{Data: "Activation Success"})
	return
}

func ForgotPassword(c *gin.Context) {
	var request models.ForgotPasswordRequest
	err := c.ShouldBindJSON(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Response{Data: err.Error()})
		return
	}

	url := os.Getenv("FRONTEND_URL") + "/reset-password/"
	_, err = services.ForgotPassword(request.Email, url)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Response{Data: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.Response{Data: "Reset email has been sent"})
	return
}

func UpdatePassword(c *gin.Context) {
	var request models.UpdatePasswordRequest
	err := c.ShouldBindJSON(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Response{Data: err.Error()})
		return
	}

	_, err = services.UpdatePassword(request.Token, request.Password)
	if err != nil {
		c.JSON(http.StatusNotFound, models.Response{Data: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.Response{Data: "Success"})
	return
}

func GetProfile(c *gin.Context) {
	result := services.GetCurrentUser(c.Request)

	if result == nil {
		c.JSON(http.StatusNotFound, models.Response{Data: "User Not Found"})
		return
	}

	c.JSON(http.StatusOK, models.Response{Data: result})
	return
}

func UpdateProfile(c *gin.Context) {
	res := services.GetCurrentUser(c.Request)
	if res == nil {
		c.JSON(http.StatusNotFound, models.Response{Data: "User Not Found"})
		return
	}

	var request models.User

	err := c.ShouldBindJSON(&request)
	if err != nil {
		c.JSON(http.StatusNotFound, models.Response{Data: err.Error()})
		return
	}

	request.Id = res.Id
	request.LastActive = time.Now()

	_, err = services.UpdateUser(res.Id, request)
	if err != nil {
		c.JSON(http.StatusNotFound, models.Response{Data: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.Response{Data: "Success"})
	return
}
