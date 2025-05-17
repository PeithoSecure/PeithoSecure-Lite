package models

// LoginRequest represents the payload for login
type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
	Email    string `json:"email,omitempty"`
}

// LoginResponse represents the token returned after login
type LoginResponse struct {
	AccessToken   string `json:"access_token"`
	RefreshToken  string `json:"refresh_token"`
	TokenType     string `json:"token_type"`
	ExpiresIn     int    `json:"expires_in"`
	EmailVerified bool   `json:"email_verified"` // <-- Added email verification info
}

// User represents a user record in the database
type User struct {
	ID            int    `json:"id"`
	Username      string `json:"username"`
	Email         string `json:"email"`
	Role          string `json:"role"`
	EmailVerified int    `json:"email_verified"`
}
