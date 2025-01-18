package models

type Post struct {
	Id        string `json:"id"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	Slug      string `json:"slug"`
	BasicDate `bson:",inline"`
}
