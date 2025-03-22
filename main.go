package main

import (
	"golang-social-chat/config"
	"golang-social-chat/controllers"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func main() {
	config.ConnectDB()
	config.ConnectRedis()

	r := chi.NewRouter()

	r.Post("/register", controllers.RegisterUser)
	r.Post("/login", controllers.LoginUser)
	r.Post("/chat/send", controllers.SendMessage)
	r.Get("/chat/history", controllers.GetMessages)
	r.Post("/post", controllers.AddPost)
	r.Get("/feed", controllers.GetPosts)

	log.Println("Starting server on :8080")
	http.ListenAndServe(":8080", r)
}
