package golang_social_chat

import (
	"encoding/json"
	"fmt"
	"golang-social-chat/config"
	"io"
	"net/http"
)

func SendMessage(w http.ResponseWriter, r *http.Request) {
	var msg map[string]string
	json.NewDecoder(r.Body).Decode(&msg)

	config.RDB.RPush("chat:room", fmt.Sprintf("%s: %s", msg["user"], msg["message"]))
	res := NewResult(http.StatusCreated, "Message sent")
	SendSuccess(w, r, res)
}

func GetMessages(w http.ResponseWriter, r *http.Request) {
	messages, _ := config.RDB.LRange("chat:room", 0, -1).Result()
	res := NewResult(http.StatusOK, "retrieved chat").Add("chat-list", messages)
	SendSuccess(w, r, res)
}

func SendMessageHandler(w http.ResponseWriter, r *http.Request) {
	_, err := io.ReadAll(r.Body)
	if err != nil {
		_err := NewUserError(http.StatusBadRequest, "unable to read the request body", err)
		SendError(w, r, _err)
		return
	}

}
