package models

import (
	"time"
)

type User struct {
	Id            string    `json:"id" bson:"id"`
	RoleId        string    `json:"roleId" bson:"roleId"`
	Name          string    `json:"name" bson:"name"`
	Email         string    `json:"email" bson:"email"`
	Password      string    `json:"password" bson:"password"`
	Image         string    `json:"image" bson:"image"`
	Active        bool      `json:"active" bson:"active" default:"false"`
	LastActive    time.Time `json:"lastActive" bson:"lastActive"`
	IsAdminSystem bool      `json:"isAdminSystem" bson:"isAdminSystem" default:"false"`
	CurrentSocial string    `json:"currentSocial" bson:"currentSocial"` // userSocialId
	UserInfo      UserInfo  `json:"info" bson:"info"`
	Status        bool      `json:"status" default:"true"` // Enable or disable by admin
	Role          Role      `json:"role" bson:"-"`
	Permissions   []string  `json:"permissions" bson:"-"`
	BasicDate     `bson:",inline"`
}

type UserSocial struct {
	Id         string      `json:"id"`
	UserId     string      `json:"userId"`
	Type       string      `json:"type"` // google, github, etc
	SocialId   string      `json:"socialId"`
	Additional interface{} `json:"additional"`
}

type UserInfo struct {
	Phone   string `json:"phone" bson:"phone"`
	Country string `json:"country" bson:"country"`
	State   string `json:"state" bson:"state"`
	City    string `json:"city" bson:"city"`
	Address string `json:"address" bson:"address"`
	ZipCode string `json:"zipCode"`
}

type AuthRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
	IsSocial bool   `json:"isSocial"`
	Social   string `json:"social"` // google, facebook, etc
	SocialId string `json:"socialId"`
	Image    string `json:"image"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email"`
}

type UpdatePasswordRequest struct {
	Token    string `json:"token"`
	Password string `json:"password"`
}

type UserRequest struct {
	Active     bool        `json:"active"`
	Email      string      `json:"email"`
	Image      string      `json:"image"`
	Main       bool        `json:"main"`
	Name       string      `json:"name"`
	Owner      bool        `json:"owner"`
	Password   string      `json:"password"`
	Permission interface{} `json:"permission"`
	Status     bool        `json:"status"`
}

type LoginResult struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Token string `json:"token"`
}
