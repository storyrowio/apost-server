package models

import (
	"github.com/dgrijalva/jwt-go"
)

type Claims struct {
	UserId 			string 		`json:"userId"`
	Email  			string   	`json:"email"`
	jwt.StandardClaims
}

type TokenEndExpire struct {
	Token  string
	Expire string
}

