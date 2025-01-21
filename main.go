package main

import (
	"apost/config"
	"apost/controllers"
	"apost/database"
	"apost/models"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/jasonlvhit/gocron"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file, will use system environment")
	}

	log.Println("Version: ", os.Getenv("VERSION"))

	if !database.Init() {
		log.Printf("Connected to MongoDB URI: Failure")
		return
	}

	router := gin.New()
	router.Use(gin.Logger())

	router.Use(static.Serve("/", static.LocalFile("./dist", true)))

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AllowMethods = []string{"POST", "GET", "PATCH", "OPTIONS", "DELETE"}
	corsConfig.AllowHeaders = []string{"Origin", "Content-Type", "Authorization", "Accept", "User-Agent", "Cache-Control", "Pragma"}
	corsConfig.ExposeHeaders = []string{"Content-Length"}
	corsConfig.AllowCredentials = true
	corsConfig.MaxAge = 12 * time.Hour
	router.Use(cors.New(corsConfig))

	api := router.Group("/api")
	{
		api.GET("/version", func(c *gin.Context) {
			c.JSON(http.StatusOK, models.Response{
				Data: "APost Api v" + os.Getenv("VERSION"),
			})
			return
		})

		api.GET("/default", controllers.CreateDefaultData)

		api.POST("/register", controllers.SignUp)
		api.POST("/login", controllers.SignIn)
		api.GET("/refresh-token", controllers.RefreshToken)
		api.POST("/activate", controllers.Activate)
		api.POST("/forgot-password", controllers.ForgotPassword)
		api.PATCH("/update-password", controllers.UpdatePassword)

		protected := api.Group("/", config.AuthMiddleware())
		{
			protected.GET("/profile", controllers.GetProfile)
			protected.PATCH("/profile", controllers.UpdateProfile)

			protected.GET("/app", controllers.GetApps)
			protected.POST("/app", controllers.CreateManyApp)
			protected.GET("/app/:id", controllers.GetAppById)
			protected.PATCH("/app/:id", controllers.UpdateApp)
			protected.DELETE("/app/:id", controllers.DeleteApp)

			protected.GET("/role", controllers.GetRoles)
			protected.POST("/role", controllers.CreateRole)
			protected.GET("/role/:id", controllers.GetRoleById)
			protected.PATCH("/role/:id", controllers.UpdateRole)
			protected.DELETE("/role/:id", controllers.DeleteRole)
			protected.POST("/role/attach-permission", controllers.AttachPermissionsToRole)

			protected.GET("/topic", controllers.GetTopics)
			protected.POST("/topic", controllers.CreateTopic)
			protected.GET("/topic/:id", controllers.GetTopicById)
			protected.PATCH("/topic/:id", controllers.UpdateTopic)
			protected.DELETE("/topic/:id", controllers.DeleteTopic)

			protected.GET("/user", controllers.GetUsers)
			protected.POST("/user", controllers.CreateUser)
			protected.GET("/user/:id", controllers.GetUserById)
			protected.PATCH("/user/:id", controllers.UpdateUser)
			protected.DELETE("/user/:id", controllers.DeleteUser)
		}
	}

	port := "8000"
	if os.Getenv("PORT") != "" {
		port = os.Getenv("PORT")
	}
	go func() {
		err := router.Run(":" + port)
		if err != nil {
			return
		}
	}()

	s := gocron.NewScheduler()
	//s.Every(5).Seconds().Do(controllers.RunAutoPost)
	log.Println(time.Now())
	//s.Every(10).Seconds().Do(controllers.RunAutoPost)
	s.Every(1).Tuesday().At("09:00").Do(controllers.RunAutoPost)
	<-s.Start()
}
