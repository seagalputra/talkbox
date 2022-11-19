package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/seagalputra/talkbox/assets"
)

var wsupgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	Error: func(w http.ResponseWriter, r *http.Request, status int, reason error) {
		log.Printf("connect error: %v", reason)
		http.Error(w, "Could not open websocket connection", http.StatusBadRequest)
	},
	// TODO: should check an origin
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func APIRoutes() {
}

func wshandler(w http.ResponseWriter, r *http.Request) {
	conn, err := wsupgrader.Upgrade(w, r, w.Header())
	if err != nil {
		log.Printf("wshandler: %v", err)
		http.Error(w, "Could not open websocket connection", http.StatusBadRequest)
	}

	for {
		t, _, err := conn.ReadMessage()
		if err != nil {
			log.Printf("wshandler: %v", err)
			break
		}
		conn.WriteMessage(t, []byte("pong"))
	}
}

func main() {
	err := assets.LoadAppConfig()
	if err != nil {
		log.Println("Unable to load env file, the app will be use system config")
	}

	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"*"},
	}))
	r.GET("/", func(ctx *gin.Context) {
		ctx.String(200, "Welcome to Talbox Application")
	})
	r.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	r.GET("/ws", func(ctx *gin.Context) {
		wshandler(ctx.Writer, ctx.Request)
	})

	srv := &http.Server{
		Addr:    ":12312",
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}

	log.Println("Server exiting")
}
