package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/seagalputra/talkbox/comment"
)

func APIRoutes() *chi.Mux {
	route := chi.NewRouter()

	commentHandler := comment.DefaultHandler()
	route.Get("/comments", commentHandler.FindAll)
	route.Post("/comments/{post_id}", commentHandler.Insert)
	route.Patch("/comments/{id}", commentHandler.Update)
	route.Delete("/comments/{id}", commentHandler.Delete)

	return route
}
