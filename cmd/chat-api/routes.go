package main

import (
	gsc "golang-social-chat"

	"github.com/go-chi/chi/v5"
)

func DefineChatRoutes(r *chi.Mux) {
	r.Post("/register", gsc.RegisterUser)
	r.Post("/login", gsc.LoginUser)
	r.Post("/logout", gsc.LogoutUser)
	r.Post("/chat/send", gsc.SendMessage)
	r.Get("/chat/history", gsc.GetMessages)
	r.Post("/post", gsc.AddPost)
	r.Get("/feed", gsc.GetPosts)
}
