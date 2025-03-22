package controllers

import (
	"encoding/json"
	"fmt"
	"golang-social-chat/config"
	"net/http"
)

func SendMessage(w http.ResponseWriter, r *http.Request) {
	var msg map[string]string
	json.NewDecoder(r.Body).Decode(&msg)

	config.RDB.RPush("chat:room", fmt.Sprintf("%s: %s", msg["user"], msg["message"]))
	w.Write([]byte("Message sent"))
}

func GetMessages(w http.ResponseWriter, r *http.Request) {
	messages, _ := config.RDB.LRange("chat:room", 0, -1).Result()
	json.NewEncoder(w).Encode(messages)
}
