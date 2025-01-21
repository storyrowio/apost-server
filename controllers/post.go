package controllers

import (
	"apost/models"
	"apost/services"
	"go.mongodb.org/mongo-driver/bson"
	"log"
)

func RunAutoPost() {
	log.Println("Run auto post")
	apps := services.GetApps(bson.M{}, nil)

	for _, app := range apps {
		topic := services.GetTopicForAutoPost(app.Id)
		res, err := services.GeneratePost(topic.Topic)
		if err != nil {
			log.Println("Error Generate Post:", err.Error())
		}

		if len(res.Choices) > 0 && res.Choices[0].Message.Content != "" {
			_, err = services.SendPostToApp(app, models.Post{
				Title:   topic.Topic,
				Content: res.Choices[0].Message.Content,
			})

			_, err := services.DeleteTopic(topic.Id)
			if err != nil {
				log.Println("Error Delete Topic:", err.Error())
			}

			log.Println("Create post for " + app.Name + ": success")
		} else {
			log.Println("Create post for " + app.Name + ": failed")
		}
	}

	//c.JSON(http.StatusOK, res)
}
