package golang_social_chat

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"golang-social-chat/config"
	"io"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// 🔹 Register User
func RegisterUser(w http.ResponseWriter, r *http.Request) {
	// Read request body
	req, err := io.ReadAll(r.Body)
	if err != nil {
		SendError(w, r, NewUserError(http.StatusBadRequest, "Unable to read request body", err))
		return
	}
	defer r.Body.Close() // Close request body when function ends

	// Parse JSON
	var user User
	if err = json.Unmarshal(req, &user); err != nil {
		SendError(w, r, NewUserError(http.StatusBadRequest, "Invalid JSON format", err))
		return
	}

	// Validate input fields
	if user.Username == "" || user.Email == "" || user.Password == "" {
		SendError(w, r, NewUserError(http.StatusBadRequest, "Username, Email, and Password are required", nil))
		return
	}

	// 🔹 Check if username or email already exists
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

	// 🔹 Hash Password
	hashedPassword, err := HashPassword(user.Password)
	if err != nil {
		SendError(w, r, NewSystemError(http.StatusInternalServerError, "Error hashing password", err))
		return
	}

	// 🔹 Insert User into Database
	res, err := config.DB.Exec("INSERT INTO users (username, email, password_hash) VALUES (?, ?, ?)", user.Username, user.Email, hashedPassword)
	if err != nil {
		SendError(w, r, NewSystemError(http.StatusInternalServerError, "Failed to register user", err))
		return
	}

	// Get newly inserted user ID
	userID, _ := res.LastInsertId()

	// ✅ Successful Registration Response
	response := NewResult(http.StatusCreated, "User registered successfully").
		Add("user_id", userID).
		Add("username", user.Username)
	SendSuccess(w, r, response)
}

// 🔹 Login User
func LoginUser(w http.ResponseWriter, r *http.Request) {
	// Read request body
	req, err := io.ReadAll(r.Body)
	if err != nil {
		SendError(w, r, NewUserError(http.StatusBadRequest, "Unable to read request body", err))
		return
	}
	defer r.Body.Close() // Close request body when function ends

	// Parse JSON
	var user User
	if err = json.Unmarshal(req, &user); err != nil {
		SendError(w, r, NewUserError(http.StatusBadRequest, "Invalid JSON format", err))
		return
	}

	// Validate input fields
	if user.Username == "" || user.Password == "" {
		SendError(w, r, NewUserError(http.StatusBadRequest, "Username and Password are required", nil))
		return
	}

	// 🔹 Fetch user_id and hashed password
	var userID int
	var hashedPassword string
	err = config.DB.QueryRow("SELECT id, password_hash FROM users WHERE username = ?", user.Username).Scan(&userID, &hashedPassword)
	if err == sql.ErrNoRows {
		SendError(w, r, NewUserError(http.StatusUnauthorized, "Invalid username or password", nil))
		return
	} else if err != nil {
		SendError(w, r, NewSystemError(http.StatusInternalServerError, "Database error", err))
		return
	}

	// 🔹 Verify Password
	if err := CheckPassword(hashedPassword, user.Password); err != nil {
		SendError(w, r, NewUserError(http.StatusUnauthorized, "Invalid username or password", nil))
		return
	}

	// 🔹 Generate JWT Token with user_id
	token, err := GenerateToken(userID, user.Username)
	if err != nil {
		SendError(w, r, NewSystemError(http.StatusInternalServerError, "Error generating token", err))
		return
	}

	// 🔹 Set User as Online in `user_status` Table
	_, err = config.DB.Exec(`
		INSERT INTO user_status (user_id, is_online, last_seen) 
		VALUES (?, ?, NOW()) 
		ON DUPLICATE KEY UPDATE is_online = VALUES(is_online), last_seen = NOW()
	`, userID, true)
	if err != nil {
		SendError(w, r, NewSystemError(http.StatusInternalServerError, "Error updating user status", err))
		return
	}

	// ✅ Send Success Response
	response := NewResult(http.StatusOK, "Login successful").
		Add("token", token).
		Add("user_id", userID).
		Add("username", user.Username)
	SendSuccess(w, r, response)
}

// LogoutUser sets the user status to offline
func LogoutUser(w http.ResponseWriter, r *http.Request) {
	// Extract JWT Token from Authorization Header
	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		SendError(w, r, NewUserError(http.StatusUnauthorized, "Missing Authorization Token", nil))
		return
	}

	// Validate and Parse Token
	claims, err := ParseToken(tokenString)
	if err != nil {
		SendError(w, r, NewUserError(http.StatusUnauthorized, "Invalid or expired token", err))
		return
	}

	// Extract user_id from claims
	userID, ok := claims["user_id"].(float64)
	if !ok {
		SendError(w, r, NewUserError(http.StatusUnauthorized, "Invalid token payload", nil))
		return
	}

	// 🔹 Set user status to offline
	_, err = config.DB.Exec(`
		UPDATE user_status SET is_online = ?, last_seen = NOW() WHERE user_id = ?
	`, false, int(userID))
	if err != nil {
		SendError(w, r, NewSystemError(http.StatusInternalServerError, "Error updating user status", err))
		return
	}

	// ✅ Send Success Response
	SendSuccess(w, r, NewResult(http.StatusOK, "Logout successful"))
}

// ParseToken verifies and extracts claims from the JWT token
func ParseToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	// Extract claims
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, fmt.Errorf("invalid token")
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

// 🔹 Generate JWT Token (Includes `user_id`)
func GenerateToken(userID int, username string) (string, error) {
	claims := jwt.MapClaims{
		"user_id":  userID,
		"username": username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(), // Expires in 24 hours
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}
