package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	"github.com/seagalputra/talkbox/comment"
	_ "github.com/seagalputra/talkbox/docs"
	httpSwagger "github.com/swaggo/http-swagger"
)

func APIRoutes() *chi.Mux {
	route := chi.NewRouter()

	commentHandler := comment.DefaultHandler()
	route.Post("/comments/{post_id}", commentHandler.Insert)

	return route
}

// @title       Talkbox
// @version     1.0
// @description This is a backend for Talkbox comment system.

// @contact.name  API Support
// @contact.url   http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url  http://www.apache.org/licenses/LICENSE-2.0.html

// @host     localhost:3000
// @BasePath /api

// @securityDefinitions.basic BasicAuth
func Handler() *chi.Mux {
	route := chi.NewRouter()
	route.Use(middleware.Logger)
	route.Use(middleware.Recoverer)

	route.Mount("/api", APIRoutes())

	route.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/docs/", http.StatusPermanentRedirect)
	})
	route.Get("/docs/*", httpSwagger.Handler())

	return route
}

func StartServer() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Unable to load application configuration file: %v", err)
	}

	host := fmt.Sprintf("%s:%s", os.Getenv("HOST"), os.Getenv("PORT"))
	server := &http.Server{
		Addr:    host,
		Handler: Handler(),
	}

	ctx, stopCtx := context.WithCancel(context.Background())
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-sig

		shutdownCtx, _ := context.WithTimeout(ctx, 30*time.Second)
		go func() {
			<-shutdownCtx.Done()
			if shutdownCtx.Err() == context.DeadlineExceeded {
				log.Fatalf("Graceful shutdown time out, forcing exit...")
			}
		}()

		err := server.Shutdown(shutdownCtx)
		if err != nil {
			log.Fatalf("Unable to shutdown gracefully, server forced exit")
		}

		stopCtx()
	}()

	log.Printf("Server running and listening on %s", host)
	err = server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Fatalf("Failed to run server: %v", err)
	}

	<-ctx.Done()
}
