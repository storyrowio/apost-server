package models

type ApiKey struct {
	Id        string `json:"id"`
	AppId     string `json:"appId" bson:"appId"`
	ApiKey    string `json:"apiKey" bson:"apiKey"`
	BasicDate `bson:",inline"`
}
