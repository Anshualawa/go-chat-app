package routes

import (
	"golang-social-chat/controllers"

	"github.com/go-chi/chi/v5"
)

func DefineChatRoutes(r *chi.Mux) {
	r.Post("/register", controllers.RegisterUser)
	r.Post("/login", controllers.LoginUser)
	r.Post("/chat/send", controllers.SendMessage)
	r.Get("/chat/history", controllers.GetMessages)
	r.Post("/post", controllers.AddPost)
	r.Get("/feed", controllers.GetPosts)
}
