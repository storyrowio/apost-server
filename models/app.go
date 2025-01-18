package models

type App struct {
	Id           string `json:"id"`
	Name         string `json:"name"`
	BaseUrl      string `json:"baseUrl" bson:"baseUrl"`
	ApiUrl       string `json:"apiUrl" bson:"apiUrl"`
	PostEndpoint string `json:"postEndpoint" bson:"postEndpoint"`
	BasicDate    `bson:",inline"`
}
