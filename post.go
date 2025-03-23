package golang_social_chat

import (
	"encoding/json"
	"fmt"
	"golang-social-chat/config"
	"net/http"
)

func AddPost(w http.ResponseWriter, r *http.Request) {
	var post map[string]string
	json.NewDecoder(r.Body).Decode(&post)

	config.RDB.LPush("social:feed", fmt.Sprintf("%s: %s", post["user"], post["message"]))
	w.Write([]byte("Post Added"))
}

func GetPosts(w http.ResponseWriter, r *http.Request) {
	posts, _ := config.RDB.LRange("social:feed", 0, -1).Result()
	json.NewEncoder(w).Encode(posts)
}
