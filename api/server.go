package api

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var wsUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// TODO: should check the origin request from client
		return true
	},
}

func wsHandler(rw http.ResponseWriter, r *http.Request) {
	conn, err := wsUpgrader.Upgrade(rw, r, nil)
	if err != nil {
		log.Printf("wsHandler: %v", err)
		return
	}

	for {
		t, msg, err := conn.ReadMessage()
		if err != nil {
			break
		}
		conn.WriteMessage(t, msg)
	}
}

func StartServer() error {
	if err := LoadAppConfig(); err != nil {
		log.Printf("Unable to load config file: %v", err)
	}

	if err := ConnectDatabase(); err != nil {
		log.Fatalf("StartServer: %v", err)
	}
	ConnectToRedis()

	userHandler := UserDefaultHandler()
	r := gin.Default()

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"http://127.0.0.1:3000", "http://localhost:3000"}
	corsConfig.AllowCredentials = true
	r.Use(cors.New(corsConfig))
	r.Use(ParseAuthCookies())

	api := r.Group("api")
	v1 := api.Group("v1")
	{
		v1.POST("/auth/register", userHandler.RegisterUserHandler)
		v1.POST("/auth/login", userHandler.LoginHandler)
		v1.GET("/users/confirm_account", userHandler.ConfirmUserAccountHandler)
		v1.GET("/ping", AuthenticateUser(), func(ctx *gin.Context) {
			userCtx, _ := ctx.Get("user")
			user := userCtx.(*User)
			log.Printf("user info: %v", user)

			ctx.JSON(200, gin.H{
				"message": "pong",
			})
			return
		})
	}

	r.GET("/ws", func(ctx *gin.Context) {
		wsHandler(ctx.Writer, ctx.Request)
	})

	port := fmt.Sprintf(":%s", AppConfig.ServerPort)
	if err := r.Run(port); err != nil {
		return err
	}

	return nil
}
