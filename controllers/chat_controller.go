package controllers

import (
	"encoding/json"
	"fmt"
	"golang-social-chat/config"
	"golang-social-chat/utils"
	"net/http"
)

func SendMessage(w http.ResponseWriter, r *http.Request) {
	var msg map[string]string
	json.NewDecoder(r.Body).Decode(&msg)

	config.RDB.RPush("chat:room", fmt.Sprintf("%s: %s", msg["user"], msg["message"]))
	res := utils.NewResult(http.StatusCreated, "Message sent")
	utils.SendSuccess(w, r, res)
}

func GetMessages(w http.ResponseWriter, r *http.Request) {
	messages, _ := config.RDB.LRange("chat:room", 0, -1).Result()
	res := utils.NewResult(http.StatusOK, "retrieved chat").Add("chat-list", messages)
	utils.SendSuccess(w, r, res)
}
