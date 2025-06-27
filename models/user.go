package models

// User represents user data from client
type User struct {
	ID       int    `json:"id,omitempty"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password,omitempty"` // omitempty for JSON responses
}

// LoginRequest represents login form data
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Response represents standard JSON response
type Response struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	UserID  int    `json:"userId,omitempty"`
	Name    string `json:"name,omitempty"`
}
