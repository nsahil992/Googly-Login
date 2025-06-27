package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"calico-go-project/database"
	"calico-go-project/models"
	"calico-go-project/utils"
)

// RegisterHandler handles user registration
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse request body
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		sendJSONResponse(w, false, "Invalid input", http.StatusBadRequest)
		return
	}

	// Check if user already exists
	exists, err := database.UserExists(user.Email)
	if err != nil {
		sendJSONResponse(w, false, "Server error", http.StatusInternalServerError)
		return
	}

	if exists {
		sendJSONResponse(w, false, "User with this email already exists", http.StatusConflict)
		return
	}

	// Hash the password
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		sendJSONResponse(w, false, "Error processing password", http.StatusInternalServerError)
		return
	}

	// Create the user
	err = database.CreateUser(user.Name, user.Email, hashedPassword)
	if err != nil {
		sendJSONResponse(w, false, "Error creating user", http.StatusInternalServerError)
		return
	}

	sendJSONResponse(w, true, "User registered successfully", http.StatusCreated)
}

// LoginHandler handles user login
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse request body
	var loginReq models.LoginRequest
	err := json.NewDecoder(r.Body).Decode(&loginReq)
	if err != nil {
		sendJSONResponse(w, false, "Invalid input", http.StatusBadRequest)
		return
	}

	// For debugging purposes
	fmt.Println("Login attempt:", loginReq.Email)

	// Get user from database
	id, name, hashedPassword, err := database.GetUserByEmail(loginReq.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			sendJSONResponse(w, false, "User not found", http.StatusUnauthorized)
		} else {
			fmt.Println("Database error:", err)
			sendJSONResponse(w, false, "Server error", http.StatusInternalServerError)
		}
		return
	}

	// Check password
	passwordMatch := utils.CheckPasswordHash(loginReq.Password, hashedPassword)
	fmt.Println("Password match:", passwordMatch)

	if !passwordMatch {
		sendJSONResponse(w, false, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// If we get here, login is successful
	response := models.Response{
		Success: true,
		Message: "Login successful",
		UserID:  id,
		Name:    name,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// Helper function to send JSON responses
func sendJSONResponse(w http.ResponseWriter, success bool, message string, statusCode int) {
	response := models.Response{
		Success: success,
		Message: message,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}
