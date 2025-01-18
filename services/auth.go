package services

import (
	"apost/database"
	"apost/lib"
	"apost/models"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

func GenerateToken(email string) (*string, error) {
	expire := time.Now().Add(24 * time.Hour)

	claims := models.Claims{
		Email: email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expire.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET_KEY")))

	if err == nil {
		return &tokenString, nil
	}

	return nil, errors.New(err.Error())
}

func CheckHeader(r *http.Request) (string, error) {
	header := r.Header["Authorization"]

	if header == nil {
		return "", errors.New("unauthorized")
	}

	split := strings.Split(header[0], " ")
	if len(split) != 2 || strings.ToLower(split[0]) != "bearer" {
		return "", errors.New("unauthorized")
	}

	return split[1], nil
}

func CheckHeaderApiKey(r *http.Request) (string, error) {
	apiKeyHeader := r.Header.Get("ApiKey")
	if apiKeyHeader == "" {
		return "", errors.New("unauthorized")
	}

	return apiKeyHeader, nil
}

func VerifyToken(tokenString string) (string, error) {
	claims := &models.Claims{}

	_, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET_KEY")), nil
	})

	if err != nil {
		log.Println(err)
		return "", errors.New("Token invalid")
	}

	return claims.Email, nil
}

func AuthWithSocial(user *models.User, request models.AuthRequest) error {
	isExist := false
	if user != nil {
		userSocial := GetUserSocials(user.Id)
		for _, val := range userSocial {
			if val.Type == request.Social && val.SocialId == request.SocialId {
				isExist = true
			}
		}
	}

	if !isExist {
		roleId := ""
		role := GetRole(bson.M{"code": models.AdminRole}, nil)
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

		if user == nil {
			newUser := models.User{
				Id:            newUserId,
				RoleId:        roleId,
				Name:          request.Name,
				Email:         request.Email,
				Image:         request.Image,
				Active:        true,
				LastActive:    time.Now(),
				CurrentSocial: userSocial.Id,
				BasicDate: models.BasicDate{
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
			}
			_, err := CreateUser(newUser)
			if err != nil {
				return err
			}
		} else {
			if user.Image == "" {
				user.Image = request.Image
			}
			user.CurrentSocial = userSocial.Id
			user.Active = true
			user.LastActive = time.Now()
			user.UpdatedAt = time.Now()
			_, err := UpdateUser(user.Id, user)
			if err != nil {
				log.Println("Error Update User Auth Social:", err.Error())
			}
		}

		_, err := CreateUserSocial(userSocial)
		if err != nil {
			return err
		}
	}

	return nil
}

func Activate(token string) (bool, error) {
	email, err := VerifyToken(token)

	if err != nil {
		return false, err
	}

	filter := bson.M{"email": email}

	var user models.User
	cursor := database.FindOne(UserCollection, filter, nil)

	err = cursor.Decode(&user)
	if err != nil {
		return false, err
	}

	user.Active = true

	_, err = database.UpdateOne(UserCollection, filter, user)
	if err != nil {
		return false, err
	}

	return true, nil
}

func ForgotPassword(email string, url string) (bool, error) {
	user := GetUser(bson.M{"email": email}, nil)
	if user == nil {
		return false, errors.New("user not found")
	}

	token, err := GenerateToken(email)

	if err != nil {
		return false, err
	}

	tokenValue := *token

	data := models.VerificationMail{
		Name: user.Name,
		Link: url + tokenValue,
	}

	err = lib.SendEmailVerification("forgot-password", user.Email, data)
	if err != nil {
		return false, err
	}

	return true, nil
}

func UpdatePassword(token string, password string) (bool, error) {
	email, err := VerifyToken(token)

	if err != nil {
		return false, err
	}

	filter := bson.M{"email": email}

	var user models.User
	cursor := database.FindOne(UserCollection, filter, nil)

	err = cursor.Decode(&user)
	if err != nil {
		return false, err
	}

	user.Password = lib.HashAndSalt(password)

	_, err = database.UpdateOne(UserCollection, filter, user)
	if err != nil {
		return false, err
	}

	return true, nil
}
