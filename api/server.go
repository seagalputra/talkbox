package api

import (
	"fmt"
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func StartServer() error {
	LoadAppConfig()

	if err := ConnectDatabase(); err != nil {
		log.Fatalf("[StartServer] %v", err)
	}
	ConnectToRedis()

	if err := LoadCloudinary(AppConfig.CloudinaryCloudName, AppConfig.CloudinaryAPIKey, AppConfig.CloudinaryAPISecret); err != nil {
		log.Fatalf("[StartServer] %v", err)
	}

	userHandler := UserDefaultHandler()
	messageHandler := MessageDefaultHandler()
	roomHandler := RoomDefaultHandler()
	r := gin.Default()

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{
		"http://127.0.0.1:3000",
		"http://localhost:3000",
		"https://talkbox.fly.dev",
		"http://talkbox.fly.dev",
	}
	corsConfig.AllowCredentials = true
	corsConfig.AllowHeaders = []string{"Authorization", "Content-Type"}
	r.Use(cors.New(corsConfig))

	api := r.Group("api")
	v1 := api.Group("v1")
	{
		v1.POST("/auth/register", userHandler.RegisterUserHandler)
		v1.POST("/auth/login", userHandler.LoginHandler)
		v1.GET("/users/confirm_account", userHandler.ConfirmUserAccountHandler)
		v1.GET("/users/profile", AuthenticateUser(), userHandler.GetProfileHandler)
		v1.PATCH("/users", AuthenticateUser(), userHandler.UpdateProfileHandler)
		v1.POST("/users/avatar", AuthenticateUser(), userHandler.UploadUserAvatarHandler)
		v1.GET("/rooms", AuthenticateUser(), roomHandler.GetRoomsHandler)
		v1.GET("/rooms/:room_id/messages", AuthenticateUser(), messageHandler.GetMessagesHandler)
	}

	r.GET("/rooms/:room_id", AuthenticateWS(), messageHandler.WSHandler)

	port := fmt.Sprintf(":%s", AppConfig.ServerPort)
	if err := r.Run(port); err != nil {
		return err
	}

	return nil
}
