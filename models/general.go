package models

import (
	"mime/multipart"
	"time"
)

const MediaCollection = "medias"

type Response struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type File struct {
	File      *multipart.FileHeader   `form:"file"`
	MultiFile []*multipart.FileHeader `form:"file"`
}

type Media struct {
	Id  string `json:"id"`
	Url string `json:"url"`

	CreatedAt time.Time `json:"createdAt"`
	CreatedBy string    `json:"createdBy"`
}

type BasicDate struct {
	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt"`
}

type AuthResult struct {
	Token string `json:"token"`
	Id    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type DefaultData struct {
	Users       []User       `json:"users"`
	Permissions []Permission `json:"permissions"`
	Roles       []Role       `json:"roles"`
	Settings    []Setting    `json:"settings"`
}
