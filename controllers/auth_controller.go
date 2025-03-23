package controllers

import (
	"database/sql"
	"encoding/json"
	"golang-social-chat/config"
	"golang-social-chat/models"
	"golang-social-chat/utils"
	"io"
	"net/http"
)

func RegisterUser(w http.ResponseWriter, r *http.Request) {
	req, err := io.ReadAll(r.Body)
	if err != nil {
		_err := utils.NewUserError(http.StatusBadRequest, "unable to read the request body", err)
		utils.SendError(w, r, _err)
		return
	}

	var user models.User
	if err = json.Unmarshal(req, &user); err != nil {
		_err := utils.NewUserError(http.StatusBadRequest, "unable to Unmarshal request body", err)
		utils.SendError(w, r, _err)
		return
	}
	if user.Username == "" || user.Password == "" {
		_err := utils.NewUserError(http.StatusBadRequest, "username and password is mandatory", nil)
		utils.SendError(w, r, _err)
		return
	}
	hashedPassword, _ := utils.HashPassword(user.Password)
	_, err = config.DB.Exec("INSERT INTO users (username, password) VALUES (?,?)", user.Username, hashedPassword)

	if err != nil {
		_err := utils.NewSystemError(http.StatusInternalServerError, "unable to perfomr database operation", err)
		utils.SendError(w, r, _err)
		return
	}

	res := utils.NewResult(http.StatusCreated, "User registered successfully")
	utils.SendSuccess(w, r, res)
}

func LoginUser(w http.ResponseWriter, r *http.Request) {
	req, err := io.ReadAll(r.Body)
	if err != nil {
		_err := utils.NewUserError(http.StatusBadRequest, "unable to read request body", err)
		utils.SendError(w, r, _err)
		return
	}

	var user models.User
	if err = json.Unmarshal(req, &user); err != nil {
		_err := utils.NewUserError(http.StatusBadRequest, "unable to unmarshal request body", err)
		utils.SendError(w, r, _err)
		return
	}

	if user.Username == "" || user.Password == "" {
		_err := utils.NewUserError(http.StatusBadRequest, "invalid input check username and password", nil)
		utils.SendError(w, r, _err)
		return
	}
	var hashedPassword string
	err = config.DB.QueryRow("SELECT password FROM users WHERE username = ?", user.Username).Scan(&hashedPassword)
	if err == sql.ErrNoRows || utils.CheckPassword(hashedPassword, user.Password) != nil {
		_err := utils.NewUserError(http.StatusUnauthorized, "Invalid username or password", err)
		utils.SendError(w, r, _err)
		return
	}

	token, _ := utils.GenerateToken(user.Username)
	res := utils.NewResult(http.StatusOK, "Login success").Add("token", token)
	utils.SendSuccess(w, r, res)
}
