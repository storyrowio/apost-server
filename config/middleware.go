package config

import (
	"apost/models"
	"apost/services"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"strings"
)

func CheckToken(tokenString string) (string, error) {
	claims := &models.Claims{}

	_, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET_KEY")), nil
	})

	if err != nil {
		return "", errors.New("token invalid")
	}

	return claims.Email, nil
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.Request.Header["Authorization"]
		apiKeyHeader := c.Request.Header.Get("ApiKey")

		if header == nil && apiKeyHeader == "" {
			c.Abort()
			c.Writer.WriteHeader(http.StatusUnauthorized)
			_, err := c.Writer.Write([]byte("unauthorized"))
			if err != nil {
				return
			}
			return
		}

		if header != nil {
			split := strings.Split(header[0], " ")
			if strings.Contains(header[0], "Bearer") {
				if len(split) != 2 || strings.ToLower(split[0]) != "bearer" {
					c.Abort()
					c.Writer.WriteHeader(http.StatusUnauthorized)
					_, err := c.Writer.Write([]byte("bearer token format needed"))
					if err != nil {
						return
					}
					return
				}

				_, err := CheckToken(split[1])
				if err != nil {
					c.Abort()
					c.Writer.WriteHeader(http.StatusUnauthorized)
					_, err := c.Writer.Write([]byte("token invalid"))
					if err != nil {
						return
					}
					return
				}
			}
		} else if apiKeyHeader != "" {
			apiKey := services.GetApiKey(apiKeyHeader)
			if apiKey == nil {
				c.Abort()
				c.Writer.WriteHeader(http.StatusUnauthorized)
				_, err := c.Writer.Write([]byte("api key is invalid"))
				if err != nil {
					return
				}
				return
			}
		}

	}
}

func CorsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, PUT, PATCH, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
