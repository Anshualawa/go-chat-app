package controllers

import (
	"database/sql"
	"encoding/json"
	"golang-social-chat/config"
	"golang-social-chat/models"
	"golang-social-chat/utils"
	"net/http"
)

func RegisterUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	json.NewDecoder(r.Body).Decode(&user)

	hashedPassword, _ := utils.HashPassword(user.Password)
	_, err := config.DB.Exec("INSERT INTO users (username, password) VALUES (?,?)", user.Username, hashedPassword)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("User created successfully"))
}

func LoginUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	json.NewDecoder(r.Body).Decode(&user)

	var hashedPassword string
	err := config.DB.QueryRow("SELECT password FROM users WHERE username = ?", user.Username).Scan(&hashedPassword)
	if err == sql.ErrNoRows || utils.CheckPassword(hashedPassword, user.Password) != nil {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	token, _ := utils.GenerateToken(user.Username)
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}
