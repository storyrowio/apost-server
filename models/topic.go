package models

type Topic struct {
	Id        string `json:"id"`
	AppId     string `json:"appId" bson:"appId"`
	Topic     string `json:"topic"`
	BasicDate `bson:",inline"`
}
