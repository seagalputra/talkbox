package api

import (
	"fmt"
	"log"
	"net/http"

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
	conn, err := wsUpgrader.Upgrade(rw, r, r.Header)
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

	r := gin.Default()
	r.GET("/", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "Hello from Talkbox",
		})
	})
	r.GET("/ws", func(ctx *gin.Context) {
		wsHandler(ctx.Writer, ctx.Request)
	})

	port := fmt.Sprintf(":%s", AppConfig.ServerPort)
	if err := r.Run(port); err != nil {
		return err
	}

	return nil
}
