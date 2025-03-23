package golang_social_chat

import (
	"database/sql"
	"encoding/json"
	"golang-social-chat/config"
	"io"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// 🔹 Register User
func RegisterUser(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close() // Close request body when function ends

	req, err := io.ReadAll(r.Body)
	if err != nil {
		SendError(w, r, NewUserError(http.StatusBadRequest, "Unable to read request body", err))
		return
	}

	var user User
	if err = json.Unmarshal(req, &user); err != nil {
		SendError(w, r, NewUserError(http.StatusBadRequest, "Invalid JSON format", err))
		return
	}

	if user.Username == "" || user.Email == "" || user.Password == "" {
		SendError(w, r, NewUserError(http.StatusBadRequest, "Username, Email, and Password are required", nil))
		return
	}

	// 🔹 Check if username or email is already registered
	var count int
	err = config.DB.QueryRow("SELECT COUNT(*) FROM users WHERE username = ? OR email = ?", user.Username, user.Email).Scan(&count)
	if err != nil {
		SendError(w, r, NewSystemError(http.StatusInternalServerError, "Database error", err))
		return
	}
	if count > 0 {
		SendError(w, r, NewUserError(http.StatusBadRequest, "Username or Email is already registered", nil))
		return
	}

	// 🔹 Hash password before storing
	hashedPassword, err := HashPassword(user.Password)
	if err != nil {
		SendError(w, r, NewSystemError(http.StatusInternalServerError, "Error hashing password", err))
		return
	}

	// 🔹 Insert new user into database
	_, err = config.DB.Exec("INSERT INTO users (username, email, password_hash) VALUES (?, ?, ?)", user.Username, user.Email, hashedPassword)
	if err != nil {
		SendError(w, r, NewSystemError(http.StatusInternalServerError, "Failed to register user", err))
		return
	}

	// ✅ Registration Success Response (Do NOT return the password!)
	res := NewResult(http.StatusCreated, "User registered successfully").Add("username", user.Username)
	SendSuccess(w, r, res)
}

// 🔹 Login User
func LoginUser(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close() // Close request body when function ends

	req, err := io.ReadAll(r.Body)
	if err != nil {
		SendError(w, r, NewUserError(http.StatusBadRequest, "Unable to read request body", err))
		return
	}

	var user User
	if err = json.Unmarshal(req, &user); err != nil {
		SendError(w, r, NewUserError(http.StatusBadRequest, "Invalid JSON format", err))
		return
	}

	if user.Username == "" || user.Password == "" {
		SendError(w, r, NewUserError(http.StatusBadRequest, "Username and Password are required", nil))
		return
	}

	// 🔹 Fetch hashed password from database
	var hashedPassword string
	err = config.DB.QueryRow("SELECT password_hash FROM users WHERE username = ?", user.Username).Scan(&hashedPassword)
	if err == sql.ErrNoRows {
		SendError(w, r, NewUserError(http.StatusUnauthorized, "Invalid username or password", nil))
		return
	} else if err != nil {
		SendError(w, r, NewSystemError(http.StatusInternalServerError, "Database error", err))
		return
	}

	// 🔹 Check if password matches
	if err := CheckPassword(hashedPassword, user.Password); err != nil {
		SendError(w, r, NewUserError(http.StatusUnauthorized, "Invalid username or password", nil))
		return
	}

	// 🔹 Generate JWT Token
	token, err := GenerateToken(user.Username)
	if err != nil {
		SendError(w, r, NewSystemError(http.StatusInternalServerError, "Error generating token", err))
		return
	}

	// ✅ Login Success Response
	res := NewResult(http.StatusOK, "Login successful").Add("token", token)
	SendSuccess(w, r, res)
}

// 🔹 Hash Password
func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashedPassword), err
}

// 🔹 Check Password
func CheckPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

// 🔹 JWT Secret Key
var secretKey = []byte("my_secret_key")

// 🔹 Generate JWT Token
func GenerateToken(username string) (string, error) {
	claims := jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(), // Expires in 24 hours
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}
